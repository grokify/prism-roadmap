package mel

import (
	"path/filepath"
	"testing"
)

// The examples in the repository load and pass validation, so they stay in
// step with the types.
func TestExamplesValidate(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.mel.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no MEL examples found")
	}
	for _, f := range files {
		p, err := ReadFile(f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if errs := Errors(p.Validate(DefaultValidationOptions())); len(errs) > 0 {
			t.Errorf("%s: %v", f, errs)
		}
	}
}

// The full example meets the strict options, so it shows a plan ready to run.
func TestProductExampleIsStrict(t *testing.T) {
	p, err := ReadFile("../../examples/product.mel.json")
	if err != nil {
		t.Fatal(err)
	}
	if errs := p.Validate(StrictValidationOptions()); len(errs) > 0 {
		t.Errorf("product.mel.json under strict options: %v", errs)
	}
}
