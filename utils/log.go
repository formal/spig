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
	"fmt"
	"sync"
	"os"
	"time"
)

var logEnabled bool = true
var fileName string = ""
var nextLog int = 0
var logMutex sync.Mutex

func mkFileName() (s string) {
	t := time.Now().Local()
	m := fmt.Sprintf("00%d",t.Month)
	d := fmt.Sprintf("00%d",t.Day)
	s = fmt.Sprintf("%d%s%s.log",t.Year,m[len(m)-2:],d[len(d)-2:])
	return
}

func setFileName() {
	for {
		fileName = mkFileName()
		time.Sleep(14400000000000); // 4 hours
	}
}

func init() {
//	logEnabled = MyServer()
	fileName = mkFileName()
	if logEnabled {
		go setFileName()
	}
}

func openFile() (*os.File,error) {
	return os.OpenFile(fileName,os.O_CREATE | os.O_APPEND | os.O_WRONLY,0700)
}


type Log interface {
	Printf(format string, a...interface{})
	PrintBuffer(buffer []byte,format string, a...interface{})
}

type logStructure struct {
	id string
}

func NewLog(id string) *logStructure {
	logMutex.Lock()
	n := nextLog
	nextLog++
	logMutex.Unlock()
	id1 := fmt.Sprintf("%s [%d]",id,n)
	return &logStructure{id1}
}

func (log *logStructure) Printf(format string, a...interface{}) {
	if logEnabled {
		logMutex.Lock()
		defer func() {
			logMutex.Unlock()
		}()

		file,e := openFile()
		if e != nil {
			return
		}
		defer func() {
        	file.Close()
		}()

		fmt.Fprintf(file,"\n%s %s\n",log.id, time.Now().Local().Format(time.UnixDate))
		fmt.Fprintf(file,"%s ",log.id)
		fmt.Fprintf(file,format,a)
		fmt.Fprint(file,"\n")
	}
}


func (log *logStructure) PrintBuffer(buffer []byte,format string, a...interface{}) {
	if logEnabled {

		logMutex.Lock()
		defer func() {
			logMutex.Unlock()
		}()

		file,e := openFile()
		if e != nil {
			return
		}
		defer func() {
        	file.Close()
		}()

		fmt.Fprintf(file,"\n%s %s\n",log.id, time.Now().Local().Format(time.UnixDate))
		fmt.Fprintf(file,"%s ",log.id)
		fmt.Fprintf(file,format,a)
		fmt.Fprint(file,"\n")
		index := 0
		for {
			if index >= len(buffer) {
				break;
			}
			fmt.Fprintf(file,"%s : ",log.id)
			for i:=0; i<16; i++ {
				if index+i >= len(buffer) {
					fmt.Fprint(file,"   ")
				} else {
					fmt.Fprintf(file,"%s ",hex(buffer[index+i]))
				}
			}
			for i:=0; i<16; i++ {
				if index+i >= len(buffer) {
					fmt.Fprint(file," ")
				} else {
					fmt.Fprintf(file,"%c",char(buffer[index+i]))
				}
			}
			fmt.Fprint(file,"\n")
			index += 16;
		}
	}
}
