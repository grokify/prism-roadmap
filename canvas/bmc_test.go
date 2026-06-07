package canvas

import (
	"encoding/json"
	"testing"
)

func TestBusinessModelCanvas(t *testing.T) {
	bmc := NewBusinessModelCanvas("bmc-1", "SaaS Platform BMC")

	bmc.CustomerSegments = []CustomerSegment{
		{ID: "cs-1", Name: "Enterprise", Size: "500+ employees"},
		{ID: "cs-2", Name: "SMB", Size: "10-500 employees"},
	}
	bmc.ValuePropositions = []ValueProposition{
		{
			ID:            "vp-1",
			Description:   "Reduce operational costs by 30%",
			CustomerPains: []string{"Manual processes", "High error rates"},
			GainCreators:  []string{"Automation", "Real-time analytics"},
			SegmentRefs:   []string{"cs-1", "cs-2"},
		},
	}
	bmc.Channels = []Channel{
		{ID: "ch-1", Name: "Direct Sales", Type: "direct"},
		{ID: "ch-2", Name: "Partner Network", Type: "indirect"},
	}
	bmc.CustomerRelationships = []CustomerRelation{
		{ID: "cr-1", Type: "dedicated", Description: "Account managers"},
	}
	bmc.RevenueStreams = []RevenueStream{
		{ID: "rs-1", Description: "Subscription fees", Type: "subscription"},
	}
	bmc.KeyResources = []Resource{
		{ID: "kr-1", Name: "Platform", Type: "intellectual", Critical: true},
		{ID: "kr-2", Name: "Engineering Team", Type: "human"},
	}
	bmc.KeyActivities = []Activity{
		{ID: "ka-1", Name: "Platform Development", Category: "production"},
	}
	bmc.KeyPartnerships = []Partnership{
		{ID: "kp-1", Partner: "Cloud Provider", Type: "buyer-supplier"},
	}
	bmc.CostStructure = []Cost{
		{ID: "cost-1", Description: "Infrastructure", Type: "variable"},
		{ID: "cost-2", Description: "Personnel", Type: "fixed"},
	}

	// Test metadata
	if bmc.Metadata.ID != "bmc-1" {
		t.Errorf("ID = %v, want %v", bmc.Metadata.ID, "bmc-1")
	}
	if bmc.Metadata.Version != VersionBMC1 {
		t.Errorf("Version = %v, want %v", bmc.Metadata.Version, VersionBMC1)
	}

	// Test AllSegmentIDs
	segIDs := bmc.AllSegmentIDs()
	if len(segIDs) != 2 {
		t.Errorf("AllSegmentIDs() length = %v, want 2", len(segIDs))
	}
	if segIDs[0] != "cs-1" {
		t.Errorf("First segment ID = %v, want cs-1", segIDs[0])
	}

	// Test AllValuePropositionIDs
	vpIDs := bmc.AllValuePropositionIDs()
	if len(vpIDs) != 1 {
		t.Errorf("AllValuePropositionIDs() length = %v, want 1", len(vpIDs))
	}

	// Test counts
	if len(bmc.Channels) != 2 {
		t.Errorf("Channels length = %v, want 2", len(bmc.Channels))
	}
	if len(bmc.KeyResources) != 2 {
		t.Errorf("KeyResources length = %v, want 2", len(bmc.KeyResources))
	}
	if len(bmc.CostStructure) != 2 {
		t.Errorf("CostStructure length = %v, want 2", len(bmc.CostStructure))
	}
}

func TestBMCJSON(t *testing.T) {
	original := &BusinessModelCanvas{
		Metadata: Metadata{
			ID:      "bmc-json-test",
			Title:   "JSON Test BMC",
			Version: VersionBMC1,
		},
		CustomerSegments: []CustomerSegment{
			{ID: "cs-1", Name: "Test Segment"},
		},
		ValuePropositions: []ValueProposition{
			{ID: "vp-1", Description: "Test Value"},
		},
		Channels: []Channel{
			{ID: "ch-1", Name: "Test Channel"},
		},
		CustomerRelationships: []CustomerRelation{
			{ID: "cr-1", Type: "self-service"},
		},
		RevenueStreams: []RevenueStream{
			{ID: "rs-1", Description: "Test Revenue"},
		},
		KeyResources: []Resource{
			{ID: "kr-1", Name: "Test Resource"},
		},
		KeyActivities: []Activity{
			{ID: "ka-1", Name: "Test Activity"},
		},
		KeyPartnerships: []Partnership{
			{ID: "kp-1", Partner: "Test Partner"},
		},
		CostStructure: []Cost{
			{ID: "cost-1", Description: "Test Cost"},
		},
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}

	// Unmarshal
	var parsed BusinessModelCanvas
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal() error: %v", err)
	}

	// Validate structure
	if parsed.Metadata.ID != "bmc-json-test" {
		t.Errorf("Metadata.ID = %v, want %v", parsed.Metadata.ID, "bmc-json-test")
	}
	if len(parsed.CustomerSegments) != 1 {
		t.Errorf("CustomerSegments length = %v, want 1", len(parsed.CustomerSegments))
	}
	if len(parsed.KeyPartnerships) != 1 {
		t.Errorf("KeyPartnerships length = %v, want 1", len(parsed.KeyPartnerships))
	}
}

func TestBMCPRDReference(t *testing.T) {
	bmc := NewBusinessModelCanvas("bmc-prd", "PRD Linked BMC")
	bmc.PRDRef = &PRDReference{
		PRDID:      "prd-456",
		PersonaIDs: []string{"persona-1", "persona-2"},
	}
	bmc.CustomerSegments = []CustomerSegment{
		{ID: "cs-1", Name: "Enterprise", PersonaRef: "persona-1"},
	}

	ref := bmc.GetPRDReference()
	if ref == nil {
		t.Fatal("GetPRDReference() returned nil")
	}
	if ref.PRDID != "prd-456" {
		t.Errorf("PRDID = %v, want %v", ref.PRDID, "prd-456")
	}
	if len(ref.PersonaIDs) != 2 {
		t.Errorf("PersonaIDs length = %v, want 2", len(ref.PersonaIDs))
	}

	if bmc.CustomerSegments[0].PersonaRef != "persona-1" {
		t.Errorf("PersonaRef = %v, want persona-1", bmc.CustomerSegments[0].PersonaRef)
	}
}

func TestValuePropositionReferences(t *testing.T) {
	vp := ValueProposition{
		ID:          "vp-1",
		Description: "Test value prop",
		SegmentRefs: []string{"cs-1", "cs-2"},
	}

	if len(vp.SegmentRefs) != 2 {
		t.Errorf("SegmentRefs length = %v, want 2", len(vp.SegmentRefs))
	}
}

func TestResourceTypes(t *testing.T) {
	resources := []Resource{
		{ID: "r1", Name: "Platform", Type: "intellectual", Critical: true},
		{ID: "r2", Name: "Team", Type: "human", Critical: false},
		{ID: "r3", Name: "Servers", Type: "physical", Critical: true},
		{ID: "r4", Name: "Capital", Type: "financial", Critical: false},
	}

	criticalCount := 0
	for _, r := range resources {
		if r.Critical {
			criticalCount++
		}
	}

	if criticalCount != 2 {
		t.Errorf("Critical resource count = %v, want 2", criticalCount)
	}
}
