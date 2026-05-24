# Cascading Goals Framework TRD

**Technical Requirements Document**

**Version:** 1.0.0
**Date:** 2026-05-22
**Author:** PRISM Team
**Status:** Draft

## Overview

This document defines the technical architecture and implementation details for the cascading goals framework supporting both OKR and V2MOM methodologies.

## Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                             │
│  prism goals [cascade|align|import|suggest|plan]            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Goals Service Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Cascade   │  │   Align     │  │   Suggest/Plan      │  │
│  │   Service   │  │   Service   │  │   Service           │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Unified Framework Interface                  │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              goals.Framework                         │    │
│  │  - Type() string                                     │    │
│  │  - GoalItems() []GoalItem                           │    │
│  │  - ResultItems() []ResultItem                       │    │
│  │  - Validate(opts) []ValidationError                 │    │
│  │  - Cascade(parent) (Framework, error)               │    │
│  │  - AlignmentScore(parent) float64                   │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│   OKR Implementation    │     │  V2MOM Implementation   │
│   goals/okr/            │     │  goals/v2mom/           │
│                         │     │                         │
│  - OKRDocument          │     │  - V2MOM                │
│  - Objective            │     │  - Method               │
│  - KeyResult            │     │  - Measure              │
└─────────────────────────┘     └─────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  PRISM Maturity Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │    Goals    │  │    SLIs     │  │    State    │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### Package Structure

```
prism-roadmap/
├── goals/
│   ├── framework.go       # Unified Framework interface
│   ├── goal_item.go       # GoalItem, ResultItem types
│   ├── cascade.go         # Cascade service
│   ├── align.go           # Alignment service
│   ├── suggest.go         # Suggestion service
│   ├── okr/
│   │   ├── okr.go         # OKR types (extended)
│   │   ├── cascade.go     # OKR-specific cascading
│   │   ├── framework.go   # Framework implementation
│   │   └── validate.go    # OKR validation
│   └── v2mom/
│       ├── v2mom.go       # V2MOM types (existing)
│       ├── cascade.go     # V2MOM cascading (existing ParentID)
│       ├── framework.go   # Framework implementation
│       └── validate.go    # V2MOM validation
├── cli/
│   └── goals/
│       ├── cascade.go     # cascade command
│       ├── align.go       # align command
│       ├── import.go      # import command
│       ├── suggest.go     # suggest command
│       └── plan.go        # plan command
```

## Data Model

### OKR Extensions

```go
// File: goals/okr/okr.go

package okr

// Objective represents a qualitative goal with cascading support.
type Objective struct {
    ID          string       `json:"id"`
    Title       string       `json:"title"`
    Description string       `json:"description,omitempty"`
    Rationale   string       `json:"rationale,omitempty"`
    Category    string       `json:"category,omitempty"`
    Owner       string       `json:"owner,omitempty"`
    Status      string       `json:"status,omitempty"`
    Progress    float64      `json:"progress,omitempty"`
    KeyResults  []KeyResult  `json:"keyResults,omitempty"`
    Risks       []Risk       `json:"risks,omitempty"`

    // Cascading support (NEW)
    ParentObjectiveID string   `json:"parentObjectiveId,omitempty"`
    CompanyOKRIDs     []string `json:"companyOkrIds,omitempty"`
    ChildObjectiveIDs []string `json:"childObjectiveIds,omitempty"`

    // Alignment metadata (NEW)
    AlignmentScore    float64  `json:"alignmentScore,omitempty"`
    AlignmentStatus   string   `json:"alignmentStatus,omitempty"` // aligned, partial, unaligned
}

// KeyResult represents a measurable outcome with alignment tracking.
type KeyResult struct {
    ID         string   `json:"id"`
    Title      string   `json:"title"`
    Baseline   string   `json:"baseline,omitempty"`
    Target     string   `json:"target"`
    Current    string   `json:"current,omitempty"`
    Score      float64  `json:"score,omitempty"`
    Confidence string   `json:"confidence,omitempty"`
    Status     string   `json:"status,omitempty"`

    // Alignment tracking (NEW)
    SupportsObjectiveIDs []string `json:"supportsObjectiveIds,omitempty"`
    SLIIDs               []string `json:"sliIds,omitempty"` // Link to PRISM SLIs

    // Phase targets (existing)
    PhaseTargets []PhaseTarget `json:"phaseTargets,omitempty"`
}
```

### Unified Framework Interface

