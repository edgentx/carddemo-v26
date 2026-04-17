package shared

import (
	"math/rand"
	"time"
)

// GenerateID generates a random integer ID for testing/simulation.
func GenerateID() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63()
}
