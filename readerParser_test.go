package bencode_test

import (
	"crypto/rand"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stefanovazzocell/bencode"
)

// A reader that never returns anything
type NoProgressReader struct{}

func (r *NoProgressReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// A reader that always returns a negative number
type NegativeReader struct{}

func (r *NegativeReader) Read(p []byte) (n int, err error) {
	return -1, nil
}

// A reader that might or might not misbehave
type ChaosReader struct {
	sr strings.Reader
}

func (r *ChaosReader) Read(p []byte) (n int, err error) {
	b := make([]byte, 1)
	n, err = rand.Reader.Read(b)
	if n != 1 || err != nil {
		panic(fmt.Sprintf("ChaosReader: got (%d, %v) from rand.Reader.Read", n, err))
	}
	if b[0] <= 60 {
		return 0, nil
	}
	if b[0] <= 120 && len(p) > 10 {
		return r.sr.Read(p[:10])
	}
	return r.sr.Read(p)
}

func newChaosReader(str string) *ChaosReader {
	return &ChaosReader{
		sr: *strings.NewReader(str),
	}
}

// Helper to automatically perform all testing
func ReaderParserTestHelper(t *testing.T, testCase string, expected interface{}) {
	t.Logf("Test case: %q", testCase)
	// Attempt parsing
	actual, err := bencode.NewParserFromReader(strings.NewReader(testCase)).AsInterface()
	if err != nil {
		t.Fatalf("Failed to parse as interface, got error: %v", err)
	}
	// Compare results
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, but got %v", expected, actual)
	}
	// Attempt specific parsing
	t.Logf("Attempting specific parsing as %T", actual)
	switch actual := actual.(type) {
	case string:
		newActual, err := bencode.NewParserFromReader(strings.NewReader(testCase)).AsString()
		if err != nil {
			t.Fatalf("Failed to parse as string, got error: %v", err)
		}
		if newActual != actual {
			t.Fatalf("Specific parsing as string returned %q, but originally was parsed as %q", newActual, actual)
		}
	case int:
		newActual, err := bencode.NewParserFromReader(strings.NewReader(testCase)).AsInt()
		if err != nil {
			t.Fatalf("Failed to parse as int, got error: %v", err)
		}
		if newActual != actual {
			t.Fatalf("Specific parsing as int returned %d, but originally was parsed as %d", newActual, actual)
		}
	case []interface{}:
		newActual, err := bencode.NewParserFromReader(strings.NewReader(testCase)).AsList()
		if err != nil {
			t.Fatalf("Failed to parse as []interface{}, got error: %v", err)
		}
		if !reflect.DeepEqual(actual, newActual) {
			t.Fatalf("Specific parsing as []interface{} returned %v, but originally was parsed as %v", newActual, actual)
		}
	case map[string]interface{}:
		newActual, err := bencode.NewParserFromReader(strings.NewReader(testCase)).AsDict()
		if err != nil {
			t.Fatalf("Failed to parse as map[string]interface{}, got error: %v", err)
		}
		if !reflect.DeepEqual(actual, newActual) {
			t.Fatalf("Specific parsing as map[string]interface{} returned %v, but originally was parsed as %v", newActual, actual)
		}
	default:
		t.Fatalf("Got an unexpected type back: %T", actual)
	}
}

