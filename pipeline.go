package onthewire

import (
	"bytes"
	"crypto/rsa"
	"io"
	"reflect"
	"slices"
	"time"
)

// Represents a pipeline of read/write operations for any type T. The underlying operations act upon byte slices and return byte slices and error if one occured.
type Pipeline[T any] struct {
	readPipeline  *ReadPipeline[T]
	writePipeline *WritePipeline[T]
}

// Represents a pipeline of write operations for any type W. The underlying operations act upon byte slices and return byte slices and error if one occured.
type WritePipeline[W any] struct {
	encoder         func(W) ([]byte, error)
	writeOperations []func([]byte) ([]byte, error)
	useTimeout      bool
	timeoutDuration time.Duration
}

// Creates a new empty pipeline supporting write operations.
//
// This pipeline can be build and used immediately. By default, if no encoder has been selected, the Gob encoder will be used at build time.
func NewWritePipeline[W any]() *WritePipeline[W] {
	return &WritePipeline[W]{
		writeOperations: make([]func([]byte) ([]byte, error), 0),
		encoder:         nil,
		useTimeout:      false,
		timeoutDuration: time.Second,
	}
}

// Represents a pipeline of read operations for any type R. The underlying operations act upon byte slices and return byte slices and error if one occured.
type ReadPipeline[R any] struct {
	readOperations  []func([]byte) ([]byte, error)
	decoder         func([]byte) (R, error)
	useTimeout      bool
	timeoutDuration time.Duration
}

// Creates a new empty pipeline supporting read operations.
//
// This pipeline can be build and used immediately. By default, if no encoder has been selected, the Gob encoder will be used at build time.
func NewReadPipeline[R any]() *ReadPipeline[R] {
	return &ReadPipeline[R]{
		readOperations:  make([]func([]byte) ([]byte, error), 0),
		decoder:         nil,
		useTimeout:      false,
		timeoutDuration: time.Second,
	}
}

// Creates a new empty pipeline with read and write pipelines underpinning it's operation.
//
// This pipeline can be build and used immediately. By default, if no encoder has been selected, the Gob encoder will be used at build time.
func New[T any]() *Pipeline[T] {
	return &Pipeline[T]{
		readPipeline:  NewReadPipeline[T](),
		writePipeline: NewWritePipeline[T](),
	}
}

// Compiles the pipline into a read func and write func following the specification of the pipeline operations and selected encoders.
func (p *Pipeline[T]) Build() (func(io.Reader) (T, error), func(T, io.Writer) error) {
	return p.readPipeline.Build(), p.writePipeline.Build()
}

// Custom Operations allow the consumer to define their own write and read operations to append to the pipeline.
//
// Both functions take in a []byte and output a modified []byte or error. It is ok to return a new []byte or modify the existing one.
func (p *Pipeline[T]) UseCustomOperation(readFn, writeFn func([]byte) ([]byte, error)) *Pipeline[T] {
	p.readPipeline.UseCustomOperation(readFn)
	p.writePipeline.UseCustomOperation(writeFn)
	return p
}

// Appends a compression step to the write operations and a decompression strep to the read operations.//
//
// Compression is done with the `compress/zlib` library.
func (p *Pipeline[T]) UseCompression() *Pipeline[T] {
	p.writePipeline.UseCompression()
	p.readPipeline.UseCompression()
	return p
}

// Enables on-boarding and off-boarding to the pipeline using Go's native Go Object Encoding.
//
// Gob Encoding requires that structs export their fields to be transmitted. No exported fields will result in an error on write.
func (p *Pipeline[T]) UseGobEncoding() *Pipeline[T] {
	p.readPipeline.UseGobEncoding()
	p.writePipeline.UseGobEncoding()
	return p
}

// Enables on-boarding and off-boarding to the pipeline using JSON Encoding.
func (p *Pipeline[T]) UseJSONEncoding() *Pipeline[T] {
	p.readPipeline.UseJSONEncoding()
	p.writePipeline.UseJSONEncoding()
	return p
}

// Use RSA asymmetric encryption for encrypting and decrypting data.
//
// It is up to the consumer of the library to provide callback functions that return the public and private keys. The functions will only be used during read and write operations, not during the building of the pipeline.
func (p *Pipeline[T]) UseAsymmetricEncryption(publicKeyFn func() *rsa.PublicKey, privateKeyFn func() *rsa.PrivateKey) *Pipeline[T] {
	p.readPipeline.UseAsymmetricEncryption(privateKeyFn)
	p.writePipeline.UseAsymmetricEncryption(publicKeyFn)
	return p
}

// Use RSA asymmetric encryption for signing and verifying data being sent. The write operation to the pipeline appends a []byte containing the signature. The read operation will cut the signature from the pipeline, verify it, and either continue processing or error if the signature fails to validate.
//
// It is up to the consumer of the library to provide callback functions that return the public and private keys. The functions will only be used during read and write operations, not during the building of the pipeline.
func (p *Pipeline[T]) UseSigning(publicKeyFn func() *rsa.PublicKey, privateKeyFn func() *rsa.PrivateKey) *Pipeline[T] {
	p.readPipeline.UseSigning(publicKeyFn)
	p.writePipeline.UseSigning(privateKeyFn)
	return p
}

