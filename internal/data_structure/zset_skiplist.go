package data_structure

type ZSet struct {
	zskiplist *Skiplist
	dict      map[string]float64
}

func CreateZSet() *ZSet {
	return &ZSet{
		zskiplist: CreateSkiplist(),
		dict:      make(map[string]float64),
	}
}

func (zs *ZSet) Add(score float64, ele string) int {
	if len(ele) == 0 {
		return 0
	}
	if curScore, exist := zs.dict[ele]; exist {
		if curScore != score {
			znode := zs.zskiplist.UpdateScore(curScore, ele, score)
			zs.dict[ele] = znode.score
		}
		return 1
	}
	znode := zs.zskiplist.Insert(score, ele)
	zs.dict[ele] = znode.score
	return 1
}

func (zs *ZSet) Delete(ele string) int {
	score, exist := zs.dict[ele]
	if !exist {
		return 0
	}
	delete(zs.dict, ele)
	zs.zskiplist.Delete(score, ele)
	return 1
}

/*
Returns the 0-based rank of the object or -1 if the object does not exist
If reverse if false, rank is computed considering as first element the one
with the lowest score. If reverse is true, rank is computed considering as element
with rank 0 the one with the highest score
*/
func (zs *ZSet) GetRank(ele string, reverse bool) (int64, float64) {
	size := zs.zskiplist.length
	score, exist := zs.dict[ele]
	if !exist {
		return -1, 0
	}
	rank := int64(zs.zskiplist.GetRank(score, ele))
	if reverse {
		rank = int64(size) - rank
	} else {
		rank--
	}
	return rank, score
}

func (zs *ZSet) GetScore(ele string) (float64, bool) {
	score, exist := zs.dict[ele]
	if !exist {
		return 0, false
	}
	return score, true
}

func (zs *ZSet) Len() int {
	return int(zs.zskiplist.length)
}
