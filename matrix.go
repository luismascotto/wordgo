package main

import (
	"bufio"
	"fmt"
	"os"
)

// LetterMatrix representa a matriz de letras onde as palavras serão buscadas
type LetterMatrix struct {
	matrix [][]rune
	rows   int
	cols   int
}

// NewLetterMatrix cria uma nova matriz de letras a partir de um arquivo
func NewLetterMatrix(filename string) (*LetterMatrix, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%s da matriz: %w", ErrFileOpen, err)
	}
	defer file.Close()

	var matrix [][]rune
	var maxCols int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			row := []rune(line)
			if len(row) > maxCols {
				maxCols = len(row)
			}
			matrix = append(matrix, row)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%s da matriz: %w", ErrFileRead, err)
	}

	if len(matrix) == 0 {
		return nil, ErrEmptyMatrix
	}

	// Padronizar todas as linhas para ter o mesmo comprimento
	for i, row := range matrix {
		if len(row) < maxCols {
			// Right-pad com espaços
			paddedRow := make([]rune, maxCols)
			copy(paddedRow, row)
			for j := len(row); j < maxCols; j++ {
				paddedRow[j] = ' '
			}
			matrix[i] = paddedRow
		}
	}

	cols := maxCols

	return &LetterMatrix{
		matrix: matrix,
		rows:   len(matrix),
		cols:   cols,
	}, nil
}

// GetMatrix retorna a matriz de letras
func (lm *LetterMatrix) GetMatrix() [][]rune {
	return lm.matrix
}

// GetDimensions retorna as dimensões da matriz
func (lm *LetterMatrix) GetDimensions() (int, int) {
	return lm.rows, lm.cols
}

// PrintMatrix imprime a matriz de letras
func (lm *LetterMatrix) PrintMatrix() {
	fmt.Println("Matriz de Letras:")
	fmt.Printf("Dimensões: %dx%d\n\n", lm.rows, lm.cols)

	for i, row := range lm.matrix {
		fmt.Printf("%2d: %s\n", i, string(row))
	}
}
