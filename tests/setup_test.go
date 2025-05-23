package onthewire_test

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	otw "github.com/EddisonKing/on-the-wire"
)

type TestStruct struct {
	I int
	B bool
	S string
	F float64
}

var someStruct = TestStruct{
	I: 100,
	B: false,
	S: "Hello",
	F: 3.141592654,
}

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

func getKeys() (func() *rsa.PublicKey, func() *rsa.PrivateKey) {
	if privateKey == nil {
		key, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
		if err != nil {
			panic(err)
		}

		privateKey = key
	}

	if publicKey == nil {
		publicKey = &privateKey.PublicKey
	}

	return func() *rsa.PublicKey {
			return publicKey
		}, func() *rsa.PrivateKey {
			return privateKey
		}
}

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	otw.SetLogger(logger)
}

const charSet string = "abcdefghijklmnopqrstuvwxzy0123456789"

func randomString() string {
	res := strings.Builder{}
	for range 20 {
		res.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return res.String()
}

type DelayBuffer struct {
	rw io.ReadWriter
	t  time.Duration
}

func NewDelayBuffer(rw io.ReadWriter, duration time.Duration) *DelayBuffer {
	return &DelayBuffer{
		rw: rw,
		t:  duration,
	}
}

func (sb *DelayBuffer) Read(p []byte) (int, error) {
	time.Sleep(sb.t)
	return sb.rw.Read(p)
}

func (sb *DelayBuffer) Write(d []byte) (int, error) {
	time.Sleep(sb.t)
	return sb.rw.Write(d)
}
