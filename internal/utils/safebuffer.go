package utils

import (
	"bytes"
	"sync"
)

// Buffer as a replacement from bytes.Buffer
// We need to use this method because we need a safety buffer operations running inside goroutines
// ref: https://stackoverflow.com/questions/19646717/is-the-go-bytes-buffer-thread-safe
type Buffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *Buffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}
