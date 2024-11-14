package json

import (
	"gcloc/pkg/file"
	"gcloc/pkg/language"
)

// LanguagesResult defines the result of the analysis in JSON format.
type LanguagesResult struct {
	Languages []language.GClocLanguage `json:"languages"`
	Total     language.GClocLanguage   `json:"total"`
}

// FilesResult defines the result of the analysis(by files) in JSON format.
type FilesResult struct {
	Files file.GClocFiles        `json:"files"`
	Total language.GClocLanguage `json:"total"`
}

// NewLanguagesResultFromGCloc creates a new LanguagesResult from a language.Languages.
func NewLanguagesResultFromGCloc(total *language.Language, sortedLanguages language.Languages) LanguagesResult {
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

	t := language.GClocLanguage{
		FileCount: total.Total,
		Codes:     total.Codes,
		Comments:  total.Comments,
		Blanks:    total.Blanks,
	}

	return LanguagesResult{
		Languages: langs,
		Total:     t,
	}
}

// NewFilesResultFromGCloc creates a new FilesResult from a file.GClocFiles.
func NewFilesResultFromGCloc(total *language.Language, sortedFiles file.GClocFiles) FilesResult {
	t := language.GClocLanguage{
		FileCount: total.Total,
		Codes:     total.Codes,
		Comments:  total.Comments,
		Blanks:    total.Blanks,
	}

	return FilesResult{
		Files: sortedFiles,
		Total: t,
	}
}
