package file

import (
	"bufio"
	"gcloc/pkg/bspool"
	"gcloc/pkg/language"
	"gcloc/pkg/option"
	log "gcloc/pkg/simplelog"
	"gcloc/pkg/utils"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

// GClocFile represents a file with its lines of code, comments and blanks.
type GClocFile struct {
	Name     string `xml:"name,attr" json:"name"`         // Name of the file
	Language string `xml:"language,attr" json:"language"` // Language of the file
	Codes    uint32 `xml:"codes,attr" json:"codes"`       // Number of lines of code
	Comments uint32 `xml:"comments,attr" json:"comments"` // Number of lines of comments]
	Blanks   uint32 `xml:"blanks,attr" json:"blanks"`     // Number of blank lines
}

// GClocFiles is a slice of GClocFile.
type GClocFiles []GClocFile

// SortByName sorts the files by name. (ASC)
func (gf GClocFiles) SortByName() {
	sort.Slice(gf, func(i, j int) bool {
		return gf[i].Name < gf[j].Name
	})
}

// SortByCodes sorts the files by number of lines of code. (DESC)
func (gf GClocFiles) SortByCodes() {
	sort.Slice(gf, func(i, j int) bool {
		return gf[i].Codes > gf[j].Codes
	})
}

// SortByComments sorts the files by number of lines of comments. (DESC)
func (gf GClocFiles) SortByComments() {
	sort.Slice(gf, func(i, j int) bool {
		if gf[i].Comments == gf[j].Comments {
			return gf[i].Codes > gf[j].Codes
		}
		return gf[i].Comments > gf[j].Comments
	})
}

// SortByBlanks sorts the files by number of blank lines. (DESC)
func (gf GClocFiles) SortByBlanks() {
	sort.Slice(gf, func(i, j int) bool {
		if gf[i].Blanks == gf[j].Blanks {
			return gf[i].Codes > gf[j].Codes
		}
		return gf[i].Blanks > gf[j].Blanks
	})
}

// Len returns the number of files.
func (gf GClocFiles) Len() int {
	return len(gf)
}

// AnalyzeFile analyzes a file and returns a GClocFile.
func AnalyzeFile(filename string, language *language.Language, opts *option.GClocOptions) *GClocFile {
	fp, err := os.Open(filename)
	if err != nil {
		// ignore error
		return &GClocFile{Name: filename}
	}

	defer func(fp *os.File) {
		err := fp.Close()
		if err != nil {
			log.Error("Failed to close file %s: %v", filename, err)
		}
	}(fp)

	return AnalyzeReader(filename, language, fp, opts)
}

func AnalyzeReader(filename string, language *language.Language, file io.Reader, opts *option.GClocOptions) *GClocFile {
	if opts.Debug {
		log.Info("Analyzing %s", filename)
	}

	gClocFile := &GClocFile{
		Name:     filename,
		Language: language.Name,
	}

	isFirstLine := true
	var inComments [][2]string

	buf := bspool.GetByteSlice()   // Get a buffer from the pool
	defer bspool.PutByteSlice(buf) // Return the buffer to the pool

	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf.Bytes(), opts.MaxLineLength)

	for scanner.Scan() {
		lineOrg := scanner.Text()
		line := strings.TrimSpace(lineOrg)

		if len(strings.TrimSpace(line)) == 0 {
			onBlank(gClocFile, opts, len(inComments) > 0, line, lineOrg)
			continue
		}

		// shebang line is counted as a code line
		if isFirstLine && strings.HasPrefix(line, "#!") {
			isFirstLine = false
			onCode(gClocFile, opts, len(inComments) > 0, line, lineOrg)
			continue
		}

		if len(inComments) == 0 {
			if isFirstLine {
				line = utils.TrimBOM(line)
			}

			if len(language.RegexLineComments) > 0 {
				if handleSingleLineCommentsRegex(language, gClocFile, opts, line, lineOrg, inComments) {
					continue
				}
			} else {
				if handleSingleLineComments(language, gClocFile, opts, line, lineOrg, inComments) {
					continue
				}
			}

			if len(language.MultipleLines) == 0 {
				onCode(gClocFile, opts, len(inComments) > 0, line, lineOrg)
				continue
			}
		}

		if len(inComments) == 0 && !utils.ContainComment(line, language.MultipleLines) {
			onCode(gClocFile, opts, len(inComments) > 0, line, lineOrg)
			continue
		}

		if len(language.MultipleLines) == 1 &&
			len(language.MultipleLines[0]) == 2 &&
			language.MultipleLines[0][0] == "" {
			onCode(gClocFile, opts, len(inComments) > 0, line, lineOrg)
			continue
		}

		lenLine := len(line)
		codeFlags := make([]bool, len(language.MultipleLines))
		for pos := 0; pos < lenLine; {
			for idx, ml := range language.MultipleLines {
				begin, end := ml[0], ml[1]
				lenBegin := len(begin)

				if pos+lenBegin <= lenLine &&
					strings.HasPrefix(line[pos:], begin) &&
					(begin != end || len(inComments) == 0) {
					pos += lenBegin
					inComments = append(inComments, [2]string{begin, end})
					continue
				}

				if n := len(inComments); n > 0 {
					last := inComments[n-1]
					if pos+len(last[1]) <= lenLine && strings.HasPrefix(line[pos:], last[1]) {
						inComments = inComments[:n-1]
						pos += len(last[1])
					}
				} else if pos < lenLine && !unicode.IsSpace(utils.NextRune(line[pos:])) {
					codeFlags[idx] = true
				}
			}
			pos++
		}

		isCode := true
		for _, b := range codeFlags {
			if !b {
				isCode = false
			}
		}

		if isCode {
			onCode(gClocFile, opts, len(inComments) > 0, line, lineOrg)
		} else {
			onComment(gClocFile, opts, len(inComments) > 0, line, lineOrg)
		}
	}

	return gClocFile
}

