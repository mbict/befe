package bufferpool

import (
	"net/http/httputil"
	"sync"
)

type syncBufferPool struct {
	pool sync.Pool
}

func New() httputil.BufferPool {
	return &syncBufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				// The Pool's New function should generally only return pointer
				// types, since a pointer can be put into the return interface
				// value without an allocation:
				//return new(bytes.Buffer)
				return make([]byte, 250*1024) //every request gets a default 250k byte buffer
			},
		},
	}
}

func (bp *syncBufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

func (bp *syncBufferPool) Put(buffer []byte) {
	bp.pool.Put(buffer)
}
