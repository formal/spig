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
	"bytes"
)

//
// byte streams
//

type IByteStream struct {
	source []byte
	index  int
}

func MakeIByteStream(data []byte) *IByteStream {
	return &IByteStream{data,0}
}

func (stream *IByteStream) ReadByte() (byte,error) {
	if stream.index >= len(stream.source) {
		return 0,errors.New("Byte stream empty")
	}
	b :=  stream.source[stream.index]
	stream.index++
	return b,nil
}

type OByteStream struct {
	destination *bytes.Buffer
	index  int
}

func MakeOByteStream(size int) *OByteStream {
	buff := bytes.NewBuffer(make([]byte,1024))
	buff.Reset()
	return &OByteStream{buff,0}
}

func (stream *OByteStream) WriteByte(b byte) error {
	stream.destination.WriteByte(b)
	return nil
}

func (stream *OByteStream) GetBuffer() ([]byte,error) {
	return stream.destination.Bytes(),nil
}
