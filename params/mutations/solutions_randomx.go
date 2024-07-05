package mutations

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"unsafe"
)

// CheckSolutionWithTarget verifies a mining solution using RandomX
func CheckSolutionWithTarget(vm unsafe.Pointer, blockHeader []byte, nonce []byte, solution []byte, target []byte) (bool, error) {
	if vm == nil {
		return false, errors.New("RandomX VM is not initialized")
	}

	// Combine block header and nonce
	input := append(blockHeader, nonce...)

	// Calculate the hash
	hash := CalculateHash(vm, input)

	// Compare the hash with the provided solution
	if !bytes.Equal(hash, solution) {
		return false, errors.New("solution does not match calculated hash")
	}

	// Check if the hash meets the target difficulty
	if bytes.Compare(hash, target) > 0 {
		return false, errors.New("hash does not meet target difficulty")
	}

	return true, nil
}

// CheckSolution verifies a mining solution using RandomX
func CheckSolution(vm unsafe.Pointer, blockHeader []byte, nonce []byte, solution []byte, difficulty *big.Int) (bool, error) {
	if vm == nil {
		return false, errors.New("RandomX VM is not initialized")
	}

	// Combine block header and nonce
	input := append(blockHeader, nonce...)

	// Calculate the hash
	hash := CalculateHash(vm, input)

	// Compare the hash with the provided solution
	if !bytes.Equal(hash, solution) {
		return false, errors.New("solution does not match calculated hash")
	}

	// Convert hash to big.Int for comparison with difficulty
	hashInt := new(big.Int).SetBytes(hash)

	  // Calculate the maximum target (2^248 - 1)
	  // WARNING DONT FORGET TO ADJUST TO ETICA SMART CONTRACT DIIFICULTY BEFORE RELEASE:
	  // big.NewInt(248) should be replaced by big.NewInt(224)
	  maxTarget := new(big.Int).Exp(big.NewInt(2), big.NewInt(248), nil)
	  maxTarget.Sub(maxTarget, big.NewInt(1))
  
	  // Calculate the current target based on the difficulty
	  currentTarget := new(big.Int).Div(maxTarget, difficulty)
  
	  // Check if the hash is less than or equal to the current target
	  if hashInt.Cmp(currentTarget) > 0 {
		  return false, errors.New("hash does not meet target difficulty")
	  }

	return true, nil
}

// Helper function to convert hex string to bytes
func hexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
