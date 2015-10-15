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
	"github.com/formal/spig/aeskey"
	"fmt"
	"net"
	"os"
	"flag"
)

const N = 10240;

const AMP = "AES/CBC/pkcs5padding"

func main() {
	var help = flag.Bool("h", false, "help")
	var ip = flag.String("i", "127.0.0.1", "ip address")
	var port = flag.String("p", "8003", "port")
	flag.Parse()
    if *help || flag.NArg() != 1 {
		fmt.Printf("USAGE: encrypt1 <string>\n")
		flag.PrintDefaults()
    	return;
    }

	utils.Version()

	debug := utils.NewDebug(utils.USER,"encrypt1")

	debug.Printf("Alg/Mode/Padding = %s",AMP)

	key,e := aeskey.Key()
	if e != nil {
		fmt.Printf("AES key error: %v\n",e)
		os.Exit(1)
	}
	debug.PrintBuffer(key,"Key = ")

	iv,e := aeskey.Iv()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}
	debug.PrintBuffer(iv,"IV = ")

	laddr := "" + *ip + ":"  + *port

	addr, e := net.ResolveTCPAddr("tcp",laddr)
	if e != nil {
		fmt.Printf("Cannot resolve address %s\n",laddr)
		os.Exit(1)
	}
	conn, e := net.DialTCP("tcp",nil,addr)
	if e != nil {
		fmt.Printf("Dialed failed on address %s\n",laddr)
		os.Exit(2)
	}

	defer func() {
		conn.Close()
	}()

	fmt.Printf("Connected to remote address %s\n",conn.RemoteAddr())
	fmt.Printf("Connected from local address %s\n",conn.LocalAddr())

	sbuff := utils.MakeByteOEncoding(2048)

	e = sbuff.WriteString(flag.Arg(0))
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

	e = sbuff.WriteBinary([]byte(flag.Arg(0)))
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

    body,e := sbuff.GetBuffer()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

	pbuff := utils.MakeByteOEncoding(2048)

	e = pbuff.WriteStructured(body)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

    plaintext,e := pbuff.GetBuffer()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

	debug.PrintBuffer(plaintext,"Plaintext encoding for {T1,B} =")

    ciphertext,e := utils.Encrypt(AMP,iv,key[0:],plaintext)
	if e != nil {
		fmt.Printf("Encryption error: %v\n",e)
	}
	debug.PrintBuffer(ciphertext,"Ciphertext of encoding for {T1,B} =")

	obuff := utils.MakeTcpOEncoding(conn)
	debug.Printf("Sending ciphertext as binary encoding")
	e = obuff.WriteBinary(ciphertext)
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

	ibuff := utils.MakeTcpIEncoding(conn)
	debug.Printf("Reading ciphertext as binary encoding")
	ciphertext,e = ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

    plaintext,e = utils.Decrypt(AMP,iv,key[0:],ciphertext)
	if e != nil {
		fmt.Printf("Decryption error: %v\n",e)
	}
	debug.PrintBuffer(plaintext,"Plaintext encoding of {B,T2} = ")

	tbuff := utils.MakeByteIEncoding(plaintext)

	b,e := tbuff.ReadStructured()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}
	debug.PrintBuffer(b,"Encoding of B,T2 = ")

	cbuff := utils.MakeByteIEncoding(b)

	debug.Printf("Reading buffer B")
	b,e = cbuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}
	debug.PrintBuffer(b,"B = ")

	debug.Printf("Reading string T2")
	s,e := cbuff.ReadString()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}
	debug.Printf("T2 = %s",s)

	fmt.Printf("String received = %s\n",s)
	fmt.Printf("Binary data received contained %d bytes\n",len(b))

}
