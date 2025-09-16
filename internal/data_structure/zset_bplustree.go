package data_structure

const DEFAULT_BPLUS_TREE_DEGREE = 4

type ZSetBPlusTree struct {
	Tree         *BPlusTree
	MemberScores map[string]float64
}

func NewZSetBPlusTree(degree int) *ZSetBPlusTree {
	return &ZSetBPlusTree{
		Tree:         NewBPlusTree(degree),
		MemberScores: make(map[string]float64),
	}
}

func (ss *ZSetBPlusTree) Add(score float64, member string) int {
	return ss.Tree.Add(score, member)
}

func (ss *ZSetBPlusTree) GetScore(member string) (float64, bool) {
	return ss.Tree.Score(member)
}

func (ss *ZSetBPlusTree) GetRank(member string) int {
	return ss.Tree.GetRank(member)
}
