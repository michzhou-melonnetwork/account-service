package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddress_UpdatePrimary_setFalse(t *testing.T) {
	testEntity := Address{
		ID:        "test1",
		Name:      "address1",
		Owner:     "test_account",
		Pubkey:    "9ZNTfG4NyQgxy2SWjSiQoUyBPEvXT2xo7fKc5hPYYJ7b",
		Currency:  "SOL",
		IsPrimary: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testEntity.UpdatePrimary(false)
	assert.Equal(t, false, testEntity.IsPrimary)
}

func TestAddress_UpdatePrimary_setTrue(t *testing.T) {
	testEntity := Address{
		ID:        "test1",
		Name:      "address1",
		Owner:     "test_account",
		Pubkey:    "9ZNTfG4NyQgxy2SWjSiQoUyBPEvXT2xo7fKc5hPYYJ7b",
		Currency:  "SOL",
		IsPrimary: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testEntity.UpdatePrimary(true)
	assert.Equal(t, true, testEntity.IsPrimary)
}
