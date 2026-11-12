package mdshift

import "testing"

func TestShift(t *testing.T) {
	for _, tt := range []struct {
		name string
		md   string
		n    int
		want string
	}{
		{"zero leaves it unchanged", "# A\n\n## B\n", 0, "# A\n\n## B\n"},
		{"negative leaves it unchanged", "# A\n", -2, "# A\n"},
		{"empty", "", 2, ""},
		{"one level", "# A\n\ntext\n\n## B\n", 1, "## A\n\ntext\n\n### B\n"},
		{"two levels", "## A\n### B\n", 2, "#### A\n##### B\n"},
		{"capped at six", "#### A\n##### B\n###### C\n", 3, "###### A\n###### B\n###### C\n"},
		{"empty heading", "#\n", 1, "##\n"},
		{"indented heading", "   # A\n", 1, "   ## A\n"},
		{"tab after hashes", "#\tA\n", 1, "##\tA\n"},
		{"not a heading without a space", "#hashtag\n", 1, "#hashtag\n"},
		{"seven hashes is not a heading", "####### A\n", 1, "####### A\n"},
		{"code indent is not a heading", "    # A\n", 1, "    # A\n"},
		{"inline hash", "text # not a heading\n", 1, "text # not a heading\n"},
		{"backtick fence", "# A\n```sh\n# comment\n```\n# B\n", 1, "## A\n```sh\n# comment\n```\n## B\n"},
		{"tilde fence", "~~~\n# comment\n~~~\n# B\n", 1, "~~~\n# comment\n~~~\n## B\n"},
		{"shorter fence does not close", "````\n```\n# in code\n````\n# B\n", 1, "````\n```\n# in code\n````\n## B\n"},
		{"other fence char does not close", "```\n~~~\n# in code\n```\n# B\n", 1, "```\n~~~\n# in code\n```\n## B\n"},
		{"unclosed fence runs to the end", "# A\n```\n# in code\n", 1, "## A\n```\n# in code\n"},
		{"inline code span is not a fence", "``` a ` b\n# A\n", 1, "``` a ` b\n## A\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Shift(tt.md, tt.n); got != tt.want {
				t.Errorf("Shift(%q, %d) = %q, want %q", tt.md, tt.n, got, tt.want)
			}
		})
	}
}
