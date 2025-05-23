package onthewire

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
)

func asymmetricEncrypt(publicKeyFn func() *rsa.PublicKey) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		logger.Debug("Retrieving public key...")
		publicKey := publicKeyFn()
		logger.Debug("Public key retrieved")

		buffer := bytes.NewBuffer(nil)
		logger.Debug("Encrypting using public key...")

		for i := 0; i < len(data); i += 117 {
			start := i
			end := i + 117

			var encrypted []byte
			var err error
			if end > len(data) {
				encrypted, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, data[start:])
			} else {
				encrypted, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, data[start:end])
			}

			if err != nil {
				logger.Error("Failed to encrypt with public key", "Error", err)
				return nil, err
			}

			writeLV(encrypted, buffer)
		}

		writeLV([]byte{}, buffer)

		logger.Debug("Successfully encrypted using public key", "ByteCount", len(buffer.Bytes()))
		return buffer.Bytes(), nil
	}
}

func asymmetricDecrypt(privateKeyFn func() *rsa.PrivateKey) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		logger.Debug("Retrieving private key...")
		privateKey := privateKeyFn()
		logger.Debug("Private key retrieved")

		dataReader := bytes.NewReader(data)
		buffer := bytes.NewBuffer(nil)
		logger.Debug("Decrypting using private key...", "ByteCount", len(data))

		for {
			chunk, _, err := readLV(dataReader)
			if err != nil {
				logger.Error("Failed to read encrypted chunk", "Error", err)
				return nil, err
			}

			if len(chunk) == 0 {
				break
			}

			decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, chunk)
			if err != nil {
				logger.Error("Failed to decrypt using private key", "Error", err)
				return nil, err
			}

			buffer.Write(decrypted)
		}

		logger.Debug("Successfully decrypted using private key", "ByteCount", len(buffer.Bytes()))
		return buffer.Bytes(), nil
	}
}
