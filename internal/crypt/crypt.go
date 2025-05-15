package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/zalando/go-keyring"
)

const service = "dev.frankmayer.request-review"

func Encrypt(plaintext string) (string, error) {
	key, err := getKey()
	if err != nil {
		return "", errors.Join(errors.New("Error getting AES key"), err)
	}
	enc, err := encryptAESGCM(key, plaintext)
	if err != nil {
		return "", errors.Join(errors.New("Error encrypting plaintext"), err)
	}
	return enc, nil
}

func Decrypt(ciphertext string) (string, error) {
	key, err := getKey()
	if err != nil {
		return "", errors.Join(errors.New("Error getting AES key"), err)
	}
	dec, err := decryptAESGCM(key, ciphertext)
	if err != nil {
		return "", errors.Join(errors.New("Error decrypting ciphertext"), err)
	}
	return dec, nil
}

func getKey() (string, error) {
	if k, err := keyring.Get(service, "crypt-key"); err == nil {
		return k, nil
	}

	k, err := genRandomKey(32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating AES key: %v\n", err)
		panic(err)
	}

	err = keyring.Set(service, "crypt-key", k)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting AES key: %v\n", err)
		panic(err)
	}

	return k, nil
}

func genRandomKey(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func encryptAESGCM(b64key, plaintext string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return "", err
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

	ct := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	out := append(nonce, ct...)
	return base64.StdEncoding.EncodeToString(out), nil
}

// decryptAESGCM decrypts base64(nonce|ciphertext) using a base64 key.
func decryptAESGCM(b64key, b64data string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ct := data[:nonceSize], data[nonceSize:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}
