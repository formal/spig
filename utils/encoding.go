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
	"flag"
	"fmt"
	"net"
	"math/big"
)

var encoding int = NT

func init() {
	  str := fmt.Sprintf("encoding NT(%d) TLV(%d)",NT,TLC)
		flag.IntVar(&encoding, "e", NT ,str)
}

//
// Input decoding
//

type IEncoding interface {
	ReadUint64() (uint64,error)
	ReadString() (string,error)
	ReadInteger() (string,error)
	ReadBig() (*big.Int,error)
	ReadBinary() (b []byte,err error)
	ReadStructured() ([]byte,error)
}

//
// Output encoding
//

type OEncoding interface {
	WriteUint64(x uint64) error
	WriteString(str string) error
	WriteInteger(str string) error
	WriteBig(*big.Int) error
	WriteBinary(data []byte) error
	WriteBinaryWithDelay(delay int64,data []byte) error
	WriteStructured(data []byte) error
	GetBuffer() ([]byte,error)
}



func makeIEncoding(s IStream) IEncoding {
	if encoding == NT {
		return MakeINtEncoding(s)
	}
	return MakeITlcEncoding(s)
}

func makeOEncoding(s OStream) OEncoding {
	if encoding == NT {
		return MakeONtEncoding(s)
	}
	return MakeOTlcEncoding(s)
}

func MakeByteIEncoding(data []byte) IEncoding {
	s := MakeIByteStream(data);
	return makeIEncoding(s)
}

func MakeByteOEncoding(size int) OEncoding {
	stream := MakeOByteStream(size);
	return makeOEncoding(stream)
}


func MakeTcpIEncoding(conn *net.TCPConn) IEncoding {
	stream := MakeITcpStream(conn);
	return makeIEncoding(stream)
}

func MakeTcpOEncoding(conn *net.TCPConn) OEncoding {
	stream := MakeOTcpStream(conn);
	return makeOEncoding(stream)
}
