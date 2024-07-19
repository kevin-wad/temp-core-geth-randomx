package mutations

import (
	"encoding/binary"
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

func TestCheckSolutionWithRxSlowHash(t *testing.T) {
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
	blockHeaderHex := "101096a5a1b4061274d1d8e13640eff7416062d3366960171731b703b31244d20c252d090c9d97000000008f3f41a03692ea66f71676a3eae82c215be3347b447fd2545b0cfd2c7b850ad837"

	// Convert hex to bytes
	blockHeader, err := hex.DecodeString(blockHeaderHex)
	if err != nil {
		fmt.Printf("Error decoding block header: %v\n", err)
		return
	}

	fmt.Printf("--- > Block header (bytes): %v\n", blockHeader)
	fmt.Printf("--- > Block header length: %d bytes\n", len(blockHeader))
	fmt.Printf("Block header: %x\n", blockHeader)

	// Use an integer nonce instead of a hex string
	nonce := uint32(229382)
	fmt.Printf("Nonce (integer): %d\n", nonce)

	// Convert the nonce to a 4-byte slice (little-endian)
	nonceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceBytes, nonce)
	fmt.Printf("Nonce (bytes): %x\n", nonceBytes)

	// Create a copy of the block header and insert the nonce at the correct offset
	blobWithNonce := make([]byte, len(blockHeader))
	copy(blobWithNonce, blockHeader)
	copy(blobWithNonce[nonceOffset:], nonceBytes)
	// Log the original blob and the blob with the new nonce
	fmt.Printf("Original Blob: %x\n", blockHeader)
	fmt.Printf("Blob with Nonce: %x\n", blobWithNonce)

	blockHeight := uint64(3182000) // WARNING: use hardcoded value for tests need to implement get it from tx inputs
	seedHeight := RxSeedHeight(blockHeight)
	fmt.Printf("seedHeight: %x\n", seedHeight)

	// Use the difficulty value we calculated earlier
	difficulty := big.NewInt(480045)
	fmt.Printf("Difficulty: %s\n", difficulty.String())

	// Calculate max target
	maxTarget := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	maxTarget.Sub(maxTarget, big.NewInt(1))

	// Calculate current target based on difficulty
	currentTarget := new(big.Int).Div(maxTarget, difficulty)
	fmt.Printf("Current target: %x\n", currentTarget.Bytes())

	// Seed hash (as used in VerifyEticaTransaction)
	seedHashString := "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da788"
	seedHash, _ := hex.DecodeString(seedHashString)
	fmt.Printf("Seed hash: %x\n", seedHash)

	mhash := calculateRandomXHash(blobWithNonce, seedHash)
	fmt.Printf("Calculated mhash is: %x\n", mhash)

	// Test with correct solution
	fmt.Println("Testing with correct solution")
	valid, err := CheckSolutionWithRxSlowHash(vm, blobWithNonce, mhash, difficulty, blockHeight, seedHash)
	if !valid || err != nil {
		t.Errorf("Correct solution check failed: %v", err)
	}
	fmt.Printf("Correct solution check result: valid=%v, err=%v\n", valid, err)

	seedHashString2 := "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da799"
	seedHash2, _ := hex.DecodeString(seedHashString2)
	fmt.Printf("Seedhash2: %x\n", seedHash2)

	mhash2 := calculateRandomXHash(blobWithNonce, seedHash2)

	// Test with solution not meeting target
	fmt.Println("Testing with solution not meeting target")
	hardDifficulty := new(big.Int).Mul(difficulty, big.NewInt(1000)) // 1000 times harder
	valid, err = CheckSolutionWithRxSlowHash(vm, blobWithNonce, mhash2, hardDifficulty, blockHeight, seedHash)
	if valid || err == nil {
		t.Error("Target difficulty check failed")
	}
	fmt.Printf("Hard difficulty check result: valid=%v, err=%v\n", valid, err)

	fmt.Println("TestCheckSolution completed")
}

