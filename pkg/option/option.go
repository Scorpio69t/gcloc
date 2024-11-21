package option

import "regexp"

// GClocOptions is a struct that holds the options for the gcloc command.
type GClocOptions struct {
	Debug            bool                // Debug mode
	SkipDuplicated   bool                // Skip duplicated files
	MaxLineLength    int                 // Maximum line length
	ExcludeExts      map[string]struct{} // Excluded extensions
	ExcludeLanguages map[string]struct{} // Excluded languages
	IncludeLanguages map[string]struct{} // Included languages
	ReNotMatch       *regexp.Regexp      // Regular expression for not matching files
	ReMatch          *regexp.Regexp      // Regular expression for matching files
	ReNotMatchDir    *regexp.Regexp      // Regular expression for not matching directories
	ReMatchDir       *regexp.Regexp      // Regular expression for matching directories

	// OnCode is triggered for each line of code.
	OnCode func(line string)
	// OnBlack is triggered for each blank line.
	OnBlank func(line string)
	// OnComment is triggered for each line of comments.
	OnComment func(line string)
}

// NewGClocOptions returns a new GClocOptions struct.
func NewGClocOptions() *GClocOptions {
	return &GClocOptions{
		Debug:            false,
		SkipDuplicated:   false,
		MaxLineLength:    1024 * 1024, // 1MB
		ExcludeExts:      make(map[string]struct{}),
		ExcludeLanguages: make(map[string]struct{}),
		IncludeLanguages: make(map[string]struct{}),
	}
}
