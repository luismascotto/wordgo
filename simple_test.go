package main

import (
	"os"
	"strings"
	"testing"
)

// AI_GENERATED_CODE_START
// [AI Generated] Data: 19/12/2024
// Descrição: Testes simples para verificar funcionalidade básica
// Gerado por: Cursor AI
// Versão: Go 1.21
// AI_GENERATED_CODE_END

// TestBasicWordSearch tests basic word search functionality
func TestBasicWordSearch(t *testing.T) {
	// Create a simple test matrix
	testMatrix := "CAT\nDOG"

	// Create a test dictionary
	testDict := "CAT\nDOG\nAT\nGO"

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_basic_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_basic_*.txt", testDict)
	defer dictFile.Close()
	defer os.Remove(dictFile.Name())

	// Load matrix and dictionary
	matrix, err := NewLetterMatrixFromFile(matrixFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(dictFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSimpleSearcher(matrix, dict)

	// Search from all positions in all directions
	searcher.SearchAllWords(1) // Use single worker for deterministic results

	// Get results
	results := searcher.GetResults()

	// Should find some words
	if len(results) == 0 {
		t.Error("Expected to find some words, but found none")
	}

	// Log results
	t.Logf("Found %d words in basic test", len(results))
	for _, result := range results {
		t.Logf("'%s' at (%d,%d) in direction %s",
			result.Word, result.StartRow, result.StartCol, result.Direction)
	}

	// Verify we found at least some expected words
	foundWords := make(map[string]bool)
	for _, result := range results {
		foundWords[result.Word] = true
	}

	// Check for basic words
	if !foundWords["CAT"] {
		t.Error("Expected to find 'CAT' but it was not found")
	}
	if !foundWords["DOG"] {
		t.Error("Expected to find 'DOG' but it was not found")
	}
}

// TestWordSearchWithSpacesSimple tests word search with spaces using a simpler matrix
func TestWordSearchWithSpacesSimple(t *testing.T) {
	// Create a test matrix with intentional spaces
	testMatrix := "A B\nCD"

	// Create a test dictionary
	testDict := "AB\nCD\nA\nB\nC\nD"

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_spaces_simple_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_spaces_simple_*.txt", testDict)
	defer dictFile.Close()
	defer os.Remove(dictFile.Name())

	// Load matrix and dictionary
	matrix, err := NewLetterMatrixFromFile(matrixFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test matrix: %v", err)
	}

	dict, err := NewDictionary(dictFile.Name())
	if err != nil {
		t.Fatalf("Failed to load test dictionary: %v", err)
	}

	// Create word searcher
	searcher := NewWordSimpleSearcher(matrix, dict)

	// Search from all positions
	searcher.SearchAllWords(1)

	// Get results
	results := searcher.GetResults()

	// Should find words
	if len(results) == 0 {
		t.Error("Expected to find some words, but found none")
	}

	// Log results
	t.Logf("Found %d words in matrix with spaces", len(results))
	for _, result := range results {
		t.Logf("'%s' at (%d,%d) in direction %s",
			result.Word, result.StartRow, result.StartCol, result.Direction)
	}

	// Verify that spaces are handled correctly (search should stop at spaces)
	for _, result := range results {
		if strings.Contains(result.Word, " ") {
			t.Errorf("Found word with space: '%s'", result.Word)
		}
	}
}