```go
// File: goals/framework.go

package goals

import "time"

// Framework abstracts over OKR and V2MOM goal-setting methodologies.
type Framework interface {
    // Type returns "okr" or "v2mom"
    Type() string

    // Metadata returns document-level information
    Metadata() FrameworkMetadata

    // GoalItems returns all goals (Objectives or Methods) in unified format
    GoalItems() []GoalItem

    // ResultItems returns all results (Key Results or Measures) in unified format
    ResultItems() []ResultItem

    // Validate checks document structure and relationships
    Validate(opts ValidationOptions) []ValidationError

    // Cascade generates a child framework from this parent
    Cascade(opts CascadeOptions) (Framework, error)

    // AlignmentScore calculates alignment with parent framework (0.0-1.0)
    AlignmentScore(parent Framework) float64

    // Import converts to PRISM goals format
    ToPRISMGoals() ([]PRISMGoal, error)
}

// FrameworkMetadata contains document-level information.
type FrameworkMetadata struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Owner     string    `json:"owner,omitempty"`
    Team      string    `json:"team,omitempty"`
    Period    string    `json:"period,omitempty"`
    Status    string    `json:"status,omitempty"`
    ParentID  string    `json:"parentId,omitempty"`
    CreatedAt time.Time `json:"createdAt,omitempty"`
    UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// GoalItem represents a goal in framework-agnostic format.
type GoalItem struct {
    ID          string       `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description,omitempty"`
    Owner       string       `json:"owner,omitempty"`
    Status      string       `json:"status,omitempty"`
    Progress    float64      `json:"progress,omitempty"`
    Priority    int          `json:"priority,omitempty"`
    ParentID    string       `json:"parentId,omitempty"`
    ChildIDs    []string     `json:"childIds,omitempty"`
    Results     []ResultItem `json:"results,omitempty"`

    // Source tracking
    SourceType string `json:"sourceType"` // "objective" or "method"
    SourceID   string `json:"sourceId"`
}

// ResultItem represents a measurable result in framework-agnostic format.
type ResultItem struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Baseline string  `json:"baseline,omitempty"`
    Target   string  `json:"target"`
    Current  string  `json:"current,omitempty"`
    Score    float64 `json:"score,omitempty"`
    Status   string  `json:"status,omitempty"`
    GoalID   string  `json:"goalId"`

    // PRISM integration
    SLIIDs []string `json:"sliIds,omitempty"`

    // Source tracking
    SourceType string `json:"sourceType"` // "keyResult" or "measure"
    SourceID   string `json:"sourceId"`
}

// CascadeOptions configures child framework generation.
type CascadeOptions struct {
    Team            string   `json:"team"`
    Owner           string   `json:"owner,omitempty"`
    Period          string   `json:"period,omitempty"`
    IncludeGoalIDs  []string `json:"includeGoalIds,omitempty"`  // Filter to specific goals
    ExcludeGoalIDs  []string `json:"excludeGoalIds,omitempty"`  // Exclude specific goals
    ResultsPerGoal  int      `json:"resultsPerGoal,omitempty"`  // Max results per goal (0=all)
    InheritResults  bool     `json:"inheritResults,omitempty"`  // Copy parent results as templates
}

// ValidationOptions configures validation behavior.
type ValidationOptions struct {
    RequireParentAlignment bool `json:"requireParentAlignment,omitempty"`
    MaxResultsPerGoal      int  `json:"maxResultsPerGoal,omitempty"`
    RequireSLILinks        bool `json:"requireSliLinks,omitempty"`
}

// ValidationError represents a validation issue.
type ValidationError struct {
    Code     string `json:"code"`
    Message  string `json:"message"`
    Path     string `json:"path,omitempty"`
    Severity string `json:"severity"` // error, warning, info
}
```

### Cascade Service

```go
// File: goals/cascade.go

package goals

// CascadeService generates child frameworks from parents.
type CascadeService struct {
    // Optional: template customization
    GoalNameTemplate   string // e.g., "[{Team}] {ParentGoalName}"
    ResultNameTemplate string // e.g., "Support: {ParentResultName}"
}

// Cascade creates a child framework from parent.
func (s *CascadeService) Cascade(parent Framework, opts CascadeOptions) (Framework, error) {
    switch parent.Type() {
    case "okr":
        return s.cascadeOKR(parent, opts)
    case "v2mom":
        return s.cascadeV2MOM(parent, opts)
    default:
        return nil, fmt.Errorf("unknown framework type: %s", parent.Type())
    }
}

