package main

import (
	"crypto/rand"
	//"compress/gzip"
	"errors"
	//"os"
)

const BLOCKSIZE = 4096 // In bytes

func (f *FS) readBlock(n int64) (b []byte, err error) {
	_, err = f.ReadAt(b, n*BLOCKSIZE)
	if err != nil {
		return
	}
	return
}

func (f *FS) writeBlock(b []byte, n int64) (err error) {
	size := BLOCKSIZE - f.ciph.NonceSize() - f.ciph.Overhead()
	if len(b) > size {
		return errors.New("Block too large")
	}
	dst := make([]byte, BLOCKSIZE)
	nonce := dst[:f.ciph.NonceSize()]
	ciphertext := dst[f.ciph.NonceSize():]
	_, err = rand.Read(nonce)
	if err != nil {
		return
	}
	f.ciph.Seal(ciphertext, nonce, b, nil)
	return
}
