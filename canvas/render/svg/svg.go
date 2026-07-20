// Package svg provides native SVG renderers for canvas types.
package svg

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

// SVGRenderer renders canvas types to native SVG.
type SVGRenderer struct{}

// NewSVGRenderer creates a new native SVG renderer.
func NewSVGRenderer() *SVGRenderer {
	return &SVGRenderer{}
}

// Format returns the output format.
func (r *SVGRenderer) Format() render.Format {
	return render.FormatSVG
}

// FileExtension returns the file extension for SVG files.
func (r *SVGRenderer) FileExtension() string {
	return ".svg"
}

// Supports returns true for supported canvas types.
func (r *SVGRenderer) Supports(canvasType canvas.CanvasType) bool {
	switch canvasType {
	case canvas.CanvasTypeOpportunity,
		canvas.CanvasTypeLeanUX,
		canvas.CanvasTypeBMC,
		canvas.CanvasTypeOpportunitySpec:
		return true
	}
	return false
}

// Render converts a canvas to SVG format.
func (r *SVGRenderer) Render(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	switch c.Type {
	case canvas.CanvasTypeOpportunity:
		if opts.GridLayout {
			return r.renderOpportunityGrid(c.Opportunity, opts)
		}
		return nil, fmt.Errorf("SVG renderer only supports grid layout for Opportunity Canvas")
	case canvas.CanvasTypeLeanUX:
		if opts.GridLayout {
			return r.renderLeanUXGrid(c.LeanUX, opts)
		}
		return nil, fmt.Errorf("SVG renderer only supports grid layout for Lean UX Canvas")
	case canvas.CanvasTypeBMC:
		if c.BMC == nil {
			return nil, fmt.Errorf("canvas type %s has no BMC data", c.Type)
		}
		return r.renderBMCGrid(c.BMC, opts)
	case canvas.CanvasTypeOpportunitySpec:
		if c.OpportunitySpec == nil {
			return nil, fmt.Errorf("canvas type %s has no OpportunitySpec data", c.Type)
		}
		return r.renderOpportunitySpecGrid(c.OpportunitySpec, opts)
	default:
		return nil, fmt.Errorf("SVG renderer does not support canvas type: %s", c.Type)
	}
}

// Color scheme matching the reference SVG
type colorScheme struct {
	fill         string
	stroke       string
	titleFill    string
	subtitleFill string
}

var (
	blueScheme = colorScheme{
		fill:         "rgb(12, 68, 124)",
		stroke:       "rgb(133, 183, 235)",
		titleFill:    "rgb(181, 212, 244)",
		subtitleFill: "rgb(133, 183, 235)",
	}
	greenScheme = colorScheme{
		fill:         "rgb(8, 80, 65)",
		stroke:       "rgb(93, 202, 165)",
		titleFill:    "rgb(159, 225, 203)",
		subtitleFill: "rgb(93, 202, 165)",
	}
	orangeScheme = colorScheme{
		fill:         "rgb(99, 56, 6)",
		stroke:       "rgb(239, 159, 39)",
		titleFill:    "rgb(250, 199, 117)",
		subtitleFill: "rgb(239, 159, 39)",
	}
	grayScheme = colorScheme{
		fill:         "rgb(68, 68, 65)",
		stroke:       "rgb(180, 178, 169)",
		titleFill:    "rgb(211, 209, 199)",
		subtitleFill: "rgb(180, 178, 169)",
	}
	// Lean UX Canvas color schemes
	purpleScheme = colorScheme{
		fill:         "rgb(60, 52, 137)",
		stroke:       "rgb(175, 169, 236)",
		titleFill:    "rgb(206, 203, 246)",
		subtitleFill: "rgb(175, 169, 236)",
	}
	coralScheme = colorScheme{
		fill:         "rgb(113, 43, 19)",
		stroke:       "rgb(240, 153, 123)",
		titleFill:    "rgb(245, 196, 179)",
		subtitleFill: "rgb(240, 153, 123)",
	}
)

