// Package mermaid provides Mermaid diagram renderers for canvas types.
package mermaid

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

// MermaidRenderer renders canvas types to Mermaid diagram syntax.
type MermaidRenderer struct{}

// NewMermaidRenderer creates a new Mermaid renderer.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{}
}

// Format returns the output format.
func (r *MermaidRenderer) Format() render.Format {
	return render.FormatMermaid
}

// FileExtension returns the file extension for Mermaid files.
func (r *MermaidRenderer) FileExtension() string {
	return ".mmd"
}

// Supports returns true for all canvas types.
func (r *MermaidRenderer) Supports(canvasType canvas.CanvasType) bool {
	return true
}

// Render converts a canvas to Mermaid format.
func (r *MermaidRenderer) Render(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
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

func (r *MermaidRenderer) renderOST(ost *canvas.OpportunitySolutionTree, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TD\n")
	sb.WriteString(fmt.Sprintf("    subgraph OST[\"%s\"]\n", ost.Metadata.Title))

	// Outcome at root
	outcomeID := sanitizeID(ost.Outcome.ID)
	sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", outcomeID, ost.Outcome.Description))

	// Opportunities
	for _, opp := range ost.Outcome.Opportunities {
		oppID := sanitizeID(opp.ID)
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", oppID, truncate(opp.Description, 40)))
		sb.WriteString(fmt.Sprintf("    %s --> %s\n", outcomeID, oppID))

		// Solutions
		for _, sol := range opp.Solutions {
			solID := sanitizeID(sol.ID)
			label := truncate(sol.Description, 35)
			if sol.Status != "" {
				label = fmt.Sprintf("%s\\n[%s]", label, sol.Status)
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", solID, label))
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", oppID, solID))

			// Experiments
			for _, exp := range sol.Experiments {
				expID := sanitizeID(exp.ID)
				expLabel := truncate(exp.Hypothesis, 30)
				if exp.Status != "" {
					expLabel = fmt.Sprintf("%s\\n[%s]", expLabel, exp.Status)
				}
				sb.WriteString(fmt.Sprintf("    %s{{\"%s\"}}\n", expID, expLabel))
				sb.WriteString(fmt.Sprintf("    %s --> %s\n", solID, expID))
			}
		}
	}

	sb.WriteString("    end\n")

	// Styling
	sb.WriteString("    style OST fill:#f9f9f9\n")
	sb.WriteString(fmt.Sprintf("    style %s fill:#E8F5E9\n", outcomeID))

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderBMC(bmc *canvas.BusinessModelCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph LR\n")
	sb.WriteString(fmt.Sprintf("    subgraph BMC[\"%s\"]\n", bmc.Metadata.Title))
	sb.WriteString("    direction TB\n")

	// Key Partners
	sb.WriteString("    subgraph KP[\"Key Partners\"]\n")
	for _, p := range bmc.KeyPartnerships {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(p.ID), p.Partner))
	}
	sb.WriteString("    end\n")

	// Key Activities
	sb.WriteString("    subgraph KA[\"Key Activities\"]\n")
	for _, a := range bmc.KeyActivities {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(a.ID), a.Name))
	}
	sb.WriteString("    end\n")

	// Value Propositions
	sb.WriteString("    subgraph VP[\"Value Propositions\"]\n")
	for _, vp := range bmc.ValuePropositions {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(vp.ID), truncate(vp.Description, 30)))
	}
	sb.WriteString("    end\n")

	// Customer Relationships
	sb.WriteString("    subgraph CR[\"Customer Relationships\"]\n")
	for _, cr := range bmc.CustomerRelationships {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(cr.ID), cr.Type))
	}
	sb.WriteString("    end\n")

	// Customer Segments
	sb.WriteString("    subgraph CS[\"Customer Segments\"]\n")
	for _, seg := range bmc.CustomerSegments {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(seg.ID), seg.Name))
	}
	sb.WriteString("    end\n")

	// Connections
	sb.WriteString("    KP --> KA\n")
	sb.WriteString("    KA --> VP\n")
	sb.WriteString("    VP --> CR\n")
	sb.WriteString("    CR --> CS\n")

	sb.WriteString("    end\n")

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderOpportunity(opp *canvas.OpportunityCanvas, opts *render.Options) ([]byte, error) {
	if opts.GridLayout {
		return r.renderOpportunityGrid(opp, opts)
	}
	return r.renderOpportunityFlow(opp, opts)
}

// renderOpportunityGrid renders Opportunity Canvas in BMC-style grid layout (no arrows).
func (r *MermaidRenderer) renderOpportunityGrid(opp *canvas.OpportunityCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph OC[\"%s\"]\n", opp.Metadata.Title))
	sb.WriteString("    direction LR\n")

	// Row 1: Users & Customers | Problems | Solution Ideas
	sb.WriteString("    subgraph ROW1[\" \"]\n")
	sb.WriteString("    direction LR\n")

	// Users & Customers
	sb.WriteString("    subgraph USERS[\"Users & Customers\"]\n")
	for _, u := range opp.Users {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(u.ID), u.Name))
	}
	sb.WriteString("    end\n")

	// Problems
	sb.WriteString("    subgraph PROBLEMS[\"Problems\"]\n")
	for _, p := range opp.Problems {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(p.ID), truncate(p.Description, 30)))
	}
	sb.WriteString("    end\n")

	// Solution Ideas
	sb.WriteString("    subgraph SOLUTIONS[\"Solution Ideas\"]\n")
	for i, s := range opp.SolutionIdeas {
		sb.WriteString(fmt.Sprintf("    si%d[\"%s\"]\n", i+1, truncate(s, 30)))
	}
	sb.WriteString("    end\n")

	sb.WriteString("    end\n") // ROW1

	// Row 2: Solutions Today | User Value | Adoption Strategy
	sb.WriteString("    subgraph ROW2[\" \"]\n")
	sb.WriteString("    direction LR\n")

	// Solutions Today
	sb.WriteString("    subgraph CURRENT[\"Solutions Today\"]\n")
	for _, s := range opp.CurrentSolutions {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(s.ID), truncate(s.Name, 30)))
	}
	sb.WriteString("    end\n")

	// User Value
	sb.WriteString("    subgraph UVAL[\"User Value\"]\n")
	for i, v := range opp.UserValue {
		sb.WriteString(fmt.Sprintf("    uv%d[\"%s\"]\n", i+1, truncate(v, 30)))
	}
	sb.WriteString("    end\n")

	// Adoption Strategy
	sb.WriteString("    subgraph ADOPT[\"Adoption Strategy\"]\n")
	for i, a := range opp.AdoptionStrategy {
		sb.WriteString(fmt.Sprintf("    as%d[\"%s\"]\n", i+1, truncate(a, 30)))
	}
	sb.WriteString("    end\n")

	sb.WriteString("    end\n") // ROW2

	// Row 3: User Metrics | Business Problem | Business Metrics
	sb.WriteString("    subgraph ROW3[\" \"]\n")
	sb.WriteString("    direction LR\n")

	// User Metrics
	sb.WriteString("    subgraph UMET[\"User Metrics\"]\n")
	for i, m := range opp.UserMetrics {
		sb.WriteString(fmt.Sprintf("    um%d[\"%s\"]\n", i+1, truncate(m, 30)))
	}
	sb.WriteString("    end\n")

	// Business Problem
	sb.WriteString("    subgraph BPROB[\"Business Problem\"]\n")
	if opp.BusinessProblem != "" {
		sb.WriteString(fmt.Sprintf("    bp[\"%s\"]\n", truncate(opp.BusinessProblem, 35)))
	}
	sb.WriteString("    end\n")

	// Business Metrics
	sb.WriteString("    subgraph BMET[\"Business Metrics\"]\n")
	for i, m := range opp.BusinessMetrics {
		sb.WriteString(fmt.Sprintf("    bm%d[\"%s\"]\n", i+1, truncate(m, 30)))
	}
	sb.WriteString("    end\n")

	sb.WriteString("    end\n") // ROW3

	// Row 4: Budget (full width)
	if opp.Budget != nil {
		sb.WriteString("    subgraph BUDGET[\"Budget\"]\n")
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
		if len(budgetItems) > 0 {
			sb.WriteString(fmt.Sprintf("    budget[\"%s\"]\n", strings.Join(budgetItems, " | ")))
		}
		sb.WriteString("    end\n")
	}

	sb.WriteString("    end\n") // OC

	// Styling
	sb.WriteString("    style ROW1 fill:none,stroke:none\n")
	sb.WriteString("    style ROW2 fill:none,stroke:none\n")
	sb.WriteString("    style ROW3 fill:none,stroke:none\n")

	return []byte(sb.String()), nil
}

