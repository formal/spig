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
	defer func() {
		fmt.Printf("... %s worker finished.\n",name)
		conn.Close()
	}()

	debug := utils.NewDebug(utils.USER,name)

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
	debug.PrintBuffer(plaintext,"Plaintext encoding of {T1,B} = ")

	sbuff := utils.MakeByteIEncoding(plaintext)
	debug.Printf("Reading structured data {T1,B}")
	body,e := sbuff.ReadStructured()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(body,"Encoding for T1,B = ")

	cbuff := utils.MakeByteIEncoding(body)

	s,e := cbuff.ReadString()
	debug.Printf("Reading string T1")
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


	tbuff := utils.MakeByteOEncoding(2048)
	e = tbuff.WriteBinary(b)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}


	e = tbuff.WriteString("Along come the scientists and make the words of our fathers into folklore.")
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	body,e = tbuff.GetBuffer()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}

	pbuff := utils.MakeByteOEncoding(2048)
	e = pbuff.WriteStructured(body)
	if e != nil {
		fmt.Printf("%s Error: %v\n",e)
		return
	}


    plaintext,e = pbuff.GetBuffer()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.PrintBuffer(plaintext,"plaintext encoding of {B,T2} = ")

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
    workers.AddWorker("encrypt1",3,worker);
}
