package bsc

import (
	"path/filepath"
	"testing"
)

// The examples in the repository load and pass validation, so they stay in
// step with the types.
func TestExamplesValidate(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.bsc.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no Balanced Scorecard examples found")
	}
	for _, f := range files {
		b, err := ReadFile(f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if errs := Errors(b.Validate(DefaultValidationOptions())); len(errs) > 0 {
			t.Errorf("%s: %v", f, errs)
		}
	}
}

// The full example meets the strict options, so it shows a complete Balanced Scorecard.
func TestProductExampleIsStrict(t *testing.T) {
	b, err := ReadFile("../../examples/product.bsc.json")
	if err != nil {
		t.Fatal(err)
	}
	if errs := b.Validate(StrictValidationOptions()); len(errs) > 0 {
		t.Errorf("product.bsc.json under strict options: %v", errs)
	}
}
