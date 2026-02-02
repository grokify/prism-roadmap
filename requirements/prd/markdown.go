package prd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/grokify/structured-plan/common"
)

// MarkdownOptions configures markdown generation.
type MarkdownOptions struct {
	// IncludeFrontmatter adds YAML frontmatter for Pandoc
	IncludeFrontmatter bool
	// Margin sets the page margin (e.g., "2cm")
	Margin string
	// MainFont sets the main font family
	MainFont string
	// SansFont sets the sans-serif font family
	SansFont string
	// MonoFont sets the monospace font family
	MonoFont string
	// FontFamily sets the LaTeX font family (e.g., "helvet")
	FontFamily string
	// DescriptionMaxLen sets the max length for description fields in tables (default: 0, no limit)
	DescriptionMaxLen int
	// IncludeSwimlaneTable adds a swimlane view of the roadmap (phases as columns, deliverable types as rows)
	IncludeSwimlaneTable bool
	// RoadmapTableOptions configures the swimlane/roadmap table generation
	RoadmapTableOptions *RoadmapTableOptions
	// IncludeTOC adds a Table of Contents with internal links (default: true)
	IncludeTOC *bool
}

// DefaultDescriptionMaxLen is the default maximum length for description fields in tables.
// A value of 0 means no truncation (full text is displayed).
const DefaultDescriptionMaxLen = 0

// DefaultMarkdownOptions returns sensible defaults for markdown generation.
// By default, no text truncation is applied (DescriptionMaxLen = 0).
// Fonts default to DejaVu family for Unicode/emoji support in Pandoc PDF output.
func DefaultMarkdownOptions() MarkdownOptions {
	return MarkdownOptions{
		IncludeFrontmatter: true,
		Margin:             "2cm",
		MainFont:           "DejaVu Sans",
		SansFont:           "DejaVu Sans",
		MonoFont:           "DejaVu Sans Mono",
		FontFamily:         "",
		DescriptionMaxLen:  DefaultDescriptionMaxLen,
	}
}

// ToMarkdown converts a PRD Document to markdown format.
func (d *Document) ToMarkdown(opts MarkdownOptions) string {
	var sb strings.Builder

	// YAML Frontmatter
	if opts.IncludeFrontmatter {
		sb.WriteString(d.generateFrontmatter(opts))
	}

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", d.Metadata.Title))

	// Metadata table
	sb.WriteString(d.generateMetadataTable())

	// Table of Contents (default: enabled)
	includeTOC := opts.IncludeTOC == nil || *opts.IncludeTOC
	if includeTOC {
		sb.WriteString(d.generateTableOfContents(opts))
	}

	// Executive Summary
	sb.WriteString(d.generateExecutiveSummary())

	// Objectives
	sb.WriteString(d.generateObjectives())

	// Personas
	sb.WriteString(d.generatePersonas())

	// User Stories
	sb.WriteString(d.generateUserStories())

	// Requirements
	sb.WriteString(d.generateRequirements(opts))

	// Roadmap
	sb.WriteString(d.generateRoadmap(opts))

	// Optional sections
	if d.TechArchitecture != nil {
		sb.WriteString(d.generateTechArchitecture())
	}

	if d.Assumptions != nil {
		sb.WriteString(d.generateAssumptions())
	}

	if len(d.InScope) > 0 {
		sb.WriteString(d.generateInScope())
	}

	if len(d.OutOfScope) > 0 {
		sb.WriteString(d.generateOutOfScope())
	}

	if len(d.Risks) > 0 {
		sb.WriteString(d.generateRisks())
	}

	if len(d.OpenItems) > 0 {
		sb.WriteString(d.generateOpenItems())
	}

	if d.CurrentState != nil {
		sb.WriteString(d.generateCurrentState())
	}

	if d.SecurityModel != nil {
		sb.WriteString(d.generateSecurityModel())
	}

	if len(d.Appendices) > 0 {
		sb.WriteString(d.generateAppendices())
	}

	if len(d.Glossary) > 0 {
		sb.WriteString(d.generateGlossary())
	}

	if len(d.RelatedDocuments) > 0 {
		sb.WriteString(d.generateRelatedDocuments())
	}

	// Custom sections
	if len(d.CustomSections) > 0 {
		sb.WriteString(d.generateCustomSections())
	}

	// Footer
	sb.WriteString("\n---\n\n*Generated from structured PRD JSON format*\n")

	return sb.String()
}

func (d *Document) generateFrontmatter(opts MarkdownOptions) string {
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %q\n", d.Metadata.Title))

	// Authors
	if len(d.Metadata.Authors) > 0 {
		names := make([]string, len(d.Metadata.Authors))
		for i, a := range d.Metadata.Authors {
			names[i] = a.Name
		}
		sb.WriteString(fmt.Sprintf("author: %q\n", strings.Join(names, ", ")))
	}

	// Date
	if !d.Metadata.UpdatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("date: %q\n", d.Metadata.UpdatedAt.Format("2006-01-02")))
	} else if !d.Metadata.CreatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("date: %q\n", d.Metadata.CreatedAt.Format("2006-01-02")))
	}

	sb.WriteString(fmt.Sprintf("version: %q\n", d.Metadata.Version))
	sb.WriteString(fmt.Sprintf("status: %q\n", d.Metadata.Status))

	// Pandoc/LaTeX settings
	if opts.Margin != "" {
		sb.WriteString(fmt.Sprintf("geometry: margin=%s\n", opts.Margin))
	}
	if opts.MainFont != "" {
		sb.WriteString(fmt.Sprintf("mainfont: %q\n", opts.MainFont))
	}
	if opts.SansFont != "" {
		sb.WriteString(fmt.Sprintf("sansfont: %q\n", opts.SansFont))
	}
	if opts.MonoFont != "" {
		sb.WriteString(fmt.Sprintf("monofont: %q\n", opts.MonoFont))
	}
	if opts.FontFamily != "" {
		sb.WriteString(fmt.Sprintf("fontfamily: %s\n", opts.FontFamily))
	}

	sb.WriteString("header-includes:\n")
	sb.WriteString("  - \\usepackage{fontspec}\n")
	sb.WriteString("  - \\setmainfont{DejaVu Sans}\n")
	sb.WriteString("  - \\setsansfont{DejaVu Sans}\n")
	sb.WriteString("  - \\setmonofont{DejaVu Sans Mono}\n")
	sb.WriteString("---\n\n")

	return sb.String()
}

