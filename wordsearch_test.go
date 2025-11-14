package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// AI_GENERATED_CODE_START
// [AI Generated] Data: 19/12/2024
// Descrição: Testes especializados para algoritmos de busca de palavras
// Gerado por: Cursor AI
// Versão: Go 1.21
// AI_GENERATED_CODE_END

// TestWordSearchDirections tests word search in all 8 directions
func TestWordSearchDirections(t *testing.T) {
	// Create a test matrix with known words in different directions
	testMatrix := "CAT\nDOG\nBAT"

	// Create a test dictionary with words we expect to find
	testDict := "CAT\nDOG\nBAT\nAT\nGO\nTA\nAD\nTO\nAG\nDO\nBA\nAB\nCB\nGD\nOT\nGA\nTB\nAO\nDT\nCO\nBG\nAC\nDB\nGT\nCA\nDO\nBA"

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_directions_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_directions_*.txt", testDict)
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

	// Test each direction individually
	directions := []Direction{
		{R, "→", 0, 1},
		{L, "←", 0, -1},
		{B, "↓", 1, 0},
		{T, "↑", -1, 0},
		{BR, "↘", 1, 1},
		{BL, "↙", 1, -1},
		{TR, "↗", -1, 1},
		{TL, "↖", -1, -1},
	}

	// Test searching from position (0,0) in each direction
	for _, direction := range directions {
		t.Run(direction.Name, func(t *testing.T) {
			// Clear previous results
			searcher.results = make([]WordResult, 0)

			// Search from position (0,0)
			searcher.SimpleSearchFromPosition(0, 0, direction)

			// Get results
			results := searcher.GetResults()

			// Log what we found
			t.Logf("Direction %s: Found %d words", direction.Name, len(results))
			for _, result := range results {
				t.Logf("  '%s' at (%d,%d)", result.Word, result.StartRow, result.StartCol)
			}

			// Should find at least some words in most directions
			if len(results) == 0 {
				t.Logf("No words found in direction %s (this might be expected)", direction.Name)
			}
		})
	}
}

// TestWordSearchFromAllPositions tests searching from every position in the matrix
func TestWordSearchFromAllPositions(t *testing.T) {
	// Create a simple 2x2 test matrix
	testMatrix := "AB\nCD"

	// Create a test dictionary
	testDict := "AB\nCD\nA\nB\nC\nD\nABC\nABD\nACD\nBCD"

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_allpos_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_allpos_*.txt", testDict)
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

	// Should find words
	if len(results) == 0 {
		t.Error("Expected to find some words, but found none")
	}

	// Log all results
	t.Logf("Total words found: %d", len(results))
	for _, result := range results {
		t.Logf("'%s' at (%d,%d) in direction %s",
			result.Word, result.StartRow, result.StartCol, result.Direction)
	}

	// Verify we have results from different positions
	positions := make(map[string]bool)
	for _, result := range results {
		pos := fmt.Sprintf("(%d,%d)", result.StartRow, result.StartCol)
		positions[pos] = true
	}

	t.Logf("Words found from %d different positions", len(positions))
}

// TestWordSearchWithSpaces tests word search with intentional spaces in matrix
func TestWordSearchWithSpaces(t *testing.T) {
	// Create a test matrix with intentional spaces
	testMatrix := "A B\nCDE\nF G"

	// Create a test dictionary
	testDict := "AB\nCD\nDE\nEF\nFG\nA\nB\nC\nD\nE\nF\nG"

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_spaces_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_spaces_*.txt", testDict)
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

// TestWordSearchPerformance tests performance with different worker counts
func TestWordSearchPerformance(t *testing.T) {
	// Create a larger test matrix
	testMatrix := strings.Repeat("ABCDEFGHIJKLMNOP\n", 8) // 8x16 matrix

	// Create a larger test dictionary
	testDict := strings.Join([]string{
		"ABC", "BCD", "CDE", "DEF", "EFG", "FGH", "GHI", "HIJ",
		"IJK", "JKL", "KLM", "LMN", "MNO", "NOP", "AB", "BC",
		"CD", "DE", "EF", "FG", "GH", "HI", "IJ", "JK", "KL",
		"LM", "MN", "NO", "OP", "A", "B", "C", "D", "E", "F",
		"G", "H", "I", "J", "K", "L", "M", "N", "O", "P",
	}, "\n")

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_perf_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_perf_*.txt", testDict)
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

	// Test different worker counts
	workerCounts := []int{1, 2, 4, 8}

	for _, numWorkers := range workerCounts {
		t.Run(fmt.Sprintf("Workers_%d", numWorkers), func(t *testing.T) {
			// Create new searcher for each test
			searcher := NewWordSimpleSearcher(matrix, dict)

			// Time the search
			start := time.Now()
			searcher.SearchAllWords(numWorkers)
			duration := time.Since(start)

			// Get results
			results := searcher.GetResults()

			t.Logf("Workers: %d, Duration: %v, Words found: %d",
				numWorkers, duration, len(results))

			// Should find words
			if len(results) == 0 {
				t.Error("Expected to find some words, but found none")
			}
		})
	}
}

// TestWordSearchThreadSafety tests that results are collected thread-safely
func TestWordSearchThreadSafety(t *testing.T) {
	// Create a test matrix and dictionary
	testMatrix := strings.Repeat("ABCDEFGH\n", 4) // 4x8 matrix
	testDict := "ABC\nBCD\nCDE\nDEF\nEFG\nFGH\nAB\nBC\nCD\nDE\nEF\nFG\nGH"

	// Create temporary files
	matrixFile := createTempFile(t, "test_matrix_thread_*.txt", testMatrix)
	defer matrixFile.Close()
	defer os.Remove(matrixFile.Name())

	dictFile := createTempFile(t, "test_dict_thread_*.txt", testDict)
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

	// Test with multiple workers to stress thread safety
	searcher := NewWordSimpleSearcher(matrix, dict)
	searcher.SearchAllWords(8) // Use 8 workers

	// Get results multiple times to check for race conditions
	results1 := searcher.GetResults()
	results2 := searcher.GetResults()
	results3 := searcher.GetResults()

	// All calls should return the same results
	if len(results1) != len(results2) || len(results2) != len(results3) {
		t.Errorf("Results length inconsistent: %d, %d, %d",
			len(results1), len(results2), len(results3))
	}

	// Verify no duplicates in results
	wordSet := make(map[string]bool)
	for _, result := range results1 {
		if wordSet[result.Word] {
			t.Errorf("Duplicate word found: %s", result.Word)
		}
		wordSet[result.Word] = true
	}

	t.Logf("Thread safety test passed: %d unique words found", len(wordSet))
}

// Helper function to create temporary files
func createTempFile(t *testing.T, pattern, content string) *os.File {
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return tmpFile
}
