package canvas

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCanvasWrapper(t *testing.T) {
	tests := []struct {
		name       string
		canvas     *Canvas
		wantType   CanvasType
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "valid BMC",
			canvas: &Canvas{
				Type: CanvasTypeBMC,
				BMC: &BusinessModelCanvas{
					Metadata: Metadata{
						ID:    "bmc-1",
						Title: "Test BMC",
					},
				},
			},
			wantType: CanvasTypeBMC,
			wantErr:  false,
		},
		{
			name: "valid OST",
			canvas: &Canvas{
				Type: CanvasTypeOST,
				OST: &OpportunitySolutionTree{
					Metadata: Metadata{
						ID:    "ost-1",
						Title: "Test OST",
					},
				},
			},
			wantType: CanvasTypeOST,
			wantErr:  false,
		},
		{
			name:       "nil canvas",
			canvas:     nil,
			wantErr:    true,
			wantErrMsg: "canvas is nil",
		},
		{
			name: "empty type",
			canvas: &Canvas{
				BMC: &BusinessModelCanvas{},
			},
			wantErr:    true,
			wantErrMsg: "canvas type is required",
		},
		{
			name: "type mismatch",
			canvas: &Canvas{
				Type: CanvasTypeBMC,
				OST:  &OpportunitySolutionTree{},
			},
			wantErr:    true,
			wantErrMsg: "type is 'bmc' but BMC field is nil",
		},
		{
			name: "multiple canvases set",
			canvas: &Canvas{
				Type: CanvasTypeBMC,
				BMC:  &BusinessModelCanvas{},
				OST:  &OpportunitySolutionTree{},
			},
			wantErr:    true,
			wantErrMsg: "multiple inner canvases set; exactly one required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.canvas.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
					return
				}
				if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
					t.Errorf("Validate() error = %v, want %v", err, tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
				return
			}
			if tt.canvas.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", tt.canvas.Type, tt.wantType)
			}
		})
	}
}

func TestCanvasTypeGuards(t *testing.T) {
	bmc := NewBMC(&BusinessModelCanvas{})
	ost := NewOST(&OpportunitySolutionTree{})
	feature := NewFeature(&FeatureCanvas{})
	leanux := NewLeanUX(&LeanUXCanvas{})
	opp := NewOpportunity(&OpportunityCanvas{})

	if !bmc.IsBMC() {
		t.Error("IsBMC() should return true for BMC canvas")
	}
	if bmc.IsOST() {
		t.Error("IsOST() should return false for BMC canvas")
	}

	if !ost.IsOST() {
		t.Error("IsOST() should return true for OST canvas")
	}
	if ost.IsBMC() {
		t.Error("IsBMC() should return false for OST canvas")
	}

	if !feature.IsFeature() {
		t.Error("IsFeature() should return true for Feature canvas")
	}

	if !leanux.IsLeanUX() {
		t.Error("IsLeanUX() should return true for LeanUX canvas")
	}

	if !opp.IsOpportunity() {
		t.Error("IsOpportunity() should return true for Opportunity canvas")
	}
}

func TestShapeUpCanvasTypes(t *testing.T) {
	// Test ShapeUpPitch
	pitch := NewShapeUpPitchCanvas(&ShapeUpPitch{
		Metadata: Metadata{ID: "pitch-1", Title: "Test Pitch"},
		Problem: SUProblem{
			Statement: "Users can't do X",
			Evidence:  []string{"Support tickets show 50% asking for this"},
		},
		Appetite: SUAppetite{
			Weeks: 6,
			Size:  "big-batch",
		},
	})
	if !pitch.IsShapeUpPitch() {
		t.Error("IsShapeUpPitch() should return true for ShapeUpPitch canvas")
	}
	if pitch.IsOST() {
		t.Error("IsOST() should return false for ShapeUpPitch canvas")
	}
	if err := pitch.Validate(); err != nil {
		t.Errorf("ShapeUpPitch Validate() error: %v", err)
	}
	meta := pitch.GetMetadata()
	if meta == nil || meta.ID != "pitch-1" {
		t.Error("GetMetadata() should return pitch metadata")
	}

	// Test ShapeUpBet
	bet := NewShapeUpBetCanvas(&ShapeUpBet{
		Metadata: Metadata{ID: "bet-1", Title: "Test Bet"},
	})
	if !bet.IsShapeUpBet() {
		t.Error("IsShapeUpBet() should return true for ShapeUpBet canvas")
	}
	if err := bet.Validate(); err != nil {
		t.Errorf("ShapeUpBet Validate() error: %v", err)
	}

	// Test ShapeUpScope
	scope := NewShapeUpScopeCanvas(&ShapeUpScope{
		Metadata: Metadata{ID: "scope-1", Title: "Test Scope"},
	})
	if !scope.IsShapeUpScope() {
		t.Error("IsShapeUpScope() should return true for ShapeUpScope canvas")
	}
	if err := scope.Validate(); err != nil {
		t.Errorf("ShapeUpScope Validate() error: %v", err)
	}
}

