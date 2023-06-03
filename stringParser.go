package bencode

import (
	"io"
	"strconv"
	"strings"
)

// A bencode parser that uses a string as source
type stringParser struct {
	bencode string
	i       int
}

// Returns the next byte
func (sp *stringParser) readByte() (byte, error) {
	if sp.i >= len(sp.bencode) {
		return 0, io.EOF
	}
	sp.i++
	return sp.bencode[sp.i-1], nil
}

// Backtracks by 1 byte
//
// MUST only be called at most once immediately following a call to nextByte()
func (sp *stringParser) undoReadByte() {
	sp.i--
	if sp.i < 0 {
		panic("stringParser: invalid undoReadByte")
	}
}

// Reads a number until an end byte
func (sp *stringParser) readIntTo(separator byte) (int, error) {
	i := sp.i
	if i >= len(sp.bencode) {
		return 0, io.ErrUnexpectedEOF
	}
	index := strings.IndexByte(sp.bencode[i:], separator)
	if index == -1 {
		return 0, io.ErrUnexpectedEOF
	}
	sp.i += index + 1
	return strconv.Atoi(sp.bencode[i : i+index])
}

// Returns a string of a given length
func (sp *stringParser) readString(length int) (string, error) {
	sp.i += length
	if sp.i > len(sp.bencode) {
		return "", io.ErrUnexpectedEOF
	}
	return sp.bencode[sp.i-length : sp.i], nil
}

// Returns a bencode decoder from a given string
func NewParserFromString(bencode string) *decoder {
	return &decoder{
		&stringParser{
			bencode: bencode,
			i:       0,
		},
	}
}