/*
func TestCheckSolutionWithRxSlowHash(t *testing.T) {
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
	blockHeaderHex := "101096a5a1b4061274d1d8e13640eff7416062d3366960171731b703b31244d20c252d090c9d97000000008f3f41a03692ea66f71676a3eae82c215be3347b447fd2545b0cfd2c7b850ad837"

	// Convert hex to bytes
	blockHeader, err := hex.DecodeString(blockHeaderHex)
	if err != nil {
		fmt.Printf("Error decoding block header: %v\n", err)
		return
	}

	// Verify the length (68 bytes as we calculated earlier), not sure since nonce can have variable length
	/* expectedLength := 68
	 if len(blockHeader) != expectedLength {
		 fmt.Printf("Invalid block header length: expected %d bytes, got %d\n", expectedLength, len(blockHeader))
		 return
	 } * /

	fmt.Printf("--- > Block header (bytes): %v\n", blockHeader)
	fmt.Printf("--- > Block header length: %d bytes\n", len(blockHeader))

	fmt.Printf("Block header: %x\n", blockHeader)

	/*
		nonce, err := hexToBytes("00000100")
		if err != nil {
			fmt.Printf("Error decoding nonce: %v\n", err)
			return
		}
		fmt.Printf("Nonce: %x\n", nonce)

		// Create a copy of the block header and insert the nonce at the correct offset
		blobWithNonce := make([]byte, len(blockHeader))
		copy(blobWithNonce, blockHeader)
		copy(blobWithNonce[nonceOffset:], nonce) * /

	// Use an integer nonce instead of a hex string
	nonce := uint32(65536)
	fmt.Printf("Nonce (integer): %d\n", nonce)

	// Convert the nonce to a 4-byte slice (little-endian)
	nonceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceBytes, nonce)
	fmt.Printf("Nonce (bytes): %x\n", nonceBytes)

	// Create a copy of the block header and insert the nonce at the correct offset
	blobWithNonce := make([]byte, len(blockHeader))
	copy(blobWithNonce, blockHeader)
	copy(blobWithNonce[nonceOffset:], nonceBytes)
	// Log the original blob and the blob with the new nonce
	fmt.Printf("Original Blob: %x\n", blockHeader)
	fmt.Printf("Blob with Nonce: %x\n", blobWithNonce)

	/* EXPECTED RESULT
	XMRig params are: {
	321|index  |   id: '1',
	321|index  |   job_id: '83510',
	321|index  |   nonce: '06800300',
	321|index  |   result: '8f2460d90ef6a1b5a0d7e2fa53f4d8e461dd661eccc25f9993f77131ab79f557'
	321|index  | }
	* /

	blockHeight := uint64(3182000) // WARNING: use hardcoded value for tests need to implement get it from tx inputs
	seedHeight := RxSeedHeight(blockHeight)
	fmt.Printf("seedHeight: %x\n", seedHeight)

	// Use the difficulty value we calculated earlier
	difficulty := big.NewInt(480045)
	fmt.Printf("Difficulty: %s\n", difficulty.String())

	// Calculate max target
	maxTarget := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	maxTarget.Sub(maxTarget, big.NewInt(1))

	// Calculate current target based on difficulty
	currentTarget := new(big.Int).Div(maxTarget, difficulty)
	fmt.Printf("Current target: %x\n", currentTarget.Bytes())

	// Seed hash (as used in VerifyEticaTransaction)
	seedHashString := "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da788"
	seedHash, _ := hex.DecodeString(seedHashString)
	fmt.Printf("Seed hash: %x\n", seedHash)

	seedHashString2 := "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da799"
	seedHash2, _ := hex.DecodeString(seedHashString2)
	fmt.Printf("Seedhash2: %x\n", seedHash2)

	OxseedHashString := "0x25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da788"
	OxseedHash, _ := hex.DecodeString(OxseedHashString)
	fmt.Printf("OxSeed hash: %x\n", OxseedHash)

	correctSolution := RxSlowHash(blockHeight, seedHeight, seedHash, blobWithNonce)
	fmt.Printf("RxSlowHash result hash: %x\n", correctSolution)
	fmt.Printf("correctSolution: %x\n", correctSolution)

	correctSolution2 := RxSlowHash(1000, 500, seedHash, blobWithNonce)
	fmt.Printf("RxSlowHash result2 hash: %x\n", correctSolution2)
	fmt.Printf("correctSolution2: %x\n", correctSolution2)

	correctSolution3 := RxSlowHash(blockHeight, seedHeight, seedHash2, blobWithNonce)
	fmt.Printf("RxSlowHash result3 hash: %x\n", correctSolution3)
	fmt.Printf("correctSolution3: %x\n", correctSolution3)

	correctSolution4 := RxSlowHash(1000, 500, seedHash2, blobWithNonce)
	fmt.Printf("RxSlowHash result2 hash: %x\n", correctSolution4)
	fmt.Printf("correctSolution2: %x\n", correctSolution4)

	mhash := calculateRandomXHash(blobWithNonce, seedHash)
	fmt.Printf("Calculated mhash is: %x\n", mhash)

	mhash2 := calculateRandomXHash(blobWithNonce, seedHash2)
	fmt.Printf("Calculated mhash2 is: %x\n", mhash2)

	/*mhash := calculateRandomXHash(blockHeader, nonce, seedHash)
	// Print the calculated hash
	fmt.Printf("Calculated mhash is: %x\n", mhash) */

