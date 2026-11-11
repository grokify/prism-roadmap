package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// resetFlags sets every flag back to its default, since cobra keeps flag
// values between runs in one process.
func resetFlags(t *testing.T, cmd *cobra.Command) {
	t.Helper()
	for _, fs := range []*pflag.FlagSet{cmd.Flags(), cmd.PersistentFlags()} {
		fs.VisitAll(func(f *pflag.Flag) {
			if _, isSlice := f.Value.(pflag.SliceValue); isSlice {
				return // none of the goals commands take slices
			}
			if err := f.Value.Set(f.DefValue); err != nil {
				t.Fatalf("reset --%s: %v", f.Name, err)
			}
			f.Changed = false
		})
	}
	for _, c := range cmd.Commands() {
		resetFlags(t, c)
	}
}

func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	resetFlags(t, RootCmd)
	var out bytes.Buffer
	RootCmd.SetOut(&out)
	RootCmd.SetErr(&out)
	RootCmd.SetArgs(args)
	err := RootCmd.Execute()
	return out.String(), err
}

func TestDetectGoalFramework(t *testing.T) {
	for file, want := range map[string]string{
		"product.vsem.json": "vsem", "product.mel.json": "mel", "manufacturing.dmaic.json": "dmaic",
		"product.v2mom.json": "v2mom", "minimal.v2mom.json": "v2mom",
		"product.ogsm.json": "ogsm", "minimal.ogsm.json": "ogsm",
	} {
		data, err := os.ReadFile(filepath.Join("..", "examples", file))
		if err != nil {
			t.Fatal(err)
		}
		if got, err := DetectGoalFramework(data); err != nil || got != want {
			t.Errorf("%s: framework %q, %v; want %q", file, got, err, want)
		}
	}
	if got, err := DetectGoalFramework([]byte(`{"objectives":[]}`)); err != nil || got != "okr" {
		t.Errorf("objectives: %q, %v", got, err)
	}
	if _, err := DetectGoalFramework([]byte(`{"x":1}`)); err == nil {
		t.Error("an unrecognized document should be refused")
	}
}

func TestGoalsValidate(t *testing.T) {
	out, err := run(t, "goals", "validate", "../examples/product.mel.json", "--strict")
	if err != nil || !strings.Contains(out, "is a valid MEL document") {
		t.Errorf("valid MEL: %v\n%s", err, out)
	}
	bad := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(bad, []byte(`{"strategies":[{"id":"s"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err = run(t, "goals", "validate", bad)
	if err == nil || !strings.Contains(out, "Errors:") || !strings.Contains(out, "strategy name is required") {
		t.Errorf("invalid VSEM should fail with its errors: %v\n%s", err, out)
	}
	if _, err := run(t, "goals", "validate", bad, "--framework", "nope"); err == nil {
		t.Error("an unknown framework should be refused")
	}
}

// Every framework's template, from either init command, validates; the
// VSEM, DMAIC, and MEL templates meet their strict options.
func TestGoalsInitTemplatesValidate(t *testing.T) {
	dir := t.TempDir()
	for _, fw := range GoalFrameworks {
		file := filepath.Join(dir, fw+".json")
		if out, err := run(t, "goals", "init", "--framework", fw, "--team", "Onboarding", "-o", file); err != nil {
			t.Fatalf("init %s: %v\n%s", fw, err, out)
		}
		if out, err := run(t, "goals", "validate", file); err != nil {
			t.Errorf("%s template: %v\n%s", fw, err, out)
		}
	}
	for _, fw := range []string{"vsem", "ogsm", "dmaic", "mel"} {
		file := filepath.Join(dir, "sub-"+fw+".json")
		if out, err := run(t, "goals", fw, "init", "--name", "Plan", "-o", file); err != nil {
			t.Fatalf("%s init: %v\n%s", fw, err, out)
		}
		if out, err := run(t, "goals", fw, "validate", file, "--strict"); err != nil {
			t.Errorf("%s template under strict options: %v\n%s", fw, err, out)
		}
		if _, err := run(t, "goals", fw, "init", "-o", file); err == nil {
			t.Errorf("%s init should not overwrite %s", fw, file)
		}
	}
}

func TestGoalsGenerate(t *testing.T) {
	dir := t.TempDir()
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"goals", "vsem", "generate", "markdown", "../examples/product.vsem.json"}, "## Strategy"},
		{[]string{"goals", "vsem", "generate", "marp", "../examples/product.vsem.json", "--theme", "corporate"}, "marp: true"},
		{[]string{"goals", "ogsm", "generate", "markdown", "../examples/product.ogsm.json"}, "## On One Page"},
		{[]string{"goals", "ogsm", "generate", "marp", "../examples/product.ogsm.json"}, "## The Plan on One Page"},
		{[]string{"goals", "mel", "generate", "markdown", "../examples/product.mel.json"}, "## Indicator Reference Sheets"},
		{[]string{"goals", "mel", "generate", "pandoc", "../examples/product.mel.json", "--summary"}, "toc: true"},
		{[]string{"goals", "dmaic", "generate", "marp", "../examples/manufacturing.dmaic.json"}, "marp: true"},
		{[]string{"goals", "dmaic", "generate", "pandoc", "../examples/manufacturing.dmaic.json"}, "geometry:"},
	} {
		file := filepath.Join(dir, strings.Join(tc.args[1:4], "-")+".md")
		if out, err := run(t, append(tc.args, "-o", file)...); err != nil {
			t.Fatalf("%v: %v\n%s", tc.args, err, out)
		}
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), tc.want) {
			t.Errorf("%v: output lacks %q", tc.args, tc.want)
		}
	}
	out, err := run(t, "goals", "mel", "generate", "markdown", "../examples/product.mel.json", "--summary", "-o", filepath.Join(dir, "summary.md"))
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	data, err := os.ReadFile(filepath.Join(dir, "summary.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "## Indicator Reference Sheets") {
		t.Error("--summary should leave out the reference sheets")
	}
}

func TestSlug(t *testing.T) {
	for in, want := range map[string]string{"FY2027 Product Plan": "fy2027-product-plan", "  A--b  ": "a-b", "MEL": "mel"} {
		if got := slug(in); got != want {
			t.Errorf("slug(%q) = %q, want %q", in, got, want)
		}
	}
}
