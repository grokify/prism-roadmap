package rmi

import (
	"fmt"
	"sort"
	"time"

	"github.com/grokify/prism-roadmap/prioritization"
)

// RoadmapItemSet represents a collection of roadmap items.
type RoadmapItemSet struct {
	Items       []RoadmapItem `json:"items"`
	Description string        `json:"description,omitempty"`
	Quarter     string        `json:"quarter,omitempty"` // e.g., "Q3 2026"
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}

// NewRoadmapItemSet creates a new set.
func NewRoadmapItemSet() *RoadmapItemSet {
	now := time.Now()
	return &RoadmapItemSet{
		Items:     []RoadmapItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Add adds an item to the set.
func (s *RoadmapItemSet) Add(item RoadmapItem) {
	s.Items = append(s.Items, item)
	s.UpdatedAt = time.Now()
}

// GetByID returns an item by ID.
func (s *RoadmapItemSet) GetByID(id string) *RoadmapItem {
	for i := range s.Items {
		if s.Items[i].ID == id {
			return &s.Items[i]
		}
	}
	return nil
}

// GetByMoSCoW returns items matching the given MoSCoW priority.
func (s *RoadmapItemSet) GetByMoSCoW(moscow prioritization.MoSCoWPriority) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		if item.MoSCoW == moscow {
			result = append(result, item)
		}
	}
	return result
}

// GetByStatus returns items matching the given status.
func (s *RoadmapItemSet) GetByStatus(status RMIStatus) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		if item.Status == status {
			result = append(result, item)
		}
	}
	return result
}

// MustHaves returns all Must Have items.
func (s *RoadmapItemSet) MustHaves() []RoadmapItem {
	return s.GetByMoSCoW(prioritization.MoSCoWMustHave)
}

// ShouldHaves returns all Should Have items.
func (s *RoadmapItemSet) ShouldHaves() []RoadmapItem {
	return s.GetByMoSCoW(prioritization.MoSCoWShouldHave)
}

// CouldHaves returns all Could Have items.
func (s *RoadmapItemSet) CouldHaves() []RoadmapItem {
	return s.GetByMoSCoW(prioritization.MoSCoWCouldHave)
}

// Actionable returns all actionable items.
func (s *RoadmapItemSet) Actionable() []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		if item.IsActionable() {
			result = append(result, item)
		}
	}
	return result
}

// Blocked returns items with blocking dependencies.
func (s *RoadmapItemSet) Blocked() []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		if item.HasBlockingDependencies() {
			result = append(result, item)
		}
	}
	return result
}

// SortByPriority sorts items by MoSCoW priority (Must first).
func (s *RoadmapItemSet) SortByPriority() {
	sort.Slice(s.Items, func(i, j int) bool {
		return s.Items[i].MoSCoW.Weight() > s.Items[j].MoSCoW.Weight()
	})
}

// SortByPriorityScore sorts items by composite priority score (highest first).
func (s *RoadmapItemSet) SortByPriorityScore() {
	sort.Slice(s.Items, func(i, j int) bool {
		return s.Items[i].PriorityScore() > s.Items[j].PriorityScore()
	})
}

// SortByRICE sorts items by RICE score (highest first).
// Items without RICE scores are placed at the end.
func (s *RoadmapItemSet) SortByRICE() {
	sort.Slice(s.Items, func(i, j int) bool {
		scoreI := 0.0
		scoreJ := 0.0
		if s.Items[i].RICE != nil {
			scoreI = s.Items[i].RICE.Score
		}
		if s.Items[j].RICE != nil {
			scoreJ = s.Items[j].RICE.Score
		}
		return scoreI > scoreJ
	})
}

// SortByMarketSignal sorts items by market signal score (highest first).
func (s *RoadmapItemSet) SortByMarketSignal() {
	sort.Slice(s.Items, func(i, j int) bool {
		scoreI := 0.0
		scoreJ := 0.0
		if s.Items[i].MarketSignal != nil {
			scoreI = s.Items[i].MarketSignal.Score
		}
		if s.Items[j].MarketSignal != nil {
			scoreJ = s.Items[j].MarketSignal.Score
		}
		return scoreI > scoreJ
	})
}

// TopN returns the top N items by priority score.
func (s *RoadmapItemSet) TopN(n int) []RoadmapItem {
	s.SortByPriorityScore()
	if n > len(s.Items) {
		n = len(s.Items)
	}
	return s.Items[:n]
}

// TotalEffortDays returns the total effort in person-days for all items.
func (s *RoadmapItemSet) TotalEffortDays() int {
	total := 0
	for _, item := range s.Items {
		total += item.EffectiveEffortDays()
	}
	return total
}

// MoSCoWSummary returns a count of items by MoSCoW priority.
func (s *RoadmapItemSet) MoSCoWSummary() map[prioritization.MoSCoWPriority]int {
	counts := make(map[prioritization.MoSCoWPriority]int)
	for _, item := range s.Items {
		counts[item.MoSCoW]++
	}
	return counts
}

// StatusSummary returns a count of items by status.
func (s *RoadmapItemSet) StatusSummary() map[RMIStatus]int {
	counts := make(map[RMIStatus]int)
	for _, item := range s.Items {
		counts[item.Status]++
	}
	return counts
}

// Validate returns an error if the set is invalid.
func (s *RoadmapItemSet) Validate() error {
	seenIDs := make(map[string]bool)
	for i, item := range s.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d (%s): %w", i, item.ID, err)
		}
		if seenIDs[item.ID] {
			return fmt.Errorf("duplicate item ID: %s", item.ID)
		}
		seenIDs[item.ID] = true
	}
	return nil
}

// FilterByCapability returns items that reference the given capability ID.
func (s *RoadmapItemSet) FilterByCapability(capID string) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		for _, ref := range item.CapabilityRefs {
			if ref == capID {
				result = append(result, item)
				break
			}
		}
	}
	return result
}

// FilterByGoal returns items that reference the given goal ID.
func (s *RoadmapItemSet) FilterByGoal(goalID string) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		for _, ref := range item.GoalRefs {
			if ref == goalID {
				result = append(result, item)
				break
			}
		}
	}
	return result
}

// FilterByIdea returns items that reference the given idea ID.
func (s *RoadmapItemSet) FilterByIdea(ideaID string) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		for _, ref := range item.IdeaRefs {
			if ref == ideaID {
				result = append(result, item)
				break
			}
		}
	}
	return result
}

// FilterByQuarter returns items for the given quarter.
func (s *RoadmapItemSet) FilterByQuarter(quarter string) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		if item.Quarter == quarter {
			result = append(result, item)
		}
	}
	return result
}

// FilterByTag returns items that have the given tag.
func (s *RoadmapItemSet) FilterByTag(tag string) []RoadmapItem {
	var result []RoadmapItem
	for _, item := range s.Items {
		for _, t := range item.Tags {
			if t == tag {
				result = append(result, item)
				break
			}
		}
	}
	return result
}
