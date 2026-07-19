package templates

import (
	"slices"
	"testing"
)

func TestGetCanonicalTemplates(t *testing.T) {
	for _, name := range []string{"bmc", "opportunity-spec"} {
		content, err := Get(name)
		if err != nil {
			t.Errorf("Get(%q): %v", name, err)
			continue
		}
		if len(content) == 0 {
			t.Errorf("Get(%q) returned empty content", name)
		}
	}
}

func TestGetAcceptsSuffix(t *testing.T) {
	withSuffix, err := Get("bmc.md")
	if err != nil {
		t.Fatalf("Get with suffix: %v", err)
	}
	withoutSuffix, err := Get("bmc")
	if err != nil {
		t.Fatalf("Get without suffix: %v", err)
	}
	if withSuffix != withoutSuffix {
		t.Error("suffix and non-suffix forms returned different content")
	}
}

func TestNamesIncludesBMC(t *testing.T) {
	names := Names()
	if !slices.Contains(names, "bmc") {
		t.Errorf("Names() = %v, want to include bmc", names)
	}
}
