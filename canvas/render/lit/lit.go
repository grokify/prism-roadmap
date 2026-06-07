// Package lit provides Lit web component data renderers for canvas types.
package lit

import (
	"encoding/json"
	"fmt"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

// LitData provides structured data for Lit web components.
type LitData struct {
	Component  string        `json:"component"`        // Web component name
	CanvasType string        `json:"canvasType"`       // Canvas type identifier
	Title      string        `json:"title,omitempty"`  // Canvas title
	Data       any           `json:"data"`             // Canvas-specific data
	Layout     *LayoutConfig `json:"layout,omitempty"` // Layout configuration
	Theme      *ThemeConfig  `json:"theme,omitempty"`  // Theme configuration
}

// LayoutConfig configures the visual layout.
type LayoutConfig struct {
	Rows      int    `json:"rows,omitempty"`
	Columns   int    `json:"columns,omitempty"`
	Direction string `json:"direction,omitempty"` // "down", "right", "left", "up"
	Grid      bool   `json:"grid,omitempty"`
}

// ThemeConfig configures the visual theme.
type ThemeConfig struct {
	Name       string            `json:"name,omitempty"`
	FontFamily string            `json:"fontFamily,omitempty"`
	Colors     map[string]string `json:"colors,omitempty"`
}

// LitRenderer renders canvas types to Lit component JSON data.
type LitRenderer struct{}

// NewLitRenderer creates a new Lit renderer.
func NewLitRenderer() *LitRenderer {
	return &LitRenderer{}
}

// Format returns the output format.
func (r *LitRenderer) Format() render.Format {
	return render.FormatLit
}

// FileExtension returns the file extension for Lit data files.
func (r *LitRenderer) FileExtension() string {
	return ".json"
}

// Supports returns true for all canvas types.
func (r *LitRenderer) Supports(canvasType canvas.CanvasType) bool {
	return true
}

// Render converts a canvas to Lit component JSON data.
func (r *LitRenderer) Render(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	data := &LitData{
		Component:  "canvas-" + string(c.Type),
		CanvasType: string(c.Type),
		Data:       c.GetInnerCanvas(),
		Layout:     layoutFromOptions(c.Type, opts),
		Theme:      themeFromOptions(opts),
	}

	// Get title from metadata
	if meta := c.GetMetadata(); meta != nil {
		data.Title = meta.Title
	}

	return json.MarshalIndent(data, "", "  ")
}

// layoutFromOptions creates a LayoutConfig from render Options.
func layoutFromOptions(canvasType canvas.CanvasType, opts *render.Options) *LayoutConfig {
	layout := &LayoutConfig{
		Direction: opts.Direction,
		Grid:      opts.GridLayout,
	}

	// Set default rows/columns based on canvas type
	switch canvasType {
	case canvas.CanvasTypeBMC:
		layout.Rows = 3
		layout.Columns = 5
		layout.Grid = true
	case canvas.CanvasTypeOpportunity:
		layout.Rows = 3
		layout.Columns = 3
		layout.Grid = true
	case canvas.CanvasTypeFeature:
		layout.Rows = 4
		layout.Columns = 2
		layout.Grid = true
	case canvas.CanvasTypeLeanUX:
		layout.Rows = 3
		layout.Columns = 3
		layout.Grid = true
	case canvas.CanvasTypeOST:
		layout.Direction = "down"
		layout.Grid = false
	}

	return layout
}

// themeFromOptions creates a ThemeConfig from render Options.
func themeFromOptions(opts *render.Options) *ThemeConfig {
	return &ThemeConfig{
		Name:       opts.Theme,
		FontFamily: opts.FontFamily,
		Colors:     opts.Colors,
	}
}

// BMCBlockData provides structured data for BMC web component blocks.
type BMCBlockData struct {
	KeyPartners       []BlockItem `json:"keyPartners"`
	KeyActivities     []BlockItem `json:"keyActivities"`
	KeyResources      []BlockItem `json:"keyResources"`
	ValuePropositions []BlockItem `json:"valuePropositions"`
	CustomerRelations []BlockItem `json:"customerRelations"`
	Channels          []BlockItem `json:"channels"`
	CustomerSegments  []BlockItem `json:"customerSegments"`
	CostStructure     []BlockItem `json:"costStructure"`
	RevenueStreams    []BlockItem `json:"revenueStreams"`
}

// OSTTreeData provides structured data for OST web component.
type OSTTreeData struct {
	Outcome       NodeData   `json:"outcome"`
	Opportunities []NodeData `json:"opportunities"`
}

// BlockItem represents an item in a canvas block.
type BlockItem struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

// NodeData represents a node in a tree structure.
type NodeData struct {
	ID       string     `json:"id"`
	Label    string     `json:"label"`
	Type     string     `json:"type,omitempty"` // outcome, opportunity, solution, experiment
	Status   string     `json:"status,omitempty"`
	Children []NodeData `json:"children,omitempty"`
}

// RenderBMCBlocks renders BMC-specific block data.
func (r *LitRenderer) RenderBMCBlocks(bmc *canvas.BusinessModelCanvas) (*BMCBlockData, error) {
	data := &BMCBlockData{
		KeyPartners:       make([]BlockItem, 0, len(bmc.KeyPartnerships)),
		KeyActivities:     make([]BlockItem, 0, len(bmc.KeyActivities)),
		KeyResources:      make([]BlockItem, 0, len(bmc.KeyResources)),
		ValuePropositions: make([]BlockItem, 0, len(bmc.ValuePropositions)),
		CustomerRelations: make([]BlockItem, 0, len(bmc.CustomerRelationships)),
		Channels:          make([]BlockItem, 0, len(bmc.Channels)),
		CustomerSegments:  make([]BlockItem, 0, len(bmc.CustomerSegments)),
		CostStructure:     make([]BlockItem, 0, len(bmc.CostStructure)),
		RevenueStreams:    make([]BlockItem, 0, len(bmc.RevenueStreams)),
	}

	for _, p := range bmc.KeyPartnerships {
		data.KeyPartners = append(data.KeyPartners, BlockItem{
			ID:          p.ID,
			Label:       p.Partner,
			Description: p.Description,
		})
	}

	for _, a := range bmc.KeyActivities {
		data.KeyActivities = append(data.KeyActivities, BlockItem{
			ID:          a.ID,
			Label:       a.Name,
			Description: a.Description,
		})
	}

	for _, res := range bmc.KeyResources {
		data.KeyResources = append(data.KeyResources, BlockItem{
			ID:          res.ID,
			Label:       res.Name,
			Description: res.Description,
		})
	}

	for _, vp := range bmc.ValuePropositions {
		data.ValuePropositions = append(data.ValuePropositions, BlockItem{
			ID:    vp.ID,
			Label: vp.Description,
		})
	}

	for _, cr := range bmc.CustomerRelationships {
		data.CustomerRelations = append(data.CustomerRelations, BlockItem{
			ID:          cr.ID,
			Label:       cr.Type,
			Description: cr.Description,
		})
	}

	for _, ch := range bmc.Channels {
		data.Channels = append(data.Channels, BlockItem{
			ID:          ch.ID,
			Label:       ch.Name,
			Description: ch.Description,
		})
	}

	for _, seg := range bmc.CustomerSegments {
		data.CustomerSegments = append(data.CustomerSegments, BlockItem{
			ID:          seg.ID,
			Label:       seg.Name,
			Description: seg.Description,
		})
	}

	for _, cost := range bmc.CostStructure {
		data.CostStructure = append(data.CostStructure, BlockItem{
			ID:    cost.ID,
			Label: cost.Description,
		})
	}

	for _, rev := range bmc.RevenueStreams {
		data.RevenueStreams = append(data.RevenueStreams, BlockItem{
			ID:    rev.ID,
			Label: rev.Description,
		})
	}

	return data, nil
}

// RenderOSTTree renders OST-specific tree data.
func (r *LitRenderer) RenderOSTTree(ost *canvas.OpportunitySolutionTree) (*OSTTreeData, error) {
	data := &OSTTreeData{
		Outcome: NodeData{
			ID:    ost.Outcome.ID,
			Label: ost.Outcome.Description,
			Type:  "outcome",
		},
		Opportunities: make([]NodeData, 0, len(ost.Outcome.Opportunities)),
	}

	for _, opp := range ost.Outcome.Opportunities {
		oppNode := NodeData{
			ID:       opp.ID,
			Label:    opp.Description,
			Type:     "opportunity",
			Children: make([]NodeData, 0, len(opp.Solutions)),
		}

		for _, sol := range opp.Solutions {
			solNode := NodeData{
				ID:       sol.ID,
				Label:    sol.Description,
				Type:     "solution",
				Status:   sol.Status,
				Children: make([]NodeData, 0, len(sol.Experiments)),
			}

			for _, exp := range sol.Experiments {
				expNode := NodeData{
					ID:     exp.ID,
					Label:  exp.Hypothesis,
					Type:   "experiment",
					Status: exp.Status,
				}
				solNode.Children = append(solNode.Children, expNode)
			}

			oppNode.Children = append(oppNode.Children, solNode)
		}

		data.Opportunities = append(data.Opportunities, oppNode)
	}

	// Add opportunities as children of outcome
	data.Outcome.Children = data.Opportunities

	return data, nil
}

// RenderWithBlockData renders a canvas with enhanced block data structure.
func (r *LitRenderer) RenderWithBlockData(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	var enhancedData any
	var err error

	switch c.Type {
	case canvas.CanvasTypeBMC:
		enhancedData, err = r.RenderBMCBlocks(c.BMC)
	case canvas.CanvasTypeOST:
		enhancedData, err = r.RenderOSTTree(c.OST)
	default:
		// Use default rendering for other types
		return r.Render(c, opts)
	}

	if err != nil {
		return nil, fmt.Errorf("rendering enhanced data: %w", err)
	}

	data := &LitData{
		Component:  "canvas-" + string(c.Type),
		CanvasType: string(c.Type),
		Data:       enhancedData,
		Layout:     layoutFromOptions(c.Type, opts),
		Theme:      themeFromOptions(opts),
	}

	if meta := c.GetMetadata(); meta != nil {
		data.Title = meta.Title
	}

	return json.MarshalIndent(data, "", "  ")
}

func init() {
	render.Register(NewLitRenderer())
}
