package mutations

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params/vars"
)

func CalculateDigest(challengeNumber string, sender common.Address, nonce *big.Int) [32]byte {
	// Convert challengeNumber to bytes
	challengeBytes := []byte(challengeNumber)

	// Convert sender address to bytes
	senderBytes := sender.Bytes()

	// Convert nonce to bytes
	nonceBytes := nonce.Bytes()

	// Concatenate all bytes (this is equivalent to abi.encodePacked in Solidity)
	packed := append(challengeBytes, senderBytes...)
	packed = append(packed, nonceBytes...)

	// Calculate keccak256 hash
	hash := crypto.Keccak256(packed)

	// Convert to [32]byte
	var digest [32]byte
	copy(digest[:], hash)

	return digest
}

func VerifyEticaTransactions(block *types.Block) error {
	fmt.Printf("Verifying Etica transactions for block %d\n", block.NumberU64())
	fmt.Println("==== Begin Block Verification ====")
	fmt.Printf("Block Number: %d\n", block.NumberU64())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", block.ParentHash().Hex())
	fmt.Printf("Timestamp: %v\n", block.Time())
	fmt.Printf("Difficulty: %s\n", block.Difficulty().String())
	fmt.Printf("Gas Limit: %d\n", block.GasLimit())
	fmt.Printf("Gas Used: %d\n", block.GasUsed())

	// Initialize RandomX (do this once, not for every transaction)
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		return fmt.Errorf("failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)

	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		return fmt.Errorf("failed to create RandomX")

	}
	defer DestroyVM(vm)

	for i, tx := range block.Transactions() {
		fmt.Printf("Analyzing transaction %d: %s\n", i, tx.Hash().Hex())

		// ... (previous logging code remains the same)

		if tx.To() != nil && *tx.To() == vars.EticaSmartContractAddress {
			fmt.Println("  Transaction is to Etica smart contract")

			blockHeader, nonce, difficulty, err := ExtractSolutionData(tx.Data())
			if err != nil {
				if err.Error() == "Invalid function selector" {
					fmt.Println("  Transaction is not a mintrandomX call")
					continue
				}
				fmt.Printf("  Error extracting solution data: %v\n", err)
				return err
			}

			fmt.Println("  Transaction is a mintrandomX call")
			fmt.Printf("  Extracted block header: %x\n", blockHeader)
			fmt.Printf("  Extracted nonce: %x\n", nonce)
			fmt.Printf("  Extracted difficulty: %s\n", difficulty.String())
			input := append(blockHeader, nonce...)
			fmt.Printf("Input for solution: %x\n", input)
			correctSolution := CalculateHash(vm, input)
			fmt.Println(" ****** Performing RandomX verification... ********")
			valid, err := CheckSolution(vm, blockHeader, nonce, correctSolution, difficulty)
			if err != nil {
				fmt.Printf("  RandomX verification error: %v\n", err)
				return err
			}
			if valid {
				fmt.Println("*********************** -------- RandomX verification passed ---------- ***********************")
			} else {
				fmt.Println("*************************** -------- RandomX verification failed ---------- ***********************")
			}
		} else {
			fmt.Println("  Transaction is not to Etica smart contract")
		}

		fmt.Println("  Transaction analysis complete")
		fmt.Println("------------------------------------")
	}

	fmt.Printf("Verification complete for block %d\n", block.NumberU64())
	return nil
}

