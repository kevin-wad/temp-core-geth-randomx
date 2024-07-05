package mutations

/*
#cgo CFLAGS: -I./include
#cgo linux LDFLAGS: -L./lib -lrandomx -lstdc++
#cgo darwin LDFLAGS: -L./lib -lrandomx -lstdc++
#cgo windows LDFLAGS: -L./lib -lrandomx -lstdc++ -lws2_32 -ladvapi32
#include "randomx.h"
*/
import "C"

import (
	"unsafe"
)

// RandomXFlags type to represent RandomX flags
type RandomXFlags uint32

const (
	// FlagDefault is the default RandomX flag
	FlagDefault RandomXFlags = C.RANDOMX_FLAG_DEFAULT
)

// InitRandomX initializes RandomX with the given flags
func InitRandomX(flags RandomXFlags) unsafe.Pointer {
	return unsafe.Pointer(C.randomx_alloc_cache(C.randomx_flags(flags)))
}

// DestroyRandomX frees the RandomX cache
func DestroyRandomX(cache unsafe.Pointer) {
	C.randomx_release_cache((*C.struct_randomx_cache)(cache))
}

// CreateVM creates a new RandomX VM instance
func CreateVM(cache unsafe.Pointer, flags RandomXFlags) unsafe.Pointer {
	return unsafe.Pointer(C.randomx_create_vm(C.randomx_flags(flags), (*C.struct_randomx_cache)(cache), nil))
}

// DestroyVM destroys a RandomX VM instance
func DestroyVM(vm unsafe.Pointer) {
	C.randomx_destroy_vm((*C.struct_randomx_vm)(vm))
}

// CalculateHash calculates a RandomX hash
func CalculateHash(vm unsafe.Pointer, input []byte) []byte {
	output := make([]byte, C.RANDOMX_HASH_SIZE)
	C.randomx_calculate_hash((*C.struct_randomx_vm)(vm), unsafe.Pointer(&input[0]), C.size_t(len(input)), unsafe.Pointer(&output[0]))
	return output
}

// HashSize returns the size of the RandomX hash
func HashSize() int {
	return int(C.RANDOMX_HASH_SIZE)
}
