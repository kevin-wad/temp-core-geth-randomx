// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var EticaGenesisHash = common.HexToHash("0x1decfa888f36b53e47d0a90f0178f63152eabe765f70351d6f2b09d4bcde98fd")

// EticaGenesisBlock returns the Etica genesis block.
func DefaultEticaGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     EticaChainConfig,
		Nonce:      85,
		ExtraData:  hexutil.MustDecode("0x"),
		GasLimit:   50000000,
		Difficulty: big.NewInt(2000000),
		Timestamp:  1634447976,
		Alloc:      genesisT.GenesisAlloc{},
	}
}