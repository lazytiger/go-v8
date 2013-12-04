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

type embedable struct {
	data interface{}
}

func (this embedable) GetPrivateData() interface{} {
	return this.data
}

func (this *embedable) SetPrivateData(data interface{}) {
	this.data = data
}
