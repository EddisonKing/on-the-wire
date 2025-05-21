package onthewire

import (
	"bytes"
	"encoding/json"
	"reflect"
)

func jsonEncode[T any](t T) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	jsonifer := json.NewEncoder(buffer)
	if err := jsonifer.Encode(t); err != nil {
		logger.Error("Failed to JSON encode", "Error", err)
		return nil, err
	}

	logger.Debug("JSON encoded", "Type", reflect.TypeOf(t), "ByteCount", len(buffer.Bytes()), "Bytes", buffer.Bytes())
	return buffer.Bytes(), nil
}

func jsonDecode[T any](data []byte) (T, error) {
	jsonifier := json.NewDecoder(bytes.NewReader(data))

	t := *new(T)
	if err := jsonifier.Decode(&t); err != nil {
		logger.Debug("Failed to JSON decode", "Error", err)
		return t, err
	}

	logger.Debug("JSON decoded", "Type", reflect.TypeOf(t), "Instance", t)
	return t, nil
}
