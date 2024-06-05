package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	authDevivery "ya-gophkeeper-client/internal/auth/delivery/http"
	authUseCase "ya-gophkeeper-client/internal/auth/usecase"

	storeDevivery "ya-gophkeeper-client/internal/store/delivery/http"
	storeUseCase "ya-gophkeeper-client/internal/store/usecase"

	"ya-gophkeeper-client/internal/auth/entity"
	"ya-gophkeeper-client/internal/config"
	"ya-gophkeeper-client/internal/tui"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

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

func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

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

func Start() {
	conf := config.New()
	authD := authDevivery.New(conf)
	authUC := authUseCase.New(authD)
	storeD := storeDevivery.New(conf)
	storeUC := storeUseCase.New(storeD)
	token := ""

	tui := tui.New()

	tui.Register("auth", func() {
		email, err := tui.GetUserPrint("Email: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		pass, err := tui.GetHideUserPrint("Pass: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		token, err = authUC.Login(&entity.LoginDto{
			Email:    email,
			Password: pass,
		})
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("token ", token)
	})

	tui.Register("save", func() {
		message, err := tui.GetUserPrint("Message: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		key, _, err := DeriveKey([]byte("pass"), []byte("salt"))
		if err != nil {
			logrus.Error(err)
			return
		}
		data, err := Encrypt(key, []byte(message))
		if err != nil {
			logrus.Error(err)
			return
		}
		err = storeUC.Save(data)
		if err != nil {
			logrus.Error(err)
			return
		}
	})

	tui.Register("get", func() {
		id, err := tui.GetUserPrint("ID: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		data, err := storeUC.GetById(id)
		if err != nil {
			logrus.Error(err)
			return
		}
		key, _, err := DeriveKey([]byte("pass"), []byte("salt"))
		if err != nil {
			logrus.Error(err)
			return
		}
		data2, err := Decrypt(key, data.Data)
		if err != nil {
			logrus.Error(err)
			return
		}
		tui.Println(string(data2))
	})

	tui.Register("gen", func() {
		pass, err := tui.GetUserPrint("Pass: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		salt, err := tui.GetUserPrint("Salt: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		message, err := tui.GetUserPrint("Message: ")
		if err != nil {
			logrus.Error(err)
			return
		}
		key, _, err := DeriveKey([]byte(pass), []byte(salt))
		if err != nil {
			logrus.Error(err)
			return
		}
		data, err := Encrypt(key, []byte(message))
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("Encrypt: ", data)
		data2, err := Decrypt(key, data)
		if err != nil {
			logrus.Error(err)
			return
		}
		logrus.Info("Decrypt: ", string(data2))
	})

	tui.Register("crypt", func() {
		coder := base64.StdEncoding
		secret := sha256.Sum256([]byte("hello"))
		s := coder.EncodeToString(secret[:])
		fmt.Println(s)
	})

	tui.Run()
}
