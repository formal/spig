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

package orworkers

import (
	"github.com/formal/spig/utils"
	"github.com/formal/spig/aeskey"
	"github.com/formal/spig/ortokens"
	"github.com/formal/spig/workers"
	"fmt"
	"net"
	"crypto/rand"
	"bytes"
	"strconv"
)


func bworker(name string,conn *net.TCPConn) {
	debug := utils.NewDebug(utils.USER,name)
	defer func() {
		debug.Printf("... %s worker finished.",name)
		conn.Close()
	}()

	debug.Printf("%s worker connected to remote address %s",name,conn.RemoteAddr())

// Obtain keys etc.

	keyB,e := aeskey.KeyB()
	if e != nil {
		fmt.Printf("%s AES key error: %v\n",name,e)
		return
	}

	ivB,e := aeskey.IvB()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}
	debug.PrintBuffer(ivB,"B's IV = ")

//Get input from TCP stream

	ibuff := utils.MakeTcpIEncoding(conn)

	debug.Printf("Reading nonce N")
	nonce,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(nonce,"Nonce N = ")

	debug.Printf("Reading A")
	a,e := ibuff.ReadString()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.Printf("A = %v",a)

	//if a != "student" {
	//	fmt.Printf("Incorrect name for A\n")
	//	return
	//}

	debug.Printf("Reading B")
	b,e := ibuff.ReadString()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.Printf("B = %v",b)

	if b != "lecturer" {
		fmt.Printf("Incorrect name for B\n")
		return
	}

	debug.Printf("Reading A's Token")
	tokenA,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(tokenA,"Token Ciphertext = ")


// Send output to the server

	laddr := "127.0.0.1:8005"

	addr, e := net.ResolveTCPAddr("tcp",laddr)
	if e != nil {
		fmt.Printf("Cannot resolve address %s\n",laddr)
		return
	}
	sconn, e := net.DialTCP("tcp",nil,addr)
	if e != nil {
		fmt.Printf("Dialed failed on address %s\n",laddr)
		return
	}

	defer func() {
		sconn.Close()
	}()

	sobuff := utils.MakeTcpOEncoding(sconn)

	e = sobuff.WriteBinary(nonce)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

	e = sobuff.WriteString(a)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

	e = sobuff.WriteString(b)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

	e = sobuff.WriteBinary(tokenA)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

// Set up & send B's Token

	var tokenB ortokens.UserToken

	usernonce := make([]byte,16)
	_,_ = rand.Read(usernonce)

	tokenB.UserNonce = usernonce[0:]
	tokenB.Nonce = nonce[0:]
	tokenB.A = a
	tokenB.B = b

	e = ortokens.WriteUserToken(ivB,keyB,&tokenB,sobuff)

// Read Server Response

	sibuff := utils.MakeTcpIEncoding(sconn)

	debug.Printf("Reading nonce N")
	rnonce,e := sibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(rnonce,"Nonce N = ")

	if !bytes.Equal(rnonce,nonce) {
		fmt.Printf("Invalid nonce\n")
		return
	}

	debug.Printf("Reading A's Key Token")
	keytokenA,e := sibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(keytokenA,"Key Token Ciphertext = ")

	keytokenB,e := ortokens.ReadKeyToken(debug,"B",ivB,keyB,sibuff)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	if !bytes.Equal(keytokenB.UserNonce,tokenB.UserNonce) {
		fmt.Printf("Invalid nonce\n")
		return
	}

// Reply to A

	obuff := utils.MakeTcpOEncoding(conn)

	e = obuff.WriteBinary(nonce)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

	e = obuff.WriteBinary(keytokenA)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}


// Get cipphertext message

	iv,e := aeskey.Iv()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}

	debug.Printf("Reading protocol message ciphertext")
	ciphertext,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

	debug.Printf("Decrypting ciphertext")
	t,e := utils.Decrypt(ortokens.AMP,iv,keytokenB.Key[0:],ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(t,"Plaintext = ")

	sbuff := utils.MakeByteIEncoding(t)

	debug.Printf("Reading message")
	msg,e := sbuff.ReadString()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.Printf("Message = %s",msg)

// Send response

	pbuff := utils.MakeByteOEncoding(2048)

//      e = pbuff.WriteString(strconv.Itoa(len(msg)))
        e = pbuff.WriteInteger(strconv.Itoa(len(msg)))
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

    plaintext,e := pbuff.GetBuffer()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

    ciphertext,e = utils.Encrypt(ortokens.AMP,iv,keytokenB.Key[0:],plaintext)
	if e != nil {
		fmt.Printf("Encryption error: %v\n",e)
		return
	}

	e = obuff.WriteBinary(ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

}

func init() {
    workers.AddWorker("OR B",6,bworker);
}
