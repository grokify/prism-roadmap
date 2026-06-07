package canvas

import (
	"testing"
)

func TestNewShapeUpPitch(t *testing.T) {
	pitch := NewShapeUpPitch("pitch-001", "Improve Checkout Flow")
	if pitch.Metadata.ID != "pitch-001" {
		t.Errorf("expected ID pitch-001, got %s", pitch.Metadata.ID)
	}
	if pitch.Metadata.Title != "Improve Checkout Flow" {
		t.Errorf("expected title 'Improve Checkout Flow', got %s", pitch.Metadata.Title)
	}
	if pitch.Metadata.Version != VersionShapeUp1 {
		t.Errorf("expected version %s, got %s", VersionShapeUp1, pitch.Metadata.Version)
	}
}

func TestShapeUpPitchAppetite(t *testing.T) {
	pitch := NewShapeUpPitch("pitch-002", "Small Feature")
	pitch.Appetite = SUAppetite{
		Weeks: 2,
		Size:  "small-batch",
	}
	if !pitch.IsSmallBatch() {
		t.Error("expected IsSmallBatch() to be true for 2-week appetite")
	}
	if pitch.IsBigBatch() {
		t.Error("expected IsBigBatch() to be false for 2-week appetite")
	}

	pitch.Appetite = SUAppetite{
		Weeks: 6,
		Size:  "big-batch",
	}
	if pitch.IsSmallBatch() {
		t.Error("expected IsSmallBatch() to be false for 6-week appetite")
	}
	if !pitch.IsBigBatch() {
		t.Error("expected IsBigBatch() to be true for 6-week appetite")
	}
}

func TestShapeUpPitchBettingStatus(t *testing.T) {
	pitch := NewShapeUpPitch("pitch-003", "Test Feature")
	pitch.BettingStatus = "bet"
	if !pitch.IsBet() {
		t.Error("expected IsBet() to be true")
	}
	if pitch.IsDeclined() {
		t.Error("expected IsDeclined() to be false")
	}

	pitch.BettingStatus = "declined"
	if pitch.IsBet() {
		t.Error("expected IsBet() to be false")
	}
	if !pitch.IsDeclined() {
		t.Error("expected IsDeclined() to be true")
	}
}

func TestShapeUpPitchRabbitHoles(t *testing.T) {
	pitch := NewShapeUpPitch("pitch-004", "Complex Feature")
	if pitch.HasRabbitHoles() {
		t.Error("expected HasRabbitHoles() to be false when empty")
	}

	pitch.RabbitHoles = []SURabbitHole{
		{
			ID:           "rh-1",
			Description:  "Custom file format support",
			WhyDangerous: "Could take weeks to implement properly",
			Avoidance:    "Use standard formats only",
		},
	}
	if !pitch.HasRabbitHoles() {
		t.Error("expected HasRabbitHoles() to be true")
	}
}

func TestShapeUpScopeProgress(t *testing.T) {
	scope := NewShapeUpScope("scope-001", "Checkout Scopes")
	scope.Scopes = []SUScope{
		{ID: "s1", Name: "Cart UI", HillPosition: 30},      // Uphill
		{ID: "s2", Name: "Payment Flow", HillPosition: 70}, // Downhill
		{ID: "s3", Name: "Confirmation", HillPosition: 100, Status: "done"},
	}

	// Test overall progress
	progress := scope.OverallProgress()
	expected := (30 + 70 + 100) / 3 // 66
	if progress != expected {
		t.Errorf("expected progress %d, got %d", expected, progress)
	}

	// Test uphill scopes
	uphill := scope.UphillScopes()
	if len(uphill) != 1 {
		t.Errorf("expected 1 uphill scope, got %d", len(uphill))
	}
	if uphill[0].Name != "Cart UI" {
		t.Errorf("expected uphill scope 'Cart UI', got %s", uphill[0].Name)
	}

	// Test downhill scopes
	downhill := scope.DownhillScopes()
	if len(downhill) != 2 {
		t.Errorf("expected 2 downhill scopes, got %d", len(downhill))
	}

	// Test done scopes
	done := scope.DoneScopes()
	if len(done) != 1 {
		t.Errorf("expected 1 done scope, got %d", len(done))
	}
}

func TestNewShapeUpBet(t *testing.T) {
	bet := NewShapeUpBet("bet-001", "Checkout Improvement Bet")
	if bet.Metadata.ID != "bet-001" {
		t.Errorf("expected ID bet-001, got %s", bet.Metadata.ID)
	}
	if bet.Metadata.Version != VersionShapeUp1 {
		t.Errorf("expected version %s, got %s", VersionShapeUp1, bet.Metadata.Version)
	}
}
