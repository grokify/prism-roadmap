package export

import "github.com/grokify/prism-roadmap/canvas"

// ExternalRef tracks a canvas's mapping to an external system.
type ExternalRef struct {
	Provider   string `json:"provider"`   // e.g., "aha", "miro"
	ExternalID string `json:"externalId"` // ID in the external system
	URL        string `json:"url"`        // URL to view in external system
}

// ComponentMapping maps a canvas block to an external system component.
type ComponentMapping struct {
	CanvasField string // Field name in the canvas (e.g., "Users", "Problems")
	ExternalID  string // Component ID in external system
	ExternalKey string // Component key/name in external system
}

// CanvasTypeMapping defines how a canvas type maps to external system concepts.
type CanvasTypeMapping struct {
	CanvasType   canvas.CanvasType
	ExternalKind string             // e.g., "Opportunity" for Aha strategic model kind
	Components   []ComponentMapping // Block-to-component mappings
}

// OpportunityCanvasBlocks defines the standard block names for Opportunity Canvas.
var OpportunityCanvasBlocks = []string{
	"Users & Customers",
	"Problems",
	"Solution Ideas",
	"Solutions Today",
	"User Value",
	"Adoption Strategy",
	"User Metrics",
	"Business Problem",
	"Business Metrics",
	"Budget",
}

// LeanUXCanvasBlocks defines the standard block names for Lean UX Canvas.
var LeanUXCanvasBlocks = []string{
	"Business Problem",
	"Business Outcomes",
	"Users",
	"Benefits",
	"Solutions",
	"Hypotheses",
	"Riskiest Assumption",
	"Smallest Experiment",
}

// BMCCanvasBlocks defines the standard block names for Business Model Canvas.
var BMCCanvasBlocks = []string{
	"Key Partners",
	"Key Activities",
	"Key Resources",
	"Value Propositions",
	"Customer Relationships",
	"Channels",
	"Customer Segments",
	"Cost Structure",
	"Revenue Streams",
}
