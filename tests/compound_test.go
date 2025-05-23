package onthewire_test

import (
	"bytes"
	"testing"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

func TestAsymmetricEncryptionThenCompressionGobEncoding(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().
		UseTimeout(testTimeout).
		UseGobEncoding().
		UseAsymmetricEncryption(getKeys()).
		UseSigning(getKeys()).
		UseCompression().
		Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}

func TestCompressionThenAsymmetricEncryptionGobEncoding(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().
		UseTimeout(testTimeout).
		UseGobEncoding().
		UseCompression().
		UseAsymmetricEncryption(getKeys()).
		UseSigning(getKeys()).
		Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}

func TestAsymmetricEncryptionThenCompressionJSONEncoding(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().
		UseTimeout(testTimeout).
		UseJSONEncoding().
		UseAsymmetricEncryption(getKeys()).
		UseSigning(getKeys()).
		UseCompression().
		Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}

func TestCompressionThenAsymmetricEncryptionJSONEncoding(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().
		UseTimeout(testTimeout).
		UseJSONEncoding().
		UseCompression().
		UseAsymmetricEncryption(getKeys()).
		UseSigning(getKeys()).
		Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}
