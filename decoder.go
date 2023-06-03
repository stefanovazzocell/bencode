package bencode

import "errors"

var (
	// Error for when the user tries to read a specific type of bencode
	// but the bencode doesn't provide that type
	ErrInvalidType = errors.New("invalid bencode output type")
	// Error indicating an invalid string length such as "-1:"
	ErrInvalidStringLen = errors.New("invalid bencode: string length can't be negative")
)

// A generic bencode reader interface
type bencodeReader interface {
	// Returns the next byte
	readByte() (byte, error)
	// Backtracks by 1 byte
	//
	// MUST only be called at most once after a call to nextByte()
	undoReadByte()
	// Reads a number until an end byte
	readIntTo(separator byte) (int, error)
	// Returns a string of a given length
	readString(length int) (string, error)
}

// A bencode decoder
type decoder struct {
	bencodeReader
}

// Reads a single integer from the decoder.
func (d *decoder) AsInt() (int, error) {
	if b, err := d.readByte(); err != nil {
		return 0, err
	} else if b != 'i' {
		return 0, ErrInvalidType
	}
	return d.readIntTo('e')
}

// Reads a single string from the decoder.
func (d *decoder) AsString() (string, error) {
	length, err := d.readIntTo(':')
	if err != nil {
		return "", err
	}
	if length < 0 {
		return "", ErrInvalidStringLen
	}
	return d.readString(length)
}

// Reads a single list from the decoder.
// Assumes that the first 'l' has been read.
func (d *decoder) asList() ([]interface{}, error) {
	list := []interface{}{}
	for {
		// Read type
		t, err := d.readByte()
		if err != nil {
			return list, err
		}
		if t == 'e' {
			break
		}
		// Read object
		obj, err := d.asInterface(t)
		if err != nil {
			return list, err
		}
		// Add to list
		list = append(list, obj)
	}
	return list, nil
}

// Reads a single list from the decoder.
func (d *decoder) AsList() ([]interface{}, error) {
	if b, err := d.readByte(); err != nil {
		return nil, err
	} else if b != 'l' {
		return nil, ErrInvalidType
	}
	return d.asList()
}

// Reads a single dictionary from the decoder.
// Assumes that the first 'd' has been read.
func (d *decoder) asDict() (map[string]interface{}, error) {
	dict := map[string]interface{}{}
	for {
		// Check if end
		if t, err := d.readByte(); err != nil {
			return dict, err
		} else if t == 'e' {
			break
		}
		d.undoReadByte()
		// Read key
		key, err := d.AsString()
		if err != nil {
			return dict, err
		}
		// Read value's type
		t, err := d.readByte()
		if err != nil {
			return dict, err
		}
		// Read value
		value, err := d.asInterface(t)
		if err != nil {
			return dict, err
		}
		// Add to map
		dict[key] = value
	}
	return dict, nil
}

// Reads a single dictionary from the decoder.
func (d *decoder) AsDict() (map[string]interface{}, error) {
	if b, err := d.readByte(); err != nil {
		return nil, err
	} else if b != 'd' {
		return nil, ErrInvalidType
	}
	return d.asDict()
}

// Reads from the decoder and returns an interface.
// Expects a type byte to be provided.
func (d *decoder) asInterface(t byte) (interface{}, error) {
	if t == 'i' {
		return d.readIntTo('e')
	} else if t == 'l' {
		return d.asList()
	} else if t == 'd' {
		return d.asDict()
	}
	d.undoReadByte()
	return d.AsString()
}

// Reads from the decoder and returns an interface.
func (d *decoder) AsInterface() (interface{}, error) {
	t, err := d.readByte()
	if err != nil {
		return nil, err
	}
	return d.asInterface(t)
}
