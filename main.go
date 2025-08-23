package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

// Fixed errors for reuse
var (
	ErrOutOfBoundariesLeft   = errors.New("out of boundaries <")
	ErrOutOfBoundariesRight  = errors.New("out of boundaries >")
	ErrOutOfBoundariesTop    = errors.New("out of boundaries ^")
	ErrOutOfBoundariesBottom = errors.New("out of boundaries v")
	ErrEmptyMatrix           = errors.New("matriz vazia")
	ErrFileOpen              = errors.New("erro ao abrir arquivo")
	ErrFileRead              = errors.New("erro ao ler arquivo")
)

// LetterMatrix representa a matriz de letras onde as palavras serão buscadas
type LetterMatrix struct {
	matrix [][]rune
	rows   int
	cols   int
}

// Dictionary representa o dicionário de palavras para busca
type Dictionary struct {
	words map[string]bool
	trie  *TrieNode
}

// TrieNode representa um nó na árvore trie para busca de prefixos
type TrieNode struct {
	children map[rune]*TrieNode
	isWord   bool
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

// NewDictionary cria um novo dicionário a partir de um arquivo
func NewDictionary(filename string) (*Dictionary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%s do dicionário: %w", ErrFileOpen, err)
	}
	defer file.Close()

	dict := &Dictionary{
		words: make(map[string]bool),
		trie:  &TrieNode{children: make(map[rune]*TrieNode)},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(strings.ToUpper(scanner.Text()))
		if word != "" && len(word) >= 3 {
			dict.words[word] = true
			dict.insertIntoTrie(word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%s do dicionário: %w", ErrFileRead, err)
	}

	return dict, nil
}

// insertIntoTrie insere uma palavra na árvore trie
func (d *Dictionary) insertIntoTrie(word string) {
	node := d.trie
	for _, char := range word {
		if node.children[char] == nil {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isWord = true
}

// Contains verifica se uma palavra existe no dicionário
func (d *Dictionary) Contains(word string) bool {
	return d.words[strings.ToUpper(word)]
}

// Contains verifica se uma palavra existe no dicionário
func (d *Dictionary) ContainsUpped(uppedWord string) bool {
	return d.words[uppedWord]
}

// IsPrefix verifica se uma sequência é prefixo de alguma palavra válida
func (d *Dictionary) IsPrefix(sequence string) bool {
	node := d.trie
	for _, char := range sequence {
		if node.children[char] == nil {
			return false // Não é prefixo
		}
		node = node.children[char]
	}
	return true // É prefixo (pode ter filhos)
}

// IsWord verifica se uma sequência é uma palavra completa
func (d *Dictionary) IsWord(sequence string) bool {
	node := d.trie
	for _, char := range sequence {
		if node.children[char] == nil {
			return false // Não existe
		}
		node = node.children[char]
	}
	return node.isWord // É uma palavra completa
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

// PrintDictionaryStats imprime estatísticas do dicionário
func (d *Dictionary) PrintDictionaryStats() {
	fmt.Printf("Dicionário carregado com %d palavras\n", len(d.words))
}

// WordResult representa uma palavra encontrada na matriz
type WordResult struct {
	Word      string
	StartRow  int
	StartCol  int
	Direction string
	Length    int
}

// Direction representa uma direção de busca
type Direction struct {
	Name     string
	Symbol   string
	DeltaRow int
	DeltaCol int
}

// WordSearcher representa o sistema de busca de palavras
type WordSearcher struct {
	matrix     *LetterMatrix
	dictionary *Dictionary
	directions []Direction
	results    []WordResult
	seen       map[string]bool // Para evitar duplicatas
	mutex      sync.Mutex
}

// T B L R (Top, Bottom, Left, Right)

type Cell struct {
	Letter      rune
	DoubleScore bool //Only on V2
	RowEater    bool //Only on V2
}
type Word struct {
	word        []rune
	coordinates []Coord
}
type Coord struct {
	X int
	Y int
}

var board [][]Cell
var chFoundWords chan []rune
var chFinalize chan bool

// IsValidWord checks if a sequence of runes forms a valid word
func IsValidWord(runes []rune) bool {
	word := string(runes)
	// For now, just check if it's at least 3 characters
	// You can integrate this with your dictionary later
	return len(word) >= 3
}

// IsValidWordPreffix checks if a sequence of runes is a valid prefix
func IsValidWordPreffix(runes []rune) bool {
	word := string(runes)
	// For now, just check if it's at least 1 character
	// You can integrate this with your dictionary later
	return len(word) >= 1
}

func MyMainFunc() {
	start := &Word{}

	start.word = make([]rune, 0)
	start.coordinates = make([]Coord, 0)
	coord := &Coord{
		X: rand.Intn(len(board)),
		Y: rand.Intn(len(board[0])),
	}

	start.word = append(start.word, board[coord.X][coord.Y].Letter)
	start.coordinates = append(start.coordinates, *coord)

	go toWalk(*start)

	<-chFinalize
}

func toWalk(word Word) {
	if word.canWalk(L) {
		go toWalk(word)
	}
	if word.canWalk(TL) {
		go toWalk(word)
	}
	if word.canWalk(T) {
		go toWalk(word)
	}
	if word.canWalk(TR) {
		go toWalk(word)
	}
	if word.canWalk(R) {
		go toWalk(word)
	}
	if word.canWalk(BR) {
		go toWalk(word)
	}
	if word.canWalk(B) {
		go toWalk(word)
	}
	if word.canWalk(BL) {
		go toWalk(word)
	}
}

// T B L R
func (w *Word) canWalk(pos string) bool {
	newCoord, err := w.coordinates[len(w.coordinates)-1].next(pos)

	if err != nil {
		return false
	}

	if w.hasVisitedCell(*newCoord) {
		return false
	}

	w.word = append(w.word, board[newCoord.X][newCoord.Y].Letter)
	w.coordinates = append(w.coordinates, *newCoord)

	if IsValidWord(w.word) {
		chFoundWords <- w.word
	}

	if !IsValidWordPreffix(w.word) {
		return false
	}

	return true
}

// hasVisitedCell checks if a coordinate was already visited by walking backwards through the path
func (w *Word) hasVisitedCell(coord Coord) bool {
	// Walk backwards through the coordinates to check for repeated visits
	for i := len(w.coordinates) - 1; i >= 0; i-- {
		if w.coordinates[i].X == coord.X && w.coordinates[i].Y == coord.Y {
			return true
		}
	}
	return false
}

func (c *Coord) next(pos string) (*Coord, error) {
	// Use simple variables for updates and validation
	newX, newY := c.X, c.Y

	if strings.Contains(pos, L) {
		newX--
		if newX < 0 {
			return nil, ErrOutOfBoundariesLeft
		}
	}
	if strings.Contains(pos, R) {
		newX++
		if newX >= len(board) {
			return nil, ErrOutOfBoundariesRight
		}
	}
	if strings.Contains(pos, T) {
		newY--
		if newY < 0 {
			return nil, ErrOutOfBoundariesTop
		}
	}
	if strings.Contains(pos, B) {
		newY++
		if newY >= len(board[0]) {
			return nil, ErrOutOfBoundariesBottom
		}
	}

	// Only create Coord object if all validations pass
	return &Coord{X: newX, Y: newY}, nil
}

const (
	L  = "L"
	R  = "R"
	T  = "T"
	B  = "B"
	TL = "TL"
	TR = "TR"
	BL = "BL"
	BR = "BR"
)

// NewWordSearcher cria um novo buscador de palavras
func NewWordSearcher(matrix *LetterMatrix, dictionary *Dictionary) *WordSearcher {
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

func main() {
	fmt.Println("=== WordGo - Buscador de Palavras em Matriz de Letras ===")

	// Carregar matriz de letras
	fmt.Println("Carregando matriz de letras...")
	matrix, err := NewLetterMatrix("res/example.txt")
	if err != nil {
		log.Fatalf("Erro ao carregar matriz: %v", err)
	}
	matrix.PrintMatrix()
	fmt.Println()

	// Carregar dicionário
	fmt.Println("Carregando dicionário...")
	dict, err := NewDictionary("res/words.txt")
	if err != nil {
		log.Fatalf("Erro ao carregar dicionário: %v", err)
	}
	dict.PrintDictionaryStats()
	fmt.Println()

	// Iniciar busca de palavras
	fmt.Println("\n=== Iniciando Busca de Palavras ===")
	searcher := NewWordSearcher(matrix, dict)

	// Usar 4 workers para processamento paralelo
	numWorkers := 4
	fmt.Printf("Iniciando busca com %d workers...\n", numWorkers)
	searcher.SearchAllWords(numWorkers)

	// Exibir resultados
	searcher.PrintResults()
}
