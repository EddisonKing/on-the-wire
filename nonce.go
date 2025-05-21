package onthewire

import (
	"bytes"
	"fmt"
)

var ErrNonceInvalid = fmt.Errorf("nonce is invalid")

func checkNonce(check func(int) bool) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		dataReader := bytes.NewReader(data)

		logger.Debug("Reading nonce bytes...")
		nonceBytes, _, err := readLV(dataReader)
		if err != nil {
			logger.Error("Failed to read nonce bytes", "Error", err)
			return nil, err
		}

		nonce := bytesToInt(nonceBytes)
		logger.Debug("Read nonce successfully", "Nonce", nonce)

		logger.Debug("Checking if nonce is valid...")
		if !check(nonce) {
			logger.Error("Failed to validate nonce", "Nonce", nonce)
			return nil, ErrNonceInvalid
		}
		logger.Debug("Nonce is valid")

		remainingData, _, err := readLV(dataReader)
		if err != nil {
			logger.Error("Failed to read remaining data after verifying nonce", "Error", err)
			return nil, err
		}

		return remainingData, nil
	}
}

func setNonce(set func() int) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		nonce := set()
		nonceBytes := intToBytes(nonce)

		buffer := bytes.NewBuffer(nil)

		logger.Debug("Writing nonce", "Nonce", nonce)
		if _, err := writeLV(nonceBytes, buffer); err != nil {
			logger.Error("Failed to write nonce bytes", "Error", err)
			return nil, err
		}

		if _, err := writeLV(data, buffer); err != nil {
			logger.Error("Failed to write data to buffer after nonce", "Error", err)
			return nil, err
		}

		return buffer.Bytes(), nil
	}
}
