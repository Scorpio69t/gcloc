package utils

import (
	"crypto/md5"
	"fmt"
	"gcloc/pkg/language"
	"gcloc/pkg/option"
	log "gcloc/pkg/simplelog"
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
func CheckMD5Sum(path string, fileCache map[string]struct{}) bool {
	// Get the file's MD5 sum
	md5Sum, err := MD5Sum(path)
	if err != nil {
		return true // ignore the error
	}

	// Check if the file's MD5 sum is in the cache
	hash := fmt.Sprintf("%x", md5Sum)
	if _, ok := fileCache[hash]; ok {
		return true
	}

	// Add the file's MD5 sum to the cache
	fileCache[hash] = struct{}{}
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

// shouldIgnore returns true if the path should be ignored.
func shouldIgnore(path string, info os.FileInfo, vcsInRoot bool, opts *option.GClocOptions) bool {
	if CheckDefaultIgnore(path, info, vcsInRoot) {
		return true
	}
	if !CheckOptionMatch(path, info, opts) {
		return true
	}
	return false
}

// processFile processes the file.
func processFile(path, ext string, languages *language.DefinedLanguages, opts *option.GClocOptions, result map[string]*language.Language, fileCache map[string]struct{}) {
	if targetExt, ok := language.FileExtensions[ext]; ok {
		if _, ok := opts.ExcludeExts[targetExt]; ok {
			return
		}
		if len(opts.IncludeLanguages) != 0 {
			if _, ok := opts.IncludeLanguages[targetExt]; !ok {
				return
			}
		}
		if !opts.SkipDuplicated {
			if CheckMD5Sum(path, fileCache) {
				if opts.Debug {
					log.Info("[ignore=%v] find same md5", path)
				}
				return
			}
		}
		addFileToResult(path, targetExt, languages, result)
	}
}

// addFileToResult adds the file to the result.
func addFileToResult(path, targetExt string, languages *language.DefinedLanguages, result map[string]*language.Language) {
	if _, ok := result[targetExt]; !ok {
		definedLang := language.NewLanguage(
			languages.Langs[targetExt].Name,
			languages.Langs[targetExt].LineComments,
			languages.Langs[targetExt].MultipleLines,
		)
		if len(languages.Langs[targetExt].RegexLineComments) > 0 {
			definedLang.RegexLineComments = languages.Langs[targetExt].RegexLineComments
		}
		result[targetExt] = definedLang
	}
	result[targetExt].Files = append(result[targetExt].Files, path)
}

// GetAllFiles return all the files to be analyzed in paths.
func GetAllFiles(paths []string, languages *language.DefinedLanguages, opts *option.GClocOptions) (result map[string]*language.Language, err error) {
	result = make(map[string]*language.Language)
	fileCache := make(map[string]struct{})

	for _, root := range paths {
		vcsInRoot := IsVCSDir(root)
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
				return nil
			}

			if shouldIgnore(path, info, vcsInRoot, opts) {
				return nil
			}

			if ext, ok := language.GetFileType(path, opts); ok {
				processFile(path, ext, languages, opts, result, fileCache)
			}
			return nil
		})

		if err != nil {
			if opts.Debug {
				log.Error("error: %v", err)
			}
			return nil, err
		}
	}

	return
}
