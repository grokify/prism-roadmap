// Package canvas provides strategic planning canvas types including
// Business Model Canvas, Opportunity Canvas, Feature Canvas, Lean UX Canvas,
// and Opportunity Solution Tree.
package canvas

import (
	"time"

	"github.com/grokify/prism-roadmap/common"
)

// Version tracks canvas schema version for forward compatibility.
type Version string

// Canvas version constants.
const (
	VersionBMC1                   Version = "bmc/1.0"
	VersionOpportunity1           Version = "opportunity/1.0"
	VersionOpportunityAssessment1 Version = "opportunity-assessment/1.0" // Marty Cagan SVPG
	VersionOpportunitySpec1       Version = "opportunity-spec/1.0"       // Merged Patton + Cagan
	VersionFeature1               Version = "feature/1.0"
	VersionLeanUX2                Version = "leanux/2.0" // v2 per Gothelf
	VersionOST1                   Version = "ost/1.0"
	VersionShapeUp1               Version = "shapeup/1.0"   // Basecamp Shape Up (Ryan Singer)
	VersionDiscovery1             Version = "discovery/1.0" // Continuous Discovery (Teresa Torres)
)

// Metadata contains fields common to all canvas types.
type Metadata struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Version Version   `json:"version"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	Authors []Person  `json:"authors,omitempty"`
	Tags    []string  `json:"tags,omitempty"`

	// PRD integration
	PRDRef     string `json:"prdRef,omitempty"`     // PRD document ID
	FeatureRef string `json:"featureRef,omitempty"` // PRD feature/requirement ID
}

// Person reuses common.Person for canvas authors.
type Person = common.Person

// PRDReference links canvas to PRD elements.
type PRDReference struct {
	PRDID          string   `json:"prdId,omitempty"`
	FeatureIDs     []string `json:"featureIds,omitempty"`
	RequirementIDs []string `json:"requirementIds,omitempty"`
	UserStoryIDs   []string `json:"userStoryIds,omitempty"`
	PersonaIDs     []string `json:"personaIds,omitempty"`
}

// HasPRDLink returns true if the PRDReference has any links.
func (r *PRDReference) HasPRDLink() bool {
	if r == nil {
		return false
	}
	return r.PRDID != "" ||
		len(r.FeatureIDs) > 0 ||
		len(r.RequirementIDs) > 0 ||
		len(r.UserStoryIDs) > 0 ||
		len(r.PersonaIDs) > 0
}
