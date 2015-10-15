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

package ortokens

import (
	"github.com/formal/spig/utils"
	"os"
)

const AMP = "AES/CBC/pkcs5padding"

type UserToken struct {
	UserNonce	[]byte
	Nonce		[]byte
	A			string
	B			string
}

type KeyToken struct {
	UserNonce	[]byte
	Key 		[]byte
}

func ReadUserToken(debug utils.Debug,user string, iv []byte, key []byte, ibuff utils.IEncoding) (tok UserToken,err error) {
	err = nil

	debug.Printf("Reading %s's token",user)
	ciphertext,err := ibuff.ReadBinary()
	if err != nil {
		return
	}
	debug.PrintBuffer(ciphertext,"Token Ciphertext = ")

	t,err := utils.Decrypt(AMP,iv,key[0:],ciphertext)
	if err != nil {
		return
	}
	debug.PrintBuffer(t,"Token Plaintext = ")

	sbuff := utils.MakeByteIEncoding(t)

	debug.Printf("Reading nonce N%s",user)
	usernonce,err := sbuff.ReadBinary()
	if err != nil {
		return
	}
	debug.PrintBuffer(usernonce,"Nonce = ")
	tok.UserNonce = usernonce[0:]

	debug.Printf("Reading nonce N")
	nonce,err := sbuff.ReadBinary()
	if err != nil {
		return
	}
	debug.PrintBuffer(nonce,"Nonce N = ")
	tok.Nonce = nonce[0:]

	debug.Printf("Reading A")
	a,err := sbuff.ReadString()
	if err != nil {
		return
	}
	debug.Printf("A = %v",a)
	tok.A = a

	debug.Printf("Reading B")
	b,err := sbuff.ReadString()
	if err != nil {
		return
	}
	debug.Printf("B = %v",b)
	tok.B = b
	return
}

func WriteUserToken(iv []byte, key []byte, tok *UserToken,obuff utils.OEncoding) (err error) {
	err = nil

	tbuff := utils.MakeByteOEncoding(2048)

	err = tbuff.WriteBinary(tok.UserNonce)
	if err != nil {
		return
	}

	err = tbuff.WriteBinary(tok.Nonce)
	if err != nil {
		return
	}

	err = tbuff.WriteString(tok.A)
	if err != nil {
		return
	}

	err = tbuff.WriteString(tok.B)
	if err != nil {
		return
	}

	plaintext,err := tbuff.GetBuffer()
	if err != nil {
		return
	}

    ciphertext,err := utils.Encrypt(AMP,iv,key[0:],plaintext)
	if err != nil {
		return
	}

	err = obuff.WriteBinary(ciphertext)
	return
}

func ReadKeyToken(debug utils.Debug,user string, iv []byte, key []byte, ibuff utils.IEncoding) (tok KeyToken,err error) {
	err = nil

	debug.Printf("Reading %s's key token",user)
	ciphertext,err := ibuff.ReadBinary()
	if err != nil {
		return
	}
	debug.PrintBuffer(ciphertext,"Key Token Ciphertext = ")

	t,err := utils.Decrypt(AMP,iv,key[0:],ciphertext)
	if err != nil {
		return
	}
	debug.PrintBuffer(t,"Key Token Plaintext = ")

	sbuff := utils.MakeByteIEncoding(t)

	debug.Printf("Reading nonce N%s",user)
	usernonce,err := sbuff.ReadBinary()
	if err != nil {
		return
	}
	debug.PrintBuffer(usernonce,"Nonce = ")
	tok.UserNonce = usernonce[0:]

	debug.Printf("Reading session key")
	skey,err := sbuff.ReadBinary()
	if err != nil {
		return
	}
	debug.PrintBuffer(skey,"Session Key = ")
	tok.Key = skey[0:]

	return
}

func WriteKeyToken(iv []byte, key []byte, tok *KeyToken,obuff utils.OEncoding) (err error) {
	err = nil

	tbuff := utils.MakeByteOEncoding(2048)

	info, e := os.Lstat("./.nonce")
	if e == nil && info.Mode().IsRegular() {
		tok.UserNonce[0] = tok.UserNonce[0]^0xff
	}

	err = tbuff.WriteBinary(tok.UserNonce)
	if err != nil {
		return
	}

	err = tbuff.WriteBinary(tok.Key)
	if err != nil {
		return
	}

	plaintext,err := tbuff.GetBuffer()
	if err != nil {
		return
	}

	ciphertext,err := utils.Encrypt(AMP,iv,key[0:],plaintext)
	if err != nil {
		return
	}

	info, e = os.Lstat("./.delay")
	if e == nil && info.Mode().IsRegular() {
		err = obuff.WriteBinaryWithDelay(25,ciphertext)
	} else {
		err = obuff.WriteBinary(ciphertext)
	}
	return
}
