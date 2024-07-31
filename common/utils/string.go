package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func GetRandomString(length int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	var (
		result []byte
		b      []byte
		r      *rand.Rand
	)
	b = []byte(str)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, b[r.Intn(len(b))])
	}
	return string(result)
}

func ConvertToFloat(num int64) string {
	if num > 10000 {
		return fmt.Sprintf("%.1fk", float64(num)/1000.0) // 超过 1w 显示 k
	}
	return fmt.Sprintf("%d", num)
}

func CheckStringPattern(str, pattern string) bool {
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return true
	} else {
		return false
	}
}
