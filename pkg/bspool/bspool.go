package bspool

import (
	"bytes"
	"sync"
)

// bsPool is a pool of bytes.Buffer.
var bsPool = sync.Pool{New: func() any {
	return new(bytes.Buffer)
}}

// GetByteSlice returns a byte slice from the pool.
func GetByteSlice() *bytes.Buffer {
	v := bsPool.Get().(*bytes.Buffer)
	v.Reset()

	return v
}

// PutByteSlice returns a byte slice to the pool.
func PutByteSlice(bs *bytes.Buffer) {
	bsPool.Put(bs)
}
