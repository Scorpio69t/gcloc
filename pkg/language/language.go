package language

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Scorpio69t/gcloc/pkg/option"
	log "github.com/Scorpio69t/gcloc/pkg/simplelog"
	"github.com/Scorpio69t/gcloc/pkg/syncmap"
	"github.com/Scorpio69t/gcloc/pkg/utils"
	"github.com/go-enry/go-enry/v2"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// GClocLanguage is provided for xml-cloc and json format.
type GClocLanguage struct {
	Name      string `xml:"name,attr" json:"name,omitempty"`   // Name of the language
	FileCount uint32 `xml:"file_count,attr" json:"file_count"` // Number of files
	Codes     uint32 `xml:"codes,attr" json:"codes"`           // Number of lines of code
	Comments  uint32 `xml:"comments,attr" json:"comments"`     // Number of lines of comments
	Blanks    uint32 `xml:"blanks,attr" json:"blanks"`         // Number of blank lines
}

// Language is the struct for the language.
type Language struct {
	Name              string
	LineComments      []string
	RegexLineComments []*regexp.Regexp
	MultipleLines     [][]string
	Files             []string
	Codes             uint32
	Comments          uint32
	Blanks            uint32
	Total             uint32
}

// DefinedLanguages is the struct for the defined languages.
type DefinedLanguages struct {
	Langs map[string]*Language
}

// Languages is the slice of Language.
type Languages []Language

var (
	reShebangEnv  = regexp.MustCompile(`^#! *(\S+/env) ([a-zA-Z]+)`)
	reShebangLang = regexp.MustCompile(`^#! *[.a-zA-Z/]+/([a-zA-Z]+)`)
)

// SortByName sorts the languages by name. (ASC)
func (ls Languages) SortByName() {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Name < ls[j].Name
	})
}

// SortByFiles sorts the languages by files. (DESC)
func (ls Languages) SortByFiles() {
	sort.Slice(ls, func(i, j int) bool {
		iLen := len(ls[i].Files)
		jLen := len(ls[j].Files)

		if iLen == jLen {
			return ls[i].Codes > ls[j].Codes
		}

		return iLen > jLen
	})
}

// SortByCodes sorts the languages by codes. (DESC)
func (ls Languages) SortByCodes() {
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Codes > ls[j].Codes
	})
}

// SortByComments sorts the languages by comments. (DESC)
func (ls Languages) SortByComments() {
	sort.Slice(ls, func(i, j int) bool {
		if ls[i].Comments == ls[j].Comments {
			return ls[i].Codes > ls[j].Codes
		}

		return ls[i].Comments > ls[j].Comments
	})
}

// SortByBlanks sorts the languages by blanks. (DESC)
func (ls Languages) SortByBlanks() {
	sort.Slice(ls, func(i, j int) bool {
		if ls[i].Blanks == ls[j].Blanks {
			return ls[i].Codes > ls[j].Codes
		}

		return ls[i].Blanks > ls[j].Blanks
	})
}

// Len returns the length of the languages.
func (ls Languages) Len() int {
	return len(ls)
}

// NewLanguage creates a new language.
func NewLanguage(name string, lineComments []string, multipleLines [][]string) *Language {
	l := &Language{
		Name:          name,
		LineComments:  lineComments,
		MultipleLines: multipleLines,
		Files:         []string{},
	}

	l.RegexLineComments = make([]*regexp.Regexp, 0)

	return l
}

// AddRegexLineComments adds the regex line comments.
func (l *Language) AddRegexLineComments(regexLineComments []string) *Language {
	for _, regex := range regexLineComments {
		l.RegexLineComments = append(l.RegexLineComments, regexp.MustCompile(regex))
	}
	return l
}

