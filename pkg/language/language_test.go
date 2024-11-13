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
