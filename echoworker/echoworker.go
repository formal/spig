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

package echoworker

import (
	"github.com/formal/spig/utils"
	"github.com/formal/spig/workers"
	"fmt"
	"net"
)

func worker(name string,conn *net.TCPConn) {

	debug := utils.NewDebug(utils.USER,name)

	defer func() {
		fmt.Printf("... %s worker finished.\n",name)
		conn.Close()
	}()

	fmt.Printf("%s worker connected to remote address %s\n",name,conn.RemoteAddr())
	ibuff := utils.MakeTcpIEncoding(conn)

	debug.Printf("Reading string T")
	s,e := ibuff.ReadString()
	if e != nil {
		fmt.Printf("%s Error: %s\n",name,e)
		return
	}
	debug.Printf("T = %s",s)

	debug.Printf("Reading buffer B")
	b,e := ibuff.ReadBinary()
	if e != nil {
		fmt.Printf("%s Error: %s\n",name,e)
		return
	}
	debug.PrintBuffer(b,"B = ")

	fmt.Printf("You sent the string \"%s\"\n",s)
	fmt.Printf("and the binary data of length %d\n",len(b))

	obuff := utils.MakeTcpOEncoding(conn)
	debug.Printf("Sending buffer B")
	e = obuff.WriteBinary(b)
	if e != nil {
		fmt.Printf("%s Error: %s\n",name,e)
		return
	}
	debug.Printf("Sending string T")
	e = obuff.WriteString(s)
	if e != nil {
		fmt.Printf("%s Error: %s\n",name,e)
		return
	}
}

func init() {
    workers.AddWorker("echo",0,worker);
}