func VerifyEticaTransaction(tx *types.Transaction, statedb *state.StateDB) error {
	fmt.Printf("*-*-*-*-**-*-*-*-*-*-Verifying Etica transaction *-*-*-*-*-**-*-*-*-*-*-*-*-*-")
	fmt.Printf("Verifying Etica transaction: %s\n", tx.Hash().Hex())
	//fmt.Printf("Verifying Etica transaction data: %s\n", tx.Data())

	// Initialize RandomX (you might want to do this once and reuse it)
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		return fmt.Errorf("failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)

	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		return fmt.Errorf("failed to create RandomX")
	}
	defer DestroyVM(vm)
	fmt.Printf("EticaSmartContractAddress: %s\n", vars.EticaSmartContractAddress)
	if tx.To() == nil || *tx.To() != vars.EticaSmartContractAddress {
		fmt.Println("Transaction is not to Etica smart contract")
		return nil
	} else {
		fmt.Println("Transaction is to Etica smart contract")
	}

	// Check if the transaction is calling the mintrandomX() function
	if !IsSolutionProposal(tx.Data()) {
		fmt.Println("Transaction is not a mintrandomX call")
		return nil
	}

	fmt.Println("Transaction is to Etica smart contract")

	blockHeader, nonce, difficulty, err := ExtractSolutionData(tx.Data())
	if err != nil {
		if err.Error() == "Invalid function selector" {
			fmt.Println("Transaction is not a mintrandomX call")
			return nil // Not an error, just not the transaction we're looking for
		}
		fmt.Printf("Error extracting solution data: %v\n", err)
		return err
	}

	fmt.Println("Transaction is a mintrandomX call")
	fmt.Printf("Extracted block header: %x\n", blockHeader)
	fmt.Printf("Extracted nonce: %x\n", nonce)
	fmt.Printf("Extracted difficulty: %s\n", difficulty.String())
	input := append(blockHeader, nonce...)
	fmt.Printf("Input for solution: %x\n", input)
	correctSolution := CalculateHash(vm, input)
	fmt.Println("Performing RandomX verification...")
	fmt.Printf("Input mintrandomX() correctSolution: %x\n", correctSolution)
	valid, err := CheckSolution(vm, blockHeader, nonce, correctSolution, difficulty)
	if err != nil {
		fmt.Printf("RandomX verification error: %v\n", err)
		return err
	}
	if valid {
		fmt.Println("RandomX verification passed")

		// Convert blockHeader to [32]byte
		var challengeNumber [32]byte
		copy(challengeNumber[:], blockHeader)

		// Convert correctSolution to [32]byte
		var solution [32]byte
		copy(solution[:], correctSolution)

		// Get the sender's address
		from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
		if err != nil {
			return fmt.Errorf("failed to get transaction sender: %v", err)
		}

		fmt.Printf("Miner from: %x\n", from)
		// Convert nonce to BigInt
		nonceBigInt := new(big.Int).SetBytes(nonce)
		fmt.Printf("nonceBigInt is: %x\n", nonceBigInt)
		// Update the RandomX state
		updateRandomXState(statedb, challengeNumber, nonceBigInt, from, difficulty)

	} else {
		fmt.Println("RandomX verification failed")
		return fmt.Errorf("invalid RandomX solution")
	}

	return nil
}

func IsSolutionProposal(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	functionSelector := data[:4]
	fmt.Printf("Function selector: %s\n", hex.EncodeToString(functionSelector))
	// Replace with actual selector for mintrandomX
	expectedSelector := []byte{0xee, 0x7a, 0xb6, 0x30} // Actual selector for mintrandomX
	return bytes.Equal(functionSelector, expectedSelector)
}

func ExtractSolutionData(data []byte) (blockHeader []byte, nonce []byte, difficulty *big.Int, err error) {
	// Check if the data is long enough to contain all required fields
	if len(data) < 4+32+32+32 { // 4 bytes function selector + 32 bytes nonce + 32 bytes blockHeader + 32 bytes difficulty
		return nil, nil, nil, errors.New("Data too short to contain solution data")
	}

	// The first 4 bytes are the function selector, which we can check
	functionSelector := data[:4]
	expectedSelector := []byte{0xee, 0x7a, 0xb6, 0x30} // Replace with actual selector for mintrandomX
	if !bytes.Equal(functionSelector, expectedSelector) {
		return nil, nil, nil, errors.New("Invalid function selector")
	}

	// Extract nonce (uint256)
	nonce = make([]byte, 32)
	copy(nonce, data[4:36])

	// Remove leading zeros from nonce
	nonce = bytes.TrimLeft(nonce, "\x00")

	// Extract blockHeader (bytes32)
	blockHeader = make([]byte, 32)
	copy(blockHeader, data[36:68])

	// Extract difficulty (uint)
	difficulty = new(big.Int).SetBytes(data[68:100])

	return blockHeader, nonce, difficulty, nil
}

/*


Given these considerations, here's what we should do:

    Use a fixed-size array in Go:
    To more closely match Solidity's bytes32 and ensure consistent storage calculations, we should use a fixed-size array in Go. This will prevent any unexpected behavior due to varying lengths.
    Update our functions:
    We should modify our functions to use byte instead of []byte for the challengeNumber.

Here's how we can update our code:


*/

