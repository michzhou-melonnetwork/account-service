package address

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Melon-Network-Inc/account-service/internal/entity"
	"github.com/Melon-Network-Inc/account-service/internal/test"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "address")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	err = repo.Add(ctx, entity.Address{
		ID:        "test1",
		Name:      "address1",
		Owner:     "test_account",
		Pubkey:    "9ZNTfG4NyQgxy2SWjSiQoUyBPEvXT2xo7fKc5hPYYJ7b",
		Currency: "SOL",
		IsPrimary: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get
	address, err := repo.Get(ctx, "test1")
	assert.Nil(t, err)
	assert.Equal(t, "address1", address.Name)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, entity.Address{
		ID:        "test1",
		Name:      "address1 updated",
		Owner:     "test_account",
		Pubkey:    "9ZNTfG4NyQgxy2SWjSiQoUyBPEvXT2xo7fKc5hPYYJ7b",
		Currency: "SOL",
		IsPrimary: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	assert.Nil(t, err)
	address, _ = repo.Get(ctx, "test1")
	assert.Equal(t, "address1 updated", address.Name)

	// query
	addresses, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, len(addresses))

	// delete
	err = repo.Delete(ctx, "test1")
	assert.Nil(t, err)
	_, err = repo.Get(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
}
