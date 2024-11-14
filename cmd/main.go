package main

import (
	"encoding/json"
	"fmt"
	"gcloc/pkg/file"
	"gcloc/pkg/language"
	"gcloc/pkg/option"
	log "gcloc/pkg/simplelog"
	"gcloc/pkg/xml"
	"os"
	"regexp"
	"strings"

	gjson "gcloc/pkg/json"

	"gcloc"
	"github.com/jessevdk/go-flags"
)

// Version is version string for gcloc command
var Version string

// GitCommit is git commit hash string for gcloc command
var GitCommit string

// OutputTypeDefault is cloc's text output format for --output-type option
const OutputTypeDefault string = "default"

// OutputTypeClocXML is Cloc's XML output format for --output-type option
const OutputTypeClocXML string = "gcloc-xml"

// OutputTypeSloccount is Sloccount output format for --output-type option
const OutputTypeSloccount string = "sloccount"

// OutputTypeJSON is JSON output format for --output-type option
const OutputTypeJSON string = "json"

const fileHeader string = "File"
const languageHeader string = "Language"
const commonHeader string = "files          blank        comment           code"
const defaultOutputSeparator string = "----------------------------------------------------------------------------" +
	"----------------------------------------------------------------------------" +
	"----------------------------------------------------------------------------"

var rowLen = 79

// CmdOptions is gcloc command options.
// It is necessary to use notation that follows go-flags.
type CmdOptions struct {
	ByFile         bool   `long:"by-file" description:"report results for every encountered source file"`
	SortTag        string `long:"sort" default:"codes" description:"sort based on a certain column" choice:"name" choice:"files" choice:"blanks" choice:"comments" choice:"codes"`
	OutputType     string `long:"output-type" default:"default" description:"output type [values: default,gcloc-xml,sloccount,json]"`
	ExcludeExt     string `long:"exclude-ext" description:"exclude file name extensions (separated commas)"`
	IncludeLang    string `long:"include-lang" description:"include language name (separated commas)"`
	Match          string `long:"match" description:"include file name (regex)"`
	NotMatch       string `long:"not-match" description:"exclude file name (regex)"`
	MatchDir       string `long:"match-d" description:"include dir name (regex)"`
	NotMatchDir    string `long:"not-match-d" description:"exclude dir name (regex)"`
	Debug          bool   `long:"debug" description:"dump debug log for developer"`
	SkipDuplicated bool   `long:"skip-duplicated" description:"skip duplicated files"`
	ShowLang       bool   `long:"show-lang" description:"print about all languages and extensions"`
	ShowVersion    bool   `long:"version" description:"print version info"`
}

type outputBuilder struct {
	opts   *CmdOptions
	result *gcloc.Result
}

func newOutputBuilder(result *gcloc.Result, opts *CmdOptions) *outputBuilder {
	return &outputBuilder{
		opts,
		result,
	}
}

func (o *outputBuilder) WriteHeader() {
	maxPathLen := o.result.MaxPathLength
	headerLen := 28
	header := languageHeader

	if o.opts.ByFile {
		headerLen = maxPathLen + 1
		rowLen = maxPathLen + len(commonHeader) + 2
		header = fileHeader
	}
	if o.opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
		fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, commonHeader)
		fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
	}
}

func (o *outputBuilder) WriteFooter() {
	total := o.result.Total
	maxPathLen := o.result.MaxPathLength

	if o.opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
		if o.opts.ByFile {
			fmt.Printf("%-[1]*[2]v %6[3]v %14[4]v %14[5]v %14[6]v\n",
				maxPathLen, "Total", total.Total, total.Blanks, total.Comments, total.Codes)
		} else {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				"Total", total.Total, total.Blanks, total.Comments, total.Codes)
		}
		fmt.Printf("%.[2]*[1]s\n", defaultOutputSeparator, rowLen)
	}
}

func writeResultWithByFile(opts *CmdOptions, result *gcloc.Result) {
	gClocFiles := result.Files
	total := result.Total
	maxPathLen := result.MaxPathLength

	var sortedFiles file.GClocFiles
	for _, cFile := range gClocFiles {
		sortedFiles = append(sortedFiles, *cFile)
	}

	if sortedFiles != nil {
		switch opts.SortTag {
		case "name":
			sortedFiles.SortByName()
		case "comments":
			sortedFiles.SortByComments()
		case "blanks":
			sortedFiles.SortByBlanks()
		default:
			sortedFiles.SortByCodes()
		}
	}

	switch opts.OutputType {
	case OutputTypeClocXML:
		t := xml.XMLTotalFiles{
			Codes:    total.Codes,
			Comments: total.Comments,
			Blanks:   total.Blanks,
		}

		f := &xml.XMLResultFiles{
			Files: sortedFiles,
			Total: t,
		}
		xmlResult := xml.XMLResult{
			XMLFiles: f,
		}
		xmlResult.Encode()
	case OutputTypeSloccount:
		for _, sf := range sortedFiles {
			p := ""
			if strings.HasPrefix(sf.Name, "./") || string(sf.Name[0]) == "/" {
				splitPaths := strings.Split(sf.Name, string(os.PathSeparator))
				if len(splitPaths) >= 3 {
					p = splitPaths[1]
				}
			}
			fmt.Printf("%v\t%v\t%v\t%v\n",
				sf.Codes, sf.Language, p, sf.Name)
		}
	case OutputTypeJSON:
		jsonResult := gjson.NewFilesResultFromGCloc(total, sortedFiles)
		buf, err := json.Marshal(jsonResult)
		if err != nil {
			fmt.Println(err)
			panic("json marshal error")
		}

		_, err = os.Stdout.Write(buf)
		if err != nil {
			log.Error("write json result error: %v", err)
			return
		}
	default:
		for _, sf := range sortedFiles {
			clocFile := sf
			fmt.Printf("%-[1]*[2]s %21[3]v %14[4]v %14[5]v\n",
				maxPathLen, sf.Name, clocFile.Blanks, clocFile.Comments, clocFile.Codes)
		}
	}
}

