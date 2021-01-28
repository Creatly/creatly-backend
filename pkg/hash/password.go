package hash

import (
	"crypto/sha1"
	"fmt"
)

// PasswordHasher provides hashing logic to securely store passwords
type PasswordHasher interface {
	Hash(password string) string
}

// SHA1Hasher uses SHA1 to hash passwords with provided salt
type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

// Hash creates SHA1 hash of given password
func (h *SHA1Hasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
