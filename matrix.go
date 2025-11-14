package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// LetterMatrix representa a matriz de letras onde as palavras serão buscadas
type LetterMatrix struct {
	matrix          [][]rune
	rows            int
	cols            int
	specials        []string
	special_letters []SpecialLetter
}

type SpecialType int

const (
	SpecialTypeMandatory = iota
	SpecialTypeOptional
)

type SpecialLetter struct {
	special_type string
	coordinate   Coord
}

// NewLetterMatrixFromFile cria uma nova matriz de letras a partir de um arquivo
func NewLetterMatrixFromFile(filename string) (*LetterMatrix, error) {
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

	return newLetterMatrixFromRuneMatrix(matrix, maxCols), nil
}

func NewLetterMatrixFromString(matrixString string) (*LetterMatrix, error) {
	matrixLines := strings.Split(matrixString, "\n")
	var matrix [][]rune
	var maxCols int
	for _, line := range matrixLines {
		row := []rune(line)
		if len(row) > maxCols {
			maxCols = len(row)
		}
		matrix = append(matrix, []rune(line))
	}

	if len(matrix) == 0 {
		return nil, ErrEmptyMatrix
	}

	return newLetterMatrixFromRuneMatrix(matrix, maxCols), nil
}

func newLetterMatrixFromRuneMatrix(matrix [][]rune, maxCols int) *LetterMatrix {
	specialCellsCoordStrings := []string{}

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

		for pos, cell := range matrix[i] {
			//runes.upper(cell)
			if unicode.IsUpper(cell) {
				specialCellsCoordStrings = append(specialCellsCoordStrings, fmt.Sprintf("(%d,%d)", i+1, pos+1))
			}
		}
	}

	return &LetterMatrix{
		matrix:   matrix,
		rows:     len(matrix),
		cols:     len(matrix[0]),
		specials: specialCellsCoordStrings,
	}
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

// RemoverLetras remove as letras e reorganiza a matriz
func (lm *LetterMatrix) RemoveLetters(coordinates []Coord) {
	for _, coord := range coordinates {
		lm.matrix[coord.X][coord.Y] = '-'
	}
	// collapse Y moving down the letters
	for i := 0; i < lm.rows; i++ {
		for j := 0; j < lm.cols; j++ {
			if lm.matrix[i][j] == '-' {
				for k := i; k > 0; k-- {
					lm.matrix[k][j] = lm.matrix[k-1][j]
				}
				lm.matrix[0][j] = ' '
			}
		}
	}
}
