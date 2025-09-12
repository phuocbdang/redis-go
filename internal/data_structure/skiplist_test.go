package data_structure

import (
	"testing"
)

func TestCreateSkiplist(t *testing.T) {
	sl := CreateSkiplist()
	if sl.length != 0 {
		t.Fail()
	}
	if sl.level != 1 {
		t.Fail()
	}
	if len(sl.head.levels) != SKIPLIST_MAX_LEVEL {
		t.Fail()
	}
	if sl.tail != nil {
		t.Fail()
	}
}

func TestInsertSkiplist(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")

	if sl.length != 4 {
		t.Fail()
	}
	if sl.head.levels[0].forward.score != 10 {
		t.Fail()
	}
	if sl.head.levels[0].forward.ele != "k1" {
		t.Fail()
	}
	if sl.head.levels[0].forward.levels[0].forward.score != 10 {
		t.Fail()
	}
	if sl.head.levels[0].forward.levels[0].forward.ele != "k2" {
		t.Fail()
	}
	if sl.head.levels[0].forward.levels[0].forward.levels[0].forward.score != 20 {
		t.Fail()
	}
	if sl.head.levels[0].forward.levels[0].forward.levels[0].forward.ele != "k3" {
		t.Fail()
	}
	if sl.head.levels[0].forward.backward != nil {
		t.Fail()
	}
	if sl.head.levels[0].forward.levels[0].forward.levels[0].forward.levels[0].forward != sl.tail {
		t.Fail()
	}
	if sl.head.levels[0].forward.levels[0].forward.levels[0].forward != sl.tail.backward {
		t.Fail()
	}
	for i := sl.level; i < SKIPLIST_MAX_LEVEL; i++ {
		if sl.head.levels[i].span != 0 {
			t.Fail()
		}
		if sl.head.levels[i].forward != nil {
			t.Fail()
		}
	}
}

func TestGetRankSkiplist(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(50, "k5")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")
	sl.Insert(50, "k6")

	if sl.GetRank(10, "k1") != 1 {
		t.Fail()
	}
	if sl.GetRank(10, "k2") != 2 {
		t.Fail()
	}
	if sl.GetRank(20, "k3") != 3 {
		t.Fail()
	}
	if sl.GetRank(40, "k4") != 4 {
		t.Fail()
	}
	if sl.GetRank(50, "k5") != 5 {
		t.Fail()
	}
	if sl.GetRank(50, "k6") != 6 {
		t.Fail()
	}
}

func TestDeleteSkiplist(t *testing.T) {
	sl := CreateSkiplist()
	sl.Insert(10, "k1")
	sl.Insert(20, "k3")
	sl.Insert(40, "k4")
	sl.Insert(10, "k2")

	res := sl.Delete(10, "k5")
	if res != 0 {
		t.Errorf("expected 0, got %v", res)
	}

	res = sl.Delete(30, "k5")
	if res != 0 {
		t.Errorf("expected 0, got %v", res)
	}

	res = sl.Delete(20, "k3")
	if res != 1 {
		t.Errorf("expected 1, got %v", res)
	}
	if sl.length != 3 {
		t.Errorf("expected length 3, got %v", sl.length)
	}
	node1 := sl.head.levels[0].forward
	node2 := node1.levels[0].forward
	node3 := node2.levels[0].forward
	if node1.score != 10 {
		t.Errorf("expected node1.score 10, got %v", node1.score)
	}
	if node1.ele != "k1" {
		t.Errorf("expected node1.ele 'k1', got %v", node1.ele)
	}
	if node2.score != 10 {
		t.Errorf("expected node2.score 10, got %v", node2.score)
	}
	if node2.ele != "k2" {
		t.Errorf("expected node2.ele 'k2', got %v", node2.ele)
	}
	if node3.score != 40 {
		t.Errorf("expected node3.score 40, got %v", node3.score)
	}
	if node3.ele != "k4" {
		t.Errorf("expected node3.ele 'k4', got %v", node3.ele)
	}
	if node3.backward.score != 10 {
		t.Errorf("expected node3.backward.score 10, got %v", node3.backward.score)
	}
	if sl.tail.backward.score != 10 {
		t.Errorf("expected sl.tail.backward.score 10, got %v", sl.tail.backward.score)
	}
	if node3 != sl.tail {
		t.Errorf("expected node3 == sl.tail")
	}
	for i := sl.level; i < SKIPLIST_MAX_LEVEL; i++ {
		if sl.head.levels[i].forward != nil {
			t.Errorf("expected sl.head.levels[%v].forward == nil", i)
		}
	}

	res = sl.Delete(10, "k1")
	if res != 1 {
		t.Errorf("expected 1, got %v", res)
	}
	if sl.length != 2 {
		t.Errorf("expected length 2, got %v", sl.length)
	}
	
	node1 = sl.head.levels[0].forward
	node2 = node1.levels[0].forward
	if node1.score != 10 {
		t.Errorf("expected node1.score 10, got %v", node1.score)
	}
	if node1.ele != "k2" {
		t.Errorf("expected node1.ele 'k2', got %v", node1.ele)
	}
	if node2.score != 40 {
		t.Errorf("expected node2.score 40, got %v", node2.score)
	}
	if node2.ele != "k4" {
		t.Errorf("expected node2.ele 'k4', got %v", node2.ele)
	}
	if node2 != sl.tail {
		t.Errorf("expected node2 == sl.tail")
	}
	for i := sl.level; i < SKIPLIST_MAX_LEVEL; i++ {
		if sl.head.levels[i].forward != nil {
			t.Errorf("expected sl.head.levels[%v].forward == nil", i)
		}
	}
}
