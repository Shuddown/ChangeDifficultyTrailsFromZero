package main

import (
	"os"

	zstd "github.com/klauspost/compress/zstd"
)

var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))

// Decompress a buffer. We don't supply a destination buffer,
// so it will be allocated by the decoder.
func Decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}

var encoder, _ = zstd.NewWriter(nil, zstd.WithEncoderCRC(false))

func Compress(src []byte) []byte {
	return encoder.EncodeAll(src, make([]byte, 0, len(src)))
}

func isZstd(filepath string) (bool, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return false, err
	}

	defer f.Close()

	decoder, err := zstd.NewReader(f)
	if err != nil {
		return false, err
	}
	decoder.Close()
	return true, nil
}
