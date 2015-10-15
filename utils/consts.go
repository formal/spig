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



const (
	STRING = 1
	INTEGER = 2
	BINARY = 3
	STRUCTURED = 4

	TLC = 1 // type-length-content
	NT  = 0  // null terminated

	MAX_UINT64        uint64 = 0xFFFFFFFFFFFFFFFF
	MAX_UINT64_DIV_10 uint64 = MAX_UINT64/10
	MAX_UINT64_MOD_10 uint64 = MAX_UINT64%10
)