func (o *outputBuilder) WriteResult() {
	// write header
	o.WriteHeader()

	clocLanguages := o.result.Languages
	total := o.result.Total

	if o.opts.ByFile {
		writeResultWithByFile(o.opts, o.result)
	} else {
		var sortedLanguages language.Languages
		for _, lang := range clocLanguages {
			if len(lang.Files) != 0 {
				sortedLanguages = append(sortedLanguages, *lang)
			}
		}

		if sortedLanguages != nil {
			switch o.opts.SortTag {
			case "name":
				sortedLanguages.SortByName()
			case "files":
				sortedLanguages.SortByFiles()
			case "comment":
				sortedLanguages.SortByComments()
			case "blank":
				sortedLanguages.SortByBlanks()
			default:
				sortedLanguages.SortByCodes()
			}
		}

		switch o.opts.OutputType {
		case OutputTypeClocXML:
			xmlResult := xml.NewXMLResultFromCloc(total, sortedLanguages, xml.XMLResultWithLangs)
			xmlResult.Encode()
		case OutputTypeJSON:
			jsonResult := gjson.NewLanguagesResultFromGCloc(total, sortedLanguages)
			buf, err := json.Marshal(jsonResult)
			if err != nil {
				fmt.Println(err)
				panic("json marshal error")
			}

			_, err = os.Stdout.Write(buf)
			if err != nil {
				log.Error("write json result error: %v", err)
				return
			}
		default:
			for _, lang := range sortedLanguages {
				fmt.Printf("%-27v %6v %14v %14v %14v\n",
					lang.Name, len(lang.Files), lang.Blanks, lang.Comments, lang.Codes)
			}
		}
	}

	// write footer
	o.WriteFooter()
}

func main() {
	var opts CmdOptions
	clocOpts := option.NewGClocOptions()
	// parse command line options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gcloc"
	parser.Usage = "[OPTIONS] PATH[...]"

	paths, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	// value for language result
	languages := language.NewDefinedLanguages()

	if opts.ShowVersion {
		fmt.Printf("%s (%s)\n", Version, GitCommit)
		return
	}

	if opts.ShowLang {
		fmt.Println(languages.GetFormattedString())
		return
	}

	if len(paths) <= 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	// check sort tag option with other options
	if opts.ByFile && opts.SortTag == "files" {
		fmt.Println("`--sort files` option cannot be used in conjunction with the `--by-file` option")
		os.Exit(1)
	}

	// setup option for exclude extensions
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		e, ok := language.FileExtensions[ext]
		if ok {
			clocOpts.ExcludeExts[e] = struct{}{}
		} else {
			clocOpts.ExcludeExts[ext] = struct{}{}
		}
	}

	// directory and file matching options
	if opts.Match != "" {
		clocOpts.ReMatch = regexp.MustCompile(opts.Match)
	}
	if opts.NotMatch != "" {
		clocOpts.ReNotMatch = regexp.MustCompile(opts.NotMatch)
	}
	if opts.MatchDir != "" {
		clocOpts.ReMatchDir = regexp.MustCompile(opts.MatchDir)
	}
	if opts.NotMatchDir != "" {
		clocOpts.ReNotMatchDir = regexp.MustCompile(opts.NotMatchDir)
	}

	// setup option for include languages
	for _, lang := range strings.Split(opts.IncludeLang, ",") {
		if _, ok := languages.Langs[lang]; ok {
			clocOpts.IncludeLanguages[lang] = struct{}{}
		}
	}

	clocOpts.Debug = opts.Debug
	clocOpts.SkipDuplicated = opts.SkipDuplicated

	processor := gcloc.NewParser(languages, clocOpts)
	result, err := processor.Analyze(paths)
	if err != nil {
		fmt.Printf("fail gcloc analyze. error: %v\n", err)
		return
	}

	builder := newOutputBuilder(result, &opts)
	builder.WriteResult()
}
