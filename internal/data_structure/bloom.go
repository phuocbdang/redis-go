package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

const BfDefaultInitCapacity = 100
const BfDefaultErrRate = 0.01

const Ln2 float64 = 0.693147180559945
const Ln2Square float64 = 0.480453013918201
const ABigSeed uint32 = 0x9747b28c

type Bloom struct {
	hashes      int
	entries     uint64
	errRate     float64
	bitPerEntry float64
	bf          []bool
	bits        uint64
}

type HashValue struct {
	a uint64
	b uint64
}

/* http://en.wikipedia.org/wiki/Bloom_filter
 * - Optimal number of bits is: bits = (entries * ln(errRate)) / ln(2)^2
 * - bitPerEntry = bits/entries
 * - Optimal number of hash function is: hashes = bitPerEntry * ln(2)
 */
func CreateBloomFilter(entires uint64, errRate float64) *Bloom {
	bloom := &Bloom{
		entries: entires,
		errRate: errRate,
	}
	bloom.bits = uint64(math.Abs(float64(entires)*math.Log(errRate)) / Ln2Square)
	bloom.bitPerEntry = float64(bloom.bits) / float64(bloom.entries)
	bloom.hashes = int(math.Ceil(Ln2 * float64(bloom.bitPerEntry)))
	bloom.bf = make([]bool, bloom.bits)
	return bloom
}

func (b *Bloom) CalcHash(entry string) HashValue {
	hasher := murmur3.New128WithSeed(ABigSeed)
	hasher.Write([]byte(entry))
	x, y := hasher.Sum128()
	return HashValue{
		a: x,
		b: y,
	}
}

func (b *Bloom) Add(entry string) int {
	var hash uint64
	added := 0
	initHash := b.CalcHash(entry)
	for i := 0; i < b.hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		if !b.bf[hash] {
			b.bf[hash] = true
			added = 1
		}
	}
	return added
}

func (b *Bloom) Exist(entry string) bool {
	var hash uint64
	initHash := b.CalcHash(entry)
	for i := 0; i < b.hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		if !b.bf[hash] {
			return false
		}
	}
	return true
}

func (b *Bloom) AddHash(initHash HashValue) {
	var hash uint64
	for i := 0; i < b.hashes; i++ {
		hash = (initHash.a + initHash.b*(uint64(i))) % b.bits
		b.bf[hash] = true
	}
}

func (b *Bloom) ExistHash(initHash HashValue) bool {
	var hash uint64
	for i := 0; i < b.hashes; i++ {
		hash = (initHash.a + initHash.b*uint64(i)) % b.bits
		if !b.bf[hash] {
			return false
		}
	}
	return true
}
