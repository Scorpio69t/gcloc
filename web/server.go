package web

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kataras/iris/v12"

	gcloc "github.com/Scorpio69t/gcloc"
	"github.com/Scorpio69t/gcloc/pkg/file"
	"github.com/Scorpio69t/gcloc/pkg/json"
	"github.com/Scorpio69t/gcloc/pkg/language"
	"github.com/Scorpio69t/gcloc/pkg/option"
	"github.com/Scorpio69t/gcloc/pkg/utils"
)

// AnalyzeRequest defines the parameters accepted by the analyze endpoint.
type AnalyzeRequest struct {
	Paths          []string `json:"paths"`
	ByFile         bool     `json:"byFile"`
	Sort           string   `json:"sort"`
	ExcludeExt     string   `json:"excludeExt"`
	ExcludeLang    string   `json:"excludeLang"`
	IncludeLang    string   `json:"includeLang"`
	Match          string   `json:"match"`
	NotMatch       string   `json:"notMatch"`
	MatchDir       string   `json:"matchDir"`
	NotMatchDir    string   `json:"notMatchDir"`
	Debug          bool     `json:"debug"`
	SkipDuplicated bool     `json:"skipDuplicated"`
}

type analyzeResponse struct {
	json.LanguagesResult `json:",omitempty"`
	json.FilesResult     `json:",omitempty"`
	TimeUsed             string `json:"time_used"`
}

// Start runs the iris server on the given address.
func Start(addr string) error {
	app := iris.New()

	app.Get("/languages", languagesHandler)
	app.Post("/analyze", analyzeHandler)

	return app.Listen(addr)
}

func languagesHandler(ctx iris.Context) {
	langs := language.NewDefinedLanguages()
	result := make([]string, 0, len(langs.Langs))
	for k := range langs.Langs {
		result = append(result, k)
	}
	ctx.JSON(result)
}

func analyzeHandler(ctx iris.Context) {
	var req AnalyzeRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}

	var paths []string
	var repoPaths []string
	for _, p := range req.Paths {
		if utils.IsGitURL(p) {
			temp, err := utils.CloneGitRepo(p)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				_, _ = ctx.WriteString(fmt.Sprintf("clone repo %s failed: %v", p, err))
				return
			}
			paths = append(paths, temp)
			repoPaths = append(repoPaths, temp)
		} else {
			paths = append(paths, p)
		}
	}
	defer func() {
		for _, p := range repoPaths {
			_ = os.RemoveAll(p)
		}
	}()

	langs := language.NewDefinedLanguages()
	opts := option.NewGClocOptions()
	setupOptionsFromRequest(opts, langs, &req)

	parser := gcloc.NewParser(langs, opts)
	start := time.Now()
	result, err := parser.Analyze(paths)
	elapsed := time.Since(start)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}

	if req.ByFile {
		var files file.GClocFiles
		for _, f := range result.Files {
			files = append(files, *f)
		}
		resp := analyzeResponse{
			FilesResult: json.NewFilesResultFromGCloc(result.Total, files),
			TimeUsed:    fmt.Sprintf("%0.2fs", elapsed.Seconds()),
		}
		ctx.JSON(resp)
		return
	}

	var langsRes language.Languages
	for _, l := range result.Languages {
		if len(l.Files) != 0 {
			langsRes = append(langsRes, *l)
		}
	}

	resp := analyzeResponse{
		LanguagesResult: json.NewLanguagesResultFromGCloc(result.Total, langsRes),
		TimeUsed:        fmt.Sprintf("%0.2fs", elapsed.Seconds()),
	}
	ctx.JSON(resp)
}

func setupOptionsFromRequest(clocOpts *option.GClocOptions, langs *language.DefinedLanguages, req *AnalyzeRequest) {
	for _, ext := range strings.Split(req.ExcludeExt, ",") {
		if ext != "" {
			clocOpts.ExcludeExts[ext] = struct{}{}
		}
	}
	for _, l := range strings.Split(req.ExcludeLang, ",") {
		if _, ok := langs.Langs[l]; ok {
			clocOpts.ExcludeLanguages[l] = struct{}{}
		}
	}
	for _, l := range strings.Split(req.IncludeLang, ",") {
		if _, ok := langs.Langs[l]; ok {
			clocOpts.IncludeLanguages[l] = struct{}{}
		}
	}
	if req.Match != "" {
		clocOpts.ReMatch = regexp.MustCompile(req.Match)
	}
	if req.NotMatch != "" {
		clocOpts.ReNotMatch = regexp.MustCompile(req.NotMatch)
	}
	if req.MatchDir != "" {
		clocOpts.ReMatchDir = regexp.MustCompile(req.MatchDir)
	}
	if req.NotMatchDir != "" {
		clocOpts.ReNotMatchDir = regexp.MustCompile(req.NotMatchDir)
	}

	clocOpts.Debug = req.Debug
	clocOpts.SkipDuplicated = req.SkipDuplicated
}
