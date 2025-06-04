package web

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	gcloc "github.com/Scorpio69t/gcloc"
	"github.com/Scorpio69t/gcloc/pkg/file"
	"github.com/Scorpio69t/gcloc/pkg/json"
	"github.com/Scorpio69t/gcloc/pkg/language"
	"github.com/Scorpio69t/gcloc/pkg/option"
	"github.com/Scorpio69t/gcloc/pkg/utils"
)

var uploadDirs sync.Map // map[string]string

// AnalyzeRequest defines the parameters accepted by the analyze endpoint.
type AnalyzeRequest struct {
	Paths          []string `json:"paths"`
	UploadID       string   `json:"uploadId"`
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

	app.HandleDir("/", iris.Dir("./web/ui"))

	app.Get("/languages", languagesHandler)
	app.Post("/analyze", analyzeHandler)
	app.Post("/upload", uploadHandler)
	app.Get("/tree", treeHandler)

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
	if req.UploadID != "" {
		if p, ok := uploadDirs.Load(req.UploadID); ok {
			paths = append(paths, p.(string))
		} else {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString("invalid uploadId")
			return
		}
	}
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

type UploadResponse struct {
	ID string `json:"id"`
}

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"is_dir"`
	Children []FileNode `json:"children,omitempty"`
}

func uploadHandler(ctx iris.Context) {
	file, _, err := ctx.FormFile("file")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	defer file.Close()

	tmpZip, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	if _, err := io.Copy(tmpZip, file); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_ = os.Remove(tmpZip.Name())
		_, _ = ctx.WriteString(err.Error())
		return
	}
	tmpZip.Close()

	dest, err := os.MkdirTemp("", "upload-dir-*")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		_ = os.Remove(tmpZip.Name())
		return
	}

	if err := extractZip(tmpZip.Name(), dest); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		_ = os.RemoveAll(dest)
		_ = os.Remove(tmpZip.Name())
		return
	}
	_ = os.Remove(tmpZip.Name())

	id := uuid.New().String()
	uploadDirs.Store(id, dest)
	ctx.JSON(UploadResponse{ID: id})
}

func treeHandler(ctx iris.Context) {
	id := ctx.URLParam("id")
	if id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("id required")
		return
	}
	p, ok := uploadDirs.Load(id)
	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		_, _ = ctx.WriteString("not found")
		return
	}
	root := p.(string)
	node, err := buildFileTree(root, root)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString(err.Error())
		return
	}
	ctx.JSON(node)
}

func extractZip(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return err
		}
		out.Close()
		rc.Close()
	}
	return nil
}

func buildFileTree(path, base string) (FileNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileNode{}, err
	}
	rel := strings.TrimPrefix(strings.TrimPrefix(path, base), string(os.PathSeparator))
	node := FileNode{Name: info.Name(), Path: rel, IsDir: info.IsDir()}
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return node, err
		}
		for _, e := range entries {
			child, err := buildFileTree(filepath.Join(path, e.Name()), base)
			if err == nil {
				node.Children = append(node.Children, child)
			}
		}
	}
	return node, nil
}
