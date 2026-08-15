package rmi

import (
	"sort"

	"github.com/grokify/prism-roadmap/goals/okr"
)

// ObjectiveRICE summarizes the RICE-scored roadmap items linked to a single OKR
// objective. It answers the planning question "how much prioritized (RICE)
// effort is backing each objective?" by joining roadmap items to objectives
// through each item's ObjectiveRefs.
type ObjectiveRICE struct {
	ObjectiveID    string        `json:"objective_id"`
	ObjectiveTitle string        `json:"objective_title,omitempty"`
	Items          []RoadmapItem `json:"items,omitempty"` // items linked to this objective
	ScoredCount    int           `json:"scored_count"`    // linked items with a RICE score > 0
	TotalRICE      float64       `json:"total_rice"`      // sum of linked RICE scores
	MeanRICE       float64       `json:"mean_rice"`       // mean RICE over scored items
	TopRICE        float64       `json:"top_rice"`        // highest single RICE score
}

// RICEByObjective joins RICE-scored roadmap items to OKR objectives via each
// item's ObjectiveRefs and returns one ObjectiveRICE per objective. Results are
// ordered by TotalRICE descending, so objectives backed by the most prioritized
// work rank first. Objectives with no linked items are included with zero
// aggregates so gaps in coverage are visible.
func RICEByObjective(items []RoadmapItem, objectives []okr.Objective) []ObjectiveRICE {
	rollups := make([]ObjectiveRICE, 0, len(objectives))
	for _, obj := range objectives {
		roll := ObjectiveRICE{
			ObjectiveID:    obj.ID,
			ObjectiveTitle: obj.Title,
		}
		for _, item := range items {
			if !item.SupportsObjective(obj.ID) {
				continue
			}
			roll.Items = append(roll.Items, item)
			if item.RICE != nil && item.RICE.Score > 0 {
				roll.ScoredCount++
				roll.TotalRICE += item.RICE.Score
				if item.RICE.Score > roll.TopRICE {
					roll.TopRICE = item.RICE.Score
				}
			}
		}
		if roll.ScoredCount > 0 {
			roll.MeanRICE = roll.TotalRICE / float64(roll.ScoredCount)
		}
		rollups = append(rollups, roll)
	}
	sort.SliceStable(rollups, func(i, j int) bool {
		return rollups[i].TotalRICE > rollups[j].TotalRICE
	})
	return rollups
}

// UnlinkedScoredItems returns RICE-scored items that are not linked to any OKR
// objective. These are prioritized work items with no stated objective — useful
// for surfacing roadmap work that is not tied to a goal.
func UnlinkedScoredItems(items []RoadmapItem) []RoadmapItem {
	var orphans []RoadmapItem
	for _, item := range items {
		if len(item.ObjectiveRefs) == 0 && item.RICE != nil && item.RICE.Score > 0 {
			orphans = append(orphans, item)
		}
	}
	return orphans
}
