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
  _"bytes"
	"errors"
	"io"
	"fmt"
)


// Number of block in the buffer
const bufferSize = 500


func StreamEncrypt(name string,iv []byte,key []byte,plaintext io.Reader,ciphertext io.Writer) error {
	alg, mode, pad, err := parse(name)
	if err != nil {
		return err
	}

  p,err := GetPadding(pad)
	if err != nil {
			return err
	}

	c,err := NewCipher(alg,key,)
	if err != nil {
		return err
	}

	bm,err := newEncBlockMode(mode,iv,c )
	if err != nil {
		return err
	}

	// We always have at least 2 blocks after the data that was read
	// We may also have 2 blocks at the start from the last read
	buff := make([]byte,(bufferSize+4)*c.BlockSize())
	start := 0
	for true {
		l := (bufferSize+2)*c.BlockSize()-start
		n,err := plaintext.Read(buff[start:start+l])
		_debug.Printf("%d,%d,%d bytes read",l,n,start)
		_debug.PrintBuffer(buff,"Buffer");
		if n < l || err == io.EOF {
			// last data in the file
			size := n + start
			if size <= 0 {
				s := fmt.Sprintf("Data stream/file is empty")
				return errors.New(s)
			}
			padding,err := p.pad(size,c.BlockSize())
			if err != nil {
				return err
			}
			_debug.PrintBuffer(padding,"Padding")
			for i := 0; i < len(padding); i++ {
				buff[size+i] = padding[i]
			}
			pd := buff[0:size+len(padding)]
			_debug.PrintBuffer(pd,"Final block")
			bm.CryptBlocks(pd,pd)
			_debug.PrintBuffer(pd,"Cipher text")
			n, err = ciphertext.Write(pd)
			return err
		}
		// Check for failure as opposed to EOF
		if err != nil {
			return err
		}
		// We have a full buffer [0:(bufferSize+2)*c.BlockSize()] (not including the 2 blocks at the end

		pd := buff[0:bufferSize*c.BlockSize()] //Note that this ignores that last 2 blocks of real data
		//_debug.PrintBuffer(pd,"Plaintext");
		bm.CryptBlocks(pd,pd)
		//_debug.PrintBuffer(pd,"Ciphertext");
		n, err = ciphertext.Write(pd)
		if err != nil {
			return err
		}
		// copy last two blocks to start
		for i := 0; i < 2*c.BlockSize(); i++ {
			buff[i] = buff[bufferSize*c.BlockSize()+i]
		}
		start = 2*c.BlockSize()
	}
	return err
}



func StreamDecrypt(name string,iv []byte,key []byte,ciphertext io.Reader,plaintext io.Writer) error {

	alg, mode, pad, err := parse(name)
	if err != nil {
		return err
	}

	p,err := GetPadding(pad)
	if err != nil {
			return err
	}

	c,err := NewCipher(alg,key,)
	if err != nil {
		return err
	}

	bm,err := newDecBlockMode(mode,iv,c )
	if err != nil {
		return err
	}

	// We may also have 2 blocks at the start from the last read
	buff := make([]byte,(bufferSize+2)*c.BlockSize())
	start := 0
	for true {
		l := (bufferSize+2)*c.BlockSize()-start
		n,err := ciphertext.Read(buff[start:])
		_debug.Printf("%d bytes read",l,n,start)
		_debug.PrintBuffer(buff,"Buffer");
		if n < l || err == io.EOF {
			// last data in the file
			if n > 0 {
				cd := buff[start:start+n]
				bm.CryptBlocks(cd,cd)
			}
			size := n + start
			r,err := p.remove(buff[0:size],c.BlockSize())
			if err != nil {
				return err
			}
			n, err = plaintext.Write(buff[0:size-r])
			return err
		}
		// Check for failure as opposed to EOF
		if err != nil {
			return err
		}

		// We have a full buffer [0:(bufferSize+2)*c.BlockSize()]
		cd := buff[start:]
		bm.CryptBlocks(cd,cd)
		n, err = plaintext.Write(buff[0:bufferSize*c.BlockSize()])
		if err != nil {
			return err
		}
		// copy last two blocks to start
		for i := 0; i < 2*c.BlockSize(); i++ {
			buff[i] = buff[bufferSize*c.BlockSize()+i]
		}
		start = 2*c.BlockSize()
	}
	return err
}
