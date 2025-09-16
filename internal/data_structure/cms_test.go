package data_structure

import (
	"math"
	"testing"
)

func TestCreateCMS(t *testing.T) {
	w, d := uint32(10), uint32(5)
	cms := CreateCMS(w, d)
	if cms.width != w {
		t.Errorf("Expected width %d, got %d", w, cms.width)
	}
	if cms.depth != d {
		t.Errorf("Expected depth %d, got %d", d, cms.depth)
	}
	if len(cms.counter) != int(d) {
		t.Errorf("Expected counter rows %d, got %d", d, len(cms.counter))
	}
	for i := range cms.counter {
		if len(cms.counter[i]) != int(w) {
			t.Errorf("Expected counter columns %d, got %d", w, len(cms.counter[i]))
		}
	}
}

func TestCalcCMSDim(t *testing.T) {
	errRate := 0.01
	errProb := 0.001
	w, d := CalcCMSDim(errRate, errProb)
	if w == 0 || d == 0 {
		t.Errorf("Width and depth should be non-zero, got w=%d, d=%d", w, d)
	}
}

func TestIncrByAndCount(t *testing.T) {
	cms := CreateCMS(20, 5)
	item := "apple"
	count := cms.Count(item)
	if count != 0 {
		t.Errorf("Expected initial count 0, got %d", count)
	}

	// Increment by 3
	cms.IncrBy(item, 3)
	count = cms.Count(item)
	if count < 3 {
		t.Errorf("Expected count at least 3, got %d", count)
	}

	// Increment by 2
	cms.IncrBy(item, 2)
	count = cms.Count(item)
	if count < 5 {
		t.Errorf("Expected count at least 5, got %d", count)
	}
}

func TestIncrByOverflow(t *testing.T) {
	cms := CreateCMS(10, 3)
	item := "banana"
	// Set counter near max
	for i := uint32(0); i < cms.depth; i++ {
		hash := cms.calcHash(item, i)
		j := hash % cms.width
		cms.counter[i][j] = math.MaxUint32 - 1
	}
	cms.IncrBy(item, 10)
	count := cms.Count(item)
	if count != math.MaxUint32 {
		t.Errorf("Expected count to be capped at MaxUint32, got %d", count)
	}
}

func TestMultipleItems(t *testing.T) {
	cms := CreateCMS(50, 7)
	items := []string{"a", "b", "c", "d"}
	for i, item := range items {
		for j := 0; j <= i; j++ {
			cms.IncrBy(item, 1)
		}
	}
	for i, item := range items {
		count := cms.Count(item)
		if count < uint32(i+1) {
			t.Errorf("Expected count at least %d for item %s, got %d", i+1, item, count)
		}
	}
}

func TestCountNonExistentItem(t *testing.T) {
	cms := CreateCMS(10, 4)
	count := cms.Count("not-inserted")
	if count != 0 {
		t.Errorf("Expected count 0 for non-existent item, got %d", count)
	}
}