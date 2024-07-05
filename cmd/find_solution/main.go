package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/mutations"
	// Adjust this import path as necessary
)

/*
func main() {
	blockHeader, err := hex.DecodeString("0123456789abcdef0123456789abcdef")
	if err != nil {
		log.Fatalf("Failed to decode block header: %v", err)
	}

	target, err := hex.DecodeString("00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	if err != nil {
		log.Fatalf("Failed to decode target: %v", err)
	}

	nonce, hash, err := mutations.FindValidSolution(blockHeader, target)
	if err != nil {
		log.Fatalf("Failed to find valid solution: %v", err)
	}

	fmt.Println("Found valid solution!")
	fmt.Printf("Nonce: %x\n", nonce)
	fmt.Printf("Hash: %x\n", hash)
} */

func main() {
	challengeNumber := "0xa76de21cc71ed487f47d0fe296c0cf71360bc86b8610975f6c43b2aef7c6818d"
	difficulty := new(big.Int)
	difficulty.SetString("452312848583266388373324160190187140051835877600158453279131187530910662656", 10)
	minerAddress := common.HexToAddress("0x2434e3552573c1B3FC55AEC409A1691dF840CAcb") // Replace with your address

	nonce, randomXHash, err := mutations.FindValidSolution2(challengeNumber, difficulty, minerAddress)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found valid solution!\n")
	fmt.Printf("BlockHeader: 0x%s\n", challengeNumber)
	fmt.Printf("Nonce: %s\n", nonce.String())
	fmt.Printf("RandomX Hash: 0x%x\n", randomXHash)

	// Verify the solution meets the difficulty
	digestInt := new(big.Int).SetBytes(randomXHash)
	if digestInt.Cmp(difficulty) <= 0 {
		fmt.Println("Solution is valid and meets the difficulty requirement")
	} else {
		fmt.Println("Solution does not meet the difficulty requirement")
	}
}
