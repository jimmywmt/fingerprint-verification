package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func DecryptFingerprint(sharedSecret string, nonceHex string, ciphertextHex string) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", fmt.Errorf("invalid ciphertext: %w", err)
	}

	nonceBytes, err := hex.DecodeString(nonceHex)
	if err != nil {
		return "", fmt.Errorf("invalid nonce: %w", err)
	}

	key := sha256.Sum256([]byte(sharedSecret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, nonceBytes, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	parts := strings.SplitN(string(plaintext), "@@", 2)
	if len(parts) != 2 || parts[1] != sharedSecret {
		return "", fmt.Errorf("shared secret mismatch in fingerprint")
	}

	return parts[0], nil
}
