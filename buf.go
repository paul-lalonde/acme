package main

import (
	"io/ioutil"
	"os"
)


// I'm chosing not to maintain nc directly, which is the number of elements
// in the buffer, but exposing Buffer.nc(), guaranteed to stay in sync.
type Buffer struct {
	buf      []rune
}

func (b *Buffer) Insert(q0 uint, r []rune) {
	if q0 > uint(len(b.buf)) {
		panic("internal error: buffer.Insert: Out of range insertion")
	}
	b.buf = append(b.buf[:q0], append(r, b.buf[q0:]...)...)
}

func (b *Buffer) Delete(q0, q1 uint) {
	if q0 > uint(len(b.buf)) || q1 > uint(len(b.buf)) {
		panic("internal error: buffer.Delete: Out-of-range Delete")
	}
	copy(b.buf[q0:], b.buf[q1:])
	b.buf = b.buf[:uint(len(b.buf))-(q1-q0)] // Reslice to length
}

func (b *Buffer) Load(q0 uint, fd *os.File) int {
	// TODO(flux): Innefficient to load the file, then copy into the slice,
	// but I need the UTF-8 interpretation.  I could fix this by using a
	// UTF-8 -> []rune reader on top of the os.File instead.

	d, err := ioutil.ReadAll(fd)
	if err != nil {
		warning(nil, "read error in Buffer.Load")
	}
	s := string(d)
	runes := []rune(s)
	b.Insert(q0, runes)
	return len(runes)
}

func (b *Buffer) Read(q0, n uint) (r []rune) {
	// TODO(flux): Can I just reslice here, or do I need to copy?
	if !(q0 <= uint(len(b.buf)) && q0+n <= uint(len(b.buf))) {
		panic("internal error: Buffer.Read")
	}

	return b.buf[q0:q0+n]
}

func (b *Buffer) Close() {
	b.Reset()

}

func (b *Buffer) Reset() {
	b.buf = b.buf[0:0]
}

func (b *Buffer) nc() uint {
	return uint(len(b.buf))
}