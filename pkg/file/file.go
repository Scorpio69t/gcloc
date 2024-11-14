package file

import (
	"gcloc/pkg/option"
	log "gcloc/pkg/simplelog"
	"sort"
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