func calculateStorageSlot(challengeNumber [32]byte, minerAddress common.Address) common.Hash {
	// The slot of randomXValidatedSolutions is 70
	baseSlot := big.NewInt(70)

	// For the first level of mapping (challengeNumber)
	outerLocation := crypto.Keccak256Hash(
		challengeNumber[:],
		common.LeftPadBytes(baseSlot.Bytes(), 32),
	)

	// For the second level of mapping (minerAddress)
	finalSlot := crypto.Keccak256Hash(
		common.LeftPadBytes(minerAddress.Bytes(), 32),
		outerLocation.Bytes(),
	)

	return finalSlot
}

/*
func calculateStorageSlot(challengeNumber [32]byte, minerAddress common.Address) common.Hash {
	// The slot of randomXValidatedSolutions is 70
	baseSlot := big.NewInt(70)

	// Pad the baseSlot to 32 bytes
	paddedSlot := common.LeftPadBytes(baseSlot.Bytes(), 32)

	// Calculate the slot for the first level of mapping
	challengeNumberHash := crypto.Keccak256Hash(
		challengeNumber[:],
		paddedSlot,
	)

	// Calculate the slot for the second level of mapping
	finalSlot := crypto.Keccak256Hash(
		common.LeftPadBytes(minerAddress.Bytes(), 32),
		challengeNumberHash.Bytes(),
	)

	return finalSlot
} */

/*
func calculateStorageSlot(challengeNumber [32]byte, minerAddress common.Address) common.Hash {
	baseSlot := common.BigToHash(big.NewInt(70))
	challengeNumberHash := crypto.Keccak256Hash(
		challengeNumber[:],
		baseSlot.Bytes(),
	)
	finalSlot := crypto.Keccak256Hash(
		minerAddress.Bytes(),
		challengeNumberHash.Bytes(),
	)
	return finalSlot
} */

/*
func updateRandomXState(statedb *state.StateDB, challengeNumber [32]byte, solution [32]byte, miner common.Address, difficulty *big.Int) {
	solutionSlot := calculateStorageSlot(challengeNumber, miner)
	fmt.Printf("solutionSlot: %s\n", solutionSlot)
	existingSolution := statedb.GetState(vars.EticaSmartContractAddress, solutionSlot)
	fmt.Printf("existingSolution: %s\n", existingSolution)

	// Ensure nonce is 32 bytes
	paddedNonce := common.LeftPadBytes(solution[:], 32)

	// Combine nonce and difficulty similar to abi.encodePacked in Solidity
	combinedBytes := append(paddedNonce, common.LeftPadBytes(difficulty.Bytes(), 32)...)

	solutionHash := crypto.Keccak256Hash(combinedBytes)
	fmt.Printf("Solution Hash: %s\n", solutionHash.Hex())

	if existingSolution != (common.Hash{}) {
		fmt.Println("randomXValidatedSolutions already exists, not updating")
		return
	}

	statedb.SetState(vars.EticaSmartContractAddress, solutionSlot, solutionHash)

	fmt.Printf("Updated randomXValidatedSolutions:\n")
	challengeHex := "0x" + hex.EncodeToString(challengeNumber[:])
	fmt.Printf("Challenge Number: %s\n", challengeHex)
	fmt.Printf("Challenge Number: %x\n", challengeNumber)
	fmt.Printf("Miner Address: %s\n", miner.Hex())
	fmt.Printf("Solution Hash: %s\n", solutionHash.Hex())
} */

