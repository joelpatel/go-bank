package utils

import (
	"math/rand"
	"strings"
)

const (
	alphabet       = "abcdefghijklmnopqrstuvwxyz"
	alphabetLength = 26
)

// Generates a random integer between min and max (inclusive).
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generates a random string of length n.
func RandomString(n int) string {
	var builder strings.Builder

	for i := 0; i < n; i++ {
		builder.WriteByte(alphabet[rand.Intn(alphabetLength)])
	}

	return builder.String()
}

// Generates a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// Generates a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
