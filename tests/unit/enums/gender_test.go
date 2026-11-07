package enums_test

import (
	"encoding/json"
	"testing"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

func TestGenderIdx_ScanAndValue(t *testing.T) {
	var g enums.GenderIdx

	// Test scan int16
	if err := g.Scan(int16(1)); err != nil || g != enums.GenderMaleIdx {
		t.Fatalf("expected GenderMaleIdx, got %v, err: %v", g, err)
	}

	// Test scan string
	if err := g.Scan("female"); err != nil || g != enums.GenderFemaleIdx {
		t.Fatalf("expected GenderFemaleIdx, got %v, err: %v", g, err)
	}

	// Test value
	v, err := g.Value()
	if err != nil || v != int64(2) {
		t.Fatalf("expected int64(2), got %v, err: %v", v, err)
	}

	// Test IsValidIdx
	if !g.IsValidIdx() {
		t.Fatalf("expected valid idx for female")
	}

	g = enums.GenderUnknownIdx
	if g.IsValidIdx() {
		t.Fatalf("expected invalid idx for unknown")
	}
}

func TestGenderIdx_JSON(t *testing.T) {
	type Sample struct {
		Gender enums.GenderIdx `json:"gender"`
	}

	// Test unmarshal integer
	var s1 Sample
	if err := json.Unmarshal([]byte(`{"gender": 1}`), &s1); err != nil || s1.Gender != enums.GenderMaleIdx {
		t.Fatalf("expected GenderMaleIdx, got %v, err: %v", s1.Gender, err)
	}

	// Test unmarshal string
	var s2 Sample
	if err := json.Unmarshal([]byte(`{"gender": "female"}`), &s2); err != nil || s2.Gender != enums.GenderFemaleIdx {
		t.Fatalf("expected GenderFemaleIdx, got %v, err: %v", s2.Gender, err)
	}

	// Test marshal
	b, err := json.Marshal(s2)
	if err != nil || string(b) != `{"gender":2}` {
		t.Fatalf("expected {\"gender\":2}, got %s", string(b))
	}
}
