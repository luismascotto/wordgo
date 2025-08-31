package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
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

const (
	SIMULATE_SINGLE_THREAD = false
	MAX_GOROUTINES         = 32
)

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
	matrix      *LetterMatrix
	dictionary  *Dictionary
}
type Coord struct {
	X int
	Y int
}

//var chFoundWords chan []rune

//var chRoutineFinalize chan bool

var foundWords map[string]bool
var foundWordsMutex sync.Mutex

func addFoundWord(word string) {
	foundWordsMutex.Lock()
	defer foundWordsMutex.Unlock()
	foundWords[word] = true
}

func toWalk(word Word, limitGoroutines chan struct{}) {
	//word.PrintBreadCrumb()
	var wg sync.WaitGroup
	for _, dir := range Directions {
		//Clone the word
		newWord := Word{
			word:        make([]rune, len(word.word)),
			coordinates: make([]Coord, len(word.coordinates)),
			matrix:      word.matrix,
			dictionary:  word.dictionary,
		}
		copy(newWord.word, word.word)
		copy(newWord.coordinates, word.coordinates)
		if newWord.canWalk(dir) {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf(">")
			limitGoroutines <- struct{}{}
			wg.Go(func() {
				toWalk(newWord, limitGoroutines)
				fmt.Printf("\b")
				<-limitGoroutines
				time.Sleep(10 * time.Millisecond)
			})
			wg.Wait()
		}
	}
}

// T B L R
func (w *Word) canWalk(toPosition string) bool {
	rows, cols := w.matrix.GetDimensions()
	newCoord, err := w.coordinates[len(w.coordinates)-1].next(toPosition, rows, cols)

	if err != nil {
		return false
	}

	if w.hasVisitedCell(*newCoord) {
		return false
	}

	w.word = append(w.word, w.matrix.GetMatrix()[newCoord.X][newCoord.Y])
	w.coordinates = append(w.coordinates, *newCoord)
	stringWord := string(w.word)
	if w.dictionary.IsWord(stringWord) {
		addFoundWord(stringWord)
	}

	return w.dictionary.IsPrefix(stringWord)
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

func (c *Coord) next(pos string, rows int, cols int) (*Coord, error) {
	// Use simple variables for updates and validation
	newRow, newCol := c.X, c.Y

	if strings.Contains(pos, T) {
		newRow--
		if newRow < 0 {
			return nil, ErrOutOfBoundariesLeft
		}
	}
	if strings.Contains(pos, B) {
		newRow++
		if newRow >= rows {
			return nil, ErrOutOfBoundariesRight
		}
	}
	if strings.Contains(pos, L) {
		newCol--
		if newCol < 0 {
			return nil, ErrOutOfBoundariesTop
		}
	}
	if strings.Contains(pos, R) {
		newCol++
		if newCol >= cols {
			return nil, ErrOutOfBoundariesBottom
		}
	}

	// Only create Coord object if all validations pass
	return &Coord{X: newRow, Y: newCol}, nil
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

// slice of directions
var Directions []string

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

	if os.Getenv("CFG_SIMPLE") == "true" {

		// Usar 4 workers para processamento paralelo
		numWorkers := 4
		fmt.Printf("Iniciando busca com %d workers...\n", numWorkers)
		searcher.SearchAllWords(numWorkers)

		// Exibir resultados
		searcher.PrintResults()
	} else {
		Directions = make([]string, 8)
		Directions = append(Directions, L)
		Directions = append(Directions, TL)
		Directions = append(Directions, T)
		Directions = append(Directions, TR)
		Directions = append(Directions, R)
		Directions = append(Directions, BR)
		Directions = append(Directions, B)
		Directions = append(Directions, BL)

		for i := 0; i < 10; i++ {

			start := &Word{}

			start.word = make([]rune, 0)
			start.coordinates = make([]Coord, 0)
			start.matrix = matrix
			start.dictionary = dict
			coord := &Coord{}
			dimX, dimY := matrix.GetDimensions()

			for {
				coord.X, coord.Y = rand.Intn(dimX), rand.Intn(dimY)
				if matrix.GetMatrix()[coord.X][coord.Y] != ' ' {
					break
				}
			}
			//coord.X, coord.Y = rand.Intn(coord.X), rand.Intn(coord.Y)

			start.word = append(start.word, matrix.GetMatrix()[coord.X][coord.Y])
			start.coordinates = append(start.coordinates, *coord)

			// Initialize foundWords map
			foundWords = make(map[string]bool)

			var wg sync.WaitGroup
			limitGoroutines := make(chan struct{}, MAX_GOROUTINES)
			wg.Go(func() {
				toWalk(*start, limitGoroutines)
			})
			fmt.Println("Waiting for words to be found...")
			wg.Wait()

			if len(foundWords) > 0 {
				// Print all found words
				fmt.Println("Found words:")
				foundWordsList := make([]string, 0, len(foundWords))
				for word := range foundWords {
					foundWordsList = append(foundWordsList, word)
				}
				sort.Slice(foundWordsList, func(i, j int) bool {
					return len(foundWordsList[i]) > len(foundWordsList[j])
				})
				//Print in 3 columns, paddind with 2 spaces based on the length of the longest word
				maxLength := len(foundWordsList[0])
				for count, word := range foundWordsList {
					fmt.Printf("%-*.*s ", maxLength, maxLength, word)
					if (count+1)%3 == 0 {
						fmt.Println()
					}
				}
				fmt.Println()
				fmt.Println()
				time.Sleep(1000 * time.Millisecond)
			} else {
				fmt.Println("No words found")
			}
		}
	}
}
