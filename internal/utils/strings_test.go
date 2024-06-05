package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizingString(t *testing.T) {
	normalized := NormalizeString("Hello WoRld")
	assert.Equal(t, "hello_world", normalized)
}
