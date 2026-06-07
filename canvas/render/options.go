package render

// Options contains rendering options common to all renderers.
type Options struct {
	// Layout
	Direction  string // "right", "down", "left", "up"
	GridLayout bool   // For grid-based canvases (BMC, Feature)

	// Styling
	Theme      string            // "default", "corporate", "minimal"
	Colors     map[string]string // Block-specific colors
	FontFamily string

	// Content
	ShowMetadata bool // Include metadata section
	ShowPRDRefs  bool // Include PRD reference links
	MaxItems     int  // Truncate lists (0 = no limit)

	// D2-specific
	D2Theme  int    // D2 theme number (0-8)
	D2Layout string // "dagre", "elk", "tala"

	// SVG-specific
	Width  int
	Height int

	// Custom renderer-specific options
	Custom map[string]any
}

// DefaultOptions returns sensible default rendering options.
func DefaultOptions() *Options {
	return &Options{
		Direction:    "down",
		GridLayout:   false,
		Theme:        "default",
		Colors:       make(map[string]string),
		ShowMetadata: true,
		ShowPRDRefs:  false,
		MaxItems:     0,
		D2Theme:      0,
		D2Layout:     "dagre",
		Custom:       make(map[string]any),
	}
}

// BMCOptions returns options optimized for Business Model Canvas.
func BMCOptions() *Options {
	opts := DefaultOptions()
	opts.GridLayout = true
	opts.Direction = "right"
	opts.Colors = map[string]string{
		"keyPartners":       "#E3F2FD",
		"keyActivities":     "#FFF3E0",
		"keyResources":      "#FFF3E0",
		"valuePropositions": "#E8F5E9",
		"customerRelations": "#FCE4EC",
		"channels":          "#FCE4EC",
		"customerSegments":  "#F3E5F5",
		"costStructure":     "#FFEBEE",
		"revenueStreams":    "#E8F5E9",
	}
	return opts
}

// OSTOptions returns options optimized for Opportunity Solution Tree.
func OSTOptions() *Options {
	opts := DefaultOptions()
	opts.Direction = "down"
	opts.D2Layout = "dagre"
	opts.Colors = map[string]string{
		"outcome":     "#E8F5E9",
		"opportunity": "#E3F2FD",
		"solution":    "#FFF3E0",
		"experiment":  "#F3E5F5",
	}
	return opts
}

// OpportunityOptions returns options optimized for Opportunity Canvas (with arrows).
func OpportunityOptions() *Options {
	opts := DefaultOptions()
	opts.GridLayout = false
	opts.Direction = "down"
	return opts
}

// OpportunityGridOptions returns options for Opportunity Canvas in BMC-style grid.
func OpportunityGridOptions() *Options {
	opts := DefaultOptions()
	opts.GridLayout = true
	opts.Direction = "right"
	opts.Colors = map[string]string{
		"users":            "#E3F2FD", // Blue - Who
		"problems":         "#FFEBEE", // Red - Pains
		"solutionIdeas":    "#E8F5E9", // Green - Solutions
		"currentSolutions": "#FFF3E0", // Orange - Workarounds
		"userValue":        "#E8F5E9", // Green - Benefits
		"adoptionStrategy": "#F3E5F5", // Purple - Strategy
		"userMetrics":      "#E3F2FD", // Blue - User metrics
		"businessProblem":  "#FFEBEE", // Red - Business problem
		"businessMetrics":  "#E8F5E9", // Green - Business metrics
		"budget":           "#FFFDE7", // Yellow - Budget
	}
	return opts
}

// FeatureOptions returns options optimized for Feature Canvas.
func FeatureOptions() *Options {
	opts := DefaultOptions()
	opts.GridLayout = true
	opts.Direction = "right"
	opts.Colors = map[string]string{
		"ideaStatement": "#E8F5E9",
		"situations":    "#E3F2FD",
		"problems":      "#FFEBEE",
		"value":         "#E8F5E9",
		"capabilities":  "#FFF3E0",
		"restrictions":  "#FCE4EC",
		"limitations":   "#F3E5F5",
	}
	return opts
}

// LeanUXOptions returns options optimized for Lean UX Canvas.
func LeanUXOptions() *Options {
	opts := DefaultOptions()
	opts.GridLayout = true
	opts.Direction = "down"
	return opts
}

// WithDirection returns a copy of options with the specified direction.
func (o *Options) WithDirection(direction string) *Options {
	copy := *o
	copy.Direction = direction
	return &copy
}

// WithTheme returns a copy of options with the specified theme.
func (o *Options) WithTheme(theme string) *Options {
	copy := *o
	copy.Theme = theme
	return &copy
}

// WithD2Theme returns a copy of options with the specified D2 theme.
func (o *Options) WithD2Theme(theme int) *Options {
	copy := *o
	copy.D2Theme = theme
	return &copy
}

// WithD2Layout returns a copy of options with the specified D2 layout.
func (o *Options) WithD2Layout(layout string) *Options {
	copy := *o
	copy.D2Layout = layout
	return &copy
}

// WithColors returns a copy of options with additional colors.
func (o *Options) WithColors(colors map[string]string) *Options {
	copy := *o
	copy.Colors = make(map[string]string)
	for k, v := range o.Colors {
		copy.Colors[k] = v
	}
	for k, v := range colors {
		copy.Colors[k] = v
	}
	return &copy
}

// WithGridLayout returns a copy of options with grid layout enabled/disabled.
func (o *Options) WithGridLayout(enabled bool) *Options {
	copy := *o
	copy.GridLayout = enabled
	return &copy
}

// WithMetadata returns a copy of options with metadata display enabled/disabled.
func (o *Options) WithMetadata(show bool) *Options {
	copy := *o
	copy.ShowMetadata = show
	return &copy
}

// WithPRDRefs returns a copy of options with PRD reference display enabled/disabled.
func (o *Options) WithPRDRefs(show bool) *Options {
	copy := *o
	copy.ShowPRDRefs = show
	return &copy
}

// WithMaxItems returns a copy of options with the specified max items.
func (o *Options) WithMaxItems(max int) *Options {
	copy := *o
	copy.MaxItems = max
	return &copy
}
