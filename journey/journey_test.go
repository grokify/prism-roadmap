package journey

import (
	"encoding/json"
	"testing"
)

func TestJourneyRoadmapJSON(t *testing.T) {
	roadmap := JourneyRoadmap{
		ID:     "roadmap.plexusone.agent-platform",
		Type:   "journey_roadmap",
		Name:   "PlexusOne Agent Platform Evolution",
		Vision: "Enable portable, composable, omnichannel multi-agent systems.",
		Scope: &Scope{
			Type: "platform",
			ID:   "platform.plexusone",
		},
		TimeModel: &TimeModel{
			Type:       TimeModelQuarterly,
			FiscalYear: "FY2026",
			Periods: []Period{
				{ID: "now", Label: "Current State", IsCurrent: true},
				{ID: "2026-q3", Label: "Q3 2026"},
				{ID: "2026-q4", Label: "Q4 2026"},
				{ID: "2027-q1", Label: "Q1 2027"},
			},
		},
	}

	data, err := json.MarshalIndent(roadmap, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal roadmap: %v", err)
	}

	var decoded JourneyRoadmap
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal roadmap: %v", err)
	}

	if decoded.ID != roadmap.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.ID, roadmap.ID)
	}
	if decoded.TimeModel == nil {
		t.Fatal("TimeModel is nil after unmarshal")
	}
	if len(decoded.TimeModel.Periods) != 4 {
		t.Errorf("Period count: got %d, want 4", len(decoded.TimeModel.Periods))
	}
}

func TestCapabilityJourney(t *testing.T) {
	journey := CapabilityJourney{
		ID:           "capability-journey.multi-agent",
		CapabilityID: "capability.multi-agent-orchestration",
		Name:         "Multi-Agent Orchestration",
		CurrentState: &MaturityState{
			PeriodID:      "now",
			MaturityLevel: MaturityM1,
			Summary:       "Multi-agent behavior is implemented directly in individual applications.",
		},
		TargetStates: []TargetState{
			{
				PeriodID:      "2026-q3",
				MaturityLevel: MaturityM2,
				Summary:       "Agent teams are defined using Multi-Agent Spec.",
				Initiatives:   []string{"initiative.multi-agent-spec-omniagent"},
				Confidence:    0.85,
				Commitment:    CommitmentPlanned,
			},
			{
				PeriodID:      "2026-q4",
				MaturityLevel: MaturityM3,
				Summary:       "Multiple OmniAgent instances run as one coordinated team.",
				Initiatives:   []string{"initiative.agent-team-stats-reference"},
				Confidence:    0.70,
				Commitment:    CommitmentTargeted,
			},
		},
		DesiredEndState: &MaturityState{
			MaturityLevel: MaturityM4,
			Summary:       "Portable and observable multi-agent systems with interchangeable runtimes.",
		},
	}

	chain := BuildTransitionChain(journey)
	if chain == nil {
		t.Fatal("BuildTransitionChain returned nil")
	}

	if chain.StartMaturity != MaturityM1 {
		t.Errorf("StartMaturity: got %q, want %q", chain.StartMaturity, MaturityM1)
	}

	if len(chain.Transitions) != 2 {
		t.Errorf("Transition count: got %d, want 2", len(chain.Transitions))
	}

	if chain.Transitions[0].From.Maturity != MaturityM1 {
		t.Errorf("First transition from: got %q, want %q", chain.Transitions[0].From.Maturity, MaturityM1)
	}

	if chain.Transitions[0].To.Maturity != MaturityM2 {
		t.Errorf("First transition to: got %q, want %q", chain.Transitions[0].To.Maturity, MaturityM2)
	}
}

