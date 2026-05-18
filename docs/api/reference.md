# API Reference

Complete API reference for PRISM Execution packages.

## Package Overview

```
github.com/grokify/prism-execution/
├── requirements/
│   ├── prd/    # Product Requirements Document
│   ├── mrd/    # Market Requirements Document
│   └── trd/    # Technical Requirements Document
├── goals/
│   ├── okr/    # Objectives and Key Results
│   └── v2mom/  # Vision, Values, Methods, Obstacles, Measures
└── roadmap/    # Roadmap with phases and deliverables
```

## PRD Package

```go
import "github.com/grokify/prism-execution/requirements/prd"
```

### Document Operations

#### New

```go
func New(id, title string, authors ...Person) *Document
```

Creates a new PRD with required fields initialized.

#### GenerateID

```go
func GenerateID() string
```

Generates a PRD ID based on current date: `PRD-YYYY-DDD`.

#### GenerateIDWithPrefix

```go
func GenerateIDWithPrefix(prefix string) string
```

Generates an ID with custom prefix: `PREFIX-YYYY-DDD`.

#### Load

```go
func Load(path string) (*Document, error)
```

Reads a PRD from a JSON file.

#### Save

```go
func Save(doc *Document, path string) error
```

Writes a PRD to a JSON file.

### Validation

#### Validate

```go
func Validate(doc *Document) *ValidationResult
```

Checks PRD for structural and content issues.

```go
type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationWarning
}
```

### Scoring

#### Score

```go
func Score(doc *Document) *ScoringResult
```

Evaluates PRD quality and returns scoring results.

#### DefaultWeights

```go
func DefaultWeights() []CategoryWeight
```

Returns standard category weights for scoring.

### Views

#### GeneratePMView

```go
func GeneratePMView(doc *Document) *PMView
```

Creates a PM-friendly view of the PRD.

#### RenderPMMarkdown

```go
func RenderPMMarkdown(view *PMView) string
```

Generates markdown output for PM view.

#### GenerateExecView

```go
func GenerateExecView(doc *Document, scores *ScoringResult) *ExecView
```

Creates an executive-friendly view.

#### RenderExecMarkdown

```go
func RenderExecMarkdown(view *ExecView) string
```

Generates markdown output for exec view.

#### GenerateSixPagerView

```go
func GenerateSixPagerView(doc *Document) *SixPagerView
```

Creates Amazon-style 6-pager view.

#### RenderSixPagerMarkdown

```go
func RenderSixPagerMarkdown(view *SixPagerView) string
```

Generates markdown output for 6-pager.

#### GeneratePRFAQView

```go
func GeneratePRFAQView(doc *Document) *PRFAQView
```

Creates PR/FAQ view (subset of 6-pager).

#### RenderPRFAQMarkdown

```go
func RenderPRFAQMarkdown(view *PRFAQView) string
```

Generates markdown output for PR/FAQ.

### Slide Renderers

```go
import "github.com/grokify/prism-execution/requirements/prd/render/marp"
```

#### NewPRDRenderer

```go
func NewPRDRenderer() *PRDRenderer
```

Creates a new PRD Marp slide renderer.

#### NewPRDGoalsRenderer

```go
func NewPRDGoalsRenderer() *PRDGoalsRenderer
```

Creates a new PRD+Goals Marp slide renderer with expanded V2MOM and OKR sections.

#### Renderer Interface

```go
type Renderer interface {
    Format() string
    FileExtension() string
    Render(doc *prd.Document, opts *render.Options) ([]byte, error)
}
```

#### Render Options

```go
type Options struct {
    Theme               string            // "default", "corporate", "minimal"
    IncludeGoals        bool              // Include goals alignment section
    IncludeRoadmap      bool              // Include roadmap slide
    IncludeRisks        bool              // Include risks slide
    IncludeRequirements bool              // Include requirements slide
    MaxPersonas         int               // Limit personas shown (0 = all)
    MaxRequirements     int               // Limit requirements shown (0 = all)
    CustomCSS           string            // Custom CSS for Marp
    Metadata            map[string]string // Additional metadata
}
```

#### DefaultOptions

```go
func DefaultOptions() *Options
```

Returns sensible default rendering options.

#### ExecutiveOptions

```go
func ExecutiveOptions() *Options
```

Returns options for executive-focused slides (fewer details).

### Persona Library

#### NewPersonaLibrary

```go
func NewPersonaLibrary() *PersonaLibrary
```

Creates a new empty persona library.

#### LoadPersonaLibrary

```go
func LoadPersonaLibrary(path string) (*PersonaLibrary, error)
```

Reads a persona library from JSON file.

#### PersonaLibrary Methods

```go
func (lib *PersonaLibrary) Add(persona LibraryPersona) error
func (lib *PersonaLibrary) Get(id string) (*LibraryPersona, bool)
func (lib *PersonaLibrary) GetByName(name string) (*LibraryPersona, bool)
func (lib *PersonaLibrary) List() []LibraryPersona
func (lib *PersonaLibrary) ListByTag(tag string) []LibraryPersona
func (lib *PersonaLibrary) Update(persona LibraryPersona) error
func (lib *PersonaLibrary) Remove(id string) error
func (lib *PersonaLibrary) Save(path string) error
func (lib *PersonaLibrary) ImportTo(doc *Document, ids ...string) error
func (lib *PersonaLibrary) ExportFrom(doc *Document, id string) error
func (lib *PersonaLibrary) ExportAllFrom(doc *Document) error
func (lib *PersonaLibrary) SyncFromLibrary(doc *Document) error
```

