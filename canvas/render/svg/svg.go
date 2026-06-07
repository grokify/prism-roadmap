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
	return canvasType == canvas.CanvasTypeOpportunity || canvasType == canvas.CanvasTypeLeanUX
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
