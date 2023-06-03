package bencode_test

import (
	"bytes"
	"fmt"
	"math"
	"runtime"
	"strings"
	"testing"

	"github.com/stefanovazzocell/bencode"
)

// Testing helper
func GenericEncoderTester(t *testing.T, v interface{}, expected string) {
	t.Logf("Test case (%T): %v", v, v)
	// Fetch encoder
	encoder, err := bencode.NewEncoderFromInterface(v)
	if err != nil {
		t.Fatalf("NewEncoderFromInterface returned error: %v", err)
	}
	// Retrieve result
	buf := bytes.Buffer{}
	if n, err := encoder.WriteTo(&buf); err != nil {
		t.Fatalf("Failed to write to a buffer: %v", err)
	} else if n != int64(buf.Len()) {
		t.Fatalf("Wrote %d, but only has %d Len()", n, buf.Len())
	}
	actualReader := buf.Bytes()
	actualStr := encoder.String()
	actualBytes := encoder.Bytes()
	if !bytes.Equal([]byte(actualStr), actualBytes) {
		t.Fatalf("Return values don't match between bytes %q (%x) and string %q (%s)", actualBytes, actualBytes, actualStr, actualStr)
	}
	if !bytes.Equal(actualReader, actualBytes) {
		t.Fatalf("Return values don't match between the reader output %q (%x) and bytes %q (%s)", actualReader, actualReader, actualBytes, actualBytes)
	}
	// Compare result
	if expected != actualStr {
		t.Fatalf("Expected %q (%x) doesn't match actual %q (%x)", expected, expected, actualStr, actualStr)
	}
}

func TestInvalidEncoding(t *testing.T) {
	for _, invalid := range invalidTestCases {
		t.Logf("Testing with type %T: %v", invalid, invalid)
		if encoder, err := bencode.NewEncoderFromInterface(invalid); err == nil {
			t.Fatalf("Encoded invalid type without error as %v", encoder)
		}
	}
}

func TestStringEncoding(t *testing.T) {
	for _, str := range stringsTestCases {
		// Generate expected
		expected := fmt.Sprintf("%d:%s", len(str), str)
		// Run Test
		GenericEncoderTester(t, str, expected)
		// Check string-specific encoding
		encoder := bencode.NewEncoderFromString(str)
		actualStr := encoder.String()
		if actualStr != expected {
			t.Fatalf("Expected %q (%x) from string-specific encoder doesn't match actual %q (%x)", expected, expected, actualStr, actualStr)
		}
	}
}

func TestIntEncoding(t *testing.T) {
	for i, expected := range int64TestCases {
		t.Logf("Testing %d with expected value %q", i, expected)
		// Run Tests
		if strings.Contains(runtime.GOARCH, "64") {
			// Poor man's way to check if it's 64bit
			GenericEncoderTester(t, int(i), expected)
		}
		GenericEncoderTester(t, int64(i), expected)
		if math.MinInt32 <= i && i <= math.MaxInt32 {
			GenericEncoderTester(t, int32(i), expected)
		}
		if math.MinInt16 <= i && i <= math.MaxInt16 {
			GenericEncoderTester(t, int16(i), expected)
		}
		if math.MinInt8 <= i && i <= math.MaxInt8 {
			GenericEncoderTester(t, int8(i), expected)
		}
		// Check string-specific encoding
		encoder := bencode.NewEncoderFromInt(i)
		actualStr := encoder.String()
		if actualStr != expected {
			t.Fatalf("Expected %q (%x) from int-specific encoder doesn't match actual %q (%x)", expected, expected, actualStr, actualStr)
		}
	}
}

func TestUintEncoding(t *testing.T) {
	for i, expected := range uintsTestCases {
		t.Logf("Testing %d with expected value %q", i, expected)
		// Run Tests
		if math.MaxInt64 == (1<<63 - 1) {
			// Poor man's way to check if it's 64bit
			GenericEncoderTester(t, uint(i), expected)
		}
		GenericEncoderTester(t, uint64(i), expected)
		if i <= math.MaxUint32 {
			GenericEncoderTester(t, uint32(i), expected)
		}
		if i <= math.MaxUint16 {
			GenericEncoderTester(t, uint16(i), expected)
		}
		if i <= math.MaxUint8 {
			GenericEncoderTester(t, uint8(i), expected)
		}
		// Check string-specific encoding
		encoder := bencode.NewEncoderFromUint(i)
		actualStr := encoder.String()
		if actualStr != expected {
			t.Fatalf("Expected %q (%x) from uint-specific encoder doesn't match actual %q (%x)", expected, expected, actualStr, actualStr)
		}
	}
}

func TestSliceEncoding(t *testing.T) {
	for expected, list := range slicesTestCases {
		// Run Test
		GenericEncoderTester(t, list, expected)
		// Check string-specific encoding
		encoder, err := bencode.NewEncoderFromSlice(list)
		if err != nil {
			t.Fatalf("Got unexpected error from slice-specific encoder: %v", err)
		}
		actualStr := encoder.String()
		if actualStr != expected {
			t.Fatalf("Expected %q (%x) from string-specific encoder doesn't match actual %q (%x)", expected, expected, actualStr, actualStr)
		}
	}
}

func TestMapEncoding(t *testing.T) {
	for expected, m := range mapTestCases {
		// Run Test
		GenericEncoderTester(t, m, expected)
		// Check string-specific encoding
		encoder, err := bencode.NewEncoderFromMap(m)
		if err != nil {
			t.Fatalf("Got unexpected error from slice-specific encoder: %v", err)
		}
		actualStr := encoder.String()
		if actualStr != expected {
			t.Fatalf("Expected %q (%x) from string-specific encoder doesn't match actual %q (%x)", expected, expected, actualStr, actualStr)
		}
	}
}

func BenchmarkEncoder(b *testing.B) {
	for benchName, testInterface := range encoderBenchmarks {
		b.Run(benchName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bencode.NewEncoderFromInterface(testInterface)
			}
		})
	}
}
