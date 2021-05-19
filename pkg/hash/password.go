package hash

import (
	"crypto/sha1"
	"fmt"
)

// PasswordHasher provides hashing logic to securely store passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
}

// SHA1Hasher uses SHA1 to hash passwords with provided salt.
type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

// Hash creates SHA1 hash of given password.
func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
