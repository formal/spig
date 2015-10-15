// Copyright 2009 The Go Authors. All rights reserved.
package utils

import (
	"crypto/cipher"
)

func dup(p []byte) []byte {
	q := make([]byte, len(p))
	copy(q, p)
	return q
}

type rcbc struct {
	b         cipher.Block
	blockSize int
	iv        []byte
	tmp       []byte
}

func newRCBC(b cipher.Block, iv []byte) *rcbc {
	return &rcbc{
		b:         b,
		blockSize: b.BlockSize(),
		iv:        dup(iv),
		tmp:       make([]byte, b.BlockSize()),
	}
}

type rcbcEncrypter rcbc

func NewRCBCEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return (*rcbcEncrypter)(newRCBC(b, iv))
}

func (x *rcbcEncrypter) BlockSize() int {
	return x.blockSize
}

func (x *rcbcEncrypter) CryptBlocks(dst, src []byte) {
	for len(src) > 0 {
		for i := 0; i < x.blockSize; i++ {
			x.tmp[i] = x.iv[x.blockSize - i -1] ^ src[i]
		}
		x.b.Encrypt(x.iv, x.tmp)
		for i := 0; i < x.blockSize; i++ {
			dst[i] = x.iv[i]
		}
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type rcbcDecrypter rcbc

func NewRCBCDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return (*rcbcDecrypter)(newRCBC(b, iv))
}

func (x *rcbcDecrypter) BlockSize() int {
	return x.blockSize
}

func (x *rcbcDecrypter) CryptBlocks(dst, src []byte) {
	for len(src) > 0 {
		x.b.Decrypt(x.tmp, src[:x.blockSize])
		for i := 0; i < x.blockSize; i++ {
			x.tmp[i] ^= x.iv[x.blockSize - i -1]
		}
		for i := 0; i < x.blockSize; i++ {
			x.iv[i] = src[i]
			dst[i] = x.tmp[i]
		}

		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
