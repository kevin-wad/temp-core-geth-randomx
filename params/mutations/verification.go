package mutations

import (
	"errors"
	"fmt"
)

// Define constants with high enough default values
const (
	MaxBlockWeightLimit      = 1000000 // Adjust this value as needed
	BLOCK_SIZE_SANITY_LEEWAY = 100000  // Adjust this value as needed
	RX_BLOCK_VERSION         = 10      // Adjust this value to the current or desired block version
)

// Check block blob size
func CheckIncomingBlockSize(blockBlob []byte) error {
	maxSize := MaxBlockWeightLimit + BLOCK_SIZE_SANITY_LEEWAY
	if len(blockBlob) > maxSize {
		return fmt.Errorf("block blob size is too big: %d", len(blockBlob))
	}
	return nil
}

// Compute PoW hash using RandomX
func ComputePowHash(blockBlob []byte, seedHash []byte, height uint64, majorVersion int) ([]byte, error) {
	var hash []byte

	if majorVersion >= RX_BLOCK_VERSION {
		// Use RandomX to compute the hash
		flags := FlagDefault
		cache := InitRandomX(flags)
		if cache == nil {
			return nil, errors.New("failed to allocate RandomX cache")
		}
		defer DestroyRandomX(cache)

		InitCache(cache, seedHash)
		vm := CreateVM(cache, flags)
		if vm == nil {
			return nil, errors.New("failed to create RandomX VM")
		}
		defer DestroyVM(vm)

		CalculateHashFirst(vm, blockBlob)
		hash = CalculateHashLast(vm)
	} else {
		// Use CryptoNight to compute the hash (not implemented in this example)
		// hash = cryptoNightHash(blockBlob, majorVersion, height)
	}

	return hash, nil
}
