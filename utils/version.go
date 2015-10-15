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
)

const v = 0
const s = 1
const n = 0

func Version()  {
	fmt.Printf("SPIG Reference Implementation. Version %d.%d.%d\n",v,s,n)
}

func CurrentVersion() (uint64,uint64,uint64) {
	return v,s,n
}
