package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

type Argon2idParams struct {
	Time        uint32 `json:"time"`
	Memory      uint32 `json:"memory"`
	Parallelism uint8  `json:"parallelism"`
	KeyLength   uint32 `json:"key_length"`
}
type PasswordHash struct {
	Hash          string
	Salt          []byte
	Argon2idParam Argon2idParams
}

func HashPassword(password string) (*PasswordHash, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	time := uint32(3)
	memory := uint32(64 * 1024)
	parallelism := uint8(1)
	keyLen := uint32(32)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		parallelism,
		keyLen,
	)

	return &PasswordHash{
		Hash: base64.RawStdEncoding.EncodeToString(hash),
		Salt: salt,
		Argon2idParam: Argon2idParams{
			time,
			memory,
			parallelism,
			keyLen,
		},
	}, nil

}

func VerifyPassword(password string, stored *PasswordHash) bool {
	salt := stored.Salt
	expectedHash, _ := base64.RawStdEncoding.DecodeString(stored.Hash)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		stored.Argon2idParam.Time,
		stored.Argon2idParam.Memory,
		stored.Argon2idParam.Parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}
