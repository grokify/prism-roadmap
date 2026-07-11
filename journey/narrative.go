package journey

import "strconv"

// RoadmapNarrative provides the story layer for a roadmap.
// Without narrative, a roadmap is just a data table.
// The narrative makes it communicable to executives and stakeholders.
type RoadmapNarrative struct {
	Title        string           `json:"title"`
	CurrentState string           `json:"currentState,omitempty"` // Where we are today
	TurningPoint string           `json:"turningPoint,omitempty"` // What changed/enables the journey
	Journey      []JourneyChapter `json:"journey,omitempty"`      // Story per period
	Destination  string           `json:"destination,omitempty"`  // Where we're going
	CallToAction string           `json:"callToAction,omitempty"` // What we need from the audience
}

// JourneyChapter is a story segment for a specific period.
type JourneyChapter struct {
	PeriodID   string   `json:"periodId"`
	Headline   string   `json:"headline"`             // 2-4 word summary
	Story      string   `json:"story"`                // 1-2 sentence narrative
	KeyChanges []string `json:"keyChanges,omitempty"` // Bullet points
	Milestones []string `json:"milestones,omitempty"` // Key deliverables
	UserImpact string   `json:"userImpact,omitempty"` // How users are affected
	Risks      []string `json:"risks,omitempty"`      // Key risks this period
}

// StoryboardCard represents a single card in a roadmap storyboard visualization.
// Each card summarizes a period in a visually digestible format.
type StoryboardCard struct {
	PeriodID          string           `json:"periodId"`
	PeriodLabel       string           `json:"periodLabel"`
	Headline          string           `json:"headline"`
	MaturityChanges   []MaturityChange `json:"maturityChanges,omitempty"`
	MajorInitiatives  []string         `json:"majorInitiatives,omitempty"`
	UserImpact        string           `json:"userImpact,omitempty"`
	SuccessEvidence   []string         `json:"successEvidence,omitempty"`
	UnresolvedRisks   []string         `json:"unresolvedRisks,omitempty"`
	OverallConfidence float64          `json:"overallConfidence,omitempty"` // 0.0-1.0
}

// MaturityChange summarizes a capability's maturity shift for a storyboard.
type MaturityChange struct {
	CapabilityName string `json:"capabilityName"`
	From           string `json:"from"` // Maturity level
	To             string `json:"to"`   // Maturity level
}

// BuildStoryboard creates storyboard cards from a journey roadmap.
func BuildStoryboard(roadmap *JourneyRoadmap) []StoryboardCard {
	if roadmap == nil || roadmap.TimeModel == nil {
		return nil
	}

	cards := make([]StoryboardCard, 0, len(roadmap.TimeModel.Periods))

	for _, period := range roadmap.TimeModel.Periods {
		card := StoryboardCard{
			PeriodID:    period.ID,
			PeriodLabel: period.Label,
		}

		// Get narrative chapter if available
		if roadmap.Narrative != nil {
			for _, chapter := range roadmap.Narrative.Journey {
				if chapter.PeriodID == period.ID {
					card.Headline = chapter.Headline
					card.UserImpact = chapter.UserImpact
					card.UnresolvedRisks = chapter.Risks
					break
				}
			}
		}

		// Collect maturity changes for this period
		for _, cj := range roadmap.CapabilityJourneys {
			for i, target := range cj.TargetStates {
				if target.PeriodID == period.ID {
					var fromLevel string
					if i == 0 && cj.CurrentState != nil {
						fromLevel = cj.CurrentState.MaturityLevel
					} else if i > 0 {
						fromLevel = cj.TargetStates[i-1].MaturityLevel
					}

					if fromLevel != "" && fromLevel != target.MaturityLevel {
						card.MaturityChanges = append(card.MaturityChanges, MaturityChange{
							CapabilityName: cj.Name,
							From:           fromLevel,
							To:             target.MaturityLevel,
						})
					}
				}
			}
		}

		// Collect initiatives active in this period
		for _, init := range roadmap.Initiatives {
			for _, p := range init.Periods {
				if p == period.ID {
					card.MajorInitiatives = append(card.MajorInitiatives, init.Name)
					break
				}
			}
		}

		// Calculate average confidence
		var totalConf float64
		var confCount int
		for _, cj := range roadmap.CapabilityJourneys {
			for _, target := range cj.TargetStates {
				if target.PeriodID == period.ID && target.Confidence > 0 {
					totalConf += target.Confidence
					confCount++
				}
			}
		}
		if confCount > 0 {
			card.OverallConfidence = totalConf / float64(confCount)
		}

		cards = append(cards, card)
	}

	return cards
}

// ExecutiveSummary generates a brief executive summary of the roadmap.
type ExecutiveSummary struct {
	Vision               string   `json:"vision"`
	CurrentStateOverview string   `json:"currentStateOverview"`
	KeyTransformations   []string `json:"keyTransformations"`
	ExpectedOutcomes     []string `json:"expectedOutcomes"`
	TopRisks             []string `json:"topRisks"`
	ResourceRequirements string   `json:"resourceRequirements,omitempty"`
	Timeline             string   `json:"timeline"`
	CallToAction         string   `json:"callToAction,omitempty"`
}

// GenerateExecutiveSummary creates a summary from a journey roadmap.
func GenerateExecutiveSummary(roadmap *JourneyRoadmap) *ExecutiveSummary {
	if roadmap == nil {
		return nil
	}

	summary := &ExecutiveSummary{
		Vision: roadmap.Vision,
	}

	// Timeline
	if roadmap.TimeModel != nil && len(roadmap.TimeModel.Periods) > 0 {
		first := roadmap.TimeModel.Periods[0].Label
		last := roadmap.TimeModel.Periods[len(roadmap.TimeModel.Periods)-1].Label
		summary.Timeline = first + " → " + last
	}

	// Key transformations (major maturity jumps)
	for _, cj := range roadmap.CapabilityJourneys {
		if cj.CurrentState != nil && cj.DesiredEndState != nil {
			from := cj.CurrentState.MaturityLevel
			to := cj.DesiredEndState.MaturityLevel
			if from != to {
				summary.KeyTransformations = append(summary.KeyTransformations,
					cj.Name+": "+from+" → "+to)
			}
		}
	}

	// Expected outcomes
	for _, oj := range roadmap.OutcomeJourneys {
		if oj.CurrentState != nil && len(oj.TargetStates) > 0 {
			last := oj.TargetStates[len(oj.TargetStates)-1]
			summary.ExpectedOutcomes = append(summary.ExpectedOutcomes,
				oj.Name+": "+formatValue(oj.CurrentState.Value, oj.CurrentState.Unit)+" → "+formatValue(last.Value, last.Unit))
		}
	}

	// Top risks
	for _, risk := range roadmap.Risks {
		if risk.Impact == "high" || risk.Impact == "critical" {
			summary.TopRisks = append(summary.TopRisks, risk.Description)
		}
	}

	// Narrative elements
	if roadmap.Narrative != nil {
		summary.CurrentStateOverview = roadmap.Narrative.CurrentState
		summary.CallToAction = roadmap.Narrative.CallToAction
	}

	return summary
}

func formatValue(value float64, unit string) string {
	// Simple formatting - real implementation might be more sophisticated
	if unit == "" {
		return formatNumber(value)
	}
	return formatNumber(value) + " " + unit
}

func formatNumber(v float64) string {
	// Use strconv for proper formatting
	if v == float64(int(v)) {
		return strconv.Itoa(int(v))
	}
	return strconv.FormatFloat(v, 'f', 1, 64)
}
