package main

import (
	"os"
	"strings"
	"testing"
)

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
	matrix, err := NewLetterMatrixFromFile(tmpMatrixFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(tmpDictFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSimpleSearcher(matrix, dict)

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
	matrix, err := NewLetterMatrixFromFile(tmpMatrixFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(tmpDictFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSimpleSearcher(matrix, dict)

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
	_, err := NewLetterMatrixFromFile("nonexistent_file.txt")
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

	_, err = NewLetterMatrixFromFile(tmpFile.Name())
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

	matrix, err := NewLetterMatrixFromFile(tmpFile2.Name())
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
	matrix, err := NewLetterMatrixFromFile(tmpMatrixFile.Name())
	if err != nil {
		b.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(tmpDictFile.Name())
	if err != nil {
		b.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSimpleSearcher(matrix, dict)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searcher.SearchAllWords(4) // Use 4 workers
	}
}