func (r *SVGRenderer) renderOpportunityGrid(opp *canvas.OpportunityCanvas, _ *render.Options) ([]byte, error) {
	const (
		width        = 680
		cellWidth    = 190
		cellHeight   = 100 // Increased to fit content
		budgetHeight = 70
		gap          = 10
		startX       = 40
		startY       = 20
	)

	// Calculate total height
	totalHeight := startY + 3*(cellHeight+gap) + budgetHeight + 20

	var sb strings.Builder

	// SVG header
	sb.WriteString(fmt.Sprintf(`<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, width, totalHeight))
	sb.WriteString(fmt.Sprintf(`<title>%s</title>`, escapeXML(opp.Metadata.Title)))
	sb.WriteString(`<desc>Opportunity Canvas - Jeff Patton's 9-block structure</desc>`)

	// Row 1: Users & customers | Problems | Solution ideas
	row1Y := startY
	sb.WriteString(r.renderCell(startX, row1Y, cellWidth, cellHeight,
		"Users & customers", "Who has the problem",
		formatUsers(opp.Users), blueScheme))
	sb.WriteString(r.renderCell(startX+cellWidth+gap, row1Y, cellWidth, cellHeight,
		"Problems", "Pains to address",
		formatProblems(opp.Problems), blueScheme))
	sb.WriteString(r.renderCell(startX+2*(cellWidth+gap), row1Y, cellWidth, cellHeight,
		"Solution ideas", "Ways to solve it",
		formatStrings(opp.SolutionIdeas), blueScheme))

	// Row 2: Solutions today | User value | Adoption strategy
	row2Y := row1Y + cellHeight + gap
	sb.WriteString(r.renderCell(startX, row2Y, cellWidth, cellHeight,
		"Solutions today", "Current workarounds",
		formatSolutions(opp.CurrentSolutions), greenScheme))
	sb.WriteString(r.renderCell(startX+cellWidth+gap, row2Y, cellWidth, cellHeight,
		"User value", "Benefit to users",
		formatStrings(opp.UserValue), greenScheme))
	sb.WriteString(r.renderCell(startX+2*(cellWidth+gap), row2Y, cellWidth, cellHeight,
		"Adoption strategy", "How they'll find it",
		formatStrings(opp.AdoptionStrategy), greenScheme))

	// Row 3: User metrics | Business problem | Business metrics
	row3Y := row2Y + cellHeight + gap
	sb.WriteString(r.renderCell(startX, row3Y, cellWidth, cellHeight,
		"User metrics", "Behaviour to track",
		formatStrings(opp.UserMetrics), orangeScheme))
	sb.WriteString(r.renderCell(startX+cellWidth+gap, row3Y, cellWidth, cellHeight,
		"Business problem", "Why it matters to us",
		[]string{opp.BusinessProblem}, orangeScheme))
	sb.WriteString(r.renderCell(startX+2*(cellWidth+gap), row3Y, cellWidth, cellHeight,
		"Business metrics", "Outcome to measure",
		formatStrings(opp.BusinessMetrics), orangeScheme))

	// Row 4: Budget (full width)
	row4Y := row3Y + cellHeight + gap
	budgetWidth := 3*cellWidth + 2*gap
	budgetContent := formatBudget(opp.Budget)
	sb.WriteString(r.renderCell(startX, row4Y, budgetWidth, budgetHeight,
		"Budget", "What you're willing to invest to learn whether this is worth building",
		budgetContent, grayScheme))

	sb.WriteString(`</svg>`)

	return []byte(sb.String()), nil
}

func (r *SVGRenderer) renderLeanUXGrid(lux *canvas.LeanUXCanvas, _ *render.Options) ([]byte, error) {
	const (
		width      = 680
		startX     = 40
		startY     = 30
		gap        = 10
		row1Height = 80 // Business problem / Business outcomes
		row2Height = 90 // Users / Benefits / Solutions
		row3Height = 70 // Hypotheses
		row4Height = 90 // Riskiest assumption / Smallest experiment
	)

	// Calculate column widths
	halfWidth := (width - 2*startX - gap) / 2    // 295
	thirdWidth := (width - 2*startX - 2*gap) / 3 // 190 (approx)
	fullWidth := width - 2*startX                // 600

	// Calculate total height
	totalHeight := startY + row1Height + gap + row2Height + gap + row3Height + gap + row4Height + 30

	var sb strings.Builder

	// SVG header
	sb.WriteString(fmt.Sprintf(`<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, width, totalHeight))
	sb.WriteString(fmt.Sprintf(`<title>%s</title>`, escapeXML(lux.Metadata.Title)))
	sb.WriteString(`<desc>Lean UX Canvas - Jeff Gothelf's 8-block structure</desc>`)

	// Row 1: Business problem | Business outcomes (Purple)
	row1Y := startY
	sb.WriteString(r.renderCell(startX, row1Y, halfWidth, row1Height,
		"1 · Business problem", "The problem to solve",
		[]string{lux.BusinessProblem}, purpleScheme))
	sb.WriteString(r.renderCell(startX+halfWidth+gap, row1Y, halfWidth, row1Height,
		"5 · Business outcomes", "Signals of success",
		formatOutcomes(lux.BusinessOutcomes), purpleScheme))

	// Row 2: Users | Benefits | Solutions (Green)
	// Use explicit positions to ensure consistent gaps
	row2Y := row1Y + row1Height + gap
	col1X := startX
	col2X := startX + thirdWidth + gap
	col3X := col2X + thirdWidth + gap + 10 // Account for middle column's extra width
	sb.WriteString(r.renderCell(col1X, row2Y, thirdWidth, row2Height,
		"2 · Users", "Who we serve",
		formatLeanUXUsers(lux.Users), greenScheme))
	sb.WriteString(r.renderCell(col2X, row2Y, thirdWidth+10, row2Height,
		"3 · Benefits", "User outcomes",
		formatOutcomes(lux.UserOutcomes), greenScheme))
	sb.WriteString(r.renderCell(col3X, row2Y, thirdWidth, row2Height,
		"4 · Solutions", "Ideas to try",
		formatLeanUXSolutions(lux.Solutions), greenScheme))

	// Row 3: Hypotheses (Orange - full width)
	row3Y := row2Y + row2Height + gap
	sb.WriteString(r.renderCell(startX, row3Y, fullWidth, row3Height,
		"6 · Hypotheses", "We believe [solution] will achieve [outcome] for [users]",
		formatHypotheses(lux.Hypotheses), orangeScheme))

	// Row 4: Riskiest assumption | Smallest experiment (Coral)
	row4Y := row3Y + row3Height + gap
	sb.WriteString(r.renderCell(startX, row4Y, halfWidth, row4Height,
		"7 · Riskiest assumption", "Most important thing to learn first",
		[]string{lux.RiskiestAssumption}, coralScheme))
	sb.WriteString(r.renderCell(startX+halfWidth+gap, row4Y, halfWidth, row4Height,
		"8 · Smallest experiment", "Least work needed to learn it",
		formatExperiment(lux.Experiment), coralScheme))

	sb.WriteString(`</svg>`)

	return []byte(sb.String()), nil
}

func (r *SVGRenderer) renderCell(x, y, w, h int, title, subtitle string, content []string, scheme colorScheme) string {
	const (
		rx        = 8 // Corner radius
		fontStyle = `font-family:-apple-system, "system-ui", "Segoe UI", sans-serif`
	)

	var sb strings.Builder

	centerX := x + w/2

	// Group
	sb.WriteString(`<g>`)

	// Rectangle
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" rx="%d" `+
		`style="fill:%s;stroke:%s;stroke-width:0.5"/>`,
		x, y, w, h, rx, scheme.fill, scheme.stroke))

	// Title
	titleY := y + 18
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" `+
		`style="fill:%s;font-size:13px;font-weight:600;%s">%s</text>`,
		centerX, titleY, scheme.titleFill, fontStyle, escapeXML(title)))

	// Subtitle
	subtitleY := titleY + 14
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" `+
		`style="fill:%s;font-size:10px;font-weight:400;font-style:italic;%s">%s</text>`,
		centerX, subtitleY, scheme.subtitleFill, fontStyle, escapeXML(subtitle)))

	// Content items (if any)
	if len(content) > 0 && content[0] != "" {
		contentY := subtitleY + 16

		// For wide cells (budget), show items horizontally
		if w > 400 {
			// Horizontal layout for budget row
			itemWidth := w / 4
			for i, item := range content {
				if i >= 4 {
					break
				}
				if item == "" {
					continue
				}
				itemX := x + 12 + i*itemWidth
				sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="start" `+
					`style="fill:%s;font-size:10px;font-weight:400;%s">• %s</text>`,
					itemX, contentY, scheme.subtitleFill, fontStyle, escapeXML(item)))
			}
		} else {
			// Vertical layout for regular cells
			contentX := x + 12        // Left-align content with padding
			maxItems := (h - 50) / 13 // Calculate max items that fit
			if maxItems > 4 {
				maxItems = 4
			}
			for i, item := range content {
				if i >= maxItems {
					break
				}
				if item == "" {
					continue
				}
				// Truncate long items based on cell width
				maxLen := (w - 24) / 5 // Approximate chars that fit
				displayItem := item
				if len(displayItem) > maxLen {
					displayItem = displayItem[:maxLen-3] + "..."
				}
				sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="start" `+
					`style="fill:%s;font-size:10px;font-weight:400;%s">• %s</text>`,
					contentX, contentY+i*13, scheme.subtitleFill, fontStyle, escapeXML(displayItem)))
			}
		}
	}

	sb.WriteString(`</g>`)

	return sb.String()
}