// onBlank is called when a blank line is found.
func onBlank(gClocFile *GClocFile, opts *option.GClocOptions, isInComments bool, line, lineOrg string) {
	gClocFile.Blanks++

	if opts.OnBlank != nil {
		opts.OnBlank(line)
	}

	if opts.Debug {
		log.Info("[BLANK, codes:%d, comments:%d, blanks:%d, isInComments:%v] %s",
			gClocFile.Codes, gClocFile.Comments, gClocFile.Blanks, isInComments, lineOrg)
	}
}

// onComment is called when a comment line is found.
func onComment(gClocFile *GClocFile, opts *option.GClocOptions, isInComments bool, line, lineOrg string) {
	gClocFile.Comments++

	if opts.OnComment != nil {
		opts.OnComment(line)
	}

	if opts.Debug {
		log.Info("[COMMENT, codes:%d, comments:%d, blanks:%d, isInComments:%v] %s",
			gClocFile.Codes, gClocFile.Comments, gClocFile.Blanks, isInComments, lineOrg)
	}
}

// onCode is called when a code line is found.
func onCode(gClocFile *GClocFile, opts *option.GClocOptions, isInComments bool, line, lineOrg string) {
	gClocFile.Codes++

	if opts.OnCode != nil {
		opts.OnCode(line)
	}

	if opts.Debug {
		log.Info("[CODE, codes:%d, comments:%d, blanks:%d, isInComments:%v] %s",
			gClocFile.Codes, gClocFile.Comments, gClocFile.Blanks, isInComments, lineOrg)
	}
}

// handleSingleLineCommentsRegex handles single line comments with regular expressions.
func handleSingleLineCommentsRegex(language *language.Language, gClocFile *GClocFile,
	opts *option.GClocOptions, line, lineOrg string, inComments [][2]string) bool {
	for _, singleCommentRegex := range language.RegexLineComments {
		if singleCommentRegex.MatchString(line) {
			// check if single comment is a prefix of multi comment
			for _, ml := range language.MultipleLines {
				if ml[0] != "" && strings.HasPrefix(line, ml[0]) {
					return false
				}
			}
			onComment(gClocFile, opts, len(inComments) > 0, line, lineOrg)
			return true
		}
	}
	return false
}

// handleSingleLineComments handles single line comments.
func handleSingleLineComments(language *language.Language, gClocFile *GClocFile,
	opts *option.GClocOptions, line, lineOrg string, inComments [][2]string) bool {
	for _, singleComment := range language.LineComments {
		if strings.HasPrefix(line, singleComment) {
			// check if single comment is a prefix of multi comment
			for _, ml := range language.MultipleLines {
				if ml[0] != "" && strings.HasPrefix(line, ml[0]) {
					return false
				}
			}
			onComment(gClocFile, opts, len(inComments) > 0, line, lineOrg)
			return true
		}
	}
	return false
}
