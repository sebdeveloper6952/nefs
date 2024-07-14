package main

import "math"

func splitString(s string, MaxPlaintextSize int) []string {
	need := int(math.Ceil(float64(len(s)) / float64(MaxPlaintextSize)))
	parts := make([]string, need)

	for i := 0; i < need; i++ {
		start := i * MaxPlaintextSize
		end := start + MaxPlaintextSize
		if len(s) < end {
			end = len(s)
		}

		parts[i] = s[i*MaxPlaintextSize : end]
	}

	return parts
}
