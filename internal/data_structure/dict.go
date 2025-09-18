package data_structure

import (
	"log"
	"redisgo/internal/config"
	"time"
)

type Obj struct {
	Value          interface{}
	Ttl            int64
	LastAccessTime uint32
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func CreateDict() *Dict {
	return &Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
}

func now() uint32 {
	return uint32(time.Now().Unix())
}

func (d *Dict) NewObj(key string, value interface{}, ttlMs int64) *Obj {
	obj := &Obj{
		Value:          value,
		Ttl:            ttlMs,
		LastAccessTime: now(),
	}
	if ttlMs > 0 {
		d.SetExpiry(key, ttlMs)
	}
	return obj
}

func (d *Dict) Get(key string) *Obj {
	value := d.dictStore[key]
	if value != nil {
		value.LastAccessTime = now()
		if d.HasExpired(key) {
			d.Del(key)
			return nil
		}
	}
	return value
}

func (d *Dict) Set(key string, obj *Obj) {
	if len(d.dictStore) == config.EvictionMaxKeyNumber {
		d.evict()
	}
	v := d.dictStore[key]
	if v == nil {
		HashKeySpaceStat.Key++
	}
	d.dictStore[key] = obj
}

func (d *Dict) Del(key string) bool {
	if _, exist := d.dictStore[key]; exist {
		delete(d.dictStore, key)
		HashKeySpaceStat.Key--
		if _, exist := d.dictStore[key]; exist {
			delete(d.expiredDictStore, key)
			HashKeySpaceStat.Expire--
		}

		return true
	}
	return false
}

func (d *Dict) HasExpired(key string) bool {
	exp, exist := d.expiredDictStore[key]
	if !exist {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func (d *Dict) SetExpiry(key string, ttlMs int64) {
	if _, exist := d.expiredDictStore[key]; !exist {
		HashKeySpaceStat.Expire++
	}
	d.expiredDictStore[key] = uint64(time.Now().UnixMilli()) + uint64(ttlMs)
}

func (d *Dict) GetExpiry(key string) (uint64, bool) {
	exp, exist := d.expiredDictStore[key]
	return exp, exist
}

func (d *Dict) GetExpireDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) evictRandom() {
	log.Println("Trigger random eviction")
	evictCount := int64(config.EvictionRation * float64(config.EvictionMaxKeyNumber))
	for k := range d.dictStore {
		d.Del(k)
		evictCount--
		if evictCount == 0 {
			break
		}
	}
}

func (d *Dict) GetAvgTtl() int64 {
	var sumTtl int64 = 0
	var countKey int64 = 0
	for k := range d.expiredDictStore {
		sumTtl += d.dictStore[k].Ttl
		countKey += 1
		if countKey == int64(config.AvgTtlRandomSampleSize) {
			break
		}
	}
	if countKey == 0 {
		return 0
	}
	return sumTtl / countKey
}

func (d *Dict) populateEpool() {
	remain := config.EpoolLruSampleSize
	for k := range d.dictStore {
		ePool.Push(k, d.dictStore[k].LastAccessTime)
		remain--
		if remain == 0 {
			break
		}
	}
}

func (d *Dict) evictLru() {
	log.Println("Trigger lru eviction")
	d.populateEpool()
	evictCount := int64(config.EvictionRation * float64(config.EvictionMaxKeyNumber))
	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item != nil {
			d.Del(item.key)
		}
	}
}

func (d *Dict) evict() {
	switch config.EvictionPolicy {
	case "allkeys-random":
		d.evictRandom()
	case "allkeys-lru":
		d.evictLru()
	}
}
