package bencode

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// A bencode encoder
type encoder struct {
	buffer bytes.Buffer
}

// Writes data to a writer
func (e encoder) WriteTo(writer io.Writer) (n int64, err error) {
	return e.buffer.WriteTo(writer)
}

// Return a string representation of the encoded bencode
func (e encoder) String() string {
	return e.buffer.String()
}

// Return a bytes representation of the encoded bencode
func (e encoder) Bytes() []byte {
	return e.buffer.Bytes()
}

// Writes an integer to the encoder output
func (e *encoder) writeInt(i int64) {
	e.buffer.WriteRune('i')
	e.buffer.WriteString(strconv.FormatInt(i, 10))
	e.buffer.WriteRune('e')
}

// Writes an unsigned integer to the encoder output
func (e *encoder) writeUint(u uint64) {
	e.buffer.WriteRune('i')
	e.buffer.WriteString(strconv.FormatUint(u, 10))
	e.buffer.WriteRune('e')
}

// Writes a string to the encoder output
func (e *encoder) writeString(s string) {
	e.buffer.Grow(len(s) + 10)
	e.buffer.WriteString(strconv.FormatInt(int64(len(s)), 10))
	e.buffer.WriteRune(':')
	e.buffer.WriteString(s)
}

// Writes a map[string]interface to the encoder output
func (e *encoder) writeMap(m map[string]interface{}) error {
	e.buffer.WriteByte('d')
	for k, v := range m {
		e.writeString(k)
		if err := e.writeAuto(v); err != nil {
			return err
		}
	}
	e.buffer.WriteByte('e')
	return nil
}

// Writes a slice []interface to the encoder output
func (e *encoder) writeSlice(l []interface{}) error {
	e.buffer.WriteByte('l')
	for _, v := range l {
		if err := e.writeAuto(v); err != nil {
			return err
		}
	}
	e.buffer.WriteByte('e')
	return nil
}

// Automatically identify the type of an interface and write it to the encoder output
func (e *encoder) writeAuto(v interface{}) error {
	switch v := v.(type) {
	case string:
		e.writeString(v)
	case []interface{}:
		return e.writeSlice(v)
	case map[string]interface{}:
		return e.writeMap(v)
	case int:
		e.writeInt(int64(v))
	case uint:
		e.writeUint(uint64(v))
	case int64:
		e.writeInt(v)
	case uint64:
		e.writeUint(v)
	case int32:
		e.writeInt(int64(v))
	case uint32:
		e.writeUint(uint64(v))
	case int16:
		e.writeInt(int64(v))
	case uint16:
		e.writeUint(uint64(v))
	case int8:
		e.writeInt(int64(v))
	case uint8:
		e.writeUint(uint64(v))
	default:
		return fmt.Errorf("can't format interface of type %T: %v", v, v)
	}
	return nil
}

// Returns a bencode encoder from a given int64
func NewEncoderFromUint(u uint64) *encoder {
	e := new(encoder)
	e.writeUint(u)
	return e
}

// Returns a bencode encoder from a given int64
func NewEncoderFromInt(i int64) *encoder {
	e := new(encoder)
	e.writeInt(i)
	return e
}

// Returns a bencode encoder from a given string
func NewEncoderFromString(s string) *encoder {
	e := new(encoder)
	e.writeString(s)
	return e
}

// Returns a bencode encoder from a given slice
func NewEncoderFromSlice(l []interface{}) (*encoder, error) {
	e := new(encoder)
	return e, e.writeSlice(l)
}

// Returns a bencode encoder from a given map[string]interface{}
func NewEncoderFromMap(m map[string]interface{}) (*encoder, error) {
	e := new(encoder)
	return e, e.writeMap(m)
}

// Returns a bencode encoder from a given object
func NewEncoderFromInterface(v interface{}) (*encoder, error) {
	e := new(encoder)
	return e, e.writeAuto(v)
}
