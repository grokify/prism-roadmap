// Package markdown provides Markdown table renderers for canvas types.
package markdown

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

// MarkdownRenderer renders canvas types to Markdown table format.
type MarkdownRenderer struct{}

// NewMarkdownRenderer creates a new Markdown renderer.
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// Format returns the output format.
func (r *MarkdownRenderer) Format() render.Format {
	return render.FormatMarkdown
}

// FileExtension returns the file extension for Markdown files.
func (r *MarkdownRenderer) FileExtension() string {
	return ".md"
}

// Supports returns true for all canvas types.
func (r *MarkdownRenderer) Supports(_ canvas.CanvasType) bool {
	return true
}

// Render converts a canvas to Markdown format.
func (r *MarkdownRenderer) Render(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	switch c.Type {
	case canvas.CanvasTypeBMC:
		return r.renderBMC(c.BMC, opts)
	case canvas.CanvasTypeOpportunity:
		return r.renderOpportunity(c.Opportunity, opts)
	case canvas.CanvasTypeFeature:
		return r.renderFeature(c.Feature, opts)
	case canvas.CanvasTypeLeanUX:
		return r.renderLeanUX(c.LeanUX, opts)
	case canvas.CanvasTypeOST:
		return r.renderOST(c.OST, opts)
	case canvas.CanvasTypeShapeUpPitch:
		return r.renderShapeUpPitch(c.ShapeUpPitch, opts)
	case canvas.CanvasTypeShapeUpBet:
		return r.renderShapeUpBet(c.ShapeUpBet, opts)
	case canvas.CanvasTypeShapeUpScope:
		return r.renderShapeUpScope(c.ShapeUpScope, opts)
	case canvas.CanvasTypeDiscoverySnapshot:
		return r.renderDiscoverySnapshot(c.DiscoverySnapshot, opts)
	case canvas.CanvasTypeAssumptionMap:
		return r.renderAssumptionMap(c.AssumptionMap, opts)
	case canvas.CanvasTypeExperienceMap:
		return r.renderExperienceMap(c.ExperienceMap, opts)
	case canvas.CanvasTypeLeanStartup:
		return r.renderLeanStartup(c.LeanStartup, opts)
	case canvas.CanvasTypeDesignThinking:
		return r.renderDesignThinking(c.DesignThinking, opts)
	case canvas.CanvasTypeJTBD:
		return r.renderJTBD(c.JTBD, opts)
	default:
		return nil, fmt.Errorf("unsupported canvas type: %s", c.Type)
	}
}

