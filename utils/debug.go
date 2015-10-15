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

package utils

import (
	"os"
	"fmt"
	"flag"
	"sync"
)

var debugEnabled  bool = false
var systemEnabled bool = false

var nextDebug int = 0
var debugMutex sync.Mutex
var Default Debug

func init() {
	flag.BoolVar(&debugEnabled, "d", false,"enable debugging")
	f,e := os.Open("DEBUG");
	if e != nil {
		systemEnabled = false
	} else {
		systemEnabled = true
		f.Close()
	}
	Default = NewDebug(SYSTEM,"System")
}

type Debug interface {
	Printf(format string, a...interface{})
	PrintBuffer(buffer []byte,format string, a...interface{})
}

//
// Dummy debugger
//

type dummyStructure struct {
	id int
}

func (debug *dummyStructure)  Printf(format string, a...interface{}) {
}

func (debug *dummyStructure) PrintBuffer(buffer []byte,format string, a...interface{}) {
}

//
// The real debugger
//

type debugStructure struct {
	id 			string
	isSystem 	bool
}

func (debug *debugStructure) Printf(format string, a...interface{}) {
	if !debugEnabled {
		return
	}
	if debug.isSystem && !systemEnabled {
		return
	}

	debugMutex.Lock()
	defer func() {
		debugMutex.Unlock()
	}()

	fmt.Printf("%s : ",debug.id)
	fmt.Printf(format,a...)
	fmt.Print("\n")
}

func hex(b byte) string {
	s := "0" + fmt.Sprintf("%x",int(b))
	return s[len(s)-2:len(s)]
}

func char(b byte) (c int) {
	if b >= 32 && b <= 127 {
		c = int(b)
	} else {
		c = 32
	}
	return
}

func (debug *debugStructure) PrintBuffer(buffer []byte,format string, a...interface{}) {
	if !debugEnabled {
		return
	}
	if debug.isSystem && !systemEnabled {
		return
	}

	debugMutex.Lock()
	defer func() {
		debugMutex.Unlock()
	}()

	fmt.Printf("%s : ",debug.id)
	fmt.Printf(format,a...)
	fmt.Print("\n")
	if len(buffer) == 0 {
		fmt.Printf("%s : ",debug.id)
		fmt.Print("**** Empty Buffer! ****\n")
		return
	}
	index := 0
	for {
		if index >= len(buffer) {
			break;
		}
		fmt.Printf("%s : ",debug.id)
		for i:=0; i<16; i++ {
			if index+i >= len(buffer) {
				fmt.Print("   ")
			} else {
				fmt.Printf("%s ",hex(buffer[index+i]))
			}
		}
		for i:=0; i<16; i++ {
			if index+i >= len(buffer) {
				fmt.Print(" ")
			} else {
				fmt.Printf("%c",char(buffer[index+i]))
			}
		}
		fmt.Print("\n")
		index += 16;
	}
}

//
//
//

type DebugType byte
const (
    DUMMY = iota
    SYSTEM
    USER
)

func NewDebug(t DebugType,id string) Debug {
	if t == DUMMY {
		return &dummyStructure{0}
	}
	debugMutex.Lock()
	defer func() {
		debugMutex.Unlock()
	}()
	n := nextDebug
	nextDebug++
	id1 := fmt.Sprintf("[%d] %s ",n,id)
	if t == SYSTEM {
		return &debugStructure{id1,true}
	}
	return &debugStructure{id1,false}
}
