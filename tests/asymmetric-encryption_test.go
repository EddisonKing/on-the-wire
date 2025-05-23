package onthewire_test

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

func TestAsymmetricEncryptionPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someNumber, i)
}

func TestAsymmetricEncryptionPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someFloat, i)
}

func TestAsymmetricEncryptionPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someBool, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someBool, i)
}

func TestAsymmetricEncryptionPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someStr, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someStr, i)
}

func TestAsymmetricEncryptionPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someStruct, i)
}

func TestAsymmetricEncryptionJsonEncodedPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseJSONEncoding().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someNumber, i)
}

func TestAsymmetricEncryptionJsonEncodedPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseJSONEncoding().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someFloat, i)
}

func TestAsymmetricEncryptionJsonEncodedPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseJSONEncoding().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someBool, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someBool, i)
}

func TestAsymmetricEncryptionJsonEncodedPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseJSONEncoding().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someStr, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someStr, i)
}

func TestAsymmetricEncryptionJsonEncodedPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().UseJSONEncoding().UseAsymmetricEncryption(getKeys()).Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someStruct, i)
}

func TestAsymmetricEncryptionLongPayload(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[string]().UseAsymmetricEncryption(getKeys()).Build()

	longPayload := strings.Repeat("A", 2048)

	err := write(longPayload, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, longPayload, i)
}
