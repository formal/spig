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
	"net"
)

//
// tcp streams
//


type ITcpStream struct {
	conn  *net.TCPConn
}

func MakeITcpStream(conn *net.TCPConn) *ITcpStream {
	return &ITcpStream{conn}
}

func (stream *ITcpStream) ReadByte() (byte,error) {
	var buff [1]byte
	_,e := stream.conn.Read(buff[0:1])
	return buff[0],e
}

type OTcpStream struct {
	conn  *net.TCPConn
}

func MakeOTcpStream(conn *net.TCPConn) *OTcpStream {
	return &OTcpStream{conn}
}

func (stream *OTcpStream) WriteByte(b byte) error {
	var buff [1]byte
	buff[0] = b
	_, e := stream.conn.Write(buff[0:1])
	return e
}


func (stream *OTcpStream) GetBuffer() ([]byte,error) {
	return nil,errors.New("Cannot retrieve byte array from TCP stream")
}
