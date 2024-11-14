package xml

import (
	"encoding/xml"
	"fmt"
	"gcloc/pkg/file"
	"gcloc/pkg/language"
)

// XMLResultType is the result type in XML format.
type XMLResultType int8

const (
	// XMLResultWithLangs is the result type for each language in XML format
	XMLResultWithLangs XMLResultType = iota
	// XMLResultWithFiles is the result type for each file in XML format
	XMLResultWithFiles
)

// XMLTotalLanguages is the total result in XML format.
type XMLTotalLanguages struct {
	SumFiles uint32 `xml:"sum_files,attr"`
	Codes    uint32 `xml:"codes,attr"`
	Comments uint32 `xml:"comments,attr"`
	Blanks   uint32 `xml:"blanks,attr"`
}

// XMLResultLanguages stores the results in XML format.
type XMLResultLanguages struct {
	Languages []language.GClocLanguage `xml:"language"`
	Total     XMLTotalLanguages        `xml:"total"`
}

// XMLTotalFiles is the total result per file in XML format.
type XMLTotalFiles struct {
	Codes    uint32 `xml:"codes,attr"`
	Comments uint32 `xml:"comments,attr"`
	Blanks   uint32 `xml:"blanks,attr"`
}

// XMLResultFiles stores per file results in XML format.
type XMLResultFiles struct {
	Files file.GClocFiles `xml:"file"`
	Total XMLTotalFiles   `xml:"total"`
}

// XMLResult stores the results in XML format.
type XMLResult struct {
	XMLName      xml.Name            `xml:"results"`
	XMLFiles     *XMLResultFiles     `xml:"files,omitempty"`
	XMLLanguages *XMLResultLanguages `xml:"languages,omitempty"`
}

// Encode outputs XMLResult in a human readable format.
func (x *XMLResult) Encode() {
	if output, err := xml.MarshalIndent(x, "", "  "); err == nil {
		fmt.Printf(xml.Header)
		fmt.Println(string(output))
	}
}

// NewXMLResultFromCloc returns XMLResult with default data set.
func NewXMLResultFromCloc(total *language.Language, sortedLanguages language.Languages, _ XMLResultType) *XMLResult {
	var langs []language.GClocLanguage
	for _, lang := range sortedLanguages {
		c := language.GClocLanguage{
			Name:      lang.Name,
			FileCount: uint32(len(lang.Files)),
			Codes:     lang.Codes,
			Comments:  lang.Comments,
			Blanks:    lang.Blanks,
		}
		langs = append(langs, c)
	}

	t := XMLTotalLanguages{
		Codes:    total.Codes,
		Comments: total.Comments,
		Blanks:   total.Blanks,
		SumFiles: total.Total,
	}

	f := &XMLResultLanguages{
		Languages: langs,
		Total:     t,
	}

	return &XMLResult{
		XMLLanguages: f,
	}
}