// cascadeOKR creates child OKR from parent OKR.
func (s *CascadeService) cascadeOKR(parent Framework, opts CascadeOptions) (Framework, error) {
    parentOKR := parent.(*okr.OKRDocument)

    child := &okr.OKRDocument{
        Metadata: okr.Metadata{
            ID:       generateID(),
            Name:     fmt.Sprintf("%s OKRs", opts.Team),
            Team:     opts.Team,
            Owner:    opts.Owner,
            Period:   opts.Period,
            ParentID: parentOKR.Metadata.ID,
            Status:   "draft",
        },
    }

    // Generate child objectives from parent
    for _, parentObj := range parentOKR.Objectives {
        if !s.shouldInclude(parentObj.ID, opts) {
            continue
        }

        childObj := okr.Objective{
            ID:                generateID(),
            Title:             s.formatGoalName(parentObj.Title, opts),
            ParentObjectiveID: parentObj.ID,
            Status:            "draft",
        }

        // Optionally inherit key results as templates
        if opts.InheritResults {
            for _, parentKR := range parentObj.KeyResults {
                childKR := okr.KeyResult{
                    ID:                   generateID(),
                    Title:                s.formatResultName(parentKR.Title, opts),
                    Target:               parentKR.Target,
                    SupportsObjectiveIDs: []string{parentObj.ID},
                }
                childObj.KeyResults = append(childObj.KeyResults, childKR)
            }
        }

        child.Objectives = append(child.Objectives, childObj)
    }

    return child, nil
}
```

### Alignment Service

```go
// File: goals/align.go

package goals

// AlignmentService validates and scores goal alignment.
type AlignmentService struct{}

// AlignmentResult contains alignment analysis.
type AlignmentResult struct {
    OverallScore     float64           `json:"overallScore"`     // 0.0-1.0
    Status           string            `json:"status"`           // aligned, partial, unaligned
    GoalAlignments   []GoalAlignment   `json:"goalAlignments"`
    MissingSupport   []string          `json:"missingSupport"`   // Parent goals without child support
    UnalignedGoals   []string          `json:"unalignedGoals"`   // Child goals without parent link
    Recommendations  []string          `json:"recommendations"`
}

// GoalAlignment represents alignment between child and parent goal.
type GoalAlignment struct {
    ChildGoalID    string   `json:"childGoalId"`
    ChildGoalName  string   `json:"childGoalName"`
    ParentGoalID   string   `json:"parentGoalId,omitempty"`
    ParentGoalName string   `json:"parentGoalName,omitempty"`
    Score          float64  `json:"score"`
    Status         string   `json:"status"`
    Issues         []string `json:"issues,omitempty"`
}

// Align validates child framework alignment with parent.
func (s *AlignmentService) Align(child, parent Framework) (*AlignmentResult, error) {
    result := &AlignmentResult{
        Status: "aligned",
    }

    parentGoals := parent.GoalItems()
    childGoals := child.GoalItems()

    // Build parent goal lookup
    parentMap := make(map[string]GoalItem)
    for _, g := range parentGoals {
        parentMap[g.ID] = g
    }

    // Track which parent goals are supported
    supported := make(map[string]bool)

    // Analyze each child goal
    for _, childGoal := range childGoals {
        alignment := GoalAlignment{
            ChildGoalID:   childGoal.ID,
            ChildGoalName: childGoal.Name,
            Score:         1.0,
            Status:        "aligned",
        }

        if childGoal.ParentID == "" {
            alignment.Score = 0.0
            alignment.Status = "unaligned"
            alignment.Issues = append(alignment.Issues, "No parent goal linked")
            result.UnalignedGoals = append(result.UnalignedGoals, childGoal.ID)
        } else if parentGoal, ok := parentMap[childGoal.ParentID]; ok {
            alignment.ParentGoalID = parentGoal.ID
            alignment.ParentGoalName = parentGoal.Name
            supported[parentGoal.ID] = true

            // Check result alignment
            if !s.hasAlignedResults(childGoal, parentGoal) {
                alignment.Score = 0.5
                alignment.Status = "partial"
                alignment.Issues = append(alignment.Issues,
                    "Results don't explicitly support parent objectives")
            }
        } else {
            alignment.Score = 0.0
            alignment.Status = "unaligned"
            alignment.Issues = append(alignment.Issues,
                fmt.Sprintf("Parent goal %s not found", childGoal.ParentID))
        }

        result.GoalAlignments = append(result.GoalAlignments, alignment)
    }

    // Check for unsupported parent goals
    for _, parentGoal := range parentGoals {
        if !supported[parentGoal.ID] {
            result.MissingSupport = append(result.MissingSupport, parentGoal.ID)
            result.Recommendations = append(result.Recommendations,
                fmt.Sprintf("Consider adding goals to support '%s'", parentGoal.Name))
        }
    }

    // Calculate overall score
    result.OverallScore = s.calculateOverallScore(result)
    if result.OverallScore < 0.5 {
        result.Status = "unaligned"
    } else if result.OverallScore < 0.9 {
        result.Status = "partial"
    }

    return result, nil
}
```

### Suggestion Service

```go
// File: goals/suggest.go

