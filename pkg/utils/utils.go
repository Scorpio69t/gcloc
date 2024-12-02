package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/Scorpio69t/gcloc/pkg/option"
	"github.com/Scorpio69t/gcloc/pkg/syncmap"
	"os"
	"path/filepath"
	"strings"
)

const (
	bom = "\xef\xbb\xbf" // UTF-8 BOM
)

// FileExits returns true if the file exists.
func FileExits(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// TrimBOM trims the UTF-8 BOM from the line.
func TrimBOM(line string) string {
	if strings.HasPrefix(line, bom) {
		return line[len(bom):]
	}

	return line
}

// ContainComment returns true if the line contains a comment.
func ContainComment(line string, multipleLines [][]string) bool {
	for _, comments := range multipleLines {
		for _, comment := range comments {
			if strings.Contains(line, comment) {
				// check if the comment is in the "string"
				if strings.Contains(line, "\""+comment) || strings.Contains(line, comment+"\"") {
					continue
				}
				return true
			}
		}
	}

	return false
}

// NextRune returns the next rune in the string.
func NextRune(str string) rune {
	for _, r := range str {
		return r
	}

	return 0
}

// MD5Sum returns the MD5 sum of the file.
func MD5Sum(path string) ([16]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return [16]byte{}, err
	}

	return md5.Sum(content), nil
}

// CheckMD5Sum checks if the file's MD5 sum is in the cache.
func CheckMD5Sum(path string, fileCache *syncmap.SyncMap[string, struct{}]) bool {
	// Get the file's MD5 sum
	md5Sum, err := MD5Sum(path)
	if err != nil {
		return true // ignore the error
	}

	// Check if the file's MD5 sum is in the cache
	hash := fmt.Sprintf("%x", md5Sum)
	if ok := fileCache.Has(hash); ok {
		return true
	}

	// Add the file's MD5 sum to the cache
	fileCache.Store(hash, struct{}{})
	return false
}

// IsVCSDir returns true if the path is a VCS directory.
func IsVCSDir(path string) bool {
	path = strings.TrimPrefix(path, string(os.PathSeparator))
	vcsDirs := []string{".bzr", ".cvs", ".hg", ".git", ".svn"}
	for _, dir := range vcsDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

// CheckDefaultIgnore returns true if the path should be ignored.
func CheckDefaultIgnore(path string, info os.FileInfo, isVCS bool) bool {
	if info.IsDir() {
		// directory is ignored
		return true
	}

	if !isVCS && IsVCSDir(path) {
		// vcs file or directory is ignored
		return true
	}

	return false
}

// CheckOptionMatch returns true if the path matches the options.
func CheckOptionMatch(path string, info os.FileInfo, opts *option.GClocOptions) bool {
	if opts.ReNotMatch != nil && opts.ReNotMatch.MatchString(info.Name()) {
		return false
	}

	if opts.ReMatch != nil && !opts.ReMatch.MatchString(info.Name()) {
		return false
	}

	dir := filepath.Dir(path)
	if opts.ReNotMatchDir != nil && opts.ReNotMatchDir.MatchString(dir) {
		return false
	}

	if opts.ReMatchDir != nil && !opts.ReMatchDir.MatchString(dir) {
		return false
	}

	return true
}

// IsGitURL returns true if the path is a git URL.
func IsGitURL(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasSuffix(path, ".git")
}