func formatUsers(users []canvas.User) []string {
	result := make([]string, 0, len(users))
	for _, u := range users {
		result = append(result, u.Name)
	}
	return result
}

func formatProblems(problems []canvas.Problem) []string {
	result := make([]string, 0, len(problems))
	for _, p := range problems {
		result = append(result, p.Description)
	}
	return result
}

func formatSolutions(solutions []canvas.Solution) []string {
	result := make([]string, 0, len(solutions))
	for _, s := range solutions {
		result = append(result, s.Name)
	}
	return result
}

func formatStrings(items []string) []string {
	return items
}

func formatBudget(budget *canvas.Budget) []string {
	if budget == nil {
		return nil
	}
	var items []string
	if budget.TimeEstimate != "" {
		items = append(items, "Time: "+budget.TimeEstimate)
	}
	if budget.TeamSize != "" {
		items = append(items, "Team: "+budget.TeamSize)
	}
	if budget.CostEstimate != "" {
		items = append(items, "Cost: "+budget.CostEstimate)
	}
	return items
}

func formatOutcomes(outcomes []canvas.Outcome) []string {
	result := make([]string, 0, len(outcomes))
	for _, o := range outcomes {
		result = append(result, o.Description)
	}
	return result
}

func formatLeanUXUsers(users []canvas.LeanUXUser) []string {
	result := make([]string, 0, len(users))
	for _, u := range users {
		result = append(result, u.Name)
	}
	return result
}

