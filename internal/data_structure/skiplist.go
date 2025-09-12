package data_structure

import (
	"math"
	"math/rand"
	"strings"
)

const SKIPLIST_MAX_LEVEL = 32

type SkiplistLevel struct {
	forward *SkiplistNode
	span    uint32
}

type SkiplistNode struct {
	ele      string
	score    float64
	backward *SkiplistNode
	levels   []SkiplistLevel
}

type Skiplist struct {
	head   *SkiplistNode
	tail   *SkiplistNode
	length uint32
	level  int
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Intn(2) == 1 {
		level++
	}
	if level > SKIPLIST_MAX_LEVEL {
		return SKIPLIST_MAX_LEVEL
	}
	return level
}

func (sl *Skiplist) CreateNode(level int, score float64, ele string) *SkiplistNode {
	return &SkiplistNode{
		ele:      ele,
		score:    score,
		backward: nil,
		levels:   make([]SkiplistLevel, level),
	}
}

func CreateSkiplist() *Skiplist {
	sl := Skiplist{
		length: 0,
		level:  1,
	}
	sl.head = sl.CreateNode(SKIPLIST_MAX_LEVEL, math.Inf(-1), "")
	sl.head.backward = nil
	sl.tail = nil
	return &sl
}

/*
Insert a new element to the Skiplist, allow duplicated scores
Caller should check if ele is already inserted or not
*/
// Insert adds a new node with a given score and element to the Skiplist
// The new node's level is determined probabilistically
func (sl *Skiplist) Insert(score float64, ele string) *SkiplistNode {
	// `update` stores the nodes that need to have their 'forward' pointers updated
	// at each level to insert the new node
	// `rank` stores the number of nodes visited at each level while searching
	// for the insertion position
	update := [SKIPLIST_MAX_LEVEL]*SkiplistNode{}
	rank := [SKIPLIST_MAX_LEVEL]uint32{}
	x := sl.head

	// Traverse the Skiplist from the highest level down to find the insertion point
	// This loop determines the `update` and `rank` arrays
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		// Move forward on the current level as long as the next node's score is less than
		// or equal to the new node's score
		// `strings.Compare` handless the case of equal scores to maintain a stable sort order
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score &&
					strings.Compare(x.levels[i].forward.ele, ele) == -1)) {
			// Accumulate the 'span' of each node to calculate the rank
			rank[i] += x.levels[i].span
			// Move to the next node
			x = x.levels[i].forward
		}
		// Store the last node visited at the level before dropping down
		update[i] = x
	}

	// Determine the level of the new node using a probabilistic method
	level := sl.randomLevel()
	// If the new node's level is higher than the current highest level of the Skiplist
	// update the Skiplist's state
	if level > sl.level {
		// For the new levels, the insertion point is the head of the Skiplist
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			// The span for new levels from the head to the end is the entire list length
			update[i].levels[i].span = sl.length
		}
		// Update the Skiplist's highest level
		sl.level = level
	}

	// Create new node and insert
	x = sl.CreateNode(level, score, ele)
	// Link the new node into the Skiplist at all its levels
	for i := 0; i < level; i++ {
		// Update the forward pointers to insert the new node
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x
		// Calculate the span for the new node
		x.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		// Update the span for the node before the new node
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	// Increase span for untouched level because we have a new node
	// For levels higher than the new node's level, the span of the `update` nodes
	// (which are the nodes before the insertion point) needs to be increased by one
	// This accounts for the new node being inserted "below" them
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	// Update the backward pointer for the new node, which is at the bottom level (0)
	if update[0] == sl.head {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	// Update the backward pointer of the node that comes after the new node
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x
	} else {
		// If the new node is the last one in the list, update the tail
		sl.tail = x
	}

	// Increase the total length of the Skiplist
	sl.length++
	// Return the newly inserted node
	return x
}

/*
Find the rank for an element by both score and key
Return 0 when the element cannot be found, rank otherwise
Note that the rank is 1-based due to the span of zsl->header to the
first element
*/
func (sl *Skiplist) GetRank(score float64, ele string) uint32 {
	x := sl.head
	var rank uint32 = 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score &&
					strings.Compare(x.levels[i].forward.ele, ele) <= 0)) {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x.score == score && strings.Compare(x.ele, ele) == 0 {
			return rank
		}
	}
	return 0
}

func (sl *Skiplist) DeleteNode(x *SkiplistNode, update [SKIPLIST_MAX_LEVEL]*SkiplistNode) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].forward == x {
			update[i].levels[i].span += x.levels[i].span - 1
			update[i].levels[i].forward = x.levels[i].forward
		} else {
			update[i].levels[i].span--
		}
	}
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x.backward
	} else {
		// x is tail
		sl.tail = x.backward
	}
	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *Skiplist) Delete(score float64, ele string) int {
	update := [SKIPLIST_MAX_LEVEL]*SkiplistNode{}
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score &&
					strings.Compare(x.levels[i].forward.ele, ele) == -1)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}
	x = x.levels[0].forward
	if x != nil && x.score == score && strings.Compare(x.ele, ele) == 0 {
		sl.DeleteNode(x, update)
		return 1
	}
	return 0
}

func (sl *Skiplist) UpdateScore(curScore float64, ele string, newScore float64) *SkiplistNode {
	update := [SKIPLIST_MAX_LEVEL]*SkiplistNode{}
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < curScore ||
				(x.levels[i].forward.score == curScore &&
					strings.Compare(x.levels[i].forward.ele, ele) == -1)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}
	x = x.levels[0].forward
	if (x.backward == nil || x.backward.score < newScore) &&
		(x.levels[0].forward != nil || newScore < x.levels[0].forward.score) {
		x.score = newScore
		return x
	}
	sl.DeleteNode(x, update)
	newNode := sl.Insert(newScore, ele)
	return newNode
}