func NewDefinedLanguages() *DefinedLanguages {
	return &DefinedLanguages{
		Langs: map[string]*Language{
			"ActionScript":        NewLanguage("ActionScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Ada":                 NewLanguage("Ada", []string{"--"}, [][]string{{"", ""}}),
			"Alda":                NewLanguage("Alda", []string{"#"}, [][]string{{"", ""}}),
			"Ant":                 NewLanguage("Ant", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"ANTLR":               NewLanguage("ANTLR", []string{"//"}, [][]string{{"/*", "*/"}}),
			"AsciiDoc":            NewLanguage("AsciiDoc", []string{}, [][]string{{"", ""}}),
			"Assembly":            NewLanguage("Assembly", []string{"//", ";", "#", "@", "|", "!"}, [][]string{{"/*", "*/"}}),
			"ATS":                 NewLanguage("ATS", []string{"//"}, [][]string{{"/*", "*/"}, {"(*", "*)"}}),
			"AutoHotkey":          NewLanguage("AutoHotkey", []string{";"}, [][]string{{"", ""}}),
			"Awk":                 NewLanguage("Awk", []string{"#"}, [][]string{{"", ""}}),
			"Arduino Sketch":      NewLanguage("Arduino Sketch", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Ballerina":           NewLanguage("Ballerina", []string{"//"}, [][]string{{"", ""}}),
			"Batch":               NewLanguage("Batch", []string{"REM", "rem"}, [][]string{{"", ""}}),
			"Berry":               NewLanguage("Berry", []string{"#"}, [][]string{{"#-", "-#"}}),
			"BASH":                NewLanguage("BASH", []string{"#"}, [][]string{{"", ""}}),
			"Bicep":               NewLanguage("Bicep", []string{"//"}, [][]string{{"/*", "*/"}}),
			"BitBake":             NewLanguage("BitBake", []string{"#"}, [][]string{{"", ""}}),
			"C":                   NewLanguage("C", []string{"//"}, [][]string{{"/*", "*/"}}),
			"C Header":            NewLanguage("C Header", []string{"//"}, [][]string{{"/*", "*/"}}),
			"C Shell":             NewLanguage("C Shell", []string{"#"}, [][]string{{"", ""}}),
			"Cairo":               NewLanguage("Cairo", []string{"//"}, [][]string{{"", ""}}),
			"Carbon":              NewLanguage("Carbon", []string{"//"}, [][]string{{"", ""}}),
			"Cap'n Proto":         NewLanguage("Cap'n Proto", []string{"#"}, [][]string{{"", ""}}),
			"Carp":                NewLanguage("Carp", []string{";"}, [][]string{{"", ""}}),
			"C#":                  NewLanguage("C#", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Chapel":              NewLanguage("Chapel", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Circom":              NewLanguage("Circom", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Clojure":             NewLanguage("Clojure", []string{"#", "#_"}, [][]string{{"", ""}}),
			"COBOL":               NewLanguage("COBOL", []string{"*", "/"}, [][]string{{"", ""}}),
			"CoffeeScript":        NewLanguage("CoffeeScript", []string{"#"}, [][]string{{"###", "###"}}),
			"Coq":                 NewLanguage("Coq", []string{"(*"}, [][]string{{"(*", "*)"}}),
			"ColdFusion":          NewLanguage("ColdFusion", []string{}, [][]string{{"<!---", "--->"}}),
			"ColdFusion CFScript": NewLanguage("ColdFusion CFScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"CMake":               NewLanguage("CMake", []string{"#"}, [][]string{{"", ""}}),
			"C++":                 NewLanguage("C++", []string{"//"}, [][]string{{"/*", "*/"}}),
			"C++ Header":          NewLanguage("C++ Header", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Crystal":             NewLanguage("Crystal", []string{"#"}, [][]string{{"", ""}}),
			"CSS":                 NewLanguage("CSS", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Cython":              NewLanguage("Cython", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}}),
			"CUDA":                NewLanguage("CUDA", []string{"//"}, [][]string{{"/*", "*/"}}),
			"D":                   NewLanguage("D", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Dart":                NewLanguage("Dart", []string{"//", "///"}, [][]string{{"/*", "*/"}}),
			"Dhall":               NewLanguage("Dhall", []string{"--"}, [][]string{{"{-", "-}"}}),
			"DTrace":              NewLanguage("DTrace", []string{}, [][]string{{"/*", "*/"}}),
			"Device Tree":         NewLanguage("Device Tree", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Eiffel":              NewLanguage("Eiffel", []string{"--"}, [][]string{{"", ""}}),
			"Elm":                 NewLanguage("Elm", []string{"--"}, [][]string{{"{-", "-}"}}),
			"Elixir":              NewLanguage("Elixir", []string{"#"}, [][]string{{"", ""}}),
			"Erlang":              NewLanguage("Erlang", []string{"%"}, [][]string{{"", ""}}),
			"Expect":              NewLanguage("Expect", []string{"#"}, [][]string{{"", ""}}),
			"Fish":                NewLanguage("Fish", []string{"#"}, [][]string{{"", ""}}),
			"Frege":               NewLanguage("Frege", []string{"--"}, [][]string{{"{-", "-}"}}),
			"F*":                  NewLanguage("F*", []string{"(*", "//"}, [][]string{{"(*", "*)"}}),
			"F#":                  NewLanguage("F#", []string{"(*"}, [][]string{{"(*", "*)"}}),
			"Lean":                NewLanguage("Lean", []string{"--"}, [][]string{{"/-", "-/"}}),
			"Logtalk":             NewLanguage("Logtalk", []string{"%"}, [][]string{{"", ""}}),
			"Lua":                 NewLanguage("Lua", []string{"--"}, [][]string{{"--[[", "]]"}}),
			"Lilypond":            NewLanguage("Lilypond", []string{"%"}, [][]string{{"", ""}}),
			"LISP":                NewLanguage("LISP", []string{";;"}, [][]string{{"#|", "|#"}}),
			"LiveScript":          NewLanguage("LiveScript", []string{"#"}, [][]string{{"/*", "*/"}}),
			"Factor":              NewLanguage("Factor", []string{"! "}, [][]string{{"", ""}}),
			"FORTRAN Legacy":      NewLanguage("FORTRAN Legacy", []string{"c", "C", "!", "*"}, [][]string{{"", ""}}),
			"FORTRAN Modern":      NewLanguage("FORTRAN Modern", []string{"!"}, [][]string{{"", ""}}),
			"Gherkin":             NewLanguage("Gherkin", []string{"#"}, [][]string{{"", ""}}),
			"Gleam":               NewLanguage("Gleam", []string{"//"}, [][]string{{"", ""}}),
			"GLSL":                NewLanguage("GLSL", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Go":                  NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Groovy":              NewLanguage("Groovy", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Handlebars":          NewLanguage("Handlebars", []string{}, [][]string{{"<!--", "-->"}, {"{{!", "}}"}}),
			"Haskell":             NewLanguage("Haskell", []string{"--"}, [][]string{{"{-", "-}"}}),
			"Haxe":                NewLanguage("Haxe", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Hare":                NewLanguage("Hare", []string{"//"}, [][]string{{"", ""}}),
			"HLSL":                NewLanguage("HLSL", []string{"//"}, [][]string{{"/*", "*/"}}),
			"HTML":                NewLanguage("HTML", []string{"//", "<!--"}, [][]string{{"<!--", "-->"}}),
			"Hy":                  NewLanguage("Hy", []string{";"}, [][]string{{"", ""}}),
			"Idris":               NewLanguage("Idris", []string{"--"}, [][]string{{"{-", "-}"}}),
			"Imba":                NewLanguage("Imba", []string{"#"}, [][]string{{"###", "###"}}),
			"Io":                  NewLanguage("Io", []string{"//", "#"}, [][]string{{"/*", "*/"}}),
			"SKILL":               NewLanguage("SKILL", []string{";"}, [][]string{{"/*", "*/"}}),
			"JAI":                 NewLanguage("JAI", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Janet":               NewLanguage("Janet", []string{"#"}, [][]string{{"", ""}}),
			"Java":                NewLanguage("Java", []string{"//"}, [][]string{{"/*", "*/"}}),
			"JSP":                 NewLanguage("JSP", []string{"//"}, [][]string{{"/*", "*/"}}),
			"JavaScript":          NewLanguage("JavaScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Julia":               NewLanguage("Julia", []string{"#"}, [][]string{{"#:=", ":=#"}}),
			"Jupyter Notebook":    NewLanguage("Jupyter Notebook", []string{"#"}, [][]string{{"", ""}}),
			"Just":                NewLanguage("Just", []string{"#"}, [][]string{{"", ""}}).AddRegexLineComments([]string{`^#[^!].*`}),
			"JSON":                NewLanguage("JSON", []string{}, [][]string{{"", ""}}),
			"JSX":                 NewLanguage("JSX", []string{"//"}, [][]string{{"/*", "*/"}}),
			"KakouneScript":       NewLanguage("KakouneScript", []string{"#"}, [][]string{{"", ""}}),
			"Koka":                NewLanguage("Koka", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Kotlin":              NewLanguage("Kotlin", []string{"//"}, [][]string{{"/*", "*/"}}),
			"LD Script":           NewLanguage("LD Script", []string{"//"}, [][]string{{"/*", "*/"}}),
			"LESS":                NewLanguage("LESS", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Objective-C":         NewLanguage("Objective-C", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Markdown":            NewLanguage("Markdown", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Motoko":              NewLanguage("Motoko", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Nearley":             NewLanguage("Nearley", []string{"#"}, [][]string{{"", ""}}),
			"Nix":                 NewLanguage("Nix", []string{"#"}, [][]string{{"/*", "*/"}}),
			"NSIS":                NewLanguage("NSIS", []string{"#", ";"}, [][]string{{"/*", "*/"}}),
			"Nu":                  NewLanguage("Nu", []string{";", "#"}, [][]string{{"", ""}}),
			"OCaml":               NewLanguage("OCaml", []string{}, [][]string{{"(*", "*)"}}),
			"Objective-C++":       NewLanguage("Objective-C++", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Makefile":            NewLanguage("Makefile", []string{"#"}, [][]string{{"", ""}}),
			"MATLAB":              NewLanguage("MATLAB", []string{"%"}, [][]string{{"%{", "}%"}}),
			"Mercury":             NewLanguage("Mercury", []string{"%"}, [][]string{{"/*", "*/"}}),
			"Maven":               NewLanguage("Maven", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Meson":               NewLanguage("Meson", []string{"#"}, [][]string{{"", ""}}),
			"Mojo":                NewLanguage("Mojo", []string{"#"}, [][]string{{"", ""}}),
			"Move":                NewLanguage("Move", []string{"//"}, [][]string{{"", ""}}),
			"Mustache":            NewLanguage("Mustache", []string{}, [][]string{{"{{!", "}}"}}),
			"M4":                  NewLanguage("M4", []string{"#"}, [][]string{{"", ""}}),
			"Nim":                 NewLanguage("Nim", []string{"#"}, [][]string{{"#[", "]#"}}),
			"Nunjucks":            NewLanguage("Nunjucks", []string{}, [][]string{{"{#", "#}"}, {"<!--", "-->"}}),
			"lex":                 NewLanguage("lex", []string{}, [][]string{{"/*", "*/"}}),
			"Odin":                NewLanguage("Odin", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Ohm":                 NewLanguage("Ohm", []string{"//"}, [][]string{{"/*", "*/"}}),
			"PHP":                 NewLanguage("PHP", []string{"#", "//"}, [][]string{{"/*", "*/"}}),
			"Pascal":              NewLanguage("Pascal", []string{"//"}, [][]string{{"{", ")"}}),
			"Perl":                NewLanguage("Perl", []string{"#"}, [][]string{{":=", ":=cut"}}),
			"Plain Text":          NewLanguage("Plain Text", []string{}, [][]string{{"", ""}}),
			"Plan9 Shell":         NewLanguage("Plan9 Shell", []string{"#"}, [][]string{{"", ""}}),
			"Pony":                NewLanguage("Pony", []string{"//"}, [][]string{{"/*", "*/"}}),
			"PowerShell":          NewLanguage("PowerShell", []string{"#"}, [][]string{{"<#", "#>"}}),
			"Polly":               NewLanguage("Polly", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Protocol Buffers":    NewLanguage("Protocol Buffers", []string{"//"}, [][]string{{"", ""}}),
			"PRQL":                NewLanguage("PRQL", []string{"#"}, [][]string{{"", ""}}),
			"Python":              NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}}),
			"Q":                   NewLanguage("Q", []string{"/ "}, [][]string{{"\\", "/"}, {"/", "\\"}}),
			"QML":                 NewLanguage("QML", []string{"//"}, [][]string{{"/*", "*/"}}),
			"R":                   NewLanguage("R", []string{"#"}, [][]string{{"", ""}}),
			"Rebol":               NewLanguage("Rebol", []string{";"}, [][]string{{"", ""}}),
			"Red":                 NewLanguage("Red", []string{";"}, [][]string{{"", ""}}),
			"Rego":                NewLanguage("Rego", []string{"#"}, [][]string{{"", ""}}),
			"RMarkdown":           NewLanguage("RMarkdown", []string{}, [][]string{{"", ""}}),
			"RAML":                NewLanguage("RAML", []string{"#"}, [][]string{{"", ""}}),
			"Racket":              NewLanguage("Racket", []string{";"}, [][]string{{"#|", "|#"}}),
			"ReStructuredText":    NewLanguage("ReStructuredText", []string{}, [][]string{{"", ""}}),
			"Ring":                NewLanguage("Ring", []string{"#", "//"}, [][]string{{"/*", "*/"}}),
			"Ruby":                NewLanguage("Ruby", []string{"#"}, [][]string{{":=begin", ":=end"}}),
			"Ruby HTML":           NewLanguage("Ruby HTML", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Rust":                NewLanguage("Rust", []string{"//", "///", "//!"}, [][]string{{"/*", "*/"}}),
			"Scala":               NewLanguage("Scala", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Sass":                NewLanguage("Sass", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Scheme":              NewLanguage("Scheme", []string{";"}, [][]string{{"#|", "|#"}}),
			"sed":                 NewLanguage("sed", []string{"#"}, [][]string{{"", ""}}),
			"Stan":                NewLanguage("Stan", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Starlark":            NewLanguage("Starlark", []string{"#"}, [][]string{{"", ""}}),
			"Solidity":            NewLanguage("Solidity", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Bourne Shell":        NewLanguage("Bourne Shell", []string{"#"}, [][]string{{"", ""}}),
			"Standard ML":         NewLanguage("Standard ML", []string{}, [][]string{{"(*", "*)"}}),
			"SQL":                 NewLanguage("SQL", []string{"--"}, [][]string{{"/*", "*/"}}),
			"Svelte":              NewLanguage("Svelte", []string{"//"}, [][]string{{"/*", "*/"}, {"<!--", "-->"}}),
			"Swift":               NewLanguage("Swift", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Terra":               NewLanguage("Terra", []string{"--"}, [][]string{{"--[[", "]]"}}),
			"TeX":                 NewLanguage("TeX", []string{"%"}, [][]string{{"", ""}}),
			"Inno Setup":          NewLanguage("Inno Setup", []string{";"}, [][]string{{"", ""}}),
			"Isabelle":            NewLanguage("Isabelle", []string{}, [][]string{{"(*", "*)"}}),
			"TLA":                 NewLanguage("TLA", []string{"\\*"}, [][]string{{"(*", "*)"}}),
			"Tcl/Tk":              NewLanguage("Tcl/Tk", []string{"#"}, [][]string{{"", ""}}),
			"TOML":                NewLanguage("TOML", []string{"#"}, [][]string{{"", ""}}),
			"TypeScript":          NewLanguage("TypeScript", []string{"//"}, [][]string{{"/*", "*/"}}),
			"HCL":                 NewLanguage("HCL", []string{"#", "//"}, [][]string{{"/*", "*/"}}),
			"Umka":                NewLanguage("Umka", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Unity-Prefab":        NewLanguage("Unity-Prefab", []string{}, [][]string{{"", ""}}),
			"MSBuild script":      NewLanguage("MSBuild script", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Vala":                NewLanguage("Vala", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Verilog":             NewLanguage("Verilog", []string{"//"}, [][]string{{"/*", "*/"}}),
			"VimL":                NewLanguage("VimL", []string{`"`}, [][]string{{"", ""}}),
			"Visual Basic":        NewLanguage("Visual Basic", []string{"'"}, [][]string{{"", ""}}),
			"Vue":                 NewLanguage("Vue", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"Vyper":               NewLanguage("Vyper", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}}),
			"WiX":                 NewLanguage("WiX", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XML":                 NewLanguage("XML", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XML resource":        NewLanguage("XML resource", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XSLT":                NewLanguage("XSLT", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"XSD":                 NewLanguage("XSD", []string{"<!--"}, [][]string{{"<!--", "-->"}}),
			"YAML":                NewLanguage("YAML", []string{"#"}, [][]string{{"", ""}}),
			"Yacc":                NewLanguage("Yacc", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Yul":                 NewLanguage("Yul", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Zephir":              NewLanguage("Zephir", []string{"//"}, [][]string{{"/*", "*/"}}),
			"Zig":                 NewLanguage("Zig", []string{"//", "///"}, [][]string{{"", ""}}),
			"Zsh":                 NewLanguage("Zsh", []string{"#"}, [][]string{{"", ""}}),
		},
	}
}

// GetFormattedString return DefinedLanguages as a human-readable string.
func (l *DefinedLanguages) GetFormattedString() string {
	var buf bytes.Buffer
	var printLangs []string

	for _, lang := range l.Langs {
		printLangs = append(printLangs, lang.Name)
	}

	sort.Strings(printLangs)
	for _, lang := range printLangs {
		buf.WriteString(fmt.Sprintf("%-30v (%s)\n", lang, lang2Exts(lang)))
	}

	return buf.String()
}

// loadFileExtsFromJson loads the file extensions from the JSON file.
func loadFileExtsFromJson(filePath string) (map[string]string, error) {
	if !utils.FileExits(filePath) {
		return nil, errors.New("file extensions file does not exist")
	}

	var fileExts map[string]string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(file).Decode(&fileExts)
	if err != nil {
		return nil, err
	}

	return fileExts, nil
}

// FileExtensions is the map of file extensions.
var FileExtensions = map[string]string{
	"as":          "ActionScript",
	"ada":         "Ada",
	"adb":         "Ada",
	"ads":         "Ada",
	"alda":        "Alda",
	"Ant":         "Ant",
	"adoc":        "AsciiDoc",
	"asciidoc":    "AsciiDoc",
	"asm":         "Assembly",
	"S":           "Assembly",
	"s":           "Assembly",
	"dats":        "ATS",
	"sats":        "ATS",
	"hats":        "ATS",
	"ahk":         "AutoHotkey",
	"awk":         "Awk",
	"bal":         "Ballerina",
	"bat":         "Batch",
	"btm":         "Batch",
	"bicep":       "Bicep",
	"bb":          "BitBake",
	"be":          "Berry",
	"cairo":       "Cairo",
	"carbon":      "Carbon",
	"cbl":         "COBOL",
	"cmd":         "Batch",
	"bash":        "BASH",
	"sh":          "Bourne Shell",
	"c":           "C",
	"carp":        "Carp",
	"csh":         "C Shell",
	"ec":          "C",
	"erl":         "Erlang",
	"hrl":         "Erlang",
	"pgc":         "C",
	"capnp":       "Cap'n Proto",
	"chpl":        "Chapel",
	"circom":      "Circom",
	"cs":          "C#",
	"clj":         "Clojure",
	"coffee":      "CoffeeScript",
	"cfm":         "ColdFusion",
	"cfc":         "ColdFusion CFScript",
	"cmake":       "CMake",
	"cc":          "C++",
	"cpp":         "C++",
	"cxx":         "C++",
	"pcc":         "C++",
	"c++":         "C++",
	"cr":          "Crystal",
	"css":         "CSS",
	"cu":          "CUDA",
	"d":           "D",
	"dart":        "Dart",
	"dhall":       "Dhall",
	"dtrace":      "DTrace",
	"dts":         "Device Tree",
	"dtsi":        "Device Tree",
	"e":           "Eiffel",
	"elm":         "Elm",
	"el":          "LISP",
	"exp":         "Expect",
	"ex":          "Elixir",
	"exs":         "Elixir",
	"feature":     "Gherkin",
	"factor":      "Factor",
	"fish":        "Fish",
	"fr":          "Frege",
	"fst":         "F*",
	"F#":          "F#",   // deplicated F#/GLSL
	"GLSL":        "GLSL", // both use ext '.fs'
	"vs":          "GLSL",
	"shader":      "HLSL",
	"cg":          "HLSL",
	"cginc":       "HLSL",
	"hlsl":        "HLSL",
	"lean":        "Lean",
	"hlean":       "Lean",
	"lgt":         "Logtalk",
	"lisp":        "LISP",
	"lsp":         "LISP",
	"lua":         "Lua",
	"ls":          "LiveScript",
	"sc":          "LISP",
	"f":           "FORTRAN Legacy",
	"F":           "FORTRAN Legacy",
	"f77":         "FORTRAN Legacy",
	"for":         "FORTRAN Legacy",
	"ftn":         "FORTRAN Legacy",
	"pfo":         "FORTRAN Legacy",
	"f90":         "FORTRAN Modern",
	"F90":         "FORTRAN Modern",
	"f95":         "FORTRAN Modern",
	"f03":         "FORTRAN Modern",
	"f08":         "FORTRAN Modern",
	"gleam":       "Gleam",
	"g4":          "ANTLR",
	"go":          "Go",
	"go2":         "Go",
	"groovy":      "Groovy",
	"gradle":      "Groovy",
	"h":           "C Header",
	"hbs":         "Handlebars",
	"hs":          "Haskell",
	"hpp":         "C++ Header",
	"hh":          "C++ Header",
	"html":        "HTML",
	"ha":          "Hare",
	"hx":          "Haxe",
	"hxx":         "C++ Header",
	"hy":          "Hy",
	"idr":         "Idris",
	"imba":        "Imba",
	"il":          "SKILL",
	"ino":         "Arduino Sketch",
	"io":          "Io",
	"iss":         "Inno Setup",
	"ipynb":       "Jupyter Notebook",
	"jai":         "JAI",
	"java":        "Java",
	"jsp":         "JSP",
	"js":          "JavaScript",
	"jl":          "Julia",
	"janet":       "Janet",
	"json":        "JSON",
	"jsx":         "JSX",
	"just":        "Just",
	"kak":         "KakouneScript",
	"kk":          "Koka",
	"kt":          "Kotlin",
	"kts":         "Kotlin",
	"lds":         "LD Script",
	"less":        "LESS",
	"ly":          "Lilypond",
	"Objective-C": "Objective-C", // deplicated Obj-C/Matlab/Mercury
	"Matlab":      "MATLAB",      // both use ext '.m'
	"Mercury":     "Mercury",     // use ext '.m'
	"md":          "Markdown",
	"markdown":    "Markdown",
	"mo":          "Motoko",
	"Motoko":      "Motoko",
	"ne":          "Nearley",
	"nix":         "Nix",
	"nsi":         "NSIS",
	"nsh":         "NSIS",
	"nu":          "Nu",
	"ML":          "OCaml",
	"ml":          "OCaml",
	"mli":         "OCaml",
	"mll":         "OCaml",
	"mly":         "OCaml",
	"mm":          "Objective-C++",
	"maven":       "Maven",
	"makefile":    "Makefile",
	"meson":       "Meson",
	"mustache":    "Mustache",
	"m4":          "M4",
	"mojo":        "Mojo",
	"🔥":           "Mojo",
	"move":        "Move",
	"l":           "lex",
	"nim":         "Nim",
	"njk":         "Nunjucks",
	"odin":        "Odin",
	"ohm":         "Ohm",
	"php":         "PHP",
	"pas":         "Pascal",
	"PL":          "Perl",
	"pl":          "Perl",
	"pm":          "Perl",
	"plan9sh":     "Plan9 Shell",
	"pony":        "Pony",
	"ps1":         "PowerShell",
	"text":        "Plain Text",
	"txt":         "Plain Text",
	"polly":       "Polly",
	"proto":       "Protocol Buffers",
	"prql":        "PRQL",
	"py":          "Python",
	"pxd":         "Cython",
	"pyx":         "Cython",
	"q":           "Q",
	"qml":         "QML",
	"r":           "R",
	"R":           "R",
	"raml":        "RAML",
	"Rebol":       "Rebol",
	"red":         "Red",
	"rego":        "Rego",
	"Rmd":         "RMarkdown",
	"rake":        "Ruby",
	"rb":          "Ruby",
	"resx":        "XML resource", // ref: https://docs.microsoft.com/en-us/dotnet/framework/resources/creating-resource-files-for-desktop-apps#ResxFiles
	"ring":        "Ring",
	"rkt":         "Racket",
	"rhtml":       "Ruby HTML",
	"rs":          "Rust",
	"rst":         "ReStructuredText",
	"sass":        "Sass",
	"scala":       "Scala",
	"scss":        "Sass",
	"scm":         "Scheme",
	"sed":         "sed",
	"stan":        "Stan",
	"star":        "Starlark",
	"sml":         "Standard ML",
	"sol":         "Solidity",
	"sql":         "SQL",
	"svelte":      "Svelte",
	"swift":       "Swift",
	"t":           "Terra",
	"tex":         "TeX",
	"thy":         "Isabelle",
	"tla":         "TLA",
	"sty":         "TeX",
	"tcl":         "Tcl/Tk",
	"toml":        "TOML",
	"TypeScript":  "TypeScript",
	"tsx":         "TypeScript",
	"tf":          "HCL",
	"um":          "Umka",
	"mat":         "Unity-Prefab",
	"prefab":      "Unity-Prefab",
	"Coq":         "Coq",
	"vala":        "Vala",
	"Verilog":     "Verilog",
	"csproj":      "MSBuild script",
	"vbproj":      "MSBuild script",
	"vcproj":      "MSBuild script",
	"vb":          "Visual Basic",
	"vim":         "VimL",
	"vue":         "Vue",
	"vy":          "Vyper",
	"xml":         "XML",
	"XML":         "XML",
	"xsd":         "XSD",
	"xsl":         "XSLT",
	"xslt":        "XSLT",
	"wxs":         "WiX",
	"yaml":        "YAML",
	"yml":         "YAML",
	"y":           "Yacc",
	"yul":         "Yul",
	"zep":         "Zephir",
	"zig":         "Zig",
	"zsh":         "Zsh",
	"mk":          "Makefile",
}

var shebang2Ext = map[string]string{
	"gosh":    "scm",
	"make":    "make",
	"perl":    "pl",
	"rc":      "plan9sh",
	"python":  "py",
	"ruby":    "rb",
	"escript": "erl",
}

// GetShebang returns the language from the shebang line.
func GetShebang(line string) (shebangLang string, ok bool) {
	ret := reShebangEnv.FindAllStringSubmatch(line, -1)
	if ret != nil && len(ret[0]) == 3 {
		shebangLang = ret[0][2]
		if sl, ok := shebang2Ext[shebangLang]; ok {
			return sl, ok
		}
		return shebangLang, true
	}

	ret = reShebangLang.FindAllStringSubmatch(line, -1)
	if ret != nil && len(ret[0]) >= 2 {
		shebangLang = ret[0][1]
		if sl, ok := shebang2Ext[shebangLang]; ok {
			return sl, ok
		}
		return shebangLang, true
	}

	return "", false
}

// GetFileTypeByShebang returns the language from the shebang line.
func GetFileTypeByShebang(path string) (shebangLang string, ok bool) {
	f, err := os.Open(path)
	if err != nil {
		return // ignore error
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	reader := bufio.NewReader(f)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return
	}
	line = bytes.TrimLeftFunc(line, unicode.IsSpace)

	if len(line) > 2 && line[0] == '#' && line[1] == '!' {
		return GetShebang(string(line))
	}
	return
}

// GetFileType get the file type from the file path.
func GetFileType(path string, opts *option.GClocOptions) (ext string, ok bool) {
	ext = filepath.Ext(path)
	base := filepath.Base(path)

	switch ext {
	case ".m", ".v", ".fs", ".r", ".ts":
		content, err := os.ReadFile(path)
		if err != nil {
			return "", false
		}
		lang := enry.GetLanguage(path, content)
		if opts.Debug {
			fmt.Printf("path=%v, lang=%v\n", path, lang)
		}
		return lang, true
	case ".mo":
		content, err := os.ReadFile(path)
		if err != nil {
			return "", false
		}
		lang := enry.GetLanguage(path, content)
		if opts.Debug {
			fmt.Printf("path=%v, lang=%v\n", path, lang)
		}
		if lang != "" {
			return "Motoko", true
		}
		return lang, true
	}

	switch base {
	case "meson.build", "meson_options.txt":
		return "meson", true
	case "CMakeLists.txt":
		return "cmake", true
	case "configure.ac":
		return "m4", true
	case "Makefile.am":
		return "makefile", true
	case "build.xml":
		return "Ant", true
	case "pom.xml":
		return "maven", true
	}

	switch strings.ToLower(base) {
	case "justfile":
		return "just", true
	case "makefile":
		return "makefile", true
	case "nukefile":
		return "nu", true
	case "rebar": // skip
		return "", false
	}

	shebangLang, ok := GetFileTypeByShebang(path)
	if ok {
		return shebangLang, true
	}

	if len(ext) >= 2 {
		return ext[1:], true
	}

	return ext, ok
}

func lang2Exts(lang string) (exts string) {
	var es []string
	for ext, l := range FileExtensions {
		if lang == l {
			switch lang {
			case "Objective-C", "MATLAB", "Mercury":
				ext = "m"
			case "F#":
				ext = "fs"
			case "GLSL":
				if ext == "GLSL" {
					ext = "fs"
				}
			case "TypeScript":
				ext = "ts"
			case "Motoko":
				ext = "mo"
			}
			es = append(es, ext)
		}
	}
	return strings.Join(es, ", ")
}

// shouldIgnore returns true if the path should be ignored.
func shouldIgnore(path string, info os.FileInfo, vcsInRoot bool, opts *option.GClocOptions) bool {
	if utils.CheckDefaultIgnore(path, info, vcsInRoot) {
		return true
	}
	if !utils.CheckOptionMatch(path, info, opts) {
		return true
	}
	return false
}

// processFile processes the file.
func processFile(path, ext string, languages *DefinedLanguages, opts *option.GClocOptions,
	result *syncmap.SyncMap[string, *Language], fileCache *syncmap.SyncMap[string, struct{}], mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	if targetExt, ok := FileExtensions[ext]; ok {
		if _, ok := opts.ExcludeExts[ext]; ok {
			return
		}

		// exclude languages
		if len(opts.ExcludeLanguages) != 0 {
			if _, ok := opts.ExcludeLanguages[targetExt]; ok {
				return
			}
		}

		// include languages
		if len(opts.IncludeLanguages) != 0 {
			if _, ok := opts.IncludeLanguages[targetExt]; !ok {
				return
			}
		}

		if !opts.SkipDuplicated {
			if utils.CheckMD5Sum(path, fileCache) {
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
func addFileToResult(path, targetExt string, languages *DefinedLanguages,
	result *syncmap.SyncMap[string, *Language]) {
	if ok := result.Has(targetExt); !ok {
		definedLang := NewLanguage(
			languages.Langs[targetExt].Name,
			languages.Langs[targetExt].LineComments,
			languages.Langs[targetExt].MultipleLines,
		)

		if len(languages.Langs[targetExt].RegexLineComments) > 0 {
			definedLang.RegexLineComments = languages.Langs[targetExt].RegexLineComments
		}
		result.Store(targetExt, definedLang)
	}

	lang, _ := result.Load(targetExt)
	lang.Files = append(lang.Files, path)
	result.Store(targetExt, lang)
}

// GetAllFiles return all the files to be analyzed in paths.
func GetAllFiles(paths []string, languages *DefinedLanguages, opts *option.GClocOptions) (map[string]*Language, error) {
	result := syncmap.NewSyncMap[string, *Language](0)
	fileCache := syncmap.NewSyncMap[string, struct{}](0)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, root := range paths {
		vcsInRoot := utils.IsVCSDir(root)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Error("error: %v", err)
				return nil
			}

			if shouldIgnore(path, info, vcsInRoot, opts) {
				return nil
			}

			wg.Add(1)

			go func(p string) {
				defer wg.Done()

				if ext, ok := GetFileType(p, opts); ok {
					processFile(p, ext, languages, opts, result, fileCache, &mu)
				}
			}(path)

			return nil
		})

		if err != nil {
			if opts.Debug {
				log.Error("error: %v", err)
			}
			return nil, err
		}
	}

	wg.Wait()

	return result.ToMap(), nil
}
