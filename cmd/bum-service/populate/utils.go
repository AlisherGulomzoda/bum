//nolint:gosec //it's ok for generating test data.
package populate

import (
	"math/rand"
)

func randomIntFromInterval(min, max int) int {
	return rand.Intn(max-min) + min
}

func getPointOfInt8(val int8) *int8 {
	return &val
}
