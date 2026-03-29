package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/box"
)

// KeyPair represents an Ed25519 keypair (for signatures)
type KeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

// EncryptionKeyPair represents an X25519 keypair (for encryption)
type EncryptionKeyPair struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
}

// GenerateKeyPair generates an Ed25519 keypair
func GenerateKeyPair() (*KeyPair, error) {
	pubkey, privkey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 keypair: %w", err)
	}
	return &KeyPair{
		PrivateKey: privkey,
		PublicKey:  pubkey,
	}, nil
}

// GenerateEncryptionKeyPair generates an X25519 keypair
func GenerateEncryptionKeyPair() (*EncryptionKeyPair, error) {
	pubkey, privkey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate X25519 keypair: %w", err)
	}
	return &EncryptionKeyPair{
		PrivateKey: *privkey,
		PublicKey:  *pubkey,
	}, nil
}

// DeriveX25519 derives an X25519 public key from Ed25519 public key
// This allows the same keypair to be used for both signing and encryption
func DeriveX25519(ed25519PubKey ed25519.PublicKey) ([32]byte, error) {
	var x25519PubKey [32]byte
	curve25519.ScalarBaseMult(&x25519PubKey, (*[32]byte)(ed25519PubKey))
	return x25519PubKey, nil
}

// HashPassword hashes a password using Argon2id (memory-hard)
func HashPassword(password string, salt []byte) ([]byte, error) {
	// Argon2id parameters (memory-hard for modern hardware)
	timeCost    := uint32(3)
	memoryCost  := uint32(64 * 1024) // 64 MB
	parallelism := uint8(4)
	keyLength   := uint32(32)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		timeCost,
		memoryCost,
		parallelism,
		keyLength,
	)
	return hash, nil
}

// GenerateSalt generates a random salt for password hashing
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// BoxEncrypt encrypts a message using X25519 public key (NaCl box)
func BoxEncrypt(message []byte, nonce *[24]byte, recipientPubKey, senderPrivKey *[32]byte) ([]byte, error) {
	var encrypted []byte
	encrypted = box.Seal(encrypted, message, nonce, recipientPubKey, senderPrivKey)
	return encrypted, nil
}

// BoxDecrypt decrypts a message using X25519 private key
func BoxDecrypt(encrypted []byte, nonce *[24]byte, senderPubKey, recipientPrivKey *[32]byte) ([]byte, bool) {
	var decrypted []byte
	decrypted, success := box.Open(decrypted, encrypted, nonce, senderPubKey, recipientPrivKey)
	return decrypted, success
}

// Sign signs a message using Ed25519 private key
func Sign(message []byte, privkey ed25519.PrivateKey) []byte {
	return ed25519.Sign(privkey, message)
}

// Verify verifies a signature using Ed25519 public key
func Verify(message, signature []byte, pubkey ed25519.PublicKey) bool {
	return ed25519.Verify(pubkey, message, signature)
}

// GenerateNonce generates a random nonce for NaCl box encryption
func GenerateNonce() (*[24]byte, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return &nonce, nil
}

// HashKey creates a blind index for a public key
func HashKey(pubkey []byte) []byte {
	hash := sha256.Sum256(pubkey)
	return hash[:]
}

// EncodeKey encodes a key to base64
func EncodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// DecodeKey decodes a base64 key
func DecodeKey(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password string, hash, salt []byte) (bool, error) {
	newHash, err := HashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(newHash, hash) == 1, nil
}
