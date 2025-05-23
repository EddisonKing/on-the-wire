package onthewire_test

import (
	"bytes"
	"math/rand"
	"testing"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

func TestSigningPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseSigning(getKeys()).Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someNumber, i)
}

func TestSigningPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseSigning(getKeys()).Build()

	err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someFloat, i)
}

func TestSigningPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseSigning(getKeys()).Build()

	err := write(someBool, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someBool, i)
}

func TestSigningPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseSigning(getKeys()).Build()

	err := write(someStr, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStr, i)
}

func TestSigningPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStruct := TestStruct{
		I: 100,
		B: false,
		S: "Hello",
		F: 3.141592654,
	}

	read, write := otw.New[TestStruct]().UseSigning(getKeys()).Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}

func TestSigningJsonEncodedPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseJSONEncoding().UseSigning(getKeys()).Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someNumber, i)
}

func TestSigningJsonEncodedPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseJSONEncoding().UseSigning(getKeys()).Build()

	err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someFloat, i)
}

func TestSigningJsonEncodedPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseJSONEncoding().UseSigning(getKeys()).Build()

	err := write(someBool, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someBool, i)
}

func TestSigningJsonEncodedPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseJSONEncoding().UseSigning(getKeys()).Build()

	err := write(someStr, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStr, i)
}

func TestSigningJsonEncodedPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStruct := TestStruct{
		I: 100,
		B: false,
		S: "Hello",
		F: 3.141592654,
	}

	read, write := otw.New[TestStruct]().UseJSONEncoding().UseSigning(getKeys()).Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}
