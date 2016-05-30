package main

import (
	"bufio"
	"io"
)

type P1 struct {
	rd *bufio.Reader
}

// NewReader returns a Reader whose read telegrams from a Reader object.
func NewP1(rd io.Reader) P1 {
	return P1{rd: bufio.NewReader(rd)}
}

// Read reads data into p. It returns the number of bytes read into p. It reads
// till it encounters a `!` byte.
func (p1 *P1) Read(p []byte) (n int, err error) {
	b, err := p1.rd.ReadBytes('!')
	p = append(p, b...)

	return len(b), err
}
