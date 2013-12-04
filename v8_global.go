package v8

import (
	"sync"
)

var (
	gAllocator *ArrayBufferAllocator
	gMutex     sync.Mutex
)

func init() {
	gAllocator = newArrayBufferAllocator()
}
