package data_structure

import "testing"

func TestZSetGetRank(t *testing.T) {
	ss := NewSortedSet(3)
	ss.Add(20.0, "k2")
	ss.Add(40.0, "k4")
	ss.Add(10.0, "k1")
	ss.Add(30.0, "k3")
	ss.Add(50.0, "k5")
	ss.Add(60.0, "k6")
	ss.Add(80.0, "k8")
	ss.Add(70.0, "k7")

	rank := ss.GetRank("k1")
	score, _ := ss.GetScore("k1")
	if rank != 0 || score != 10.0 {
		t.Fail()
	}

	rank = ss.GetRank("k2")
	score, _ = ss.GetScore("k2")
	if rank != 1 || score != 20.0 {
		t.Fail()
	}

	rank = ss.GetRank("k3")
	score, _ = ss.GetScore("k3")
	if rank != 2 || score != 30.0 {
		t.Fail()
	}

	rank = ss.GetRank("k4")
	score, _ = ss.GetScore("k4")
	if rank != 3 || score != 40.0 {
		t.Fail()
	}

	rank = ss.GetRank("k5")
	score, _ = ss.GetScore("k5")
	if rank != 4 || score != 50.0 {
		t.Fail()
	}

	rank = ss.GetRank("k6")
	score, _ = ss.GetScore("k6")
	if rank != 5 || score != 60.0 {
		t.Fail()
	}

	rank = ss.GetRank("k7")
	score, _ = ss.GetScore("k7")
	if rank != 6 || score != 70.0 {
		t.Fail()
	}

	rank = ss.GetRank("k8")
	score, _ = ss.GetScore("k8")
	if rank != 7 || score != 80.0 {
		t.Fail()
	}

}
