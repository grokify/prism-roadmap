// Package d2 provides D2 diagram language renderers for canvas types.
package d2

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

// D2Renderer renders canvas types to D2 diagram language.
type D2Renderer struct{}

// NewD2Renderer creates a new D2 text renderer.
func NewD2Renderer() *D2Renderer {
	return &D2Renderer{}
}

// Format returns the output format.
func (r *D2Renderer) Format() render.Format {
	return render.FormatD2
}

// FileExtension returns the file extension for D2 files.
func (r *D2Renderer) FileExtension() string {
	return ".d2"
}

// Supports returns true for all canvas types.
func (r *D2Renderer) Supports(canvasType canvas.CanvasType) bool {
	return true
}

// Render converts a canvas to D2 format.
func (r *D2Renderer) Render(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
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
	case canvas.CanvasTypeShapeUpScope:
		return r.renderShapeUpScope(c.ShapeUpScope, opts)
	case canvas.CanvasTypeAssumptionMap:
		return r.renderAssumptionMap(c.AssumptionMap, opts)
	case canvas.CanvasTypeDiscoverySnapshot:
		return r.renderDiscoverySnapshot(c.DiscoverySnapshot, opts)
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

func (r *D2Renderer) renderBMC(bmc *canvas.BusinessModelCanvas, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("# Business Model Canvas: " + bmc.Metadata.Title + "\n\n")

	if opts.GridLayout {
		sb.WriteString("grid-rows: 3\n")
		sb.WriteString("grid-columns: 5\n\n")
	}

	// Key Partners
	sb.WriteString("keyPartners: Key Partners {\n")
	sb.WriteString(r.styleBlock("keyPartners", opts))
	for _, p := range bmc.KeyPartnerships {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(p.ID), d2str(p.Partner)))
	}
	sb.WriteString("}\n\n")

	// Key Activities
	sb.WriteString("keyActivities: Key Activities {\n")
	sb.WriteString(r.styleBlock("keyActivities", opts))
	for _, a := range bmc.KeyActivities {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(a.ID), d2str(a.Name)))
	}
	sb.WriteString("}\n\n")

	// Key Resources
	sb.WriteString("keyResources: Key Resources {\n")
	sb.WriteString(r.styleBlock("keyResources", opts))
	for _, res := range bmc.KeyResources {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(res.ID), d2str(res.Name)))
	}
	sb.WriteString("}\n\n")

	// Value Propositions
	sb.WriteString("valuePropositions: Value Propositions {\n")
	sb.WriteString(r.styleBlock("valuePropositions", opts))
	for _, vp := range bmc.ValuePropositions {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(vp.ID), truncate(vp.Description, 50)))
	}
	sb.WriteString("}\n\n")

	// Customer Relationships
	sb.WriteString("customerRelations: Customer Relationships {\n")
	sb.WriteString(r.styleBlock("customerRelations", opts))
	for _, cr := range bmc.CustomerRelationships {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(cr.ID), d2str(cr.Type)))
	}
	sb.WriteString("}\n\n")

	// Channels
	sb.WriteString("channels: Channels {\n")
	sb.WriteString(r.styleBlock("channels", opts))
	for _, ch := range bmc.Channels {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(ch.ID), d2str(ch.Name)))
	}
	sb.WriteString("}\n\n")

	// Customer Segments
	sb.WriteString("customerSegments: Customer Segments {\n")
	sb.WriteString(r.styleBlock("customerSegments", opts))
	for _, seg := range bmc.CustomerSegments {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(seg.ID), d2str(seg.Name)))
	}
	sb.WriteString("}\n\n")

	// Cost Structure (bottom left)
	sb.WriteString("costStructure: Cost Structure {\n")
	sb.WriteString(r.styleBlock("costStructure", opts))
	for _, cost := range bmc.CostStructure {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(cost.ID), truncate(cost.Description, 40)))
	}
	sb.WriteString("}\n\n")

	// Revenue Streams (bottom right)
	sb.WriteString("revenueStreams: Revenue Streams {\n")
	sb.WriteString(r.styleBlock("revenueStreams", opts))
	for _, rev := range bmc.RevenueStreams {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(rev.ID), truncate(rev.Description, 40)))
	}
	sb.WriteString("}\n\n")

	// Connections
	sb.WriteString("# Connections\n")
	sb.WriteString("keyPartners -> keyActivities\n")
	sb.WriteString("keyActivities -> valuePropositions\n")
	sb.WriteString("keyResources -> valuePropositions\n")
	sb.WriteString("valuePropositions -> customerRelations\n")
	sb.WriteString("valuePropositions -> channels\n")
	sb.WriteString("customerRelations -> customerSegments\n")
	sb.WriteString("channels -> customerSegments\n")
	sb.WriteString("keyPartners -> costStructure\n")
	sb.WriteString("keyActivities -> costStructure\n")
	sb.WriteString("keyResources -> costStructure\n")
	sb.WriteString("valuePropositions -> revenueStreams\n")
	sb.WriteString("customerSegments -> revenueStreams\n")

	return []byte(sb.String()), nil
}

func (r *D2Renderer) renderOST(ost *canvas.OpportunitySolutionTree, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("# Opportunity Solution Tree: " + ost.Metadata.Title + "\n\n")
	sb.WriteString("direction: down\n\n")

	// Outcome (root)
	outcomeID := sanitizeID(ost.Outcome.ID)
	sb.WriteString(fmt.Sprintf("%s: {\n", outcomeID))
	sb.WriteString(fmt.Sprintf("  label: \"%s\"\n", d2str(ost.Outcome.Description)))
	sb.WriteString(r.styleBlock("outcome", opts))
	if ost.Outcome.Metric != "" {
		sb.WriteString(fmt.Sprintf("  metric: \"%s\"\n", d2str(ost.Outcome.Metric)))
	}
	sb.WriteString("}\n\n")

	// Opportunities
	for _, opp := range ost.Outcome.Opportunities {
		oppID := sanitizeID(opp.ID)
		sb.WriteString(fmt.Sprintf("%s: {\n", oppID))
		sb.WriteString(fmt.Sprintf("  label: \"%s\"\n", truncate(opp.Description, 60)))
		sb.WriteString(r.styleBlock("opportunity", opts))
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("%s -> %s\n\n", outcomeID, oppID))

		// Solutions
		for _, sol := range opp.Solutions {
			solID := sanitizeID(sol.ID)
			sb.WriteString(fmt.Sprintf("%s: {\n", solID))
			sb.WriteString(fmt.Sprintf("  label: \"%s\"\n", truncate(sol.Description, 50)))
			sb.WriteString(r.styleBlock("solution", opts))
			if sol.Status != "" {
				sb.WriteString(fmt.Sprintf("  status: \"%s\"\n", d2str(sol.Status)))
			}
			sb.WriteString("}\n")
			sb.WriteString(fmt.Sprintf("%s -> %s\n\n", oppID, solID))

			// Experiments
			for _, exp := range sol.Experiments {
				expID := sanitizeID(exp.ID)
				sb.WriteString(fmt.Sprintf("%s: {\n", expID))
				sb.WriteString(fmt.Sprintf("  label: \"%s\"\n", truncate(exp.Hypothesis, 40)))
				sb.WriteString(r.styleBlock("experiment", opts))
				if exp.Status != "" {
					sb.WriteString(fmt.Sprintf("  status: \"%s\"\n", d2str(exp.Status)))
				}
				sb.WriteString("}\n")
				sb.WriteString(fmt.Sprintf("%s -> %s\n\n", solID, expID))
			}
		}
	}

	return []byte(sb.String()), nil
}