func (d *Document) generateMetadataTable() string {
	var sb strings.Builder
	sb.WriteString("| Field | Value |\n")
	sb.WriteString("|-------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **ID** | %s |\n", d.Metadata.ID))
	sb.WriteString(fmt.Sprintf("| **Version** | %s |\n", d.Metadata.Version))
	sb.WriteString(fmt.Sprintf("| **Status** | %s |\n", d.Metadata.Status))

	if !d.Metadata.CreatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("| **Created** | %s |\n", d.Metadata.CreatedAt.Format("2006-01-02")))
	}
	if !d.Metadata.UpdatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("| **Updated** | %s |\n", d.Metadata.UpdatedAt.Format("2006-01-02")))
	}

	if len(d.Metadata.Authors) > 0 {
		sb.WriteString(fmt.Sprintf("| **Author(s)** | %s |\n", common.FormatPeopleMarkdown(d.Metadata.Authors)))
	}

	if len(d.Metadata.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("| **Tags** | %s |\n", strings.Join(d.Metadata.Tags, ", ")))
	}

	sb.WriteString("\n")

	if d.Metadata.SemanticVersioning {
		sb.WriteString("*This document uses [Semantic Versioning](https://semver.org/).*\n\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateTableOfContents(_ MarkdownOptions) string {
	var sb strings.Builder
	sb.WriteString("## Table of Contents\n\n")

	// Fixed sections (always present)
	sb.WriteString("1. [Executive Summary](#1-executive-summary)\n")
	sb.WriteString("2. [Objectives and Goals](#2-objectives-and-goals)\n")
	sb.WriteString("3. [Personas](#3-personas)\n")
	sb.WriteString("4. [User Stories](#4-user-stories)\n")
	sb.WriteString("5. [Functional Requirements](#5-functional-requirements)\n")
	sb.WriteString("6. [Non-Functional Requirements](#6-non-functional-requirements)\n")
	sb.WriteString("7. [Roadmap](#7-roadmap)\n")

	// Optional sections - track section number
	sectionNum := 8

	if d.TechArchitecture != nil {
		sb.WriteString(fmt.Sprintf("%d. [Technical Architecture](#technical-architecture)\n", sectionNum))
		sectionNum++
	}

	if d.Assumptions != nil {
		sb.WriteString(fmt.Sprintf("%d. [Assumptions and Constraints](#assumptions-and-constraints)\n", sectionNum))
		sectionNum++
	}

	if len(d.InScope) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [In Scope](#in-scope)\n", sectionNum))
		sectionNum++
	}

	if len(d.OutOfScope) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [Out of Scope](#out-of-scope)\n", sectionNum))
		sectionNum++
	}

	if len(d.Risks) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [Risk Assessment](#risk-assessment)\n", sectionNum))
		sectionNum++
	}

	if len(d.OpenItems) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [Open Items](#open-items)\n", sectionNum))
		sectionNum++
	}

	if d.CurrentState != nil {
		sb.WriteString(fmt.Sprintf("%d. [Current State](#current-state)\n", sectionNum))
		sectionNum++
	}

	if d.SecurityModel != nil {
		sb.WriteString(fmt.Sprintf("%d. [Security Model](#security-model)\n", sectionNum))
		sectionNum++
	}

	if len(d.Appendices) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [Appendices](#appendices)\n", sectionNum))
		sectionNum++
	}

	if len(d.Glossary) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [Glossary](#glossary)\n", sectionNum))
		sectionNum++
	}

	if len(d.RelatedDocuments) > 0 {
		sb.WriteString(fmt.Sprintf("%d. [Related Documents](#related-documents)\n", sectionNum))
		sectionNum++
	}

	// Custom sections
	for _, cs := range d.CustomSections {
		slug := toSlug(cs.Title)
		sb.WriteString(fmt.Sprintf("%d. [%s](#%s)\n", sectionNum, cs.Title, slug))
		sectionNum++
	}

	sb.WriteString("\n---\n\n")
	return sb.String()
}

// toSlug converts a string to a URL-friendly slug for markdown anchors.
func toSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (d *Document) generateExecutiveSummary() string {
	var sb strings.Builder
	sb.WriteString("## 1. Executive Summary\n\n")

	sb.WriteString("### 1.1 Problem Statement\n\n")
	sb.WriteString(d.ExecutiveSummary.ProblemStatement + "\n\n")

	sb.WriteString("### 1.2 Proposed Solution\n\n")
	sb.WriteString(d.ExecutiveSummary.ProposedSolution + "\n\n")

	if len(d.ExecutiveSummary.ExpectedOutcomes) > 0 {
		sb.WriteString("### 1.3 Expected Outcomes\n\n")
		for _, outcome := range d.ExecutiveSummary.ExpectedOutcomes {
			sb.WriteString(fmt.Sprintf("- %s\n", outcome))
		}
		sb.WriteString("\n")
	}

	if d.ExecutiveSummary.TargetAudience != "" {
		sb.WriteString("### 1.4 Target Audience\n\n")
		sb.WriteString(d.ExecutiveSummary.TargetAudience + "\n\n")
	}

	if d.ExecutiveSummary.ValueProposition != "" {
		sb.WriteString("### 1.5 Value Proposition\n\n")
		sb.WriteString(d.ExecutiveSummary.ValueProposition + "\n\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateObjectives() string {
	var sb strings.Builder
	sb.WriteString("## 2. Objectives and Goals\n\n")

	if len(d.Objectives.OKRs) > 0 {
		sb.WriteString(d.generateOKRs())
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateOKRs() string {
	var sb strings.Builder

	// Objectives overview - quick scan of all objectives
	sb.WriteString("### 2.1 Objectives Overview\n\n")
	for i, okr := range d.Objectives.OKRs {
		obj := okr.Objective
		timeframe := ""
		if obj.Timeframe != "" {
			timeframe = fmt.Sprintf(" (%s)", obj.Timeframe)
		}
		sb.WriteString(fmt.Sprintf("%d. **%s**%s\n", i+1, obj.Description, timeframe))
	}
	sb.WriteString("\n")

	// Detailed OKRs with Key Results
	sb.WriteString("### 2.2 OKRs (Objectives and Key Results)\n\n")

	for i, okr := range d.Objectives.OKRs {
		obj := okr.Objective

		// Objective header with metadata
		// Use Title if set, otherwise fall back to Description for backward compatibility
		objTitle := obj.Title
		if objTitle == "" {
			objTitle = obj.Description
		}
		timeframe := ""
		if obj.Timeframe != "" {
			timeframe = fmt.Sprintf(" (%s)", obj.Timeframe)
		}
		sb.WriteString(fmt.Sprintf("#### Objective %d: %s%s\n\n", i+1, objTitle, timeframe))

		// Objective metadata table
		if obj.Owner != "" || obj.Category != "" || len(obj.AlignedWith) > 0 {
			sb.WriteString("| Attribute | Value |\n")
			sb.WriteString("|-----------|-------|\n")
			if obj.Category != "" {
				sb.WriteString(fmt.Sprintf("| **Category** | %s |\n", obj.Category))
			}
			if obj.Owner != "" {
				sb.WriteString(fmt.Sprintf("| **Owner** | %s |\n", obj.Owner))
			}
			if len(obj.AlignedWith) > 0 {
				sb.WriteString(fmt.Sprintf("| **Aligned With** | %s |\n", strings.Join(obj.AlignedWith, ", ")))
			}
			if obj.Rationale != "" {
				sb.WriteString(fmt.Sprintf("| **Rationale** | %s |\n", obj.Rationale))
			}
			sb.WriteString("\n")
		}

		// Key Results table
		sb.WriteString("**Key Results:**\n\n")
		sb.WriteString("| KR | Description | Baseline | Target | Current | Confidence |\n")
		sb.WriteString("|----|-------------|----------|--------|---------|------------|\n")

		for j, kr := range okr.KeyResults {
			baseline := kr.Baseline
			if baseline == "" {
				baseline = "-"
			}
			current := kr.Current
			if current == "" {
				current = "-"
			}
			confidence := "-"
			if kr.Confidence != "" {
				confidence = kr.Confidence
			}

			// Format with unit if present
			target := kr.Target
			if kr.Unit != "" && target != "-" {
				target = fmt.Sprintf("%s %s", kr.Target, kr.Unit)
			}

			// Use Title if set, otherwise fall back to Description for backward compatibility
			krTitle := kr.Title
			if krTitle == "" {
				krTitle = kr.Description
			}
			sb.WriteString(fmt.Sprintf("| KR%d.%d | %s | %s | %s | %s | %s |\n",
				i+1, j+1, krTitle, baseline, target, current, confidence))
		}
		sb.WriteString("\n")

		// Phase targets if present
		for _, kr := range okr.KeyResults {
			if len(kr.PhaseTargets) > 0 {
				sb.WriteString(fmt.Sprintf("**%s - Phase Targets:**\n\n", kr.Description))
				sb.WriteString("| Phase | Target | Status | Actual | Notes |\n")
				sb.WriteString("|-------|--------|--------|--------|-------|\n")
				for _, pt := range kr.PhaseTargets {
					status := pt.Status
					if status == "" {
						status = "not_started"
					}
					actual := pt.Actual
					if actual == "" {
						actual = "-"
					}
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
						pt.PhaseID, pt.Target, status, actual, pt.Notes))
				}
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

func (d *Document) generatePersonas() string {
	var sb strings.Builder
	sb.WriteString("## 3. Personas\n\n")

	for i, p := range d.Personas {
		primary := ""
		if p.IsPrimary {
			primary = " (Primary)"
		}
		sb.WriteString(fmt.Sprintf("### 3.%d %s%s\n\n", i+1, p.Name, primary))

		sb.WriteString("| Attribute | Description |\n")
		sb.WriteString("|-----------|-------------|\n")
		sb.WriteString(fmt.Sprintf("| **Role** | %s |\n", p.Role))
		sb.WriteString(fmt.Sprintf("| **Description** | %s |\n", p.Description))
		if p.TechnicalProficiency != "" {
			sb.WriteString(fmt.Sprintf("| **Technical Proficiency** | %s |\n", p.TechnicalProficiency))
		}
		sb.WriteString("\n")

		if len(p.Goals) > 0 {
			sb.WriteString("**Goals:**\n\n")
			for _, g := range p.Goals {
				sb.WriteString(fmt.Sprintf("- %s\n", g))
			}
			sb.WriteString("\n")
		}

		if len(p.PainPoints) > 0 {
			sb.WriteString("**Pain Points:**\n\n")
			for _, pp := range p.PainPoints {
				sb.WriteString(fmt.Sprintf("- %s\n", pp))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateUserStories() string {
	var sb strings.Builder
	sb.WriteString("## 4. User Stories\n\n")

	// Group by persona
	personaStories := make(map[string][]UserStory)
	for _, us := range d.UserStories {
		personaStories[us.PersonaID] = append(personaStories[us.PersonaID], us)
	}

	sectionNum := 1
	for _, p := range d.Personas {
		stories, ok := personaStories[p.ID]
		if !ok || len(stories) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("### 4.%d %s Stories\n\n", sectionNum, p.Name))
		sb.WriteString("| ID | Story | Priority | Phase |\n")
		sb.WriteString("|------|----------------------------------------|----------|-------|\n")
		for _, us := range stories {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				us.ID, us.Story(), us.Priority, us.PhaseID))
		}
		sb.WriteString("\n")
		sectionNum++
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateRequirements(opts MarkdownOptions) string {
	var sb strings.Builder

	// Functional Requirements
	sb.WriteString("## 5. Functional Requirements\n\n")

	// Group by category
	categories := make(map[string][]FunctionalRequirement)
	for _, fr := range d.Requirements.Functional {
		categories[fr.Category] = append(categories[fr.Category], fr)
	}

	// Sort category names for consistent ordering
	var categoryNames []string
	for cat := range categories {
		categoryNames = append(categoryNames, cat)
	}
	sort.Strings(categoryNames)

	sectionNum := 1
	for _, cat := range categoryNames {
		reqs := categories[cat]
		sb.WriteString(fmt.Sprintf("### 5.%d %s\n\n", sectionNum, cat))
		sb.WriteString("| ID | Title | Description | Priority | Phase |\n")
		sb.WriteString("|------|-----------------|--------------------------------------------|----------|-------|\n")
		for _, r := range reqs {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				r.ID, r.Title, truncate(r.Description, opts.DescriptionMaxLen), r.Priority, r.PhaseID))
		}
		sb.WriteString("\n")
		sectionNum++
	}

	// Non-Functional Requirements
	sb.WriteString("## 6. Non-Functional Requirements\n\n")

	// Group by category
	nfrCategories := make(map[NFRCategory][]NonFunctionalRequirement)
	for _, nfr := range d.Requirements.NonFunctional {
		nfrCategories[nfr.Category] = append(nfrCategories[nfr.Category], nfr)
	}

	nfrCategoryDisplayNames := map[NFRCategory]string{
		NFRPerformance:      "Performance",
		NFRScalability:      "Scalability",
		NFRReliability:      "Reliability",
		NFRAvailability:     "Availability",
		NFRSecurity:         "Security",
		NFRMultiTenancy:     "Multi-Tenancy",
		NFRObservability:    "Observability",
		NFRMaintainability:  "Maintainability",
		NFRUsability:        "Usability",
		NFRCompatibility:    "Compatibility",
		NFRCompliance:       "Compliance",
		NFRDisasterRecovery: "Disaster Recovery",
		NFRCostEfficiency:   "Cost Efficiency",
	}

	// Sort NFR category keys for consistent ordering
	var nfrCategoryKeys []NFRCategory
	for cat := range nfrCategories {
		nfrCategoryKeys = append(nfrCategoryKeys, cat)
	}
	sort.Slice(nfrCategoryKeys, func(i, j int) bool {
		return string(nfrCategoryKeys[i]) < string(nfrCategoryKeys[j])
	})

	sectionNum = 1
	for _, cat := range nfrCategoryKeys {
		reqs := nfrCategories[cat]
		catName := nfrCategoryDisplayNames[cat]
		if catName == "" {
			catName = string(cat)
		}
		sb.WriteString(fmt.Sprintf("### 6.%d %s\n\n", sectionNum, catName))
		sb.WriteString("| ID | Title | Target | Priority | Phase |\n")
		sb.WriteString("|----|-------|--------|----------|-------|\n")
		for _, r := range reqs {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				r.ID, r.Title, r.Target, r.Priority, r.PhaseID))
		}
		sb.WriteString("\n")
		sectionNum++
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateRoadmap(opts MarkdownOptions) string {
	var sb strings.Builder
	sb.WriteString("## 7. Roadmap\n\n")

	// Swimlane table view (phases as columns, deliverable types as rows)
	if opts.IncludeSwimlaneTable && len(d.Roadmap.Phases) > 0 {
		sb.WriteString("### 7.1 Roadmap Overview (Swimlane View)\n\n")
		tableOpts := DefaultRoadmapTableOptions()
		if opts.RoadmapTableOptions != nil {
			tableOpts = *opts.RoadmapTableOptions
		}
		// Enable OKR swimlanes by default if OKRs with PhaseTargets exist
		if len(d.Objectives.OKRs) > 0 {
			tableOpts.IncludeOKRs = true
		}
		sb.WriteString(d.ToSwimlaneTableWithOKRs(tableOpts))
		sb.WriteString("\n")
		if tableOpts.IncludeStatus {
			sb.WriteString("**Legend:**\n\n")
			sb.WriteString(StatusLegend())
			sb.WriteString("\n")
		}
		sb.WriteString("### 7.2 Phase Details\n\n")
	}

	for _, phase := range d.Roadmap.Phases {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", phase.ID, phase.Name))

		sb.WriteString(fmt.Sprintf("**Type:** %s\n\n", phase.Type))

		if len(phase.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("**Dependencies:** %s\n\n", strings.Join(phase.Dependencies, ", ")))
		}

		if len(phase.Goals) > 0 {
			sb.WriteString("**Goals:**\n\n")
			for _, g := range phase.Goals {
				sb.WriteString(fmt.Sprintf("- %s\n", g))
			}
			sb.WriteString("\n")
		}

		if len(phase.Deliverables) > 0 {
			sb.WriteString("**Deliverables:**\n\n")
			sb.WriteString("| ID | Title | Type | Status |\n")
			sb.WriteString("|----|-------|------|--------|\n")
			for _, del := range phase.Deliverables {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					del.ID, del.Title, del.Type, del.Status))
			}
			sb.WriteString("\n")
		}

		if len(phase.SuccessCriteria) > 0 {
			sb.WriteString("**Success Criteria:**\n\n")
			for _, sc := range phase.SuccessCriteria {
				sb.WriteString(fmt.Sprintf("- %s\n", sc))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

func (d *Document) generateTechArchitecture() string {
	var sb strings.Builder
	sb.WriteString("## Technical Architecture\n\n")

	if d.TechArchitecture.Overview != "" {
		sb.WriteString("### Overview\n\n")
		sb.WriteString(d.TechArchitecture.Overview + "\n\n")
	}

	// Services
	if len(d.TechArchitecture.Services) > 0 {
		sb.WriteString("### Services\n\n")
		sb.WriteString("| ID | Name | Layer | Protocol | Language | Owner |\n")
		sb.WriteString("|----|------|-------|----------|----------|-------|\n")
		for _, svc := range d.TechArchitecture.Services {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				svc.ID, svc.Name, svc.Layer, svc.Protocol, svc.Language, svc.Owner))
		}
		sb.WriteString("\n")

		// Service details
		for _, svc := range d.TechArchitecture.Services {
			if len(svc.Responsibilities) > 0 || svc.Description != "" || len(svc.Dependencies) > 0 {
				sb.WriteString(fmt.Sprintf("#### %s\n\n", svc.Name))
				if svc.Description != "" {
					sb.WriteString(svc.Description + "\n\n")
				}
				if len(svc.Responsibilities) > 0 {
					sb.WriteString("**Responsibilities:**\n\n")
					for _, r := range svc.Responsibilities {
						sb.WriteString(fmt.Sprintf("- %s\n", r))
					}
					sb.WriteString("\n")
				}
				if svc.LanguageRationale != "" {
					sb.WriteString(fmt.Sprintf("**Language Rationale:** %s\n\n", svc.LanguageRationale))
				}
				if len(svc.Dependencies) > 0 {
					sb.WriteString(fmt.Sprintf("**Dependencies:** %s\n\n", strings.Join(svc.Dependencies, ", ")))
				}
			}
		}
	}

	// APIs
	if len(d.TechArchitecture.APIs) > 0 {
		sb.WriteString("### APIs\n\n")
		for _, api := range d.TechArchitecture.APIs {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", api.Name))
			sb.WriteString(fmt.Sprintf("**Protocol:** %s", api.Protocol))
			if api.Version != "" {
				sb.WriteString(fmt.Sprintf(" | **Version:** %s", api.Version))
			}
			if api.BasePath != "" {
				sb.WriteString(fmt.Sprintf(" | **Base Path:** %s", api.BasePath))
			}
			sb.WriteString("\n\n")

			if api.Description != "" {
				sb.WriteString(api.Description + "\n\n")
			}

			if len(api.Endpoints) > 0 {
				sb.WriteString("| Method | Path | Description | Auth |\n")
				sb.WriteString("|--------|------|-------------|------|\n")
				for _, ep := range api.Endpoints {
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
						ep.Method, ep.Path, ep.Description, ep.Auth))
				}
				sb.WriteString("\n")
			}

			if api.OpenAPISpec != "" {
				sb.WriteString(fmt.Sprintf("**OpenAPI Spec:** %s\n\n", api.OpenAPISpec))
			}
			if api.ProtobufSpec != "" {
				sb.WriteString(fmt.Sprintf("**Protobuf Spec:** %s\n\n", api.ProtobufSpec))
			}
		}
	}

	// Storage Architecture
	if len(d.TechArchitecture.StorageArchitecture) > 0 {
		sb.WriteString("### Storage Architecture\n\n")
		sb.WriteString("| Category | Purpose | Technology | Encryption | Retention |\n")
		sb.WriteString("|----------|---------|------------|------------|----------|\n")
		for _, sc := range d.TechArchitecture.StorageArchitecture {
			encryption := sc.Encryption
			if encryption == "" {
				encryption = "-"
			}
			retention := sc.Retention
			if retention == "" {
				retention = "-"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				sc.Category, sc.Purpose, sc.Technology, encryption, retention))
		}
		sb.WriteString("\n")
	}

	// GitOps
	if d.TechArchitecture.GitOps != nil && d.TechArchitecture.GitOps.Enabled {
		sb.WriteString("### GitOps Configuration\n\n")
		gitops := d.TechArchitecture.GitOps
		if gitops.Provider != "" {
			sb.WriteString(fmt.Sprintf("**Provider:** %s\n\n", gitops.Provider))
		}
		if gitops.Workflow != "" {
			sb.WriteString(fmt.Sprintf("**Workflow:** %s\n\n", gitops.Workflow))
		}
		if len(gitops.SourcesOfTruth) > 0 {
			sb.WriteString("**Sources of Truth:**\n\n")
			sb.WriteString("| Artifact | Location | GitOps Enabled | Rationale |\n")
			sb.WriteString("|----------|----------|----------------|----------|\n")
			for _, sot := range gitops.SourcesOfTruth {
				gitOpsEnabled := "No"
				if sot.GitOpsEnabled {
					gitOpsEnabled = "Yes"
				}
				rationale := sot.Rationale
				if rationale == "" {
					rationale = "-"
				}
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					sot.Artifact, sot.Location, gitOpsEnabled, rationale))
			}
			sb.WriteString("\n")
		}
	}

	// Orchestration
	if d.TechArchitecture.Orchestration != nil {
		sb.WriteString("### Workflow Orchestration\n\n")
		orch := d.TechArchitecture.Orchestration
		if orch.Description != "" {
			sb.WriteString(orch.Description + "\n\n")
		}

		if orch.ShortLived != nil {
			sb.WriteString("**Short-Lived Workflows:**\n\n")
			sb.WriteString(fmt.Sprintf("- **Engine:** %s\n", orch.ShortLived.Name))
			if orch.ShortLived.Language != "" {
				sb.WriteString(fmt.Sprintf("- **Language:** %s\n", orch.ShortLived.Language))
			}
			if orch.ShortLived.Rationale != "" {
				sb.WriteString(fmt.Sprintf("- **Rationale:** %s\n", orch.ShortLived.Rationale))
			}
			if len(orch.ShortLived.UseCases) > 0 {
				sb.WriteString("- **Use Cases:**\n")
				for _, uc := range orch.ShortLived.UseCases {
					sb.WriteString(fmt.Sprintf("  - %s\n", uc))
				}
			}
			sb.WriteString("\n")
		}

		if orch.LongRunning != nil {
			sb.WriteString("**Long-Running Workflows:**\n\n")
			sb.WriteString(fmt.Sprintf("- **Engine:** %s\n", orch.LongRunning.Name))
			if orch.LongRunning.Language != "" {
				sb.WriteString(fmt.Sprintf("- **Language:** %s\n", orch.LongRunning.Language))
			}
			if orch.LongRunning.Rationale != "" {
				sb.WriteString(fmt.Sprintf("- **Rationale:** %s\n", orch.LongRunning.Rationale))
			}
			if len(orch.LongRunning.UseCases) > 0 {
				sb.WriteString("- **Use Cases:**\n")
				for _, uc := range orch.LongRunning.UseCases {
					sb.WriteString(fmt.Sprintf("  - %s\n", uc))
				}
			}
			sb.WriteString("\n")
		}
	}

	if len(d.TechArchitecture.IntegrationPoints) > 0 {
		sb.WriteString("### Integration Points\n\n")
		sb.WriteString("| ID | Name | Type | Description | Auth Method |\n")
		sb.WriteString("|----|------|------|-------------|-------------|\n")
		for _, ip := range d.TechArchitecture.IntegrationPoints {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				ip.ID, ip.Name, ip.Type, ip.Description, ip.AuthMethod))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateAssumptions() string {
	var sb strings.Builder
	sb.WriteString("## Assumptions and Constraints\n\n")

	if len(d.Assumptions.Assumptions) > 0 {
		sb.WriteString("### Assumptions\n\n")
		sb.WriteString("| ID | Assumption | Risk if Invalid |\n")
		sb.WriteString("|----|------------|------------------|\n")
		for _, a := range d.Assumptions.Assumptions {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				a.ID, a.Description, a.Risk))
		}
		sb.WriteString("\n")
	}

	if len(d.Assumptions.Constraints) > 0 {
		sb.WriteString("### Constraints\n\n")
		sb.WriteString("| ID | Type | Constraint | Impact | Mitigation |\n")
		sb.WriteString("|----|------|------------|--------|------------|\n")
		for _, c := range d.Assumptions.Constraints {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				c.ID, c.Type, c.Description, c.Impact, c.Mitigation))
		}
		sb.WriteString("\n")
	}

	if len(d.Assumptions.Dependencies) > 0 {
		sb.WriteString("### Dependencies\n\n")
		sb.WriteString("| ID | Name | Type | Status |\n")
		sb.WriteString("|----|------|------|--------|\n")
		for _, dep := range d.Assumptions.Dependencies {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dep.ID, dep.Name, dep.Type, dep.Status))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateInScope() string {
	var sb strings.Builder
	sb.WriteString("## In Scope\n\n")

	for _, item := range d.InScope {
		sb.WriteString(fmt.Sprintf("- %s\n", item))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateOutOfScope() string {
	var sb strings.Builder
	sb.WriteString("## Out of Scope\n\n")

	for _, item := range d.OutOfScope {
		sb.WriteString(fmt.Sprintf("- %s\n", item))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateRisks() string {
	var sb strings.Builder
	sb.WriteString("## Risk Assessment\n\n")

	sb.WriteString("| ID | Risk | Probability | Impact | Mitigation | Status |\n")
	sb.WriteString("|----|------|-------------|--------|------------|--------|\n")
	for _, r := range d.Risks {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			r.ID, r.Description, r.Probability, r.Impact, r.Mitigation, r.Status))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateOpenItems() string {
	var sb strings.Builder
	sb.WriteString("## Open Items\n\n")
	sb.WriteString("*The following items require decisions. Please review the options and tradeoffs.*\n\n")

	for i, item := range d.OpenItems {
		// Item header with status
		statusBadge := ""
		switch item.Status {
		case OpenItemStatusOpen:
			statusBadge = "🔴 Open"
		case OpenItemStatusInDiscussion:
			statusBadge = "🟡 In Discussion"
		case OpenItemStatusBlocked:
			statusBadge = "⛔ Blocked"
		case OpenItemStatusResolved:
			statusBadge = "✅ Resolved"
		case OpenItemStatusDeferred:
			statusBadge = "⏸️ Deferred"
		default:
			statusBadge = "🔴 Open"
		}

		sb.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, item.Title))
		sb.WriteString(fmt.Sprintf("**Status:** %s", statusBadge))
		if item.Priority != "" {
			sb.WriteString(fmt.Sprintf(" | **Priority:** %s", item.Priority))
		}
		if item.Owner != "" {
			sb.WriteString(fmt.Sprintf(" | **Owner:** %s", item.Owner))
		}
		sb.WriteString("\n\n")

		if item.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", item.Description))
		}

		if item.Context != "" {
			sb.WriteString(fmt.Sprintf("**Context:** %s\n\n", item.Context))
		}

		// Options table
		if len(item.Options) > 0 {
			sb.WriteString("#### Options\n\n")
			sb.WriteString("| Option | Description | Effort | Risk | Recommended |\n")
			sb.WriteString("|--------|-------------|--------|------|-------------|\n")
			for _, opt := range item.Options {
				recommended := ""
				if opt.Recommended {
					recommended = "⭐ Yes"
				}
				sb.WriteString(fmt.Sprintf("| **%s** | %s | %s | %s | %s |\n",
					opt.Title, opt.Description, opt.Effort, opt.Risk, recommended))
			}
			sb.WriteString("\n")

			// Detailed pros/cons for each option
			for _, opt := range item.Options {
				if len(opt.Pros) > 0 || len(opt.Cons) > 0 {
					sb.WriteString(fmt.Sprintf("**%s**", opt.Title))
					if opt.Recommended {
						sb.WriteString(" ⭐ *Recommended*")
					}
					sb.WriteString("\n\n")

					if len(opt.Pros) > 0 {
						sb.WriteString("*Pros:*\n")
						for _, pro := range opt.Pros {
							sb.WriteString(fmt.Sprintf("- ✅ %s\n", pro))
						}
					}
					if len(opt.Cons) > 0 {
						sb.WriteString("\n*Cons:*\n")
						for _, con := range opt.Cons {
							sb.WriteString(fmt.Sprintf("- ⚠️ %s\n", con))
						}
					}
					if opt.RecommendationRationale != "" {
						sb.WriteString(fmt.Sprintf("\n*Rationale:* %s\n", opt.RecommendationRationale))
					}
					sb.WriteString("\n")
				}
			}
		}

		// Resolution (if resolved)
		if item.Resolution != nil && item.Resolution.Decision != "" {
			sb.WriteString("#### Resolution\n\n")
			sb.WriteString(fmt.Sprintf("**Decision:** %s\n\n", item.Resolution.Decision))
			if item.Resolution.Rationale != "" {
				sb.WriteString(fmt.Sprintf("**Rationale:** %s\n\n", item.Resolution.Rationale))
			}
			if item.Resolution.DecidedBy != "" {
				sb.WriteString(fmt.Sprintf("**Decided by:** %s\n\n", item.Resolution.DecidedBy))
			}
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

func (d *Document) generateCurrentState() string {
	var sb strings.Builder
	sb.WriteString("## Current State\n\n")

	cs := d.CurrentState

	// Overview
	if cs.Overview != "" {
		sb.WriteString("### Overview\n\n")
		sb.WriteString(cs.Overview + "\n\n")
	}

	// Current Approaches
	if len(cs.Approaches) > 0 {
		sb.WriteString("### Current Approaches\n\n")
		for _, approach := range cs.Approaches {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", approach.Name))
			if approach.Description != "" {
				sb.WriteString(approach.Description + "\n\n")
			}
			if approach.Usage != "" {
				sb.WriteString(fmt.Sprintf("**Usage:** %s\n\n", approach.Usage))
			}
			if approach.Owner != "" {
				sb.WriteString(fmt.Sprintf("**Owner:** %s\n\n", approach.Owner))
			}
			if len(approach.Problems) > 0 {
				sb.WriteString("**Problems:**\n\n")
				for _, p := range approach.Problems {
					sb.WriteString(fmt.Sprintf("- %s\n", p))
				}
				sb.WriteString("\n")
			}
		}
	}

	// Problems with Current State
	if len(cs.Problems) > 0 {
		sb.WriteString("### Problems\n\n")
		sb.WriteString("| ID | Problem | Impact | Frequency | Affected Users |\n")
		sb.WriteString("|----|---------|--------|-----------|----------------|\n")
		for _, p := range cs.Problems {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				p.ID, p.Description, p.Impact, p.Frequency, p.AffectedUsers))
		}
		sb.WriteString("\n")
	}

	// Target State
	if cs.TargetState != "" {
		sb.WriteString("### Target State\n\n")
		sb.WriteString(cs.TargetState + "\n\n")
	}

	// Baseline Metrics
	if len(cs.Metrics) > 0 {
		sb.WriteString("### Baseline Metrics\n\n")
		sb.WriteString("| ID | Metric | Current Value | Target Value | Measurement Method |\n")
		sb.WriteString("|----|--------|---------------|--------------|--------------------|\n")
		for _, m := range cs.Metrics {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				m.ID, m.Name, m.CurrentValue, m.TargetValue, m.MeasurementMethod))
		}
		sb.WriteString("\n")
	}

	// Diagrams
	if len(cs.Diagrams) > 0 {
		sb.WriteString("### Diagrams\n\n")
		for _, diag := range cs.Diagrams {
			sb.WriteString(fmt.Sprintf("- [%s](%s)", diag.Title, diag.URL))
			if diag.Type != "" {
				sb.WriteString(fmt.Sprintf(" (%s)", diag.Type))
			}
			if diag.Description != "" {
				sb.WriteString(fmt.Sprintf(" - %s", diag.Description))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateSecurityModel() string {
	var sb strings.Builder
	sb.WriteString("## Security Model\n\n")

	sm := d.SecurityModel

	// Overview
	if sm.Overview != "" {
		sb.WriteString("### Overview\n\n")
		sb.WriteString(sm.Overview + "\n\n")
	}

	// Threat Model
	sb.WriteString("### Threat Model\n\n")

	if len(sm.ThreatModel.Assets) > 0 {
		sb.WriteString("**Assets:**\n\n")
		for _, asset := range sm.ThreatModel.Assets {
			sb.WriteString(fmt.Sprintf("- %s\n", asset))
		}
		sb.WriteString("\n")
	}

	if len(sm.ThreatModel.ThreatActors) > 0 {
		sb.WriteString("**Threat Actors:**\n\n")
		for _, actor := range sm.ThreatModel.ThreatActors {
			sb.WriteString(fmt.Sprintf("- %s\n", actor))
		}
		sb.WriteString("\n")
	}

	if len(sm.ThreatModel.TrustBoundaries) > 0 {
		sb.WriteString("**Trust Boundaries:**\n\n")
		for _, boundary := range sm.ThreatModel.TrustBoundaries {
			sb.WriteString(fmt.Sprintf("- %s\n", boundary))
		}
		sb.WriteString("\n")
	}

	if len(sm.ThreatModel.KeyThreats) > 0 {
		sb.WriteString("**Key Threats:**\n\n")
		sb.WriteString("| ID | Category | Threat | Severity | Mitigation | Status |\n")
		sb.WriteString("|----|----------|--------|----------|------------|--------|\n")
		for _, t := range sm.ThreatModel.KeyThreats {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				t.ID, t.Category, t.Threat, t.Severity, t.Mitigation, t.Status))
		}
		sb.WriteString("\n")
	}

	// Access Control
	sb.WriteString("### Access Control\n\n")
	sb.WriteString(fmt.Sprintf("**Model:** %s\n\n", sm.AccessControl.Model))

	if sm.AccessControl.Description != "" {
		sb.WriteString(sm.AccessControl.Description + "\n\n")
	}

	if sm.AccessControl.Policies != "" {
		sb.WriteString(fmt.Sprintf("**Policy Engine:** %s\n\n", sm.AccessControl.Policies))
	}

	if len(sm.AccessControl.Layers) > 0 {
		sb.WriteString("**Layers:**\n\n")
		sb.WriteString("| Layer | Controls | Description |\n")
		sb.WriteString("|-------|----------|-------------|\n")
		for _, layer := range sm.AccessControl.Layers {
			controls := strings.Join(layer.Controls, ", ")
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				layer.Layer, controls, layer.Description))
		}
		sb.WriteString("\n")
	}

	if len(sm.AccessControl.Roles) > 0 {
		sb.WriteString("**Roles:**\n\n")
		sb.WriteString("| Role | Description | Permissions | Scope |\n")
		sb.WriteString("|------|-------------|-------------|-------|\n")
		for _, role := range sm.AccessControl.Roles {
			perms := strings.Join(role.Permissions, ", ")
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				role.Role, role.Description, perms, role.Scope))
		}
		sb.WriteString("\n")
	}

	// Encryption
	sb.WriteString("### Encryption\n\n")

	sb.WriteString("**At Rest:**\n\n")
	sb.WriteString(fmt.Sprintf("- **Method:** %s\n", sm.Encryption.AtRest.Method))
	sb.WriteString(fmt.Sprintf("- **Key Management:** %s\n", sm.Encryption.AtRest.KeyManagement))
	if sm.Encryption.AtRest.Provider != "" {
		sb.WriteString(fmt.Sprintf("- **Provider:** %s\n", sm.Encryption.AtRest.Provider))
	}
	if sm.Encryption.AtRest.Rotation != "" {
		sb.WriteString(fmt.Sprintf("- **Rotation:** %s\n", sm.Encryption.AtRest.Rotation))
	}
	sb.WriteString("\n")

	sb.WriteString("**In Transit:**\n\n")
	sb.WriteString(fmt.Sprintf("- **Method:** %s\n", sm.Encryption.InTransit.Method))
	sb.WriteString(fmt.Sprintf("- **Key Management:** %s\n", sm.Encryption.InTransit.KeyManagement))
	if sm.Encryption.InTransit.Provider != "" {
		sb.WriteString(fmt.Sprintf("- **Provider:** %s\n", sm.Encryption.InTransit.Provider))
	}
	sb.WriteString("\n")

	if sm.Encryption.FieldLevel != nil {
		sb.WriteString("**Field Level:**\n\n")
		sb.WriteString(fmt.Sprintf("- **Method:** %s\n", sm.Encryption.FieldLevel.Method))
		sb.WriteString(fmt.Sprintf("- **Key Management:** %s\n", sm.Encryption.FieldLevel.KeyManagement))
		sb.WriteString("\n")
	}

	// Audit Logging
	sb.WriteString("### Audit Logging\n\n")
	sb.WriteString(fmt.Sprintf("**Scope:** %s\n\n", sm.AuditLogging.Scope))

	if len(sm.AuditLogging.Events) > 0 {
		sb.WriteString("**Events:**\n\n")
		for _, event := range sm.AuditLogging.Events {
			sb.WriteString(fmt.Sprintf("- %s\n", event))
		}
		sb.WriteString("\n")
	}

	if sm.AuditLogging.Format != "" {
		sb.WriteString(fmt.Sprintf("**Format:** %s\n\n", sm.AuditLogging.Format))
	}
	sb.WriteString(fmt.Sprintf("**Retention:** %s\n\n", sm.AuditLogging.Retention))

	if sm.AuditLogging.Immutability != "" {
		sb.WriteString(fmt.Sprintf("**Immutability:** %s\n\n", sm.AuditLogging.Immutability))
	}
	if sm.AuditLogging.Destination != "" {
		sb.WriteString(fmt.Sprintf("**Destination:** %s\n\n", sm.AuditLogging.Destination))
	}

	// Compliance Controls
	if len(sm.ComplianceControls) > 0 {
		sb.WriteString("### Compliance Controls\n\n")

		// Sort framework names for consistent ordering
		var frameworks []string
		for framework := range sm.ComplianceControls {
			frameworks = append(frameworks, framework)
		}
		sort.Strings(frameworks)

		for _, framework := range frameworks {
			controls := sm.ComplianceControls[framework]
			sb.WriteString(fmt.Sprintf("**%s:**\n\n", framework))
			for _, ctrl := range controls {
				sb.WriteString(fmt.Sprintf("- %s\n", ctrl))
			}
			sb.WriteString("\n")
		}
	}

	// Data Classification
	if len(sm.DataClassification) > 0 {
		sb.WriteString("### Data Classification\n\n")
		sb.WriteString("| Level | Description | Handling | Examples |\n")
		sb.WriteString("|-------|-------------|----------|----------|\n")
		for _, dc := range sm.DataClassification {
			examples := strings.Join(dc.Examples, ", ")
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				dc.Level, dc.Description, dc.Handling, examples))
		}
		sb.WriteString("\n")
	}

	// Appendix References
	if len(sm.AppendixRefs) > 0 {
		sb.WriteString("### Related Appendices\n\n")
		for _, ref := range sm.AppendixRefs {
			sb.WriteString(fmt.Sprintf("- [%s](#appendix-%s)\n", ref, toSlug(ref)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	return sb.String()
}

func (d *Document) generateAppendices() string {
	var sb strings.Builder
	sb.WriteString("## Appendices\n\n")

	for i, appendix := range d.Appendices {
		// Appendix header with anchor
		sb.WriteString(fmt.Sprintf("### Appendix %s: %s {#appendix-%s}\n\n",
			indexToLetter(i), appendix.Title, toSlug(appendix.ID)))

		if appendix.Description != "" {
			sb.WriteString(fmt.Sprintf("*%s*\n\n", appendix.Description))
		}

		// Show tags if present
		if len(appendix.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(appendix.Tags, ", ")))
		}

		// Schema indicator
		if appendix.Schema != "" && appendix.Schema != AppendixSchemaCustom {
			sb.WriteString(fmt.Sprintf("**Schema:** %s\n\n", appendix.Schema))
		}

		// Content string (rendered first)
		if appendix.ContentString != "" {
			sb.WriteString(appendix.ContentString + "\n\n")
		}

		// Content table (rendered after string)
		if appendix.ContentTable != nil && len(appendix.ContentTable.Rows) > 0 {
			// Headers
			if len(appendix.ContentTable.Headers) > 0 {
				sb.WriteString("| " + strings.Join(appendix.ContentTable.Headers, " | ") + " |\n")
				sb.WriteString("|" + strings.Repeat("--------|", len(appendix.ContentTable.Headers)) + "\n")
			}
			// Rows
			for _, row := range appendix.ContentTable.Rows {
				sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
			}
			sb.WriteString("\n")

			// Caption
			if appendix.ContentTable.Caption != "" {
				sb.WriteString(fmt.Sprintf("*%s*\n\n", appendix.ContentTable.Caption))
			}
		}

		// Referenced by
		if len(appendix.ReferencedBy) > 0 {
			sb.WriteString("**Referenced by:** ")
			sb.WriteString(strings.Join(appendix.ReferencedBy, ", "))
			sb.WriteString("\n\n")
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// indexToLetter converts a 0-based index to a letter (A, B, C, ..., Z, AA, AB, ...).
func indexToLetter(i int) string {
	if i < 26 {
		return string(rune('A' + i))
	}
	// For indices >= 26, use AA, AB, etc.
	return string(rune('A'+i/26-1)) + string(rune('A'+i%26))
}

func (d *Document) generateGlossary() string {
	var sb strings.Builder
	sb.WriteString("## Glossary\n\n")

	sb.WriteString("| Term | Definition |\n")
	sb.WriteString("|------|------------|\n")
	for _, term := range d.Glossary {
		var name string
		if term.Acronym != "" {
			name = fmt.Sprintf("**%s** (%s)", term.Term, term.Acronym)
		} else {
			name = fmt.Sprintf("**%s**", term.Term)
		}
		sb.WriteString(fmt.Sprintf("| %s | %s |\n", name, term.Definition))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateRelatedDocuments() string {
	var sb strings.Builder
	sb.WriteString("## Related Documents\n\n")

	sb.WriteString("| ID | Title | Type | Relationship | Description |\n")
	sb.WriteString("|----|-------|------|--------------|-------------|\n")
	for _, doc := range d.RelatedDocuments {
		title := doc.Title
		if doc.URL != "" {
			title = fmt.Sprintf("[%s](%s)", doc.Title, doc.URL)
		} else if doc.Path != "" {
			title = fmt.Sprintf("[%s](%s)", doc.Title, doc.Path)
		}
		description := doc.Description
		if description == "" {
			description = "-"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			doc.ID, title, doc.Type, doc.Relationship, description))
	}
	sb.WriteString("\n---\n\n")

	return sb.String()
}

func (d *Document) generateCustomSections() string {
	var sb strings.Builder

	for _, cs := range d.CustomSections {
		sb.WriteString(fmt.Sprintf("## %s\n\n", cs.Title))
		if cs.Description != "" {
			sb.WriteString(cs.Description + "\n\n")
		}
		// Content is interface{}, so we just note it exists
		sb.WriteString("*See JSON source for detailed content.*\n\n")
		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// truncate shortens a string to maxLen, adding "..." if truncated.
// If maxLen is 0 or negative, the string is returned unchanged (no truncation).
func truncate(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
