package onthewire_test

import (
	"bytes"
	"math/rand"
	"testing"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

func TestGobEncodedPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().Build()

	n1, err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someNumber, i)
}

func TestGobEncodedPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().Build()

	n1, err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someFloat, i)
}

func TestGobEncodedPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().Build()

	n1, err := write(someBool, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someBool, i)
}

func TestGobEncodedPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().Build()

	n1, err := write(someStr, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStr, i)
}

func TestGobEncodedPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().Build()

	n1, err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStruct, i)
}

func TestJsonEncodedPipelineForInt(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseJSONEncoding().Build()

	n1, err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someNumber, i)
}

func TestJsonEncodedPipelineForFloat(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseJSONEncoding().Build()

	n1, err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someFloat, i)
}

func TestJsonEncodedPipelineForBool(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseJSONEncoding().Build()

	n1, err := write(someBool, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someBool, i)
}

func TestJsonEncodedPipelineForString(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someStr := randomString()

	read, write := otw.New[string]().UseJSONEncoding().Build()

	n1, err := write(someStr, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStr, i)
}

func TestJsonEncodedPipelineForStruct(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	read, write := otw.New[TestStruct]().UseJSONEncoding().Build()

	n1, err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStruct, i)
}
