package nefs

import "math"

func SplitString(s string, chunkLength int) []string {
	need := int(math.Ceil(float64(len(s)) / float64(chunkLength)))
	parts := make([]string, need)

	for i := 0; i < need; i++ {
		start := i * chunkLength
		end := start + chunkLength
		if len(s) < end {
			end = len(s)
		}

		parts[i] = s[i*chunkLength : end]
	}

	return parts
}