// Enables a timeout on all operations. This timeout is used for each operation, as well as encoding and decoding the initial and final payloads.
func (p *Pipeline[T]) UseTimeout(t time.Duration) *Pipeline[T] {
	p.readPipeline.UseTimeout(t)
	p.writePipeline.UseTimeout(t)
	return p
}

// Enables the use of integer nonces in read and write operations. The provided nonce set and check functions are callbacks for how to generate the nonce for writing and how to check them upon reading.
func (p *Pipeline[T]) UseNonce(set func() int, check func(int) bool) *Pipeline[T] {
	p.readPipeline.UseNonce(check)
	p.writePipeline.UseNonce(set)
	return p
}

// Compiles the pipline into a read func following the specification of the pipeline operations and selected encoders.
func (p *ReadPipeline[R]) Build() func(io.Reader) (R, error) {
	logger.Info("Building read pipeline", "NumberOfReadOperations", len(p.readOperations))
	if p.decoder == nil {
		logger.Warn("No encoding selected, defaulting to Gob Encoding")
		p.decoder = gobDecode
	}

	slices.Reverse(p.readOperations)

	logger.Debug("Building Read function")

	rlv := conditionalAddTimeoutReader(p.useTimeout, readLV, p.timeoutDuration)

	readFn := func(r io.Reader) (R, error) {
		t := *new(R)

		buffer := bytes.NewBuffer(nil)

		logger.Debug("Beginning to read chunks...")
		for {
			bufferSection, n, err := rlv(r)
			if err != nil {
				logger.Error("Failed to read chunk", "Error", err)
				return t, err
			}

			if len(bufferSection) == 0 {
				break
			}

			buffer.Write(bufferSection)
			logger.Debug("Read chunk bytes", "ByteCount", n)
		}

		logger.Debug("Beginning read operations")
		data := buffer.Bytes()
		for _, operation := range p.readOperations {
			d, err := operation(data)
			if err != nil {
				logger.Error("Failed to complete read pipeline. An operation failed", "Error", err)
				return t, err
			}
			data = d
		}

		t, err := p.decoder(data)
		if err != nil {
			logger.Error("Failed to decode final bytes as required type", "Error", err)
			return t, err
		}

		logger.Debug("Completed reading")
		return t, nil
	}

	return readFn
}

// Custom Operations allow the consumer to define their own read operations to append to the pipeline.
//
// Functions take in a []byte and output a modified []byte or error. It is ok to return a new []byte or modify the existing one.
func (p *ReadPipeline[R]) UseCustomOperation(readFn func([]byte) ([]byte, error)) *ReadPipeline[R] {
	p.readOperations = append(p.readOperations, conditionalAddTimeout(p.useTimeout, readFn, p.timeoutDuration))
	return p
}

// Appends a decompression step to the read operations.
//
// Compression is done with the `compress/zlib` library.
func (p *ReadPipeline[R]) UseCompression() *ReadPipeline[R] {
	p.readOperations = append(p.readOperations, conditionalAddTimeout(p.useTimeout, decompress, p.timeoutDuration))
	return p
}

// Enables off-boarding from the pipeline using Go's native Go Object Encoding.
//
// Gob Encoding requires that structs export their fields to be transmitted. No exported fields will result in an error on write.
func (p *ReadPipeline[R]) UseGobEncoding() *ReadPipeline[R] {
	p.decoder = gobDecode
	return p
}

// Enables off-boarding from the pipeline using JSON Encoding.
func (p *ReadPipeline[R]) UseJSONEncoding() *ReadPipeline[R] {
	p.decoder = jsonDecode
	return p
}

// Use RSA asymmetric encryption for decrypting data.
//
// It is up to the consumer of the library to provide a callback function to return the private key. The functions will only be used during read and write operations, not during the building of the pipeline.
func (p *ReadPipeline[R]) UseAsymmetricEncryption(privateKeyFn func() *rsa.PrivateKey) *ReadPipeline[R] {
	p.readOperations = append(p.readOperations, conditionalAddTimeout(p.useTimeout, asymmetricDecrypt(privateKeyFn), p.timeoutDuration))
	return p
}

// Use RSA asymmetric encryption for verifying data being sent. The read operation will cut the signature from the pipeline, verify it, and either continue processing or error if the signature fails to validate.
//
// It is up to the consumer of the library to provide callback functions that return the public key. The functions will only be used during read and write operations, not during the building of the pipeline.
func (p *ReadPipeline[R]) UseSigning(publicKeyFn func() *rsa.PublicKey) *ReadPipeline[R] {
	p.readOperations = append(p.readOperations, conditionalAddTimeout(p.useTimeout, verify(publicKeyFn), p.timeoutDuration))
	return p
}

// Enables a timeout on all operations. This timeout is used for each operation, as well as decoding final payloads.
func (p *ReadPipeline[R]) UseTimeout(t time.Duration) *ReadPipeline[R] {
	p.useTimeout = true
	p.timeoutDuration = t
	return p
}

