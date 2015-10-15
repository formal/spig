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

type INtEncoding struct {
	source IStream
}


func MakeINtEncoding(stream IStream) *INtEncoding {
	enc := INtEncoding{stream}
	return &enc
}

func (enc *INtEncoding) ReadByteOrEOS() (int,error) {
	v,e := enc.source.ReadByte()
	if e != nil {
		return 0,e
	}
	if v == 0 { // Check for end of segment 00 00
		    			// or NULL                  00 ??
		v,e = enc.source.ReadByte()
		if e != nil {
			return 0,e
		}
		if v == 0 {
			return -1,nil
		}
		v = 0
	}
	return int(v),nil
}

func (enc *INtEncoding) _readBuffer() (*bytes.Buffer,error) {
	buff := bytes.NewBuffer(make([]byte,1024))
	buff.Reset()
	for {
		v,e := enc.ReadByteOrEOS()
		if e != nil {
			return nil,e
		}
		if v < 0 {
			break
		}
		buff.WriteByte(byte(v))
	}
	return buff,nil
}


func (enc *INtEncoding) checkType(t int) error {
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

func (enc *INtEncoding) readBuffer(t int) (*bytes.Buffer,error) {
	e := enc.checkType(t)
	if e != nil {
		return nil,e
	}
	return enc._readBuffer()
}

func (enc *INtEncoding) ReadUint64() (uint64,error) {
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

func (enc *INtEncoding) ReadString() (string,error) {
	buff,e := enc.readBuffer(STRING)
	if e != nil {
		return "",e
	}

	return buff.String(),nil
}

func (enc *INtEncoding) ReadInteger() (string,error) {
	buff,e := enc.readBuffer(INTEGER)
	if e != nil {
		return "",e
	}

	return buff.String(),nil
}

func (enc *INtEncoding) ReadBig() (*big.Int,error) {
	x := new(big.Int)
	buff,e := enc.readBuffer(INTEGER)
	if e != nil {
	   return x,e
	}
	str := buff.String()
	x,b := x.SetString(str,10)
	if b {
		return x,nil
	}
	return x,errors.New("Problem reading big integer")
}

func (enc *INtEncoding) ReadBinary() (b []byte,err error) {
	buff,e := enc.readBuffer(BINARY)
	if e != nil {
		return nil,e
	}

	return buff.Bytes(),e
}

func (enc *INtEncoding) ReadStructured() ([]byte,error) {
	buff,e := enc.readBuffer(STRUCTURED)
	if e != nil {
		return nil,e
	}

	return buff.Bytes(),nil
}


//
// Output encoder
//

type ONtEncoding struct {
	destination OStream
}

func MakeONtEncoding(stream OStream) *ONtEncoding {
	enc := ONtEncoding{stream}
	return &enc
}


func (enc *ONtEncoding) write(b byte) error {
	e := enc.destination.WriteByte(b)
	if e != nil {
		return e;
	}
	if b == 0 {
		e := enc.destination.WriteByte(1)
		if e != nil {
			return e;
		}
	}
	return nil
}

func (enc *ONtEncoding) writeEOS() error {
	e := enc.destination.WriteByte(0)
	if e != nil {
		return e;
	}
	e = enc.destination.WriteByte(0)
	if e != nil {
		return e;
	}
	return nil
}

func (enc *ONtEncoding) writeBuffer(buff []byte) error { // we could make this a degenerate case of writeBufferWithDelay
	for i:=0; i<len(buff); i++ {
		e := enc.write(buff[i])
		if e != nil {
			return e
		}
	}
	e := enc.writeEOS()
	if e != nil {
		return e
	}
	return nil
}

func (enc *ONtEncoding) writeBufferWithDelay(delay int64,buff []byte) error {
	for i:=0; i<len(buff); i++ {
		e := enc.write(buff[i])
		if (i == len(buff)/2) {
			fmt.Printf("UNNECESSARY DELAY .....................")
			time.Sleep(time.Duration(delay * 1000000000)); // delay secondssecs
			fmt.Printf(" FINISHED\n")
		}
		if e != nil {
			return e
		}
	}
	e := enc.writeEOS()
	if e != nil {
		return e
	}
	return nil
}

func (enc *ONtEncoding) WriteUint64(x uint64) error {
	e := enc.write(INTEGER)
	if e != nil {
		return e
	}
	s := fmt.Sprint(x)
	return enc.writeBuffer([]byte(s))
}

func (enc *ONtEncoding) WriteString(str string) error {
	e := enc.write(STRING)
	if e != nil {
		return e
	}
	return enc.writeBuffer([]byte(str))
}

func (enc *ONtEncoding) WriteInteger(str string) error {
	e := enc.write(INTEGER)
	if e != nil {
		return e
	}
	return enc.writeBuffer([]byte(str))
}

func (enc *ONtEncoding) WriteBig(x *big.Int) error {
	e := enc.write(INTEGER)
	if e != nil {
		return e
	}
	s := fmt.Sprint(x)
	return enc.writeBuffer([]byte(s))
}

func (enc *ONtEncoding) WriteBinary(data []byte) error {
	e := enc.write(BINARY)
	if e != nil {
		return e
	}
	return enc.writeBuffer(data)
}

func (enc *ONtEncoding) WriteBinaryWithDelay(delay int64,data []byte) error {
	e := enc.write(BINARY)
	if e != nil {
		return e
	}
	return enc.writeBufferWithDelay(delay,data)
}

func (enc *ONtEncoding) WriteStructured(data []byte) error {
	e := enc.write(STRUCTURED)
	if e != nil {
		return e
	}
	return enc.writeBuffer(data)
}

func (enc *ONtEncoding) GetBuffer() ([]byte,error) {
	return enc.destination.GetBuffer()
}
