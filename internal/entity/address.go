package entity

import (
	"time"
)

// Address represents an address record.
type Address struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Owner     string    `json:"owenr"`
	Pubkey    string    `json:"pubkey"`
	Currency  string    `json:"currency"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (address *Address) UpdatePrimary(b bool) {
	address.IsPrimary = b
}
