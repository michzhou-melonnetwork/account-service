package address

import (
	"context"

	"github.com/Melon-Network-Inc/account-service/internal/entity"
	"github.com/Melon-Network-Inc/account-service/pkg/dbcontext"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Repository encapsulates the logic to access addresses from the data source.
type Repository interface {
	// Get returns the address with the specified address ID.
	Get(ctx context.Context, id string) (entity.Address, error)
	// List returns all addresses owned by the user.
	List(ctx context.Context) ([]entity.Address, error)
	// GetPrimaryAddress returns the primary address.
	GetPrimaryAddress(ctx context.Context) (entity.Address, error)
	// Count returns the number of addresses.
	Count(ctx context.Context) (int, error)
	// Query returns the list of addresses with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.Address, error)
	// Add saves a new address in the storage.
	Add(ctx context.Context, address entity.Address) error
	// Update updates the address with given ID in the storage.
	Update(ctx context.Context, address entity.Address) error
	// UpdatePrimaryAddress updates the address with the given isPrimary value.
	UpdatePrimaryAddress(ctx context.Context, isPrimary bool) error
	// Delete removes the address with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists addresses in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new address repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the address with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Address, error) {
	var address entity.Address
	err := r.db.With(ctx).Select().Model(id, &address)
	return address, err
}

// List returns all addresses owned by the user.
func (r repository) List(ctx context.Context) ([]entity.Address, error) {
	owner := ctx.Value("user")
	var addresses []entity.Address
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).All(addresses)
	return addresses, err
}

// Get reads the address with the specified ID from the database.
func (r repository) GetPrimaryAddress(ctx context.Context) (entity.Address, error) {
	owner := ctx.Value("user")
	var address entity.Address
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"owner": owner, "is_primary": true}).
		One(&address)
	return address, err
}

// Create saves a new address record in the database.
// It returns the ID of the newly inserted address record.
func (r repository) Add(ctx context.Context, address entity.Address) error {
	return r.db.With(ctx).Model(&address).Insert()
}

// Update saves the changes to an address in the database.
func (r repository) Update(ctx context.Context, address entity.Address) error {
	return r.db.With(ctx).Model(&address).Update()
}

// UpdatePrimary updates the primary field to an address in the database.
func (r repository) UpdatePrimaryAddress(ctx context.Context, isPrimary bool) error {
	var primary Address
	if primary, err := r.GetPrimaryAddress(ctx); err != nil {
		primary.UpdatePrimary(isPrimary)
	}
	return r.Update(ctx, primary.Address)
}

// Delete deletes an address with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	address, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&address).Delete()
}

// Count returns the number of the address records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("address").Row(&count)
	return count, err
}

// Query retrieves the address records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Address, error) {
	var addresses []entity.Address
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&addresses)
	return addresses, err
}
