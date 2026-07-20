// Package integration provides utilities for integrating prism-roadmap with other
// PRISM modules like productcontext (ideas), prism-maturity (goals), and prism-capability.
package integration

import (
	"time"

	"github.com/grokify/prism-roadmap/signal"
)

// IdeaData represents the minimal idea data needed to create a MarketSignal.
// This interface allows importing idea data from productcontext without creating
// a direct dependency on that package.
type IdeaData struct {
	ID            string   `json:"id"`
	Source        string   `json:"source"`
	Votes         int      `json:"votes"`
	CustomerCount int      `json:"customerCount"`
	ARRImpact     int64    `json:"arrImpact"` // in cents
	CustomerIDs   []string `json:"customerIds,omitempty"`
}

// IdeaToMarketSignal converts an IdeaData to a MarketSignal.
func IdeaToMarketSignal(idea IdeaData) *signal.MarketSignal {
	return signal.NewMarketSignalFromIdea(
		idea.Votes,
		idea.CustomerCount,
		idea.ARRImpact,
		idea.Source,
	)
}

// IdeasToMarketSignal aggregates multiple ideas into a single MarketSignal.
func IdeasToMarketSignal(ideas []IdeaData) *signal.MarketSignal {
	signals := make([]*signal.MarketSignal, 0, len(ideas))
	for _, idea := range ideas {
		signals = append(signals, IdeaToMarketSignal(idea))
	}
	result := signal.AggregateMarketSignals(signals...)
	result.UpdatedAt = time.Now()
	return result
}

// IdeaConversion tracks the conversion of an idea to a roadmap item.
type IdeaConversion struct {
	IdeaID      string    `json:"ideaId"`
	IdeaSource  string    `json:"ideaSource"`
	RMIID       string    `json:"rmiId"`
	ConvertedAt time.Time `json:"convertedAt"`
	ConvertedBy string    `json:"convertedBy,omitempty"`
}

// ConversionBatch represents a batch of idea-to-RMI conversions.
type ConversionBatch struct {
	Conversions []IdeaConversion `json:"conversions"`
	CreatedAt   time.Time        `json:"createdAt"`
}

// NewConversionBatch creates a new conversion batch.
func NewConversionBatch() *ConversionBatch {
	return &ConversionBatch{
		Conversions: []IdeaConversion{},
		CreatedAt:   time.Now(),
	}
}

// Add records a new conversion.
func (b *ConversionBatch) Add(ideaID, ideaSource, rmiID, convertedBy string) {
	b.Conversions = append(b.Conversions, IdeaConversion{
		IdeaID:      ideaID,
		IdeaSource:  ideaSource,
		RMIID:       rmiID,
		ConvertedAt: time.Now(),
		ConvertedBy: convertedBy,
	})
}

// GetByIdeaID finds conversion by idea ID.
func (b *ConversionBatch) GetByIdeaID(ideaID string) *IdeaConversion {
	for i := range b.Conversions {
		if b.Conversions[i].IdeaID == ideaID {
			return &b.Conversions[i]
		}
	}
	return nil
}

// GetByRMIID finds conversions by RMI ID.
func (b *ConversionBatch) GetByRMIID(rmiID string) []IdeaConversion {
	var result []IdeaConversion
	for _, c := range b.Conversions {
		if c.RMIID == rmiID {
			result = append(result, c)
		}
	}
	return result
}

// RefValidation holds the result of validating entity references.
type RefValidation struct {
	Valid   bool     `json:"valid"`
	Missing []string `json:"missing,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// RefValidator validates entity references.
type RefValidator struct {
	ValidIdeaIDs       map[string]bool
	ValidGoalIDs       map[string]bool
	ValidCapabilityIDs map[string]bool
}

// NewRefValidator creates a new reference validator.
func NewRefValidator() *RefValidator {
	return &RefValidator{
		ValidIdeaIDs:       make(map[string]bool),
		ValidGoalIDs:       make(map[string]bool),
		ValidCapabilityIDs: make(map[string]bool),
	}
}

// RegisterIdea adds an idea ID to the valid set.
func (v *RefValidator) RegisterIdea(id string) {
	v.ValidIdeaIDs[id] = true
}

// RegisterGoal adds a goal ID to the valid set.
func (v *RefValidator) RegisterGoal(id string) {
	v.ValidGoalIDs[id] = true
}

// RegisterCapability adds a capability ID to the valid set.
func (v *RefValidator) RegisterCapability(id string) {
	v.ValidCapabilityIDs[id] = true
}

// ValidateIdeaRefs validates a list of idea references.
func (v *RefValidator) ValidateIdeaRefs(refs []string) RefValidation {
	result := RefValidation{Valid: true}
	for _, ref := range refs {
		if !v.ValidIdeaIDs[ref] {
			result.Valid = false
			result.Missing = append(result.Missing, ref)
		}
	}
	return result
}

// ValidateGoalRefs validates a list of goal references.
func (v *RefValidator) ValidateGoalRefs(refs []string) RefValidation {
	result := RefValidation{Valid: true}
	for _, ref := range refs {
		if !v.ValidGoalIDs[ref] {
			result.Valid = false
			result.Missing = append(result.Missing, ref)
		}
	}
	return result
}

// ValidateCapabilityRefs validates a list of capability references.
func (v *RefValidator) ValidateCapabilityRefs(refs []string) RefValidation {
	result := RefValidation{Valid: true}
	for _, ref := range refs {
		if !v.ValidCapabilityIDs[ref] {
			result.Valid = false
			result.Missing = append(result.Missing, ref)
		}
	}
	return result
}
