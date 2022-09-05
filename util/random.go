package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates random string of length n
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabets)
	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generate random owner
func RandomUsername() string {
	return RandomString(6)
}

// Generate random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
