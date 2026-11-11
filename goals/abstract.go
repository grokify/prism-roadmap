package goals

// GoalItem represents a high-level goal that can be either:
// - An Objective (OKR framework)
// - A Method (V2MOM framework)
// - A Strategy (VSEM framework)
//
// This abstraction enables framework-agnostic rendering for roadmaps,
// reports, and presentations.
type GoalItem struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Owner       string   `json:"owner,omitempty"`
	Status      string   `json:"status,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ResultItem represents a measurable result that can be either:
// - A Key Result (OKR framework)
// - A Measure (V2MOM framework)
// - A Metric (VSEM framework)
//
// This abstraction enables framework-agnostic rendering for tracking
// and reporting.
type ResultItem struct {
	ID          string  `json:"id,omitempty"`
	Title       string  `json:"title"`
	Metric      string  `json:"metric,omitempty"`
	Baseline    string  `json:"baseline,omitempty"`
	Target      string  `json:"target,omitempty"`
	Current     string  `json:"current,omitempty"`
	Unit        string  `json:"unit,omitempty"`
	Status      string  `json:"status,omitempty"`
	Score       float64 `json:"score,omitempty"` // 0.0-1.0 achievement score
	GoalID      string  `json:"goalId,omitempty"`
	PhaseID     string  `json:"phaseId,omitempty"`     // For roadmap alignment
	PhaseTarget string  `json:"phaseTarget,omitempty"` // Target for specific phase
}

// GoalItems returns all goals as a unified slice.
// Returns Objectives for OKR, Methods for V2MOM, and Strategies for VSEM.
func (g *Goals) GoalItems() []GoalItem {
	if g == nil {
		return nil
	}

	switch g.Framework {
	case FrameworkOKR:
		if g.OKR == nil {
			return nil
		}
		items := make([]GoalItem, 0, len(g.OKR.OKRs))
		for _, okrItem := range g.OKR.OKRs {
			obj := okrItem.Objective
			title := obj.Title
			if title == "" {
				title = obj.Description
			}
			items = append(items, GoalItem{
				ID:          obj.ID,
				Title:       title,
				Description: obj.Description,
				Owner:       obj.Owner,
				Status:      obj.Status,
				Tags:        obj.Tags,
			})
		}
		return items

	case FrameworkV2MOM:
		if g.V2MOM == nil {
			return nil
		}
		items := make([]GoalItem, 0, len(g.V2MOM.Methods))
		for _, method := range g.V2MOM.Methods {
			items = append(items, GoalItem{
				ID:          method.ID,
				Title:       method.Name,
				Description: method.Description,
				Owner:       method.Owner,
				Status:      method.Status,
				Priority:    method.Priority,
			})
		}
		return items

	case FrameworkVSEM:
		if g.VSEM == nil {
			return nil
		}
		items := make([]GoalItem, 0, len(g.VSEM.Strategies))
		for _, s := range g.VSEM.Strategies {
			items = append(items, GoalItem{
				ID:          s.ID,
				Title:       s.Name,
				Description: s.Description,
				Owner:       s.Owner,
				Priority:    s.Priority,
			})
		}
		return items

	default:
		return nil
	}
}

// ResultItems returns all measurable results as a unified slice.
// Returns Key Results for OKR, Measures for V2MOM, and Metrics for VSEM. A
// VSEM metric that measures the vision as a whole has no GoalID.
func (g *Goals) ResultItems() []ResultItem {
	if g == nil {
		return nil
	}

	switch g.Framework {
	case FrameworkOKR:
		if g.OKR == nil {
			return nil
		}
		var items []ResultItem
		for _, okrItem := range g.OKR.OKRs {
			goalID := okrItem.Objective.ID
			krs := okrItem.KeyResults
			if len(krs) == 0 {
				krs = okrItem.Objective.KeyResults
			}
			for _, kr := range krs {
				title := kr.Title
				if title == "" {
					title = kr.Description
				}
				item := ResultItem{
					ID:       kr.ID,
					Title:    title,
					Metric:   kr.Metric,
					Baseline: kr.Baseline,
					Target:   kr.Target,
					Current:  kr.Current,
					Unit:     kr.Unit,
					Status:   kr.Status,
					Score:    kr.Score,
					GoalID:   goalID,
				}
				// Add phase targets if present
				for _, pt := range kr.PhaseTargets {
					ptItem := item
					ptItem.PhaseID = pt.PhaseID
					ptItem.PhaseTarget = pt.Target
					if pt.Status != "" {
						ptItem.Status = pt.Status
					}
					items = append(items, ptItem)
				}
				// If no phase targets, add the base item
				if len(kr.PhaseTargets) == 0 {
					items = append(items, item)
				}
			}
		}
		return items

	case FrameworkV2MOM:
		if g.V2MOM == nil {
			return nil
		}
		var items []ResultItem
		// Global measures
		for _, m := range g.V2MOM.Measures {
			items = append(items, ResultItem{
				ID:       m.ID,
				Title:    m.Name,
				Metric:   m.Name,
				Baseline: m.Baseline,
				Target:   m.Target,
				Current:  m.Current,
				Unit:     m.Unit,
				Status:   m.Status,
				Score:    m.Progress,
			})
		}
		// Nested measures under methods
		for _, method := range g.V2MOM.Methods {
			for _, m := range method.Measures {
				items = append(items, ResultItem{
					ID:       m.ID,
					Title:    m.Name,
					Metric:   m.Name,
					Baseline: m.Baseline,
					Target:   m.Target,
					Current:  m.Current,
					Unit:     m.Unit,
					Status:   m.Status,
					Score:    m.Progress,
					GoalID:   method.ID,
				})
			}
		}
		return items

	case FrameworkVSEM:
		if g.VSEM == nil {
			return nil
		}
		items := make([]ResultItem, 0, len(g.VSEM.Metrics))
		for _, m := range g.VSEM.Metrics {
			items = append(items, ResultItem{
				ID:       m.ID,
				Title:    m.Name,
				Metric:   m.Name,
				Baseline: m.Baseline,
				Target:   m.Target,
				Current:  m.Current,
				Unit:     m.Unit,
				Status:   m.Status,
				Score:    m.Progress,
				GoalID:   m.StrategyID,
			})
		}
		return items

	default:
		return nil
	}
}

// GoalLabel returns the display label for goals.
// Returns "Objectives" for OKR, "Methods" for V2MOM, "Strategies" for VSEM.
func (g *Goals) GoalLabel() string {
	if g == nil {
		return "Goals"
	}
	switch g.Framework {
	case FrameworkOKR:
		return "Objectives"
	case FrameworkV2MOM:
		return "Methods"
	case FrameworkVSEM:
		return "Strategies"
	default:
		return "Goals"
	}
}

// ResultLabel returns the display label for results.
// Returns "Key Results" for OKR, "Measures" for V2MOM, "Metrics" for VSEM.
func (g *Goals) ResultLabel() string {
	if g == nil {
		return "Results"
	}
	switch g.Framework {
	case FrameworkOKR:
		return "Key Results"
	case FrameworkV2MOM:
		return "Measures"
	case FrameworkVSEM:
		return "Metrics"
	default:
		return "Results"
	}
}

// GoalCount returns the number of goals.
func (g *Goals) GoalCount() int {
	return len(g.GoalItems())
}

// ResultCount returns the number of results.
func (g *Goals) ResultCount() int {
	return len(g.ResultItems())
}

// ResultItemsByPhase returns results grouped by phase ID.
func (g *Goals) ResultItemsByPhase() map[string][]ResultItem {
	results := g.ResultItems()
	byPhase := make(map[string][]ResultItem)
	for _, r := range results {
		if r.PhaseID != "" {
			byPhase[r.PhaseID] = append(byPhase[r.PhaseID], r)
		}
	}
	return byPhase
}
