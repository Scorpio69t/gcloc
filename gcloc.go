package gcloc

import (
	"github.com/Scorpio69t/gcloc/pkg/file"
	"github.com/Scorpio69t/gcloc/pkg/language"
	"github.com/Scorpio69t/gcloc/pkg/option"
	log "github.com/Scorpio69t/gcloc/pkg/simplelog"
	"github.com/Scorpio69t/gcloc/pkg/syncmap"
	"sync"
)

// Parser is the main struct for parsing files.
type Parser struct {
	languages *language.DefinedLanguages
	opts      *option.GClocOptions
}

// Result is the main struct for the result of parsing files.
type Result struct {
	Total         *language.Language
	Files         map[string]*file.GClocFile
	Languages     map[string]*language.Language
	MaxPathLength int
}

// NewParser creates a new Parser.
func NewParser(languages *language.DefinedLanguages, opts *option.GClocOptions) *Parser {
	return &Parser{
		languages: languages,
		opts:      opts,
	}
}

// Analyze analyzes the files in the given paths.
func (p *Parser) Analyze(paths []string) (*Result, error) {
	total := language.NewLanguage("Total", []string{}, [][]string{{"", ""}}) // Create a new language for the total.
	languages, err := language.GetAllFiles(paths, p.languages, p.opts)
	if err != nil {
		log.Error("Error getting all files: %v", err)
		return nil, err
	}

	maxPathLen := 0
	num := 0
	for _, lang := range languages {
		num += len(lang.Files)
		for _, f := range lang.Files {
			l := len(f)
			if maxPathLen < l {
				maxPathLen = l
			}
		}
	}

	gClocFiles := syncmap.NewSyncMap[string, *file.GClocFile](num)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, lang := range languages {
		wg.Add(1)

		go func(l *language.Language) {
			defer wg.Done()

			for _, filename := range l.Files {
				cf := file.AnalyzeFile(filename, l, p.opts)
				cf.Language = l.Name

				l.Codes += cf.Codes
				l.Comments += cf.Comments
				l.Blanks += cf.Blanks
				gClocFiles.Store(filename, cf)
			}

			mu.Lock()
			defer mu.Unlock()

			files := uint32(len(l.Files))
			if files > 0 {
				total.Total += files
				total.Blanks += l.Blanks
				total.Comments += l.Comments
				total.Codes += l.Codes
			}
		}(lang)
	}

	wg.Wait()

	return &Result{
		Total:         total,
		Files:         gClocFiles.ToMap(), // Convert syncmap to map.
		Languages:     languages,
		MaxPathLength: maxPathLen,
	}, nil
}
