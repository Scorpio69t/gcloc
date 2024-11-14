package language

import "testing"

func TestLoadFileExts(t *testing.T) {
	exts, err := loadFileExtsFromJson("../../config/exts.json")
	if err != nil {
		t.Error(err)
	}

	if len(exts) == 0 {
		t.Error("file extensions are empty")
	}

	if _, ok := exts["go"]; !ok {
		t.Error("file extensions do not contain go")
	}

	t.Logf("file extensions count: %d", len(exts))
}

func TestGetShebang(t *testing.T) {
	lang := "py"
	shebang := "#!/usr/bin/env python"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangWithSpace(t *testing.T) {
	lang := "py"
	shebang := "#! /usr/bin/env python"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangBashWithEnv(t *testing.T) {
	lang := "bash"
	shebang := "#!/usr/bin/env bash"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangBash(t *testing.T) {
	lang := "bash"
	shebang := "#!/usr/bin/bash"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangBashWithSpace(t *testing.T) {
	lang := "bash"
	shebang := "#! /usr/bin/bash"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangPlan9Shell(t *testing.T) {
	lang := "plan9sh"
	shebang := "#!/usr/rc"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangStartDot(t *testing.T) {
	lang := "pl"
	shebang := "#!./perl -o"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}

func TestGetShebangMk(t *testing.T) {
	lang := "make"
	shebang := "#!/usr/bin/make -f"

	s, ok := GetShebang(shebang)
	if !ok {
		t.Errorf("invalid logic. shebang=[%v]", shebang)
	}

	if lang != s {
		t.Errorf("invalid logic. lang=[%v] shebang=[%v]", lang, s)
	}

	t.Logf("lang=[%v] shebang=[%v]", lang, s)
}
