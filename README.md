# On The Wire

## Intro
When writing data between services, whether across the network or inter-process communication, it's common to encode and compress data or perform other byte manipulations when transmitting. This library provides a builder pattern for constructing those Read and Write operations. Operations include:
- Gob encoding using Golang's native object encoding
- JSON encoding
- Compression using `compress/zlib` 
- Encryption/decryption using `crypto/rsa`
- Signing/verifying using `crypto/rsa`
- Applying and validating nonces
- Handling timeouts
- Custom []byte -> []byte transforms during reading and writing

## Installation and Import
```bash
go get "github.com/EddisonKing/on-the-wire"
```

```go
import (
  otw "github.com/EddisonKing/on-the-wire"
)
```

## The Pipeline
The core concept is a `Pipeline` (actually a wrapper of a paired `ReadPipeline` and `WritePipeline`) which is in effect just a list of intended actions to perform on data during a write or read. A pipeline does not perform these actions, the `Build()` function of the pipeline produces a function capable of running these instructions.

The simplest `Pipeline` can be defined as:
```go
read, write := otw.New[T].Build()
```

The `New()` function will produce a `Pipeline[T]` where `T` is `any`. It can be any Golang type, whether it be a base type like a string or int, or more complicated types like slices, maps and structs. The `Build()` function produces two functions for reading and writing that implement the desired effects, in the non-functional example above, would just perform `Gob` encoding by default to serialise and deserialise the data.

If the type `T` is not the same for both reading and writing for whatever reason, you can construct two separate pipelines using equivalent methods below but on `ReadPipeline[R]` and `WritePipeline[W]`.

### Encoding/Decoding
Currently, there are two supported encodings, but others are planned. The default is `encoding/gob`, but can be requested explicitly:
```go
read, write := otw.New[T].UseGobEncoding().Build()
```

If JSON is preferred, it can be requested:
```go
read, write := otw.New[T].UseJSONEncoding().Build()
```

### Compression
To enable compression in the pipeline:
```go
read, write := otw.New[T].UseCompression().Build()
```

Compression use the `compress/zlib` library, but future plans will allow this to be switched out if this isn't suitable.

### Encryption/Decryption
Encryption and decryption is performed using `crypto/rsa` and currently only supports RSA for asymmetric encryption/decryption. AES for symmetric encryption is on the roadmap.

```go
pubKeyFn := func() *rsa.PublicKey { ... }
privKeyFn := func() *rsa.PrivateKey { ... }

read, write := otw.New[T].UseAsymmetricEncryption(pubKeyFn, privKeyFn).Build()
```

The `UseAsymmetricEncryption()` function takes in functions that are used to callback to during read and write operations to get the private and public keys respectively. This is more ergonomic since the keys won't get baked into the `read` and `write` functions the pipelines create and the keys are free to change over time.

### Signing/Verification
Like encryption and decryption, the `crypto/rsa` library is used. The function to add signing behaves similar to the encryption and decryption as well since it gets the keys during each `read` and `write` operation and the keys aren't baked into the functions at `Build()` time.

```go
pubKeyFn := func() *rsa.PublicKey { ... }
privKeyFn := func() *rsa.PrivateKey { ... }

read, write := otw.New[T].UseSigning(pubKeyFn, privKeyFn).Build()
```

### Timeouts
```go
read, write := otw.New[T].UseTimeout(time.Duration).Build()
```

Since the underlying `io.Reader` and `io.Writer` that the `read` and `write` functions could act upon may be delayed for whatever reason, it's sometimes necessary to add a timeout.

### Nonces
Nonces are simple `int` that are written and read from the data. This is often used as an additional security measure to prevent replays of messages. `UseNonce()` is provided as a mechanism for the user to hook in a custom nonce check if this additional security is required. 

```go
set := func() int { ... }
check := func(int) bool { ... }

read, write := otw.New[T].UseNonce(set, check).Build()
```

In one scenario, the `set` function would do something like adding the generated nonce to a map, then the `check` function would verify the nonce it read exists in the map. Any other nonce could be suspicious.
In a different scenario, where the sender is not the same as the receiver, the `check` function might just test that it hasn't seen the same nonce twice, thereby protecting against replays.
The usage and relevance of a nonce is dependant on the context obviously, so it's up to the library consumer to determine whether this feature should be used.

### Custom Operations
```go
writeTransformer := func([]byte) ([]byte, error) { ... }
readTransformer := func([]byte) ([]byte, error) { ... }

read, write := otw.New[T].UseCustomOperation(readTransformer, writeTransformer).Build()
```

Custom operations are provided as an additional way to add functionality to the pipeline that doesn't already exist.

### Logging
```go
otw.SetLogger(*slog.Logger)
```

Most of the logging is very verbose in this library so only setting the Level to Debug will produce significant output and shouldn't be necessary unless attempting to debug a custom operation.

## Putting It All Together!
```go
type Test struct {
  someMessage string
  someInt int
  someBool bool
}

pubKeyFn := func() *rsa.PublicKey { ... }
privKeyFn := func() *rsa.PrivateKey { ... }

set := func() int { ... }
check := func(int) bool { ... }

read, write := otw.New[Test].
  UseJSONEncoding().
  UseNonce(set, check).
  UseAsymmetricEncryption(pubKeyFn, privKeyFn).
  UseSigning(pubKeyFn, privKeyFn).
  UseCompression().
  UseTimeout(time.Second * 15).
  Build()
```

In this example, `write` would be a function that takes a `Test` struct and some `io.Writer`, and does the following to it:
1. Converts `Test` to JSON
2. Prepends a nonce to the data, created from the `set` function
3. Encrypts all of those bytes with the public key produced by `pubKeyFn`
4. Compresses the bytes using `compress/zlib`
5. (optional) Times out if write takes too long

And for `read` it does the inverse:
1. Decompresses using `compress/zlib`
2. Decrypts the bytes using the private key produced by `privKeyFn`
3. Strips the nonce off the bytes and uses the `check` function to verify the nonce
4. Converts the bytes, in JSON format, back into the `Test` type
5. (optional) Times out if read takes too long
