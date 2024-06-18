// Package crypto пакет для шифрования и засшифрования данных
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"

	"golang.org/x/crypto/scrypt"
)

// EncryptOAEP получить зашифрованные данные по ключу
func EncryptOAEP(public *rsa.PublicKey, msg []byte) ([]byte, error) {
	msgLen := len(msg)
	hash := sha512.New()
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, rand.Reader, public, msg[start:finish], nil)
		if err != nil {
			return nil, err
		}
		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}
	return encryptedBytes, nil
}

// DecryptOAEP получить расшифрованные данные по ключу
func DecryptOAEP(private *rsa.PrivateKey, msg []byte) ([]byte, error) {
	msgLen := len(msg)
	hash := sha512.New()
	step := private.PublicKey.Size()
	var decryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, rand.Reader, private, msg[start:finish], nil)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}
	return decryptedBytes, nil
}

// GetKey сгенерировать новый ключ используя пароль и соль
func GetKey(password, salt []byte) ([]byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Encrypt получить зашифрованные данные по ключу
func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt получить расшифрованные данные по ключу
func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