func formatLeanUXSolutions(solutions []canvas.LeanUXSolution) []string {
	result := make([]string, 0, len(solutions))
	for _, s := range solutions {
		result = append(result, s.Description)
	}
	return result
}

func formatHypotheses(hypotheses []canvas.Hypothesis) []string {
	result := make([]string, 0, len(hypotheses))
	for _, h := range hypotheses {
		result = append(result, h.WeBelieve)
	}
	return result
}

func formatExperiment(exp *canvas.Experiment) []string {
	if exp == nil {
		return nil
	}
	var items []string
	if exp.Name != "" {
		items = append(items, exp.Name)
	} else if exp.Description != "" {
		items = append(items, exp.Description)
	}
	if exp.Method != "" {
		items = append(items, "Method: "+exp.Method)
	}
	if exp.Status != "" {
		items = append(items, "Status: "+string(exp.Status))
	}
	return items
}

// renderBMCGrid renders a Business Model Canvas in Osterwalder's nine-block
// layout: five columns across the top (two of them split into upper and lower
// cells), with Cost Structure and Revenue Streams spanning the bottom.
func (r *SVGRenderer) renderBMCGrid(bmc *canvas.BusinessModelCanvas, _ *render.Options) ([]byte, error) {
	const (
		startX     = 20
		startY     = 20
		gap        = 8
		colWidth   = 196
		topHeight  = 230
		botHeight  = 90
		halfHeight = (topHeight - gap) / 2
	)

	width := startX*2 + 5*colWidth + 4*gap
	totalHeight := startY + topHeight + gap + botHeight + 20
	colX := func(i int) int { return startX + i*(colWidth+gap) }

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, width, totalHeight))
	sb.WriteString(fmt.Sprintf(`<title>%s</title>`, escapeXML(bmc.Metadata.Title)))
	sb.WriteString(`<desc>Business Model Canvas - Osterwalder's 9-block structure</desc>`)

	// Column 1: Key Partnerships (full height).
	sb.WriteString(r.renderCell(colX(0), startY, colWidth, topHeight,
		"Key Partnerships", "Who helps us",
		mapStrings(bmc.KeyPartnerships, func(p canvas.Partnership) string { return p.Description }), blueScheme))

	// Column 2: Key Activities (top) / Key Resources (bottom).
	sb.WriteString(r.renderCell(colX(1), startY, colWidth, halfHeight,
		"Key Activities", "What we do",
		mapStrings(bmc.KeyActivities, func(a canvas.Activity) string { return a.Name }), blueScheme))
	sb.WriteString(r.renderCell(colX(1), startY+halfHeight+gap, colWidth, halfHeight,
		"Key Resources", "What we need",
		mapStrings(bmc.KeyResources, func(x canvas.Resource) string { return x.Name }), blueScheme))

	// Column 3: Value Propositions (full height, centre).
	sb.WriteString(r.renderCell(colX(2), startY, colWidth, topHeight,
		"Value Propositions", "The value we deliver",
		mapStrings(bmc.ValuePropositions, func(v canvas.ValueProposition) string { return v.Description }), greenScheme))

	// Column 4: Customer Relationships (top) / Channels (bottom).
	sb.WriteString(r.renderCell(colX(3), startY, colWidth, halfHeight,
		"Customer Relationships", "How we engage",
		mapStrings(bmc.CustomerRelationships, func(c canvas.CustomerRelation) string { return c.Description }), orangeScheme))
	sb.WriteString(r.renderCell(colX(3), startY+halfHeight+gap, colWidth, halfHeight,
		"Channels", "How we reach them",
		mapStrings(bmc.Channels, func(c canvas.Channel) string { return c.Name }), orangeScheme))

	// Column 5: Customer Segments (full height).
	sb.WriteString(r.renderCell(colX(4), startY, colWidth, topHeight,
		"Customer Segments", "Who we serve",
		mapStrings(bmc.CustomerSegments, func(c canvas.CustomerSegment) string { return c.Name }), orangeScheme))

	// Bottom row: Cost Structure | Revenue Streams.
	botY := startY + topHeight + gap
	botWidth := (5*colWidth + 4*gap - gap) / 2
	sb.WriteString(r.renderCell(startX, botY, botWidth, botHeight,
		"Cost Structure", "What it costs to operate",
		mapStrings(bmc.CostStructure, func(c canvas.Cost) string { return c.Description }), grayScheme))
	sb.WriteString(r.renderCell(startX+botWidth+gap, botY, botWidth, botHeight,
		"Revenue Streams", "How we earn",
		mapStrings(bmc.RevenueStreams, func(x canvas.RevenueStream) string { return x.Description }), grayScheme))

	sb.WriteString(`</svg>`)
	return []byte(sb.String()), nil
}

