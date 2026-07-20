// Package signal provides market signal types for aggregating customer demand.
package signal

import "time"

// MarketSignal aggregates customer/market demand signals for a roadmap item.
// Signals can come from ideas, feature requests, or direct customer feedback.
type MarketSignal struct {
	// Vote aggregation from linked ideas/requests
	TotalVotes int `json:"total_votes"`

	// Customer impact
	CustomerCount int    `json:"customer_count"`
	TotalARR      int64  `json:"total_arr"`          // in cents (USD by default)
	Currency      string `json:"currency,omitempty"` // e.g., "USD", "EUR"

	// Calculated score
	Score float64 `json:"score"`

	// Source tracking
	IdeaCount int      `json:"idea_count"`
	Sources   []string `json:"sources,omitempty"` // e.g., ["aha", "manual", "intercom"]

	// Customer details (optional)
	CustomerIDs []string `json:"customer_ids,omitempty"` // IDs of requesting customers

	// Idea tracking (optional)
	IdeaIDs []string `json:"idea_ids,omitempty"` // IDs of source ideas

	// Metadata
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Notes     string    `json:"notes,omitempty"`
}

// CalculateScore computes the market signal score.
// Formula: (votes × 0.1) + (customer_count × 1.0) + (arr / 10_000_000)
// ARR is divided by 10M cents ($100k) to normalize contribution.
func (m *MarketSignal) CalculateScore() float64 {
	m.Score = float64(m.TotalVotes)*0.1 +
		float64(m.CustomerCount)*1.0 +
		float64(m.TotalARR)/10_000_000 // cents to $100k units
	return m.Score
}

// ARRInDollars returns the ARR in dollars (from cents).
func (m *MarketSignal) ARRInDollars() float64 {
	return float64(m.TotalARR) / 100
}

// AvgARRPerCustomer returns the average ARR per requesting customer.
func (m *MarketSignal) AvgARRPerCustomer() float64 {
	if m.CustomerCount == 0 {
		return 0
	}
	return float64(m.TotalARR) / float64(m.CustomerCount)
}

// AvgVotesPerIdea returns the average votes per idea.
func (m *MarketSignal) AvgVotesPerIdea() float64 {
	if m.IdeaCount == 0 {
		return 0
	}
	return float64(m.TotalVotes) / float64(m.IdeaCount)
}

// HasSignal returns true if there's any signal data.
func (m *MarketSignal) HasSignal() bool {
	return m.TotalVotes > 0 || m.CustomerCount > 0 || m.TotalARR > 0 || m.IdeaCount > 0
}

// Merge combines another MarketSignal into this one.
func (m *MarketSignal) Merge(other *MarketSignal) {
	if other == nil {
		return
	}
	m.TotalVotes += other.TotalVotes
	m.CustomerCount += other.CustomerCount
	m.TotalARR += other.TotalARR
	m.IdeaCount += other.IdeaCount

	// Merge sources (deduplicated)
	sourceSet := make(map[string]bool)
	for _, s := range m.Sources {
		sourceSet[s] = true
	}
	for _, s := range other.Sources {
		if !sourceSet[s] {
			m.Sources = append(m.Sources, s)
			sourceSet[s] = true
		}
	}

	// Merge customer IDs (deduplicated)
	customerSet := make(map[string]bool)
	for _, c := range m.CustomerIDs {
		customerSet[c] = true
	}
	for _, c := range other.CustomerIDs {
		if !customerSet[c] {
			m.CustomerIDs = append(m.CustomerIDs, c)
			customerSet[c] = true
		}
	}

	// Merge idea IDs (deduplicated)
	ideaSet := make(map[string]bool)
	for _, id := range m.IdeaIDs {
		ideaSet[id] = true
	}
	for _, id := range other.IdeaIDs {
		if !ideaSet[id] {
			m.IdeaIDs = append(m.IdeaIDs, id)
			ideaSet[id] = true
		}
	}

	// Recalculate score
	m.CalculateScore()
}

// NewMarketSignal creates a new MarketSignal with default values.
func NewMarketSignal() *MarketSignal {
	return &MarketSignal{
		Currency:  "USD",
		Sources:   []string{},
		UpdatedAt: time.Now(),
	}
}

// NewMarketSignalFromIdea creates a MarketSignal from idea data.
func NewMarketSignalFromIdea(votes int, customerCount int, arrCents int64, source string) *MarketSignal {
	signal := &MarketSignal{
		TotalVotes:    votes,
		CustomerCount: customerCount,
		TotalARR:      arrCents,
		IdeaCount:     1,
		Sources:       []string{source},
		Currency:      "USD",
		UpdatedAt:     time.Now(),
	}
	signal.CalculateScore()
	return signal
}

// AggregateMarketSignals combines multiple signals into one.
func AggregateMarketSignals(signals ...*MarketSignal) *MarketSignal {
	result := NewMarketSignal()
	for _, s := range signals {
		result.Merge(s)
	}
	return result
}
