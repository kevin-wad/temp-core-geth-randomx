package mutations

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// FindValidSolution searches for a valid nonce that produces a hash meeting the target difficulty
func FindValidSolution(blockHeader []byte, target []byte) (nonce []byte, hash []byte, err error) {
	// Initialize RandomX
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		return nil, nil, fmt.Errorf("failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)

	// Create VM
	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		return nil, nil, fmt.Errorf("failed to create RandomX VM")
	}
	defer DestroyVM(vm)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000000; i++ { // Limit the number of attempts
		// Generate a random nonce
		nonce = make([]byte, 8)
		rand.Read(nonce)

		// Combine block header and nonce
		input := append(blockHeader, nonce...)

		// Calculate the hash
		hash = CalculateHash(vm, input)

		// Check if the hash meets the target
		if bytes.Compare(hash, target) <= 0 {
			return nonce, hash, nil
		}
	}

	return nil, nil, fmt.Errorf("failed to find valid solution after 1,000,000 attempts")
}

func FindValidSolution2(challengeNumber string, difficulty *big.Int, minerAddress common.Address) (nonce *big.Int, hash []byte, err error) {
	// Convert challenge number to bytes
	challengeBytes, err := common.ParseHexOrString(challengeNumber)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse challenge number: %v", err)
	}

	// Initialize RandomX
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		return nil, nil, fmt.Errorf("failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)

	// Create VM
	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		return nil, nil, fmt.Errorf("failed to create RandomX VM")
	}
	defer DestroyVM(vm)

	// Convert difficulty to a 32-byte array
	difficultyBytes := make([]byte, 32)
	difficulty.FillBytes(difficultyBytes)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000000; i++ { // Limit the number of attempts
		// Generate a random nonce
		nonce = new(big.Int).SetUint64(rand.Uint64())

		// Calculate RandomX hash
		input := append(challengeBytes, nonce.Bytes()...)
		hash := CalculateHash(vm, input)

		// Convert RandomX hash to a big.Int for comparison
		hashInt := new(big.Int).SetBytes(hash)

		// Check if the hash meets the difficulty
		if hashInt.Cmp(difficulty) <= 0 {
			fmt.Printf("Correct Input FindValidSolution2(): %x\n", input)
			fmt.Printf("hash: %x\n", hash)
			// Calculate Keccak256 hash of (nonce, difficulty)
			nonceBytes := nonce.Bytes()
			combinedBytes := append(nonceBytes, difficultyBytes...)
			keccakHash := crypto.Keccak256(combinedBytes)
			fmt.Printf("Keccak256(nonce, difficulty): %x\n", keccakHash)
			return nonce, hash, nil
		}
	}

	return nil, nil, fmt.Errorf("failed to find valid solution after 1,000,000 attempts")
}
