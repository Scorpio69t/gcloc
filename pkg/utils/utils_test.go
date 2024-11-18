package utils

import (
	"github.com/Scorpio69t/gcloc/pkg/option"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestContainsComment(t *testing.T) {
	if !ContainComment(`int a; /* A takes care of counts */`, [][]string{{"/*", "*/"}}) {
		t.Errorf("invalid")
	}

	if !ContainComment(`bool f; /* `, [][]string{{"/*", "*/"}}) {
		t.Errorf("invalid")
	}

	if ContainComment(`}`, [][]string{{"/*", "*/"}}) {
		t.Errorf("invalid")
	}
}

func TestCheckMD5SumIgnore(t *testing.T) {
	fileCache := make(map[string]struct{})

	if CheckMD5Sum("./utils_test.go", fileCache) {
		t.Errorf("invalid sequence")
	}

	if !CheckMD5Sum("./utils_test.go", fileCache) {
		t.Errorf("invalid sequence")
	}
}

func TestCheckDefaultIgnore(t *testing.T) {
	appFS := afero.NewMemMapFs()
	if err := appFS.Mkdir("/test", os.ModeDir); err != nil {
		t.Fatal(err)
	}
	_, _ = appFS.Create("/test/one.go")

	fileInfo, _ := appFS.Stat("/")
	if !CheckDefaultIgnore("/", fileInfo, false) {
		t.Errorf("invalid logic: this is directory")
	}

	if !CheckDefaultIgnore("/", fileInfo, true) {
		t.Errorf("invalid logic: this is vcs file or directory")
	}

	fileInfo, _ = appFS.Stat("/test/one.go")
	if CheckDefaultIgnore("/test/one.go", fileInfo, false) {
		t.Errorf("invalid logic: should not ignore this file")
	}
}

type MockFileInfo struct {
	FileName    string
	IsDirectory bool
}

func (mfi MockFileInfo) Name() string       { return mfi.FileName }
func (mfi MockFileInfo) Size() int64        { return int64(8) }
func (mfi MockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mfi MockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi MockFileInfo) IsDir() bool        { return mfi.IsDirectory }
func (mfi MockFileInfo) Sys() interface{}   { return nil }

func TestCheckOptionMatch(t *testing.T) {
	opts := &option.GClocOptions{}
	fi := MockFileInfo{FileName: "/", IsDirectory: true}
	if !CheckOptionMatch("/", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is nil")
	}

	opts.ReNotMatchDir = regexp.MustCompile("thisisdir-not-match")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if !CheckOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is nil")
	}

	opts.ReNotMatchDir = regexp.MustCompile("thisisdir")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if CheckOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is ignore")
	}

	opts = &option.GClocOptions{}
	opts.ReMatchDir = regexp.MustCompile("thisisdir")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if !CheckOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is not ignore")
	}

	opts.ReMatchDir = regexp.MustCompile("thisisdir-not-match")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if CheckOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is ignore")
	}

	opts = &option.GClocOptions{}
	opts.ReNotMatchDir = regexp.MustCompile("thisisdir-not-match")
	opts.ReMatchDir = regexp.MustCompile("thisisdir")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if !CheckOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is not ignore")
	}
}
