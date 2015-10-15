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
	"fmt"
	"bytes"
	"time"
	"math/big"
)

//
// Input encoder

type ITlcEncoding struct {
	source IStream
}

func MakeITlcEncoding(stream IStream) *ITlcEncoding {
	enc := ITlcEncoding{stream}
	return &enc
}

func (enc *ITlcEncoding) readLength() (int,error) {
	b1,e := enc.source.ReadByte()
	if e != nil {
		return 0,e
	}
	b2,e :=enc.source.ReadByte()
	if e != nil {
		return 0,e
	}
	return 256*int(b2) + int(b1),nil
}

func (enc *ITlcEncoding) _readBuffer() (*bytes.Buffer,error) {
	l,e := enc.readLength()
	if e != nil {
		return nil,e
	}

	buff := bytes.NewBuffer(make([]byte,l))
	buff.Reset()
	for i:= 0; i<l; i++ {
		v,e := enc.source.ReadByte()
		if e != nil {
			return nil,e
		}
		buff.WriteByte(v)
	}
	return buff,nil
}

func (enc *ITlcEncoding) checkType(t int) error {
	b,e := enc.source.ReadByte()
	if e != nil {
		return e
	}
	if int(b) != t {
		s := fmt.Sprintf("Invalid type %d in stream; expected %d",b,t)
		return errors.New(s)
	}
	return nil
}

func (enc *ITlcEncoding) readBuffer(t int) (*bytes.Buffer,error) {
	e := enc.checkType(t)
	if e != nil {
		return nil,e
	}
	return enc._readBuffer()
}

func (enc *ITlcEncoding) ReadUint64() (uint64,error) {
	buff,e := enc.readBuffer(INTEGER)
	if e != nil {
		return 0,e
	}

	s := buff.String()
	var x uint64 = 0
	for i:=0; i<len(s); i++ {
		if ((s[i] > '9') || (s[i] < '0')) {
			return 0,errors.New("Invalid character in integer encoding")
		}
		var d uint64 = uint64(s[i] - '0')
		if ((x > MAX_UINT64_DIV_10) || ((x == MAX_UINT64_DIV_10) && (d > MAX_UINT64_MOD_10))) {
			return 0,errors.New("Integer string value too large")
		}
		x = 10*x + d
	}
	return x,nil
}

func (enc *ITlcEncoding) ReadString() (string,error) {
	buff,e := enc.readBuffer(STRING)
	if e != nil {
		return "",e
	}

	return buff.String(),nil
}

func (enc *ITlcEncoding) ReadInteger() (string,error) {
	buff,e := enc.readBuffer(INTEGER)
	if e != nil {
		return "",e
	}

	return buff.String(),nil
}

func (enc *ITlcEncoding) ReadBig() (*big.Int,error) {
	x := new(big.Int)
	buff,e := enc.readBuffer(INTEGER)
	if e != nil {
	   return x,e
	}
	x,b := x.SetString(buff.String(),10)
	if b {
		return x,nil
	}
	return x,errors.New("Problem reading big integer")
}

func (enc *ITlcEncoding) ReadBinary() (b []byte,err error) {
	buff,e := enc.readBuffer(BINARY)
	if e != nil {
		return nil,e
	}

	return buff.Bytes(),e
}

func (enc *ITlcEncoding) ReadStructured() ([]byte,error) {
	buff,e := enc.readBuffer(STRUCTURED)
	if e != nil {
		return nil,e
	}

	return buff.Bytes(),e
}


//
// Output encoder
//

type OTlcEncoding struct {
	destination OStream
}


func MakeOTlcEncoding(stream OStream) *OTlcEncoding {
	enc := OTlcEncoding{stream}
	return &enc
}


func (enc *OTlcEncoding) writeBuffer(t int,buff []byte) error { // we could make this a degenerate case of writeBufferWithDelay
	e := enc.destination.WriteByte(byte(t))
	if e != nil {
		return e
	}

	// write length
	l := len(buff)
	e = enc.destination.WriteByte(byte(l%256))
	if e != nil {
		return e;
	}
	e = enc.destination.WriteByte(byte(l/256))
	if e != nil {
		return e;
	}

	for i:=0; i<len(buff); i++ {
		e := enc.destination.WriteByte(buff[i])
		if e != nil {
			return e
		}
	}
	return nil
}

func (enc *OTlcEncoding) writeBufferWithDelay(delay int64,t int,buff []byte) error {
	e := enc.destination.WriteByte(byte(t))
	if e != nil {
		return e
	}

	// write length
	l := len(buff)
	e = enc.destination.WriteByte(byte(l%256))
	if e != nil {
		return e;
	}
	e = enc.destination.WriteByte(byte(l/256))
	if e != nil {
		return e;
	}

	for i:=0; i<len(buff); i++ {
		e := enc.destination.WriteByte(buff[i])
		if e != nil {
			return e
		}
	}
	return nil
}

func (enc *OTlcEncoding) WriteUint64(x uint64) error {
	s := fmt.Sprint(x)
	return enc.writeBuffer(INTEGER,[]byte(s))
}

func (enc *OTlcEncoding) WriteString(str string) error {
	buff := bytes.NewBuffer(make([]byte,len(str)))
	buff.Reset()

	for i:=0; i<len(str); i++ {
		e := buff.WriteByte(str[i])
		if e != nil {
			return e
		}
	}

	e := enc.writeBuffer(STRING,buff.Bytes())
	if e != nil {
		return e
	}

	return nil
}

func (enc *OTlcEncoding) WriteInteger(str string) error {
	buff := bytes.NewBuffer(make([]byte,len(str)))
	buff.Reset()

	for i:=0; i<len(str); i++ {
		e := buff.WriteByte(str[i])
		if e != nil {
			return e
		}
	}

	e := enc.writeBuffer(INTEGER,buff.Bytes())
	return e
}

func (enc *OTlcEncoding) WriteBig(x *big.Int) error {
	s := fmt.Sprint(x)
	e := enc.writeBuffer(INTEGER,[]byte(s))
	return e
}


func (enc *OTlcEncoding) WriteBinary(data []byte) error {
	e := enc.writeBuffer(BINARY,data)
	if e != nil {
		return e
	}

	return nil
}

func (enc *OTlcEncoding) WriteBinaryWithDelay(delay int64,data []byte) error {
	e := enc.writeBufferWithDelay(delay,BINARY,data)
	if e != nil {
		return e
	}
	return nil
}


func (enc *OTlcEncoding) WriteStructured(data []byte) error {
	e := enc.writeBuffer(STRUCTURED,data)
	if e != nil {
		return e
	}

	return nil
}

func (enc *OTlcEncoding) GetBuffer() ([]byte,error) {
	return enc.destination.GetBuffer()
}
