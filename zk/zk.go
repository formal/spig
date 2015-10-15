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
	"github.com/formal/spig/zknumbers"
	"fmt"
	"net"
	"flag"
	"strconv"
	"crypto/rand"
)

var debugger utils.Debug

func coinIsHead() (bool) {
	coin := make([]byte,1)
	_,_ = rand.Read(coin)
	return int(coin[0]) % 2 == 0
}


func tryOnce(ip string,port string)(ok bool,e error) {
	ok = false
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

	//fmt.Printf("Connected to remote address %s\n",conn.RemoteAddr())
	//fmt.Printf("Connected from local address %s\n",conn.LocalAddr())
	fmt.Print(".")

	obuff := utils.MakeTcpOEncoding(conn)
	ibuff := utils.MakeTcpIEncoding(conn)

	sn,e := strconv.ParseInt(flag.Arg(0),10,64)
	if e != nil {
		debugger.Printf("Error: %v\n",e)
		return
	}
	e = obuff.WriteUint64(uint64(sn))
	if e != nil {
		debugger.Printf("Error: %v\n",e)
		return
	}

	debugger.Printf("Reading x")
	x,e := ibuff.ReadBig()
	if e != nil {
		debugger.Printf("Error: %v\n",e)
		return
	}
	debugger.Printf("x = %v\n",x)


	c := 0
	if coinIsHead() {
		c = 1
	}
	debugger.Printf("c = %d\n",c)

	e = obuff.WriteUint64(uint64(c))
	if e != nil {
		debugger.Printf("Error: %v\n",e)
		return
	}

	debugger.Printf("Reading y")
	y,e := ibuff.ReadBig()
	if e != nil {
		debugger.Printf("Error: %v\n",e)
		return
	}
	debugger.Printf("y = %v\n",y)

	// Check the result
	y.Mul(y,y)
	y.Mod(y,zknumbers.N)

	if (c == 1) {
		x.Mul(x,zknumbers.X)
		x.Mod(x,zknumbers.N)
	}
	debugger.Printf("y**2 = %v\n",y)
	debugger.Printf("x or x*X =%v\n",x )
	return x.Cmp(y) == 0,nil
}


func test(ip string,port string,tries int) (correct, wrong int) {
	correct, wrong = 0,0
	for i := 0; i < tries; i++ {
		b,e := tryOnce(ip,port)
		if e == nil {
			if b {
				correct++;
			} else {
				wrong++;
			}
		}
	}
	return
}


func main() {
	var help = flag.Bool("h", false, "help")
	var ip = flag.String("i", "127.0.0.1", "ip address")
	var attemps = flag.Int("a", 20, "attempts")
	flag.Parse()
	if *help || flag.NArg() != 1 {
		fmt.Printf("USAGE: zk <student number>\n")
		flag.PrintDefaults()
		return;
	}

	utils.Version()

	debugger = utils.NewDebug(utils.USER,"ZK")
	fmt.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
	c,w := test(*ip,"8010",*attemps)
	fmt.Printf("Correct: %d; Wrong: %d\n",c,w)
	fmt.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
	c,w = test(*ip,"8011",*attemps)
	fmt.Printf("Correct: %d; Wrong: %d\n",c,w)
	fmt.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
	c,w = test(*ip,"8012",*attemps)
	fmt.Printf("Correct: %d; Wrong: %d\n",c,w)
	fmt.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
	c,w = test(*ip,"8013",*attemps)
	fmt.Printf("Correct: %d; Wrong: %d\n",c,w)
}
