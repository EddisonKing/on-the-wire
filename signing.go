package onthewire

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func sign(privateKeyFn func() *rsa.PrivateKey) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		logger.Debug("Retrieving private key...")
		privateKey := privateKeyFn()
		logger.Debug("Private key retrieved")

		hash := sha256.Sum256(data)
		logger.Debug("Hashed data", "Hash", hex.EncodeToString(hash[:]))

		logger.Debug("Signing hash...")
		signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
		if err != nil {
			logger.Error("Failed to sign hash", "Error", err)
			return nil, err
		}
		logger.Debug("Signed hash", "Signature", hex.EncodeToString(signature))

		buffer := bytes.NewBuffer(nil)

		if _, err := writeLV(data, buffer); err != nil {
			logger.Error("Failed to write data before signature", "Error", err)
			return nil, err
		}

		if _, err := writeLV(signature, buffer); err != nil {
			logger.Error("Failed to write signature", "Error", err)
			return nil, err
		}

		logger.Debug("Successfully signed data", "ByteCount", len(buffer.Bytes()))
		return buffer.Bytes(), nil
	}
}

func verify(publicKeyFn func() *rsa.PublicKey) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		logger.Debug("Retrieving public key...")
		publicKey := publicKeyFn()
		logger.Debug("Public key retrieved")

		logger.Debug("Verifying data...")
		dataReader := bytes.NewReader(data)

		signedData, _, err := readLV(dataReader)
		if err != nil {
			logger.Error("Failed to read signed data", "Error", err)
			return nil, err
		}

		signature, _, err := readLV(dataReader)
		if err != nil {
			logger.Error("Failed to read signature", "Error", err)
			return nil, err
		}

		hash := sha256.Sum256(signedData)
		logger.Debug("Data hashed", "Hash", hex.EncodeToString(hash[:]))
		valid := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature) == nil

		if valid {
			logger.Debug("Successfully verified signature", "Signature", hex.EncodeToString(signature))
			return signedData, nil
		}

		logger.Error("Failed to verify signature")
		return nil, fmt.Errorf("failed to verify signature")
	}
}
