package bencode_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stefanovazzocell/bencode"
)

func TestStringParser(t *testing.T) {
	// Test strings
	for _, actual := range stringsTestCases {
		ParserTestingHelper(t, fmt.Sprintf("%d:%s", len(actual), actual), actual)
	}
	// Test int
	for test, actual := range intsTestCases {
		ParserTestingHelper(t, test, actual)
	}
	// Test slices
	for test, actual := range slicesTestCases {
		ParserTestingHelper(t, test, actual)
	}
	// Test maps
	for test, actual := range complexMapTestCases {
		ParserTestingHelper(t, test, actual)
	}
	// Test Invalid Parse
	for _, invalid := range invalidParserInputs {
		if out, err := bencode.NewParserFromString(invalid).AsInterface(); err == nil {
			t.Fatalf("Expected invalid %q (AsInterface) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromString(invalid).AsInt(); err == nil {
			t.Fatalf("Expected invalid %q (AsInt) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromString(invalid).AsString(); err == nil {
			t.Fatalf("Expected invalid %q (AsString) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromString(invalid).AsList(); err == nil {
			t.Fatalf("Expected invalid %q (AsList) to fail.\nInstead got %v", invalid, out)
		}
		if out, err := bencode.NewParserFromString(invalid).AsDict(); err == nil {
			t.Fatalf("Expected invalid %q (AsDict) to fail.\nInstead got %v", invalid, out)
		}
	}
	// Test Invalid Type
	for invalid, trueType := range invalidTypeParse {
		if trueType != 'i' {
			if out, err := bencode.NewParserFromString(invalid).AsInt(); err == nil {
				t.Fatalf("Expected invalid %q (AsInt) to fail.\nInstead got %v", invalid, out)
			}
		}
		if trueType != 's' {
			if out, err := bencode.NewParserFromString(invalid).AsString(); err == nil {
				t.Fatalf("Expected invalid %q (AsString) to fail.\nInstead got %v", invalid, out)
			}
		}
		if trueType != 'l' {
			if out, err := bencode.NewParserFromString(invalid).AsList(); err == nil {
				t.Fatalf("Expected invalid %q (AsList) to fail.\nInstead got %v", invalid, out)
			}
		}
		if trueType != 'd' {
			if out, err := bencode.NewParserFromString(invalid).AsDict(); err == nil {
				t.Fatalf("Expected invalid %q (AsDict) to fail.\nInstead got %v", invalid, out)
			}
		}
	}
}

func BenchmarkStringParser(b *testing.B) {
	for benchName, testString := range parserBenchmarks {
		b.Run(benchName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bencode.NewParserFromString(testString).AsInterface()
			}
		})
	}
}

func FuzzStringParser(f *testing.F) {
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
		bencode.NewParserFromString(test).AsInt()
		bencode.NewParserFromString(test).AsString()
		bencode.NewParserFromString(test).AsList()
		bencode.NewParserFromString(test).AsDict()
		// Now parse as generic interface
		parsedObj, err := bencode.NewParserFromString(test).AsInterface()
		if err != nil {
			t.SkipNow()
		}
		// Try to re-encode
		reEncoded, err := bencode.NewEncoderFromInterface(parsedObj)
		if err != nil {
			t.Fatalf("Failed to re-encode %v: %v", parsedObj, err)
		}
		// Parse encoded for final check
		reParsedObj, err := bencode.NewParserFromString(reEncoded.String()).AsInterface()
		if err != nil {
			t.Fatalf("Failed to re-decode %q: %v", reEncoded.String(), err)
		}
		// Compare original
		if !reflect.DeepEqual(parsedObj, reParsedObj) {
			t.Fatalf("Originally parsed as %v, encoded as %q, and re-parsed as different %v", parsedObj, reEncoded, reParsedObj)
		}
	})
}
