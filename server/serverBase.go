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
	"github.com/formal/spig/workers"
	"github.com/formal/spig/utils"
	"fmt"
	"net"
	"os"
	"time"
	"flag"
)

func listener(name string,port int,fn workers.WorkerFunc) {
	//fmt.Printf("name = %s; port = %d\n",name,port)
	laddr := fmt.Sprintf("0.0.0.0:%d",port)
	addr, e := net.ResolveTCPAddr("tcp",laddr)
	if e != nil {
		fmt.Printf("Cannot resolve address %s\n",laddr)
		os.Exit(1)
	}
	socket,e := net.ListenTCP("tcp",addr)
	if e != nil {
		fmt.Printf("Failed to listen on address %s\n",laddr)
		os.Exit(2)
	}
	fmt.Printf("%s server listening on port %d\n",name,port)

	log := utils.NewLog(name)
	for {
		conn, e := socket.AcceptTCP()
		if e != nil {
			fmt.Printf("Accept failed on address %s\n",laddr)
			//break
			time.Sleep(60000000000); // 30 min
		}
		log.Printf("Connected to remote address %s",conn.RemoteAddr())
		go fn(name,conn)
	}

}

func Start() {
	var p = flag.Int("p", 8000, "base port")
	var help = flag.Bool("h", false, "help")
	flag.Parse()
    if *help || flag.NArg() != 0 {
		fmt.Printf("Usage: server\n")
		flag.PrintDefaults()
    	return;
    }

	log := utils.NewLog("server");
	log.Printf("Server started.")
	utils.Version()
	workers.Apply(listener,*p)
	for{
		fmt.Printf("++++ %s\n",time.Now().Local().Format(time.UnixDate))
		time.Sleep(1800000000000); // 30 min
	}
}
