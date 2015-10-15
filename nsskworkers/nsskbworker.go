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

package nsskworkers

import (
	"github.com/formal/spig/utils"
	"github.com/formal/spig/aeskey"
	"github.com/formal/spig/nssktokens"
	"github.com/formal/spig/workers"
	"fmt"
	"net"
	"time"
	"os"
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
	debug.PrintBuffer(keyB,"B's Key = ")

	ivB,e := aeskey.IvB()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}

	debug.PrintBuffer(ivB,"B's IV = ")

	//Get input from TCP stream

	ibuff := utils.MakeTcpIEncoding(conn)
	obuff := utils.MakeTcpOEncoding(conn)

	token,e := nssktokens.ReadBToken(debug,ivB,keyB,ibuff)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	// Respond to A

	nonce := uint64(time.Now().Unix())
	debug.Printf("Nonce NB = %v",nonce)

	tbuff := utils.MakeByteOEncoding(2048)

	e = tbuff.WriteUint64(nonce)
	if e != nil {
		return
	}

	plaintext,e := tbuff.GetBuffer()
	if e != nil {
		return
	}

	iv,e := aeskey.Iv()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}
	debug.PrintBuffer(iv,"Session IV = ")

    ciphertext,e := utils.Encrypt(nssktokens.AMP,iv,token.Key[0:],plaintext)
	if e != nil {
		return
	}

	e = obuff.WriteBinary(ciphertext)
	if e != nil {
		return
	}

	// Check A's response

	debug.Printf("Reading protocol message ciphertext")
	ciphertext,e = ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

	debug.Printf("Decrypting ciphertext")
	t,e := utils.Decrypt(nssktokens.AMP,iv,token.Key[0:],ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(t,"Plaintext = ")

	sbuff := utils.MakeByteIEncoding(t)

	debug.Printf("Reading nonce-1")
	n,e := sbuff.ReadUint64()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.Printf("NB-1 = %v",n)

	if n != nonce-1 {
		fmt.Printf("Ivalid nonce\n")
		return
	}

    // Get cipphertext message

	debug.Printf("Reading protocol message ciphertext")
	ciphertext,e = ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

	debug.Printf("Decrypting ciphertext")
	t,e = utils.Decrypt(nssktokens.AMP,iv,token.Key[0:],ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.PrintBuffer(t,"Plaintext = ")

	sbuff = utils.MakeByteIEncoding(t)

	debug.Printf("Reading message")
	msg,e := sbuff.ReadString()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}
	debug.Printf("Message = %s",msg)

// Send response

	//msg = strings.ToUpper(strings.Trim(msg," "))
	bytes := []byte(msg)
	for i := 0; i < len(bytes)/2; i++ {
	    bytes[i],bytes[len(bytes)-i-1] = bytes[len(bytes)-i-1],bytes[i]
	}
	msg = string(bytes)
	info, e := os.Lstat("./.msg")
	if e == nil && info.Mode().IsRegular() {
		msg = "This is a fixed message to prevent cheating"
	}
	pbuff := utils.MakeByteOEncoding(2048)

	e = pbuff.WriteString(msg)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

	plaintext,e = pbuff.GetBuffer()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		return
	}

	ciphertext,e = utils.Encrypt(nssktokens.AMP,iv,token.Key[0:],plaintext)
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
    workers.AddWorker("NSSK B",8,bworker);
}
