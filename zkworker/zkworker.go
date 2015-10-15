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

package zkworker

import (
	"github.com/formal/spig/utils"
	"github.com/formal/spig/workers"
	"github.com/formal/spig/zknumbers"

	"math/big"
	"fmt"
	"net"
	"crypto/rand"
)

func random() (x *big.Int) {
	kData := make([]byte,50)
	_,_ = rand.Read(kData)
	x = new(big.Int)
	x.Mod(x.SetBytes(kData),zknumbers.N)
	return
}

func coinIsHead() (bool) {
	coin := make([]byte,1)
	_,_ = rand.Read(coin)
	return int(coin[0]) % 2 == 0
}


func worker(name string,conn *net.TCPConn) {
	debug := utils.NewDebug(utils.USER,name)
	defer func() {
		debug.Printf("... worker finished.")
		conn.Close()
	}()

	instance := int(name[3])-int('0')

	debug.Printf("Worker connected to remote address %s",conn.RemoteAddr())

	ibuff := utils.MakeTcpIEncoding(conn)
	obuff := utils.MakeTcpOEncoding(conn)

	//debug.Printf("Reading student number")
	id,e := ibuff.ReadUint64()
	if e != nil {
		fmt.Printf("%s Error: %v\n",name,e)
		return
	}
	debug.Printf("Student Number = %d",id)

	instance = (int(id) + instance) % 4

	//
	// We know the square root
	// so we play the game properly
	//
	//debug.Printf("Instance %d\n",instance)
	if instance == 0 {
		k := random()
		x := new(big.Int).Mul(k,k)
		x.Mod(x,zknumbers.N)
		e = obuff.WriteBig(x)
		if e != nil {
			fmt.Printf("%s Error: %v\n",e)
			return
		}

		//debug.Printf("Reading challenge value")
		c,e := ibuff.ReadUint64()
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}
		debug.Printf("Challenge value = %d",c)

		if c != 0 {
			k.Mod(k.Mul(k,zknumbers.Z),zknumbers.N)
		}
		e = obuff.WriteBig(k)
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}
		return;
	}

	//
	// We don't know the square root
	// so we decide at random which challenge we will answer
	//
	if (instance == 1){
		if coinIsHead() {
			instance = 2
		} else {
			instance = 3
		}
	}

	//
	// We don't know the square root
	// so we will answer challenge 0
	//
	if instance == 2 {
		k := random()
		x := new(big.Int).Mul(k,k)
		x.Mod(x,zknumbers.N)
		e = obuff.WriteBig(x)
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}

		//debug.Printf("Reading challenge value")
		c,e := ibuff.ReadUint64()
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}
		debug.Printf("Challenge value = %d",c)

		if c != 0 {
			return
		}
		e = obuff.WriteBig(k)
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}
		return;
	}

	//
	// We don't know the square root
	// so we will answer challenge 1
	//
	if instance == 3 {
		r := random()
		x := new(big.Int).Mul(r,r)
		x.Mod(x.Mul(x,zknumbers.InverseX),zknumbers.N)
		e = obuff.WriteBig(x)
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}

		//debug.Printf("Reading challenge value")
		c,e := ibuff.ReadUint64()
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}
		debug.Printf("Challenge value = %d",c)

		if c != 1 {
			return
		}
		e = obuff.WriteBig(r)
		if e != nil {
			fmt.Printf("%s Error: %v\n",name,e)
			return
		}
		return;
	}

	fmt.Printf("Error: %s worker did not work properly\n",name)
}

func init() {
    workers.AddWorker("ZK10",10,worker);
    workers.AddWorker("ZK11",11,worker);
    workers.AddWorker("ZK12",12,worker);
    workers.AddWorker("ZK13",13,worker);
}
