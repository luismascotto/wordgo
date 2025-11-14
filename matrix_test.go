package main

import (
	"os"
	"strings"
	"testing"
)

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
	matrix, err := NewLetterMatrixFromFile(tmpFile.Name())
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
	matrix, err := NewLetterMatrixFromFile(tmpFile.Name())
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

func TestRemoveLetter(t *testing.T) {
	matrixString := "abc\ndef\nGhi"
	matrix, err := NewLetterMatrixFromString(matrixString)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}

	// Remove coord (2,0) - row 2, column 0 (the 'G')
	matrix.RemoveLetters([]Coord{{X: 2, Y: 0}})

	// Get the last row
	matrixData := matrix.GetMatrix()
	lastRow := string(matrixData[len(matrixData)-1])

	// Debug: print full matrix
	t.Logf("Matrix after removal:")
	for i, row := range matrixData {
		t.Logf("Row %d: %s", i, string(row))
	}
	t.Logf("Specials: %v", matrix.specials)

	// Verify the last row is "DHI"
	if strings.ToUpper(lastRow) != "DHI" {
		t.Errorf("Expected last row to be 'DHI', got '%s'", lastRow)
	}
}

func TestRemoveLetters(t *testing.T) {
	matrixString := "aBc\ndef\nghi"
	matrix, err := NewLetterMatrixFromString(matrixString)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}
	t.Logf("Matrix before removal:")
	for i, row := range matrix.GetMatrix() {
		t.Logf("Row %d: %s", i, string(row))
	}
	t.Logf("Specials: %v", matrix.specials)
	// Remove 'GEC'
	matrix.RemoveLetters([]Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 2}})

	// Get the last row
	matrixData := matrix.GetMatrix()

	// Debug: print full matrix
	t.Logf("Matrix after removal:")
	for i, row := range matrixData {
		t.Logf("Row %d: %s", i, string(row))
	}
	t.Logf("Specials: %v", matrix.specials)
	firstRow := string(matrixData[0])
	secondRow := string(matrixData[1])
	lastRow := string(matrixData[2])
	// Verify the first row is "ABC"
	if strings.ToUpper(firstRow) != "A  " {
		t.Errorf("Expected first row to be 'A  ', got '%s'", firstRow)
	}
	// Verify the second row is "DEF"
	if strings.ToUpper(secondRow) != "D F" {
		t.Errorf("Expected second row to be 'D F', got '%s'", secondRow)
	}
	// Verify the last row is "DHI"
	if strings.ToUpper(lastRow) != "GBI" {
		t.Errorf("Expected last row to be 'GBI', got '%s'", lastRow)
	}
}
