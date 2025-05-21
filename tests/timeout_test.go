package onthewire_test

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

var testTimeout = time.Millisecond * 50

func TestExpectedTimeoutGobEncodedPipelineForInt(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someNumber := rand.Int()

	_, write := otw.New[int]().UseTimeout(testTimeout).Build()

	_, err := write(someNumber, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutGobEncodedPipelineForFloat(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someFloat := rand.Float32()

	_, write := otw.New[float32]().UseTimeout(testTimeout).Build()

	_, err := write(someFloat, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutGobEncodedPipelineForBool(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someBool := !(rand.Float32() > 0.5)

	_, write := otw.New[bool]().UseTimeout(testTimeout).Build()

	_, err := write(someBool, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutGobEncodedPipelineForString(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someStr := randomString()

	_, write := otw.New[string]().UseTimeout(testTimeout).Build()

	_, err := write(someStr, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutGobEncodedPipelineForStruct(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	_, write := otw.New[TestStruct]().UseTimeout(testTimeout).Build()

	_, err := write(someStruct, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutJsonEncodedPipelineForInt(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someNumber := rand.Int()

	_, write := otw.New[int]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	_, err := write(someNumber, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutJsonEncodedPipelineForFloat(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someFloat := rand.Float32()

	_, write := otw.New[float32]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	_, err := write(someFloat, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutJsonEncodedPipelineForBool(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someBool := !(rand.Float32() > 0.5)

	_, write := otw.New[bool]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	_, err := write(someBool, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutJsonEncodedPipelineForString(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	someStr := randomString()

	_, write := otw.New[string]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	_, err := write(someStr, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestExpectedTimeoutJsonEncodedPipelineForStruct(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Second*1)

	_, write := otw.New[TestStruct]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	_, err := write(someStruct, buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrTimedOut, err)
}

func TestSucceedBeforeTimeoutGobEncodedPipelineForInt(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someNumber := rand.Int()

	read, write := otw.New[int]().Build()

	n1, err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someNumber, i)
}

func TestSucceedBeforeTimeoutGobEncodedPipelineForFloat(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().Build()

	n1, err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someFloat, i)
}

func TestSucceedBeforeTimeoutGobEncodedPipelineForBool(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().Build()

	n1, err := write(someBool, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someBool, i)
}

func TestSucceedBeforeTimeoutGobEncodedPipelineForString(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someStr := randomString()

	read, write := otw.New[string]().Build()

	n1, err := write(someStr, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStr, i)
}

func TestSucceedBeforeTimeoutGobEncodedPipelineForStruct(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	read, write := otw.New[TestStruct]().Build()

	n1, err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStruct, i)
}

func TestSucceedBeforeTimeoutJsonEncodedPipelineForInt(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	n1, err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someNumber, i)
}

func TestSucceedBeforeTimeoutJsonEncodedPipelineForFloat(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someFloat := rand.Float32()

	read, write := otw.New[float32]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	n1, err := write(someFloat, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someFloat, i)
}

func TestSucceedBeforeTimeoutJsonEncodedPipelineForBool(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someBool := !(rand.Float32() > 0.5)

	read, write := otw.New[bool]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	n1, err := write(someBool, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someBool, i)
}

func TestSucceedBeforeTimeoutJsonEncodedPipelineForString(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	someStr := randomString()

	read, write := otw.New[string]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	n1, err := write(someStr, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStr, i)
}

func TestSucceedBeforeTimeoutJsonEncodedPipelineForStruct(t *testing.T) {
	buffer := NewDelayBuffer(bytes.NewBuffer(nil), time.Millisecond*10)

	read, write := otw.New[TestStruct]().UseJSONEncoding().UseTimeout(testTimeout).Build()

	n1, err := write(someStruct, buffer)
	assert.Nil(t, err)

	i, n2, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, n1, n2)
	assert.Equal(t, someStruct, i)
}
