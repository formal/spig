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

package main

import (
	"github.com/formal/spig/utils"
	"github.com/formal/spig/ortokens"
	"github.com/formal/spig/aeskey"
	"fmt"
	"net"
	"bytes"
	"flag"
  "crypto/rand"
)

func main() {
	var help = flag.Bool("h", false, "help")
	var ip = flag.String("i", "127.0.0.1", "ip address")
	var port = flag.String("p", "8006", "port")
	flag.Parse()
	if *help || flag.NArg() != 1 {
		fmt.Printf("USAGE: ora <string>\n")
		flag.PrintDefaults()
	   return;
	}

	utils.Version()

	debug := utils.NewDebug(utils.USER,"OR A")

	keyA,e := aeskey.KeyA()
	if e != nil {
		fmt.Printf("AES key error: %v\n",e)
		return
	}
	debug.PrintBuffer(keyA,"A's Key = ")

	ivA,e := aeskey.IvA()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}
	debug.PrintBuffer(ivA,"A's IV = ")

	laddr := "" + *ip + ":"  + *port

	addr, e := net.ResolveTCPAddr("tcp",laddr)
	if e != nil {
		fmt.Printf("Cannot resolve address %s\n",laddr)
		return
	}

	conn, e := net.DialTCP("tcp",nil,addr)
	if e != nil {
		fmt.Printf("Dialed failed on address %s\n",laddr)
		return
	}

	defer func() {
		conn.Close()
	}()

	fmt.Printf("Connected to remote address %s\n",conn.RemoteAddr())
	fmt.Printf("Connected from local address %s\n",conn.LocalAddr())

	obuff := utils.MakeTcpOEncoding(conn)

	nonce := make([]byte,16)
	_,_ = rand.Read(nonce)

	usernonce := make([]byte,16)
	_,_ = rand.Read(usernonce)


	e = obuff.WriteBinary(nonce)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	e = obuff.WriteString("student")
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	e = obuff.WriteString("lecturer")
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}


// Set up & send A's Token

	var tokenA ortokens.UserToken

	tokenA.UserNonce = usernonce[0:]
	tokenA.Nonce = nonce[0:]
	tokenA.A = "student"
	tokenA.B = "lecturer"

	e = ortokens.WriteUserToken(ivA,keyA,&tokenA,obuff)

// Read B's Response

	ibuff := utils.MakeTcpIEncoding(conn)

	debug.Printf("Reading nonce N")
	rnonce,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}
	debug.PrintBuffer(rnonce,"Nonce N = ")

	if !bytes.Equal(rnonce,nonce) {
		fmt.Printf("Invalid nonce\n")
		return
	}

	keytokenA,e := ortokens.ReadKeyToken(debug,"A",ivA,keyA,ibuff)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	if !bytes.Equal(keytokenA.UserNonce,tokenA.UserNonce) {
		fmt.Printf("Invalid nonce\n")
		return
	}

// Send ciphertext

	iv,e := aeskey.Iv()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}

	pbuff := utils.MakeByteOEncoding(2048)

	e = pbuff.WriteString(flag.Arg(0))
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

    plaintext,e := pbuff.GetBuffer()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

    ciphertext,e := utils.Encrypt(ortokens.AMP,iv,keytokenA.Key[0:],plaintext)
	if e != nil {
		fmt.Printf("Encryption error: %v\n",e)
		return
	}

	e = obuff.WriteBinary(ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

// Get cipphertext response

	debug.Printf("Reading protocol message ciphertext")
	ciphertext,e = ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

	debug.Printf("Decrypting ciphertext")
	t,e := utils.Decrypt(ortokens.AMP,iv,keytokenA.Key[0:],ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(t,"Plaintext = ")

	sbuff := utils.MakeByteIEncoding(t)

	debug.Printf("Reading message")
	msg,e := sbuff.ReadInteger()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

	fmt.Printf("%s\n",msg)

}
