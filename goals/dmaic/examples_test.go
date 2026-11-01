package dmaic

import (
	"path/filepath"
	"testing"
)

// The examples in the repository load and pass validation, so they stay in
// step with the types.
func TestExamplesValidate(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.dmaic.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no DMAIC examples found")
	}
	for _, f := range files {
		doc, err := ReadFile(f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if errs := Errors(doc.Validate(DefaultValidationOptions())); len(errs) > 0 {
			t.Errorf("%s: %v", f, errs)
		}
	}
}