func TestContinuousDiscoveryCanvasTypes(t *testing.T) {
	// Test DiscoverySnapshot
	snapshot := NewDiscoverySnapshotCanvas(&DiscoverySnapshot{
		Metadata: Metadata{ID: "snapshot-1", Title: "Week 12 Discovery"},
		Week:     "2024-W12",
		Interviews: []CDInterview{
			{ID: "int-1", ParticipantType: "Enterprise User", Stories: []CDStory{{ID: "s-1", Situation: "When logging in"}}},
		},
	})
	if !snapshot.IsDiscoverySnapshot() {
		t.Error("IsDiscoverySnapshot() should return true for DiscoverySnapshot canvas")
	}
	if snapshot.IsBMC() {
		t.Error("IsBMC() should return false for DiscoverySnapshot canvas")
	}
	if err := snapshot.Validate(); err != nil {
		t.Errorf("DiscoverySnapshot Validate() error: %v", err)
	}
	meta := snapshot.GetMetadata()
	if meta == nil || meta.ID != "snapshot-1" {
		t.Error("GetMetadata() should return snapshot metadata")
	}

	// Test AssumptionMap
	assumptionMap := NewAssumptionMapCanvas(&AssumptionMap{
		Metadata: Metadata{ID: "am-1", Title: "Test Assumption Map"},
		Desirability: []CDAssumption{
			{ID: "a-1", Description: "Users want this", Importance: "high", Confidence: "low"},
		},
	})
	if !assumptionMap.IsAssumptionMap() {
		t.Error("IsAssumptionMap() should return true for AssumptionMap canvas")
	}
	if err := assumptionMap.Validate(); err != nil {
		t.Errorf("AssumptionMap Validate() error: %v", err)
	}

	// Test ExperienceMap
	experienceMap := NewExperienceMapCanvas(&ExperienceMap{
		Metadata: Metadata{ID: "em-1", Title: "Test Experience Map"},
		Phases: []EMPhase{
			{ID: "p-1", Name: "Awareness", Description: "User discovers product"},
		},
	})
	if !experienceMap.IsExperienceMap() {
		t.Error("IsExperienceMap() should return true for ExperienceMap canvas")
	}
	if err := experienceMap.Validate(); err != nil {
		t.Errorf("ExperienceMap Validate() error: %v", err)
	}
}

func TestCanvasGetInnerCanvasNewTypes(t *testing.T) {
	// Test ShapeUpPitch
	pitch := &ShapeUpPitch{Metadata: Metadata{ID: "pitch-inner"}}
	canvasPitch := NewShapeUpPitchCanvas(pitch)
	inner := canvasPitch.GetInnerCanvas()
	if inner == nil {
		t.Fatal("GetInnerCanvas() returned nil for ShapeUpPitch")
	}
	innerPitch, ok := inner.(*ShapeUpPitch)
	if !ok {
		t.Fatalf("GetInnerCanvas() type = %T, want *ShapeUpPitch", inner)
	}
	if innerPitch.Metadata.ID != "pitch-inner" {
		t.Errorf("ID = %v, want pitch-inner", innerPitch.Metadata.ID)
	}

	// Test DiscoverySnapshot
	snapshot := &DiscoverySnapshot{Metadata: Metadata{ID: "snapshot-inner"}}
	canvasSnapshot := NewDiscoverySnapshotCanvas(snapshot)
	inner = canvasSnapshot.GetInnerCanvas()
	if inner == nil {
		t.Fatal("GetInnerCanvas() returned nil for DiscoverySnapshot")
	}
	innerSnapshot, ok := inner.(*DiscoverySnapshot)
	if !ok {
		t.Fatalf("GetInnerCanvas() type = %T, want *DiscoverySnapshot", inner)
	}
	if innerSnapshot.Metadata.ID != "snapshot-inner" {
		t.Errorf("ID = %v, want snapshot-inner", innerSnapshot.Metadata.ID)
	}
}

