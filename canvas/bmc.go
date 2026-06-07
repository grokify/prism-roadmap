package canvas

// BusinessModelCanvas represents Osterwalder's Business Model Canvas.
// The 9-block canvas covers key aspects of a business model.
type BusinessModelCanvas struct {
	Metadata              Metadata           `json:"metadata"`
	CustomerSegments      []CustomerSegment  `json:"customerSegments"`
	ValuePropositions     []ValueProposition `json:"valuePropositions"`
	Channels              []Channel          `json:"channels"`
	CustomerRelationships []CustomerRelation `json:"customerRelationships"`
	RevenueStreams        []RevenueStream    `json:"revenueStreams"`
	KeyResources          []Resource         `json:"keyResources"`
	KeyActivities         []Activity         `json:"keyActivities"`
	KeyPartnerships       []Partnership      `json:"keyPartnerships"`
	CostStructure         []Cost             `json:"costStructure"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// CustomerSegment represents a target customer group.
type CustomerSegment struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Size        string   `json:"size,omitempty"`       // Market size (e.g., "1M users", "$50B TAM")
	Needs       []string `json:"needs,omitempty"`      // Customer needs addressed
	PersonaRef  string   `json:"personaRef,omitempty"` // Link to PRD persona
}

// ValueProposition describes the value delivered to customers.
type ValueProposition struct {
	ID            string   `json:"id"`
	Description   string   `json:"description"`
	CustomerPains []string `json:"customerPains,omitempty"` // Pains addressed
	CustomerGains []string `json:"customerGains,omitempty"` // Gains created
	PainRelievers []string `json:"painRelievers,omitempty"` // How pains are relieved
	GainCreators  []string `json:"gainCreators,omitempty"`  // How gains are created
	SegmentRefs   []string `json:"segmentRefs,omitempty"`   // Customer segment IDs served
}

// Channel represents how a company reaches its customers.
type Channel struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type,omitempty"`  // direct, indirect, owned, partner
	Phase       string   `json:"phase,omitempty"` // awareness, evaluation, purchase, delivery, after-sales
	Description string   `json:"description,omitempty"`
	SegmentRefs []string `json:"segmentRefs,omitempty"` // Customer segments reached
}

// CustomerRelation describes the type of relationship with customers.
type CustomerRelation struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"` // personal, dedicated, self-service, automated, community, co-creation
	Description string   `json:"description,omitempty"`
	SegmentRefs []string `json:"segmentRefs,omitempty"` // Customer segments
}

// RevenueStream describes how a company generates revenue.
type RevenueStream struct {
	ID            string   `json:"id"`
	Description   string   `json:"description"`
	Type          string   `json:"type,omitempty"`          // asset sale, usage fee, subscription, licensing, etc.
	PricingModel  string   `json:"pricingModel,omitempty"`  // fixed, dynamic, negotiated
	Revenue       string   `json:"revenue,omitempty"`       // Estimated revenue or percentage
	SegmentRefs   []string `json:"segmentRefs,omitempty"`   // Customer segments
	ValuePropRefs []string `json:"valuePropRefs,omitempty"` // Value propositions monetized
}

// Resource describes key resources required.
type Resource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"` // physical, intellectual, human, financial
	Description string `json:"description,omitempty"`
	Critical    bool   `json:"critical,omitempty"` // Is this a critical resource?
}

// Activity describes key activities performed.
type Activity struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category,omitempty"` // production, problem-solving, platform/network
	Description string   `json:"description,omitempty"`
	Resources   []string `json:"resources,omitempty"` // Resource IDs used
}

// Partnership describes key partnerships.
type Partnership struct {
	ID          string   `json:"id"`
	Partner     string   `json:"partner"`
	Type        string   `json:"type,omitempty"`       // strategic alliance, coopetition, joint venture, buyer-supplier
	Motivation  string   `json:"motivation,omitempty"` // optimization, risk reduction, acquisition
	Description string   `json:"description,omitempty"`
	Resources   []string `json:"resources,omitempty"`  // Resources acquired from partner
	Activities  []string `json:"activities,omitempty"` // Activities performed by partner
}

// Cost describes cost structure elements.
type Cost struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type,omitempty"`      // fixed, variable
	Category    string `json:"category,omitempty"`  // cost-driven, value-driven
	Amount      string `json:"amount,omitempty"`    // Estimated cost
	Driver      string `json:"driver,omitempty"`    // What drives this cost (resource, activity, partnership)
	DriverRef   string `json:"driverRef,omitempty"` // ID of driving element
}

// NewBusinessModelCanvas creates a new BusinessModelCanvas with defaults.
func NewBusinessModelCanvas(id, title string) *BusinessModelCanvas {
	return &BusinessModelCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionBMC1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *BusinessModelCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// AllSegmentIDs returns all customer segment IDs.
func (c *BusinessModelCanvas) AllSegmentIDs() []string {
	ids := make([]string, 0, len(c.CustomerSegments))
	for _, seg := range c.CustomerSegments {
		ids = append(ids, seg.ID)
	}
	return ids
}

// AllValuePropositionIDs returns all value proposition IDs.
func (c *BusinessModelCanvas) AllValuePropositionIDs() []string {
	ids := make([]string, 0, len(c.ValuePropositions))
	for _, vp := range c.ValuePropositions {
		ids = append(ids, vp.ID)
	}
	return ids
}
