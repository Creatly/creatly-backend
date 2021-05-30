package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provides hashing logic to securely store passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	CompareHashAndPassword(hash, password string) error
}

// BcryptHasher uses Bcrypt to hash passwords with provided cost.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates new instance of the BcryptHasher with provided cost.
// If the provided cost is lower than bcrypt.MinCost or greater than bcrypt.MaxCost
// then bcrypt.DefaultCost will be used.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *BcryptHasher) CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
