package main

import (
	"testing"
)

func TestRemoveLetter(t *testing.T) {
	matrixString := "ABC\nDEF\nGHI"
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

	// Verify the last row is "DHI"
	if lastRow != "DHI" {
		t.Errorf("Expected last row to be 'DHI', got '%s'", lastRow)
	}
}

func TestRemoveLetters(t *testing.T) {
	matrixString := "ABC\nDEF\nGHI"
	matrix, err := NewLetterMatrixFromString(matrixString)
	if err != nil {
		t.Fatalf("Failed to create matrix: %v", err)
	}

	// Remove 'GEC'
	matrix.RemoveLetters([]Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 2}})

	// Get the last row
	matrixData := matrix.GetMatrix()

	// Debug: print full matrix
	t.Logf("Matrix after removal:")
	for i, row := range matrixData {
		t.Logf("Row %d: %s", i, string(row))
	}
	firstRow := string(matrixData[0])
	secondRow := string(matrixData[1])
	lastRow := string(matrixData[2])
	// Verify the first row is "ABC"
	if firstRow != "A  " {
		t.Errorf("Expected first row to be 'A  ', got '%s'", firstRow)
	}
	// Verify the second row is "DEF"
	if secondRow != "D F" {
		t.Errorf("Expected second row to be 'D F', got '%s'", secondRow)
	}
	// Verify the last row is "DHI"
	if lastRow != "GBI" {
		t.Errorf("Expected last row to be 'GBI', got '%s'", lastRow)
	}
}
