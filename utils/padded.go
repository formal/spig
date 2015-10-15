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
	"strings"
	"crypto/cipher"
)


var _strict = true

func SetStrict(v bool) {
	_strict = v
}

/*
 * Functions for creating new cipher objects
 */

type newCipherFunc func (key []byte) (cipher.Block,error)

var newCiphers map[string] newCipherFunc = make(map[string] newCipherFunc,10)

func SetNewCipherFunc(name string,fun newCipherFunc) {
	newCiphers[name] = fun
}

func NewCipher(alg string, key []byte) (cipher.Block,error) {
	f,ok := newCiphers[alg]
	if !ok {
		s := "Unsupported algorithm: " + alg
		return nil,errors.New(s)
	}
	return f(key)
}

/*
 * Padding & removal functions
 */

type padFunc 	func(size int,blockSize int) ([]byte,error)
type removeFunc func(last []byte,blockSize int) (int,error)

type padding struct {
	pad 	padFunc
	remove 	removeFunc
}

var paddings 	map[string] padding 	= make(map[string] padding,10)

func SetPaddingFuncs(name string,pad padFunc,remove removeFunc) {
	paddings[name] = padding{pad,remove}
}

func GetPadding(pad string) (*padding,error) {
	p,ok := paddings[pad]
	if !ok {
		s := "Unsupported padding: " + pad
		return nil,errors.New(s)
	}
	return &p,nil
}


/*
 * BlockMode functions
 */

type ModeParams interface {
}

type newModeEncFunc	func(cipher cipher.Block,iv []byte) cipher.BlockMode
type newModeDecFunc func(cipher cipher.Block,iv []byte) cipher.BlockMode

var modeEncs map[string] newModeEncFunc		= make(map[string] newModeEncFunc,10)
var modeDecs map[string] newModeDecFunc		= make(map[string] newModeDecFunc,10)

func SetModeFuncs(name string,enc newModeEncFunc,dec newModeDecFunc) {
	modeEncs[name] = enc
	modeDecs[name] = dec
}

func newEncBlockMode(mode string,iv []byte,cipher cipher.Block) (cipher.BlockMode,error) {
	encFunc,ok := modeEncs[mode]
	if !ok {
		s := "Unsupported mode: " + mode
		return nil,errors.New(s)
	}
	return encFunc(cipher,iv),nil
}

func newDecBlockMode(mode string,iv []byte,cipher cipher.Block) (cipher.BlockMode,error) {
	decFunc,ok := modeDecs[mode]
	if !ok {
		s := "Unsupported mode: " + mode
		return nil,errors.New(s)
	}
	return decFunc(cipher,iv),nil
}

func parse(name string) (alg, mode, padding string,err error) {
	l := strings.Split(name,"/")
	if len(l) != 3 {
		s := "Invalid name: " + name
		return "","","",errors.New(s)
	}
	return 	strings.ToLower(l[0]),
			strings.ToLower(l[1]),
			strings.ToLower(l[2]),
			nil
}



// Note that plaintext is overwritten
func Encrypt(name string,iv []byte,key []byte,plaintext []byte) (ciphertext []byte,err error) {
	_debug.PrintBuffer(iv,"Encrypt IV")
	_debug.PrintBuffer(key,"Encrypt Key")
	_debug.PrintBuffer(plaintext,"Plaintext")
	alg, mode, pad, err := parse(name)
	if err != nil {
		return
	}

    p,err := GetPadding(pad)
	if err != nil {
			return
	}

	b,err := NewCipher(alg,key)
	if err != nil {
			return
	}

	bm, err := newEncBlockMode(mode,iv,b)
	if err != nil {
			return
	}

	// Add padding
	// Do not assuming that plaintext has enough room to add padding
	l := len(plaintext)
	padding,err := p.pad(l,b.BlockSize())
	if err != nil {
			return
	}
	ciphertext = make([]byte,l+len(padding))
	for i := 0; i<l; i++ {
			ciphertext[i] = plaintext[i]
	}
	for i := 0; i < len(padding); i++ {
		ciphertext[l+i] = padding[i]
	}
	bm.CryptBlocks(ciphertext,ciphertext);
	if err != nil {
			return
	}
	_debug.PrintBuffer(ciphertext,"Ciphertext")
	return
}


// Note that ciphertext is overwritten
func Decrypt(name string,iv []byte,key []byte,ciphertext []byte) (plaintext []byte,err error)  {
	_debug.PrintBuffer(iv,"Decrypt IV")
	_debug.PrintBuffer(key,"Decrypt Key")
	_debug.PrintBuffer(ciphertext,"Ciphertext")
	alg, mode, pad, err := parse(name)
	if err != nil {
		return
	}

	p,err := GetPadding(pad)
	if err != nil {
			return
	}

	b,err := NewCipher(alg,key)
	if err != nil {
			return
	}

	if len(ciphertext)%b.BlockSize() != 0 {
		err = errors.New("Ciphertext not a mutiple of block size")
		return
	}

	bm, err := newDecBlockMode(mode,iv,b)
	if err != nil {
			return
	}

	bm.CryptBlocks(ciphertext,ciphertext);
	if err != nil {
			return
	}

	// Remove padding
	l,err := p.remove(ciphertext,b.BlockSize())
	if err != nil {
			return
	}
	plaintext = ciphertext[:len(ciphertext)-l]
	_debug.PrintBuffer(plaintext,"Plaintext")
	return
}
