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
	"flag"
)
func main() {

	var help = flag.Bool("h", false, "help")
	flag.Parse()
    if *help || flag.NArg() != 0 {
		fmt.Printf("Usage: showor\n")
		flag.PrintDefaults()
    	return;
    }


	utils.Version()
	fmt.Printf("Otway Rees  Key for A = %s\n",aeskey.KeyAString())
	fmt.Printf("Otway Rees   IV for A = %s\n",aeskey.IvAString())

	fmt.Printf("Otway Rees  Key for B = %s\n",aeskey.KeyBString())
	fmt.Printf("Otway Rees   IV for B = %s\n",aeskey.IvBString())

	fmt.Printf("Otway Rees   IV for Session = %s\n",aeskey.IvString())
}