// renderOpportunityFlow renders Opportunity Canvas with arrows showing flow.
func (r *MermaidRenderer) renderOpportunityFlow(opp *canvas.OpportunityCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph OC[\"%s\"]\n", opp.Metadata.Title))

	// Problem Space
	sb.WriteString("    subgraph PS[\"Problem Space\"]\n")
	for _, p := range opp.Problems {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(p.ID), truncate(p.Description, 35)))
	}
	for _, u := range opp.Users {
		sb.WriteString(fmt.Sprintf("    %s((\"%s\"))\n", sanitizeID(u.ID), u.Name))
	}
	sb.WriteString("    end\n")

	// Value Proposition
	sb.WriteString(fmt.Sprintf("    VP[\"%s\"]\n", truncate(opp.ValueProposition.Statement, 40)))

	// User Value
	sb.WriteString("    subgraph UV[\"User Value\"]\n")
	for i, v := range opp.UserValue {
		sb.WriteString(fmt.Sprintf("    uv%d[\"%s\"]\n", i+1, truncate(v, 30)))
	}
	sb.WriteString("    end\n")

	// Business Value
	sb.WriteString("    subgraph BV[\"Business Value\"]\n")
	for i, v := range opp.BusinessValue {
		sb.WriteString(fmt.Sprintf("    bv%d[\"%s\"]\n", i+1, truncate(v, 30)))
	}
	sb.WriteString("    end\n")

	// Validation
	sb.WriteString("    subgraph VAL[\"Validation\"]\n")
	for _, a := range opp.Assumptions {
		shape := "[\"%s\"]"
		if a.Validated {
			shape = "([\"%s\"])"
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(a.ID), truncate(a.Description, 30)))
	}
	sb.WriteString("    end\n")

	// Connections
	sb.WriteString("    PS --> VP\n")
	sb.WriteString("    VP --> UV\n")
	sb.WriteString("    VP --> BV\n")
	sb.WriteString("    UV --> VAL\n")
	sb.WriteString("    BV --> VAL\n")

	sb.WriteString("    end\n")

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderFeature(fc *canvas.FeatureCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph FC[\"%s\"]\n", fc.Metadata.Title))

	// Idea Statement
	sb.WriteString(fmt.Sprintf("    IDEA[\"%s\"]\n", fc.IdeaStatement))

	// Situations
	sb.WriteString("    subgraph SIT[\"Situations\"]\n")
	for _, s := range fc.Situations {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(s.ID), truncate(s.Description, 35)))
	}
	sb.WriteString("    end\n")

	// Problems
	sb.WriteString("    subgraph PROB[\"Problems\"]\n")
	for _, p := range fc.Problems {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(p.ID), truncate(p.Description, 35)))
	}
	sb.WriteString("    end\n")

	// Value
	sb.WriteString("    subgraph VAL[\"Value\"]\n")
	for _, v := range fc.Value {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(v.ID), truncate(v.Description, 35)))
	}
	sb.WriteString("    end\n")

	// Capabilities
	sb.WriteString("    subgraph CAP[\"Capabilities\"]\n")
	for _, c := range fc.Capabilities {
		label := truncate(c.Description, 30)
		if c.Priority != "" {
			label = fmt.Sprintf("%s [%s]", label, c.Priority)
		}
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(c.ID), label))
	}
	sb.WriteString("    end\n")

	// Constraints
	sb.WriteString("    subgraph CON[\"Constraints\"]\n")
	for i, rest := range fc.Restrictions {
		sb.WriteString(fmt.Sprintf("    r%d[\"%s\"]\n", i+1, truncate(rest, 30)))
	}
	for i, lim := range fc.Limitations {
		sb.WriteString(fmt.Sprintf("    l%d[\"%s\"]\n", i+1, truncate(lim, 30)))
	}
	sb.WriteString("    end\n")

	// Connections
	sb.WriteString("    IDEA --> SIT\n")
	sb.WriteString("    SIT --> PROB\n")
	sb.WriteString("    PROB --> VAL\n")
	sb.WriteString("    VAL --> CAP\n")
	sb.WriteString("    CAP --> CON\n")

	sb.WriteString("    end\n")

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderLeanUX(lux *canvas.LeanUXCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph LUX[\"%s\"]\n", lux.Metadata.Title))

	// Business Problem
	sb.WriteString(fmt.Sprintf("    BP[\"%s\"]\n", truncate(lux.BusinessProblem, 40)))

	// Business Outcomes
	sb.WriteString("    subgraph BO[\"Business Outcomes\"]\n")
	for _, o := range lux.BusinessOutcomes {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(o.ID), truncate(o.Description, 35)))
	}
	sb.WriteString("    end\n")

	// Users
	sb.WriteString("    subgraph USR[\"Users\"]\n")
	for _, u := range lux.Users {
		sb.WriteString(fmt.Sprintf("    %s((\"%s\"))\n", sanitizeID(u.ID), u.Name))
	}
	sb.WriteString("    end\n")

	// User Outcomes
	sb.WriteString("    subgraph UO[\"User Outcomes\"]\n")
	for _, o := range lux.UserOutcomes {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(o.ID), truncate(o.Description, 35)))
	}
	sb.WriteString("    end\n")

	// Solutions
	sb.WriteString("    subgraph SOL[\"Solutions\"]\n")
	for _, s := range lux.Solutions {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(s.ID), truncate(s.Description, 35)))
	}
	sb.WriteString("    end\n")

	// Hypotheses
	sb.WriteString("    subgraph HYP[\"Hypotheses\"]\n")
	for _, h := range lux.Hypotheses {
		shape := "[\"%s\"]"
		if h.Validated != nil && *h.Validated {
			shape = "([\"%s\"])"
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(h.ID), truncate(h.WeBelieve, 30)))
	}
	sb.WriteString("    end\n")

	// Riskiest Assumption
	sb.WriteString(fmt.Sprintf("    RA{{\"%s\"}}\n", truncate(lux.RiskiestAssumption, 35)))

	// Experiment
	if lux.Experiment != nil {
		sb.WriteString(fmt.Sprintf("    EXP[\"%s\\n[%s]\"]\n", truncate(lux.Experiment.Description, 30), lux.Experiment.Status))
	}

	// Connections
	sb.WriteString("    BP --> BO\n")
	sb.WriteString("    USR --> UO\n")
	sb.WriteString("    UO --> SOL\n")
	sb.WriteString("    SOL --> HYP\n")
	sb.WriteString("    HYP --> RA\n")
	if lux.Experiment != nil {
		sb.WriteString("    RA --> EXP\n")
	}

	sb.WriteString("    end\n")

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderShapeUpPitch(pitch *canvas.ShapeUpPitch, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph PITCH[\"%s\"]\n", pitch.Metadata.Title))

	// Problem section
	sb.WriteString("    subgraph PROBLEM[\"Problem\"]\n")
	sb.WriteString(fmt.Sprintf("    statement[\"%s\"]\n", truncate(pitch.Problem.Statement, 40)))
	if pitch.Problem.WhyNow != "" {
		sb.WriteString(fmt.Sprintf("    whyNow[\"Why now: %s\"]\n", truncate(pitch.Problem.WhyNow, 35)))
	}
	if pitch.Problem.Audience != "" {
		sb.WriteString(fmt.Sprintf("    audience[\"Audience: %s\"]\n", truncate(pitch.Problem.Audience, 35)))
	}
	sb.WriteString("    end\n")

	// Appetite section
	sb.WriteString("    subgraph APPETITE[\"Appetite\"]\n")
	appetiteLabel := fmt.Sprintf("%d weeks (%s)", pitch.Appetite.Weeks, pitch.Appetite.Size)
	sb.WriteString(fmt.Sprintf("    appetite[\"%s\"]\n", appetiteLabel))
	if pitch.Appetite.Rationale != "" {
		sb.WriteString(fmt.Sprintf("    rationale[\"%s\"]\n", truncate(pitch.Appetite.Rationale, 35)))
	}
	sb.WriteString("    end\n")

	// Solution section
	sb.WriteString("    subgraph SOLUTION[\"Solution\"]\n")
	sb.WriteString(fmt.Sprintf("    approach[\"%s\"]\n", truncate(pitch.Solution.Approach, 40)))
	for i, must := range pitch.Solution.MustInclude {
		sb.WriteString(fmt.Sprintf("    must%d[\"✓ %s\"]\n", i+1, truncate(must, 30)))
	}
	for i, nice := range pitch.Solution.NiceToHave {
		sb.WriteString(fmt.Sprintf("    nice%d[\"○ %s\"]\n", i+1, truncate(nice, 30)))
	}
	sb.WriteString("    end\n")

	// Rabbit Holes
	if len(pitch.RabbitHoles) > 0 {
		sb.WriteString("    subgraph RABBITS[\"Rabbit Holes\"]\n")
		for _, rh := range pitch.RabbitHoles {
			sb.WriteString(fmt.Sprintf("    %s{{\"⚠ %s\"}}\n", sanitizeID(rh.ID), truncate(rh.Description, 30)))
		}
		sb.WriteString("    end\n")
	}

	// No-Gos
	if len(pitch.NoGos) > 0 {
		sb.WriteString("    subgraph NOGOS[\"No-Gos\"]\n")
		for i, nogo := range pitch.NoGos {
			sb.WriteString(fmt.Sprintf("    nogo%d[\"✗ %s\"]\n", i+1, truncate(nogo, 30)))
		}
		sb.WriteString("    end\n")
	}

	// Connections
	sb.WriteString("    PROBLEM --> APPETITE\n")
	sb.WriteString("    APPETITE --> SOLUTION\n")
	if len(pitch.RabbitHoles) > 0 {
		sb.WriteString("    SOLUTION -.-> RABBITS\n")
	}
	if len(pitch.NoGos) > 0 {
		sb.WriteString("    SOLUTION -.-> NOGOS\n")
	}

	sb.WriteString("    end\n")

	// Betting status styling
	if pitch.BettingStatus == "bet" {
		sb.WriteString("    style PITCH fill:#E8F5E9\n")
	} else if pitch.BettingStatus == "declined" {
		sb.WriteString("    style PITCH fill:#FFEBEE\n")
	}

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderShapeUpBet(bet *canvas.ShapeUpBet, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph LR\n")
	sb.WriteString(fmt.Sprintf("    subgraph BET[\"%s\"]\n", bet.Metadata.Title))

	// Pitch reference
	sb.WriteString(fmt.Sprintf("    pitch[\"%s\"]\n", truncate(bet.PitchTitle, 35)))

	// Cycle info
	sb.WriteString("    subgraph CYCLE[\"Cycle\"]\n")
	sb.WriteString(fmt.Sprintf("    cycleName[\"%s\"]\n", bet.Cycle.Name))
	sb.WriteString(fmt.Sprintf("    cycleDate[\"%s to %s\"]\n", bet.Cycle.StartDate, bet.Cycle.EndDate))
	sb.WriteString(fmt.Sprintf("    weeks[\"%d weeks\"]\n", bet.Cycle.Weeks))
	sb.WriteString("    end\n")

	// Team
	sb.WriteString("    subgraph TEAM[\"Team\"]\n")
	if bet.Team.Designer != "" {
		sb.WriteString(fmt.Sprintf("    designer((\"Designer: %s\"))\n", bet.Team.Designer))
	}
	for i, prog := range bet.Team.Programmers {
		sb.WriteString(fmt.Sprintf("    prog%d((\"Dev: %s\"))\n", i+1, prog))
	}
	if bet.Team.Lead != "" {
		sb.WriteString(fmt.Sprintf("    lead((\"Lead: %s\"))\n", bet.Team.Lead))
	}
	sb.WriteString("    end\n")

	// Decision
	decisionShape := "\"%s\""
	if bet.Decision == "bet" {
		decisionShape = "([\"%s\"])"
	} else if bet.Decision == "declined" {
		decisionShape = "{{\"%s\"}}"
	}
	sb.WriteString(fmt.Sprintf("    decision"+decisionShape+"\n", "Decision: "+bet.Decision))

	// Connections
	sb.WriteString("    pitch --> CYCLE\n")
	sb.WriteString("    CYCLE --> TEAM\n")
	sb.WriteString("    TEAM --> decision\n")

	sb.WriteString("    end\n")

	// Styling based on decision
	if bet.Decision == "bet" {
		sb.WriteString("    style decision fill:#C8E6C9\n")
	} else if bet.Decision == "declined" {
		sb.WriteString("    style decision fill:#FFCDD2\n")
	} else if bet.Decision == "deferred" {
		sb.WriteString("    style decision fill:#FFF9C4\n")
	}

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderShapeUpScope(scope *canvas.ShapeUpScope, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph LR\n")
	sb.WriteString(fmt.Sprintf("    subgraph HILL[\"%s - Hill Chart\"]\n", scope.Metadata.Title))

	// Uphill section (figuring things out)
	sb.WriteString("    subgraph UPHILL[\"⛰ Uphill (Figuring Out)\"]\n")
	for _, s := range scope.Scopes {
		if s.HillPosition < 50 {
			sb.WriteString(fmt.Sprintf("    %s[\"%s\\n%d%%\"]\n", sanitizeID(s.ID), truncate(s.Name, 25), s.HillPosition))
		}
	}
	sb.WriteString("    end\n")

	// Top of hill
	sb.WriteString("    subgraph TOP[\"🏔 Top (Figured Out)\"]\n")
	for _, s := range scope.Scopes {
		if s.HillPosition >= 50 && s.HillPosition < 60 {
			sb.WriteString(fmt.Sprintf("    %s[\"%s\\n%d%%\"]\n", sanitizeID(s.ID), truncate(s.Name, 25), s.HillPosition))
		}
	}
	sb.WriteString("    end\n")

	// Downhill section (executing)
	sb.WriteString("    subgraph DOWNHILL[\"🎿 Downhill (Executing)\"]\n")
	for _, s := range scope.Scopes {
		if s.HillPosition >= 60 && s.HillPosition < 100 {
			sb.WriteString(fmt.Sprintf("    %s[\"%s\\n%d%%\"]\n", sanitizeID(s.ID), truncate(s.Name, 25), s.HillPosition))
		}
	}
	sb.WriteString("    end\n")

	// Done section
	sb.WriteString("    subgraph DONE[\"✅ Done\"]\n")
	for _, s := range scope.Scopes {
		if s.HillPosition == 100 || s.Status == "done" {
			sb.WriteString(fmt.Sprintf("    %s([\"%s\"])\n", sanitizeID(s.ID), truncate(s.Name, 25)))
		}
	}
	sb.WriteString("    end\n")

	// Flow
	sb.WriteString("    UPHILL --> TOP\n")
	sb.WriteString("    TOP --> DOWNHILL\n")
	sb.WriteString("    DOWNHILL --> DONE\n")

	sb.WriteString("    end\n")

	// Overall progress
	progress := scope.OverallProgress()
	sb.WriteString(fmt.Sprintf("    PROGRESS[\"Overall Progress: %d%%\"]\n", progress))
	sb.WriteString("    HILL --> PROGRESS\n")

	// Styling
	sb.WriteString("    style UPHILL fill:#FFF3E0\n")
	sb.WriteString("    style TOP fill:#E3F2FD\n")
	sb.WriteString("    style DOWNHILL fill:#E8F5E9\n")
	sb.WriteString("    style DONE fill:#C8E6C9\n")

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderDiscoverySnapshot(ds *canvas.DiscoverySnapshot, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph SNAPSHOT[\"%s - Week %s\"]\n", ds.Metadata.Title, ds.Week))

	// Interviews
	sb.WriteString("    subgraph INTERVIEWS[\"📞 Interviews\"]\n")
	for _, interview := range ds.Interviews {
		label := truncate(interview.ParticipantType, 25)
		if interview.Quality != "" {
			label += fmt.Sprintf(" [%s]", interview.Quality)
		}
		sb.WriteString(fmt.Sprintf("    %s((\"%s\"))\n", sanitizeID(interview.ID), label))
	}
	if len(ds.Interviews) == 0 {
		sb.WriteString("    noInt[\"No interviews\"]\n")
	}
	sb.WriteString("    end\n")

	// Opportunities discovered
	sb.WriteString("    subgraph OPPS[\"💡 Opportunities\"]\n")
	for _, opp := range ds.OpportunitiesDiscovered {
		icon := "💡"
		if opp.Action == "strengthened" {
			icon = "⬆"
		} else if opp.Action == "weakened" {
			icon = "⬇"
		} else if opp.Action == "retired" {
			icon = "❌"
		}
		sb.WriteString(fmt.Sprintf("    %s[\"%s %s\"]\n", sanitizeID(opp.OpportunityID), icon, truncate(opp.Description, 25)))
	}
	if len(ds.OpportunitiesDiscovered) == 0 {
		sb.WriteString("    noOpp[\"No updates\"]\n")
	}
	sb.WriteString("    end\n")

	// Assumption tests
	sb.WriteString("    subgraph TESTS[\"🧪 Assumption Tests\"]\n")
	for _, test := range ds.AssumptionTests {
		statusIcon := "⏳"
		if test.Status == "completed" {
			if test.Result == "validated" {
				statusIcon = "✅"
			} else if test.Result == "invalidated" {
				statusIcon = "❌"
			} else {
				statusIcon = "❓"
			}
		}
		sb.WriteString(fmt.Sprintf("    %s{{\"%s %s\"}}\n", sanitizeID(test.ID), statusIcon, truncate(test.Assumption.Description, 25)))
	}
	if len(ds.AssumptionTests) == 0 {
		sb.WriteString("    noTest[\"No tests\"]\n")
	}
	sb.WriteString("    end\n")

	// Key learnings
	if len(ds.KeyLearnings) > 0 {
		sb.WriteString("    subgraph LEARNINGS[\"📚 Key Learnings\"]\n")
		for i, learning := range ds.KeyLearnings {
			sb.WriteString(fmt.Sprintf("    learn%d[\"%s\"]\n", i+1, truncate(learning, 35)))
		}
		sb.WriteString("    end\n")
	}

	// Decisions
	if len(ds.Decisions) > 0 {
		sb.WriteString("    subgraph DECISIONS[\"🎯 Decisions\"]\n")
		for _, dec := range ds.Decisions {
			decIcon := "→"
			if dec.Type == "pivot" {
				decIcon = "↪"
			} else if dec.Type == "kill" {
				decIcon = "✗"
			}
			sb.WriteString(fmt.Sprintf("    %s([\"%s %s\"])\n", sanitizeID(dec.ID), decIcon, truncate(dec.Description, 30)))
		}
		sb.WriteString("    end\n")
	}

	// Flow
	sb.WriteString("    INTERVIEWS --> OPPS\n")
	sb.WriteString("    OPPS --> TESTS\n")
	if len(ds.KeyLearnings) > 0 {
		sb.WriteString("    TESTS --> LEARNINGS\n")
		if len(ds.Decisions) > 0 {
			sb.WriteString("    LEARNINGS --> DECISIONS\n")
		}
	} else if len(ds.Decisions) > 0 {
		sb.WriteString("    TESTS --> DECISIONS\n")
	}

	sb.WriteString("    end\n")

	// Weekly touchpoint indicator
	if ds.HasWeeklyTouchpoint() {
		sb.WriteString("    style INTERVIEWS fill:#E8F5E9\n")
	} else {
		sb.WriteString("    style INTERVIEWS fill:#FFEBEE\n")
	}

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderAssumptionMap(am *canvas.AssumptionMap, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TB\n")
	sb.WriteString(fmt.Sprintf("    subgraph AM[\"%s\"]\n", am.Metadata.Title))

	// Desirability
	sb.WriteString("    subgraph D[\"❤️ Desirability\"]\n")
	for _, a := range am.Desirability {
		shape := "\"%s\""
		if a.Validated {
			shape = "([\"%s\"])"
		}
		label := truncate(a.Description, 30)
		if a.Importance == "high" && a.Confidence == "low" {
			label = "⚠ " + label
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(a.ID), label))
	}
	if len(am.Desirability) == 0 {
		sb.WriteString("    dNone[\"None\"]\n")
	}
	sb.WriteString("    end\n")

	// Viability
	sb.WriteString("    subgraph V[\"💰 Viability\"]\n")
	for _, a := range am.Viability {
		shape := "\"%s\""
		if a.Validated {
			shape = "([\"%s\"])"
		}
		label := truncate(a.Description, 30)
		if a.Importance == "high" && a.Confidence == "low" {
			label = "⚠ " + label
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(a.ID), label))
	}
	if len(am.Viability) == 0 {
		sb.WriteString("    vNone[\"None\"]\n")
	}
	sb.WriteString("    end\n")

	// Feasibility
	sb.WriteString("    subgraph F[\"🔧 Feasibility\"]\n")
	for _, a := range am.Feasibility {
		shape := "\"%s\""
		if a.Validated {
			shape = "([\"%s\"])"
		}
		label := truncate(a.Description, 30)
		if a.Importance == "high" && a.Confidence == "low" {
			label = "⚠ " + label
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(a.ID), label))
	}
	if len(am.Feasibility) == 0 {
		sb.WriteString("    fNone[\"None\"]\n")
	}
	sb.WriteString("    end\n")

	// Usability
	sb.WriteString("    subgraph U[\"👤 Usability\"]\n")
	for _, a := range am.Usability {
		shape := "\"%s\""
		if a.Validated {
			shape = "([\"%s\"])"
		}
		label := truncate(a.Description, 30)
		if a.Importance == "high" && a.Confidence == "low" {
			label = "⚠ " + label
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(a.ID), label))
	}
	if len(am.Usability) == 0 {
		sb.WriteString("    uNone[\"None\"]\n")
	}
	sb.WriteString("    end\n")

	// Ethical
	sb.WriteString("    subgraph E[\"⚖️ Ethical\"]\n")
	for _, a := range am.Ethical {
		shape := "\"%s\""
		if a.Validated {
			shape = "([\"%s\"])"
		}
		label := truncate(a.Description, 30)
		if a.Importance == "high" && a.Confidence == "low" {
			label = "⚠ " + label
		}
		sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(a.ID), label))
	}
	if len(am.Ethical) == 0 {
		sb.WriteString("    eNone[\"None\"]\n")
	}
	sb.WriteString("    end\n")

	sb.WriteString("    end\n")

	// High risk assumptions callout
	highRisk := am.HighRiskAssumptions()
	if len(highRisk) > 0 {
		sb.WriteString("    subgraph HIGHRISK[\"⚠️ High Risk - Test First\"]\n")
		for _, a := range highRisk {
			sb.WriteString(fmt.Sprintf("    hr_%s{{\"⚠ %s\"}}\n", sanitizeID(a.ID), truncate(a.Description, 30)))
		}
		sb.WriteString("    end\n")
		sb.WriteString("    style HIGHRISK fill:#FFF3E0\n")
	}

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderExperienceMap(em *canvas.ExperienceMap, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph LR\n")
	sb.WriteString(fmt.Sprintf("    subgraph EM[\"%s\"]\n", em.Metadata.Title))

	// Persona
	if em.PersonaDescription != "" {
		sb.WriteString(fmt.Sprintf("    persona((\"👤 %s\"))\n", truncate(em.PersonaDescription, 30)))
	}

	// Journey phases
	sb.WriteString("    subgraph JOURNEY[\"Customer Journey\"]\n")
	sb.WriteString("    direction LR\n")

	var prevPhaseID string
	for _, phase := range em.Phases {
		phaseID := sanitizeID(phase.ID)

		// Phase container
		sb.WriteString(fmt.Sprintf("    subgraph %s[\"%s\"]\n", phaseID, phase.Name))

		// Actions
		for _, action := range phase.Actions {
			sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(action.ID), truncate(action.Description, 25)))
		}

		// Feeling indicator
		if phase.Feeling != "" {
			feelIcon := "😐"
			if phase.Feeling == "positive" {
				feelIcon = "😊"
			} else if phase.Feeling == "negative" {
				feelIcon = "😟"
			}
			sb.WriteString(fmt.Sprintf("    %s_feel[\"%s %s\"]\n", phaseID, feelIcon, truncate(phase.Thinking, 20)))
		}

		// Pain points
		for i, pain := range phase.PainPoints {
			sb.WriteString(fmt.Sprintf("    %s_pain%d{{\"⚠ %s\"}}\n", phaseID, i+1, truncate(pain, 20)))
		}

		sb.WriteString("    end\n")

		// Connect phases
		if prevPhaseID != "" {
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", prevPhaseID, phaseID))
		}
		prevPhaseID = phaseID
	}

	sb.WriteString("    end\n")

	// Connect persona to journey
	if em.PersonaDescription != "" && len(em.Phases) > 0 {
		sb.WriteString(fmt.Sprintf("    persona --> %s\n", sanitizeID(em.Phases[0].ID)))
	}

	sb.WriteString("    end\n")

	// Top opportunities
	if len(em.TopOpportunities) > 0 {
		sb.WriteString("    subgraph TOPOPPS[\"💡 Top Opportunities\"]\n")
		for i, opp := range em.TopOpportunities {
			sb.WriteString(fmt.Sprintf("    opp%d[\"%s\"]\n", i+1, truncate(opp, 35)))
		}
		sb.WriteString("    end\n")
		sb.WriteString("    JOURNEY --> TOPOPPS\n")
	}

	// Styling
	for _, phase := range em.Phases {
		if phase.Feeling == "positive" {
			sb.WriteString(fmt.Sprintf("    style %s fill:#E8F5E9\n", sanitizeID(phase.ID)))
		} else if phase.Feeling == "negative" {
			sb.WriteString(fmt.Sprintf("    style %s fill:#FFEBEE\n", sanitizeID(phase.ID)))
		}
	}

	return []byte(sb.String()), nil
}