/*
		OcorrectSolution := RxSlowHash(blockHeight, seedHeight, seedHash, OxtestInput)
		fmt.Printf("RxSlowHash with 0xtestInput result hash: %x\n", OcorrectSolution)
		fmt.Printf("OcorrectSolution: %x\n", OcorrectSolution)

		OcorrectSolutio0XseedHash := RxSlowHash(blockHeight, seedHeight, OxseedHash, OxtestInput)
		fmt.Printf("RxSlowHash with 0xtestInput and OxseedHash result hash: %x\n", OcorrectSolutio0XseedHash)
		fmt.Printf("OcorrectSolutio0XseedHash: %x\n", OcorrectSolutio0XseedHash)

		correctSolutio0XseedHash := RxSlowHash(blockHeight, seedHeight, OxseedHash, testInput)
		fmt.Printf("RxSlowHash with 0xtestInput and OxseedHash result hash: %x\n", correctSolutio0XseedHash)
		fmt.Printf("correctSolutio0XseedHash: %x\n", correctSolutio0XseedHash) * /

	// Test with correct solution
	fmt.Println("Testing with correct solution")
	valid, err := CheckSolutionWithRxSlowHash(vm, blobWithNonce, correctSolution, difficulty, blockHeight, seedHash)
	if !valid || err != nil {
		t.Errorf("Correct solution check failed: %v", err)
	}
	fmt.Printf("Correct solution check result: valid=%v, err=%v\n", valid, err)

	// Test with incorrect solution
	fmt.Println("Testing with incorrect solution")
	incorrectSolution := make([]byte, len(correctSolution))
	copy(incorrectSolution, correctSolution)
	incorrectSolution[0] ^= 0xff // Flip some bits to make it incorrect
	valid, err = CheckSolutionWithRxSlowHash(vm, blobWithNonce, incorrectSolution, difficulty, blockHeight, seedHash)
	if valid || err == nil {
		t.Error("Incorrect solution check failed")
	}
	fmt.Printf("Incorrect solution check result: valid=%v, err=%v\n", valid, err)

	// Test with solution not meeting target
	fmt.Println("Testing with solution not meeting target")
	hardDifficulty := new(big.Int).Mul(difficulty, big.NewInt(1000)) // 1000 times harder
	valid, err = CheckSolutionWithRxSlowHash(vm, blobWithNonce, incorrectSolution, hardDifficulty, blockHeight, seedHash)
	if valid || err == nil {
		t.Error("Target difficulty check failed")
	}
	fmt.Printf("Hard difficulty check result: valid=%v, err=%v\n", valid, err)

	fmt.Println("TestCheckSolution completed")
} */

/*
// Function to initialize RandomX cache and VM, and calculate the hash
func calculateRandomXHash(blockHeader, nonce, seedHash []byte) []byte {
	combinedInput := append(blockHeader, nonce...)

	flags := FlagDefault
	cache := InitRandomX(flags)
	if cache == nil {
		panic("Failed to allocate RandomX cache")
	}
	defer DestroyRandomX(cache)

	InitCache(cache, seedHash)

	vm := CreateVM(cache, flags)
	if vm == nil {
		panic("Failed to create RandomX VM")
	}
	defer DestroyVM(vm)

	/*CalculateHashFirst(vm, combinedInput)
	hash := CalculateHashLast(vm) * /

	//CalculateHashFirst(vm, combinedInput)
	//hash := CalculateHashLast(vm)
	// OR
	hash := CalculateHash(vm, combinedInput)

	return hash
} */

