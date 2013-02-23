// Package leakbuf provides leaky buffer.
// It's based on the example in Effective Go.
package leakybuf

type LeakyBuf struct {
	bufSize  int // size of each buffer
	freeList chan []byte
}

// NewLeakyBuf creates a leaky buffer which can hold at most n buffer, each
// with bufSize bytes.
func NewLeakyBuf(n, bufSize int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

// Get returns a buffer from the leaky buffer or create a new buffer.
func (lb *LeakyBuf) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, lb.bufSize)
	}
	return
}

// Put add the buffer into the free buffer pool for reuse. If the buffer size
// is not the same with the leaky buffer's, Put will simply return and do not
// add the buffer.
func (lb *LeakyBuf) Put(b []byte) {
	if len(b) != lb.bufSize {
		return
	}
	select {
	case lb.freeList <- b:
	default:
	}
	return
}