func updateRandomXState(statedb *state.StateDB, challengeNumber [32]byte, nonce *big.Int, miner common.Address, difficulty *big.Int) {
	solutionSlot := calculateStorageSlot(challengeNumber, miner)
	fmt.Printf("solutionSlot: %s\n", solutionSlot)
	existingSolution := statedb.GetState(vars.EticaSmartContractAddress, solutionSlot)
	fmt.Printf("existingSolution: %s\n", existingSolution)

	fmt.Printf("updateRandomXState nonce is: %s\n", nonce)

	fmt.Printf("updateRandomXState nonce is: %s\n", nonce)
	fmt.Printf("updateRandomXState difficulty is: %s\n", difficulty)

	// Pack nonce and difficulty as uint256 (32 bytes each)
	packed := make([]byte, 64)
	nonce.FillBytes(packed[:32])
	difficulty.FillBytes(packed[32:])

	solutionHash := crypto.Keccak256Hash(packed)
	fmt.Printf("Solution Hash: %s\n", solutionHash.Hex())

	if existingSolution != (common.Hash{}) {
		fmt.Println("randomXValidatedSolutions already exists, not updating")
		return
	}

	statedb.SetState(vars.EticaSmartContractAddress, solutionSlot, solutionHash)

	fmt.Printf("Updated randomXValidatedSolutions:\n")
	challengeHex := "0x" + hex.EncodeToString(challengeNumber[:])
	fmt.Printf("Challenge Number: %s\n", challengeHex)
	fmt.Printf("Challenge Number: %x\n", challengeNumber)
	fmt.Printf("Miner Address: %s\n", miner.Hex())
	fmt.Printf("Nonce: 0x%x\n", nonce)
	fmt.Printf("Difficulty: %s\n", difficulty.String())
	fmt.Printf("Packed bytes: 0x%x\n", packed) // Add this line to see the packed bytes
	fmt.Printf("Solution Hash: %s\n", solutionHash.Hex())
}

/*
func calculateStorageSlot(challengeNumber [32]byte, minerAddress common.Address) common.Hash {
	baseSlot := common.BigToHash(big.NewInt(70))
	challengeNumberHash := crypto.Keccak256Hash(
		challengeNumber[:],
		baseSlot.Bytes(),
	)
	finalSlot := crypto.Keccak256Hash(
		minerAddress.Bytes(),
		challengeNumberHash.Bytes(),
	)
	return finalSlot
}

func updateRandomXState(statedb *state.StateDB, challengeNumber [32]byte, solution []byte, miner common.Address, difficulty *big.Int) {
	solutionSlot := calculateStorageSlot(challengeNumber, miner)
	existingSolution := statedb.GetState(vars.EticaSmartContractAddress, solutionSlot)
	if existingSolution != (common.Hash{}) {
		fmt.Println("Solution already exists, not updating")
		return
	}
	combinedBytes := append(solution, difficulty.Bytes()...)
	solutionHash := crypto.Keccak256Hash(combinedBytes)
	statedb.SetState(vars.EticaSmartContractAddress, solutionSlot, solutionHash)
	fmt.Printf("Updated randomXValidatedSolutions:\n")
	fmt.Printf("Challenge Number: %x\n", challengeNumber)
	fmt.Printf("Miner Address: %s\n", miner.Hex())
	fmt.Printf("Solution Hash: %s\n", solutionHash.Hex())
} */

/*
func calculateStorageSlot(challengeNumber *big.Int, minerAddress common.Address) common.Hash {
	// The slot of randomXValidatedSolutions is 70
	baseSlot := common.BigToHash(big.NewInt(70))

	// Calculate the slot for the first level of mapping
	challengeNumberHash := crypto.Keccak256Hash(
		challengeNumber.Bytes(),
		baseSlot.Bytes(),
	)

	// Calculate the slot for the second level of mapping
	finalSlot := crypto.Keccak256Hash(
		minerAddress.Bytes(),
		challengeNumberHash.Bytes(),
	)

	return finalSlot
}

func updateRandomXState(statedb *state.StateDB, challengeNumber *big.Int, solution []byte, miner common.Address, difficulty *big.Int) {
	// Calculate storage slot
	solutionSlot := calculateStorageSlot(challengeNumber, miner)

	// Check if solution already exists
	existingSolution := statedb.GetState(vars.EticaSmartContractAddress, solutionSlot)
	if existingSolution != (common.Hash{}) {
		// Solution already exists, do not update
		return
	}

	// Calculate Keccak256 hash of (solution, difficulty)
	combinedBytes := append(solution, difficulty.Bytes()...)
	solutionHash := crypto.Keccak256Hash(combinedBytes)

	// Update solution
	statedb.SetState(vars.EticaSmartContractAddress, solutionSlot, solutionHash)

	fmt.Printf("Updated randomXValidatedSolutions:\n")
	fmt.Printf("Challenge Number: %x\n", challengeNumber)
	fmt.Printf("Miner Address: %s\n", miner.Hex())
	fmt.Printf("Solution Hash: %s\n", solutionHash.Hex())
}
*/
