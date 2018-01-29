package secure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputingHash(t *testing.T) {
	hash := computeHash("secret100")
	expected := "6c091667070fa0d70b8bab92755dec7af8d9adce498d08e18c6446e0d71d6cd9"
	assert.Equal(t, expected, hash, "Unequal hash")
}
