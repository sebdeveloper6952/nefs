package main

import (
	"strings"
	"testing"

	"github.com/nbd-wtf/go-nostr/nip44"
	"github.com/stretchr/testify/assert"
)

func TestSplitString1Part(t *testing.T) {
	s := strings.Repeat("a", 10)

	parts := splitString(s, nip44.MaxPlaintextSize)

	assert.Equal(t, 1, len(parts), "should have 1 part")
}

func TestSplitStringMultipleParts(t *testing.T) {
	s := strings.Repeat("a", nip44.MaxPlaintextSize*5+1)

	parts := splitString(s, nip44.MaxPlaintextSize)

	assert.Equal(t, 6, len(parts), "should have 6 parts")
}
