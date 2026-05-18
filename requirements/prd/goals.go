package prd

import (
	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/goals/v2mom"
)

// GoalsAlignment represents alignment with strategic goals.
// This allows a PRD to reference or embed goals from the structured-goals package.
type GoalsAlignment struct {
	// V2MOMRef is a reference to an external V2MOM document.
	V2MOMRef *GoalReference `json:"v2mom_ref,omitempty"`

	// V2MOM is an embedded V2MOM document.
	V2MOM *v2mom.V2MOM `json:"v2mom,omitempty"`

	// OKRRef is a reference to an external OKR document.
	OKRRef *GoalReference `json:"okrRef,omitempty"`

	// OKR is an embedded OKR document.
	OKR *okr.OKRDocument `json:"okr,omitempty"`

	// AlignedObjectives maps PRD objectives to goal IDs.
	// Key is the PRD objective ID, value is the goal/method/objective ID.
	AlignedObjectives map[string]string `json:"alignedObjectives,omitempty"`
}

// GoalReference represents a reference to an external goals document.
type GoalReference struct {
	// ID is the unique identifier of the goals document.
	ID string `json:"id"`

	// Path is the file path to the goals document.
	Path string `json:"path,omitempty"`

	// URL is a URL to the goals document (e.g., Confluence, Notion).
	URL string `json:"url,omitempty"`

	// Version is the version of the goals document this PRD aligns with.
	Version string `json:"version,omitempty"`
}

// HasV2MOM returns true if V2MOM alignment is configured (ref or embedded).
func (g *GoalsAlignment) HasV2MOM() bool {
	if g == nil {
		return false
	}
	return g.V2MOM != nil || g.V2MOMRef != nil
}

// HasOKR returns true if OKR alignment is configured (ref or embedded).
func (g *GoalsAlignment) HasOKR() bool {
	if g == nil {
		return false
	}
	return g.OKR != nil || g.OKRRef != nil
}

// GetV2MOMMethodIDs returns the IDs of all V2MOM methods referenced by this PRD.
func (g *GoalsAlignment) GetV2MOMMethodIDs() []string {
	if g == nil || g.AlignedObjectives == nil {
		return nil
	}
	var ids []string
	for _, id := range g.AlignedObjectives {
		ids = append(ids, id)
	}
	return ids
}
