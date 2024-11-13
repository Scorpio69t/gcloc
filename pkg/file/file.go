package file

import "sort"

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
