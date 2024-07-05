package mutations

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

/*
func TestRandomX(t *testing.T) {
	// Initialize RandomX
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		t.Fatal("Failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)

	// Create VM
	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		t.Fatal("Failed to create RandomX VM")
	}
	defer DestroyVM(vm)

	// Test hash calculation
	input := []byte("RandomX test input")
	hash := CalculateHash(vm, input)

	// Print the hash
	t.Logf("RandomX hash: %s", hex.EncodeToString(hash))

	// Basic sanity check
	if len(hash) != HashSize() {
		t.Errorf("Unexpected hash size: got %d, want %d", len(hash), HashSize())
	}
} */

func init() {
	fmt.Println("Initializing randomx_test.go")
}

func TestRandomX(t *testing.T) {
	fmt.Println("Starting TestRandomX")
	// Initialize RandomX
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		t.Fatal("Failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)
	fmt.Println("RandomX cache initialized")

	// Create VM
	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		t.Fatal("Failed to create RandomX VM")
	}
	defer DestroyVM(vm)
	fmt.Println("RandomX VM created")

	// Test hash calculation
	input := []byte("RandomX test input")
	fmt.Printf("Input for hash calculation: %s\n", string(input))
	hash := CalculateHash(vm, input)

	// Print the hash
	fmt.Printf("RandomX hash: %s\n", hex.EncodeToString(hash))

	// Basic sanity check
	if len(hash) != HashSize() {
		t.Errorf("Unexpected hash size: got %d, want %d", len(hash), HashSize())
	}
	fmt.Println("TestRandomX completed")
}

func TestCheckSolution(t *testing.T) {
	fmt.Println("Starting TestCheckSolution")
	// Initialize RandomX
	cache := InitRandomX(FlagDefault)
	if cache == nil {
		t.Fatal("Failed to initialize RandomX cache")
	}
	defer DestroyRandomX(cache)
	fmt.Println("RandomX cache initialized for TestCheckSolution")

	// Create VM
	vm := CreateVM(cache, FlagDefault)
	if vm == nil {
		t.Fatal("Failed to create RandomX VM")
	}
	defer DestroyVM(vm)
	fmt.Println("RandomX VM created for TestCheckSolution")

	// Test data
	blockHeader, _ := hexToBytes("0123456789abcdef0123456789abcdef")
	fmt.Printf("Block header: %x\n", blockHeader)
	nonce, _ := hexToBytes("2d9093434ce50229")
	fmt.Printf("Nonce: %x\n", nonce)
	target, _ := hexToBytes("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	fmt.Printf("Target: %x\n", target)

	// Fake "correct" solution that meets the target
	fakeCorrectSolution, _ := hexToBytes("00774e65ac49fb30a1d8271ebb9a3ebb97ba6f1381681fc79d2c62e09bd44433")
	fmt.Printf("Fake correct solution: %x\n", fakeCorrectSolution)

	// Test with fake correct solution
	fmt.Println("Testing with fake correct solution")
	valid, err := CheckSolutionWithTarget(vm, blockHeader, nonce, fakeCorrectSolution, target)
	if !valid || err != nil {
		t.Errorf("Fake solution check failed: %v", err)
	}
	fmt.Printf("Fake solution check result: valid=%v, err=%v\n", valid, err)

	// Calculate the correct solution
	input := append(blockHeader, nonce...)
	fmt.Printf("Input for solution: %x\n", input)
	correctSolution := CalculateHash(vm, input)
	fmt.Printf("Input solution: %x\n", correctSolution)

	// Test with input solution
	fmt.Println("Testing with input solution")
	valid2, err := CheckSolutionWithTarget(vm, blockHeader, nonce, correctSolution, target)
	if !valid2 || err != nil {
		t.Errorf("Input solution check failed: %v", err)
	}
	fmt.Printf("Input solution check result: valid=%v, err=%v\n", valid2, err)

	// Test with incorrect solution
	fmt.Println("Testing with incorrect solution")
	incorrectSolution := make([]byte, len(correctSolution))
	copy(incorrectSolution, correctSolution)
	incorrectSolution[0] ^= 0xff // Flip some bits to make it incorrect
	fmt.Printf("Incorrect solution: %x\n", incorrectSolution)
	valid, err = CheckSolutionWithTarget(vm, blockHeader, nonce, incorrectSolution, target)
	if valid || err == nil {
		t.Error("Incorrect solution check failed")
	}
	fmt.Printf("Incorrect solution check result: valid=%v, err=%v\n", valid, err)

	// Test with solution not meeting target
	fmt.Println("Testing with solution not meeting target")
	hardTarget, _ := hexToBytes("0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	fmt.Printf("Hard target: %x\n", hardTarget)
	valid, err = CheckSolutionWithTarget(vm, blockHeader, nonce, correctSolution, hardTarget)
	if valid || err == nil {
		t.Error("Target difficulty check failed")
	}
	fmt.Printf("Hard target check result: valid=%v, err=%v\n", valid, err)

	fmt.Println("TestCheckSolution completed")
}

func TestCalculateDiggest(t *testing.T) {

	challengeNumber := "someChallenge"
	sender := common.HexToAddress("0x1234567890123456789012345678901234567890")
	nonce := big.NewInt(12345)

	digest := CalculateDigest(challengeNumber, sender, nonce)

	fmt.Printf("Digest: 0x%x\n", digest)

}
