package onthewire_test

import (
	"bytes"
	"math/rand"
	"testing"

	otw "github.com/EddisonKing/on-the-wire"
	"github.com/stretchr/testify/assert"
)

func TestValidNonceShouldPass(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseNonce(func() int { return 1 }, func(i int) bool { return i == 1 }).Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	i, err := read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, someNumber, i)
}

func TestInvalidNonceShouldFail(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	someNumber := rand.Int()

	read, write := otw.New[int]().UseNonce(func() int { return 1 }, func(i int) bool { return i == 2 }).Build()

	err := write(someNumber, buffer)
	assert.Nil(t, err)

	_, err = read(buffer)
	assert.NotNil(t, err)
	assert.Equal(t, otw.ErrNonceInvalid, err)
}
