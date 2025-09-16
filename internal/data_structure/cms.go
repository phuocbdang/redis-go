package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

// LOG_10_POINT_FIVE is a precomputed for log10(0.5)
const LOG_10_POINT_FIVE = -0.30102999566

// CMS stands for Count-min Sketch data structure
// The counter filed has been changed to a 2D slice for better clarity and indexing
type CMS struct {
	width   uint32
	depth   uint32
	counter [][]uint32
}

// CreateCMS initializes a new Count-min Sketch with a given with and depth
func CreateCMS(w uint32, d uint32) *CMS {
	cms := &CMS{
		width: w,
		depth: d,
	}
	// Initialize the 2D slice
	cms.counter = make([][]uint32, d)
	for i := uint32(0); i < d; i++ {
		cms.counter[i] = make([]uint32, w)
	}
	return cms
}

// CalcCMSDim calculates the dimensions (width and depth) of the CMS
// based on the desired error rate and probability
func CalcCMSDim(errRate float64, probability float64) (uint32, uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(probability) / LOG_10_POINT_FIVE))
	return w, d
}

// calcHash calculates a 32-bit hash for the given item and seed
func (c *CMS) calcHash(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}

// IncrBy increments the count for an item by a specific value
// It returns the estimated count for the item after the increment
func (c *CMS) IncrBy(item string, value uint32) uint32 {
	var minCount uint32 = math.MaxUint32
	// Loop through the each row of the 2D array
	for i := uint32(0); i < c.depth; i++ {
		// Calculate a new hash for each row using the row index as the seed
		hash := c.calcHash(item, i)
		// Use the hash to get the column index within the row
		j := hash % c.width

		// Safely add the value to prevent overflow
		if math.MaxUint32-c.counter[i][j] < value {
			c.counter[i][j] = math.MaxUint32
		} else {
			c.counter[i][j] += value
		}

		// Keep track of the minimum count across all rows
		minCount = min(c.counter[i][j], minCount)
	}
	return minCount
}

// Count returns the estimated count for an item
// It retrieves the minimum count across all hash functions to provide the most accurate estimate
func (c *CMS) Count(item string) uint32 {
	var minCount uint32 = math.MaxUint32
	// Loop through each row of the 2D array
	for i := uint32(0); i < c.depth; i++ {
		// Calculate the hash for this row
		hash := c.calcHash(item, i)
		// Determine the column index
		j := hash % c.width

		// Find the minimum count across all rows
		minCount = min(minCount, c.counter[i][j])
	}
	return minCount
}
