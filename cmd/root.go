/*
Copyright © 2024 Yang Ruitao yangruitao6@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Scorpio69t/gcloc"
	"github.com/Scorpio69t/gcloc/pkg/file"
	gjson "github.com/Scorpio69t/gcloc/pkg/json"
	"github.com/Scorpio69t/gcloc/pkg/language"
	"github.com/Scorpio69t/gcloc/pkg/option"
	log "github.com/Scorpio69t/gcloc/pkg/simplelog"
	"github.com/Scorpio69t/gcloc/pkg/utils"
	"github.com/Scorpio69t/gcloc/pkg/xml"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type options struct {
	ByFile         bool
	SortTag        string
	OutputType     string
	ExcludeExt     string
	ExcludeLang    string
	IncludeLang    string
	Match          string
	NotMatch       string
	MatchDir       string
	NotMatchDir    string
	Debug          bool
	SkipDuplicated bool
}

type outputBuilder struct {
	opts   *options
	result *gcloc.Result
}

// opts is a global variable to store command options
var opts options

const (
	OutputTypeDefault   string = "default"
	OutputTypeClocXML   string = "gcloc-xml"
	OutputTypeSloccount string = "sloccount"
	OutputTypeJSON      string = "json"

	fileHeader             string = "File"
	languageHeader         string = "Language"
	commonHeader           string = "files          blank        comment           code"
	defaultOutputSeparator string = "----------------------------------------------------------------------------" +
		"----------------------------------------------------------------------------" +
		"----------------------------------------------------------------------------"
)

// rowLen is the default row length
var rowLen = 79

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gcloc [flags] PATH...",
	Short: "A tool for counting source code files and lines in various programming languages",
	Long: `gcloc is a tool for counting source code files and lines in various programming languages.
It supports a variety of programming languages, and can be customized to support more languages.
It is a simple and easy-to-use tool that can help you quickly count the number of source code files and lines in a project.`,

	Args: cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runGCloc,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gcloc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// Define flags
	rootCmd.Flags().BoolVar(&opts.ByFile, "by-file", false, "report results for every encountered source file")
	rootCmd.Flags().StringVar(&opts.SortTag, "sort", "codes", "sort based on a certain column [name, files, blanks, comments, codes]")
	rootCmd.Flags().StringVar(&opts.OutputType, "output-type", "default", "output type [default, gcloc-xml, sloccount, json]")
	rootCmd.Flags().StringVar(&opts.ExcludeExt, "exclude-ext", "", "exclude file name extensions (comma-separated)")
	rootCmd.Flags().StringVar(&opts.ExcludeLang, "exclude-lang", "", "exclude language names (comma-separated)")
	rootCmd.Flags().StringVar(&opts.IncludeLang, "include-lang", "", "include language names (comma-separated)")
	rootCmd.Flags().StringVar(&opts.Match, "match", "", "include file name (regex)")
	rootCmd.Flags().StringVar(&opts.NotMatch, "not-match", "", "exclude file name (regex)")
	rootCmd.Flags().StringVar(&opts.MatchDir, "match-d", "", "include dir name (regex)")
	rootCmd.Flags().StringVar(&opts.NotMatchDir, "not-match-d", "", "exclude dir name (regex)")
	rootCmd.Flags().BoolVar(&opts.Debug, "debug", false, "dump debug log for developer")
	rootCmd.Flags().BoolVar(&opts.SkipDuplicated, "skip-duplicated", false, "skip duplicated files")
}

func runGCloc(cmd *cobra.Command, args []string) {
	if opts.ByFile && opts.SortTag == "files" {
		cmd.Println("`--sort files` option cannot be used in conjunction with the `--by-file` option")
		os.Exit(1)
	}

	var paths []string
	var repoPaths []string

	for _, arg := range args {
		if utils.IsGitURL(arg) {
			tempPath, err := utils.CloneGitRepo(arg)
			if err != nil {
				fmt.Printf("fail to clone git repo[%s]. error: %v\n", arg, err)
				return
			}
			paths = append(paths, tempPath)
			repoPaths = append(repoPaths, tempPath)
		} else {
			paths = append(paths, arg)
		}
	}

	defer func() {
		for _, path := range repoPaths {
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("fail to remove git repo[%s]. error: %v\n", path, err)
			}
		}
	}()

	languages := language.NewDefinedLanguages()
	gClocOpts := option.NewGClocOptions()

	// Setup options
	setupOptions(gClocOpts, languages)

	start := time.Now() // get current time

	parser := gcloc.NewParser(languages, gClocOpts)
	result, err := parser.Analyze(paths)
	elapsed := time.Since(start) // Count time
	if err != nil {
		fmt.Printf("fail gcloc analyze. error: %v\n", err)
		return
	}

	// Output result
	builder := newOutputBuilder(result, &opts)
	builder.WriteResult(elapsed)
}

// setupOptions setup options for cloc
func setupOptions(clocOpts *option.GClocOptions, languages *language.DefinedLanguages) {
	// Exclude extensions
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		clocOpts.ExcludeExts[ext] = struct{}{}
	}

	// Exclude languages
	for _, lang := range strings.Split(opts.ExcludeLang, ",") {
		if _, ok := languages.Langs[lang]; ok {
			clocOpts.ExcludeLanguages[lang] = struct{}{}
		}
	}

	// Include languages
	for _, lang := range strings.Split(opts.IncludeLang, ",") {
		if _, ok := languages.Langs[lang]; ok {
			clocOpts.IncludeLanguages[lang] = struct{}{}
		}
	}

	// Setup regex matches
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

	clocOpts.Debug = opts.Debug
	clocOpts.SkipDuplicated = opts.SkipDuplicated
}

func newOutputBuilder(result *gcloc.Result, opts *options) *outputBuilder {
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

func writeResultWithByFile(opts *options, result *gcloc.Result) {
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

func (o *outputBuilder) WriteResult(elapsed time.Duration) {
	clocLanguages := o.result.Languages
	total := o.result.Total

	seconds := elapsed.Seconds()
	// calculate files per second and lines per second
	filesPerSecond := float64(total.Total) / seconds
	linesPerSecond := float64(total.Codes) / seconds

	// write time elapsed with seconds
	fmt.Printf("github.com/Scorpio69t/gcloc T=%0.2f s (%0.1f files/s %0.1f lines/s)\n", seconds, filesPerSecond, linesPerSecond)

	// write header
	o.WriteHeader()

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
