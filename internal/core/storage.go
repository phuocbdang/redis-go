package core

import "redisgo/internal/data_structure"

var dictStore *data_structure.Dict
var setStore map[string]*data_structure.SimpleSet
var zsetStore map[string]*data_structure.SortedSet

func init() {
	dictStore = data_structure.CreateDict()
	setStore = make(map[string]*data_structure.SimpleSet)
	zsetStore = make(map[string]*data_structure.SortedSet)
}
