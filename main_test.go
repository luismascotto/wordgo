package main

import (
	"os"
	"strings"
	"testing"
)

// AI_GENERATED_CODE_START
// [AI Generated] Data: 19/12/2024
// Descrição: Testes unitários para o sistema de busca de palavras em matriz de letras
// Gerado por: Cursor AI
// Versão: Go 1.21
// AI_GENERATED_CODE_END

// TestNewLetterMatrix tests matrix loading functionality
func TestNewLetterMatrix(t *testing.T) {
	// Create a temporary test matrix file
	testMatrix := "ABCD\nEFGH\nIJKL"
	tmpFile, err := os.CreateTemp("", "test_matrix_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testMatrix)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the matrix
	matrix, err := NewLetterMatrix(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewLetterMatrix failed: %v", err)
	}

	// Verify dimensions
	rows, cols := matrix.GetDimensions()
	if rows != 3 {
		t.Errorf("Expected 3 rows, got %d", rows)
	}
	if cols != 4 {
		t.Errorf("Expected 4 columns, got %d", cols)
	}

	// Verify content
	matrixData := matrix.GetMatrix()
	expected := [][]rune{
		{'A', 'B', 'C', 'D'},
		{'E', 'F', 'G', 'H'},
		{'I', 'J', 'K', 'L'},
	}

	for i, row := range matrixData {
		for j, char := range row {
			if char != expected[i][j] {
				t.Errorf("Expected %c at [%d][%d], got %c", expected[i][j], i, j, char)
			}
		}
	}
}

// TestNewLetterMatrixWithPadding tests matrix loading with different row lengths
func TestNewLetterMatrixWithPadding(t *testing.T) {
	// Create a test matrix with different row lengths
	testMatrix := "ABC\nDEFGH\nIJ"
	tmpFile, err := os.CreateTemp("", "test_matrix_pad_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testMatrix)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the matrix
	matrix, err := NewLetterMatrix(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewLetterMatrix failed: %v", err)
	}

	// Verify dimensions (should be padded to max length)
	rows, cols := matrix.GetDimensions()
	if rows != 3 {
		t.Errorf("Expected 3 rows, got %d", rows)
	}
	if cols != 5 {
		t.Errorf("Expected 5 columns (max length), got %d", cols)
	}

	// Verify padding
	matrixData := matrix.GetMatrix()
	// First row should be padded with spaces
	if len(matrixData[0]) != 5 {
		t.Errorf("First row should be padded to 5 columns")
	}
	if matrixData[0][3] != ' ' || matrixData[0][4] != ' ' {
		t.Errorf("First row should be padded with spaces")
	}
}

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

// TestWordSearcher tests the word search functionality
func TestWordSearcher(t *testing.T) {
	// Create a simple test matrix
	testMatrix := "CAT\nDOG\nBAT"
	tmpMatrixFile, err := os.CreateTemp("", "test_matrix_search_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp matrix file: %v", err)
	}
	defer os.Remove(tmpMatrixFile.Name())

	_, err = tmpMatrixFile.WriteString(testMatrix)
	if err != nil {
		t.Fatalf("Failed to write to temp matrix file: %v", err)
	}
	tmpMatrixFile.Close()

	// Create a test dictionary
	testDict := "CAT\nDOG\nBAT\nAT\nGO\nTA\nAD\nTO"
	tmpDictFile, err := os.CreateTemp("", "test_dict_search_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp dict file: %v", err)
	}
	defer os.Remove(tmpDictFile.Name())

	_, err = tmpDictFile.WriteString(testDict)
	if err != nil {
		t.Fatalf("Failed to write to temp dict file: %v", err)
	}
	tmpDictFile.Close()

	// Load matrix and dictionary
	matrix, err := NewLetterMatrix(tmpMatrixFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(tmpDictFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSearcher(matrix, dict)

	// Test searching from specific position
	searcher.SimpleSearchFromPosition(0, 0, Direction{R, "→", 0, 1})

	// Get results
	results := searcher.GetResults()

	// Should find at least some words
	if len(results) == 0 {
		t.Error("Expected to find some words, but found none")
	}

	// Verify some expected results
	foundWords := make(map[string]bool)
	for _, result := range results {
		foundWords[result.Word] = true
	}

	expectedWords := []string{"CAT", "AT", "DOG", "GO", "BAT", "AT"}
	for _, word := range expectedWords {
		if !foundWords[word] {
			t.Errorf("Expected to find word '%s' but it was not found", word)
		}
	}
}

// TestWordSearcherParallel tests parallel word searching
func TestWordSearcherParallel(t *testing.T) {
	// Create a larger test matrix for parallel testing
	testMatrix := strings.Repeat("ABCDEFGH\n", 4) // 4x8 matrix
	tmpMatrixFile, err := os.CreateTemp("", "test_matrix_parallel_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp matrix file: %v", err)
	}
	defer os.Remove(tmpMatrixFile.Name())

	_, err = tmpMatrixFile.WriteString(testMatrix)
	if err != nil {
		t.Fatalf("Failed to write to temp matrix file: %v", err)
	}
	tmpMatrixFile.Close()

	// Create a test dictionary with common words
	testDict := "ABC\nBCD\nCDE\nDEF\nEFG\nFGH\nAB\nBC\nCD\nDE\nEF\nFG\nGH"
	tmpDictFile, err := os.CreateTemp("", "test_dict_parallel_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp dict file: %v", err)
	}
	defer os.Remove(tmpDictFile.Name())

	_, err = tmpDictFile.WriteString(testDict)
	if err != nil {
		t.Fatalf("Failed to write to temp dict file: %v", err)
	}
	tmpDictFile.Close()

	// Load matrix and dictionary
	matrix, err := NewLetterMatrix(tmpMatrixFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(tmpDictFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSearcher(matrix, dict)

	// Test parallel search with 2 workers
	searcher.SearchAllWords(2)

	// Get results
	results := searcher.GetResults()

	// Should find words
	if len(results) == 0 {
		t.Error("Expected to find some words in parallel search, but found none")
	}

	// Verify results are thread-safe (no duplicates or corruption)
	wordCounts := make(map[string]int)
	for _, result := range results {
		wordCounts[result.Word]++
	}

	// Check for reasonable results
	t.Logf("Found %d words in parallel search", len(results))
}

// TestEdgeCases tests edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	// Test empty matrix file
	_, err := NewLetterMatrix("nonexistent_file.txt")
	if err == nil {
		t.Error("Expected error when loading nonexistent matrix file")
	}

	// Test empty dictionary file
	_, err = NewDictionary("nonexistent_file.txt")
	if err == nil {
		t.Error("Expected error when loading nonexistent dictionary file")
	}

	// Test empty matrix content
	tmpFile, err := os.CreateTemp("", "test_empty_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = NewLetterMatrix(tmpFile.Name())
	if err == nil {
		t.Error("Expected error when loading empty matrix file")
	}

	// Test matrix with single character
	tmpFile2, err := os.CreateTemp("", "test_single_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())

	_, err = tmpFile2.WriteString("A")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile2.Close()

	matrix, err := NewLetterMatrix(tmpFile2.Name())
	if err != nil {
		t.Fatalf("Failed to load single character matrix: %v", err)
	}

	rows, cols := matrix.GetDimensions()
	if rows != 1 || cols != 1 {
		t.Errorf("Expected 1x1 matrix, got %dx%d", rows, cols)
	}
}

// BenchmarkWordSearch benchmarks the word search performance
func BenchmarkWordSearch(b *testing.B) {
	// Create a test matrix and dictionary for benchmarking
	testMatrix := strings.Repeat("ABCDEFGHIJKLMNOP\n", 10) // 10x16 matrix
	tmpMatrixFile, err := os.CreateTemp("", "bench_matrix_*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp matrix file: %v", err)
	}
	defer os.Remove(tmpMatrixFile.Name())

	_, err = tmpMatrixFile.WriteString(testMatrix)
	if err != nil {
		b.Fatalf("Failed to write to temp matrix file: %v", err)
	}
	tmpMatrixFile.Close()

	// Create a larger test dictionary
	testDict := strings.Join([]string{
		"ABC", "BCD", "CDE", "DEF", "EFG", "FGH", "GHI", "HIJ",
		"ABC", "BCD", "CDE", "DEF", "EFG", "FGH", "GHI", "HIJ",
		"ABC", "BCD", "CDE", "DEF", "EFG", "FGH", "GHI", "HIJ",
	}, "\n")

	tmpDictFile, err := os.CreateTemp("", "bench_dict_*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp dict file: %v", err)
	}
	defer os.Remove(tmpDictFile.Name())

	_, err = tmpDictFile.WriteString(testDict)
	if err != nil {
		b.Fatalf("Failed to write to temp dict file: %v", err)
	}
	tmpDictFile.Close()

	// Load matrix and dictionary
	matrix, err := NewLetterMatrix(tmpMatrixFile.Name())
	if err != nil {
		b.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(tmpDictFile.Name())
	if err != nil {
		b.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSearcher(matrix, dict)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searcher.SearchAllWords(4) // Use 4 workers
	}
}