## Core Types

### Document

```go
type Document struct {
    Metadata         Metadata
    ExecutiveSummary ExecutiveSummary
    Objectives       Objectives
    Personas         []Persona
    UserStories      []UserStory
    Requirements     Requirements
    Roadmap          Roadmap
    Assumptions      *AssumptionsConstraints
    OutOfScope       []string
    TechArchitecture *TechnicalArchitecture
    UXRequirements   *UXRequirements
    Risks            []Risk
    Glossary         []GlossaryTerm
    CustomSections   []CustomSection
    Problem          *ProblemDefinition
    Market           *MarketDefinition
    Solution         *SolutionDefinition
    Decisions        *DecisionsDefinition
    Reviews          *ReviewsDefinition
    RevisionHistory  []RevisionRecord
    Goals            *GoalsAlignment
}
```

### Metadata

```go
type Metadata struct {
    ID        string
    Title     string
    Version   string
    Status    Status
    CreatedAt time.Time
    UpdatedAt time.Time
    Authors   []Person
    Reviewers []Person
    Approvers []Approver
    Tags      []string
}
```

### Person

```go
type Person struct {
    Name  string
    Email string
    Role  string
}
```

### Persona

```go
type Persona struct {
    ID                   string
    Name                 string
    Role                 string
    IsPrimary            bool
    Description          string
    Goals                []string
    PainPoints           []string
    TechnicalProficiency TechnicalProficiency
    Demographics         *Demographics
}
```

### Requirements

```go
type Requirements struct {
    Functional    []FunctionalRequirement
    NonFunctional []NonFunctionalRequirement
}

type FunctionalRequirement struct {
    ID                 string
    Title              string
    Description        string
    Priority           Priority
    MoSCoW             MoSCoW
    AcceptanceCriteria []string
    PersonaIDs         []string
}

type NonFunctionalRequirement struct {
    ID          string
    Category    NFRCategory
    Description string
    Target      string
    Priority    Priority
}
```

### GoalsAlignment

```go
type GoalsAlignment struct {
    V2MOMRef          *GoalReference
    V2MOM             *v2mom.V2MOM
    OKRRef            *GoalReference
    OKR               *okr.OKRDocument
    AlignedObjectives map[string]string
}

type GoalReference struct {
    ID      string
    Path    string
    URL     string
    Version string
}
```

## Constants

### Status

```go
const (
    StatusDraft      Status = "draft"
    StatusInReview   Status = "in_review"
    StatusApproved   Status = "approved"
    StatusDeprecated Status = "deprecated"
)
```

### Priority

```go
const (
    PriorityCritical Priority = "critical"
    PriorityHigh     Priority = "high"
    PriorityMedium   Priority = "medium"
    PriorityLow      Priority = "low"
)
```

### MoSCoW

```go
const (
    MoSCoWMust   MoSCoW = "must"
    MoSCoWShould MoSCoW = "should"
    MoSCoWCould  MoSCoW = "could"
    MoSCoWWont   MoSCoW = "wont"
)
```

### NFR Categories

```go
const (
    NFRPerformance     NFRCategory = "performance"
    NFRScalability     NFRCategory = "scalability"
    NFRReliability     NFRCategory = "reliability"
    NFRAvailability    NFRCategory = "availability"
    NFRSecurity        NFRCategory = "security"
    NFRMultiTenancy    NFRCategory = "multi_tenancy"
    NFRObservability   NFRCategory = "observability"
    NFRMaintainability NFRCategory = "maintainability"
    NFRUsability       NFRCategory = "usability"
    NFRCompatibility   NFRCategory = "compatibility"
    NFRCompliance      NFRCategory = "compliance"
)
```

### Review Decisions

```go
const (
    ReviewApprove     ReviewDecision = "approve"
    ReviewRevise      ReviewDecision = "revise"
    ReviewReject      ReviewDecision = "reject"
    ReviewHumanReview ReviewDecision = "human_review"
)
```

## MRD Package

```go
import "github.com/grokify/prism-execution/requirements/mrd"
```

See [MRD Documentation](../documents/mrd.md) for types and functions.

## TRD Package

```go
import "github.com/grokify/prism-execution/requirements/trd"
```

See [TRD Documentation](../documents/trd.md) for types and functions.

## Goals Package

```go
import (
    "github.com/grokify/prism-execution/goals/v2mom"
    "github.com/grokify/prism-execution/goals/okr"
)
```

See [Goals Integration](../goals/overview.md) for V2MOM and OKR types.

## Error Handling

All functions that can fail return an error:

```go
doc, err := prd.Load("file.prd.json")
if err != nil {
    log.Fatalf("Failed to load PRD: %v", err)
}

if err := prd.Save(doc, "output.prd.json"); err != nil {
    log.Fatalf("Failed to save PRD: %v", err)
}
```

## Next Steps

- [Quick Start](../getting-started/quickstart.md)
- [PRD Documentation](../documents/prd.md)
- [Examples](../examples/prd-examples.md)