// renderOpportunitySpecGrid renders the OpportunitySpec in its twelve-box layout
// (four rows of three), coloured by row: Discovery, Value, Market, Validation.
func (r *SVGRenderer) renderOpportunitySpecGrid(os *canvas.OpportunitySpec, _ *render.Options) ([]byte, error) {
	const (
		startX     = 40
		startY     = 20
		gap        = 10
		cellWidth  = 200
		cellHeight = 115
	)

	width := startX*2 + 3*cellWidth + 2*gap
	totalHeight := startY + 4*(cellHeight+gap) + 10
	colX := func(i int) int { return startX + i*(cellWidth+gap) }
	rowY := func(i int) int { return startY + i*(cellHeight+gap) }

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg width="100%%" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, width, totalHeight))
	sb.WriteString(fmt.Sprintf(`<title>%s</title>`, escapeXML(os.Metadata.Title)))
	sb.WriteString(`<desc>OpportunitySpec - 12-box merge of Patton's Opportunity Canvas and Cagan's Opportunity Assessment</desc>`)

	cell := func(col, row int, title, subtitle string, content []string, scheme colorScheme) {
		sb.WriteString(r.renderCell(colX(col), rowY(row), cellWidth, cellHeight, title, subtitle, content, scheme))
	}

	// Row 1: Discovery (blue).
	cell(0, 0, "1 · Users & Problem", "Who has the problem",
		lead(os.UsersAndProblem.EstimatedUserBase, os.UsersAndProblem.AffectedPersonas), blueScheme)
	cell(1, 0, "2 · Current Solutions", "How they solve it today",
		concat(os.CurrentSolutions.Workarounds, os.CurrentSolutions.InternalSolutions), blueScheme)
	cell(2, 0, "3 · Solution Ideas", "Our concepts",
		lead(os.SolutionIdeas.RecommendedIdea, os.SolutionIdeas.AlternativesPros), blueScheme)

	// Row 2: Value (green).
	cell(0, 1, "4 · User Value", "Value to users",
		lead(os.UserValue.ValueStatement, os.UserValue.KeyBenefits), greenScheme)
	cell(1, 1, "5 · Business Value", "Value to the business",
		lead(os.BusinessValue.BusinessProblem, os.BusinessValue.BusinessOutcomes), greenScheme)
	cell(2, 1, "6 · Competitive Edge", "Why us",
		lead(os.CompetitiveEdge.Differentiator, os.CompetitiveEdge.CoreStrengths), greenScheme)

	// Row 3: Market (orange).
	cell(0, 2, "7 · Market & Timing", "Who and why now",
		lead(os.MarketAndTiming.PrimarySegment, os.MarketAndTiming.Industries), orangeScheme)
	cell(1, 2, "8 · Go-to-Market", "How we reach users",
		lead(os.GoToMarket.Strategy, os.GoToMarket.Channels), orangeScheme)
	cell(2, 2, "9 · Success Metrics", "How we measure",
		concat(os.SuccessMetrics.UserMetrics, os.SuccessMetrics.BusinessMetrics), orangeScheme)

	// Row 4: Validation (purple).
	cell(0, 3, "10 · Critical Requirements", "What must be true",
		os.CriticalRequirements.MustHaveCapabilities, purpleScheme)
	cell(1, 3, "11 · Risks & Assumptions", "What we're betting on",
		compact([]string{os.RisksAndAssumptions.RiskiestAssumption, os.RisksAndAssumptions.HighestRisk}), purpleScheme)
	cell(2, 3, "12 · Recommendation", "The decision",
		compact([]string{os.Recommendation.Decision, os.Recommendation.Confidence}), purpleScheme)

	sb.WriteString(`</svg>`)
	return []byte(sb.String()), nil
}

// mapStrings projects a slice through f, dropping empty results.
func mapStrings[T any](items []T, f func(T) string) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		if s := f(it); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// lead prepends a non-empty leading string to items.
func lead(head string, items []string) []string {
	if head == "" {
		return items
	}
	return append([]string{head}, items...)
}

// concat joins two string slices.
func concat(a, b []string) []string {
	return append(append([]string{}, a...), b...)
}

// compact drops empty strings.
func compact(items []string) []string {
	out := make([]string, 0, len(items))
	for _, s := range items {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func init() {
	render.Register(NewSVGRenderer())
}