package goals

import (
    prism "github.com/grokify/prism-maturity"
    "github.com/grokify/prism-maturity/maturity"
)

// SuggestionService generates goal suggestions from SLI gaps.
type SuggestionService struct{}

// Suggestion represents a suggested goal for maturity improvement.
type Suggestion struct {
    ID              string          `json:"id"`
    Title           string          `json:"title"`
    Description     string          `json:"description"`
    Priority        string          `json:"priority"` // P0, P1, P2, P3
    Domain          string          `json:"domain"`
    Layer           string          `json:"layer,omitempty"`
    CurrentLevel    float64         `json:"currentLevel"`
    TargetLevel     float64         `json:"targetLevel"`
    SLIGaps         []SLIGap        `json:"sliGaps"`
    SuggestedKRs    []SuggestedKR   `json:"suggestedKRs"`
    EstimatedEffort string          `json:"estimatedEffort"` // S, M, L, XL
    Impact          string          `json:"impact"`          // low, medium, high
}

// SLIGap represents a gap between current and target SLI value.
type SLIGap struct {
    SLIID        string  `json:"sliId"`
    SLIName      string  `json:"sliName"`
    CurrentValue float64 `json:"currentValue"`
    TargetValue  float64 `json:"targetValue"`
    GapPercent   float64 `json:"gapPercent"`
    Threshold    string  `json:"threshold"` // The threshold to meet
}

// SuggestedKR represents a suggested Key Result.
type SuggestedKR struct {
    Title    string `json:"title"`
    Baseline string `json:"baseline"`
    Target   string `json:"target"`
    SLIID    string `json:"sliId"`
}

// SuggestOptions configures suggestion generation.
type SuggestOptions struct {
    TargetLevel    int      `json:"targetLevel,omitempty"`    // M1-M5 (0 = next level)
    Domains        []string `json:"domains,omitempty"`        // Filter by domain
    Layers         []string `json:"layers,omitempty"`         // Filter by layer
    MaxSuggestions int      `json:"maxSuggestions,omitempty"` // Limit results
    MinGapPercent  float64  `json:"minGapPercent,omitempty"`  // Min gap to include
}

// Suggest generates goal suggestions from maturity gaps.
func (s *SuggestionService) Suggest(
    model *maturity.Spec,
    state *prism.PRISMDocument,
    opts SuggestOptions,
) ([]Suggestion, error) {
    var suggestions []Suggestion

    // Analyze each SLI
    for _, sli := range model.AllSLIs() {
        // Apply filters
        if !s.matchesFilters(sli, opts) {
            continue
        }

        // Get current value from state
        currentValue := s.getCurrentValue(state, sli.ID)
        currentLevel := s.calculateLevel(model, sli.ID, currentValue)

        // Determine target level
        targetLevel := opts.TargetLevel
        if targetLevel == 0 {
            targetLevel = int(currentLevel) + 1
            if targetLevel > 5 {
                continue // Already at max
            }
        }

        // Find gaps to reach target level
        gaps := s.findGaps(model, sli.ID, currentValue, targetLevel)
        if len(gaps) == 0 {
            continue // No gaps
        }

        // Check gap threshold
        maxGap := s.maxGapPercent(gaps)
        if opts.MinGapPercent > 0 && maxGap < opts.MinGapPercent {
            continue
        }

        // Generate suggestion
        suggestion := Suggestion{
            ID:              fmt.Sprintf("suggestion-%s", sli.ID),
            Title:           fmt.Sprintf("Improve %s to M%d", sli.Name, targetLevel),
            Description:     s.generateDescription(sli, gaps, targetLevel),
            Priority:        s.calculatePriority(gaps, sli),
            Domain:          sli.Category,
            Layer:           sli.Layer,
            CurrentLevel:    currentLevel,
            TargetLevel:     float64(targetLevel),
            SLIGaps:         gaps,
            SuggestedKRs:    s.generateKRs(sli, gaps),
            EstimatedEffort: s.estimateEffort(gaps),
            Impact:          s.assessImpact(sli, gaps),
        }

        suggestions = append(suggestions, suggestion)
    }

    // Sort by priority
    sort.Slice(suggestions, func(i, j int) bool {
        return suggestions[i].Priority < suggestions[j].Priority
    })

    // Apply limit
    if opts.MaxSuggestions > 0 && len(suggestions) > opts.MaxSuggestions {
        suggestions = suggestions[:opts.MaxSuggestions]
    }

    return suggestions, nil
}
```

## CLI Commands

### Goals Command Group

```go
// File: cli/goals/root.go

