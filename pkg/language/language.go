package language

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gcloc/pkg/option"
	"gcloc/pkg/utils"
	"github.com/go-enry/go-enry/v2"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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
