package gcloc

import (
	"github.com/Scorpio69t/gcloc/pkg/file"
	"github.com/Scorpio69t/gcloc/pkg/language"
	"github.com/Scorpio69t/gcloc/pkg/option"
	log "github.com/Scorpio69t/gcloc/pkg/simplelog"
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

	gClocFiles := make(map[string]*file.GClocFile, num)
	for _, lang := range languages {
		for _, f := range lang.Files {
			cf := file.AnalyzeFile(f, lang, p.opts)
			cf.Language = lang.Name

			lang.Codes += cf.Codes
			lang.Comments += cf.Comments
			lang.Blanks += cf.Blanks
			gClocFiles[f] = cf
		}

		files := uint32(len(lang.Files))
		if len(lang.Files) <= 0 {
			continue
		}

		total.Total += files
		total.Blanks += lang.Blanks
		total.Comments += lang.Comments
		total.Codes += lang.Codes
	}

	return &Result{
		Total:         total,
		Files:         gClocFiles,
		Languages:     languages,
		MaxPathLength: maxPathLen,
	}, nil
}
