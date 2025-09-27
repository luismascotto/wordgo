package main

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
)

// WordSearcher representa o sistema de busca de palavras
type WordSearcher struct {
	matrix     *LetterMatrix
	dictionary *Dictionary
	directions []Direction
	results    []WordResult
	seen       map[string]bool // Para evitar duplicatas
	mutex      sync.Mutex
}

// NewWordSimpleSearcher cria um novo buscador de palavras
func NewWordSimpleSearcher(matrix *LetterMatrix, dictionary *Dictionary) *WordSearcher {
	return &WordSearcher{
		matrix:     matrix,
		dictionary: dictionary,
		directions: []Direction{
			{R, "→", 0, 1},
			{L, "←", 0, -1},
			{B, "↓", 1, 0},
			{T, "↑", -1, 0},
			{BR, "↘", 1, 1},
			{BL, "↙", 1, -1},
			{TR, "↗", -1, 1},
			{TL, "↖", -1, -1},
		},
		results: make([]WordResult, 0),
		seen:    make(map[string]bool),
	}
}

// SearchFromPosition busca palavras a partir de uma posição específica em uma direção
func (ws *WordSearcher) SimpleSearchFromPosition(startRow, startCol int, direction Direction) {
	matrix := ws.matrix.GetMatrix()
	rows, cols := ws.matrix.GetDimensions()

	var currentWord strings.Builder
	row, col := startRow, startCol

	// Buscar na direção especificada
	for row >= 0 && row < rows && col >= 0 && col < cols {
		char := matrix[row][col]

		// Parar se encontrar espaço
		if char == ' ' {
			break
		}

		currentWord.WriteRune(char)
		sequence := currentWord.String()

		// Verificar se é um prefixo válido
		if !ws.dictionary.IsPrefix(sequence) {
			break // Não há palavras que começam com esta sequência
		}

		// Verificar se é uma palavra válida (mínimo 3 caracteres)
		if len(sequence) >= 3 && ws.dictionary.IsWord(sequence) {
			ws.addResult(WordResult{
				Word:      sequence,
				StartRow:  startRow,
				StartCol:  startCol,
				Direction: direction.Name,
				Length:    len(sequence),
			})
		}

		// Mover para a próxima posição na direção
		row += direction.DeltaRow
		col += direction.DeltaCol
	}
}

// addResult adiciona um resultado de forma thread-safe
func (ws *WordSearcher) addResult(result WordResult) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	// Criar chave única para evitar duplicatas
	key := fmt.Sprintf("%s_%d_%d_%s", result.Word, result.StartRow, result.StartCol, result.Direction)

	// Só adicionar se não vimos esta combinação antes
	if !ws.seen[key] {
		ws.seen[key] = true
		ws.results = append(ws.results, result)
	}
}

// SearchAllWords busca todas as palavras na matriz usando goroutines
func (ws *WordSearcher) SearchAllWords(numWorkers int) {
	rows, cols := ws.matrix.GetDimensions()

	// Canal para distribuir trabalho
	jobs := make(chan [3]int, rows*cols*len(ws.directions))

	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				startRow, startCol, dirIndex := job[0], job[1], job[2]
				ws.SimpleSearchFromPosition(startRow, startCol, ws.directions[dirIndex])
			}
		}()
	}

	// Distribuir trabalho
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			for dirIndex := range ws.directions {
				jobs <- [3]int{row, col, dirIndex}
			}
		}
	}
	close(jobs)

	// Aguardar todos os workers terminarem
	wg.Wait()
}

// GetResults retorna todos os resultados encontrados
func (ws *WordSearcher) GetResults() []WordResult {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return append([]WordResult{}, ws.results...)
}

// PrintResults imprime os resultados da busca
func (ws *WordSearcher) PrintResults() {
	results := ws.GetResults()
	fmt.Printf("\n=== Resultados da Busca ===\n")
	fmt.Printf("Total de palavras encontradas: %d\n\n", len(results))

	// Agrupar por direção
	byDirection := make(map[string][]WordResult)
	for _, result := range results {
		byDirection[result.Direction] = append(byDirection[result.Direction], result)
	}

	for direction, words := range byDirection {
		fmt.Printf("%s (%d palavras):\n", direction, len(words))
		for _, word := range words {
			fmt.Printf("  '%s' em (%d,%d) - %d letras\n",
				word.Word, word.StartRow, word.StartCol, word.Length)
		}
		fmt.Println()
	}
}

func (w *Word) PrintBreadCrumb() {
	gr := bytes.Fields(debug.Stack())[1]
	fmt.Printf("%s: ", string(gr))
	for _, coord := range w.coordinates {
		fmt.Printf("(%d,%d) ", coord.X, coord.Y)
	}
	fmt.Printf("[%s] ", string(w.word))
	fmt.Println()
}
