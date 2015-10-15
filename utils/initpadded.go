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
	"errors"
	_"fmt"
	"crypto/cipher"
	"crypto/aes"
	"golang.org/x/crypto/blowfish"
)
//
var _debug Debug

/*
 * dtg padding
 */

func dtgPad(size int,blockSize int) ([]byte,error) {
	_debug.Printf("DCU Padding for %d bytes",size)
	e := (blockSize - size%blockSize)%blockSize
	l := blockSize+e
	buff := make([]byte,l)
	for i:=0; i<e;i++ {
		buff[i] = 0
	}
	buff[e] = byte(e)
	for i:=e+1; i<len(buff);i++ {
		buff[i] = 0xff
	}
	_debug.PrintBuffer(buff,"Padding = ")
	return buff,nil
}

func dtgRemove(last []byte, blockSize int) (int,error) {
	_debug.PrintBuffer(last,"LAST")
	l := len(last)
	x := int(last[l-blockSize])
	if x >= blockSize {
		return 0,errors.New("1. Invalid DTG Padding")
	}
	if _strict {
		for i := l-x-blockSize; i < l-blockSize; i++ {
			if last[i] != 0 {
				return l,errors.New("2. Invalid DTG Padding")
			}
		}
		for i:=l-blockSize+1; i<l; i++ {
			if last[i] != 0xff {
				return l,errors.New("3. Invalid DTG Padding")
			}
		}
	}
	r := blockSize+x
	_debug.Printf("DCU Padding bytes to remove: %d\n",r)
	return r,nil
}

/*
 * PKCS5/7 Padding
 */

func pkcs5Pad(size int,blockSize int) ([]byte,error) {
	_debug.Printf("PKCS5/PKCS7 Padding for %d bytes",size)
	e := (blockSize - size%blockSize)
	buff := make([]byte,e)
	for i:=0; i<e; i++ {
		buff[i] = byte(e)
	}
	_debug.PrintBuffer(buff,"Padding = ")
	return buff,nil
}

func pkcs5Remove(last []byte, blockSize int) (int,error) {
	_debug.PrintBuffer(last,"LAST")
	x := int(last[len(last)-1])
	if x > blockSize {
		return 0,errors.New("Invalid PKCS5/PKCS7 Padding")
	}
	if _strict {
		for i:=1;i<=x;i++ {
			if last[len(last)-i] != byte(x) {
				return 0,errors.New("Invalid PKCS5/PKCS7 Padding")
			}
		}
	}
	_debug.Printf("PKCS5/PKCS7 Padding bytes to remove: %d\n",x)
	return x,nil
}


/*
 * Zero Padding
 */

func zeroPaddingPad(size int,blockSize int) ([]byte,error) {
	_debug.Printf("Padding for %d bytes",size)
	e := (blockSize - size%blockSize)%blockSize
	buff := make([]byte,e)
	for i:=0; i<e; i++ {
		buff[i] = 0
	}
	return buff,nil
}

func zeroPaddingRemove(last []byte, blockSize int) (int,error) {
	_debug.PrintBuffer(last,"LAST")
	return 0,nil
}




/*
 * init()
 */

func init() {
	_debug = NewDebug(SYSTEM,"formal/padded")
	SetNewCipherFunc("aes",
					func (key []byte) (c cipher.Block,e error) {
						c,e = aes.NewCipher(key)
						return
					})
	SetNewCipherFunc("blowfish",
					func (key []byte) (c cipher.Block,e error) {
						c,e = blowfish.NewCipher(key)
						return
					})

	SetPaddingFuncs("dtgpadding",dtgPad,dtgRemove)
	SetPaddingFuncs("pkcs5padding",pkcs5Pad,pkcs5Remove)
	SetPaddingFuncs("pkcs7padding",pkcs5Pad,pkcs5Remove)
	SetPaddingFuncs("zeropadding",zeroPaddingPad,zeroPaddingRemove)
	SetModeFuncs("ecb",NewECBEncrypter,NewECBDecrypter)
	SetModeFuncs("cbc",cipher.NewCBCEncrypter,cipher.NewCBCDecrypter)
	SetModeFuncs("rcbc",NewRCBCEncrypter,NewRCBCDecrypter)
}
