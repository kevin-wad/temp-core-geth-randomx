package mutations

import (
	"fmt"
	"sync"
	"unsafe"
)

var (
	initOnce   sync.Once
	mainCache  unsafe.Pointer
	mainVM     unsafe.Pointer
	cacheMutex sync.RWMutex
	vmMutex    sync.RWMutex
)

const (
	SEEDHASH_EPOCH_BLOCKS = 2048
	SEEDHASH_EPOCH_LAG    = 64
)

func initCache(seedHash []byte) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if mainCache != nil {
		DestroyRandomX(mainCache)
	}

	mainCache = InitRandomX(FlagDefault)
	if mainCache == nil {
		panic("Failed to initialize RandomX cache")
	}

	InitCache(mainCache, seedHash)
}

func initVM() {
	vmMutex.Lock()
	defer vmMutex.Unlock()

	if mainVM != nil {
		DestroyVM(mainVM)
	}

	mainVM = CreateVM(mainCache, FlagDefault)
	if mainVM == nil {
		panic("Failed to create RandomX VM")
	}
}

func RxSlowHash(mainHeight uint64, seedHeight uint64, seedHash []byte, data []byte) []byte {

	fmt.Printf("! ------------------- RxSlowHash called ---------------- !")

	initOnce.Do(func() {
		fmt.Printf("! ------------------- RxSlowHash initOnce called ---------------- !")
		dummySeedHash := make([]byte, 32)
		initCache(dummySeedHash)
		initVM()
	})

	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	vmMutex.RLock()
	defer vmMutex.RUnlock()

	if len(seedHash) != 32 {
		panic("Invalid seed hash length")
	}

	if mainCache == nil || len(seedHash) != 32 {
		fmt.Printf("! ------------------- RxSlowHash mainCache == nil || len(seedHash) != 32 assed true ---------------- !")
		initCache(seedHash)
		initVM()
	}

	// Check if we need to reinitialize the cache with the new seed hash
	// Future Code:
	/*if !bytes.Equal(seedHash, getCurrentSeedHash()) {
	    initCache(seedHash)
	    initVM()
	} */
	fmt.Printf("! ------------------- RxSlowHash reach call to CalculateHash ---------------- !")
	return CalculateHash(mainVM, data)
}

/* Future Code:
func getCurrentSeedHash() []byte {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if mainCache == nil {
		return nil
	}

	// Assuming there's a way to get the current seed hash from the cache
	// You might need to modify the RandomX C code to expose this functionality
	// For now, let's assume there's a function to do this
	return C.randomx_get_current_seed_hash((*C.struct_randomx_cache)(mainCache))
}
*/
