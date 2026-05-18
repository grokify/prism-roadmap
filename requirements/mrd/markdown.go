package mrd

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-execution/common"
)

// MarkdownOptions configures markdown generation.
type MarkdownOptions struct {
	IncludeFrontmatter bool
	Margin             string
	MainFont           string
	SansFont           string
	MonoFont           string
	FontFamily         string
}

// DefaultMarkdownOptions returns default options.
// Default fonts work with standard pdflatex. For Unicode/emoji support,
// use --pdf-engine=xelatex with appropriate system fonts.
func DefaultMarkdownOptions() MarkdownOptions {
	return MarkdownOptions{
		IncludeFrontmatter: true,
		Margin:             "2cm",
		MainFont:           "",
		SansFont:           "",
		MonoFont:           "",
		FontFamily:         "",
	}
}

// ToMarkdown converts the MRD to markdown format.
func (d *Document) ToMarkdown(opts MarkdownOptions) string {
	var sb strings.Builder

	if opts.IncludeFrontmatter {
		sb.WriteString(d.generateFrontmatter(opts))
	}

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", d.Metadata.Title))

	// Document info table
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|-------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **ID** | %s |\n", d.Metadata.ID))
	sb.WriteString(fmt.Sprintf("| **Version** | %s |\n", d.Metadata.Version))
	sb.WriteString(fmt.Sprintf("| **Status** | %s |\n", d.Metadata.Status))
	sb.WriteString(fmt.Sprintf("| **Created** | %s |\n", d.Metadata.CreatedAt.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("| **Updated** | %s |\n", d.Metadata.UpdatedAt.Format("2006-01-02")))

	if len(d.Metadata.Authors) > 0 {
		sb.WriteString(fmt.Sprintf("| **Author(s)** | %s |\n", common.FormatPeopleMarkdown(d.Metadata.Authors)))
	}

	if len(d.Metadata.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("| **Tags** | %s |\n", strings.Join(d.Metadata.Tags, ", ")))
	}
	sb.WriteString("\n---\n\n")

	// Executive Summary
	sb.WriteString("## 1. Executive Summary\n\n")
	sb.WriteString("### 1.1 Market Opportunity\n\n")
	sb.WriteString(d.ExecutiveSummary.MarketOpportunity)
	sb.WriteString("\n\n")

	sb.WriteString("### 1.2 Proposed Offering\n\n")
	sb.WriteString(d.ExecutiveSummary.ProposedOffering)
	sb.WriteString("\n\n")

	if len(d.ExecutiveSummary.KeyFindings) > 0 {
		sb.WriteString("### 1.3 Key Findings\n\n")
		for _, finding := range d.ExecutiveSummary.KeyFindings {
			sb.WriteString(fmt.Sprintf("- %s\n", finding))
		}
		sb.WriteString("\n")
	}

	if d.ExecutiveSummary.Recommendation != "" {
		sb.WriteString("### 1.4 Recommendation\n\n")
		sb.WriteString(d.ExecutiveSummary.Recommendation)
		sb.WriteString("\n\n")
	}

	sb.WriteString("---\n\n")

	// Market Overview
	sb.WriteString("## 2. Market Overview\n\n")
	sb.WriteString("### 2.1 Market Size\n\n")
	sb.WriteString("| Metric | Value | Year | Source |\n")
	sb.WriteString("|--------|-------|------|--------|\n")
	sb.WriteString(fmt.Sprintf("| **TAM** (Total Addressable Market) | %s | %d | %s |\n",
		d.MarketOverview.TAM.Value, d.MarketOverview.TAM.Year, d.MarketOverview.TAM.Source))
	sb.WriteString(fmt.Sprintf("| **SAM** (Serviceable Addressable Market) | %s | %d | %s |\n",
		d.MarketOverview.SAM.Value, d.MarketOverview.SAM.Year, d.MarketOverview.SAM.Source))
	sb.WriteString(fmt.Sprintf("| **SOM** (Serviceable Obtainable Market) | %s | %d | %s |\n",
		d.MarketOverview.SOM.Value, d.MarketOverview.SOM.Year, d.MarketOverview.SOM.Source))
	sb.WriteString("\n")

	if d.MarketOverview.GrowthRate != "" {
		sb.WriteString(fmt.Sprintf("**Growth Rate:** %s\n\n", d.MarketOverview.GrowthRate))
	}

	if d.MarketOverview.MarketStage != "" {
		sb.WriteString(fmt.Sprintf("**Market Stage:** %s\n\n", d.MarketOverview.MarketStage))
	}

	if len(d.MarketOverview.Trends) > 0 {
		sb.WriteString("### 2.2 Market Trends\n\n")
		sb.WriteString("| Trend | Description | Impact | Timeframe |\n")
		sb.WriteString("|-------|-------------|--------|----------|\n")
		for _, t := range d.MarketOverview.Trends {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				t.Name, t.Description, t.Impact, t.Timeframe))
		}
		sb.WriteString("\n")
	}

	if len(d.MarketOverview.Drivers) > 0 {
		sb.WriteString("### 2.3 Market Drivers\n\n")
		for _, driver := range d.MarketOverview.Drivers {
			sb.WriteString(fmt.Sprintf("- %s\n", driver))
		}
		sb.WriteString("\n")
	}

	if len(d.MarketOverview.Barriers) > 0 {
		sb.WriteString("### 2.4 Barriers to Entry\n\n")
		for _, barrier := range d.MarketOverview.Barriers {
			sb.WriteString(fmt.Sprintf("- %s\n", barrier))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Target Market
	sb.WriteString("## 3. Target Market\n\n")

	if len(d.TargetMarket.PrimarySegments) > 0 {
		sb.WriteString("### 3.1 Primary Segments\n\n")
		for _, seg := range d.TargetMarket.PrimarySegments {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", seg.Name))
			sb.WriteString(fmt.Sprintf("%s\n\n", seg.Description))
			if seg.Size != "" {
				sb.WriteString(fmt.Sprintf("- **Size:** %s\n", seg.Size))
			}
			if seg.Growth != "" {
				sb.WriteString(fmt.Sprintf("- **Growth:** %s\n", seg.Growth))
			}
			if len(seg.Needs) > 0 {
				sb.WriteString("- **Key Needs:**\n")
				for _, need := range seg.Needs {
					sb.WriteString(fmt.Sprintf("  - %s\n", need))
				}
			}
			if len(seg.Challenges) > 0 {
				sb.WriteString("- **Challenges:**\n")
				for _, ch := range seg.Challenges {
					sb.WriteString(fmt.Sprintf("  - %s\n", ch))
				}
			}
			sb.WriteString("\n")
		}
	}

	if len(d.TargetMarket.SecondarySegments) > 0 {
		sb.WriteString("### 3.2 Secondary Segments\n\n")
		for _, seg := range d.TargetMarket.SecondarySegments {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", seg.Name))
			sb.WriteString(fmt.Sprintf("%s\n\n", seg.Description))
		}
	}

	if len(d.TargetMarket.BuyerPersonas) > 0 {
		sb.WriteString("### 3.3 Buyer Personas\n\n")
		for _, p := range d.TargetMarket.BuyerPersonas {
			sb.WriteString(fmt.Sprintf("#### %s (%s)\n\n", p.Name, p.Title))
			sb.WriteString(fmt.Sprintf("%s\n\n", p.Description))
			sb.WriteString(fmt.Sprintf("- **Buying Role:** %s\n", p.BuyingRole))
			sb.WriteString(fmt.Sprintf("- **Budget Authority:** %v\n", p.BudgetAuthority))
			if len(p.PainPoints) > 0 {
				sb.WriteString("- **Pain Points:**\n")
				for _, pp := range p.PainPoints {
					sb.WriteString(fmt.Sprintf("  - %s\n", pp))
				}
			}
			if len(p.Goals) > 0 {
				sb.WriteString("- **Goals:**\n")
				for _, g := range p.Goals {
					sb.WriteString(fmt.Sprintf("  - %s\n", g))
				}
			}
			if len(p.BuyingCriteria) > 0 {
				sb.WriteString("- **Buying Criteria:**\n")
				for _, bc := range p.BuyingCriteria {
					sb.WriteString(fmt.Sprintf("  - %s\n", bc))
				}
			}
			sb.WriteString("\n")
		}
	}

	if len(d.TargetMarket.Verticals) > 0 || len(d.TargetMarket.GeographicFocus) > 0 || len(d.TargetMarket.CompanySize) > 0 {
		sb.WriteString("### 3.4 Market Focus\n\n")
		if len(d.TargetMarket.Verticals) > 0 {
			sb.WriteString(fmt.Sprintf("- **Industry Verticals:** %s\n", strings.Join(d.TargetMarket.Verticals, ", ")))
		}
		if len(d.TargetMarket.GeographicFocus) > 0 {
			sb.WriteString(fmt.Sprintf("- **Geographic Focus:** %s\n", strings.Join(d.TargetMarket.GeographicFocus, ", ")))
		}
		if len(d.TargetMarket.CompanySize) > 0 {
			sb.WriteString(fmt.Sprintf("- **Company Size:** %s\n", strings.Join(d.TargetMarket.CompanySize, ", ")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Competitive Landscape
	sb.WriteString("## 4. Competitive Landscape\n\n")
	sb.WriteString("### 4.1 Overview\n\n")
	sb.WriteString(d.CompetitiveLandscape.Overview)
	sb.WriteString("\n\n")

	if len(d.CompetitiveLandscape.Competitors) > 0 {
		sb.WriteString("### 4.2 Competitor Analysis\n\n")
		for _, c := range d.CompetitiveLandscape.Competitors {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", c.Name))
			if c.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", c.Description))
			}
			sb.WriteString("| Attribute | Value |\n")
			sb.WriteString("|-----------|-------|\n")
			if c.Category != "" {
				sb.WriteString(fmt.Sprintf("| **Category** | %s |\n", c.Category))
			}
			if c.MarketShare != "" {
				sb.WriteString(fmt.Sprintf("| **Market Share** | %s |\n", c.MarketShare))
			}
			if c.Pricing != "" {
				sb.WriteString(fmt.Sprintf("| **Pricing** | %s |\n", c.Pricing))
			}
			if c.ThreatLevel != "" {
				sb.WriteString(fmt.Sprintf("| **Threat Level** | %s |\n", c.ThreatLevel))
			}
			sb.WriteString("\n")

			if len(c.Strengths) > 0 {
				sb.WriteString("**Strengths:**\n")
				for _, s := range c.Strengths {
					sb.WriteString(fmt.Sprintf("- %s\n", s))
				}
				sb.WriteString("\n")
			}
			if len(c.Weaknesses) > 0 {
				sb.WriteString("**Weaknesses:**\n")
				for _, w := range c.Weaknesses {
					sb.WriteString(fmt.Sprintf("- %s\n", w))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(d.CompetitiveLandscape.Differentiators) > 0 {
		sb.WriteString("### 4.3 Our Differentiators\n\n")
		for _, diff := range d.CompetitiveLandscape.Differentiators {
			sb.WriteString(fmt.Sprintf("- %s\n", diff))
		}
		sb.WriteString("\n")
	}

	if len(d.CompetitiveLandscape.CompetitiveGaps) > 0 {
		sb.WriteString("### 4.4 Competitive Gaps to Address\n\n")
		for _, gap := range d.CompetitiveLandscape.CompetitiveGaps {
			sb.WriteString(fmt.Sprintf("- %s\n", gap))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Market Requirements
	sb.WriteString("## 5. Market Requirements\n\n")
	if len(d.MarketRequirements) > 0 {
		sb.WriteString("| ID | Requirement | Priority | Category | Source |\n")
		sb.WriteString("|----|-------------|----------|----------|--------|\n")
		for _, req := range d.MarketRequirements {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				req.ID, req.Title, req.Priority, req.Category, req.Source))
		}
		sb.WriteString("\n")

		// Detailed requirements
		sb.WriteString("### Requirement Details\n\n")
		for _, req := range d.MarketRequirements {
			sb.WriteString(fmt.Sprintf("#### %s: %s\n\n", req.ID, req.Title))
			sb.WriteString(fmt.Sprintf("%s\n\n", req.Description))
			if req.Validation != "" {
				sb.WriteString(fmt.Sprintf("**Validation:** %s\n\n", req.Validation))
			}
			if len(req.Segments) > 0 {
				sb.WriteString(fmt.Sprintf("**Target Segments:** %s\n\n", strings.Join(req.Segments, ", ")))
			}
		}
	}

	sb.WriteString("---\n\n")

	// Positioning
	sb.WriteString("## 6. Positioning\n\n")
	sb.WriteString("### 6.1 Positioning Statement\n\n")
	sb.WriteString(fmt.Sprintf("> %s\n\n", d.Positioning.Statement))

	sb.WriteString("| Attribute | Value |\n")
	sb.WriteString("|-----------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **Target Audience** | %s |\n", d.Positioning.TargetAudience))
	sb.WriteString(fmt.Sprintf("| **Category** | %s |\n", d.Positioning.Category))
	if d.Positioning.Tagline != "" {
		sb.WriteString(fmt.Sprintf("| **Tagline** | %s |\n", d.Positioning.Tagline))
	}
	sb.WriteString("\n")

	if len(d.Positioning.KeyBenefits) > 0 {
		sb.WriteString("### 6.2 Key Benefits\n\n")
		for _, b := range d.Positioning.KeyBenefits {
			sb.WriteString(fmt.Sprintf("- %s\n", b))
		}
		sb.WriteString("\n")
	}

	if len(d.Positioning.Differentiators) > 0 {
		sb.WriteString("### 6.3 Differentiators\n\n")
		for _, d := range d.Positioning.Differentiators {
			sb.WriteString(fmt.Sprintf("- %s\n", d))
		}
		sb.WriteString("\n")
	}

	if len(d.Positioning.ProofPoints) > 0 {
		sb.WriteString("### 6.4 Proof Points\n\n")
		for _, p := range d.Positioning.ProofPoints {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Go-to-Market (optional)
	if d.GoToMarket != nil {
		sb.WriteString("## 7. Go-to-Market Strategy\n\n")

		if d.GoToMarket.LaunchStrategy != "" {
			sb.WriteString("### 7.1 Launch Strategy\n\n")
			sb.WriteString(d.GoToMarket.LaunchStrategy)
			sb.WriteString("\n\n")
		}

		if d.GoToMarket.LaunchTiming != "" {
			sb.WriteString(fmt.Sprintf("**Target Launch:** %s\n\n", d.GoToMarket.LaunchTiming))
		}

		if d.GoToMarket.PricingStrategy != nil {
			sb.WriteString("### 7.2 Pricing Strategy\n\n")
			sb.WriteString(fmt.Sprintf("**Model:** %s\n\n", d.GoToMarket.PricingStrategy.Model))
			if d.GoToMarket.PricingStrategy.Positioning != "" {
				sb.WriteString(fmt.Sprintf("**Positioning:** %s\n\n", d.GoToMarket.PricingStrategy.Positioning))
			}
			if len(d.GoToMarket.PricingStrategy.Tiers) > 0 {
				sb.WriteString("| Tier | Price | Billing | Target Buyer |\n")
				sb.WriteString("|------|-------|---------|---------------|\n")
				for _, tier := range d.GoToMarket.PricingStrategy.Tiers {
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
						tier.Name, tier.Price, tier.Billing, tier.TargetBuyer))
				}
				sb.WriteString("\n")
			}
		}

		if len(d.GoToMarket.DistributionChannels) > 0 {
			sb.WriteString("### 7.3 Distribution Channels\n\n")
			for _, ch := range d.GoToMarket.DistributionChannels {
				sb.WriteString(fmt.Sprintf("- %s\n", ch))
			}
			sb.WriteString("\n")
		}

		if len(d.GoToMarket.Milestones) > 0 {
			sb.WriteString("### 7.4 Milestones\n\n")
			sb.WriteString("| Milestone | Target Date | Status |\n")
			sb.WriteString("|-----------|-------------|--------|\n")
			for _, m := range d.GoToMarket.Milestones {
				date := ""
				if !m.TargetDate.IsZero() {
					date = m.TargetDate.Format("2006-01-02")
				}
				sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", m.Name, date, m.Status))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	// Success Metrics
	sectionNum := 7
	if d.GoToMarket != nil {
		sectionNum = 8
	}
	sb.WriteString(fmt.Sprintf("## %d. Success Metrics\n\n", sectionNum))
	if len(d.SuccessMetrics) > 0 {
		sb.WriteString("| ID | Metric | Target | Timeframe | Measurement |\n")
		sb.WriteString("|----|--------|--------|-----------|-------------|\n")
		for _, m := range d.SuccessMetrics {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				m.ID, m.Name, m.Target, m.Timeframe, m.MeasurementMethod))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Risks
	sectionNum++
	if len(d.Risks) > 0 {
		sb.WriteString(fmt.Sprintf("## %d. Risks\n\n", sectionNum))
		sb.WriteString("| ID | Risk | Probability | Impact | Mitigation |\n")
		sb.WriteString("|----|------|-------------|--------|------------|\n")
		for _, r := range d.Risks {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				r.ID, r.Description, r.Probability, r.Impact, r.Mitigation))
		}
		sb.WriteString("\n---\n\n")
		sectionNum++
	}

	// Assumptions
	if len(d.Assumptions) > 0 {
		sb.WriteString(fmt.Sprintf("## %d. Assumptions\n\n", sectionNum))
		sb.WriteString("| ID | Assumption | Validated | Risk if Wrong |\n")
		sb.WriteString("|----|------------|-----------|---------------|\n")
		for _, a := range d.Assumptions {
			validated := "No"
			if a.Validated {
				validated = "Yes"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				a.ID, a.Description, validated, a.Risk))
		}
		sb.WriteString("\n---\n\n")
		sectionNum++
	}

	// Custom sections
	for _, cs := range d.CustomSections {
		sb.WriteString(fmt.Sprintf("## %d. %s\n\n", sectionNum, cs.Title))
		if content, ok := cs.Content.(string); ok {
			sb.WriteString(content)
			sb.WriteString("\n\n")
		}
		sb.WriteString("---\n\n")
		sectionNum++
	}

	// Glossary
	if len(d.Glossary) > 0 {
		sb.WriteString(fmt.Sprintf("## %d. Glossary\n\n", sectionNum))
		sb.WriteString("| Term | Definition |\n")
		sb.WriteString("|------|------------|\n")
		for _, g := range d.Glossary {
			term := g.Term
			if g.Acronym != "" {
				term = fmt.Sprintf("%s (%s)", g.Term, g.Acronym)
			}
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", term, g.Definition))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (d *Document) generateFrontmatter(opts MarkdownOptions) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: \"%s\"\n", d.Metadata.Title))

	if len(d.Metadata.Authors) > 0 {
		sb.WriteString(fmt.Sprintf("author: \"%s\"\n", d.Metadata.Authors[0].Name))
	}

	sb.WriteString(fmt.Sprintf("date: \"%s\"\n", d.Metadata.CreatedAt.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("version: \"%s\"\n", d.Metadata.Version))
	sb.WriteString(fmt.Sprintf("status: \"%s\"\n", d.Metadata.Status))
	sb.WriteString(fmt.Sprintf("geometry: margin=%s\n", opts.Margin))
	sb.WriteString(fmt.Sprintf("mainfont: \"%s\"\n", opts.MainFont))
	sb.WriteString(fmt.Sprintf("sansfont: \"%s\"\n", opts.SansFont))
	sb.WriteString(fmt.Sprintf("monofont: \"%s\"\n", opts.MonoFont))
	sb.WriteString(fmt.Sprintf("fontfamily: %s\n", opts.FontFamily))
	sb.WriteString("---\n\n")

	return sb.String()
}
