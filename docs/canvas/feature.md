# Feature Canvas

The Feature Canvas follows Nikita Efimov's framework (CC BY-SA 4.0) for defining and validating features before implementation.

## Overview

The Feature Canvas helps teams articulate:

- **Why** the feature is needed (problems, situations)
- **What** value it provides (benefits, capabilities)
- **What constraints** exist (restrictions, limitations)

## Grid Layout

| Idea Statement (spans full width) ||
|-----------------------------------|---|
| **Situations** | **Restrictions** |
| **Problems** | **Limitations** |
| **Value** ||
| **Capabilities** ||

## Structure

```go
type FeatureCanvas struct {
    Metadata      Metadata
    IdeaStatement string         // Top banner

    // Problem Area
    Situations    []Situation    // When/where does problem occur?
    Problems      []Problem      // What problems to solve?
    Value         []ValueItem    // What value does it provide?
    Capabilities  []Capability   // What must the solution do?

    // Solution Area
    Restrictions  []string       // Hard constraints
    Limitations   []string       // Known limitations
}

type Capability struct {
    ID          string
    Description string
    Priority    string  // must, should, could, wont (MoSCoW)
}
```

## Example

```go
fc := canvas.NewFeatureCanvas("feat-export", "Data Export Feature")

fc.IdeaStatement = "Enable users to export their data in multiple formats"

fc.Situations = []canvas.Situation{
    {ID: "s1", Description: "User needs to share data with external tools"},
}

fc.Problems = []canvas.Problem{
    {ID: "p1", Description: "No way to get data out of the system"},
}

fc.Value = []canvas.ValueItem{
    {ID: "v1", Description: "Users can integrate with existing workflows"},
}

fc.Capabilities = []canvas.Capability{
    {ID: "c1", Description: "Export to CSV", Priority: "must"},
    {ID: "c2", Description: "Export to JSON", Priority: "must"},
    {ID: "c3", Description: "Scheduled exports", Priority: "should"},
}

fc.Restrictions = []string{"Must comply with GDPR data export requirements"}
fc.Limitations = []string{"Initial release limited to 10,000 rows per export"}

// Render
c := canvas.NewFeature(fc)
d2Output, _ := render.Render(c, render.FormatD2, render.FeatureOptions())
```

## MoSCoW Priorities

| Priority | Meaning |
|----------|---------|
| `must` | Essential for release |
| `should` | Important but not critical |
| `could` | Nice to have |
| `wont` | Not in this release |

## Examples

See `examples/canvas/feature/` for complete examples.