func (r *D2Renderer) renderOpportunity(opp *canvas.OpportunityCanvas, opts *render.Options) ([]byte, error) {
	if opts.GridLayout {
		return r.renderOpportunityGrid(opp, opts)
	}
	return r.renderOpportunityFlow(opp, opts)
}

// renderOpportunityGrid renders Opportunity Canvas in BMC-style 3x3 grid + budget row.
func (r *D2Renderer) renderOpportunityGrid(opp *canvas.OpportunityCanvas, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Opportunity Canvas (Grid): " + opp.Metadata.Title + "\n\n")
	sb.WriteString("grid-rows: 4\n")
	sb.WriteString("grid-columns: 3\n\n")

	// Row 1: Users & Customers | Problems | Solution Ideas
	sb.WriteString("users: {\n")
	sb.WriteString("  label: \"Users & Customers\\nWho has the problem\"\n")
	sb.WriteString(r.styleBlock("users", opts))
	for _, u := range opp.Users {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(u.ID), d2str(u.Name)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("problems: {\n")
	sb.WriteString("  label: \"Problems\\nPains to address\"\n")
	sb.WriteString(r.styleBlock("problems", opts))
	for _, p := range opp.Problems {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(p.ID), truncate(p.Description, 40)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("solutionIdeas: {\n")
	sb.WriteString("  label: \"Solution Ideas\\nWays to solve it\"\n")
	sb.WriteString(r.styleBlock("solutionIdeas", opts))
	for i, s := range opp.SolutionIdeas {
		sb.WriteString(fmt.Sprintf("  si%d: \"%s\"\n", i+1, truncate(s, 40)))
	}
	sb.WriteString("}\n\n")

	// Row 2: Solutions Today | User Value | Adoption Strategy
	sb.WriteString("currentSolutions: {\n")
	sb.WriteString("  label: \"Solutions Today\\nCurrent workarounds\"\n")
	sb.WriteString(r.styleBlock("currentSolutions", opts))
	for _, s := range opp.CurrentSolutions {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(s.ID), d2str(s.Name)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("userValue: {\n")
	sb.WriteString("  label: \"User Value\\nBenefit to users\"\n")
	sb.WriteString(r.styleBlock("userValue", opts))
	for i, v := range opp.UserValue {
		sb.WriteString(fmt.Sprintf("  uv%d: \"%s\"\n", i+1, truncate(v, 40)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("adoptionStrategy: {\n")
	sb.WriteString("  label: \"Adoption Strategy\\nHow they'll find it\"\n")
	sb.WriteString(r.styleBlock("adoptionStrategy", opts))
	for i, a := range opp.AdoptionStrategy {
		sb.WriteString(fmt.Sprintf("  as%d: \"%s\"\n", i+1, truncate(a, 40)))
	}
	sb.WriteString("}\n\n")

	// Row 3: User Metrics | Business Problem | Business Metrics
	sb.WriteString("userMetrics: {\n")
	sb.WriteString("  label: \"User Metrics\\nBehaviour to track\"\n")
	sb.WriteString(r.styleBlock("userMetrics", opts))
	for i, m := range opp.UserMetrics {
		sb.WriteString(fmt.Sprintf("  um%d: \"%s\"\n", i+1, truncate(m, 40)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("businessProblem: {\n")
	sb.WriteString("  label: \"Business Problem\\nWhy it matters to us\"\n")
	sb.WriteString(r.styleBlock("businessProblem", opts))
	if opp.BusinessProblem != "" {
		sb.WriteString(fmt.Sprintf("  bp: \"%s\"\n", truncate(opp.BusinessProblem, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("businessMetrics: {\n")
	sb.WriteString("  label: \"Business Metrics\\nOutcome to measure\"\n")
	sb.WriteString(r.styleBlock("businessMetrics", opts))
	for i, m := range opp.BusinessMetrics {
		sb.WriteString(fmt.Sprintf("  bm%d: \"%s\"\n", i+1, truncate(m, 40)))
	}
	sb.WriteString("}\n\n")

	// Row 4: Budget (colspan 3)
	sb.WriteString("budget: {\n")
	sb.WriteString("  grid-column-span: 3\n")
	sb.WriteString("  label: \"Budget\\nWhat you're willing to invest to learn whether this is worth building\"\n")
	sb.WriteString(r.styleBlock("budget", opts))
	if opp.Budget != nil {
		var budgetItems []string
		if opp.Budget.TimeEstimate != "" {
			budgetItems = append(budgetItems, "Time: "+opp.Budget.TimeEstimate)
		}
		if opp.Budget.TeamSize != "" {
			budgetItems = append(budgetItems, "Team: "+opp.Budget.TeamSize)
		}
		if opp.Budget.CostEstimate != "" {
			budgetItems = append(budgetItems, "Cost: "+opp.Budget.CostEstimate)
		}
		if opp.Budget.Constraints != "" {
			budgetItems = append(budgetItems, "Constraints: "+opp.Budget.Constraints)
		}
		if len(budgetItems) > 0 {
			sb.WriteString(fmt.Sprintf("  details: \"%s\"\n", escapeD2(strings.Join(budgetItems, " | "))))
		}
	}
	sb.WriteString("}\n")

	return []byte(sb.String()), nil
}

// renderOpportunityFlow renders Opportunity Canvas with arrows showing flow.
func (r *D2Renderer) renderOpportunityFlow(opp *canvas.OpportunityCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Opportunity Canvas (Flow): " + opp.Metadata.Title + "\n\n")
	sb.WriteString("direction: down\n\n")

	// Row 1: Problem Space
	sb.WriteString("problems: Problems {\n")
	for _, p := range opp.Problems {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(p.ID), truncate(p.Description, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("users: Users {\n")
	for _, u := range opp.Users {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(u.ID), d2str(u.Name)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("currentSolutions: Current Solutions {\n")
	for _, s := range opp.CurrentSolutions {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(s.ID), d2str(s.Name)))
	}
	sb.WriteString("}\n\n")

	// Row 2: Solution Space
	sb.WriteString("valueProposition: Value Proposition {\n")
	sb.WriteString(fmt.Sprintf("  statement: \"%s\"\n", truncate(opp.ValueProposition.Statement, 60)))
	sb.WriteString("}\n\n")

	sb.WriteString("userValue: User Value {\n")
	for i, v := range opp.UserValue {
		sb.WriteString(fmt.Sprintf("  uv%d: \"%s\"\n", i+1, truncate(v, 40)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("businessValue: Business Value {\n")
	for i, v := range opp.BusinessValue {
		sb.WriteString(fmt.Sprintf("  bv%d: \"%s\"\n", i+1, truncate(v, 40)))
	}
	sb.WriteString("}\n\n")

	// Row 3: Validation
	sb.WriteString("assumptions: Assumptions {\n")
	for _, a := range opp.Assumptions {
		validated := ""
		if a.Validated {
			validated = " [validated]"
		}
		sb.WriteString(fmt.Sprintf("  %s: \"%s%s\"\n", sanitizeID(a.ID), truncate(a.Description, 40), validated))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("risks: Risks {\n")
	for _, risk := range opp.Risks {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(risk.ID), truncate(risk.Description, 40)))
	}
	sb.WriteString("}\n\n")

	if opp.Budget != nil {
		sb.WriteString("budget: Budget {\n")
		if opp.Budget.TimeEstimate != "" {
			sb.WriteString(fmt.Sprintf("  time: \"%s\"\n", escapeD2(opp.Budget.TimeEstimate)))
		}
		if opp.Budget.CostEstimate != "" {
			sb.WriteString(fmt.Sprintf("  cost: \"%s\"\n", escapeD2(opp.Budget.CostEstimate)))
		}
		sb.WriteString("}\n\n")
	}

	// Connections (arrows)
	sb.WriteString("# Flow\n")
	sb.WriteString("problems -> valueProposition\n")
	sb.WriteString("users -> valueProposition\n")
	sb.WriteString("valueProposition -> userValue\n")
	sb.WriteString("valueProposition -> businessValue\n")
	sb.WriteString("assumptions -> risks\n")

	return []byte(sb.String()), nil
}

func (r *D2Renderer) renderFeature(fc *canvas.FeatureCanvas, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Feature Canvas: " + fc.Metadata.Title + "\n\n")

	if opts.GridLayout {
		sb.WriteString("grid-rows: 4\n")
		sb.WriteString("grid-columns: 2\n\n")
	}

	// Top banner
	sb.WriteString("ideaStatement: {\n")
	sb.WriteString(fmt.Sprintf("  label: \"%s\"\n", d2str(fc.IdeaStatement)))
	sb.WriteString(r.styleBlock("ideaStatement", opts))
	sb.WriteString("}\n\n")

	// Left side - Problem Area
	sb.WriteString("situations: Situations {\n")
	sb.WriteString(r.styleBlock("situations", opts))
	for _, s := range fc.Situations {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(s.ID), truncate(s.Description, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("problems: Problems {\n")
	sb.WriteString(r.styleBlock("problems", opts))
	for _, p := range fc.Problems {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(p.ID), truncate(p.Description, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("value: Value {\n")
	sb.WriteString(r.styleBlock("value", opts))
	for _, v := range fc.Value {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(v.ID), truncate(v.Description, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("capabilities: Capabilities {\n")
	sb.WriteString(r.styleBlock("capabilities", opts))
	for _, c := range fc.Capabilities {
		priority := ""
		if c.Priority != "" {
			priority = fmt.Sprintf(" [%s]", c.Priority)
		}
		sb.WriteString(fmt.Sprintf("  %s: \"%s%s\"\n", sanitizeID(c.ID), truncate(c.Description, 45), priority))
	}
	sb.WriteString("}\n\n")

	// Right side - Constraints
	sb.WriteString("restrictions: Restrictions {\n")
	sb.WriteString(r.styleBlock("restrictions", opts))
	for i, rest := range fc.Restrictions {
		sb.WriteString(fmt.Sprintf("  r%d: \"%s\"\n", i+1, truncate(rest, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("limitations: Limitations {\n")
	sb.WriteString(r.styleBlock("limitations", opts))
	for i, lim := range fc.Limitations {
		sb.WriteString(fmt.Sprintf("  l%d: \"%s\"\n", i+1, truncate(lim, 50)))
	}
	sb.WriteString("}\n\n")

	// Flow
	sb.WriteString("# Flow\n")
	sb.WriteString("ideaStatement -> situations\n")
	sb.WriteString("situations -> problems\n")
	sb.WriteString("problems -> value\n")
	sb.WriteString("value -> capabilities\n")
	sb.WriteString("capabilities -> restrictions\n")
	sb.WriteString("capabilities -> limitations\n")

	return []byte(sb.String()), nil
}

func (r *D2Renderer) renderLeanUX(lux *canvas.LeanUXCanvas, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Lean UX Canvas: " + lux.Metadata.Title + "\n\n")

	if opts.GridLayout {
		sb.WriteString("grid-rows: 3\n")
		sb.WriteString("grid-columns: 3\n\n")
	}

	// Top row - Business context
	sb.WriteString("businessProblem: Business Problem {\n")
	sb.WriteString(fmt.Sprintf("  problem: \"%s\"\n", truncate(lux.BusinessProblem, 60)))
	sb.WriteString("}\n\n")

	sb.WriteString("businessOutcomes: Business Outcomes {\n")
	for _, o := range lux.BusinessOutcomes {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(o.ID), truncate(o.Description, 50)))
	}
	sb.WriteString("}\n\n")

	// Middle row - User context
	sb.WriteString("users: Users {\n")
	for _, u := range lux.Users {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(u.ID), d2str(u.Name)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("userOutcomes: User Outcomes {\n")
	for _, o := range lux.UserOutcomes {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(o.ID), truncate(o.Description, 50)))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("solutions: Solutions {\n")
	for _, s := range lux.Solutions {
		sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(s.ID), truncate(s.Description, 50)))
	}
	sb.WriteString("}\n\n")

	// Bottom row - Validation
	sb.WriteString("hypotheses: Hypotheses {\n")
	for _, h := range lux.Hypotheses {
		validated := ""
		if h.Validated != nil {
			if *h.Validated {
				validated = " [validated]"
			} else {
				validated = " [invalidated]"
			}
		}
		sb.WriteString(fmt.Sprintf("  %s: \"We believe %s%s\"\n", sanitizeID(h.ID), truncate(h.WeBelieve, 40), validated))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("riskiestAssumption: Riskiest Assumption {\n")
	sb.WriteString(fmt.Sprintf("  assumption: \"%s\"\n", truncate(lux.RiskiestAssumption, 60)))
	sb.WriteString("}\n\n")

	if lux.Experiment != nil {
		sb.WriteString("experiment: Experiment {\n")
		sb.WriteString(fmt.Sprintf("  description: \"%s\"\n", truncate(lux.Experiment.Description, 50)))
		sb.WriteString(fmt.Sprintf("  method: \"%s\"\n", d2str(lux.Experiment.Method)))
		sb.WriteString(fmt.Sprintf("  status: \"%s\"\n", d2str(string(lux.Experiment.Status))))
		sb.WriteString("}\n\n")
	}

	// Flow
	sb.WriteString("# Flow\n")
	sb.WriteString("businessProblem -> businessOutcomes\n")
	sb.WriteString("users -> userOutcomes\n")
	sb.WriteString("userOutcomes -> solutions\n")
	sb.WriteString("solutions -> hypotheses\n")
	sb.WriteString("hypotheses -> riskiestAssumption\n")
	if lux.Experiment != nil {
		sb.WriteString("riskiestAssumption -> experiment\n")
	}

	return []byte(sb.String()), nil
}

// renderShapeUpPitch renders a Shape Up Pitch.
func (r *D2Renderer) renderShapeUpPitch(pitch *canvas.ShapeUpPitch, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Shape Up Pitch: " + pitch.Metadata.Title + "\n\n")
	sb.WriteString("direction: down\n\n")

	// Problem section
	sb.WriteString("problem: {\n")
	sb.WriteString("  label: \"Problem\"\n")
	sb.WriteString(r.styleBlock("problem", opts))
	if pitch.Problem.RawIdea != "" {
		sb.WriteString(fmt.Sprintf("  rawIdea: \"%s\"\n", truncate(pitch.Problem.RawIdea, 60)))
	}
	if pitch.Problem.Statement != "" {
		sb.WriteString(fmt.Sprintf("  statement: \"%s\"\n", truncate(pitch.Problem.Statement, 60)))
	}
	if pitch.Problem.WhyNow != "" {
		sb.WriteString(fmt.Sprintf("  whyNow: \"%s\"\n", truncate(pitch.Problem.WhyNow, 50)))
	}
	sb.WriteString("}\n\n")

	// Appetite section
	sb.WriteString("appetite: {\n")
	sb.WriteString("  label: \"Appetite\"\n")
	sb.WriteString(r.styleBlock("appetite", opts))
	sb.WriteString(fmt.Sprintf("  weeks: \"%d weeks (%s)\"\n", pitch.Appetite.Weeks, pitch.Appetite.Size))
	if pitch.Appetite.Rationale != "" {
		sb.WriteString(fmt.Sprintf("  rationale: \"%s\"\n", truncate(pitch.Appetite.Rationale, 50)))
	}
	sb.WriteString("}\n\n")

	// Solution section
	sb.WriteString("solution: {\n")
	sb.WriteString("  label: \"Solution\"\n")
	sb.WriteString(r.styleBlock("solution", opts))
	if pitch.Solution.Approach != "" {
		sb.WriteString(fmt.Sprintf("  approach: \"%s\"\n", truncate(pitch.Solution.Approach, 60)))
	}
	for i, must := range pitch.Solution.MustInclude {
		sb.WriteString(fmt.Sprintf("  must%d: \"✓ %s\"\n", i+1, truncate(must, 50)))
	}
	for i, nice := range pitch.Solution.NiceToHave {
		sb.WriteString(fmt.Sprintf("  nice%d: \"○ %s\"\n", i+1, truncate(nice, 50)))
	}
	sb.WriteString("}\n\n")

	// Rabbit Holes
	if len(pitch.RabbitHoles) > 0 {
		sb.WriteString("rabbitHoles: {\n")
		sb.WriteString("  label: \"Rabbit Holes\"\n")
		sb.WriteString(r.styleBlock("rabbitHoles", opts))
		for i, rh := range pitch.RabbitHoles {
			sb.WriteString(fmt.Sprintf("  rh%d: \"⚠ %s\"\n", i+1, truncate(rh.Description, 50)))
		}
		sb.WriteString("}\n\n")
	}

	// No-Gos
	if len(pitch.NoGos) > 0 {
		sb.WriteString("noGos: {\n")
		sb.WriteString("  label: \"No-Gos\"\n")
		sb.WriteString(r.styleBlock("noGos", opts))
		for i, ng := range pitch.NoGos {
			sb.WriteString(fmt.Sprintf("  ng%d: \"✗ %s\"\n", i+1, truncate(ng, 50)))
		}
		sb.WriteString("}\n\n")
	}

	// Flow
	sb.WriteString("# Flow\n")
	sb.WriteString("problem -> appetite\n")
	sb.WriteString("appetite -> solution\n")
	if len(pitch.RabbitHoles) > 0 {
		sb.WriteString("solution -> rabbitHoles\n")
	}
	if len(pitch.NoGos) > 0 {
		sb.WriteString("solution -> noGos\n")
	}

	return []byte(sb.String()), nil
}

// renderShapeUpScope renders a Shape Up Scope with hill chart visualization.
func (r *D2Renderer) renderShapeUpScope(scope *canvas.ShapeUpScope, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Shape Up Scope: " + scope.Metadata.Title + "\n\n")

	// Hill chart visualization
	sb.WriteString("hill: {\n")
	sb.WriteString("  label: \"Hill Chart\"\n")
	sb.WriteString("  style.fill: \"#f5f5f5\"\n")

	// Create phases
	sb.WriteString("  uphill: {\n")
	sb.WriteString("    label: \"Figuring Out (Uphill)\"\n")
	sb.WriteString("    style.fill: \"#e3f2fd\"\n")
	sb.WriteString("  }\n")
	sb.WriteString("  top: {\n")
	sb.WriteString("    label: \"Figured Out (Top)\"\n")
	sb.WriteString("    style.fill: \"#c8e6c9\"\n")
	sb.WriteString("  }\n")
	sb.WriteString("  downhill: {\n")
	sb.WriteString("    label: \"Executing (Downhill)\"\n")
	sb.WriteString("    style.fill: \"#fff3e0\"\n")
	sb.WriteString("  }\n")
	sb.WriteString("}\n\n")

	// Place scopes on the hill
	for _, s := range scope.Scopes {
		scopeID := sanitizeID(s.ID)
		sb.WriteString(fmt.Sprintf("%s: {\n", scopeID))
		sb.WriteString(fmt.Sprintf("  label: \"%s\"\n", truncate(s.Name, 40)))

		// Determine hill position
		position := s.HillPosition
		if position == 0 {
			position = 50 // Default to middle
		}

		var phase string
		if position < 40 {
			phase = "Uphill"
			sb.WriteString("  style.fill: \"#bbdefb\"\n")
		} else if position < 60 {
			phase = "Top"
			sb.WriteString("  style.fill: \"#a5d6a7\"\n")
		} else {
			phase = "Downhill"
			sb.WriteString("  style.fill: \"#ffcc80\"\n")
		}
		sb.WriteString(fmt.Sprintf("  position: \"%d%% (%s)\"\n", position, phase))

		// Show done status
		doneCount := 0
		for _, t := range s.Tasks {
			if t.Done {
				doneCount++
			}
		}
		if len(s.Tasks) > 0 {
			sb.WriteString(fmt.Sprintf("  tasks: \"%d/%d done\"\n", doneCount, len(s.Tasks)))
		}

		sb.WriteString("}\n\n")
	}

	// Connect scopes to hill phases
	for _, s := range scope.Scopes {
		scopeID := sanitizeID(s.ID)
		position := s.HillPosition
		if position == 0 {
			position = 50
		}

		if position < 40 {
			sb.WriteString(fmt.Sprintf("hill.uphill -> %s: {\n  style.stroke-dash: 3\n}\n", scopeID))
		} else if position < 60 {
			sb.WriteString(fmt.Sprintf("hill.top -> %s: {\n  style.stroke-dash: 3\n}\n", scopeID))
		} else {
			sb.WriteString(fmt.Sprintf("hill.downhill -> %s: {\n  style.stroke-dash: 3\n}\n", scopeID))
		}
	}

	return []byte(sb.String()), nil
}

// renderAssumptionMap renders a Continuous Discovery Assumption Map.
func (r *D2Renderer) renderAssumptionMap(am *canvas.AssumptionMap, opts *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Assumption Map: " + am.Metadata.Title + "\n\n")

	if opts.GridLayout {
		sb.WriteString("grid-rows: 2\n")
		sb.WriteString("grid-columns: 5\n\n")
	}

	// Desirability column
	sb.WriteString("desirability: {\n")
	sb.WriteString("  label: \"Desirability\\n(Will they want it?)\"\n")
	sb.WriteString("  style.fill: \"#e3f2fd\"\n")
	for i, a := range am.Desirability {
		status := r.assumptionStatus(a)
		sb.WriteString(fmt.Sprintf("  d%d: \"%s%s\"\n", i+1, truncate(a.Description, 40), status))
	}
	sb.WriteString("}\n\n")

	// Viability column
	sb.WriteString("viability: {\n")
	sb.WriteString("  label: \"Viability\\n(Will it work for us?)\"\n")
	sb.WriteString("  style.fill: \"#f3e5f5\"\n")
	for i, a := range am.Viability {
		status := r.assumptionStatus(a)
		sb.WriteString(fmt.Sprintf("  v%d: \"%s%s\"\n", i+1, truncate(a.Description, 40), status))
	}
	sb.WriteString("}\n\n")

	// Feasibility column
	sb.WriteString("feasibility: {\n")
	sb.WriteString("  label: \"Feasibility\\n(Can we build it?)\"\n")
	sb.WriteString("  style.fill: \"#e8f5e9\"\n")
	for i, a := range am.Feasibility {
		status := r.assumptionStatus(a)
		sb.WriteString(fmt.Sprintf("  f%d: \"%s%s\"\n", i+1, truncate(a.Description, 40), status))
	}
	sb.WriteString("}\n\n")

	// Usability column
	sb.WriteString("usability: {\n")
	sb.WriteString("  label: \"Usability\\n(Can they use it?)\"\n")
	sb.WriteString("  style.fill: \"#fff3e0\"\n")
	for i, a := range am.Usability {
		status := r.assumptionStatus(a)
		sb.WriteString(fmt.Sprintf("  u%d: \"%s%s\"\n", i+1, truncate(a.Description, 40), status))
	}
	sb.WriteString("}\n\n")

	// Ethical column
	sb.WriteString("ethical: {\n")
	sb.WriteString("  label: \"Ethical\\n(Should we build it?)\"\n")
	sb.WriteString("  style.fill: \"#ffebee\"\n")
	for i, a := range am.Ethical {
		status := r.assumptionStatus(a)
		sb.WriteString(fmt.Sprintf("  e%d: \"%s%s\"\n", i+1, truncate(a.Description, 40), status))
	}
	sb.WriteString("}\n\n")

	// Test first section (high importance, low confidence)
	highRisk := am.HighRiskAssumptions()
	if len(highRisk) > 0 {
		sb.WriteString("testFirst: {\n")
		sb.WriteString("  label: \"Test First (High Importance + Low Confidence)\"\n")
		sb.WriteString("  style.fill: \"#ffcdd2\"\n")
		for i, a := range highRisk {
			sb.WriteString(fmt.Sprintf("  tf%d: \"%s [%s]\"\n", i+1, truncate(a.Description, 35), a.Type))
		}
		sb.WriteString("}\n")
	}

	return []byte(sb.String()), nil
}

func (r *D2Renderer) assumptionStatus(a canvas.CDAssumption) string {
	if a.Validated {
		return " ✓"
	}
	if a.Importance == "high" && a.Confidence == "low" {
		return " ⚠"
	}
	return ""
}

// renderDiscoverySnapshot renders a Continuous Discovery weekly snapshot.
func (r *D2Renderer) renderDiscoverySnapshot(ds *canvas.DiscoverySnapshot, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Discovery Snapshot: " + ds.Metadata.Title + "\n\n")
	sb.WriteString(fmt.Sprintf("# Week: %s\n\n", ds.Week))
	sb.WriteString("direction: down\n\n")

	// Interviews section
	if len(ds.Interviews) > 0 {
		sb.WriteString("interviews: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Interviews (%d)\"\n", len(ds.Interviews)))
		sb.WriteString("  style.fill: \"#e3f2fd\"\n")
		for i, int := range ds.Interviews {
			sb.WriteString(fmt.Sprintf("  i%d: \"%s - %d stories\"\n", i+1, int.ParticipantType, len(int.Stories)))
		}
		sb.WriteString("}\n\n")
	}

	// Opportunities discovered
	if len(ds.OpportunitiesDiscovered) > 0 {
		sb.WriteString("opportunities: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Opportunities (%d)\"\n", len(ds.OpportunitiesDiscovered)))
		sb.WriteString("  style.fill: \"#e8f5e9\"\n")
		for i, opp := range ds.OpportunitiesDiscovered {
			sb.WriteString(fmt.Sprintf("  o%d: \"[%s] %s\"\n", i+1, opp.Action, truncate(opp.Description, 40)))
		}
		sb.WriteString("}\n\n")
	}

	// Assumption tests
	if len(ds.AssumptionTests) > 0 {
		sb.WriteString("tests: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Assumption Tests (%d)\"\n", len(ds.AssumptionTests)))
		sb.WriteString("  style.fill: \"#fff3e0\"\n")
		for i, test := range ds.AssumptionTests {
			status := test.Status
			sb.WriteString(fmt.Sprintf("  t%d: \"[%s] %s\"\n", i+1, status, truncate(test.Hypothesis, 40)))
		}
		sb.WriteString("}\n\n")
	}

	// Key learnings
	if len(ds.KeyLearnings) > 0 {
		sb.WriteString("learnings: {\n")
		sb.WriteString("  label: \"Key Learnings\"\n")
		sb.WriteString("  style.fill: \"#f3e5f5\"\n")
		for i, learning := range ds.KeyLearnings {
			sb.WriteString(fmt.Sprintf("  l%d: \"%s\"\n", i+1, truncate(learning, 50)))
		}
		sb.WriteString("}\n\n")
	}

	// Decisions
	if len(ds.Decisions) > 0 {
		sb.WriteString("decisions: {\n")
		sb.WriteString("  label: \"Decisions Made\"\n")
		sb.WriteString("  style.fill: \"#c8e6c9\"\n")
		for i, dec := range ds.Decisions {
			sb.WriteString(fmt.Sprintf("  d%d: \"%s\"\n", i+1, truncate(dec.Description, 50)))
		}
		sb.WriteString("}\n\n")
	}

	// Flow
	sb.WriteString("# Flow\n")
	if len(ds.Interviews) > 0 {
		if len(ds.OpportunitiesDiscovered) > 0 {
			sb.WriteString("interviews -> opportunities\n")
		}
		if len(ds.KeyLearnings) > 0 {
			sb.WriteString("interviews -> learnings\n")
		}
	}
	if len(ds.AssumptionTests) > 0 && len(ds.KeyLearnings) > 0 {
		sb.WriteString("tests -> learnings\n")
	}
	if len(ds.KeyLearnings) > 0 && len(ds.Decisions) > 0 {
		sb.WriteString("learnings -> decisions\n")
	}

	return []byte(sb.String()), nil
}

func (r *D2Renderer) styleBlock(name string, opts *render.Options) string {
	if color, ok := opts.Colors[name]; ok {
		return fmt.Sprintf("  style.fill: \"%s\"\n", color)
	}
	return ""
}

// sanitizeID makes a string safe for use as a D2 identifier.
func sanitizeID(id string) string {
	// Replace spaces and special characters
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	return id
}

// truncate shortens a string to max length, adding ellipsis if needed.
// Also escapes special D2 characters.
func truncate(s string, max int) string {
	s = escapeD2(s)
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// escapeD2 escapes special characters for D2 labels.
func escapeD2(s string) string {
	// Escape $ to prevent D2 substitution
	s = strings.ReplaceAll(s, "$", "\\$")
	// Escape backticks
	s = strings.ReplaceAll(s, "`", "\\`")
	return s
}

// d2str escapes a string for D2 without truncation.
func d2str(s string) string {
	return escapeD2(s)
}

// renderLeanStartup renders a Lean Startup canvas with Build-Measure-Learn loop.
func (r *D2Renderer) renderLeanStartup(ls *canvas.LeanStartupCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Lean Startup Canvas: " + ls.Metadata.Title + "\n\n")
	sb.WriteString("direction: down\n\n")

	// Vision and Strategy
	sb.WriteString("vision: {\n")
	sb.WriteString("  label: \"Vision & Strategy\"\n")
	sb.WriteString("  style.fill: \"#e3f2fd\"\n")
	if ls.Vision != "" {
		sb.WriteString(fmt.Sprintf("  vision: \"%s\"\n", truncate(ls.Vision, 60)))
	}
	if ls.TargetCustomer != "" {
		sb.WriteString(fmt.Sprintf("  customer: \"Target: %s\"\n", truncate(ls.TargetCustomer, 50)))
	}
	if ls.ProblemHypothesis != "" {
		sb.WriteString(fmt.Sprintf("  problem: \"Problem: %s\"\n", truncate(ls.ProblemHypothesis, 50)))
	}
	sb.WriteString("}\n\n")

	// Hypotheses
	sb.WriteString("hypotheses: {\n")
	sb.WriteString("  label: \"Core Hypotheses\"\n")
	sb.WriteString("  style.fill: \"#fff3e0\"\n")
	if ls.ValueHypothesis != nil {
		status := "⏳"
		if ls.ValueHypothesis.Validated != nil {
			if *ls.ValueHypothesis.Validated {
				status = "✓"
			} else {
				status = "✗"
			}
		}
		sb.WriteString(fmt.Sprintf("  value: \"Value [%s]: %s\"\n", status, truncate(ls.ValueHypothesis.Statement, 45)))
	}
	if ls.GrowthHypothesis != nil {
		status := "⏳"
		if ls.GrowthHypothesis.Validated != nil {
			if *ls.GrowthHypothesis.Validated {
				status = "✓"
			} else {
				status = "✗"
			}
		}
		sb.WriteString(fmt.Sprintf("  growth: \"Growth [%s]: %s (%s)\"\n", status, truncate(ls.GrowthHypothesis.Statement, 35), ls.GrowthHypothesis.GrowthModel))
	}
	sb.WriteString("}\n\n")

	// MVP iterations
	if len(ls.MVPs) > 0 {
		sb.WriteString("mvps: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"MVP Iterations (%d)\"\n", len(ls.MVPs)))
		sb.WriteString("  style.fill: \"#e8f5e9\"\n")
		for _, mvp := range ls.MVPs {
			statusIcon := "📋"
			switch mvp.Status {
			case "building":
				statusIcon = "🔨"
			case "measuring":
				statusIcon = "📊"
			case "learning":
				statusIcon = "💡"
			case "complete":
				statusIcon = "✓"
			}
			sb.WriteString(fmt.Sprintf("  %s: \"%s %s (%s)\"\n", sanitizeID(mvp.ID), statusIcon, truncate(mvp.Name, 35), mvp.Type))
		}
		sb.WriteString("}\n\n")
	}

	// Build-Measure-Learn cycle
	sb.WriteString("bml: {\n")
	sb.WriteString("  label: \"Build-Measure-Learn Loop\"\n")
	sb.WriteString("  style.fill: \"#f3e5f5\"\n")
	sb.WriteString("  build: Build { style.fill: \"#bbdefb\" }\n")
	sb.WriteString("  measure: Measure { style.fill: \"#c8e6c9\" }\n")
	sb.WriteString("  learn: Learn { style.fill: \"#ffe0b2\" }\n")
	sb.WriteString("  build -> measure -> learn -> build\n")
	sb.WriteString("}\n\n")

	// Experiments
	if len(ls.Experiments) > 0 {
		sb.WriteString("experiments: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Experiments (%d)\"\n", len(ls.Experiments)))
		sb.WriteString("  style.fill: \"#e1f5fe\"\n")
		for _, exp := range ls.Experiments {
			status := string(exp.Status)
			sb.WriteString(fmt.Sprintf("  %s: \"[%s] %s\"\n", sanitizeID(exp.ID), status, truncate(exp.LearnHypothesis, 40)))
		}
		sb.WriteString("}\n\n")
	}

	// Pivots
	if len(ls.Pivots) > 0 {
		sb.WriteString("pivots: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Pivot History (%d)\"\n", len(ls.Pivots)))
		sb.WriteString("  style.fill: \"#ffcdd2\"\n")
		for _, pivot := range ls.Pivots {
			sb.WriteString(fmt.Sprintf("  %s: \"%s: %s → %s\"\n", sanitizeID(pivot.ID), pivot.Type, truncate(pivot.FromState, 20), truncate(pivot.ToState, 20)))
		}
		sb.WriteString("}\n\n")
	}

	// Innovation Accounting Metrics
	if len(ls.Metrics) > 0 {
		sb.WriteString("metrics: {\n")
		sb.WriteString("  label: \"Innovation Accounting\"\n")
		sb.WriteString("  style.fill: \"#dcedc8\"\n")
		for _, m := range ls.Metrics {
			trend := ""
			if m.Trend != "" {
				switch m.Trend {
				case "improving":
					trend = " ↑"
				case "declining":
					trend = " ↓"
				}
			}
			sb.WriteString(fmt.Sprintf("  %s: \"%s: %s%s\"\n", sanitizeID(m.ID), m.Name, m.Current, trend))
		}
		sb.WriteString("}\n\n")
	}

	// Flow connections
	sb.WriteString("# Flow\n")
	sb.WriteString("vision -> hypotheses\n")
	if len(ls.MVPs) > 0 {
		sb.WriteString("hypotheses -> mvps\n")
		sb.WriteString("mvps -> bml\n")
	} else {
		sb.WriteString("hypotheses -> bml\n")
	}
	if len(ls.Experiments) > 0 {
		sb.WriteString("bml -> experiments\n")
	}
	if len(ls.Pivots) > 0 && len(ls.Experiments) > 0 {
		sb.WriteString("experiments -> pivots: pivot?\n")
	}
	if len(ls.Metrics) > 0 {
		sb.WriteString("bml -> metrics\n")
	}

	return []byte(sb.String()), nil
}

// renderDesignThinking renders a Design Thinking canvas with five phases.
func (r *D2Renderer) renderDesignThinking(dt *canvas.DesignThinkingCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Design Thinking Canvas: " + dt.Metadata.Title + "\n\n")
	sb.WriteString("direction: right\n\n")

	// Phase indicator
	sb.WriteString("phases: {\n")
	sb.WriteString("  label: \"Design Thinking Phases\"\n")
	sb.WriteString("  direction: right\n")
	sb.WriteString("  empathize: Empathize { style.fill: \"#e3f2fd\" }\n")
	sb.WriteString("  define: Define { style.fill: \"#f3e5f5\" }\n")
	sb.WriteString("  ideate: Ideate { style.fill: \"#e8f5e9\" }\n")
	sb.WriteString("  prototype: Prototype { style.fill: \"#fff3e0\" }\n")
	sb.WriteString("  test: Test { style.fill: \"#ffebee\" }\n")
	sb.WriteString("  empathize -> define -> ideate -> prototype -> test\n")

	// Highlight current phase
	switch dt.CurrentPhase {
	case canvas.DTPhaseEmpathize:
		sb.WriteString("  empathize.style.stroke-width: 3\n")
	case canvas.DTPhaseDefine:
		sb.WriteString("  define.style.stroke-width: 3\n")
	case canvas.DTPhaseIdeate:
		sb.WriteString("  ideate.style.stroke-width: 3\n")
	case canvas.DTPhasePrototype:
		sb.WriteString("  prototype.style.stroke-width: 3\n")
	case canvas.DTPhaseTest:
		sb.WriteString("  test.style.stroke-width: 3\n")
	}
	sb.WriteString("}\n\n")

	// Empathy Maps
	if len(dt.EmpathyMaps) > 0 {
		sb.WriteString("empathyMaps: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Empathy Maps (%d)\"\n", len(dt.EmpathyMaps)))
		sb.WriteString("  style.fill: \"#e3f2fd\"\n")
		for _, em := range dt.EmpathyMaps {
			sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", sanitizeID(em.ID), truncate(em.PersonaName, 30)))
		}
		sb.WriteString("}\n\n")
	}

	// Problem Statement (POV)
	if dt.ProblemStatement != "" {
		sb.WriteString("pov: {\n")
		sb.WriteString("  label: \"Point of View\"\n")
		sb.WriteString("  style.fill: \"#f3e5f5\"\n")
		sb.WriteString(fmt.Sprintf("  statement: \"%s\"\n", truncate(dt.ProblemStatement, 60)))
		sb.WriteString("}\n\n")
	}

	// How Might We questions
	if len(dt.HowMightWe) > 0 {
		sb.WriteString("hmw: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"How Might We (%d)\"\n", len(dt.HowMightWe)))
		sb.WriteString("  style.fill: \"#f3e5f5\"\n")
		for i, q := range dt.HowMightWe {
			sb.WriteString(fmt.Sprintf("  hmw%d: \"%s\"\n", i+1, truncate(q, 50)))
		}
		sb.WriteString("}\n\n")
	}

	// Ideas
	if len(dt.Ideas) > 0 {
		sb.WriteString("ideas: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Ideas (%d)\"\n", len(dt.Ideas)))
		sb.WriteString("  style.fill: \"#e8f5e9\"\n")
		for _, idea := range dt.Ideas {
			selected := ""
			if idea.Selected {
				selected = " ★"
			}
			sb.WriteString(fmt.Sprintf("  %s: \"%s%s\"\n", sanitizeID(idea.ID), truncate(idea.Title, 30), selected))
		}
		sb.WriteString("}\n\n")
	}

	// Prototypes
	if len(dt.Prototypes) > 0 {
		sb.WriteString("prototypes: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Prototypes (%d)\"\n", len(dt.Prototypes)))
		sb.WriteString("  style.fill: \"#fff3e0\"\n")
		for _, p := range dt.Prototypes {
			sb.WriteString(fmt.Sprintf("  %s: \"%s (%s, %s)\"\n", sanitizeID(p.ID), truncate(p.Name, 25), p.Type, p.Fidelity))
		}
		sb.WriteString("}\n\n")
	}

	// Tests
	if len(dt.Tests) > 0 {
		sb.WriteString("tests: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Tests (%d)\"\n", len(dt.Tests)))
		sb.WriteString("  style.fill: \"#ffebee\"\n")
		for _, t := range dt.Tests {
			result := "⏳"
			if len(t.Findings) > 0 {
				if t.ShouldContinue {
					result = "✓"
				} else {
					result = "↩"
				}
			}
			sb.WriteString(fmt.Sprintf("  %s: \"[%s] %s\"\n", sanitizeID(t.ID), result, t.Method))
		}
		sb.WriteString("}\n\n")
	}

	// Flow connections
	sb.WriteString("# Flow\n")
	if len(dt.EmpathyMaps) > 0 {
		sb.WriteString("phases.empathize -> empathyMaps\n")
		if dt.ProblemStatement != "" {
			sb.WriteString("empathyMaps -> pov\n")
		}
	}
	if dt.ProblemStatement != "" && len(dt.HowMightWe) > 0 {
		sb.WriteString("pov -> hmw\n")
	}
	if len(dt.HowMightWe) > 0 && len(dt.Ideas) > 0 {
		sb.WriteString("hmw -> ideas\n")
	} else if len(dt.Ideas) > 0 {
		sb.WriteString("phases.ideate -> ideas\n")
	}
	if len(dt.Ideas) > 0 && len(dt.Prototypes) > 0 {
		sb.WriteString("ideas -> prototypes\n")
	}
	if len(dt.Prototypes) > 0 && len(dt.Tests) > 0 {
		sb.WriteString("prototypes -> tests\n")
	}

	// Iteration loop
	if dt.IterationCount > 0 {
		sb.WriteString(fmt.Sprintf("tests -> prototypes: \"iterate (%d)\" {\n  style.stroke-dash: 3\n}\n", dt.IterationCount))
	}

	return []byte(sb.String()), nil
}

// renderJTBD renders a Jobs-to-be-Done canvas.
func (r *D2Renderer) renderJTBD(jtbd *canvas.JTBDCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("# Jobs-to-be-Done Canvas: " + jtbd.Metadata.Title + "\n\n")
	sb.WriteString("direction: down\n\n")

	// Main Job
	if jtbd.MainJob != nil {
		sb.WriteString("mainJob: {\n")
		sb.WriteString("  label: \"Main Job\"\n")
		sb.WriteString("  style.fill: \"#e3f2fd\"\n")
		sb.WriteString(fmt.Sprintf("  statement: \"%s\"\n", truncate(jtbd.MainJob.Statement, 60)))
		sb.WriteString(fmt.Sprintf("  type: \"Type: %s\"\n", jtbd.MainJob.Type))
		if jtbd.MainJob.Importance != "" {
			sb.WriteString(fmt.Sprintf("  importance: \"Importance: %s\"\n", jtbd.MainJob.Importance))
		}
		sb.WriteString("}\n\n")
	}

	// Related Jobs
	if len(jtbd.RelatedJobs) > 0 {
		sb.WriteString("relatedJobs: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Related Jobs (%d)\"\n", len(jtbd.RelatedJobs)))
		sb.WriteString("  style.fill: \"#e8f5e9\"\n")
		for _, job := range jtbd.RelatedJobs {
			icon := "📋"
			switch job.Type {
			case canvas.JobTypeFunctional:
				icon = "⚙️"
			case canvas.JobTypeEmotional:
				icon = "💭"
			case canvas.JobTypeSocial:
				icon = "👥"
			case canvas.JobTypeConsumption:
				icon = "📦"
			}
			sb.WriteString(fmt.Sprintf("  %s: \"%s %s\"\n", sanitizeID(job.ID), icon, truncate(job.Statement, 45)))
		}
		sb.WriteString("}\n\n")
	}

	// Job Map (stages)
	if len(jtbd.JobMap) > 0 {
		sb.WriteString("jobMap: {\n")
		sb.WriteString("  label: \"Job Map (Universal Job Map)\"\n")
		sb.WriteString("  style.fill: \"#fff3e0\"\n")
		sb.WriteString("  direction: right\n")
		for _, stage := range jtbd.JobMap {
			stageID := sanitizeID(stage.ID)
			sb.WriteString(fmt.Sprintf("  %s: \"%s\"\n", stageID, truncate(stage.Name, 20)))
		}
		// Connect stages in sequence
		for i := 0; i < len(jtbd.JobMap)-1; i++ {
			sb.WriteString(fmt.Sprintf("  %s -> %s\n", sanitizeID(jtbd.JobMap[i].ID), sanitizeID(jtbd.JobMap[i+1].ID)))
		}
		sb.WriteString("}\n\n")
	}

	// Desired Outcomes
	if len(jtbd.DesiredOutcomes) > 0 {
		sb.WriteString("desiredOutcomes: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Desired Outcomes (%d)\"\n", len(jtbd.DesiredOutcomes)))
		sb.WriteString("  style.fill: \"#c8e6c9\"\n")
		for _, out := range jtbd.DesiredOutcomes {
			opportunity := ""
			if out.Opportunity > 0 {
				opportunity = fmt.Sprintf(" [%.1f]", out.Opportunity)
			}
			sb.WriteString(fmt.Sprintf("  %s: \"%s %s%s\"\n", sanitizeID(out.ID), out.Direction, truncate(out.Statement, 40), opportunity))
		}
		sb.WriteString("}\n\n")
	}

	// Undesired Outcomes
	if len(jtbd.UndesiredOutcomes) > 0 {
		sb.WriteString("undesiredOutcomes: {\n")
		sb.WriteString(fmt.Sprintf("  label: \"Undesired Outcomes (%d)\"\n", len(jtbd.UndesiredOutcomes)))
		sb.WriteString("  style.fill: \"#ffcdd2\"\n")
		for _, out := range jtbd.UndesiredOutcomes {
			sb.WriteString(fmt.Sprintf("  %s: \"%s %s\"\n", sanitizeID(out.ID), out.Direction, truncate(out.Statement, 45)))
		}
		sb.WriteString("}\n\n")
	}

	// Forces Analysis
	hasForces := len(jtbd.PushForces) > 0 || len(jtbd.PullForces) > 0 || len(jtbd.Anxieties) > 0 || len(jtbd.Habits) > 0
	if hasForces {
		sb.WriteString("forces: {\n")
		sb.WriteString("  label: \"Forces Analysis\"\n")
		sb.WriteString("  style.fill: \"#f3e5f5\"\n")

		// Push forces
		if len(jtbd.PushForces) > 0 {
			sb.WriteString("  push: {\n")
			sb.WriteString("    label: \"Push (away from current)\"\n")
			sb.WriteString("    style.fill: \"#ffcdd2\"\n")
			for i, f := range jtbd.PushForces {
				strength := ""
				if f.Strength != "" {
					strength = fmt.Sprintf(" [%s]", f.Strength)
				}
				sb.WriteString(fmt.Sprintf("    p%d: \"← %s%s\"\n", i+1, truncate(f.Description, 35), strength))
			}
			sb.WriteString("  }\n")
		}

		// Pull forces
		if len(jtbd.PullForces) > 0 {
			sb.WriteString("  pull: {\n")
			sb.WriteString("    label: \"Pull (toward new)\"\n")
			sb.WriteString("    style.fill: \"#c8e6c9\"\n")
			for i, f := range jtbd.PullForces {
				strength := ""
				if f.Strength != "" {
					strength = fmt.Sprintf(" [%s]", f.Strength)
				}
				sb.WriteString(fmt.Sprintf("    p%d: \"→ %s%s\"\n", i+1, truncate(f.Description, 35), strength))
			}
			sb.WriteString("  }\n")
		}

		// Anxieties
		if len(jtbd.Anxieties) > 0 {
			sb.WriteString("  anxieties: {\n")
			sb.WriteString("    label: \"Anxieties\"\n")
			sb.WriteString("    style.fill: \"#fff9c4\"\n")
			for i, f := range jtbd.Anxieties {
				sb.WriteString(fmt.Sprintf("    a%d: \"⚠ %s\"\n", i+1, truncate(f.Description, 40)))
			}
			sb.WriteString("  }\n")
		}

		// Habits
		if len(jtbd.Habits) > 0 {
			sb.WriteString("  habits: {\n")
			sb.WriteString("    label: \"Habits\"\n")
			sb.WriteString("    style.fill: \"#e1f5fe\"\n")
			for i, f := range jtbd.Habits {
				sb.WriteString(fmt.Sprintf("    h%d: \"↻ %s\"\n", i+1, truncate(f.Description, 40)))
			}
			sb.WriteString("  }\n")
		}

		sb.WriteString("}\n\n")
	}

	// Hiring/Firing Solutions
	if len(jtbd.HiringSolutions) > 0 || len(jtbd.FiringSolutions) > 0 {
		sb.WriteString("solutions: {\n")
		sb.WriteString("  label: \"Current Solutions\"\n")
		sb.WriteString("  style.fill: \"#e0e0e0\"\n")
		for _, sol := range jtbd.HiringSolutions {
			why := ""
			if sol.WhyHired != "" {
				why = fmt.Sprintf(" - %s", truncate(sol.WhyHired, 25))
			}
			sb.WriteString(fmt.Sprintf("  %s: \"✓ %s%s\"\n", sanitizeID(sol.ID), truncate(sol.Name, 25), why))
		}
		for i, sol := range jtbd.FiringSolutions {
			sb.WriteString(fmt.Sprintf("  fired%d: \"✗ %s\"\n", i+1, truncate(sol, 40)))
		}
		sb.WriteString("}\n\n")
	}

	// Opportunity Scores
	if len(jtbd.OpportunityScores) > 0 {
		sb.WriteString("opportunityScores: {\n")
		sb.WriteString("  label: \"Opportunity Scores (ODI)\"\n")
		sb.WriteString("  style.fill: \"#dcedc8\"\n")
		for i, score := range jtbd.OpportunityScores {
			sb.WriteString(fmt.Sprintf("  score%d: \"Imp: %.1f | Sat: %.1f | Opp: %.1f\"\n", i+1, score.Importance, score.Satisfaction, score.Opportunity))
		}
		sb.WriteString("}\n\n")
	}

	// Flow connections
	sb.WriteString("# Flow\n")
	if jtbd.MainJob != nil && len(jtbd.RelatedJobs) > 0 {
		sb.WriteString("mainJob -> relatedJobs\n")
	}
	if jtbd.MainJob != nil && len(jtbd.JobMap) > 0 {
		sb.WriteString("mainJob -> jobMap\n")
	}
	if len(jtbd.JobMap) > 0 && len(jtbd.DesiredOutcomes) > 0 {
		sb.WriteString("jobMap -> desiredOutcomes\n")
	}
	if len(jtbd.DesiredOutcomes) > 0 && len(jtbd.OpportunityScores) > 0 {
		sb.WriteString("desiredOutcomes -> opportunityScores\n")
	}
	if hasForces && (len(jtbd.HiringSolutions) > 0 || len(jtbd.FiringSolutions) > 0) {
		sb.WriteString("forces -> solutions\n")
	}

	return []byte(sb.String()), nil
}

func init() {
	render.Register(NewD2Renderer())
}
