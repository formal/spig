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

package aeskey

import (
	"strings"
	"strconv"
	"crypto/rand"
)

var keyString = "4f 01 b6 81 72 1c e4 e1 c4 17 bc 5d 82 41 f2 35"
var ivString =  "41 78 b6 00 9e ff c1 f4 37 34 67 23 2d dd 33 f6"

// NSSK / Otway Rees Keys & IVs
var keyAString = "ff 81 cd 46 a9 a9 7b b6 38 9c 7a ce 7c 6b cf 75"
var keyBString = "45 6b 1d 8e 81 2a f3 3c 60 f1 4b 31 45 21 fc db"
var ivAString =  "41 78 b6 45 9e ff c1 f4 37 37 67 77 2d 7d 33 f6"
var ivBString =  "41 23 b6 00 9e ff c1 f4 67 34 67 23 2d dd 11 f6"

func makeBytes(s string) ([]byte,error) {
	list := strings.Split(s," ")
	b := make([]byte,len(list))
	for i := range(list) {
		x, err := strconv.ParseUint(list[i],16,8)
		b[i] = byte(x)
		if err != nil {
			return nil,err
		}
	}
	return b,nil
}

func convert(val byte) byte {
	if val < 10 {
		return byte(val + '0')
	}
	return byte(val + 'a' - 10)
}

func MakeString(key []byte) string {
	str := make([]byte,3*16)
	for i,b := range(key) {
		str[3*i] = convert(b >> 4)
		str[3*i+1] = convert(b % 16)
		str[3*i+2] = ' '
	}
	return string(str[0:len(str)-1])
}

func SessionKey() ([]byte,error) {
	key := make([]byte,16)
	_,err := rand.Read(key)
	return key,err
}

func Key() ([]byte,error) {
	return makeBytes(keyString)
}

func KeyA() ([]byte,error) {
	return makeBytes(keyAString)
}

func KeyB() ([]byte,error) {
	return makeBytes(keyBString)
}

func Iv() ([]byte,error) {
	return makeBytes(ivString)
}

func IvA() ([]byte,error) {
	return makeBytes(ivAString)
}

func IvB() ([]byte,error) {
	return makeBytes(ivBString)
}

func keyForString(key string) (s string) {
	for _,c := range key {
		if c != ' ' {
			s += string(c)
		}
	}
	return
}

func KeyString() (s string) {
	return keyForString(keyString)
}

func KeyAString() (s string) {
	return keyForString(keyAString)
}

func KeyBString() (s string) {
	return keyForString(keyBString)
}

func IvString() (s string) {
	return keyForString(ivString)
}

func IvAString() (s string) {
	return keyForString(ivAString)
}

func IvBString() (s string) {
	return keyForString(ivBString)
}
