package onthewire

import (
	"encoding/binary"
	"io"
)

func intToBytes(i int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}

func bytesToInt(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}

func writeLV(data []byte, w io.Writer) (int, error) {
	sizeBytes := intToBytes(len(data))
	sn, err := w.Write(sizeBytes)
	if err != nil {
		return 0, err
	}

	dn, err := w.Write(data)
	if err != nil {
		return 0, nil
	}

	return sn + dn, nil
}

func readLV(r io.Reader) ([]byte, int, error) {
	sizeBytes := make([]byte, 4)
	if _, err := r.Read(sizeBytes); err != nil {
		return nil, 0, err
	}

	dataSize := bytesToInt(sizeBytes)

	if dataSize == 0 {
		return []byte{}, 4, nil
	}

	data := make([]byte, dataSize)
	if _, err := r.Read(data); err != nil {
		return nil, 4, err
	}

	return data, 4 + len(data), nil
}
