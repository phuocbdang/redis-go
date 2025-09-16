package core

import "time"

const ActiveExpireFrequency = 100 * time.Millisecond
const ActiveExpireSampleSize = 20
const ActiveExpireThreshold = 0.1

func ActiveDeleteExpiredKeys() {
	for {
		expiredCount := 0
		sampleCountRemain := ActiveExpireSampleSize
		for key, expireTime := range dictStore.GetExpireDictStore() {
			sampleCountRemain--
			if sampleCountRemain < 0 {
				break
			}
			if time.Now().UnixMilli() > int64(expireTime) {
				dictStore.Del(key)
				expiredCount++
			}
		}

		if float64(expiredCount)/float64(ActiveExpireSampleSize) < ActiveExpireThreshold {
			break
		}
	}
}
