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
	"fmt"
	"net"
	"os"
	"flag"
)

func main() {
	var help = flag.Bool("h", false, "help")
	var ip = flag.String("i", "127.0.0.1", "ip address")
	var port = flag.String("p", "8000", "port")
	flag.Parse()
    if *help || flag.NArg() != 1 {
		fmt.Printf("USAGE: echo <string>\n")
		flag.PrintDefaults()
    	return;
    }

	debug := utils.NewDebug(utils.USER,"echo")

	utils.Version()

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

	obuff := utils.MakeTcpOEncoding(conn)

	debug.Printf("Sending string T")
	e = obuff.WriteString(flag.Arg(0))
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

	debug.Printf("Sending buffer B")
	e = obuff.WriteBinary([]byte(flag.Arg(0)))
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}

	ibuff := utils.MakeTcpIEncoding(conn)
	debug.Printf("Reading buffer B")
	b,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}
	debug.PrintBuffer(b,"B = ")

	debug.Printf("Reading string T")
	s,e := ibuff.ReadString()
	if e != nil {
		fmt.Printf("Error: %s\n",e)
		os.Exit(1)
	}
	debug.Printf("T = %s",s)

	fmt.Printf("String received = %s\n",s)
	fmt.Printf("Binary data received contained %d bytes\n",len(b))

}
