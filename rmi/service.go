// Package rmi provides service operations for roadmap item management.
package rmi

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/grokify/prism-roadmap/prioritization"
)

// ReadFile loads a RoadmapItemSet from a JSON file.
func ReadFile(path string) (*RoadmapItemSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	var set RoadmapItemSet
	if err := json.Unmarshal(data, &set); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &set, nil
}

// WriteFile saves a RoadmapItemSet to a JSON file.
func (s *RoadmapItemSet) WriteFile(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// CreateInput holds parameters for creating a new roadmap item.
type CreateInput struct {
	ID          string
	Name        string
	Description string
	MoSCoW      prioritization.MoSCoWPriority
	Quarter     string
	Owner       string
	Tags        []string
}

// UpdateInput holds optional parameters for updating a roadmap item.
type UpdateInput struct {
	Name        *string
	Description *string
	MoSCoW      *prioritization.MoSCoWPriority
	Status      *RMIStatus
	Quarter     *string
	Owner       *string
	Progress    *int
	Tags        []string // nil means no change, empty slice clears
}

// ListFilter defines criteria for filtering roadmap items.
type ListFilter struct {
	Status   RMIStatus
	MoSCoW   prioritization.MoSCoWPriority
	Quarter  string
	Tag      string
	Owner    string
	Limit    int
	SortBy   string // "priority", "rice", "market_signal", "created"
	SortDesc bool
}

// Service provides roadmap item management operations.
type Service struct {
	set  *RoadmapItemSet
	path string
}

// NewService creates a new service with an empty set.
func NewService() *Service {
	return &Service{
		set: NewRoadmapItemSet(),
	}
}

// NewServiceFromFile creates a service by loading from a file.
func NewServiceFromFile(path string) (*Service, error) {
	set, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Service{
		set:  set,
		path: path,
	}, nil
}

// Set returns the underlying RoadmapItemSet.
func (s *Service) Set() *RoadmapItemSet {
	return s.set
}

// Save persists the set to its file (if loaded from file).
func (s *Service) Save() error {
	if s.path == "" {
		return fmt.Errorf("no file path set")
	}
	return s.set.WriteFile(s.path)
}

// SaveAs persists the set to a new file path.
func (s *Service) SaveAs(path string) error {
	if err := s.set.WriteFile(path); err != nil {
		return err
	}
	s.path = path
	return nil
}

// Create adds a new roadmap item.
func (s *Service) Create(input CreateInput) (*RoadmapItem, error) {
	// Check for duplicate ID
	if existing := s.set.GetByID(input.ID); existing != nil {
		return nil, fmt.Errorf("item %s already exists", input.ID)
	}

	if input.MoSCoW == "" {
		input.MoSCoW = prioritization.MoSCoWShouldHave // default
	}

	item := NewRoadmapItem(input.ID, input.Name, input.MoSCoW)
	item.Description = input.Description
	item.Quarter = input.Quarter
	item.Owner = input.Owner
	item.Tags = input.Tags

	s.set.Add(*item)
	s.set.UpdatedAt = time.Now()

	return item, nil
}

// Get retrieves an item by ID.
func (s *Service) Get(id string) (*RoadmapItem, error) {
	item := s.set.GetByID(id)
	if item == nil {
		return nil, fmt.Errorf("item %s not found", id)
	}
	return item, nil
}

// List returns items matching the filter.
func (s *Service) List(filter ListFilter) []RoadmapItem {
	result := make([]RoadmapItem, 0, len(s.set.Items))

	for _, item := range s.set.Items {
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.MoSCoW != "" && item.MoSCoW != filter.MoSCoW {
			continue
		}
		if filter.Quarter != "" && item.Quarter != filter.Quarter {
			continue
		}
		if filter.Owner != "" && item.Owner != filter.Owner {
			continue
		}
		if filter.Tag != "" {
			hasTag := false
			for _, t := range item.Tags {
				if t == filter.Tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		result = append(result, item)
	}

	// Sort
	switch filter.SortBy {
	case "priority":
		sort.Slice(result, func(i, j int) bool {
			if filter.SortDesc {
				return result[i].PriorityScore() < result[j].PriorityScore()
			}
			return result[i].PriorityScore() > result[j].PriorityScore()
		})
	case "rice":
		sort.Slice(result, func(i, j int) bool {
			scoreI := 0.0
			scoreJ := 0.0
			if result[i].RICE != nil {
				scoreI = result[i].RICE.Score
			}
			if result[j].RICE != nil {
				scoreJ = result[j].RICE.Score
			}
			if filter.SortDesc {
				return scoreI < scoreJ
			}
			return scoreI > scoreJ
		})
	case "market_signal":
		sort.Slice(result, func(i, j int) bool {
			scoreI := 0.0
			scoreJ := 0.0
			if result[i].MarketSignal != nil {
				scoreI = result[i].MarketSignal.Score
			}
			if result[j].MarketSignal != nil {
				scoreJ = result[j].MarketSignal.Score
			}
			if filter.SortDesc {
				return scoreI < scoreJ
			}
			return scoreI > scoreJ
		})
	case "created":
		sort.Slice(result, func(i, j int) bool {
			if filter.SortDesc {
				return result[i].CreatedAt.After(result[j].CreatedAt)
			}
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		})
	default:
		// Default: sort by MoSCoW weight
		sort.Slice(result, func(i, j int) bool {
			return result[i].MoSCoW.Weight() > result[j].MoSCoW.Weight()
		})
	}

	// Apply limit
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result
}

// Update modifies an existing item.
func (s *Service) Update(id string, input UpdateInput) (*RoadmapItem, bool, error) {
	item := s.set.GetByID(id)
	if item == nil {
		return nil, false, fmt.Errorf("item %s not found", id)
	}

	updated := false

	if input.Name != nil && *input.Name != item.Name {
		item.Name = *input.Name
		updated = true
	}
	if input.Description != nil && *input.Description != item.Description {
		item.Description = *input.Description
		updated = true
	}
	if input.MoSCoW != nil && *input.MoSCoW != item.MoSCoW {
		if !prioritization.IsValidMoSCoWPriority(*input.MoSCoW) {
			return nil, false, fmt.Errorf("invalid moscow priority: %s", *input.MoSCoW)
		}
		item.MoSCoW = *input.MoSCoW
		updated = true
	}
	if input.Status != nil && *input.Status != item.Status {
		item.Status = *input.Status
		updated = true
	}
	if input.Quarter != nil && *input.Quarter != item.Quarter {
		item.Quarter = *input.Quarter
		updated = true
	}
	if input.Owner != nil && *input.Owner != item.Owner {
		item.Owner = *input.Owner
		updated = true
	}
	if input.Progress != nil && (item.Progress == nil || *input.Progress != *item.Progress) {
		item.Progress = input.Progress
		updated = true
	}
	if input.Tags != nil {
		item.Tags = input.Tags
		updated = true
	}

	if updated {
		item.UpdatedAt = time.Now()
		s.set.UpdatedAt = time.Now()
	}

	return item, updated, nil
}

// Delete removes an item by ID.
func (s *Service) Delete(id string) error {
	newItems := make([]RoadmapItem, 0, len(s.set.Items))
	found := false
	for _, item := range s.set.Items {
		if item.ID == id {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}
	if !found {
		return fmt.Errorf("item %s not found", id)
	}
	s.set.Items = newItems
	s.set.UpdatedAt = time.Now()
	return nil
}

// Summary provides aggregated statistics.
type Summary struct {
	TotalItems       int                                   `json:"total_items"`
	TotalEffortDays  int                                   `json:"total_effort_days"`
	MoSCoWCounts     map[prioritization.MoSCoWPriority]int `json:"moscow_counts"`
	StatusCounts     map[RMIStatus]int                     `json:"status_counts"`
	QuarterCounts    map[string]int                        `json:"quarter_counts,omitempty"`
	ActionableCount  int                                   `json:"actionable_count"`
	BlockedCount     int                                   `json:"blocked_count"`
	AvgPriorityScore float64                               `json:"avg_priority_score"`
}

// Summary returns aggregated statistics for all items.
func (s *Service) Summary() *Summary {
	summary := &Summary{
		TotalItems:      len(s.set.Items),
		TotalEffortDays: s.set.TotalEffortDays(),
		MoSCoWCounts:    s.set.MoSCoWSummary(),
		StatusCounts:    s.set.StatusSummary(),
		QuarterCounts:   make(map[string]int),
	}

	var totalScore float64
	for _, item := range s.set.Items {
		if item.IsActionable() {
			summary.ActionableCount++
		}
		if item.HasBlockingDependencies() {
			summary.BlockedCount++
		}
		if item.Quarter != "" {
			summary.QuarterCounts[item.Quarter]++
		}
		totalScore += item.PriorityScore()
	}

	if len(s.set.Items) > 0 {
		summary.AvgPriorityScore = totalScore / float64(len(s.set.Items))
	}

	return summary
}

// TopByPriority returns the top N items by priority score.
func (s *Service) TopByPriority(n int) []RoadmapItem {
	return s.set.TopN(n)
}

// TopByRICE returns the top N items by RICE score.
func (s *Service) TopByRICE(n int) []RoadmapItem {
	items := make([]RoadmapItem, len(s.set.Items))
	copy(items, s.set.Items)

	sort.Slice(items, func(i, j int) bool {
		scoreI := 0.0
		scoreJ := 0.0
		if items[i].RICE != nil {
			scoreI = items[i].RICE.Score
		}
		if items[j].RICE != nil {
			scoreJ = items[j].RICE.Score
		}
		return scoreI > scoreJ
	})

	if n > 0 && n < len(items) {
		items = items[:n]
	}
	return items
}

// TopByMarketSignal returns the top N items by market signal score.
func (s *Service) TopByMarketSignal(n int) []RoadmapItem {
	items := make([]RoadmapItem, len(s.set.Items))
	copy(items, s.set.Items)

	sort.Slice(items, func(i, j int) bool {
		scoreI := 0.0
		scoreJ := 0.0
		if items[i].MarketSignal != nil {
			scoreI = items[i].MarketSignal.Score
		}
		if items[j].MarketSignal != nil {
			scoreJ = items[j].MarketSignal.Score
		}
		return scoreI > scoreJ
	})

	if n > 0 && n < len(items) {
		items = items[:n]
	}
	return items
}

// ActionableItems returns all items that can be worked on.
func (s *Service) ActionableItems() []RoadmapItem {
	return s.set.Actionable()
}

// BlockedItems returns items with blocking dependencies.
func (s *Service) BlockedItems() []RoadmapItem {
	return s.set.Blocked()
}
