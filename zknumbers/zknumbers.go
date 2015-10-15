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

package zknumbers

import (
	"math/big"
)

var N *big.Int 			= big.NewInt(0)
var Z *big.Int 			= big.NewInt(0)
var X *big.Int 			= big.NewInt(0)
var InverseX *big.Int 	= big.NewInt(0)


func init() {
	N.SetString("46202163959188865788186592555816982554618337464084392899700696882997367398909",10);
	Z.SetString("32978234982374982374982365606627168723423423466662137468234",10);
	X.Mod(X.Mul(Z,Z),N)
	InverseX.ModInverse(X,N);
}
