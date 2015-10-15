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

package encrypt0worker

import (
	"github.com/formal/spig/utils"
	"github.com/formal/spig/workers"
	"github.com/formal/spig/aeskey"
	"fmt"
	"net"
)

const AMP = "AES/CBC/pkcs5padding"

func worker(name string,conn *net.TCPConn) {

	debug := utils.NewDebug(utils.SYSTEM,name)

	defer func() {
		fmt.Printf("... %s worker finished.\n",name)
		conn.Close()
	}()


	key,e := aeskey.Key()
	if e != nil {
		fmt.Printf("%s AES key error: %v\n",name,e)
		return
	}
	iv,e := aeskey.Iv()
	if e != nil {
		fmt.Printf("%s AES IV error: %v\n",name,e)
		return
	}

	fmt.Printf("%s worker connected to remote address %s\n",name,conn.RemoteAddr())
	ibuff := utils.MakeTcpIEncoding(conn)

	debug.Printf("Reading ciphertext as binary encoding")
	ciphertext,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

	plaintext,e := utils.Decrypt(AMP,iv,key[0:],ciphertext)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(plaintext,"Plaintext encoding of T1,B = ")

	cbuff := utils.MakeByteIEncoding(plaintext)

	debug.Printf("Reading string T1")
	s,e := cbuff.ReadString()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.Printf("T1 = %s",s)

	debug.Printf("Reading buffer B")
	b,e := cbuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(b,"B = ")

	fmt.Printf("You sent the string \"%s\"\n",s)
	fmt.Printf("and the binary data of length %d\n",len(b))

	obuff := utils.MakeTcpOEncoding(conn)

	pbuff := utils.MakeByteOEncoding(2048)
	e = pbuff.WriteBinary(b)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	e = pbuff.WriteString("God is alive. He just doesn't want to get involved.")
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

    plaintext,e = pbuff.GetBuffer()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(plaintext,"plaintext encoding of B,T2 = ")

    ciphertext,e = utils.Encrypt(AMP,iv,key[0:],plaintext)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	debug.PrintBuffer(ciphertext,"Sending binary encoding of ciphertext =")
	e = obuff.WriteBinary(ciphertext)
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

}

func init() {
    workers.AddWorker("encrypt0",2,worker);
}
