package main

import "math"

func splitString(s string, maxPlaintextSize int) []string {
	need := int(math.Ceil(float64(len(s)) / float64(maxPlaintextSize)))
	parts := make([]string, need)

	for i := 0; i < need; i++ {
		start := i * maxPlaintextSize
		end := start + maxPlaintextSize
		if len(s) < end {
			end = len(s)
		}

		parts[i] = s[i*maxPlaintextSize : end]
	}

	return parts
}
