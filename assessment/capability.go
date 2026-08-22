package assessment

import (
	"fmt"

	core "github.com/grokify/prism-core"
)

// CapabilityReference and its relation constants are re-exported from
// prism-core — the PRISM ecosystem's shared primitives module every PRISM
// module (prism-capability, prism-roadmap, prism-maturity) already
// depends on — rather than redefined locally, so a capability reference
// means exactly the same thing everywhere in the family
// (~/go/src/github.com/grokify/prism composes all three together; using
// them together, not independently, is the design priority). Aliasing
// (not wrapping) keeps CapabilityReference fully type-identical to
// core.CapabilityRef: a []core.CapabilityRef produced by another PRISM
// module is directly usable here with no conversion.
//
// This is the Engineering-facing counterpart to portfolio dimensions
// (Kano/MIH): where those explain customer/market value, capability
// references show what organizational/platform capability stack the
// roadmap is building or leaning on — the basis for a "here's what this
// roadmap does for Engineering" capability-stack view (ideation doc), not a
// ranking input.
type CapabilityReference = core.CapabilityRef

const (
	CapabilityEnables   = core.CapabilityRelationEnables
	CapabilityImproves  = core.CapabilityRelationImproves
	CapabilityDependsOn = core.CapabilityRelationDependsOn
)

// ValidateCapabilityReference returns an error if required fields are
// missing or the relation is not recognized. A free function rather than a
// method — CapabilityReference is a prism-core type alias and so cannot
// carry package-local methods.
func ValidateCapabilityReference(c CapabilityReference) error {
	if c.CapabilityID == "" {
		return fmt.Errorf("capabilityId is required")
	}
	if !core.ValidCapabilityRelation(c.Relation) {
		return fmt.Errorf("invalid relation: %q", c.Relation)
	}
	return nil
}
