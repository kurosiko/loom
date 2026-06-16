package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

// KeyDerivationParams defines parameters for Argon2 key derivation
type KeyDerivationParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

// DefaultKeyDerivationParams provides secure default parameters for Argon2
var DefaultKeyDerivationParams = KeyDerivationParams{
	Time:    3,
	Memory:  64 * 1024, // 64 MB
	Threads: 4,
	KeyLen:  32, // 256 bits for AES-256
	SaltLen: 16,
}

// GenerateSalt generates a cryptographically secure random salt
func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey derives a cryptographic key from a password using Argon2id
func DeriveKey(password string, salt []byte, params KeyDerivationParams) ([]byte, error) {
	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLen)
	return key, nil
}

// Encrypt encrypts plaintext using AES-256-GCM with authenticated encryption
func Encrypt(key, plaintext []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM with authenticated encryption
func Decrypt(key []byte, ciphertextStr string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// HashPassword hashes a password using Argon2id with automatic salt generation
func HashPassword(password string) (string, error) {
	salt, err := GenerateSalt(DefaultKeyDerivationParams.SaltLen)
	if err != nil {
		return "", err
	}

	key, err := DeriveKey(password, salt, DefaultKeyDerivationParams)
	if err != nil {
		return "", err
	}

	// Encode salt and hash together for storage
	encoded := base64.StdEncoding.EncodeToString(salt) + ":" + base64.StdEncoding.EncodeToString(key)
	return encoded, nil
}

// VerifyPassword verifies a password against a stored hash
func VerifyPassword(password, storedHash string) bool {
	parts := splitHash(storedHash)
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	storedKey, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	params := DefaultKeyDerivationParams
	params.KeyLen = uint32(len(storedKey))

	derivedKey, err := DeriveKey(password, salt, params)
	if err != nil {
		return false
	}

	return string(derivedKey) == string(storedKey)
}

func splitHash(hash string) []string {
	result := make([]string, 0)
	current := ""
	for _, c := range hash {
		if c == ':' {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

// HMAC computes HMAC-SHA256 of data with key
func HMAC(key, data []byte) []byte {
	h := sha256.New()
	h.Write(key)
	h.Write(data)
	return h.Sum(nil)
}
