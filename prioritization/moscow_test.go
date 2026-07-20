package prioritization

import "testing"

func TestMoSCoWPriorityWeight(t *testing.T) {
	tests := []struct {
		priority MoSCoWPriority
		want     int
	}{
		{MoSCoWMustHave, 4},
		{MoSCoWShouldHave, 3},
		{MoSCoWCouldHave, 2},
		{MoSCoWWontHave, 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		got := tt.priority.Weight()
		if got != tt.want {
			t.Errorf("MoSCoWPriority(%s).Weight() = %d, want %d", tt.priority, got, tt.want)
		}
	}
}

func TestMoSCoWPriorityIsActionable(t *testing.T) {
	tests := []struct {
		priority MoSCoWPriority
		want     bool
	}{
		{MoSCoWMustHave, true},
		{MoSCoWShouldHave, true},
		{MoSCoWCouldHave, true},
		{MoSCoWWontHave, false},
		{"", false},
	}

	for _, tt := range tests {
		got := tt.priority.IsActionable()
		if got != tt.want {
			t.Errorf("MoSCoWPriority(%s).IsActionable() = %v, want %v", tt.priority, got, tt.want)
		}
	}
}

func TestMoSCoWItemValidate(t *testing.T) {
	tests := []struct {
		name    string
		item    MoSCoWItem
		wantErr bool
	}{
		{
			name: "valid item",
			item: MoSCoWItem{
				ItemID:   "item-1",
				Priority: MoSCoWMustHave,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			item: MoSCoWItem{
				Priority: MoSCoWMustHave,
			},
			wantErr: true,
		},
		{
			name: "missing priority",
			item: MoSCoWItem{
				ItemID: "item-1",
			},
			wantErr: true,
		},
		{
			name: "invalid priority",
			item: MoSCoWItem{
				ItemID:   "item-1",
				Priority: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MoSCoWItem.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMoSCoWSetSortByPriority(t *testing.T) {
	set := NewMoSCoWSet()
	set.Add(MoSCoWItem{ItemID: "could", Priority: MoSCoWCouldHave})
	set.Add(MoSCoWItem{ItemID: "must", Priority: MoSCoWMustHave})
	set.Add(MoSCoWItem{ItemID: "should", Priority: MoSCoWShouldHave})
	set.Add(MoSCoWItem{ItemID: "wont", Priority: MoSCoWWontHave})

	set.SortByPriority()

	expected := []string{"must", "should", "could", "wont"}
	for i, item := range set.Items {
		if item.ItemID != expected[i] {
			t.Errorf("After sort, position %d = %s, want %s", i, item.ItemID, expected[i])
		}
	}
}

func TestMoSCoWSetMustHaves(t *testing.T) {
	set := NewMoSCoWSet()
	set.Add(MoSCoWItem{ItemID: "could-1", Priority: MoSCoWCouldHave})
	set.Add(MoSCoWItem{ItemID: "must-1", Priority: MoSCoWMustHave})
	set.Add(MoSCoWItem{ItemID: "must-2", Priority: MoSCoWMustHave})
	set.Add(MoSCoWItem{ItemID: "should-1", Priority: MoSCoWShouldHave})

	musts := set.MustHaves()
	if len(musts) != 2 {
		t.Errorf("MustHaves() returned %d items, want 2", len(musts))
	}
}

func TestMoSCoWSetSummary(t *testing.T) {
	set := NewMoSCoWSet()
	set.Add(MoSCoWItem{ItemID: "must-1", Priority: MoSCoWMustHave})
	set.Add(MoSCoWItem{ItemID: "must-2", Priority: MoSCoWMustHave})
	set.Add(MoSCoWItem{ItemID: "should-1", Priority: MoSCoWShouldHave})
	set.Add(MoSCoWItem{ItemID: "could-1", Priority: MoSCoWCouldHave})

	summary := set.Summary()

	if summary[MoSCoWMustHave] != 2 {
		t.Errorf("Summary[MustHave] = %d, want 2", summary[MoSCoWMustHave])
	}
	if summary[MoSCoWShouldHave] != 1 {
		t.Errorf("Summary[ShouldHave] = %d, want 1", summary[MoSCoWShouldHave])
	}
	if summary[MoSCoWCouldHave] != 1 {
		t.Errorf("Summary[CouldHave] = %d, want 1", summary[MoSCoWCouldHave])
	}
}

func TestParseMoSCoWPriority(t *testing.T) {
	tests := []struct {
		input   string
		want    MoSCoWPriority
		wantErr bool
	}{
		{"must_have", MoSCoWMustHave, false},
		{"must", MoSCoWMustHave, false},
		{"Must", MoSCoWMustHave, false},
		{"MUST", MoSCoWMustHave, false},
		{"M", MoSCoWMustHave, false},
		{"should_have", MoSCoWShouldHave, false},
		{"should", MoSCoWShouldHave, false},
		{"S", MoSCoWShouldHave, false},
		{"could_have", MoSCoWCouldHave, false},
		{"could", MoSCoWCouldHave, false},
		{"C", MoSCoWCouldHave, false},
		{"wont_have", MoSCoWWontHave, false},
		{"won't", MoSCoWWontHave, false},
		{"W", MoSCoWWontHave, false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		got, err := ParseMoSCoWPriority(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseMoSCoWPriority(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseMoSCoWPriority(%q) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestValidMoSCoWPriorities(t *testing.T) {
	priorities := ValidMoSCoWPriorities()
	if len(priorities) != 4 {
		t.Errorf("ValidMoSCoWPriorities() returned %d priorities, want 4", len(priorities))
	}
}

func TestIsValidMoSCoWPriority(t *testing.T) {
	if !IsValidMoSCoWPriority(MoSCoWMustHave) {
		t.Error("IsValidMoSCoWPriority(MustHave) = false, want true")
	}
	if IsValidMoSCoWPriority("invalid") {
		t.Error("IsValidMoSCoWPriority(invalid) = true, want false")
	}
}