var goalsCmd = &cobra.Command{
    Use:   "goals",
    Short: "Goal framework management (OKR and V2MOM)",
    Long:  `Unified commands for managing OKR and V2MOM goal frameworks.`,
}

func init() {
    goalsCmd.AddCommand(cascadeCmd)
    goalsCmd.AddCommand(alignCmd)
    goalsCmd.AddCommand(importCmd)
    goalsCmd.AddCommand(suggestCmd)
    goalsCmd.AddCommand(planCmd)
}
```

### Cascade Command

```go
// File: cli/goals/cascade.go

var cascadeCmd = &cobra.Command{
    Use:   "cascade <parent-file>",
    Short: "Generate child goals from parent framework",
    Long: `Generate a child OKR or V2MOM document from a parent document.

The child document will have goals linked to parent goals via ParentID fields.

Examples:
  # Cascade company OKRs to team
  prism goals cascade company-okrs.json --team "Platform Team" -o team-okrs.json

  # Cascade V2MOM with specific goals only
  prism goals cascade company-v2mom.json --team "SRE" --include-goals "reliability,performance"

  # Inherit parent KRs as templates
  prism goals cascade parent.json --team "DevOps" --inherit-results`,
    RunE: runCascade,
}

var (
    cascadeTeam          string
    cascadeOwner         string
    cascadePeriod        string
    cascadeOutput        string
    cascadeIncludeGoals  []string
    cascadeExcludeGoals  []string
    cascadeInheritResult bool
)

func init() {
    cascadeCmd.Flags().StringVar(&cascadeTeam, "team", "", "Team name for child document (required)")
    cascadeCmd.Flags().StringVar(&cascadeOwner, "owner", "", "Owner for child document")
    cascadeCmd.Flags().StringVar(&cascadePeriod, "period", "", "Period (e.g., Q2 2026)")
    cascadeCmd.Flags().StringVarP(&cascadeOutput, "output", "o", "", "Output file (default: stdout)")
    cascadeCmd.Flags().StringSliceVar(&cascadeIncludeGoals, "include-goals", nil, "Goal IDs to include")
    cascadeCmd.Flags().StringSliceVar(&cascadeExcludeGoals, "exclude-goals", nil, "Goal IDs to exclude")
    cascadeCmd.Flags().BoolVar(&cascadeInheritResult, "inherit-results", false, "Copy parent results as templates")

    cascadeCmd.MarkFlagRequired("team")
}
```

### Suggest Command

```go
// File: cli/goals/suggest.go

var suggestCmd = &cobra.Command{
    Use:   "suggest",
    Short: "Suggest goals from SLI maturity gaps",
    Long: `Analyze current SLI state against maturity thresholds and suggest goals.

Examples:
  # Suggest goals to reach next maturity level
  prism goals suggest --model model.json --state state.json

  # Suggest goals to reach M4
  prism goals suggest --model model.json --state state.json --target-level 4

  # Filter by domain
  prism goals suggest --model model.json --state state.json --domain reliability

  # Output as OKR format
  prism goals suggest --model model.json --state state.json --format okr -o suggested-okrs.json`,
    RunE: runSuggest,
}

var (
    suggestModel        string
    suggestState        string
    suggestTargetLevel  int
    suggestDomains      []string
    suggestLayers       []string
    suggestMaxResults   int
    suggestMinGap       float64
    suggestFormat       string
    suggestOutput       string
)

