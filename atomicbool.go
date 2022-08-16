/*
Atomic store and load of bool type.

Copyright 2017 Nicolas S. Dade
*/
package atomicbool

import (
	"sync/atomic"
	"unsafe"

	"github.com/nsd20463/cpuendian"
)

// sanity check our assumptions
func init() {
	if unsafe.Sizeof(false) != 1 {
		panic("nsd20463/atomicbool assumes bools fit in 1 byte")
	}
}

func StoreBool(addr *bool, val bool) {
	// figure out what uint32 this bool is part of, and edit the uint32
	// NOTE WELL we have to do this in such a way that gc moving the bool around will update our local vars,
	// which in turn means we can't store anything in a uintptr type except in the middle of an expression.
	// We're also going to make the (currently safe) assumption that the alignment of the bool within the uint32
	// will not change if gc moves the bool.
	p32 := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) &^ 3))
	shift := 8 * (uintptr(unsafe.Pointer(addr)) & 3)
	if cpuendian.Big {
		shift = 32 - 8 - shift
	}

	for {
		i := atomic.LoadUint32(p32)
		n := i &^ (0xff << shift)
		if val {
			n |= 1 << shift
		}
		if atomic.CompareAndSwapUint32(p32, i, n) {
			return
		}
	}
}

func LoadBool(addr *bool) (val bool) {
	// see comments in StoreBool
	p32 := (*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(addr)) &^ 3))
	shift := 8 * (uintptr(unsafe.Pointer(addr)) & 3)
	if cpuendian.Big {
		shift = 32 - 8 - shift
	}

	i := atomic.LoadUint32(p32)
	i >>= shift
	i &= 0xff
	return i != 0
}