func (r *MarkdownRenderer) renderBMC(bmc *canvas.BusinessModelCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + bmc.Metadata.Title + "\n\n")

	// Customer Segments
	sb.WriteString("## Customer Segments\n\n")
	if len(bmc.CustomerSegments) > 0 {
		sb.WriteString("| ID | Name | Size |\n")
		sb.WriteString("|---|---|---|\n")
		for _, seg := range bmc.CustomerSegments {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", seg.ID, seg.Name, seg.Size))
		}
	}
	sb.WriteString("\n")

	// Value Propositions
	sb.WriteString("## Value Propositions\n\n")
	if len(bmc.ValuePropositions) > 0 {
		sb.WriteString("| ID | Description |\n")
		sb.WriteString("|---|---|\n")
		for _, vp := range bmc.ValuePropositions {
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", vp.ID, vp.Description))
		}
	}
	sb.WriteString("\n")

	// Channels
	sb.WriteString("## Channels\n\n")
	if len(bmc.Channels) > 0 {
		sb.WriteString("| ID | Name | Type | Phase |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, ch := range bmc.Channels {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", ch.ID, ch.Name, ch.Type, ch.Phase))
		}
	}
	sb.WriteString("\n")

	// Key Resources
	sb.WriteString("## Key Resources\n\n")
	if len(bmc.KeyResources) > 0 {
		sb.WriteString("| ID | Name | Type |\n")
		sb.WriteString("|---|---|---|\n")
		for _, res := range bmc.KeyResources {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", res.ID, res.Name, res.Type))
		}
	}
	sb.WriteString("\n")

	// Key Activities
	sb.WriteString("## Key Activities\n\n")
	if len(bmc.KeyActivities) > 0 {
		sb.WriteString("| ID | Name | Category |\n")
		sb.WriteString("|---|---|---|\n")
		for _, act := range bmc.KeyActivities {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", act.ID, act.Name, act.Category))
		}
	}
	sb.WriteString("\n")

	// Key Partnerships
	sb.WriteString("## Key Partnerships\n\n")
	if len(bmc.KeyPartnerships) > 0 {
		sb.WriteString("| ID | Partner | Type |\n")
		sb.WriteString("|---|---|---|\n")
		for _, p := range bmc.KeyPartnerships {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", p.ID, p.Partner, p.Type))
		}
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderOpportunity(opp *canvas.OpportunityCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + opp.Metadata.Title + "\n\n")

	// Problems
	sb.WriteString("## Problems\n\n")
	if len(opp.Problems) > 0 {
		sb.WriteString("| ID | Description | Severity |\n")
		sb.WriteString("|---|---|---|\n")
		for _, p := range opp.Problems {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", p.ID, p.Description, p.Severity))
		}
	}
	sb.WriteString("\n")

	// Users
	sb.WriteString("## Users\n\n")
	if len(opp.Users) > 0 {
		sb.WriteString("| ID | Name | Description |\n")
		sb.WriteString("|---|---|---|\n")
		for _, u := range opp.Users {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", u.ID, u.Name, u.Description))
		}
	}
	sb.WriteString("\n")

	// Value Proposition
	sb.WriteString("## Value Proposition\n\n")
	sb.WriteString(opp.ValueProposition.Statement + "\n\n")

	// Assumptions
	sb.WriteString("## Assumptions\n\n")
	if len(opp.Assumptions) > 0 {
		sb.WriteString("| ID | Description | Validated |\n")
		sb.WriteString("|---|---|---|\n")
		for _, a := range opp.Assumptions {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", a.ID, a.Description, validated))
		}
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderFeature(fc *canvas.FeatureCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + fc.Metadata.Title + "\n\n")
	sb.WriteString("**Idea:** " + fc.IdeaStatement + "\n\n")

	// Situations
	sb.WriteString("## Situations\n\n")
	if len(fc.Situations) > 0 {
		sb.WriteString("| ID | Description | Actor |\n")
		sb.WriteString("|---|---|---|\n")
		for _, s := range fc.Situations {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", s.ID, s.Description, s.Actor))
		}
	}
	sb.WriteString("\n")

	// Problems
	sb.WriteString("## Problems\n\n")
	if len(fc.Problems) > 0 {
		sb.WriteString("| ID | Description |\n")
		sb.WriteString("|---|---|\n")
		for _, p := range fc.Problems {
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", p.ID, p.Description))
		}
	}
	sb.WriteString("\n")

	// Capabilities
	sb.WriteString("## Capabilities\n\n")
	if len(fc.Capabilities) > 0 {
		sb.WriteString("| ID | Description | Priority |\n")
		sb.WriteString("|---|---|---|\n")
		for _, c := range fc.Capabilities {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", c.ID, c.Description, c.Priority))
		}
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderLeanUX(lux *canvas.LeanUXCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + lux.Metadata.Title + "\n\n")

	// Business Problem
	sb.WriteString("## Business Problem\n\n")
	sb.WriteString(lux.BusinessProblem + "\n\n")

	// Users
	sb.WriteString("## Users\n\n")
	if len(lux.Users) > 0 {
		sb.WriteString("| ID | Name | Description |\n")
		sb.WriteString("|---|---|---|\n")
		for _, u := range lux.Users {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", u.ID, u.Name, u.Description))
		}
	}
	sb.WriteString("\n")

	// Hypotheses
	sb.WriteString("## Hypotheses\n\n")
	if len(lux.Hypotheses) > 0 {
		sb.WriteString("| ID | We Believe | Will Result In | Validated |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, h := range lux.Hypotheses {
			validated := "-"
			if h.Validated != nil {
				if *h.Validated {
					validated = "Yes"
				} else {
					validated = "No"
				}
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", h.ID, h.WeBelieve, h.WillResultIn, validated))
		}
	}
	sb.WriteString("\n")

	// Riskiest Assumption
	sb.WriteString("## Riskiest Assumption\n\n")
	sb.WriteString(lux.RiskiestAssumption + "\n\n")

	// Experiment
	if lux.Experiment != nil {
		sb.WriteString("## Experiment\n\n")
		sb.WriteString("| Field | Value |\n")
		sb.WriteString("|---|---|\n")
		sb.WriteString(fmt.Sprintf("| Description | %s |\n", lux.Experiment.Description))
		sb.WriteString(fmt.Sprintf("| Method | %s |\n", lux.Experiment.Method))
		sb.WriteString(fmt.Sprintf("| Status | %s |\n", lux.Experiment.Status))
		sb.WriteString(fmt.Sprintf("| Success Criteria | %s |\n", lux.Experiment.SuccessCriteria))
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderOST(ost *canvas.OpportunitySolutionTree, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + ost.Metadata.Title + "\n\n")

	// Outcome
	sb.WriteString("## Outcome\n\n")
	sb.WriteString(fmt.Sprintf("**%s:** %s\n\n", ost.Outcome.ID, ost.Outcome.Description))
	if ost.Outcome.Metric != "" {
		sb.WriteString(fmt.Sprintf("- **Metric:** %s\n", ost.Outcome.Metric))
	}
	if ost.Outcome.Target != "" {
		sb.WriteString(fmt.Sprintf("- **Target:** %s\n", ost.Outcome.Target))
	}
	sb.WriteString("\n")

	// Opportunities
	sb.WriteString("## Opportunities\n\n")
	for _, opp := range ost.Outcome.Opportunities {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", opp.ID, opp.Description))
		if opp.Source != "" {
			sb.WriteString(fmt.Sprintf("- **Source:** %s\n", opp.Source))
		}

		// Solutions
		if len(opp.Solutions) > 0 {
			sb.WriteString("\n**Solutions:**\n\n")
			sb.WriteString("| ID | Description | Type | Status |\n")
			sb.WriteString("|---|---|---|---|\n")
			for _, sol := range opp.Solutions {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", sol.ID, sol.Description, sol.Type, sol.Status))
			}
		}
		sb.WriteString("\n")
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderShapeUpPitch(pitch *canvas.ShapeUpPitch, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + pitch.Metadata.Title + "\n\n")

	// Problem
	sb.WriteString("## Problem\n\n")
	sb.WriteString(fmt.Sprintf("**Statement:** %s\n\n", pitch.Problem.Statement))
	if pitch.Problem.WhyNow != "" {
		sb.WriteString(fmt.Sprintf("**Why Now:** %s\n\n", pitch.Problem.WhyNow))
	}
	if pitch.Problem.Audience != "" {
		sb.WriteString(fmt.Sprintf("**Audience:** %s\n\n", pitch.Problem.Audience))
	}
	if len(pitch.Problem.Evidence) > 0 {
		sb.WriteString("**Evidence:**\n\n")
		for _, e := range pitch.Problem.Evidence {
			sb.WriteString(fmt.Sprintf("- %s\n", e))
		}
		sb.WriteString("\n")
	}

	// Appetite
	sb.WriteString("## Appetite\n\n")
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|---|---|\n")
	sb.WriteString(fmt.Sprintf("| Weeks | %d |\n", pitch.Appetite.Weeks))
	sb.WriteString(fmt.Sprintf("| Size | %s |\n", pitch.Appetite.Size))
	if pitch.Appetite.Rationale != "" {
		sb.WriteString(fmt.Sprintf("| Rationale | %s |\n", pitch.Appetite.Rationale))
	}
	sb.WriteString("\n")

	// Solution
	sb.WriteString("## Solution\n\n")
	sb.WriteString(fmt.Sprintf("**Approach:** %s\n\n", pitch.Solution.Approach))

	if len(pitch.Solution.MustInclude) > 0 {
		sb.WriteString("**Must Include:**\n\n")
		for _, item := range pitch.Solution.MustInclude {
			sb.WriteString(fmt.Sprintf("- %s\n", item))
		}
		sb.WriteString("\n")
	}

	if len(pitch.Solution.NiceToHave) > 0 {
		sb.WriteString("**Nice to Have:**\n\n")
		for _, item := range pitch.Solution.NiceToHave {
			sb.WriteString(fmt.Sprintf("- %s\n", item))
		}
		sb.WriteString("\n")
	}

	// Rabbit Holes
	if len(pitch.RabbitHoles) > 0 {
		sb.WriteString("## Rabbit Holes\n\n")
		sb.WriteString("| ID | Description | Why Dangerous |\n")
		sb.WriteString("|---|---|---|\n")
		for _, rh := range pitch.RabbitHoles {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", rh.ID, rh.Description, rh.WhyDangerous))
		}
		sb.WriteString("\n")
	}

	// No-Gos
	if len(pitch.NoGos) > 0 {
		sb.WriteString("## No-Gos\n\n")
		for _, nogo := range pitch.NoGos {
			sb.WriteString(fmt.Sprintf("- %s\n", nogo))
		}
		sb.WriteString("\n")
	}

	// Betting Status
	if pitch.BettingStatus != "" {
		sb.WriteString("## Betting Status\n\n")
		sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", pitch.BettingStatus))
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderShapeUpBet(bet *canvas.ShapeUpBet, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + bet.Metadata.Title + "\n\n")

	// Pitch Reference
	sb.WriteString("## Pitch\n\n")
	sb.WriteString(fmt.Sprintf("**Reference:** %s\n\n", bet.PitchRef))
	if bet.PitchTitle != "" {
		sb.WriteString(fmt.Sprintf("**Title:** %s\n\n", bet.PitchTitle))
	}

	// Cycle
	sb.WriteString("## Cycle\n\n")
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|---|---|\n")
	sb.WriteString(fmt.Sprintf("| ID | %s |\n", bet.Cycle.ID))
	sb.WriteString(fmt.Sprintf("| Name | %s |\n", bet.Cycle.Name))
	sb.WriteString(fmt.Sprintf("| Start | %s |\n", bet.Cycle.StartDate))
	sb.WriteString(fmt.Sprintf("| End | %s |\n", bet.Cycle.EndDate))
	sb.WriteString(fmt.Sprintf("| Weeks | %d |\n", bet.Cycle.Weeks))
	sb.WriteString("\n")

	// Team
	sb.WriteString("## Team\n\n")
	if bet.Team.Designer != "" {
		sb.WriteString(fmt.Sprintf("- **Designer:** %s\n", bet.Team.Designer))
	}
	if len(bet.Team.Programmers) > 0 {
		sb.WriteString("- **Programmers:** " + strings.Join(bet.Team.Programmers, ", ") + "\n")
	}
	if bet.Team.Lead != "" {
		sb.WriteString(fmt.Sprintf("- **Lead:** %s\n", bet.Team.Lead))
	}
	sb.WriteString("\n")

	// Decision
	sb.WriteString("## Decision\n\n")
	sb.WriteString(fmt.Sprintf("**Decision:** %s\n\n", bet.Decision))
	if bet.Rationale != "" {
		sb.WriteString(fmt.Sprintf("**Rationale:** %s\n\n", bet.Rationale))
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderShapeUpScope(scope *canvas.ShapeUpScope, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + scope.Metadata.Title + "\n\n")

	// Overall Progress
	progress := scope.OverallProgress()
	sb.WriteString(fmt.Sprintf("**Overall Progress:** %d%%\n\n", progress))
	sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", scope.Status))

	// Hill Chart Summary
	uphill := scope.UphillScopes()
	downhill := scope.DownhillScopes()
	done := scope.DoneScopes()

	sb.WriteString("## Hill Chart Summary\n\n")
	sb.WriteString("| Phase | Count |\n")
	sb.WriteString("|---|---|\n")
	sb.WriteString(fmt.Sprintf("| Uphill (Figuring Out) | %d |\n", len(uphill)))
	sb.WriteString(fmt.Sprintf("| Downhill (Executing) | %d |\n", len(downhill)-len(done)))
	sb.WriteString(fmt.Sprintf("| Done | %d |\n", len(done)))
	sb.WriteString("\n")

	// Scopes Table
	sb.WriteString("## Scopes\n\n")
	sb.WriteString("| ID | Name | Hill Position | Status |\n")
	sb.WriteString("|---|---|---|---|\n")
	for _, s := range scope.Scopes {
		phase := "Uphill"
		if s.HillPosition == 100 || s.Status == "done" {
			phase = "Done"
		} else if s.HillPosition >= 50 {
			phase = "Downhill"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %d%% (%s) | %s |\n", s.ID, s.Name, s.HillPosition, phase, s.Status))
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderDiscoverySnapshot(ds *canvas.DiscoverySnapshot, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + ds.Metadata.Title + "\n\n")
	sb.WriteString(fmt.Sprintf("**Week:** %s\n\n", ds.Week))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Interviews:** %d\n", ds.InterviewCount()))
	sb.WriteString(fmt.Sprintf("- **Stories Collected:** %d\n", ds.TotalStories()))
	sb.WriteString(fmt.Sprintf("- **Opportunities Updated:** %d\n", len(ds.OpportunitiesDiscovered)))
	sb.WriteString(fmt.Sprintf("- **Assumption Tests:** %d\n", len(ds.AssumptionTests)))
	sb.WriteString(fmt.Sprintf("- **Weekly Touchpoint:** %v\n", ds.HasWeeklyTouchpoint()))
	sb.WriteString("\n")

	// Interviews
	if len(ds.Interviews) > 0 {
		sb.WriteString("## Interviews\n\n")
		sb.WriteString("| ID | Participant Type | Interview Type | Quality |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, interview := range ds.Interviews {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				interview.ID, interview.ParticipantType, interview.InterviewType, interview.Quality))
		}
		sb.WriteString("\n")
	}

	// Opportunities Discovered
	if len(ds.OpportunitiesDiscovered) > 0 {
		sb.WriteString("## Opportunities\n\n")
		sb.WriteString("| ID | Action | Description | Evidence Count |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, opp := range ds.OpportunitiesDiscovered {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n",
				opp.OpportunityID, opp.Action, opp.Description, opp.EvidenceCount))
		}
		sb.WriteString("\n")
	}

	// Assumption Tests
	if len(ds.AssumptionTests) > 0 {
		sb.WriteString("## Assumption Tests\n\n")
		sb.WriteString("| ID | Assumption | Method | Status | Result |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, test := range ds.AssumptionTests {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				test.ID, test.Assumption.Description, test.Method, test.Status, test.Result))
		}
		sb.WriteString("\n")
	}

	// Key Learnings
	if len(ds.KeyLearnings) > 0 {
		sb.WriteString("## Key Learnings\n\n")
		for _, learning := range ds.KeyLearnings {
			sb.WriteString(fmt.Sprintf("- %s\n", learning))
		}
		sb.WriteString("\n")
	}

	// Decisions
	if len(ds.Decisions) > 0 {
		sb.WriteString("## Decisions\n\n")
		sb.WriteString("| ID | Decision | Type | Rationale |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, dec := range ds.Decisions {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dec.ID, dec.Description, dec.Type, dec.Rationale))
		}
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderAssumptionMap(am *canvas.AssumptionMap, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + am.Metadata.Title + "\n\n")

	if am.SolutionRef != "" {
		sb.WriteString(fmt.Sprintf("**Solution:** %s\n\n", am.SolutionRef))
	}

	// High Risk Summary
	highRisk := am.HighRiskAssumptions()
	if len(highRisk) > 0 {
		sb.WriteString("## High Risk Assumptions (Test First)\n\n")
		sb.WriteString("| ID | Description | Type |\n")
		sb.WriteString("|---|---|---|\n")
		for _, a := range highRisk {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", a.ID, a.Description, a.Type))
		}
		sb.WriteString("\n")
	}

	// Desirability
	if len(am.Desirability) > 0 {
		sb.WriteString("## Desirability\n\n")
		sb.WriteString("| ID | Description | Importance | Confidence | Validated |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, a := range am.Desirability {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				a.ID, a.Description, a.Importance, a.Confidence, validated))
		}
		sb.WriteString("\n")
	}

	// Viability
	if len(am.Viability) > 0 {
		sb.WriteString("## Viability\n\n")
		sb.WriteString("| ID | Description | Importance | Confidence | Validated |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, a := range am.Viability {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				a.ID, a.Description, a.Importance, a.Confidence, validated))
		}
		sb.WriteString("\n")
	}

	// Feasibility
	if len(am.Feasibility) > 0 {
		sb.WriteString("## Feasibility\n\n")
		sb.WriteString("| ID | Description | Importance | Confidence | Validated |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, a := range am.Feasibility {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				a.ID, a.Description, a.Importance, a.Confidence, validated))
		}
		sb.WriteString("\n")
	}

	// Usability
	if len(am.Usability) > 0 {
		sb.WriteString("## Usability\n\n")
		sb.WriteString("| ID | Description | Importance | Confidence | Validated |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, a := range am.Usability {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				a.ID, a.Description, a.Importance, a.Confidence, validated))
		}
		sb.WriteString("\n")
	}

	// Ethical
	if len(am.Ethical) > 0 {
		sb.WriteString("## Ethical\n\n")
		sb.WriteString("| ID | Description | Importance | Confidence | Validated |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, a := range am.Ethical {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				a.ID, a.Description, a.Importance, a.Confidence, validated))
		}
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderExperienceMap(em *canvas.ExperienceMap, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + em.Metadata.Title + "\n\n")

	sb.WriteString(fmt.Sprintf("**Experience:** %s\n\n", em.Experience))
	if em.PersonaDescription != "" {
		sb.WriteString(fmt.Sprintf("**Persona:** %s\n\n", em.PersonaDescription))
	}

	// Top Pain Points
	if len(em.TopPainPoints) > 0 {
		sb.WriteString("## Top Pain Points\n\n")
		for _, pain := range em.TopPainPoints {
			sb.WriteString(fmt.Sprintf("- %s\n", pain))
		}
		sb.WriteString("\n")
	}

	// Top Opportunities
	if len(em.TopOpportunities) > 0 {
		sb.WriteString("## Top Opportunities\n\n")
		for _, opp := range em.TopOpportunities {
			sb.WriteString(fmt.Sprintf("- %s\n", opp))
		}
		sb.WriteString("\n")
	}

	// Journey Phases
	sb.WriteString("## Journey Phases\n\n")
	for _, phase := range em.Phases {
		sb.WriteString(fmt.Sprintf("### %d. %s\n\n", phase.Order, phase.Name))
		if phase.Description != "" {
			sb.WriteString(phase.Description + "\n\n")
		}

		if phase.Thinking != "" {
			sb.WriteString(fmt.Sprintf("**Thinking:** %s\n\n", phase.Thinking))
		}
		if phase.Feeling != "" {
			sb.WriteString(fmt.Sprintf("**Feeling:** %s\n\n", phase.Feeling))
		}

		// Actions
		if len(phase.Actions) > 0 {
			sb.WriteString("**Actions:**\n\n")
			sb.WriteString("| Description | Channel | Pain | Delight |\n")
			sb.WriteString("|---|---|---|---|\n")
			for _, action := range phase.Actions {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					action.Description, action.Channel, action.Pain, action.Delight))
			}
			sb.WriteString("\n")
		}

		// Pain Points
		if len(phase.PainPoints) > 0 {
			sb.WriteString("**Pain Points:**\n\n")
			for _, pain := range phase.PainPoints {
				sb.WriteString(fmt.Sprintf("- %s\n", pain))
			}
			sb.WriteString("\n")
		}

		// Opportunities
		if len(phase.Opportunities) > 0 {
			sb.WriteString("**Opportunities:**\n\n")
			for _, opp := range phase.Opportunities {
				sb.WriteString(fmt.Sprintf("- %s\n", opp))
			}
			sb.WriteString("\n")
		}
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderLeanStartup(ls *canvas.LeanStartupCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + ls.Metadata.Title + "\n\n")

	// Vision and Strategy
	sb.WriteString("## Vision & Strategy\n\n")
	if ls.Vision != "" {
		sb.WriteString(fmt.Sprintf("**Vision:** %s\n\n", ls.Vision))
	}
	if ls.Strategy != "" {
		sb.WriteString(fmt.Sprintf("**Strategy:** %s\n\n", ls.Strategy))
	}
	if ls.TargetCustomer != "" {
		sb.WriteString(fmt.Sprintf("**Target Customer:** %s\n\n", ls.TargetCustomer))
	}
	if ls.ProblemHypothesis != "" {
		sb.WriteString(fmt.Sprintf("**Problem Hypothesis:** %s\n\n", ls.ProblemHypothesis))
	}
	if ls.SolutionHypothesis != "" {
		sb.WriteString(fmt.Sprintf("**Solution Hypothesis:** %s\n\n", ls.SolutionHypothesis))
	}
	if ls.ProductMarketFit != "" {
		sb.WriteString(fmt.Sprintf("**Product-Market Fit Status:** %s\n\n", ls.ProductMarketFit))
	}

	// Core Hypotheses
	sb.WriteString("## Core Hypotheses\n\n")

	if ls.ValueHypothesis != nil {
		sb.WriteString("### Value Hypothesis\n\n")
		sb.WriteString(fmt.Sprintf("**Statement:** %s\n\n", ls.ValueHypothesis.Statement))
		if ls.ValueHypothesis.Validated != nil {
			status := "Invalidated"
			if *ls.ValueHypothesis.Validated {
				status = "Validated"
			}
			sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", status))
		}
		if ls.ValueHypothesis.Evidence != "" {
			sb.WriteString(fmt.Sprintf("**Evidence:** %s\n\n", ls.ValueHypothesis.Evidence))
		}
		if ls.ValueHypothesis.Confidence != "" {
			sb.WriteString(fmt.Sprintf("**Confidence:** %s\n\n", ls.ValueHypothesis.Confidence))
		}
	}

	if ls.GrowthHypothesis != nil {
		sb.WriteString("### Growth Hypothesis\n\n")
		sb.WriteString(fmt.Sprintf("**Growth Model:** %s\n\n", ls.GrowthHypothesis.GrowthModel))
		sb.WriteString(fmt.Sprintf("**Statement:** %s\n\n", ls.GrowthHypothesis.Statement))
		if ls.GrowthHypothesis.Validated != nil {
			status := "Invalidated"
			if *ls.GrowthHypothesis.Validated {
				status = "Validated"
			}
			sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", status))
		}
		if ls.GrowthHypothesis.Evidence != "" {
			sb.WriteString(fmt.Sprintf("**Evidence:** %s\n\n", ls.GrowthHypothesis.Evidence))
		}
	}

	// MVPs
	if len(ls.MVPs) > 0 {
		sb.WriteString("## MVP Iterations\n\n")
		sb.WriteString("| Name | Type | Status | Goal | Decision |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, mvp := range ls.MVPs {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				mvp.Name, mvp.Type, mvp.Status, truncateStr(mvp.Goal, 40), mvp.Decision))
		}
		sb.WriteString("\n")

		// MVP Details
		for _, mvp := range ls.MVPs {
			if mvp.Results != "" || mvp.Learnings != "" {
				sb.WriteString(fmt.Sprintf("### %s\n\n", mvp.Name))
				if mvp.Description != "" {
					sb.WriteString(mvp.Description + "\n\n")
				}
				if mvp.Results != "" {
					sb.WriteString(fmt.Sprintf("**Results:** %s\n\n", mvp.Results))
				}
				if mvp.Learnings != "" {
					sb.WriteString(fmt.Sprintf("**Learnings:** %s\n\n", mvp.Learnings))
				}
			}
		}
	}

	// Experiments
	if len(ls.Experiments) > 0 {
		sb.WriteString("## Experiments (Build-Measure-Learn)\n\n")
		sb.WriteString("| Hypothesis | Status | Method | Validated |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, exp := range ls.Experiments {
			validated := "⏳"
			if exp.LearnValidated != nil {
				if *exp.LearnValidated {
					validated = "✓"
				} else {
					validated = "✗"
				}
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				truncateStr(exp.LearnHypothesis, 40), exp.Status, exp.MeasureMethod, validated))
		}
		sb.WriteString("\n")

		// Experiment Details
		for _, exp := range ls.Experiments {
			if exp.LearnInsight != "" || exp.LearnDecision != "" {
				sb.WriteString(fmt.Sprintf("### Experiment: %s\n\n", truncateStr(exp.LearnHypothesis, 50)))

				sb.WriteString("#### Build\n\n")
				sb.WriteString(fmt.Sprintf("%s\n\n", exp.BuildDescription))
				if exp.BuildOutput != "" {
					sb.WriteString(fmt.Sprintf("**Output:** %s\n\n", exp.BuildOutput))
				}

				sb.WriteString("#### Measure\n\n")
				sb.WriteString(fmt.Sprintf("**Method:** %s\n\n", exp.MeasureMethod))
				if exp.MeasureTarget != "" {
					sb.WriteString(fmt.Sprintf("**Target:** %s\n\n", exp.MeasureTarget))
				}
				if exp.MeasureActual != "" {
					sb.WriteString(fmt.Sprintf("**Actual:** %s\n\n", exp.MeasureActual))
				}

				sb.WriteString("#### Learn\n\n")
				if exp.LearnInsight != "" {
					sb.WriteString(fmt.Sprintf("**Insight:** %s\n\n", exp.LearnInsight))
				}
				if exp.LearnDecision != "" {
					sb.WriteString(fmt.Sprintf("**Decision:** %s\n\n", exp.LearnDecision))
				}
			}
		}
	}

	// Pivots
	if len(ls.Pivots) > 0 {
		sb.WriteString("## Pivot History\n\n")
		sb.WriteString("| Type | From | To | Reason | Status |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, pivot := range ls.Pivots {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				pivot.Type, truncateStr(pivot.FromState, 25), truncateStr(pivot.ToState, 25),
				truncateStr(pivot.Reason, 30), pivot.Status))
		}
		sb.WriteString("\n")
	}

	// Innovation Accounting Metrics
	if len(ls.Metrics) > 0 {
		sb.WriteString("## Innovation Accounting\n\n")
		sb.WriteString("| Metric | Type | Current | Target | Trend |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, m := range ls.Metrics {
			trend := "-"
			switch m.Trend {
			case "improving":
				trend = "↑"
			case "declining":
				trend = "↓"
			case "stable":
				trend = "→"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				m.Name, m.Type, m.Current, m.Target, trend))
		}
		sb.WriteString("\n")
	}

	// Learning Goals
	if len(ls.LearningGoals) > 0 {
		sb.WriteString("## Current Learning Goals\n\n")
		for _, goal := range ls.LearningGoals {
			sb.WriteString(fmt.Sprintf("- %s\n", goal))
		}
		sb.WriteString("\n")
	}

	return []byte(sb.String()), nil
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func (r *MarkdownRenderer) renderDesignThinking(dt *canvas.DesignThinkingCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + dt.Metadata.Title + "\n\n")

	// Current phase
	sb.WriteString(fmt.Sprintf("**Current Phase:** %s\n\n", dt.CurrentPhase))

	// Phase 1: Empathize
	sb.WriteString("## Phase 1: Empathize\n\n")

	if len(dt.EmpathyMaps) > 0 {
		sb.WriteString("### Empathy Maps\n\n")
		for _, em := range dt.EmpathyMaps {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", em.PersonaName))
			if em.Goal != "" {
				sb.WriteString(fmt.Sprintf("**Goal:** %s\n\n", em.Goal))
			}

			sb.WriteString("| Says | Thinks | Does | Feels |\n")
			sb.WriteString("|---|---|---|---|\n")

			// Get max length
			maxLen := len(em.Says)
			if len(em.Thinks) > maxLen {
				maxLen = len(em.Thinks)
			}
			if len(em.Does) > maxLen {
				maxLen = len(em.Does)
			}
			if len(em.Feels) > maxLen {
				maxLen = len(em.Feels)
			}

			for i := 0; i < maxLen; i++ {
				says := ""
				if i < len(em.Says) {
					says = em.Says[i]
				}
				thinks := ""
				if i < len(em.Thinks) {
					thinks = em.Thinks[i]
				}
				does := ""
				if i < len(em.Does) {
					does = em.Does[i]
				}
				feels := ""
				if i < len(em.Feels) {
					feels = em.Feels[i]
				}
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", says, thinks, does, feels))
			}
			sb.WriteString("\n")

			// Pains and Gains
			if len(em.Pains) > 0 {
				sb.WriteString("**Pains:**\n\n")
				for _, pain := range em.Pains {
					sb.WriteString(fmt.Sprintf("- %s\n", pain))
				}
				sb.WriteString("\n")
			}
			if len(em.Gains) > 0 {
				sb.WriteString("**Gains:**\n\n")
				for _, gain := range em.Gains {
					sb.WriteString(fmt.Sprintf("- %s\n", gain))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(dt.Interviews) > 0 {
		sb.WriteString("### Interviews\n\n")
		sb.WriteString("| Participant | Date | Duration | Key Insights |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, i := range dt.Interviews {
			insights := ""
			if len(i.Insights) > 0 {
				insights = truncateStr(i.Insights[0], 40)
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				i.ParticipantType, i.Date, i.Duration, insights))
		}
		sb.WriteString("\n")
	}

	// Phase 2: Define
	sb.WriteString("## Phase 2: Define\n\n")

	if dt.ProblemStatement != "" {
		sb.WriteString("### Point of View (POV)\n\n")
		sb.WriteString(fmt.Sprintf("> %s\n\n", dt.ProblemStatement))
	}

	if len(dt.UserNeeds) > 0 {
		sb.WriteString("### User Needs\n\n")
		sb.WriteString("| User | Need | Insight | Priority | Validated |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, n := range dt.UserNeeds {
			validated := "No"
			if n.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				n.User, n.Need, truncateStr(n.Insight, 30), n.Priority, validated))
		}
		sb.WriteString("\n")
	}

	if len(dt.HowMightWe) > 0 {
		sb.WriteString("### How Might We Questions\n\n")
		for _, q := range dt.HowMightWe {
			sb.WriteString(fmt.Sprintf("- HMW %s\n", q))
		}
		sb.WriteString("\n")
	}

	// Phase 3: Ideate
	sb.WriteString("## Phase 3: Ideate\n\n")

	if len(dt.IdeationSessions) > 0 {
		sb.WriteString("### Ideation Sessions\n\n")
		sb.WriteString("| Date | Method | HMW Question | Ideas Generated |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, s := range dt.IdeationSessions {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d |\n",
				s.Date, s.Method, truncateStr(s.HMWQuestion, 30), s.IdeaCount))
		}
		sb.WriteString("\n")
	}

	if len(dt.Ideas) > 0 {
		sb.WriteString("### Ideas\n\n")
		sb.WriteString("| Title | Category | Impact | Feasibility | Selected |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, idea := range dt.Ideas {
			selected := ""
			if idea.Selected {
				selected = "★"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				idea.Title, idea.Category, idea.Impact, idea.Feasibility, selected))
		}
		sb.WriteString("\n")
	}

	// Phase 4: Prototype
	sb.WriteString("## Phase 4: Prototype\n\n")

	if len(dt.Prototypes) > 0 {
		sb.WriteString("### Prototypes\n\n")
		sb.WriteString("| Name | Type | Fidelity | Status | Iteration |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, p := range dt.Prototypes {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d |\n",
				p.Name, p.Type, p.Fidelity, p.Status, p.Iteration))
		}
		sb.WriteString("\n")
	}

	// Phase 5: Test
	sb.WriteString("## Phase 5: Test\n\n")

	if len(dt.Tests) > 0 {
		sb.WriteString("### Test Results\n\n")
		for _, t := range dt.Tests {
			sb.WriteString(fmt.Sprintf("#### Test: %s (%s)\n\n", t.ID, t.Method))

			if len(t.Findings) > 0 {
				sb.WriteString("**Findings:**\n\n")
				for _, f := range t.Findings {
					sb.WriteString(fmt.Sprintf("- %s\n", f))
				}
				sb.WriteString("\n")
			}

			if len(t.PositiveFeedback) > 0 {
				sb.WriteString("**Positive Feedback:**\n\n")
				for _, f := range t.PositiveFeedback {
					sb.WriteString(fmt.Sprintf("- ✓ %s\n", f))
				}
				sb.WriteString("\n")
			}

			if len(t.NegativeFeedback) > 0 {
				sb.WriteString("**Negative Feedback:**\n\n")
				for _, f := range t.NegativeFeedback {
					sb.WriteString(fmt.Sprintf("- ✗ %s\n", f))
				}
				sb.WriteString("\n")
			}

			if t.NextIteration != "" {
				sb.WriteString(fmt.Sprintf("**Next Iteration:** %s\n\n", t.NextIteration))
			}

			continueStr := "No"
			if t.ShouldContinue {
				continueStr = "Yes"
			}
			sb.WriteString(fmt.Sprintf("**Continue with this direction:** %s\n\n", continueStr))
		}
	}

	// Summary
	if dt.IterationCount > 0 {
		sb.WriteString("## Summary\n\n")
		sb.WriteString(fmt.Sprintf("**Iteration Count:** %d\n\n", dt.IterationCount))
		sb.WriteString(fmt.Sprintf("**Total Ideas:** %d\n\n", len(dt.Ideas)))
		sb.WriteString(fmt.Sprintf("**Selected Ideas:** %d\n\n", dt.SelectedIdeaCount()))
		sb.WriteString(fmt.Sprintf("**Prototypes:** %d\n\n", len(dt.Prototypes)))
		sb.WriteString(fmt.Sprintf("**Tests:** %d\n\n", len(dt.Tests)))
	}

	return []byte(sb.String()), nil
}

