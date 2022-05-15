package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

var randNumSrc = rand.NewSource(time.Now().UnixNano())

func CreateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index bc 2^6
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	for i, cache, remain := n-1, randNumSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randNumSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func ToString(num int) string {
	return fmt.Sprintf("%d", num)
}

func RemoveStringFromSlice(arr *[]string, s string) {
	lag := 0
	for i := 0; i < len(*arr); i++ {
		if (*arr)[i] != s {
			(*arr)[lag] = (*arr)[i]
			lag++
		}
	}
	for lag > 0 {
		*arr = (*arr)[:len(*arr)-1]
		lag--
	}
}
