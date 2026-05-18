# Persona Library

The Persona Library allows you to maintain reusable persona definitions across multiple PRDs.

## Purpose

- **Consistency**: Same persona definition across PRDs
- **Reusability**: Define once, use many times
- **Maintenance**: Update in one place
- **Governance**: Track persona usage

## Quick Start

```go
import "github.com/grokify/prism-execution/requirements/prd"

// Create a new library
lib := prd.NewPersonaLibrary()

// Add a persona
lib.Add(prd.LibraryPersona{
    Persona: prd.Persona{
        Name:       "Enterprise Admin",
        Role:       "IT Administrator",
        IsPrimary:  false,
        Goals:      []string{"Manage users", "Configure settings"},
        PainPoints: []string{"Complex setup", "Poor documentation"},
    },
    Tags:     []string{"enterprise", "admin", "B2B"},
    Category: "IT",
})

// Save library
lib.Save("personas.json")
```

## PersonaLibrary Structure

```go
type PersonaLibrary struct {
    Metadata LibraryMetadata  `json:"metadata"`
    Personas []LibraryPersona `json:"personas"`
}

type LibraryMetadata struct {
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    Version     string    `json:"version"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type LibraryPersona struct {
    Persona
    Tags      []string  `json:"tags,omitempty"`
    Category  string    `json:"category,omitempty"`
    UsedIn    []string  `json:"used_in,omitempty"`  // PRD IDs
    CreatedAt time.Time `json:"created_at,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
}
```

## Operations

### Add Persona

```go
// Add with auto-generated ID
lib.Add(prd.LibraryPersona{
    Persona: prd.Persona{
        Name: "Developer",
        Role: "Software Engineer",
    },
})

// Add with specific ID
lib.Add(prd.LibraryPersona{
    Persona: prd.Persona{
        ID:   "PER-DEV-001",
        Name: "Developer",
    },
})
```

### Get Persona

```go
// By ID
persona, found := lib.Get("PER-001")

// By name
persona, found := lib.GetByName("Enterprise Admin")
```

### List Personas

```go
// All personas
all := lib.List()

// By tag
admins := lib.ListByTag("admin")

// By category
itPersonas := lib.ListByCategory("IT")
```

### Update Persona

```go
persona.PainPoints = append(persona.PainPoints, "New pain point")
err := lib.Update(persona)
```

### Remove Persona

```go
err := lib.Remove("PER-001")
```

## Import/Export with PRDs

### Import to PRD

Copy personas from library to PRD:

```go
// Import specific personas
err := lib.ImportTo(doc, "PER-001", "PER-002")

// Import by tags
admins := lib.ListByTag("admin")
for _, p := range admins {
    lib.ImportTo(doc, p.ID)
}
```

### Export from PRD

Save PRD personas to library:

```go
// Export specific persona
err := lib.ExportFrom(doc, "PER-001")

// Export all personas
err := lib.ExportAllFrom(doc)
```

### Sync from Library

Update PRD personas from library (preserves local changes):

```go
err := lib.SyncFromLibrary(doc)
```

## Usage Tracking

The library tracks which PRDs use each persona:

```go
persona, _ := lib.Get("PER-001")
fmt.Printf("Used in: %v\n", persona.UsedIn)  // ["PRD-001", "PRD-002"]
```

## Example Library JSON

```json
{
  "metadata": {
    "name": "Product Personas",
    "description": "Shared personas for all product PRDs",
    "version": "1.0.0",
    "updated_at": "2025-01-22T10:00:00Z"
  },
  "personas": [
    {
      "id": "PER-001",
      "name": "Enterprise Admin",
      "role": "IT Administrator",
      "is_primary": false,
      "goals": ["Manage users", "Configure settings"],
      "pain_points": ["Complex setup", "Poor documentation"],
      "tags": ["enterprise", "admin", "B2B"],
      "category": "IT",
      "used_in": ["PRD-2025-001", "PRD-2025-003"],
      "created_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "PER-002",
      "name": "Power User",
      "role": "Account Manager",
      "is_primary": true,
      "goals": ["Quick data access", "Generate reports"],
      "pain_points": ["Slow performance", "Too many clicks"],
      "tags": ["power-user", "daily-user"],
      "category": "Sales",
      "used_in": ["PRD-2025-001"],
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

## File Conventions

| File | Purpose |
|------|---------|
| `personas.json` | Default library filename |
| `*.personas.json` | Domain-specific libraries |

```go
const prd.DefaultPersonaLibraryFilename = "personas.json"
```

## Best Practices

!!! tip "Organization"
    - Use tags for cross-cutting concerns (e.g., "enterprise", "SMB")
    - Use categories for organizational structure (e.g., "Sales", "IT")
    - Keep personas focused (one role per persona)

!!! tip "Maintenance"
    - Review unused personas quarterly
    - Update pain points as product evolves
    - Document persona changes in library version

!!! warning "Avoid"
    - Duplicating personas across libraries
    - Overly generic personas ("User")
    - Mixing user types in one persona

## Workflow Example

```go
// 1. Load or create library
lib, err := prd.LoadPersonaLibrary("personas.json")
if err != nil {
    lib = prd.NewPersonaLibrary()
}

// 2. Create new PRD
doc := prd.New("PRD-2025-001", "New Feature")

// 3. Import relevant personas
lib.ImportTo(doc, "PER-001", "PER-002")

// 4. Customize for this PRD
doc.Personas[0].PainPoints = append(
    doc.Personas[0].PainPoints,
    "Feature-specific pain point",
)

// 5. Save PRD
prd.Save(doc, "feature.prd.json")

// 6. Optionally export customizations back
lib.ExportFrom(doc, "PER-001")
lib.Save("personas.json")
```

## Next Steps

- [Scoring](scoring.md)
- [Completeness Check](completeness.md)
- [PRD Documentation](../documents/prd.md)
