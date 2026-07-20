package signal

import "testing"

func TestMarketSignalCalculateScore(t *testing.T) {
	tests := []struct {
		name   string
		signal MarketSignal
		want   float64
	}{
		{
			name: "all components",
			signal: MarketSignal{
				TotalVotes:    100,     // 100 * 0.1 = 10
				CustomerCount: 5,       // 5 * 1.0 = 5
				TotalARR:      5000000, // 5M cents = $50k / 10M = 0.5
			},
			want: 15.5,
		},
		{
			name: "votes only",
			signal: MarketSignal{
				TotalVotes: 50,
			},
			want: 5.0,
		},
		{
			name: "customers only",
			signal: MarketSignal{
				CustomerCount: 10,
			},
			want: 10.0,
		},
		{
			name:   "empty signal",
			signal: MarketSignal{},
			want:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.signal.CalculateScore()
			if got != tt.want {
				t.Errorf("CalculateScore() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestMarketSignalARRInDollars(t *testing.T) {
	signal := MarketSignal{TotalARR: 12345}
	got := signal.ARRInDollars()
	want := 123.45
	if got != want {
		t.Errorf("ARRInDollars() = %f, want %f", got, want)
	}
}

func TestMarketSignalAvgARRPerCustomer(t *testing.T) {
	tests := []struct {
		name   string
		signal MarketSignal
		want   float64
	}{
		{
			name:   "normal case",
			signal: MarketSignal{TotalARR: 100000, CustomerCount: 4},
			want:   25000,
		},
		{
			name:   "zero customers",
			signal: MarketSignal{TotalARR: 100000, CustomerCount: 0},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.signal.AvgARRPerCustomer()
			if got != tt.want {
				t.Errorf("AvgARRPerCustomer() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestMarketSignalMerge(t *testing.T) {
	s1 := &MarketSignal{
		TotalVotes:    10,
		CustomerCount: 2,
		TotalARR:      50000,
		IdeaCount:     1,
		Sources:       []string{"aha"},
		CustomerIDs:   []string{"cust-1"},
	}

	s2 := &MarketSignal{
		TotalVotes:    15,
		CustomerCount: 3,
		TotalARR:      75000,
		IdeaCount:     2,
		Sources:       []string{"manual", "aha"}, // aha should be deduplicated
		CustomerIDs:   []string{"cust-2"},
	}

	s1.Merge(s2)

	if s1.TotalVotes != 25 {
		t.Errorf("Merged TotalVotes = %d, want 25", s1.TotalVotes)
	}
	if s1.CustomerCount != 5 {
		t.Errorf("Merged CustomerCount = %d, want 5", s1.CustomerCount)
	}
	if s1.TotalARR != 125000 {
		t.Errorf("Merged TotalARR = %d, want 125000", s1.TotalARR)
	}
	if s1.IdeaCount != 3 {
		t.Errorf("Merged IdeaCount = %d, want 3", s1.IdeaCount)
	}
	if len(s1.Sources) != 2 { // "aha" and "manual" (deduplicated)
		t.Errorf("Merged Sources length = %d, want 2", len(s1.Sources))
	}
	if len(s1.CustomerIDs) != 2 {
		t.Errorf("Merged CustomerIDs length = %d, want 2", len(s1.CustomerIDs))
	}
}

func TestMarketSignalHasSignal(t *testing.T) {
	empty := MarketSignal{}
	if empty.HasSignal() {
		t.Error("Empty signal HasSignal() = true, want false")
	}

	withVotes := MarketSignal{TotalVotes: 1}
	if !withVotes.HasSignal() {
		t.Error("Signal with votes HasSignal() = false, want true")
	}
}

func TestNewMarketSignalFromIdea(t *testing.T) {
	signal := NewMarketSignalFromIdea(25, 5, 100000, "aha")

	if signal.TotalVotes != 25 {
		t.Errorf("TotalVotes = %d, want 25", signal.TotalVotes)
	}
	if signal.CustomerCount != 5 {
		t.Errorf("CustomerCount = %d, want 5", signal.CustomerCount)
	}
	if signal.TotalARR != 100000 {
		t.Errorf("TotalARR = %d, want 100000", signal.TotalARR)
	}
	if signal.IdeaCount != 1 {
		t.Errorf("IdeaCount = %d, want 1", signal.IdeaCount)
	}
	if len(signal.Sources) != 1 || signal.Sources[0] != "aha" {
		t.Error("Sources not set correctly")
	}
	if signal.Score == 0 {
		t.Error("Score should be calculated")
	}
}

func TestAggregateMarketSignals(t *testing.T) {
	s1 := NewMarketSignalFromIdea(10, 2, 50000, "aha")
	s2 := NewMarketSignalFromIdea(20, 3, 75000, "manual")

	result := AggregateMarketSignals(s1, s2)

	if result.TotalVotes != 30 {
		t.Errorf("Aggregated TotalVotes = %d, want 30", result.TotalVotes)
	}
	if result.CustomerCount != 5 {
		t.Errorf("Aggregated CustomerCount = %d, want 5", result.CustomerCount)
	}
	if result.IdeaCount != 2 {
		t.Errorf("Aggregated IdeaCount = %d, want 2", result.IdeaCount)
	}
}