// sanitizeID makes a string safe for use as a Mermaid identifier.
func sanitizeID(id string) string {
	// Replace spaces and special characters
	id = strings.ReplaceAll(id, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	id = strings.ReplaceAll(id, ".", "_")
	id = strings.ReplaceAll(id, "(", "_")
	id = strings.ReplaceAll(id, ")", "_")
	return id
}

// truncate shortens a string to max length, adding ellipsis if needed.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func (r *MermaidRenderer) renderLeanStartup(ls *canvas.LeanStartupCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TD\n")
	sb.WriteString(fmt.Sprintf("    subgraph LS[\"%s\"]\n", ls.Metadata.Title))

	// Vision and Strategy
	sb.WriteString("    subgraph VISION[\"🎯 Vision & Strategy\"]\n")
	if ls.Vision != "" {
		sb.WriteString(fmt.Sprintf("    vision[\"%s\"]\n", truncate(ls.Vision, 50)))
	}
	if ls.TargetCustomer != "" {
		sb.WriteString(fmt.Sprintf("    customer[\"👤 %s\"]\n", truncate(ls.TargetCustomer, 40)))
	}
	if ls.ProblemHypothesis != "" {
		sb.WriteString(fmt.Sprintf("    problem[\"❓ %s\"]\n", truncate(ls.ProblemHypothesis, 40)))
	}
	sb.WriteString("    end\n")

	// Hypotheses
	sb.WriteString("    subgraph HYPO[\"📋 Core Hypotheses\"]\n")
	if ls.ValueHypothesis != nil {
		status := "⏳"
		if ls.ValueHypothesis.Validated != nil {
			if *ls.ValueHypothesis.Validated {
				status = "✅"
			} else {
				status = "❌"
			}
		}
		sb.WriteString(fmt.Sprintf("    valueHypo[\"%s Value: %s\"]\n", status, truncate(ls.ValueHypothesis.Statement, 35)))
	}
	if ls.GrowthHypothesis != nil {
		status := "⏳"
		if ls.GrowthHypothesis.Validated != nil {
			if *ls.GrowthHypothesis.Validated {
				status = "✅"
			} else {
				status = "❌"
			}
		}
		model := string(ls.GrowthHypothesis.GrowthModel)
		sb.WriteString(fmt.Sprintf("    growthHypo[\"%s Growth (%s): %s\"]\n", status, model, truncate(ls.GrowthHypothesis.Statement, 30)))
	}
	sb.WriteString("    end\n")

	// Build-Measure-Learn cycle
	sb.WriteString("    subgraph BML[\"🔄 Build-Measure-Learn\"]\n")
	sb.WriteString("    build([\"🔨 Build\"])\n")
	sb.WriteString("    measure([\"📊 Measure\"])\n")
	sb.WriteString("    learn([\"💡 Learn\"])\n")
	sb.WriteString("    build --> measure --> learn --> build\n")
	sb.WriteString("    end\n")

	// MVPs
	if len(ls.MVPs) > 0 {
		sb.WriteString("    subgraph MVPS[\"🚀 MVP Iterations\"]\n")
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
				statusIcon = "✅"
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s (%s)\"]\n", sanitizeID(mvp.ID), statusIcon, truncate(mvp.Name, 25), mvp.Type))
		}
		sb.WriteString("    end\n")
	}

	// Experiments
	if len(ls.Experiments) > 0 {
		sb.WriteString("    subgraph EXPS[\"🧪 Experiments\"]\n")
		for _, exp := range ls.Experiments {
			statusIcon := "⏳"
			switch exp.Status {
			case canvas.ExperimentStatusRunning:
				statusIcon = "▶️"
			case canvas.ExperimentStatusCompleted:
				if exp.LearnValidated != nil && *exp.LearnValidated {
					statusIcon = "✅"
				} else if exp.LearnValidated != nil {
					statusIcon = "❌"
				} else {
					statusIcon = "✓"
				}
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s\"]\n", sanitizeID(exp.ID), statusIcon, truncate(exp.LearnHypothesis, 30)))
		}
		sb.WriteString("    end\n")
	}

	// Pivots
	if len(ls.Pivots) > 0 {
		sb.WriteString("    subgraph PIVOTS[\"↩️ Pivots\"]\n")
		for _, pivot := range ls.Pivots {
			sb.WriteString(fmt.Sprintf("    %s{{\"🔄 %s: %s → %s\"}}\n", sanitizeID(pivot.ID), pivot.Type, truncate(pivot.FromState, 15), truncate(pivot.ToState, 15)))
		}
		sb.WriteString("    end\n")
	}

	// Metrics
	if len(ls.Metrics) > 0 {
		sb.WriteString("    subgraph METRICS[\"📈 Innovation Accounting\"]\n")
		for _, m := range ls.Metrics {
			trend := ""
			switch m.Trend {
			case "improving":
				trend = "↑"
			case "declining":
				trend = "↓"
			default:
				trend = "→"
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s: %s\"]\n", sanitizeID(m.ID), trend, m.Name, m.Current))
		}
		sb.WriteString("    end\n")
	}

	// Connections
	sb.WriteString("    VISION --> HYPO\n")
	sb.WriteString("    HYPO --> BML\n")
	if len(ls.MVPs) > 0 {
		sb.WriteString("    BML --> MVPS\n")
		if len(ls.Experiments) > 0 {
			sb.WriteString("    MVPS --> EXPS\n")
		}
	} else if len(ls.Experiments) > 0 {
		sb.WriteString("    BML --> EXPS\n")
	}
	if len(ls.Pivots) > 0 {
		if len(ls.Experiments) > 0 {
			sb.WriteString("    EXPS -.-> PIVOTS\n")
		} else if len(ls.MVPs) > 0 {
			sb.WriteString("    MVPS -.-> PIVOTS\n")
		}
	}
	if len(ls.Metrics) > 0 {
		sb.WriteString("    BML --> METRICS\n")
	}

	sb.WriteString("    end\n")

	// PMF indicator
	if ls.HasProductMarketFit() {
		sb.WriteString("    style LS fill:#E8F5E9\n")
	}

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderDesignThinking(dt *canvas.DesignThinkingCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TD\n")
	sb.WriteString(fmt.Sprintf("    subgraph DT[\"%s\"]\n", dt.Metadata.Title))

	// Five phases
	sb.WriteString("    subgraph PHASES[\"Design Thinking Phases\"]\n")
	sb.WriteString("    direction LR\n")
	sb.WriteString("    empathize([\"❤️ Empathize\"])\n")
	sb.WriteString("    define([\"🎯 Define\"])\n")
	sb.WriteString("    ideate([\"💡 Ideate\"])\n")
	sb.WriteString("    prototype([\"🔨 Prototype\"])\n")
	sb.WriteString("    test([\"🧪 Test\"])\n")
	sb.WriteString("    empathize --> define --> ideate --> prototype --> test\n")
	sb.WriteString("    test -.->|iterate| ideate\n")
	sb.WriteString("    end\n")

	// Empathy Maps
	if len(dt.EmpathyMaps) > 0 {
		sb.WriteString("    subgraph EMPATHY[\"❤️ Empathy Maps\"]\n")
		for _, em := range dt.EmpathyMaps {
			sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(em.ID), truncate(em.PersonaName, 25)))
		}
		sb.WriteString("    end\n")
	}

	// User Needs & Insights
	if len(dt.UserNeeds) > 0 || len(dt.Insights) > 0 {
		sb.WriteString("    subgraph DEFINE[\"🎯 Insights & Needs\"]\n")
		for _, need := range dt.UserNeeds {
			validated := ""
			if need.Validated {
				validated = " ✓"
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s%s\"]\n", sanitizeID(need.ID), truncate(need.Need, 30), validated))
		}
		for _, insight := range dt.Insights {
			sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(insight.ID), truncate(insight.Description, 30)))
		}
		sb.WriteString("    end\n")
	}

	// POV Statement
	if dt.ProblemStatement != "" {
		sb.WriteString(fmt.Sprintf("    pov{{\"POV: %s\"}}\n", truncate(dt.ProblemStatement, 40)))
	}

	// HMW Questions
	if len(dt.HowMightWe) > 0 {
		sb.WriteString("    subgraph HMW[\"❓ How Might We\"]\n")
		for i, q := range dt.HowMightWe {
			sb.WriteString(fmt.Sprintf("    hmw%d[\"%s\"]\n", i+1, truncate(q, 40)))
		}
		sb.WriteString("    end\n")
	}

	// Ideas
	if len(dt.Ideas) > 0 {
		sb.WriteString("    subgraph IDEAS[\"💡 Ideas\"]\n")
		for _, idea := range dt.Ideas {
			shape := "\"%s\""
			if idea.Selected {
				shape = "([\"%s ★\"])"
			}
			sb.WriteString(fmt.Sprintf("    %s"+shape+"\n", sanitizeID(idea.ID), truncate(idea.Title, 25)))
		}
		sb.WriteString("    end\n")
	}

	// Prototypes
	if len(dt.Prototypes) > 0 {
		sb.WriteString("    subgraph PROTOS[\"🔨 Prototypes\"]\n")
		for _, p := range dt.Prototypes {
			fidelity := ""
			if p.Fidelity != "" {
				fidelity = fmt.Sprintf(" [%s]", p.Fidelity)
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s%s\"]\n", sanitizeID(p.ID), truncate(p.Name, 25), fidelity))
		}
		sb.WriteString("    end\n")
	}

	// Tests
	if len(dt.Tests) > 0 {
		sb.WriteString("    subgraph TESTS[\"🧪 Tests\"]\n")
		for _, t := range dt.Tests {
			result := "⏳"
			if len(t.Findings) > 0 {
				if t.ShouldContinue {
					result = "✅"
				} else {
					result = "↩"
				}
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s\"]\n", sanitizeID(t.ID), result, t.Method))
		}
		sb.WriteString("    end\n")
	}

	// Flow connections
	if len(dt.EmpathyMaps) > 0 {
		sb.WriteString("    EMPATHY --> DEFINE\n")
	}
	if dt.ProblemStatement != "" {
		if len(dt.UserNeeds) > 0 || len(dt.Insights) > 0 {
			sb.WriteString("    DEFINE --> pov\n")
		}
		if len(dt.HowMightWe) > 0 {
			sb.WriteString("    pov --> HMW\n")
		}
	}
	if len(dt.HowMightWe) > 0 && len(dt.Ideas) > 0 {
		sb.WriteString("    HMW --> IDEAS\n")
	}
	if len(dt.Ideas) > 0 && len(dt.Prototypes) > 0 {
		sb.WriteString("    IDEAS --> PROTOS\n")
	}
	if len(dt.Prototypes) > 0 && len(dt.Tests) > 0 {
		sb.WriteString("    PROTOS --> TESTS\n")
	}
	if dt.IterationCount > 0 && len(dt.Tests) > 0 {
		sb.WriteString(fmt.Sprintf("    TESTS -.->|%d iterations| IDEAS\n", dt.IterationCount))
	}

	sb.WriteString("    end\n")

	// Style current phase
	switch dt.CurrentPhase {
	case canvas.DTPhaseEmpathize:
		sb.WriteString("    style empathize fill:#BBDEFB\n")
	case canvas.DTPhaseDefine:
		sb.WriteString("    style define fill:#E1BEE7\n")
	case canvas.DTPhaseIdeate:
		sb.WriteString("    style ideate fill:#C8E6C9\n")
	case canvas.DTPhasePrototype:
		sb.WriteString("    style prototype fill:#FFE0B2\n")
	case canvas.DTPhaseTest:
		sb.WriteString("    style test fill:#FFCDD2\n")
	}

	return []byte(sb.String()), nil
}

