/*
  bool-in-uint32 atomic store and load.

  This type exists because of the aliasing of up to 3 other byte sized variables in
  the 4 bytes containing a regular Go 'bool' trips up the go -race detection, producing
  false races, because from Go's point of view we were reading and writing the 4 bytes,
  even though we were only, and atomically, changing only one.

  Copyright 2021 Nicolas S. Dade
*/
package atomicbool

import (
	"sync/atomic"
)

// Bool32 is a boolean stored in a uint32. It's sole purpose is to avoid go -race false
// positives which otherwise happen when a regular (1 byte) Go 'bool' is adjacent other
// variables in memory.
type Bool32 uint32

// Store atomically stores x in b
func (b *Bool32) Store(x bool) {
	var v uint32
	if x {
		v = 1
	}
	atomic.StoreUint32((*uint32)(b), v)
}

// Load atomically loads and returns the boolean value of b
func (b *Bool32) Load() bool {
	var x bool
	if atomic.LoadUint32((*uint32)(b)) != 0 {
		x = true
	}
	return x
}

// CompareAndSwap atomically performs:
//  if b == old {
//    b = new
//    return true
// } else {
//    return false
// }
func (b *Bool32) CompareAndSwap(old, new bool) (swapped bool) {
	var nw, od uint32
	if new {
		nw = 1
	}
	if old {
		od = 1
	}
	return atomic.CompareAndSwapUint32((*uint32)(b), od, nw)
}
