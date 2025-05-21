package onthewire

import (
	"crypto/rand"
	"crypto/rsa"
)

func asymmetricEncrypt(publicKeyFn func() *rsa.PublicKey) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		logger.Debug("Retrieving public key...")
		publicKey := publicKeyFn()
		logger.Debug("Public key retrieved")

		logger.Debug("Encrypting using public key...")
		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
		if err != nil {
			logger.Error("Failed to encrypt with public key", "Error", err)
			return nil, err
		}

		logger.Debug("Successfully encrypted using public key", "ByteCount", len(encrypted))
		return encrypted, nil
	}
}

func asymmetricDecrypt(privateKeyFn func() *rsa.PrivateKey) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		logger.Debug("Retrieving private key...")
		privateKey := privateKeyFn()
		logger.Debug("Private key retrieved")

		logger.Debug("Decrypting using private key...")
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
		if err != nil {
			logger.Error("Failed to decrypt using private key", "Error", err)
			return nil, err
		}

		logger.Debug("Successfully decrypted using private key", "ByteCount", len(decrypted))
		return decrypted, nil
	}
}