func TestCalculateRandomXHash(t *testing.T) {
	fmt.Println("Starting TestCalculateRandomXHash")

	// Provided blockHeader and nonce
	blockHeaderHex := "101096a5a1b4061274d1d8e13640eff7416062d3366960171731b703b31244d20c252d090c9d97000000008f3f41a03692ea66f71676a3eae82c215be3347b447fd2545b0cfd2c7b850ad837"
	nonceHex := "06800300"

	blockHeader, err := hexToBytes(blockHeaderHex)
	if err != nil {
		t.Fatalf("Error decoding block header: %v", err)
	}

	nonce, err := hexToBytes(nonceHex)
	if err != nil {
		t.Fatalf("Error decoding nonce: %v", err)
	}

	// Provided seedHash
	seedHashHex := "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da788"
	seedHash, err := hexToBytes(seedHashHex)
	if err != nil {
		t.Fatalf("Error decoding seed hash: %v", err)
	}

	// Create a copy of the block header and insert the nonce at the correct offset
	blobWithNonce := make([]byte, len(blockHeader))
	copy(blobWithNonce, blockHeader)
	copy(blobWithNonce[nonceOffset:], nonce)

	// Log the original blob and the blob with the new nonce
	fmt.Printf("Original Blob: %x\n", blockHeader)
	fmt.Printf("Blob with Nonce: %x\n", blobWithNonce)

	// Calculate the RandomX hash
	hash := calculateRandomXHash(blobWithNonce, seedHash)

	// Print the calculated hash
	fmt.Printf("Calculated hash: %x\n", hash)

	expectedHashHex := "8f2460d90ef6a1b5a0d7e2fa53f4d8e461dd661eccc25f9993f77131ab79f557"
	expectedHash, err := hexToBytes(expectedHashHex)
	if err != nil {
		t.Fatalf("Error decoding expected hash: %v", err)
	}

	if !equal(hash, expectedHash) {
		t.Errorf("Hash mismatch: got %x, want %x", hash, expectedHash)
	}
	fmt.Println("TestCalculateRandomXHash completed")
}

// Helper function to compare two byte slices
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Test function for computePowHash
func TestComputePowHash(t *testing.T) {
	fmt.Println("Starting TestComputePowHash")

	// Provided blockHeader and nonce
	blockHeaderHex := "101096a5a1b4061274d1d8e13640eff7416062d3366960171731b703b31244d20c252d090c9d97000000008f3f41a03692ea66f71676a3eae82c215be3347b447fd2545b0cfd2c7b850ad837"
	nonceHex := "06800300"

	blockHeader, err := hexToBytes(blockHeaderHex)
	if err != nil {
		t.Fatalf("Error decoding block header: %v", err)
	}

	nonce, err := hexToBytes(nonceHex)
	if err != nil {
		t.Fatalf("Error decoding nonce: %v", err)
	}

	// Provided seedHash
	seedHashHex := "25314901c96d26ff28484bddf315f0a3295f30f13590d056efd65fcb6d8da788"
	seedHash, err := hexToBytes(seedHashHex)
	if err != nil {
		t.Fatalf("Error decoding seed hash: %v", err)
	}

	// Combine block header and nonce to create block blob
	blockBlob := append(blockHeader, nonce...)

	// Define height and major version
	height := uint64(3182000)
	majorVersion := RX_BLOCK_VERSION

	// Check block blob size
	err = CheckIncomingBlockSize(blockBlob)
	if err != nil {
		t.Fatalf("Block blob size check failed: %v", err)
	}

	// Compute the PoW hash using RandomX
	hash, err := ComputePowHash(blockBlob, seedHash, height, majorVersion)
	if err != nil {
		t.Fatalf("Error computing PoW hash: %v", err)
	}

	// Print the calculated hash
	fmt.Printf("Calculated PoW hash: %x\n", hash)

	expectedHashHex := "8f2460d90ef6a1b5a0d7e2fa53f4d8e461dd661eccc25f9993f77131ab79f557"
	expectedHash, err := hexToBytes(expectedHashHex)
	if err != nil {
		t.Fatalf("Error decoding expected hash: %v", err)
	}

	if !equal(hash, expectedHash) {
		t.Errorf("Hash mismatch: got %x, want %x", hash, expectedHash)
	}
	fmt.Println("TestComputePowHash completed")
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
