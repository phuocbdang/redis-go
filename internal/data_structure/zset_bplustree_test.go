package data_structure

import "testing"

func TestZSetBPlusTreeGetRank(t *testing.T) {
	zs := NewZSetBPlusTree(3)
	zs.Add(20.0, "k2")
	zs.Add(40.0, "k4")
	zs.Add(10.0, "k1")
	zs.Add(30.0, "k3")
	zs.Add(50.0, "k5")
	zs.Add(60.0, "k6")
	zs.Add(80.0, "k8")
	zs.Add(70.0, "k7")

	rank := zs.GetRank("k1")
	score, _ := zs.GetScore("k1")
	if rank != 0 || score != 10.0 {
		t.Fail()
	}

	rank = zs.GetRank("k2")
	score, _ = zs.GetScore("k2")
	if rank != 1 || score != 20.0 {
		t.Fail()
	}

	rank = zs.GetRank("k3")
	score, _ = zs.GetScore("k3")
	if rank != 2 || score != 30.0 {
		t.Fail()
	}

	rank = zs.GetRank("k4")
	score, _ = zs.GetScore("k4")
	if rank != 3 || score != 40.0 {
		t.Fail()
	}

	rank = zs.GetRank("k5")
	score, _ = zs.GetScore("k5")
	if rank != 4 || score != 50.0 {
		t.Fail()
	}

	rank = zs.GetRank("k6")
	score, _ = zs.GetScore("k6")
	if rank != 5 || score != 60.0 {
		t.Fail()
	}

	rank = zs.GetRank("k7")
	score, _ = zs.GetScore("k7")
	if rank != 6 || score != 70.0 {
		t.Fail()
	}

	rank = zs.GetRank("k8")
	score, _ = zs.GetScore("k8")
	if rank != 7 || score != 80.0 {
		t.Fail()
	}

}
