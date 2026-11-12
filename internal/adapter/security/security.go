// Package security provides the password KDF, opaque tokens, secret
// encryption and TOTP used by the auth use cases.
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Hasher implements Argon2id in the PHC string format.
type Hasher struct {
	time        uint32
	memory      uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewHasher() *Hasher {
	return &Hasher{time: 3, memory: 64 * 1024, parallelism: 2, saltLength: 16, keyLength: 32}
}

func (h *Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, h.time, h.memory, h.parallelism, h.keyLength)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.memory, h.time, h.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func (h *Hasher) Verify(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false
	}
	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	actual := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expected))) //nolint:gosec // The hash length is bounded by the PHC payload.
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

// Tokens mints opaque URL-safe tokens and deterministic digests.
type Tokens struct{}

func NewTokens() *Tokens { return &Tokens{} }

func (t *Tokens) NewToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (t *Tokens) HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

const backupAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func (t *Tokens) NewBackupCode() (string, error) {
	var builder strings.Builder
	builder.Grow(10)
	max := big.NewInt(int64(len(backupAlphabet)))
	for range 10 {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		builder.WriteByte(backupAlphabet[index.Int64()])
	}
	return builder.String(), nil
}

func (t *Tokens) HashBackupCode(code string) string {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

// Cipher encrypts TOTP secrets with AES-256-GCM. The key is derived from the
// application secret, so the database alone cannot reveal the secret.
type Cipher struct {
	aead cipher.AEAD
}

func NewCipher(secret []byte) (*Cipher, error) {
	if len(secret) == 0 {
		return nil, errors.New("secret cipher: empty key")
	}
	key := sha256.Sum256(secret)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{aead: aead}, nil
}

func (c *Cipher) Encrypt(plain string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := c.aead.Seal(nonce, nonce, []byte(plain), nil)
	return base64.RawURLEncoding.EncodeToString(out), nil
}

func (c *Cipher) Decrypt(encoded string) (string, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	size := c.aead.NonceSize()
	if len(data) < size {
		return "", errors.New("secret cipher: short payload")
	}
	plain, err := c.aead.Open(nil, data[:size], data[size:], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