func init() {
    suggestCmd.Flags().StringVar(&suggestModel, "model", "", "Maturity model file (required)")
    suggestCmd.Flags().StringVar(&suggestState, "state", "", "State file (required)")
    suggestCmd.Flags().IntVar(&suggestTargetLevel, "target-level", 0, "Target maturity level (1-5, 0=next)")
    suggestCmd.Flags().StringSliceVar(&suggestDomains, "domain", nil, "Filter by domain")
    suggestCmd.Flags().StringSliceVar(&suggestLayers, "layer", nil, "Filter by layer")
    suggestCmd.Flags().IntVar(&suggestMaxResults, "max", 10, "Maximum suggestions")
    suggestCmd.Flags().Float64Var(&suggestMinGap, "min-gap", 0, "Minimum gap percent to include")
    suggestCmd.Flags().StringVarP(&suggestFormat, "format", "f", "json", "Output format: json, okr, v2mom, markdown")
    suggestCmd.Flags().StringVarP(&suggestOutput, "output", "o", "", "Output file (default: stdout)")

    suggestCmd.MarkFlagRequired("model")
    suggestCmd.MarkFlagRequired("state")
}
```

## Testing Strategy

### Unit Tests

```go
// File: goals/cascade_test.go

func TestCascadeOKR(t *testing.T) {
    parent := &okr.OKRDocument{
        Metadata: okr.Metadata{ID: "parent-1", Name: "Company OKRs"},
        Objectives: []okr.Objective{
            {ID: "obj-1", Title: "Improve Reliability", KeyResults: []okr.KeyResult{
                {ID: "kr-1", Title: "99.9% Availability", Target: "99.9%"},
            }},
        },
    }

    svc := &CascadeService{}
    child, err := svc.Cascade(parent, CascadeOptions{
        Team:           "Platform Team",
        InheritResults: true,
    })

    require.NoError(t, err)
    require.Equal(t, "okr", child.Type())

    goals := child.GoalItems()
    require.Len(t, goals, 1)
    require.Equal(t, "obj-1", goals[0].ParentID)
    require.Len(t, goals[0].Results, 1)
}

func TestAlignmentScore(t *testing.T) {
    tests := []struct {
        name     string
        child    Framework
        parent   Framework
        expected float64
    }{
        {"fully aligned", fullyAlignedChild, parent, 1.0},
        {"partial alignment", partialChild, parent, 0.5},
        {"no alignment", unalignedChild, parent, 0.0},
    }

    svc := &AlignmentService{}
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result, err := svc.Align(tc.child, tc.parent)
            require.NoError(t, err)
            require.InDelta(t, tc.expected, result.OverallScore, 0.1)
        })
    }
}
```

### Integration Tests

```go
// File: goals/suggest_integration_test.go

func TestSuggestFromRealModel(t *testing.T) {
    model, err := maturity.ReadSpecFile("testdata/model.json")
    require.NoError(t, err)

    state, err := prism.ReadFile("testdata/state.json")
    require.NoError(t, err)

    svc := &SuggestionService{}
    suggestions, err := svc.Suggest(model, state, SuggestOptions{
        TargetLevel:    4,
        MaxSuggestions: 5,
    })

    require.NoError(t, err)
    require.LessOrEqual(t, len(suggestions), 5)

    for _, s := range suggestions {
        require.NotEmpty(t, s.Title)
        require.GreaterOrEqual(t, s.TargetLevel, s.CurrentLevel)
        require.NotEmpty(t, s.SLIGaps)
    }
}
```

## Performance Considerations

| Operation | Target | Approach |
|-----------|--------|----------|
| Cascade 50 goals | <100ms | Single pass, no I/O |
| Align 1000 goals | <1s | Index parent goals, single pass |
| Suggest 100 SLIs | <500ms | Parallel gap calculation |

## Security Considerations

- Input validation on all file paths
- JSON schema validation before processing
- No external network calls
- Sanitize user-provided strings in templates

## Migration Path

### Existing OKR Documents

Add optional fields without breaking:

```json
{
  "objectives": [{
    "id": "obj-1",
    "title": "Existing Objective",
    "parentObjectiveId": null,  // Optional, backward compatible
    "keyResults": [...]
  }]
}
```

### Existing V2MOM Documents

No changes required - `ParentID` already exists.

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/grokify/prism-maturity | latest | SLI/maturity access |
| github.com/spf13/cobra | v1.8+ | CLI framework |
| github.com/stretchr/testify | v1.8+ | Testing |

## Related Documents

- [PRD.md](PRD.md) - Product Requirements Document
- [ROADMAP.md](ROADMAP.md) - Implementation Roadmap
