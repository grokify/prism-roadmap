package prioritization

import (
	"fmt"
	"sort"
)

// MoSCoW Prioritization Framework
// Reference: https://en.wikipedia.org/wiki/MoSCoW_method
// Method for prioritizing requirements/deliverables.

// MoSCoWPriority represents a MoSCoW prioritization level.
type MoSCoWPriority string

const (
	MoSCoWMustHave   MoSCoWPriority = "must_have"   // Critical, non-negotiable
	MoSCoWShouldHave MoSCoWPriority = "should_have" // Important but not critical
	MoSCoWCouldHave  MoSCoWPriority = "could_have"  // Nice to have
	MoSCoWWontHave   MoSCoWPriority = "wont_have"   // Explicitly out of scope

	// MoSCoWUnspecified is the zero value: not yet prioritized (e.g. items
	// imported from an external PM tool before triage). It is not a valid
	// value for IsValidMoSCoWPriority — types that allow it (such as
	// rmi.RoadmapItem) treat it as "unset" rather than as a priority level.
	MoSCoWUnspecified MoSCoWPriority = ""
)

// String returns the string representation.
func (m MoSCoWPriority) String() string {
	return string(m)
}

// Weight returns the numeric weight for prioritization calculations.
// Higher weight = higher priority.
func (m MoSCoWPriority) Weight() int {
	switch m {
	case MoSCoWMustHave:
		return 4
	case MoSCoWShouldHave:
		return 3
	case MoSCoWCouldHave:
		return 2
	case MoSCoWWontHave:
		return 0
	default:
		return 0
	}
}

// Description returns a human-readable description of the priority level.
func (m MoSCoWPriority) Description() string {
	switch m {
	case MoSCoWMustHave:
		return "Critical requirement - must be delivered"
	case MoSCoWShouldHave:
		return "Important but not critical - should be included if possible"
	case MoSCoWCouldHave:
		return "Desirable but not necessary - include if time/resources allow"
	case MoSCoWWontHave:
		return "Explicitly out of scope - will not be delivered in this timeframe"
	default:
		return "Unknown priority level"
	}
}

// IsActionable returns true if items at this priority require action.
func (m MoSCoWPriority) IsActionable() bool {
	return m != MoSCoWWontHave && m != MoSCoWUnspecified
}

