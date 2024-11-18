package json

import (
	"encoding/json"
	"fmt"
	"github.com/Scorpio69t/gcloc/pkg/file"
	"github.com/Scorpio69t/gcloc/pkg/language"
	"testing"
)

func TestOutputJSON(t *testing.T) {
	total := &language.Language{}
	files := []file.GClocFile{
		{Name: "one.go", Language: "Go"},
		{Name: "two.go", Language: "Go"},
	}

	jsonResult := NewFilesResultFromGCloc(total, files)
	if jsonResult.Files[0].Name != "one.go" {
		t.Errorf("invalid result. Name: one.go")
	}

	if jsonResult.Files[1].Name != "two.go" {
		t.Errorf("invalid result. Name: two.go")
	}

	if jsonResult.Files[1].Language != "Go" {
		t.Errorf("invalid result. lang: Go")
	}

	// check output json text
	buf, err := json.Marshal(jsonResult)
	if err != nil {
		fmt.Println(err)
		t.Errorf("json marshal error")
	}

	actualJSONText := `{"files":[{"name":"one.go","language":"Go","codes":0,"comments":0,"blanks":0},{"name":"two.go","language":"Go","codes":0,"comments":0,"blanks":0}],"total":{"file_count":0,"codes":0,"comments":0,"blanks":0}}`
	resultJSONText := string(buf)
	if actualJSONText != resultJSONText {
		t.Errorf("invalid result. '%s'", resultJSONText)
	}
}
