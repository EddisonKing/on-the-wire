package onthewire

import (
	"bytes"
	"compress/zlib"
	"io"
)

func compress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	uncompressedLen := len(data)

	logger.Debug("Compressing...", "ByteCount", uncompressedLen)
	compressor := zlib.NewWriter(buffer)
	if _, err := compressor.Write(data); err != nil {
		logger.Error("Failed to compress data", "Error", err)
		return nil, err
	}
	compressor.Close()

	compressedLen := len(buffer.Bytes())

	logger.Debug("Successfully compressed", "ByteCount", compressedLen)

	if compressedLen > uncompressedLen {
		logger.Warn("Compression resulted in an increase in total byte count. Consider that the data you are transferring is already small enough to avoid compression", "UncompressedByteCount", uncompressedLen, "CompressedByteCount", compressedLen)
	}

	return buffer.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	logger.Debug("Decompressing...")
	byteReader := bytes.NewReader(data)
	decompressor, err := zlib.NewReader(byteReader)
	if err != nil {
		logger.Error("Failed to decompress data", "Error", err)
		return nil, err
	}
	defer decompressor.Close()

	if _, err := io.Copy(buffer, decompressor); err != nil {
		logger.Error("Failed to copy decompressed data", "Error", err)
		return nil, err
	}

	logger.Debug("Successfully decompressed", "ByteCount", len(buffer.Bytes()))
	return buffer.Bytes(), nil
}
