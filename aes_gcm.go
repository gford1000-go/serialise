package serialise

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// ErrInvalidDecryptionData is raised if the data to be decrypted is too short for aes-gcm
var ErrInvalidDecryptionData = errors.New("data provided for decryption is too short")

// WithAESGCMEncryption establishes the option to encrypt/decrypt
// data using aes-gcm with the specified key
func WithAESGCMEncryption(key []byte) func(opt *Options) {

	return func(opt *Options) {

		nonceSize := 12

		opt.Encryptor = func(data []byte) ([]byte, error) {

			block, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}

			nonce := make([]byte, nonceSize)
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				return nil, err
			}

			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				return nil, err
			}

			return append(nonce, aesgcm.Seal(nil, nonce, data, nil)...), nil
		}

		opt.Decryptor = func(data []byte) ([]byte, error) {

			if len(data) < nonceSize {
				return nil, ErrInvalidDecryptionData
			}

			nonce := data[0:nonceSize]
			ciphertext := data[nonceSize:]

			block, err := aes.NewCipher(key)
			if err != nil {
				return nil, err
			}

			aesgcm, err := cipher.NewGCM(block)
			if err != nil {
				return nil, err
			}

			return aesgcm.Open(nil, nonce, ciphertext, nil)
		}
	}
}