func TestTransitionMaturityDelta(t *testing.T) {
	tests := []struct {
		from, to string
		want     int
	}{
		{"M1", "M2", 1},
		{"M1", "M4", 3},
		{"M3", "M3", 0},
		{"M4", "M2", -2},
	}

	for _, tt := range tests {
		transition := CapabilityTransition{
			From: TransitionState{Maturity: tt.from},
			To:   TransitionState{Maturity: tt.to},
		}
		got := transition.MaturityDelta()
		if got != tt.want {
			t.Errorf("MaturityDelta(%s→%s): got %d, want %d", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestDependencyGraph(t *testing.T) {
	graph := DependencyGraph{
		Dependencies: []Dependency{
			{
				ID:     "dep-1",
				From:   EntityRef{Type: EntityTypeInitiative, ID: "init-sso"},
				To:     EntityRef{Type: EntityTypeTeam, ID: "team-platform"},
				Type:   DependencyResource,
				Status: DependencyStatusBlocked,
				Risk:   DependencyRiskHigh,
			},
			{
				ID:     "dep-2",
				From:   EntityRef{Type: EntityTypeCapability, ID: "cap-multi-agent"},
				To:     EntityRef{Type: EntityTypeCapability, ID: "cap-agent-runtime"},
				Type:   DependencyRequires,
				Status: DependencyStatusPending,
			},
		},
	}

	// Test GetBlockers
	blockers := graph.GetBlockers(EntityTypeInitiative, "init-sso")
	if len(blockers) != 1 {
		t.Errorf("GetBlockers: got %d, want 1", len(blockers))
	}

	// Test GetBlocking
	blocking := graph.GetBlocking(EntityTypeTeam, "team-platform")
	if len(blocking) != 1 {
		t.Errorf("GetBlocking: got %d, want 1", len(blocking))
	}

	// Test GetCriticalPath
	critical := graph.GetCriticalPath()
	if len(critical) != 1 {
		t.Errorf("GetCriticalPath: got %d, want 1", len(critical))
	}

	// Test GetResourceDependencies
	resources := graph.GetResourceDependencies()
	if len(resources) != 1 {
		t.Errorf("GetResourceDependencies: got %d, want 1", len(resources))
	}
}

func TestTeamHierarchy(t *testing.T) {
	teams := []Team{
		{ID: "org", Name: "Engineering Org", Level: TeamLevelOrganization},
		{ID: "platform-dept", Name: "Platform", Level: TeamLevelDepartment, ParentID: "org"},
		{ID: "team-infra", Name: "Infrastructure", Level: TeamLevelTeam, ParentID: "platform-dept",
			Capacity: &TeamCapacity{FTEs: 5, StoryPointsPerSprint: 40}},
		{ID: "team-runtime", Name: "Runtime", Level: TeamLevelTeam, ParentID: "platform-dept",
			Capacity: &TeamCapacity{FTEs: 4, StoryPointsPerSprint: 32}},
	}

	hierarchy := NewTeamHierarchy(teams)

	// Test GetChildren
	children := hierarchy.GetChildren("platform-dept")
	if len(children) != 2 {
		t.Errorf("GetChildren: got %d, want 2", len(children))
	}

	// Test GetParent
	parent := hierarchy.GetParent("team-infra")
	if parent == nil || parent.ID != "platform-dept" {
		t.Errorf("GetParent: got %v, want platform-dept", parent)
	}

	// Test GetDescendants
	descendants := hierarchy.GetDescendants("org")
	if len(descendants) != 3 {
		t.Errorf("GetDescendants: got %d, want 3", len(descendants))
	}

	// Test GetByLevel
	depts := hierarchy.GetByLevel(TeamLevelDepartment)
	if len(depts) != 1 {
		t.Errorf("GetByLevel(Department): got %d, want 1", len(depts))
	}

	// Test AggregateCapacity
	totalCap := hierarchy.AggregateCapacity("platform-dept")
	if totalCap == nil {
		t.Fatal("AggregateCapacity returned nil")
	}
	if totalCap.FTEs != 9 {
		t.Errorf("AggregateCapacity FTEs: got %f, want 9", totalCap.FTEs)
	}
}

func TestInitiativeSummarize(t *testing.T) {
	init := Initiative{
		ID:        "initiative.agent-team-stats",
		Name:      "Agent Team Stats Reference",
		Status:    InitiativeStatusInProgress,
		OwnerTeam: "team-platform",
		Periods:   []string{"2026-q3", "2026-q4"},
		Advances: []CapabilityAdvance{
			{CapabilityID: "cap-1", From: "M1", To: "M3"},
			{CapabilityID: "cap-2", From: "M2", To: "M3"},
		},
		ExpectedOutcomes: []string{"outcome-1", "outcome-2"},
	}

	summary := init.Summarize()

	if summary.CapabilityCount != 2 {
		t.Errorf("CapabilityCount: got %d, want 2", summary.CapabilityCount)
	}
	if summary.OutcomeCount != 2 {
		t.Errorf("OutcomeCount: got %d, want 2", summary.OutcomeCount)
	}
	if summary.MaturityGain != 3 { // (3-1) + (3-2) = 2 + 1 = 3
		t.Errorf("MaturityGain: got %d, want 3", summary.MaturityGain)
	}
}

func TestOutcomeJourneyImpact(t *testing.T) {
	oj := OutcomeJourney{
		ID:   "outcome.time-to-agent-system",
		Name: "Time to Build Multi-Agent System",
		CurrentState: &OutcomeState{
			Value: 20,
			Unit:  "developer-days",
		},
		TargetStates: []OutcomeTargetState{
			{PeriodID: "2026-q3", Value: 12, Unit: "developer-days"},
			{PeriodID: "2026-q4", Value: 5, Unit: "developer-days"},
		},
	}

	impact := oj.CalculateImpact("2026-q4")
	if impact == nil {
		t.Fatal("CalculateImpact returned nil")
	}

	if impact.FromValue != 20 {
		t.Errorf("FromValue: got %f, want 20", impact.FromValue)
	}
	if impact.ToValue != 5 {
		t.Errorf("ToValue: got %f, want 5", impact.ToValue)
	}
	if impact.AbsoluteChange != -15 {
		t.Errorf("AbsoluteChange: got %f, want -15", impact.AbsoluteChange)
	}
}

func TestBuildStoryboard(t *testing.T) {
	roadmap := &JourneyRoadmap{
		TimeModel: &TimeModel{
			Periods: []Period{
				{ID: "now", Label: "Current"},
				{ID: "q1", Label: "Q1 2026"},
			},
		},
		CapabilityJourneys: []CapabilityJourney{
			{
				Name: "Multi-Agent",
				CurrentState: &MaturityState{
					PeriodID:      "now",
					MaturityLevel: "M1",
				},
				TargetStates: []TargetState{
					{PeriodID: "q1", MaturityLevel: "M2", Confidence: 0.8},
				},
			},
		},
		Initiatives: []Initiative{
			{Name: "Init A", Periods: []string{"q1"}},
		},
		Narrative: &RoadmapNarrative{
			Journey: []JourneyChapter{
				{PeriodID: "q1", Headline: "Build Foundation"},
			},
		},
	}

	cards := BuildStoryboard(roadmap)
	if len(cards) != 2 {
		t.Fatalf("card count: got %d, want 2", len(cards))
	}

	q1Card := cards[1]
	if q1Card.Headline != "Build Foundation" {
		t.Errorf("Headline: got %q, want %q", q1Card.Headline, "Build Foundation")
	}
	if len(q1Card.MaturityChanges) != 1 {
		t.Errorf("MaturityChanges: got %d, want 1", len(q1Card.MaturityChanges))
	}
	if len(q1Card.MajorInitiatives) != 1 {
		t.Errorf("MajorInitiatives: got %d, want 1", len(q1Card.MajorInitiatives))
	}
}
