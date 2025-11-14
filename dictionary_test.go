package main

import (
	"os"
	"testing"
)

// TestNewDictionary tests dictionary loading functionality
func TestNewDictionary(t *testing.T) {
	// Create a temporary test dictionary file
	testDict := "CAT\nDOG\nBIRD\nELEPHANT\nZEBRA"
	tmpFile, err := os.CreateTemp("", "test_dict_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testDict)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the dictionary
	dict, err := NewDictionary(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewDictionary failed: %v", err)
	}

	// Test word existence
	testCases := []struct {
		word     string
		expected bool
	}{
		{"CAT", true},
		{"DOG", true},
		{"BIRD", true},
		{"ELEPHANT", true},
		{"ZEBRA", true},
		{"MOUSE", false},
		{"cat", true}, // Should be converted to uppercase
		{"DOG", true},
	}

	for _, tc := range testCases {
		result := dict.Contains(tc.word)
		if result != tc.expected {
			t.Errorf("Contains('%s') expected %v, got %v", tc.word, tc.expected, result)
		}
	}
}

// TestDictionaryPrefixAndWord tests prefix and word checking functionality
func TestDictionaryPrefixAndWord(t *testing.T) {
	// Create a test dictionary with words that have common prefixes
	testDict := "GRE\nGREEN\nGREET\nGREETING\nHELLO\nHELP"
	tmpFile, err := os.CreateTemp("", "test_dict_prefix_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testDict)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the dictionary
	dict, err := NewDictionary(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewDictionary failed: %v", err)
	}

	// Test prefix checking
	prefixTests := []struct {
		sequence string
		expected bool
	}{
		{"G", true},        // Prefix of GRE, GREEN, etc.
		{"GR", true},       // Prefix of GRE, GREEN, etc.
		{"GRE", true},      // Prefix of GREEN, GREET, etc.
		{"GREEN", true},    // Prefix of GREEN
		{"GREET", true},    // Prefix of GREETING
		{"GREETING", true}, // Word itself (no more children)
		{"H", true},        // Prefix of HELLO, HELP
		{"HE", true},       // Prefix of HELLO, HELP
		{"HEL", true},      // Prefix of HELLO, HELP
		{"HELLO", true},    // Word itself
		{"X", false},       // Not a prefix
		{"GREX", false},    // Not a prefix
	}

	for _, tc := range prefixTests {
		result := dict.IsPrefix(tc.sequence)
		if result != tc.expected {
			t.Errorf("IsPrefix('%s') expected %v, got %v", tc.sequence, tc.expected, result)
		}
	}

	// Test word checking
	wordTests := []struct {
		sequence string
		expected bool
	}{
		{"GRE", true},
		{"GREEN", true},
		{"GREET", true},
		{"GREETING", true},
		{"HELLO", true},
		{"HELP", true},
		{"G", false},   // Not a complete word
		{"GR", false},  // Not a complete word
		{"HE", false},  // Not a complete word
		{"HEL", false}, // Not a complete word
		{"XYZ", false}, // Not in dictionary
	}

	for _, tc := range wordTests {
		result := dict.IsWord(tc.sequence)
		if result != tc.expected {
			t.Errorf("IsWord('%s') expected %v, got %v", tc.sequence, tc.expected, result)
		}
	}
}
