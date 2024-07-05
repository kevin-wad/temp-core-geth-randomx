// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vars

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Etica Mainnet Smart Contract //
var EticaSmartContractAddress = common.HexToAddress("0x352033FEAAa43c1aBD33Fa457220695C32C02feD") // Etica Mainnet Smart Contract
// Crucible Testnet Smart Contract //
var CrucibleSmartContractAddress = common.HexToAddress("0x352033FEAAa43c1aBD33Fa457220695C32C02feD")

// --------- Smart contract hardfork 1, main smart contract loads bytecode from following contract ----------- //
var EticaSmartContractAddressv2 = common.HexToAddress("0xd129Ce1842d5E5c93cF91B9024d863869D33d1cc") // Etica v2, Meticulous Hardfork
var CrucibleSmartContractAddressv2 = common.HexToAddress("0xd129Ce1842d5E5c93cF91B9024d863869D33d1cc")

// --------- Smart contract hardfork 1 ----------- //

// Eticav2ForkBlockExtra is the block header extra-data field to set for the Eticav2 fork
// point and a number of consecutive blocks to allow fast/light syncers to correctly
// pick the side they want.  0x657469636176322d686172642d666f726b is hex representation of "eticav2-hard-fork".
var Eticav2ForkBlockExtra = common.FromHex("0x657469636176322d686172642d666f726b")

// Eticav2ForkExtraRange is the number of consecutive blocks from the Eticav2 fork point
// to override the extra-data in to prevent no-fork attacks.
var Eticav2ForkExtraRange = big.NewInt(10)

// --------- Etica smart contract hardfork 1 ----------- //