func TestCanvasJSON(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	original := &Canvas{
		Type: CanvasTypeBMC,
		BMC: &BusinessModelCanvas{
			Metadata: Metadata{
				ID:      "bmc-test",
				Title:   "Test Business Model",
				Version: VersionBMC1,
				Created: now,
			},
			CustomerSegments: []CustomerSegment{
				{ID: "cs-1", Name: "Enterprise"},
				{ID: "cs-2", Name: "SMB"},
			},
			ValuePropositions: []ValueProposition{
				{ID: "vp-1", Description: "Cost savings"},
			},
		},
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}

	// Unmarshal
	var parsed Canvas
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal() error: %v", err)
	}

	// Validate
	if parsed.Type != CanvasTypeBMC {
		t.Errorf("Type = %v, want %v", parsed.Type, CanvasTypeBMC)
	}
	if parsed.BMC == nil {
		t.Fatal("BMC is nil after unmarshal")
	}
	if parsed.BMC.Metadata.ID != "bmc-test" {
		t.Errorf("BMC.Metadata.ID = %v, want %v", parsed.BMC.Metadata.ID, "bmc-test")
	}
	if len(parsed.BMC.CustomerSegments) != 2 {
		t.Errorf("CustomerSegments length = %v, want 2", len(parsed.BMC.CustomerSegments))
	}
}

func TestGetMetadata(t *testing.T) {
	bmc := NewBMC(&BusinessModelCanvas{
		Metadata: Metadata{ID: "bmc-meta", Title: "BMC Title"},
	})

	meta := bmc.GetMetadata()
	if meta == nil {
		t.Fatal("GetMetadata() returned nil")
	}
	if meta.ID != "bmc-meta" {
		t.Errorf("ID = %v, want %v", meta.ID, "bmc-meta")
	}
	if meta.Title != "BMC Title" {
		t.Errorf("Title = %v, want %v", meta.Title, "BMC Title")
	}
}

func TestGetInnerCanvas(t *testing.T) {
	bmc := &BusinessModelCanvas{Metadata: Metadata{ID: "inner-test"}}
	canvas := NewBMC(bmc)

	inner := canvas.GetInnerCanvas()
	if inner == nil {
		t.Fatal("GetInnerCanvas() returned nil")
	}
	innerBMC, ok := inner.(*BusinessModelCanvas)
	if !ok {
		t.Fatalf("GetInnerCanvas() type = %T, want *BusinessModelCanvas", inner)
	}
	if innerBMC.Metadata.ID != "inner-test" {
		t.Errorf("ID = %v, want %v", innerBMC.Metadata.ID, "inner-test")
	}
}

func TestPRDReference(t *testing.T) {
	ref := &PRDReference{
		PRDID:          "prd-123",
		FeatureIDs:     []string{"feat-1", "feat-2"},
		RequirementIDs: []string{"req-1"},
	}

	if !ref.HasPRDLink() {
		t.Error("HasPRDLink() should return true when PRDID is set")
	}

	emptyRef := &PRDReference{}
	if emptyRef.HasPRDLink() {
		t.Error("HasPRDLink() should return false for empty reference")
	}

	var nilRef *PRDReference
	if nilRef.HasPRDLink() {
		t.Error("HasPRDLink() should return false for nil reference")
	}
}
