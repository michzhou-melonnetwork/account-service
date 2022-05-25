package address

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"github.com/Melon-Network-Inc/account-service/internal/entity"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateAddressRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateAddressRequest
		wantError bool
	}{
		{"success", CreateAddressRequest{Name: "test"}, false},
		{"required", CreateAddressRequest{Name: ""}, true},
		{"too long", CreateAddressRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateAddressRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateAddressRequest
		wantError bool
	}{
		{"success", UpdateAddressRequest{Name: "test"}, false},
		{"required", UpdateAddressRequest{Name: ""}, true},
		{"too long", UpdateAddressRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	address, err := s.Create(ctx, CreateAddressRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, address.ID)
	id := address.ID
	assert.Equal(t, "test", address.Name)
	assert.NotEmpty(t, address.CreatedAt)
	assert.NotEmpty(t, address.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateAddressRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateAddressRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateAddressRequest{Name: "test2"})

	// update
	address, err = s.Update(ctx, id, UpdateAddressRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", address.Name)
	_, err = s.Update(ctx, "none", UpdateAddressRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateAddressRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateAddressRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	address, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", address.Name)
	assert.Equal(t, id, address.ID)

	// query
	addresses, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(addresses))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	address, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, address.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.Address
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Address, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Address{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.Address, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, address entity.Address) error {
	if address.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, address)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, address entity.Address) error {
	if address.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == address.ID {
			m.items[i] = address
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