func TestReaderParser(t *testing.T) {
	// Test strings
	for _, actual := range stringsTestCases {
		ReaderParserTestHelper(t, fmt.Sprintf("%d:%s", len(actual), actual), actual)
	}
	// Test int
	for test, actual := range intsTestCases {
		ReaderParserTestHelper(t, test, actual)
	}
	// Test slices
	for test, actual := range slicesTestCases {
		ReaderParserTestHelper(t, test, actual)
	}
	// Test maps
	for test, actual := range complexMapTestCases {
		ReaderParserTestHelper(t, test, actual)
	}
	// Test Invalid Parse
	for _, invalid := range invalidParserInputs {
		if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsInterface(); err == nil {
			t.Fatalf("Expected invalid %q (AsInterface) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsInt(); err == nil {
			t.Fatalf("Expected invalid %q (AsInt) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsString(); err == nil {
			t.Fatalf("Expected invalid %q (AsString) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsList(); err == nil {
			t.Fatalf("Expected invalid %q (AsList) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsDict(); err == nil {
			t.Fatalf("Expected invalid %q (AsDict) to fail.\nInstead got %v", invalid, out)
		}
	}
	// Test Invalid Type
	for invalid, trueType := range invalidTypeParse {
		if trueType != 'i' {
			if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsInt(); err == nil {
				t.Fatalf("Expected invalid %q (AsInt) to fail.\nInstead got %v", invalid, out)
			}
		}
		if trueType != 's' {
			if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsString(); err == nil {
				t.Fatalf("Expected invalid %q (AsString) to fail.\nInstead got %v", invalid, out)
			}
		}
		if trueType != 'l' {
			if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsList(); err == nil {
				t.Fatalf("Expected invalid %q (AsList) to fail.\nInstead got %v", invalid, out)
			}
		}
		if trueType != 'd' {
			if out, err := bencode.NewParserFromReader(strings.NewReader(invalid)).AsDict(); err == nil {
				t.Fatalf("Expected invalid %q (AsDict) to fail.\nInstead got %v", invalid, out)
			}
		}
	}
	// Test NoProgressReader
	if obj, err := bencode.NewParserFromReader(&NoProgressReader{}).AsInterface(); err != io.ErrNoProgress {
		t.Fatalf("Got (%v, %v) from NewParserFromReader(&NoProgressReader{}).AsInterface()", obj, err)
	}
	if obj, err := bencode.NewParserFromReader(&NoProgressReader{}).AsInt(); err != io.ErrNoProgress {
		t.Fatalf("Got (%v, %v) from NewParserFromReader(&NoProgressReader{}).AsInt()", obj, err)
	}
	if obj, err := bencode.NewParserFromReader(&NoProgressReader{}).AsList(); err != io.ErrNoProgress {
		t.Fatalf("Got (%v, %v) from NewParserFromReader(&NoProgressReader{}).AsList()", obj, err)
	}
	if obj, err := bencode.NewParserFromReader(&NoProgressReader{}).AsDict(); err != io.ErrNoProgress {
		t.Fatalf("Got (%v, %v) from NewParserFromReader(&NoProgressReader{}).AsDict()", obj, err)
	}
	if obj, err := bencode.NewParserFromReader(&NoProgressReader{}).AsString(); err != io.ErrNoProgress {
		t.Fatalf("Got (%v, %v) from NewParserFromReader(&NoProgressReader{}).AsString()", obj, err)
	}
}

func TestReaderParserPanic(t *testing.T) {
	// Test NegativeReader
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("NewParserFromReader(&NegativeReader{}).AsInterface() was expected to panic, but did not")
		}
	}()
	bencode.NewParserFromReader(&NegativeReader{}).AsInterface()
	t.Fatalf("NewParserFromReader(&NegativeReader{}).AsInterface() was expected to panic, but did not")
}

func TestReaderParserChaos(t *testing.T) {
	tester := func(testCase string, expected interface{}) {
		obj, err := bencode.NewParserFromReader(newChaosReader(testCase)).AsInterface()
		if err == nil && !reflect.DeepEqual(obj, expected) {
			t.Fatalf("Got (%v, %v) from newChaosReader(%q) in AsInterface()", obj, err, testCase)
		}
		// Try with specific types
		switch expected := expected.(type) {
		case string:
			obj, err := bencode.NewParserFromReader(newChaosReader(testCase)).AsString()
			if err == nil && obj != expected {
				t.Fatalf("Got (%s, %v) from newChaosReader(%q) in AsString()", obj, err, testCase)
			}
		case int:
			obj, err := bencode.NewParserFromReader(newChaosReader(testCase)).AsInt()
			if err == nil && obj != expected {
				t.Fatalf("Got (%d, %v) from newChaosReader(%q) in AsInt()", obj, err, testCase)
			}
		case []interface{}:
			obj, err := bencode.NewParserFromReader(newChaosReader(testCase)).AsList()
			if err == nil && !reflect.DeepEqual(obj, expected) {
				t.Fatalf("Got (%v, %v) from newChaosReader(%q) in AsList()", obj, err, testCase)
			}
		case map[string]interface{}:
			obj, err := bencode.NewParserFromReader(newChaosReader(testCase)).AsDict()
			if err == nil && !reflect.DeepEqual(obj, expected) {
				t.Fatalf("Got (%v, %v) from newChaosReader(%q) in AsList()", obj, err, testCase)
			}
		}
	}
	// Run cases (multiple times to hit different breakage points)
	for i := 0; i < 1024; i++ {
		for _, actual := range stringsTestCases {
			tester(fmt.Sprintf("%d:%s", len(actual), actual), actual)
		}
		for test, expected := range intsTestCases {
			tester(test, expected)
		}
		for test, expected := range slicesTestCases {
			tester(test, expected)
		}
		for test, expected := range complexMapTestCases {
			tester(test, expected)
		}
	}
}

func BenchmarkReaderParser(b *testing.B) {
	for benchName, testString := range parserBenchmarks {
		b.Run(benchName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bencode.NewParserFromReader(strings.NewReader(testString)).AsInterface()
			}
		})
	}
	b.Run("torrentStringAsString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bencode.NewParserFromReader(strings.NewReader(fedoraMagnetParsed)).AsString()
		}
	})
}

