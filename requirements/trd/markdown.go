package trd

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/common"
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

// ToMarkdown converts the TRD to markdown format.
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
	sb.WriteString("\n")

	// Related Documents
	if len(d.Metadata.RelatedDocuments) > 0 {
		sb.WriteString("**Related Documents:**\n\n")
		sb.WriteString("| Document | Relationship | Description |\n")
		sb.WriteString("|----------|--------------|-------------|\n")
		for _, doc := range d.Metadata.RelatedDocuments {
			title := doc.Title
			if doc.URL != "" {
				title = fmt.Sprintf("[%s](%s)", doc.Title, doc.URL)
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", title, doc.Relationship, doc.Description))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Executive Summary
	sb.WriteString("## 1. Executive Summary\n\n")
	sb.WriteString("### 1.1 Purpose\n\n")
	sb.WriteString(d.ExecutiveSummary.Purpose)
	sb.WriteString("\n\n")

	sb.WriteString("### 1.2 Scope\n\n")
	sb.WriteString(d.ExecutiveSummary.Scope)
	sb.WriteString("\n\n")

	sb.WriteString("### 1.3 Technical Approach\n\n")
	sb.WriteString(d.ExecutiveSummary.TechnicalApproach)
	sb.WriteString("\n\n")

	if len(d.ExecutiveSummary.KeyDecisions) > 0 {
		sb.WriteString("### 1.4 Key Technical Decisions\n\n")
		for _, dec := range d.ExecutiveSummary.KeyDecisions {
			sb.WriteString(fmt.Sprintf("- %s\n", dec))
		}
		sb.WriteString("\n")
	}

	if len(d.ExecutiveSummary.OutOfScope) > 0 {
		sb.WriteString("### 1.5 Out of Scope\n\n")
		for _, item := range d.ExecutiveSummary.OutOfScope {
			sb.WriteString(fmt.Sprintf("- %s\n", item))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")

	// Architecture
	sb.WriteString("## 2. Architecture\n\n")
	sb.WriteString("### 2.1 Overview\n\n")
	sb.WriteString(d.Architecture.Overview)
	sb.WriteString("\n\n")

	if len(d.Architecture.Principles) > 0 {
		sb.WriteString("### 2.2 Architecture Principles\n\n")
		for _, p := range d.Architecture.Principles {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\n")
	}

	if len(d.Architecture.Patterns) > 0 {
		sb.WriteString("### 2.3 Architecture Patterns\n\n")
		for _, p := range d.Architecture.Patterns {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\n")
	}

	if len(d.Architecture.Components) > 0 {
		sb.WriteString("### 2.4 Components\n\n")
		sb.WriteString("| ID | Component | Type | Technology | Description |\n")
		sb.WriteString("|----|-----------|------|------------|-------------|\n")
		for _, c := range d.Architecture.Components {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				c.ID, c.Name, c.Type, c.Technology, c.Description))
		}
		sb.WriteString("\n")

		// Component details
		for _, c := range d.Architecture.Components {
			if len(c.Responsibilities) > 0 || len(c.Dependencies) > 0 {
				sb.WriteString(fmt.Sprintf("#### %s: %s\n\n", c.ID, c.Name))
				if len(c.Responsibilities) > 0 {
					sb.WriteString("**Responsibilities:**\n")
					for _, r := range c.Responsibilities {
						sb.WriteString(fmt.Sprintf("- %s\n", r))
					}
				}
				if len(c.Dependencies) > 0 {
					sb.WriteString(fmt.Sprintf("\n**Dependencies:** %s\n", strings.Join(c.Dependencies, ", ")))
				}
				sb.WriteString("\n")
			}
		}
	}

	if len(d.Architecture.DataFlows) > 0 {
		sb.WriteString("### 2.5 Data Flows\n\n")
		sb.WriteString("| Name | Source | Destination | Protocol | Description |\n")
		sb.WriteString("|------|--------|-------------|----------|-------------|\n")
		for _, df := range d.Architecture.DataFlows {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				df.Name, df.Source, df.Destination, df.Protocol, df.Description))
		}
		sb.WriteString("\n")
	}

	if len(d.Architecture.ArchDecisions) > 0 {
		sb.WriteString("### 2.6 Architecture Decision Records\n\n")
		for _, adr := range d.Architecture.ArchDecisions {
			sb.WriteString(fmt.Sprintf("#### ADR-%s: %s\n\n", adr.ID, adr.Title))
			sb.WriteString(fmt.Sprintf("**Status:** %s\n\n", adr.Status))
			sb.WriteString("**Context:**\n\n")
			sb.WriteString(adr.Context)
			sb.WriteString("\n\n**Decision:**\n\n")
			sb.WriteString(adr.Decision)
			sb.WriteString("\n\n")
			if len(adr.Consequences) > 0 {
				sb.WriteString("**Consequences:**\n")
				for _, c := range adr.Consequences {
					sb.WriteString(fmt.Sprintf("- %s\n", c))
				}
				sb.WriteString("\n")
			}
		}
	}

	sb.WriteString("---\n\n")

	// Technology Stack
	sb.WriteString("## 3. Technology Stack\n\n")
	d.writeTechSection(&sb, "3.1 Languages", d.TechnologyStack.Languages)
	d.writeTechSection(&sb, "3.2 Frameworks", d.TechnologyStack.Frameworks)
	d.writeTechSection(&sb, "3.3 Databases", d.TechnologyStack.Databases)
	d.writeTechSection(&sb, "3.4 Message Queues", d.TechnologyStack.MessageQueues)
	d.writeTechSection(&sb, "3.5 Caching", d.TechnologyStack.Caching)
	d.writeTechSection(&sb, "3.6 Infrastructure", d.TechnologyStack.Infrastructure)
	d.writeTechSection(&sb, "3.7 Monitoring", d.TechnologyStack.Monitoring)
	d.writeTechSection(&sb, "3.8 CI/CD", d.TechnologyStack.CICD)

	sb.WriteString("---\n\n")

	// API Specifications
	if len(d.APISpecifications) > 0 {
		sb.WriteString("## 4. API Specifications\n\n")
		for i, api := range d.APISpecifications {
			sb.WriteString(fmt.Sprintf("### 4.%d %s\n\n", i+1, api.Name))
			sb.WriteString("| Attribute | Value |\n")
			sb.WriteString("|-----------|-------|\n")
			sb.WriteString(fmt.Sprintf("| **Type** | %s |\n", api.Type))
			if api.Version != "" {
				sb.WriteString(fmt.Sprintf("| **Version** | %s |\n", api.Version))
			}
			if api.BaseURL != "" {
				sb.WriteString(fmt.Sprintf("| **Base URL** | %s |\n", api.BaseURL))
			}
			if api.Auth != "" {
				sb.WriteString(fmt.Sprintf("| **Authentication** | %s |\n", api.Auth))
			}
			if api.RateLimit != "" {
				sb.WriteString(fmt.Sprintf("| **Rate Limit** | %s |\n", api.RateLimit))
			}
			sb.WriteString("\n")

			if api.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", api.Description))
			}

			if len(api.Endpoints) > 0 {
				sb.WriteString("**Endpoints:**\n\n")
				sb.WriteString("| Method | Path | Description |\n")
				sb.WriteString("|--------|------|-------------|\n")
				for _, ep := range api.Endpoints {
					sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", ep.Method, ep.Path, ep.Description))
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("---\n\n")
	}

	sectionNum := 5
	if len(d.APISpecifications) == 0 {
		sectionNum = 4
	}

	// Security Design
	sb.WriteString(fmt.Sprintf("## %d. Security Design\n\n", sectionNum))
	sb.WriteString(fmt.Sprintf("### %d.1 Overview\n\n", sectionNum))
	sb.WriteString(d.SecurityDesign.Overview)
	sb.WriteString("\n\n")

	if d.SecurityDesign.AuthN != nil {
		sb.WriteString(fmt.Sprintf("### %d.2 Authentication\n\n", sectionNum))
		sb.WriteString(fmt.Sprintf("- **Method:** %s\n", d.SecurityDesign.AuthN.Method))
		if d.SecurityDesign.AuthN.Provider != "" {
			sb.WriteString(fmt.Sprintf("- **Provider:** %s\n", d.SecurityDesign.AuthN.Provider))
		}
		sb.WriteString(fmt.Sprintf("- **MFA Required:** %v\n", d.SecurityDesign.AuthN.MFA))
		if d.SecurityDesign.AuthN.Details != "" {
			sb.WriteString(fmt.Sprintf("\n%s\n", d.SecurityDesign.AuthN.Details))
		}
		sb.WriteString("\n")
	}

	if d.SecurityDesign.AuthZ != nil {
		sb.WriteString(fmt.Sprintf("### %d.3 Authorization\n\n", sectionNum))
		sb.WriteString(fmt.Sprintf("- **Model:** %s\n", d.SecurityDesign.AuthZ.Model))
		if len(d.SecurityDesign.AuthZ.Roles) > 0 {
			sb.WriteString(fmt.Sprintf("- **Roles:** %s\n", strings.Join(d.SecurityDesign.AuthZ.Roles, ", ")))
		}
		if d.SecurityDesign.AuthZ.Details != "" {
			sb.WriteString(fmt.Sprintf("\n%s\n", d.SecurityDesign.AuthZ.Details))
		}
		sb.WriteString("\n")
	}

	if d.SecurityDesign.Encryption != nil {
		sb.WriteString(fmt.Sprintf("### %d.4 Encryption\n\n", sectionNum))
		if d.SecurityDesign.Encryption.AtRest != "" {
			sb.WriteString(fmt.Sprintf("- **At Rest:** %s\n", d.SecurityDesign.Encryption.AtRest))
		}
		if d.SecurityDesign.Encryption.InTransit != "" {
			sb.WriteString(fmt.Sprintf("- **In Transit:** %s\n", d.SecurityDesign.Encryption.InTransit))
		}
		if d.SecurityDesign.Encryption.KeyMgmt != "" {
			sb.WriteString(fmt.Sprintf("- **Key Management:** %s\n", d.SecurityDesign.Encryption.KeyMgmt))
		}
		sb.WriteString("\n")
	}

	if len(d.SecurityDesign.Compliance) > 0 {
		sb.WriteString(fmt.Sprintf("### %d.5 Compliance\n\n", sectionNum))
		sb.WriteString(fmt.Sprintf("**Standards:** %s\n\n", strings.Join(d.SecurityDesign.Compliance, ", ")))
	}

	if len(d.SecurityDesign.ThreatModel) > 0 {
		sb.WriteString(fmt.Sprintf("### %d.6 Threat Model\n\n", sectionNum))
		sb.WriteString("| ID | Threat | Category | Likelihood | Impact | Mitigation |\n")
		sb.WriteString("|----|--------|----------|------------|--------|------------|\n")
		for _, t := range d.SecurityDesign.ThreatModel {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				t.ID, t.Name, t.Category, t.Likelihood, t.Impact, t.Mitigation))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sectionNum++

	// Performance
	sb.WriteString(fmt.Sprintf("## %d. Performance Requirements\n\n", sectionNum))
	if d.Performance.Overview != "" {
		sb.WriteString(d.Performance.Overview)
		sb.WriteString("\n\n")
	}

	if len(d.Performance.Requirements) > 0 {
		sb.WriteString("| ID | Requirement | Metric | Target | Priority |\n")
		sb.WriteString("|----|-------------|--------|--------|----------|\n")
		for _, r := range d.Performance.Requirements {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				r.ID, r.Name, r.Metric, r.Target, r.Priority))
		}
		sb.WriteString("\n")
	}

	if len(d.Performance.Optimizations) > 0 {
		sb.WriteString("**Optimization Strategies:**\n\n")
		for _, opt := range d.Performance.Optimizations {
			sb.WriteString(fmt.Sprintf("- %s\n", opt))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sectionNum++

	// Deployment
	sb.WriteString(fmt.Sprintf("## %d. Deployment\n\n", sectionNum))
	sb.WriteString(fmt.Sprintf("### %d.1 Overview\n\n", sectionNum))
	sb.WriteString(d.Deployment.Overview)
	sb.WriteString("\n\n")

	if d.Deployment.Strategy != "" {
		sb.WriteString(fmt.Sprintf("**Deployment Strategy:** %s\n\n", d.Deployment.Strategy))
	}
	if d.Deployment.Infrastructure != "" {
		sb.WriteString(fmt.Sprintf("**Infrastructure:** %s\n\n", d.Deployment.Infrastructure))
	}

	if len(d.Deployment.Environments) > 0 {
		sb.WriteString(fmt.Sprintf("### %d.2 Environments\n\n", sectionNum))
		sb.WriteString("| Environment | Purpose | URL | Resources |\n")
		sb.WriteString("|-------------|---------|-----|----------|\n")
		for _, env := range d.Deployment.Environments {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				env.Name, env.Purpose, env.URL, env.Resources))
		}
		sb.WriteString("\n")
	}

	if len(d.Deployment.Regions) > 0 {
		sb.WriteString(fmt.Sprintf("**Regions:** %s\n\n", strings.Join(d.Deployment.Regions, ", ")))
	}

	if d.Deployment.HA != "" {
		sb.WriteString(fmt.Sprintf("### %d.3 High Availability\n\n%s\n\n", sectionNum, d.Deployment.HA))
	}

	if d.Deployment.DR != "" {
		sb.WriteString(fmt.Sprintf("### %d.4 Disaster Recovery\n\n%s\n\n", sectionNum, d.Deployment.DR))
	}

	sb.WriteString("---\n\n")
	sectionNum++

	// Integrations
	if len(d.Integration) > 0 {
		sb.WriteString(fmt.Sprintf("## %d. Integrations\n\n", sectionNum))
		sb.WriteString("| Name | Type | Direction | Protocol | Auth | Description |\n")
		sb.WriteString("|------|------|-----------|----------|------|-------------|\n")
		for _, intg := range d.Integration {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				intg.Name, intg.Type, intg.Direction, intg.Protocol, intg.AuthMethod, intg.Description))
		}
		sb.WriteString("\n---\n\n")
		sectionNum++
	}

	// Testing
	if d.Testing != nil {
		sb.WriteString(fmt.Sprintf("## %d. Testing Strategy\n\n", sectionNum))
		sb.WriteString(d.Testing.Strategy)
		sb.WriteString("\n\n")

		sb.WriteString("| Test Type | Approach |\n")
		sb.WriteString("|-----------|----------|\n")
		if d.Testing.UnitTests != "" {
			sb.WriteString(fmt.Sprintf("| Unit Tests | %s |\n", d.Testing.UnitTests))
		}
		if d.Testing.Integration != "" {
			sb.WriteString(fmt.Sprintf("| Integration Tests | %s |\n", d.Testing.Integration))
		}
		if d.Testing.E2E != "" {
			sb.WriteString(fmt.Sprintf("| E2E Tests | %s |\n", d.Testing.E2E))
		}
		if d.Testing.Performance != "" {
			sb.WriteString(fmt.Sprintf("| Performance Tests | %s |\n", d.Testing.Performance))
		}
		if d.Testing.Security != "" {
			sb.WriteString(fmt.Sprintf("| Security Tests | %s |\n", d.Testing.Security))
		}
		sb.WriteString("\n")

		if d.Testing.Coverage != "" {
			sb.WriteString(fmt.Sprintf("**Coverage Requirements:** %s\n\n", d.Testing.Coverage))
		}

		sb.WriteString("---\n\n")
		sectionNum++
	}

	// Risks
	if len(d.Risks) > 0 {
		sb.WriteString(fmt.Sprintf("## %d. Technical Risks\n\n", sectionNum))
		sb.WriteString("| ID | Risk | Probability | Impact | Mitigation | Owner |\n")
		sb.WriteString("|----|------|-------------|--------|------------|-------|\n")
		for _, r := range d.Risks {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				r.ID, r.Description, r.Probability, r.Impact, r.Mitigation, r.Owner))
		}
		sb.WriteString("\n---\n\n")
		sectionNum++
	}

	// Constraints
	if len(d.Constraints) > 0 {
		sb.WriteString(fmt.Sprintf("## %d. Constraints\n\n", sectionNum))
		sb.WriteString("| ID | Type | Constraint | Impact |\n")
		sb.WriteString("|----|------|------------|--------|\n")
		for _, c := range d.Constraints {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				c.ID, c.Type, c.Description, c.Impact))
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

func (d *Document) writeTechSection(sb *strings.Builder, title string, techs []Technology) {
	if len(techs) == 0 {
		return
	}
	fmt.Fprintf(sb, "### %s\n\n", title)
	sb.WriteString("| Technology | Version | Purpose | Rationale |\n")
	sb.WriteString("|------------|---------|---------|----------|\n")
	for _, t := range techs {
		fmt.Fprintf(sb, "| %s | %s | %s | %s |\n",
			t.Name, t.Version, t.Purpose, t.Rationale)
	}
	sb.WriteString("\n")
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
