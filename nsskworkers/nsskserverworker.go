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
)


func sworker(name string,conn *net.TCPConn) {
	debug := utils.NewDebug(utils.USER,name)
	defer func() {
		debug.Printf("... %s worker finished.",name)
		conn.Close()
	}()

	debug.Printf("%s worker connected to remote address %s",name,conn.RemoteAddr())

// Obtain keys etc.

	keyA,e := aeskey.KeyA()
	if e != nil {
		fmt.Printf("%s AES key error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(keyA,"A's Key = ")

	ivA,e := aeskey.IvA()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}
	debug.PrintBuffer(ivA,"A's IV = ")

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

	sessionKey,e := aeskey.SessionKey()
	if e != nil {
		fmt.Printf("%s AES key error: %v\n",name,e)
		return
	}

//Get input from TCP stream

	ibuff := utils.MakeTcpIEncoding(conn)

	debug.Printf("Reading A")
	a,e := ibuff.ReadString()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.Printf("A = %v",a)

	if a != "student" {
		fmt.Printf("Incorrect name for A\n")
		return
	}

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

	debug.Printf("Reading nonce N")
	nonce,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(nonce,"Nonce N = ")

// Send output to TCP stream

	obuff := utils.MakeTcpOEncoding(conn)

// Set up & send B's Key Token

	var token_B nssktokens.BToken

	token_B.A = a
	token_B.Key = sessionKey[0:]

	ciphertext,e := nssktokens.WriteBToken(debug,ivB,keyB,&token_B)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

// Set up & send A's Key Token

	var token_A nssktokens.AToken

	token_A.Nonce = nonce[0:]
	token_A.B = b
	token_A.Key = sessionKey[0:]
	token_A.CipherText = ciphertext
	e = nssktokens.WriteAToken(debug,ivA,keyA,&token_A,obuff)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

}

func init() {
    workers.AddWorker("NSSK S",7,sworker);
}
