package onthewire

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func gobEncode[T any](t T) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	gobber := gob.NewEncoder(buffer)
	if err := gobber.Encode(t); err != nil {
		logger.Error("Failed to Gob encode", "Error", err)
		return nil, err
	}

	logger.Debug("Gob encoded", "Type", reflect.TypeOf(t), "ByteCount", len(buffer.Bytes()), "Bytes", buffer.Bytes())
	return buffer.Bytes(), nil
}

func gobDecode[T any](data []byte) (T, error) {
	gobber := gob.NewDecoder(bytes.NewReader(data))

	t := *new(T)
	if err := gobber.Decode(&t); err != nil {
		logger.Error("Failed to Gob decode", "Error", err)
		return t, err
	}

	logger.Debug("Gob decoded", "Type", reflect.TypeOf(t), "Instance", t)
	return t, nil
}