func FuzzReaderParser(f *testing.F) {
	// > Add examples
	// Add strings examples
	for _, actual := range stringsTestCases {
		f.Add(fmt.Sprintf("%d:%s", len(actual), actual))
	}
	// Add int examples
	for test := range intsTestCases {
		f.Add(test)
	}
	// Add slices examples
	for test := range slicesTestCases {
		f.Add(test)
	}
	// Add maps examples
	for test := range complexMapTestCases {
		f.Add(test)
	}
	// Add Invalid Parse examples
	for _, invalid := range invalidParserInputs {
		f.Add(invalid)
	}
	// Add Invalid Type examples
	for invalid := range invalidTypeParse {
		f.Add(invalid)
	}

	// > Setup test
	f.Fuzz(func(t *testing.T, test string) {
		// Look for panics in type-specific parsing
		bencode.NewParserFromReader(strings.NewReader(test)).AsInt()
		bencode.NewParserFromReader(strings.NewReader(test)).AsString()
		bencode.NewParserFromReader(strings.NewReader(test)).AsList()
		bencode.NewParserFromReader(strings.NewReader(test)).AsDict()
		// Now parse as generic interface
		parsedObj, err := bencode.NewParserFromReader(strings.NewReader(test)).AsInterface()
		if err != nil {
			t.SkipNow()
		}
		// Try to re-encode
		reEncoded, err := bencode.NewEncoderFromInterface(parsedObj)
		if err != nil {
			t.Fatalf("Failed to re-encode %v: %v", parsedObj, err)
		}
		// Parse encoded for final check
		reParsedObj, err := bencode.NewParserFromReader(strings.NewReader(reEncoded.String())).AsInterface()
		if err != nil {
			t.Fatalf("Failed to re-decode %q: %v", reEncoded.String(), err)
		}
		// Compare original
		if !reflect.DeepEqual(parsedObj, reParsedObj) {
			t.Fatalf("Originally parsed as %v, encoded as %q, and re-parsed as different %v", parsedObj, reEncoded, reParsedObj)
		}
	})
}