func (r *MarkdownRenderer) renderJTBD(jtbd *canvas.JTBDCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# " + jtbd.Metadata.Title + "\n\n")

	// Main Job
	if jtbd.MainJob != nil {
		sb.WriteString("## Main Job\n\n")
		sb.WriteString(fmt.Sprintf("**Statement:** %s\n\n", jtbd.MainJob.Statement))
		sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", jtbd.MainJob.Type))
		if jtbd.MainJob.Importance != "" {
			sb.WriteString(fmt.Sprintf("**Importance:** %s\n\n", jtbd.MainJob.Importance))
		}
		if jtbd.MainJob.Satisfaction != "" {
			sb.WriteString(fmt.Sprintf("**Current Satisfaction:** %s\n\n", jtbd.MainJob.Satisfaction))
		}
		if jtbd.MainJob.Frequency != "" {
			sb.WriteString(fmt.Sprintf("**Frequency:** %s\n\n", jtbd.MainJob.Frequency))
		}
		if jtbd.MainJob.Context != "" {
			sb.WriteString(fmt.Sprintf("**Context:** %s\n\n", jtbd.MainJob.Context))
		}
	}

	// Related Jobs
	if len(jtbd.RelatedJobs) > 0 {
		sb.WriteString("## Related Jobs\n\n")
		sb.WriteString("| ID | Statement | Type | Importance |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, job := range jtbd.RelatedJobs {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				job.ID, job.Statement, job.Type, job.Importance))
		}
		sb.WriteString("\n")
	}

	// Job Map (Universal Job Map)
	if len(jtbd.JobMap) > 0 {
		sb.WriteString("## Job Map (Universal Job Map)\n\n")
		sb.WriteString("| Stage | Name | Description |\n")
		sb.WriteString("|---|---|---|\n")
		for _, stage := range jtbd.JobMap {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				stage.Stage, stage.Name, truncateStr(stage.Description, 50)))
		}
		sb.WriteString("\n")

		// Stage details
		for _, stage := range jtbd.JobMap {
			if len(stage.Steps) > 0 || len(stage.PainPoints) > 0 || len(stage.Outcomes) > 0 {
				sb.WriteString(fmt.Sprintf("### %s: %s\n\n", stage.Stage, stage.Name))

				if len(stage.Steps) > 0 {
					sb.WriteString("**Steps:**\n\n")
					for _, step := range stage.Steps {
						sb.WriteString(fmt.Sprintf("- %s\n", step))
					}
					sb.WriteString("\n")
				}

				if len(stage.PainPoints) > 0 {
					sb.WriteString("**Pain Points:**\n\n")
					for _, pain := range stage.PainPoints {
						sb.WriteString(fmt.Sprintf("- %s\n", pain))
					}
					sb.WriteString("\n")
				}

				if len(stage.Outcomes) > 0 {
					sb.WriteString("**Desired Outcomes:**\n\n")
					for _, outcome := range stage.Outcomes {
						sb.WriteString(fmt.Sprintf("- %s\n", outcome))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	// Desired Outcomes
	if len(jtbd.DesiredOutcomes) > 0 {
		sb.WriteString("## Desired Outcomes\n\n")
		sb.WriteString("| ID | Direction | Statement | Importance | Satisfaction | Opportunity |\n")
		sb.WriteString("|---|---|---|---|---|---|\n")
		for _, out := range jtbd.DesiredOutcomes {
			opportunity := "-"
			if out.Opportunity > 0 {
				opportunity = fmt.Sprintf("%.1f", out.Opportunity)
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %d | %s |\n",
				out.ID, out.Direction, truncateStr(out.Statement, 40), out.Importance, out.Satisfaction, opportunity))
		}
		sb.WriteString("\n")
	}

	// Undesired Outcomes
	if len(jtbd.UndesiredOutcomes) > 0 {
		sb.WriteString("## Undesired Outcomes (Avoid)\n\n")
		sb.WriteString("| ID | Direction | Statement |\n")
		sb.WriteString("|---|---|---|\n")
		for _, out := range jtbd.UndesiredOutcomes {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", out.ID, out.Direction, out.Statement))
		}
		sb.WriteString("\n")
	}

	// Circumstances
	if len(jtbd.Circumstances) > 0 {
		sb.WriteString("## Job Circumstances\n\n")
		sb.WriteString("| ID | Description | Trigger | Frequency | Urgency |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, circ := range jtbd.Circumstances {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				circ.ID, truncateStr(circ.Description, 35), circ.Trigger, circ.Frequency, circ.Urgency))
		}
		sb.WriteString("\n")
	}

	// Forces Analysis
	hasForces := len(jtbd.PushForces) > 0 || len(jtbd.PullForces) > 0 || len(jtbd.Anxieties) > 0 || len(jtbd.Habits) > 0
	if hasForces {
		sb.WriteString("## Forces Analysis\n\n")
		sb.WriteString("The forces diagram shows what drives switching behavior:\n\n")

		// Push Forces
		if len(jtbd.PushForces) > 0 {
			sb.WriteString("### Push Forces (Away from Current)\n\n")
			sb.WriteString("| Description | Strength | Category |\n")
			sb.WriteString("|---|---|---|\n")
			for _, f := range jtbd.PushForces {
				sb.WriteString(fmt.Sprintf("| ← %s | %s | %s |\n", f.Description, f.Strength, f.Category))
			}
			sb.WriteString("\n")
		}

		// Pull Forces
		if len(jtbd.PullForces) > 0 {
			sb.WriteString("### Pull Forces (Toward New)\n\n")
			sb.WriteString("| Description | Strength | Category |\n")
			sb.WriteString("|---|---|---|\n")
			for _, f := range jtbd.PullForces {
				sb.WriteString(fmt.Sprintf("| → %s | %s | %s |\n", f.Description, f.Strength, f.Category))
			}
			sb.WriteString("\n")
		}

		// Anxieties
		if len(jtbd.Anxieties) > 0 {
			sb.WriteString("### Anxieties (Concerns about Switching)\n\n")
			sb.WriteString("| Description | Strength |\n")
			sb.WriteString("|---|---|\n")
			for _, f := range jtbd.Anxieties {
				sb.WriteString(fmt.Sprintf("| ⚠ %s | %s |\n", f.Description, f.Strength))
			}
			sb.WriteString("\n")
		}

		// Habits
		if len(jtbd.Habits) > 0 {
			sb.WriteString("### Habits (Comfort with Current)\n\n")
			sb.WriteString("| Description | Strength |\n")
			sb.WriteString("|---|---|\n")
			for _, f := range jtbd.Habits {
				sb.WriteString(fmt.Sprintf("| ↻ %s | %s |\n", f.Description, f.Strength))
			}
			sb.WriteString("\n")
		}
	}

	// Hiring/Firing Solutions
	if len(jtbd.HiringSolutions) > 0 {
		sb.WriteString("## Currently Hired Solutions\n\n")
		sb.WriteString("| ID | Name | Type | Why Hired | Satisfaction |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, sol := range jtbd.HiringSolutions {
			satisfaction := "-"
			if sol.SatisfactionLevel > 0 {
				satisfaction = fmt.Sprintf("%d/10", sol.SatisfactionLevel)
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				sol.ID, sol.Name, sol.Type, truncateStr(sol.WhyHired, 30), satisfaction))
		}
		sb.WriteString("\n")
	}

	if len(jtbd.FiringSolutions) > 0 {
		sb.WriteString("## Solutions Being Fired\n\n")
		for _, sol := range jtbd.FiringSolutions {
			sb.WriteString(fmt.Sprintf("- ✗ %s\n", sol))
		}
		sb.WriteString("\n")
	}

	if len(jtbd.CompetingSolutions) > 0 {
		sb.WriteString("## Competing Solutions\n\n")
		for _, sol := range jtbd.CompetingSolutions {
			sb.WriteString(fmt.Sprintf("- %s\n", sol))
		}
		sb.WriteString("\n")
	}

	// Opportunity Scores (ODI)
	if len(jtbd.OpportunityScores) > 0 {
		sb.WriteString("## Opportunity Scores (ODI)\n\n")
		sb.WriteString("Formula: Opportunity = Importance + max(Importance - Satisfaction, 0)\n\n")
		sb.WriteString("| Outcome | Importance | Satisfaction | Opportunity | Segment |\n")
		sb.WriteString("|---|---|---|---|---|\n")
		for _, score := range jtbd.OpportunityScores {
			sb.WriteString(fmt.Sprintf("| %s | %.1f | %.1f | %.1f | %s |\n",
				score.OutcomeRef, score.Importance, score.Satisfaction, score.Opportunity, score.Segment))
		}
		sb.WriteString("\n")
	}

	// Interviews
	if len(jtbd.Interviews) > 0 {
		sb.WriteString("## Switch Interviews\n\n")
		sb.WriteString("| ID | Date | Participant | Switch Context |\n")
		sb.WriteString("|---|---|---|---|\n")
		for _, interview := range jtbd.Interviews {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				interview.ID, interview.Date, interview.ParticipantType, truncateStr(interview.SwitchContext, 40)))
		}
		sb.WriteString("\n")

		// Interview details
		for _, interview := range jtbd.Interviews {
			if len(interview.KeyQuotes) > 0 || len(interview.Insights) > 0 {
				sb.WriteString(fmt.Sprintf("### Interview: %s\n\n", interview.ID))

				if interview.SwitchContext != "" {
					sb.WriteString(fmt.Sprintf("**Switch Context:** %s\n\n", interview.SwitchContext))
				}

				if interview.FirstThought != "" {
					sb.WriteString(fmt.Sprintf("**First Thought:** %s\n\n", interview.FirstThought))
				}

				if len(interview.KeyQuotes) > 0 {
					sb.WriteString("**Key Quotes:**\n\n")
					for _, quote := range interview.KeyQuotes {
						sb.WriteString(fmt.Sprintf("> \"%s\"\n\n", quote))
					}
				}

				if len(interview.Insights) > 0 {
					sb.WriteString("**Insights:**\n\n")
					for _, insight := range interview.Insights {
						sb.WriteString(fmt.Sprintf("- %s\n", insight))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	// Validation Notes
	if jtbd.ValidationNotes != "" {
		sb.WriteString("## Validation Notes\n\n")
		sb.WriteString(jtbd.ValidationNotes + "\n\n")
	}

	return []byte(sb.String()), nil
}

func init() {
	render.Register(NewMarkdownRenderer())
}
