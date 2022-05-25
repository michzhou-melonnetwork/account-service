package address

import (
	"context"
	"time"

	"github.com/Melon-Network-Inc/account-service/internal/entity"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Service encapsulates usecase logic for addresses.
type Service interface {
	Get(ctx context.Context, id string) (Address, error)
	Query(ctx context.Context, offset, limit int) ([]Address, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateAddressRequest) (Address, error)
	Update(ctx context.Context, id string, input UpdateAddressRequest) (Address, error)
	Delete(ctx context.Context, id string) (Address, error)
}

// address represents the data about an address.
type Address struct {
	entity.Address
}

// CreateAddressRequest represents an address creation request.
type CreateAddressRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateAddressRequest fields.
func (m CreateAddressRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateAddressRequest represents an address update request.
type UpdateAddressRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateAddressRequest fields.
func (m UpdateAddressRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new address service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the address with the specified the address ID.
func (s service) Get(ctx context.Context, id string) (Address, error) {
	address, err := s.repo.Get(ctx, id)
	if err != nil {
		return Address{}, err
	}
	return Address{address}, nil
}

// Create creates a new address.
func (s service) Create(ctx context.Context, req CreateAddressRequest) (Address, error) {
	if err := req.Validate(); err != nil {
		return Address{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.Address{
		ID:        id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Address{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the address with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateAddressRequest) (Address, error) {
	if err := req.Validate(); err != nil {
		return Address{}, err
	}

	address, err := s.Get(ctx, id)
	if err != nil {
		return address, err
	}
	address.Name = req.Name
	address.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, address.Address); err != nil {
		return address, err
	}
	return address, nil
}

// Delete deletes the address with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Address, error) {
	address, err := s.Get(ctx, id)
	if err != nil {
		return Address{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Address{}, err
	}
	return address, nil
}

// Count returns the number of addresses.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the addresses with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]Address, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Address{}
	for _, item := range items {
		result = append(result, Address{item})
	}
	return result, nil
}
