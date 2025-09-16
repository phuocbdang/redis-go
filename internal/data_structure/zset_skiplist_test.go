package data_structure

import (
	"testing"
)

func TestZSetSkiplistGetRank(t *testing.T) {
	zs := CreateZSet()
	zs.Add(20.0, "k2")
	zs.Add(40.0, "k4")
	zs.Add(10.0, "k1")
	zs.Add(30.0, "k3")
	zs.Add(50.0, "k5")
	zs.Add(60.0, "k6")
	zs.Add(80.0, "k8")
	zs.Add(70.0, "k7")

	tests := []struct {
		key   string
		rank  int
		score float64
	}{
		{"k1", 0, 10.0},
		{"k2", 1, 20.0},
		{"k3", 2, 30.0},
		{"k4", 3, 40.0},
		{"k5", 4, 50.0},
		{"k6", 5, 60.0},
		{"k7", 6, 70.0},
		{"k8", 7, 80.0},
	}

	for _, tt := range tests {
		rank, score := zs.GetRank(tt.key, false)
		if rank != int64(tt.rank) {
			t.Errorf("GetRank(%q) rank = %d, want %d", tt.key, rank, tt.rank)
		}
		if score != tt.score {
			t.Errorf("GetRank(%q) score = %v, want %v", tt.key, score, tt.score)
		}
	}
}
