package main

import "io"

type byteBuf struct {
	byteBuf []byte
	ptr     int
	size    int
}

func newByteBuf(maxBufSize int) *byteBuf {
	return &byteBuf{
		byteBuf: make([]byte, maxBufSize),
		ptr:     0,
		size:    0,
	}
}

func (buf *byteBuf) IsEmpty() bool {
	return buf.size == 0 || buf.ptr >= buf.size
}

func (buf *byteBuf) IsFull() bool {
	return buf.size >= cap(buf.byteBuf)
}

func (buf *byteBuf) ReadNextByte() byte {
	nextByte := buf.byteBuf[buf.ptr]
	buf.ptr += 1
	return nextByte
}

func (buf *byteBuf) WriteNextByte(next byte) {
	buf.byteBuf[buf.size] = next
	buf.size += 1
}

func (buf *byteBuf) Write(dest io.Writer) error {
	_, err := dest.Write(buf.byteBuf[:buf.size])
	buf.ptr = 0
	buf.size = 0
	return err
}

func (buf *byteBuf) Read(src io.Reader) bool {
	buf.ptr = 0
	buf.size, _ = io.ReadAtLeast(src, buf.byteBuf, cap(buf.byteBuf))
	return buf.size != 0
}

func (buf *byteBuf) Clear() {
	buf.ptr = 0
	buf.size = 0
}
