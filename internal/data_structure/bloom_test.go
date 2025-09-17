package data_structure

import (
	"testing"
)

func TestBloomExist(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	b.Add("a")
	b.Add("b")
	if b.entries != 10 {
		t.Errorf("Expected entries to be 10, got %v", b.entries)
	}
	if b.errRate != 0.01 {
		t.Errorf("Expected Error to be 0.01, got %v", b.errRate)
	}
	if !b.Exist("a") {
		t.Errorf("Expected Exist(\"a\") to be true")
	}
	if !b.Exist("b") {
		t.Errorf("Expected Exist(\"b\") to be true")
	}
	if b.Exist("c") {
		t.Errorf("Expected Exist(\"c\") to be false")
	}
	if b.Exist("d") {
		t.Errorf("Expected Exist(\"d\") to be false")
	}
}

func TestBloomCalcHash(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	x := b.CalcHash("abcdef")
	y := b.CalcHash("abcdef")
	if x.a != y.a {
		t.Errorf("Expected x.a and y.a to be equal, got %v and %v", x.a, y.a)
	}
	if x.b != y.b {
		t.Errorf("Expected x.b and y.b to be equal, got %v and %v", x.b, y.b)
	}
}

func TestBloomAddHash(t *testing.T) {
	b := CreateBloomFilter(10, 0.01)
	hash := b.CalcHash("abcdef")
	b.AddHash(hash)
	if !b.ExistHash(hash) {
		t.Errorf("Expected ExistHash(hash) to be true")
	}
	if !b.Exist("abcdef") {
		t.Errorf("Expected Exist(\"abcdef\") to be true")
	}
}
