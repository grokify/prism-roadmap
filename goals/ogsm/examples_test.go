package ogsm

import (
	"path/filepath"
	"testing"
)

// The examples in the repository load and pass validation, so they stay in
// step with the types.
func TestExamplesValidate(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.ogsm.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no OGSM examples found")
	}
	for _, f := range files {
		o, err := ReadFile(f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if errs := Errors(o.Validate(DefaultValidationOptions())); len(errs) > 0 {
			t.Errorf("%s: %v", f, errs)
		}
	}
}

// The full example meets the strict options, so it shows a complete OGSM.
func TestProductExampleIsStrict(t *testing.T) {
	o, err := ReadFile("../../examples/product.ogsm.json")
	if err != nil {
		t.Fatal(err)
	}
	if errs := o.Validate(StrictValidationOptions()); len(errs) > 0 {
		t.Errorf("product.ogsm.json under strict options: %v", errs)
	}
}
