package bencode_test

import (
	"reflect"
	"testing"

	"github.com/stefanovazzocell/bencode"
)

// Helper to automatically perform all testing
func ParserTestingHelper(t *testing.T, testCase string, expected interface{}) {
	t.Logf("Test case: %q", testCase)
	// Attempt parsing
	actual, err := bencode.NewParserFromString(testCase).AsInterface()
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
		newActual, err := bencode.NewParserFromString(testCase).AsString()
		if err != nil {
			t.Fatalf("Failed to parse as string, got error: %v", err)
		}
		if newActual != actual {
			t.Fatalf("Specific parsing as string returned %q, but originally was parsed as %q", newActual, actual)
		}
	case int:
		newActual, err := bencode.NewParserFromString(testCase).AsInt()
		if err != nil {
			t.Fatalf("Failed to parse as int, got error: %v", err)
		}
		if newActual != actual {
			t.Fatalf("Specific parsing as int returned %d, but originally was parsed as %d", newActual, actual)
		}
	case []interface{}:
		newActual, err := bencode.NewParserFromString(testCase).AsList()
		if err != nil {
			t.Fatalf("Failed to parse as []interface{}, got error: %v", err)
		}
		if !reflect.DeepEqual(actual, newActual) {
			t.Fatalf("Specific parsing as []interface{} returned %v, but originally was parsed as %v", newActual, actual)
		}
	case map[string]interface{}:
		newActual, err := bencode.NewParserFromString(testCase).AsDict()
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
