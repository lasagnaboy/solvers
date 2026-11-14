package main

import (
	"math/rand"
	"testing"
)

func TestArrangeAvoidsLastYear(t *testing.T) {
	h, err := loadHistory("santa_history.json")
	if err != nil {
		t.Fatal(err)
	}
	santas := []string{"Salim", "Nico", "John", "Nicole", "Anna", "Susan", "Tim", "Janet"}
	forbidden := forbiddenPairs(h, 2026, 1)
	for seed := int64(0); seed < 500; seed++ {
		order, err := arrange(santas, forbidden, rand.New(rand.NewSource(seed)))
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		if len(order) != len(santas) {
			t.Fatalf("seed %d: got %d santas, want %d", seed, len(order), len(santas))
		}
		seen := map[string]bool{}
		for i, giver := range order {
			recipient := order[(i+1)%len(order)]
			if seen[giver] {
				t.Fatalf("seed %d: %s appears twice", seed, giver)
			}
			seen[giver] = true
			if giver == recipient {
				t.Fatalf("seed %d: %s gives to self", seed, giver)
			}
			if h["2025"][giver] == recipient {
				t.Fatalf("seed %d: %s => %s repeats 2025", seed, giver, recipient)
			}
		}
	}
}

func TestArrangeReportsImpossible(t *testing.T) {
	// Two people who swapped last year can't avoid each other.
	h := History{"2025": {"Tim": "Janet", "Janet": "Tim"}}
	if _, err := arrange([]string{"Tim", "Janet"}, forbiddenPairs(h, 2026, 1), rand.New(rand.NewSource(1))); err == nil {
		t.Fatal("expected an error when no valid arrangement exists")
	}
}

func TestForbiddenPairsLookback(t *testing.T) {
	h := History{"2024": {"A": "B"}, "2025": {"A": "C"}, "2026": {"A": "D"}}
	f := forbiddenPairs(h, 2026, 1)
	if !f["A"]["C"] || f["A"]["B"] || f["A"]["D"] {
		t.Fatalf("lookback 1: got %v", f["A"])
	}
	f = forbiddenPairs(h, 2026, 0)
	if !f["A"]["B"] || !f["A"]["C"] || f["A"]["D"] {
		t.Fatalf("lookback 0 (all): got %v", f["A"])
	}
}
