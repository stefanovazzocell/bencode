package bencode

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	bufferSize    = 1024
	minBufferSize = 22 // The minimum buffer size to ~ store a number
	maxEmptyReads = 100
)

var (
	ErrNegativeRead = errors.New("readerParser: reader returned a negative read")
)

// A bencode reader that uses an io.Reader as source
type readerParser struct {
	buffer [bufferSize]byte
	reader io.Reader
	s, e   int
}

// Returns the length of the buffer
func (rp *readerParser) buffered() int {
	return rp.e - rp.s
}

// Fills the buffer
func (rp *readerParser) fill() error {
	b := rp.buffered()
	if b == bufferSize {
		return nil // Already full
	}
	// Swift current buffer
	if b > 0 {
		copy(rp.buffer[:], rp.buffer[rp.s:rp.e])
		rp.e -= rp.s
		rp.s = 0
	} else {
		rp.e = 0
	}
	rp.s = 0
	// Copy new data
	for i := 0; i < maxEmptyReads; i++ {
		n, err := rp.reader.Read(rp.buffer[rp.e:])
		rp.e += n
		if err == io.EOF && rp.e > 0 || rp.e == bufferSize {
			return nil
		}
		if err != nil {
			return err
		}
		if n < 0 {
			panic(ErrNegativeRead)
		}
	}
	return io.ErrNoProgress
}

// Returns the next byte
func (rp *readerParser) readByte() (byte, error) {
	// Try to buffer value
	if rp.buffered() < 1 {
		if err := rp.fill(); err != nil {
			return 0, err
		} // else we have at least 1 item in the buffer
	}
	// Read from buffer
	rp.s++
	return rp.buffer[rp.s-1], nil
}

// Backtracks by 1 byte
//
// MUST only be called at most once immediately following a call to nextByte()
func (rp *readerParser) undoReadByte() {
	rp.s--
	if rp.s < 0 {
		panic("readerParser: invalid undoReadByte")
	}
}

// Reads a number until an end byte
func (rp *readerParser) readIntTo(separator byte) (int, error) {
	// Try to buffer value
	if rp.buffered() < minBufferSize {
		if err := rp.fill(); err != nil {
			return 0, err
		} // else we have at least 1 item in the buffer
	}
	// Lookup the integer in the buffer
	index := bytes.IndexByte(rp.buffer[rp.s:rp.e], separator)
	if index == -1 {
		return 0, io.ErrUnexpectedEOF
	}
	n, err := strconv.Atoi(string(rp.buffer[rp.s : rp.s+index]))
	rp.s += index + 1
	return n, err
}

// Returns a string of a given length
func (rp *readerParser) readString(length int) (string, error) {
	// Try to get from buffer
	lb := rp.buffered()
	if rp.buffered() >= length {
		lb = length
		rp.s += lb
		return string(rp.buffer[rp.s-lb : rp.s]), nil
	}
	// Build from scratch
	i := 0
	strBuf := make([]byte, length)
	// Fetch beginning from buffer
	if lb > 0 {
		copy(strBuf, rp.buffer[rp.s:rp.e])
		rp.s += lb
		i += lb
	}
	// Fetch rest dynamically
	countEmpty := 0
	for countEmpty <= maxEmptyReads {
		n, err := rp.reader.Read(strBuf[i:])
		i += n
		if i == length && (err == nil || err != io.EOF) {
			return string(strBuf), nil
		}
		if err != nil {
			return "", err
		}
		if n < 0 {
			panic(ErrNegativeRead)
		}
		if n == 0 {
			countEmpty += 1
		} else {
			countEmpty = 0
		}
	}
	return "", io.ErrNoProgress
}

// Returns a bencode decoder from a given string
func NewParserFromReader(reader io.Reader) *decoder {
	return &decoder{
		&readerParser{
			buffer: [bufferSize]byte{},
			reader: reader,
			s:      0,
			e:      0,
		},
	}
}
