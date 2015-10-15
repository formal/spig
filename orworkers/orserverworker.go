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
	"github.com/formal/spig/workers"
	"github.com/formal/spig/ortokens"
	"fmt"
	"net"
	"bytes"
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

	ivA,e := aeskey.IvA()
	if e != nil {
		fmt.Printf("%s AES IV error: %v\n",name,e)
		return
	}

	keyB,e := aeskey.KeyB()
	if e != nil {
		fmt.Printf("%s AES key error: %v\n",name,e)
		return
	}

	ivB,e := aeskey.IvB()
	if e != nil {
		fmt.Printf("%s AES IV error: %v\n",name,e)
		return
	}


	sessionKey,e := aeskey.SessionKey()
	if e != nil {
		fmt.Printf("%s AES key error: %v\n",name,e)
		return
	}


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

	tokenA,e := ortokens.ReadUserToken(debug,"A",ivA,keyA, ibuff)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	if tokenA.A != a || tokenA.B != b || !bytes.Equal(tokenA.Nonce,nonce) {
		fmt.Printf("Invalid token for A\n")
		return
	}

	tokenB,e := ortokens.ReadUserToken(debug,"B",ivB,keyB, ibuff)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	if tokenB.A != a || tokenB.B != b || !bytes.Equal(tokenB.Nonce,nonce) {
		fmt.Printf("Invalid token for B\n")
		return
	}

// Send output to TCP stream

	obuff := utils.MakeTcpOEncoding(conn)

	e = obuff.WriteBinary(nonce)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

// Set up & send A's Key Token

	var keytokenA ortokens.KeyToken

	keytokenA.UserNonce = tokenA.UserNonce[0:]
	keytokenA.Key = sessionKey[0:]

	e = ortokens.WriteKeyToken(ivA,keyA,&keytokenA,obuff)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}

// Set up & send B's Key Token

	var keytokenB ortokens.KeyToken

	keytokenB.UserNonce = tokenB.UserNonce[0:]
	keytokenB.Key = sessionKey[0:]

	e = ortokens.WriteKeyToken(ivB,keyB,&keytokenB,obuff)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}




}

func init() {
    workers.AddWorker("OR S",5,sworker);
}
