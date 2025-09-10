package core

import "time"

const ACTIVE_EXPIRE_FREQUENCY = 100 * time.Millisecond
const ACTIVE_EXPIRE_SAMPLE_SIZE = 20
const ACTIVE_EXPIRE_THRESHOLD = 0.1

func ActiveDeleteExpiredKeys() {
	for {
		expiredCount := 0
		sampleCountRemain := ACTIVE_EXPIRE_SAMPLE_SIZE
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

		if float64(expiredCount)/float64(ACTIVE_EXPIRE_SAMPLE_SIZE) < ACTIVE_EXPIRE_THRESHOLD {
			break
		}
	}
}
