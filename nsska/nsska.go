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
	"github.com/formal/spig/nssktokens"
	"github.com/formal/spig/aeskey"

	"fmt"
	"net"
	"bytes"
	"flag"
  "crypto/rand"
)

var a string = "student"
var b string = "lecturer"


func contact_S (debug utils.Debug,ip string,port string) (token nssktokens.AToken,e error) {
	e = nil

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

	nonce := make([]byte,16)
	_,_ = rand.Read(nonce)
	debug.PrintBuffer(nonce,"Nonce N = ")

	// connect to the server

	laddr := "" + ip + ":"  + port

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

	e = obuff.WriteString(a)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	e = obuff.WriteString(b)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	e = obuff.WriteBinary(nonce)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	// Read S's Response

	ibuff := utils.MakeTcpIEncoding(conn)

	token,e = nssktokens.ReadAToken(debug,ivA,keyA,ibuff)
	if e != nil {
		fmt.Printf("Error: %v\n",e)
		return
	}

	if !bytes.Equal(token.Nonce,nonce) {
		fmt.Printf("Invalid nonce\n")
		return
	}

	if token.B != b {
		fmt.Printf("Invalid B in token\n")
		return
	}
	return
}

func contact_B (debug utils.Debug,ip string,port string,token nssktokens.AToken,message string) (e error) {
	e =nil

	// connect to the server

	laddr := "" + ip + ":"  + port

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

	e = obuff.WriteBinary(token.CipherText)
	if e != nil {
		fmt.Printf("error: %v\n",e)
		return
	}

	// Read B's Response

	iv,e := aeskey.Iv()
	if e != nil {
		fmt.Printf("AES IV error: %v\n",e)
		return
	}
	debug.PrintBuffer(iv,"Session IV = ")



	ibuff := utils.MakeTcpIEncoding(conn)

	debug.Printf("Reading B's response")
	ciphertext,e := ibuff.ReadBinary()
	if e != nil {
		return
	}
	debug.PrintBuffer(ciphertext,"Ciphertext = ")

	t,e := utils.Decrypt(nssktokens.AMP,iv,token.Key[0:],ciphertext)
	if e != nil {
		return
	}
	debug.PrintBuffer(t,"Plaintext = ")

	sbuff := utils.MakeByteIEncoding(t)

	debug.Printf("Reading nonce NB")
	nonce,e := sbuff.ReadUint64()
	if e != nil {
		return
	}
	debug.Printf("Nonce NB = %v",nonce)

	// Respond to B

	tbuff := utils.MakeByteOEncoding(2048)

	e = tbuff.WriteUint64(nonce-1)
	if e != nil {
		return
	}

	plaintext,e := tbuff.GetBuffer()
	if e != nil {
		return
	}

    ciphertext,e = utils.Encrypt(nssktokens.AMP,iv,token.Key[0:],plaintext)
	if e != nil {
		return
	}

	e = obuff.WriteBinary(ciphertext)
	if e != nil {
		return
	}

	// Send ciphertext

	pbuff := utils.MakeByteOEncoding(2048)

	e = pbuff.WriteString(message)
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

	// Get cipphertext response

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

	fmt.Printf("%s\n",msg)
	return
}

func main() {
	var help = flag.Bool("h", false, "help")
	var ip_S = flag.String("i", "127.0.0.1", "S ip address")
	var ip_B = flag.String("ib", "127.0.0.1", "B ip address")
	var port_S = flag.String("p", "8007", "S port")
	var port_B = flag.String("pb", "8008", "B port")
	flag.Parse()
    if *help || flag.NArg() != 1 {
		fmt.Printf("USAGE: nsska <string>\n")
		flag.PrintDefaults()
    	return;
    }

	utils.Version()

	debug := utils.NewDebug(utils.USER,"NSSK A")

	token_A, err := contact_S(debug,*ip_S,*port_S)
	if err != nil {
		fmt.Printf("ERROR: %v\n",err)
		return
	}

	err = contact_B(debug,*ip_B,*port_B,token_A,flag.Arg(0))
	if err != nil {
		fmt.Printf("ERROR: %v\n",err)
		return
	}


}
