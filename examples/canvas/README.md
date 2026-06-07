# Canvas Examples

This directory contains examples for each canvas type with rendered outputs in multiple formats.

## Canvas Types

| Type | Description | Directory |
|------|-------------|-----------|
| BMC | Business Model Canvas (Osterwalder) | `bmc/` |
| OST | Opportunity Solution Tree (Torres) | `ost/` |
| Opportunity | Opportunity Canvas (Patton) | `opportunity/` |
| Feature | Feature Canvas (Efimov) | `feature/` |
| Lean UX | Lean UX Canvas (Gothelf) | `leanux/` |

### Opportunity Canvas Views

The Opportunity Canvas supports two view modes:

| View | File Suffix | Description |
|------|-------------|-------------|
| Grid | `_grid` | BMC-style 3x3+1 grid layout (no arrows) |
| Flow | `_flow` | Arrow-based flow showing relationships |

**Grid Layout (Jeff Patton's 9-block structure):**

| Users & Customers | Problems | Solution Ideas |
|-------------------|----------|----------------|
| Solutions Today | User Value | Adoption Strategy |
| User Metrics | Business Problem | Business Metrics |
| Budget (full width) |||

## File Formats

Each example includes:

| Format | Extension | Description |
|--------|-----------|-------------|
| JSON | `.json` | Source canvas data |
| D2 | `.d2` | D2 diagram language |
| SVG | `.svg` | Rendered SVG (from D2) |
| Mermaid | `.mmd` | Mermaid diagram syntax |
| Lit JSON | `.lit.json` | Data for Lit web components |
| HTML | `.html` | Interactive HTML viewer with Mermaid rendering |

## Viewing Examples

### HTML Viewer

Open any `*_example.html` file in a browser to see:

- Interactive Mermaid diagram
- Data view with structured information
- Raw Mermaid code
- JSON source data

### D2 Diagrams

View SVG files directly or regenerate with:

```bash
d2 bmc/bmc_example.d2 bmc/bmc_example.svg
```

### Mermaid Diagrams

View `.mmd` files in any Mermaid-compatible viewer or use the HTML files.

## Regenerating Examples

To regenerate all outputs:

```bash
cd examples/canvas
go run generate.go
```

To regenerate SVGs (requires d2 CLI):

```bash
d2 bmc/bmc_example.d2 bmc/bmc_example.svg
d2 ost/ost_example.d2 ost/ost_example.svg
d2 opportunity/opportunity_flow_example.d2 opportunity/opportunity_flow_example.svg
d2 opportunity/opportunity_grid_example.d2 opportunity/opportunity_grid_example.svg
d2 feature/feature_example.d2 feature/feature_example.svg
d2 leanux/leanux_example.d2 leanux/leanux_example.svg
```

## Using in Your Project

### Go

```go
import (
    "github.com/grokify/prism-roadmap/canvas"
    "github.com/grokify/prism-roadmap/canvas/render"
    "github.com/grokify/prism-roadmap/canvas/render/d2"
)

// Create an OST
ost := canvas.NewOpportunitySolutionTree("ost-1", "My OST")
ost.Outcome = canvas.OSTOutcome{
    ID:          "o1",
    Description: "Increase activation to 60%",
    // ...
}

// Wrap and render
c := canvas.NewOST(ost)
d2Output, err := render.Render(c, render.FormatD2, render.OSTOptions())
```

### Web Components (Lit)

The `.lit.json` files provide structured data for Lit web components:

```javascript
import canvasData from './ost_example.lit.json';

// Use with a custom Lit component
html`<canvas-ost .data=${canvasData}></canvas-ost>`;
```
