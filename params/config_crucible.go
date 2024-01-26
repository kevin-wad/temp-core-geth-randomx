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
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

var (
	// CrucibleChainConfig is the chain parameters to run a node on the Etica main network.
	CrucibleChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID:                 61888,
		ChainID:                   big.NewInt(61888),
		Ethash:                    new(ctypes.EthashConfig),

		//HomesteadBlock: big.NewInt(0),
        //Homestead
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		EIP150Block: big.NewInt(0),
		EIP155Block:  big.NewInt(0),

		//EIP158FBlock: big.NewInt(0),
		// EIP158~
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		//ByzantiumBlock: big.NewInt(0),
		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP649FBlock:  big.NewInt(0), // added
		EIP658FBlock: big.NewInt(0),
		

		//ConstantinopleBlock: big.NewInt(0),
		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),
		EIP1234FBlock:  big.NewInt(0), // added
		EIP1283FBlock: big.NewInt(0), // added

		PetersburgBlock: big.NewInt(0),

		EticaSmartContractv2: big.NewInt(480000), // need to be set to mainnet block height + delay

		
		// EIP1283FBlock:   big.NewInt(9573000),

		/*DisposalBlock:      big.NewInt(5900000),
		ECIP1017FBlock:     big.NewInt(5000000),
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: big.NewInt(3000000),
		ECIP1010Length:     big.NewInt(2000000),
		ECBP1100FBlock:     big.NewInt(11_380_000), // ETA 09 Oct 2020 */
		
	}

)