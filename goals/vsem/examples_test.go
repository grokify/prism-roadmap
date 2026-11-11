package vsem

import (
	"path/filepath"
	"testing"
)

// The examples in the repository load and pass validation, so they stay in
// step with the types.
func TestExamplesValidate(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.vsem.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no VSEM examples found")
	}
	for _, f := range files {
		v, err := ReadFile(f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if errs := Errors(v.Validate(DefaultValidationOptions())); len(errs) > 0 {
			t.Errorf("%s: %v", f, errs)
		}
	}
}

// The full example meets the strict options, so it shows a complete VSEM.
func TestProductExampleIsStrict(t *testing.T) {
	v, err := ReadFile("../../examples/product.vsem.json")
	if err != nil {
		t.Fatal(err)
	}
	if errs := v.Validate(StrictValidationOptions()); len(errs) > 0 {
		t.Errorf("product.vsem.json under strict options: %v", errs)
	}
}
