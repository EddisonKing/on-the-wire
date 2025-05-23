package onthewire_test

import (
	"bytes"
	"math/rand"
	"testing"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

func TestCompressionPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseCompression().Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someNumber, i)
}

func TestCompressionPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseCompression().Build()

	err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someFloat, i)
}

func TestCompressionPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseCompression().Build()

	err := write(someBool, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someBool, i)
}

func TestCompressionPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseCompression().Build()

	err := write(someStr, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStr, i)
}

func TestCompressionPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().UseCompression().Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}

func TestCompressionJsonEncodedPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseJSONEncoding().UseCompression().Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someNumber, i)
}

func TestCompressionJsonEncodedPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseJSONEncoding().UseCompression().Build()

	err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someFloat, i)
}

func TestCompressionJsonEncodedPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseJSONEncoding().UseCompression().Build()

	err := write(someBool, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someBool, i)
}

func TestCompressionJsonEncodedPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseJSONEncoding().UseCompression().Build()

	err := write(someStr, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStr, i)
}

func TestCompressionJsonEncodedPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().UseJSONEncoding().UseCompression().Build()

	err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)

	assert.Equal(t, someStruct, i)
}