func (r *MermaidRenderer) renderJTBD(jtbd *canvas.JTBDCanvas, _ *render.Options) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("graph TD\n")
	sb.WriteString(fmt.Sprintf("    subgraph JTBD[\"%s\"]\n", jtbd.Metadata.Title))

	// Main Job
	if jtbd.MainJob != nil {
		sb.WriteString("    subgraph MAINJOB[\"🎯 Main Job\"]\n")
		icon := "📋"
		switch jtbd.MainJob.Type {
		case canvas.JobTypeFunctional:
			icon = "⚙️"
		case canvas.JobTypeEmotional:
			icon = "💭"
		case canvas.JobTypeSocial:
			icon = "👥"
		case canvas.JobTypeConsumption:
			icon = "📦"
		}
		importance := ""
		if jtbd.MainJob.Importance != "" {
			importance = fmt.Sprintf(" [%s]", jtbd.MainJob.Importance)
		}
		sb.WriteString(fmt.Sprintf("    mainJob[\"%s %s%s\"]\n", icon, truncate(jtbd.MainJob.Statement, 45), importance))
		sb.WriteString("    end\n")
	}

	// Related Jobs
	if len(jtbd.RelatedJobs) > 0 {
		sb.WriteString("    subgraph RELATED[\"📋 Related Jobs\"]\n")
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
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s\"]\n", sanitizeID(job.ID), icon, truncate(job.Statement, 35)))
		}
		sb.WriteString("    end\n")
	}

	// Job Map (Universal Job Map stages)
	if len(jtbd.JobMap) > 0 {
		sb.WriteString("    subgraph JOBMAP[\"🗺️ Job Map\"]\n")
		sb.WriteString("    direction LR\n")
		var prevID string
		for _, stage := range jtbd.JobMap {
			stageID := sanitizeID(stage.ID)
			sb.WriteString(fmt.Sprintf("    %s([%q])\n", stageID, truncate(stage.Name, 20)))
			if prevID != "" {
				sb.WriteString(fmt.Sprintf("    %s --> %s\n", prevID, stageID))
			}
			prevID = stageID
		}
		sb.WriteString("    end\n")
	}

	// Desired Outcomes
	if len(jtbd.DesiredOutcomes) > 0 {
		sb.WriteString("    subgraph OUTCOMES[\"✅ Desired Outcomes\"]\n")
		for _, out := range jtbd.DesiredOutcomes {
			opportunity := ""
			if out.Opportunity > 0 {
				opportunity = fmt.Sprintf(" [%.1f]", out.Opportunity)
			}
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s%s\"]\n", sanitizeID(out.ID), out.Direction, truncate(out.Statement, 30), opportunity))
		}
		sb.WriteString("    end\n")
	}

	// Undesired Outcomes
	if len(jtbd.UndesiredOutcomes) > 0 {
		sb.WriteString("    subgraph UNDESIRED[\"❌ Undesired Outcomes\"]\n")
		for _, out := range jtbd.UndesiredOutcomes {
			sb.WriteString(fmt.Sprintf("    %s[\"%s %s\"]\n", sanitizeID(out.ID), out.Direction, truncate(out.Statement, 35)))
		}
		sb.WriteString("    end\n")
	}

	// Forces Analysis
	hasForces := len(jtbd.PushForces) > 0 || len(jtbd.PullForces) > 0 || len(jtbd.Anxieties) > 0 || len(jtbd.Habits) > 0
	if hasForces {
		sb.WriteString("    subgraph FORCES[\"⚡ Forces Analysis\"]\n")

		// Push Forces
		if len(jtbd.PushForces) > 0 {
			sb.WriteString("    subgraph PUSH[\"← Push\"]\n")
			for i, f := range jtbd.PushForces {
				sb.WriteString(fmt.Sprintf("    push%d[\"%s\"]\n", i+1, truncate(f.Description, 30)))
			}
			sb.WriteString("    end\n")
		}

		// Pull Forces
		if len(jtbd.PullForces) > 0 {
			sb.WriteString("    subgraph PULL[\"→ Pull\"]\n")
			for i, f := range jtbd.PullForces {
				sb.WriteString(fmt.Sprintf("    pull%d[\"%s\"]\n", i+1, truncate(f.Description, 30)))
			}
			sb.WriteString("    end\n")
		}

		// Anxieties
		if len(jtbd.Anxieties) > 0 {
			sb.WriteString("    subgraph ANXIETY[\"😰 Anxieties\"]\n")
			for i, f := range jtbd.Anxieties {
				sb.WriteString(fmt.Sprintf("    anx%d{{\"⚠ %s\"}}\n", i+1, truncate(f.Description, 25)))
			}
			sb.WriteString("    end\n")
		}

		// Habits
		if len(jtbd.Habits) > 0 {
			sb.WriteString("    subgraph HABITS[\"↻ Habits\"]\n")
			for i, f := range jtbd.Habits {
				sb.WriteString(fmt.Sprintf("    hab%d[\"%s\"]\n", i+1, truncate(f.Description, 30)))
			}
			sb.WriteString("    end\n")
		}

		sb.WriteString("    end\n")
	}

	// Hiring/Firing Solutions
	if len(jtbd.HiringSolutions) > 0 || len(jtbd.FiringSolutions) > 0 {
		sb.WriteString("    subgraph SOLUTIONS[\"📦 Current Solutions\"]\n")
		for _, sol := range jtbd.HiringSolutions {
			sb.WriteString(fmt.Sprintf("    %s([\"✓ %s\"])\n", sanitizeID(sol.ID), truncate(sol.Name, 30)))
		}
		for i, sol := range jtbd.FiringSolutions {
			sb.WriteString(fmt.Sprintf("    fired%d{{\"✗ %s\"}}\n", i+1, truncate(sol, 30)))
		}
		sb.WriteString("    end\n")
	}

	// Opportunity Scores
	if len(jtbd.OpportunityScores) > 0 {
		sb.WriteString("    subgraph ODI[\"📊 Opportunity Scores\"]\n")
		for i, score := range jtbd.OpportunityScores {
			sb.WriteString(fmt.Sprintf("    score%d[\"I:%.1f S:%.1f → %.1f\"]\n", i+1, score.Importance, score.Satisfaction, score.Opportunity))
		}
		sb.WriteString("    end\n")
	}

	// Flow connections
	if jtbd.MainJob != nil && len(jtbd.RelatedJobs) > 0 {
		sb.WriteString("    MAINJOB --> RELATED\n")
	}
	if jtbd.MainJob != nil && len(jtbd.JobMap) > 0 {
		sb.WriteString("    MAINJOB --> JOBMAP\n")
	}
	if len(jtbd.JobMap) > 0 && len(jtbd.DesiredOutcomes) > 0 {
		sb.WriteString("    JOBMAP --> OUTCOMES\n")
	}
	if len(jtbd.DesiredOutcomes) > 0 && len(jtbd.OpportunityScores) > 0 {
		sb.WriteString("    OUTCOMES --> ODI\n")
	}
	if hasForces && (len(jtbd.HiringSolutions) > 0 || len(jtbd.FiringSolutions) > 0) {
		sb.WriteString("    FORCES --> SOLUTIONS\n")
	}

	sb.WriteString("    end\n")

	return []byte(sb.String()), nil
}

func init() {
	render.Register(NewMermaidRenderer())
}
