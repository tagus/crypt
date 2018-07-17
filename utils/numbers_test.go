package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindingMin(t *testing.T) {
	min := Min(1, 2, 3, -594849)
	assert.Equal(t, -594849, min)
}
