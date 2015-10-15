//
// Copyright 2010-2015 David Gray
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package utils

import (
	"crypto/cipher"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block, iv []byte) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

func NewECBEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b, iv))
}

func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	for start,end := 0,x.blockSize ; start < len(src); start,end = start+x.blockSize,end+x.blockSize {
		x.b.Encrypt(dst[start:end],src[start:end])
	}
}

type ecbDecrypter ecb

func NewECBDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b, iv))
}

func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	for start,end := 0,x.blockSize ; start < len(src); start,end = start+x.blockSize,end+x.blockSize {
		x.b.Decrypt(dst[start:end],src[start:end])
	}
}