// MoSCoWItem represents an item with MoSCoW prioritization.
type MoSCoWItem struct {
	// Item identification
	ItemID      string `json:"itemId"`             // Unique identifier
	ItemName    string `json:"itemName,omitempty"` // Human-readable name
	Description string `json:"description,omitempty"`

	// Priority
	Priority      MoSCoWPriority `json:"priority"`
	Justification string         `json:"justification,omitempty"` // Why this priority level

	// Metadata
	AssignedBy   string `json:"assignedBy,omitempty"`
	AssignedDate string `json:"assignedDate,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// IsComplete returns true if required fields are set.
func (m *MoSCoWItem) IsComplete() bool {
	return m.ItemID != "" && m.Priority != ""
}

// Validate returns an error if the item is invalid.
func (m *MoSCoWItem) Validate() error {
	if m.ItemID == "" {
		return fmt.Errorf("itemId is required")
	}
	if m.Priority == "" {
		return fmt.Errorf("priority is required")
	}
	if !IsValidMoSCoWPriority(m.Priority) {
		return fmt.Errorf("invalid priority: %s", m.Priority)
	}
	return nil
}

// MoSCoWSet represents a collection of MoSCoW-prioritized items.
type MoSCoWSet struct {
	Items       []MoSCoWItem `json:"items"`
	Description string       `json:"description,omitempty"`
	AssignedBy  string       `json:"assignedBy,omitempty"`
	CreatedDate string       `json:"createdDate,omitempty"`
}

// NewMoSCoWSet creates a new MoSCoW set.
func NewMoSCoWSet() *MoSCoWSet {
	return &MoSCoWSet{
		Items: []MoSCoWItem{},
	}
}

// Add adds an item to the set.
func (s *MoSCoWSet) Add(item MoSCoWItem) {
	s.Items = append(s.Items, item)
}

// GetByPriority returns items matching the given priority.
func (s *MoSCoWSet) GetByPriority(priority MoSCoWPriority) []MoSCoWItem {
	var result []MoSCoWItem
	for _, item := range s.Items {
		if item.Priority == priority {
			result = append(result, item)
		}
	}
	return result
}

// MustHaves returns all Must Have items.
func (s *MoSCoWSet) MustHaves() []MoSCoWItem {
	return s.GetByPriority(MoSCoWMustHave)
}

// ShouldHaves returns all Should Have items.
func (s *MoSCoWSet) ShouldHaves() []MoSCoWItem {
	return s.GetByPriority(MoSCoWShouldHave)
}

// CouldHaves returns all Could Have items.
func (s *MoSCoWSet) CouldHaves() []MoSCoWItem {
	return s.GetByPriority(MoSCoWCouldHave)
}

// WontHaves returns all Won't Have items.
func (s *MoSCoWSet) WontHaves() []MoSCoWItem {
	return s.GetByPriority(MoSCoWWontHave)
}

// Actionable returns all actionable items (Must, Should, Could).
func (s *MoSCoWSet) Actionable() []MoSCoWItem {
	var result []MoSCoWItem
	for _, item := range s.Items {
		if item.Priority.IsActionable() {
			result = append(result, item)
		}
	}
	return result
}

// SortByPriority sorts items by priority weight (Must first, then Should, etc.).
func (s *MoSCoWSet) SortByPriority() {
	sort.Slice(s.Items, func(i, j int) bool {
		return s.Items[i].Priority.Weight() > s.Items[j].Priority.Weight()
	})
}

// GetByID returns an item by ID.
func (s *MoSCoWSet) GetByID(itemID string) *MoSCoWItem {
	for i := range s.Items {
		if s.Items[i].ItemID == itemID {
			return &s.Items[i]
		}
	}
	return nil
}

// Summary returns a count of items by priority.
func (s *MoSCoWSet) Summary() map[MoSCoWPriority]int {
	counts := make(map[MoSCoWPriority]int)
	for _, item := range s.Items {
		counts[item.Priority]++
	}
	return counts
}

// Validate returns an error if the set is invalid.
func (s *MoSCoWSet) Validate() error {
	for i, item := range s.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i, err)
		}
	}
	return nil
}

// NewMoSCoWItem creates a new MoSCoW item with the given values.
func NewMoSCoWItem(itemID string, priority MoSCoWPriority) *MoSCoWItem {
	return &MoSCoWItem{
		ItemID:   itemID,
		Priority: priority,
	}
}

// ValidMoSCoWPriorities returns all valid MoSCoW priority levels.
func ValidMoSCoWPriorities() []MoSCoWPriority {
	return []MoSCoWPriority{
		MoSCoWMustHave,
		MoSCoWShouldHave,
		MoSCoWCouldHave,
		MoSCoWWontHave,
	}
}

// AllMoSCoWPriorities is an alias for ValidMoSCoWPriorities for consistency.
func AllMoSCoWPriorities() []MoSCoWPriority {
	return ValidMoSCoWPriorities()
}

// IsValidMoSCoWPriority returns true if the priority is valid.
func IsValidMoSCoWPriority(p MoSCoWPriority) bool {
	for _, valid := range ValidMoSCoWPriorities() {
		if p == valid {
			return true
		}
	}
	return false
}

// ParseMoSCoWPriority parses a string to MoSCoWPriority.
// Accepts various formats: "must_have", "must", "M", etc.
func ParseMoSCoWPriority(s string) (MoSCoWPriority, error) {
	switch s {
	case "must_have", "must", "Must", "MUST", "M":
		return MoSCoWMustHave, nil
	case "should_have", "should", "Should", "SHOULD", "S":
		return MoSCoWShouldHave, nil
	case "could_have", "could", "Could", "COULD", "C":
		return MoSCoWCouldHave, nil
	case "wont_have", "wont", "Wont", "WONT", "won't", "Won't", "W":
		return MoSCoWWontHave, nil
	default:
		return "", fmt.Errorf("invalid MoSCoW priority: %s", s)
	}
}