// Enables nonces during read operations. An integer nonce is checked for validity using the check callback during reading
func (p *ReadPipeline[R]) UseNonce(check func(int) bool) *ReadPipeline[R] {
	p.readOperations = append(p.readOperations, checkNonce(check))
	return p
}

// Compiles the pipline into a read func and write func following the specification of the pipeline operations and selected encoders.
func (p *WritePipeline[W]) Build() func(W, io.Writer) error {
	logger.Info("Building write pipeline", "NumberOfWriteOperations", len(p.writeOperations))
	if p.encoder == nil {
		logger.Warn("No encoding selected, defaulting to Gob Encoding")
		p.encoder = gobEncode
	}

	wlv := conditionalAddTimeoutWriter(p.useTimeout, writeLV, p.timeoutDuration)

	logger.Debug("Building Write function")
	writeFn := func(t W, w io.Writer) error {
		encoded, err := p.encoder(t)
		if err != nil {
			logger.Error("Failed to onboard data into write pipeline", "Error", err, "Type", reflect.TypeOf(t))
			return err
		}

		logger.Debug("Beginning write operations...")
		data := encoded
		for _, operation := range p.writeOperations {
			data, err = operation(data)
			if err != nil {
				logger.Error("Failed to complete write pipeline. An operation failed")
				return err
			}
		}

		logger.Debug("Beginning chunked writes...")
		for i := 0; i < len(data); i += 1024 {
			start := i
			end := i + 1024

			written := 0
			if end >= len(data) {
				logger.Debug("Writing last chunk...")
				written, err = wlv(data[start:], w)
			} else {
				logger.Debug("Writing chunk...")
				written, err = wlv(data[start:end], w)
			}

			if err != nil {
				logger.Error("Failed to write chunk", "Error", err)
				return err
			}
			logger.Debug("Written bytes", "ByteCount", written)
		}

		// Write empty buffer to indicate to stop buffering
		if _, err := wlv([]byte{}, w); err != nil {
			logger.Error("Failed to write stop chunk", "Error", err)
			return err
		}

		logger.Debug("Completed writing")
		return nil
	}

	return writeFn
}

// Custom Operations allow the consumer to define their own write operations to append to the pipeline.
//
// Functions take in a []byte and output a modified []byte or error. It is ok to return a new []byte or modify the existing one.
func (p *WritePipeline[W]) UseCustomOperation(writeFn func([]byte) ([]byte, error)) *WritePipeline[W] {
	p.writeOperations = append(p.writeOperations, conditionalAddTimeout(p.useTimeout, writeFn, p.timeoutDuration))
	return p
}

// Appends a compression step to the write operations.
//
// Compression is done with the `compress/zlib` library.
func (p *WritePipeline[W]) UseCompression() *WritePipeline[W] {
	p.writeOperations = append(p.writeOperations, conditionalAddTimeout(p.useTimeout, compress, p.timeoutDuration))
	return p
}

// Enables on-boarding to the pipeline using Go's native Go Object Encoding.
//
// Gob Encoding requires that structs export their fields to be transmitted. No exported fields will result in an error on write.
func (p *WritePipeline[W]) UseGobEncoding() *WritePipeline[W] {
	p.encoder = gobEncode
	return p
}

// Enables on-boarding to the pipeline using JSON Encoding.
func (p *WritePipeline[W]) UseJSONEncoding() *WritePipeline[W] {
	p.encoder = jsonEncode
	return p
}

// Use RSA asymmetric encryption for encrypting data.
//
// It is up to the consumer of the library to provide a callback function to return the public key. The function will only be used during write operations, not during the building of the pipeline.
func (p *WritePipeline[W]) UseAsymmetricEncryption(publicKeyFn func() *rsa.PublicKey) *WritePipeline[W] {
	p.writeOperations = append(p.writeOperations, conditionalAddTimeout(p.useTimeout, asymmetricEncrypt(publicKeyFn), p.timeoutDuration))
	return p
}

// Use RSA asymmetric encryption for signing the data being sent. The write operation to the pipeline appends a []byte containing the signature.
//
// It is up to the consumer of the library to provide the callback function to return the private key. The function will only be used during write operations, not during the building of the pipeline.
func (p *WritePipeline[W]) UseSigning(privateKeyFn func() *rsa.PrivateKey) *WritePipeline[W] {
	p.writeOperations = append(p.writeOperations, conditionalAddTimeout(p.useTimeout, sign(privateKeyFn), p.timeoutDuration))
	return p
}

// Enables a timeout on all operations. This timeout is used for each operation, as well as encoding the initial payload.
func (p *WritePipeline[W]) UseTimeout(t time.Duration) *WritePipeline[W] {
	p.useTimeout = true
	p.timeoutDuration = t
	return p
}

// Enables nonces during read operations. An integer nonce is checked for validity using the check callback during reading
func (p *WritePipeline[W]) UseNonce(set func() int) *WritePipeline[W] {
	p.writeOperations = append(p.writeOperations, setNonce(set))
	return p
}
