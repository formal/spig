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
  _ "github.com/formal/spig/echoworker"
	_ "github.com/formal/spig/nestedworker"
	_ "github.com/formal/spig/encrypt0worker"
  _ "github.com/formal/spig/encrypt1worker"
  _ "github.com/formal/spig/orworkers"
  _ "github.com/formal/spig/nsskworkers"
  _ "github.com/formal/spig/zkworker"
)



func main() {
	Start()
}
