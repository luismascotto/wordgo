package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Fixed errors for reuse
var (
	ErrEmptyMatrix = errors.New("matriz vazia")
	ErrFileOpen    = errors.New("erro ao abrir arquivo")
	ErrFileRead    = errors.New("erro ao ler arquivo")
)

const (
	SIMULATE_SINGLE_THREAD = false
	MAX_GOROUTINES         = 32
	MIN_WORD_LENGTH        = 8
	MODE_SQUARE_SEARCH     = true
)

//var chFoundWords chan []rune

//var chRoutineFinalize chan bool

var foundWords map[string]bool
var foundWordsMutex sync.Mutex

func main() {
	fmt.Println("=== WordGo - Buscador de Palavras em Matriz de Letras ===")

	// Definir flag para arquivo de matriz
	matrixFile := flag.String("matrix", "res/example.txt", "Arquivo de matriz de letras para carregar")
	flag.Parse()

	// Validar se o arquivo especificado existe
	if _, err := os.Stat(*matrixFile); os.IsNotExist(err) {
		log.Fatalf("Arquivo de matriz não encontrado: %s", *matrixFile)
	}

	// Carregar matriz de letras
	fmt.Printf("Carregando matriz de letras de: %s\n", *matrixFile)
	matrix, err := NewLetterMatrix(*matrixFile)
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

	directions := NewDirections()

	// Iniciar busca de palavras
	fmt.Println("\n=== Iniciando Busca de Palavras ===")

	if os.Getenv("CFG_SIMPLE") == "true" {
		simpleSearcher := NewWordSimpleSearcher(matrix, dict)

		// Usar 4 workers para processamento paralelo
		numWorkers := 4
		fmt.Printf("Iniciando busca com %d workers...\n", numWorkers)
		simpleSearcher.SearchAllWords(numWorkers)

		// Exibir resultados
		simpleSearcher.PrintResults()
		time.Sleep(5000 * time.Millisecond)
		return
	}

	dimX, dimY := matrix.GetDimensions()

	allFoundWordsList := make([]string, 0, 128)

	for startX := range dimX {
		for startY := range dimY {
			fmt.Printf("(%d,%d) -> ", startX+1, startY+1)
			if matrix.GetMatrix()[startX][startY] == ' ' {
				continue
			}

			start := &Word{}

			start.word = make([]rune, 0)
			start.coordinates = make([]Coord, 0)
			start.matrix = matrix
			start.dictionary = dict
			start.directions = directions
			coord := &Coord{
				X: startX,
				Y: startY,
			}
			start.word = append(start.word, matrix.GetMatrix()[startX][startY])
			start.coordinates = append(start.coordinates, *coord)

			// Initialize foundWords map
			foundWords = make(map[string]bool)

			var wg sync.WaitGroup
			limitGoroutines := make(chan struct{}, MAX_GOROUTINES)
			wg.Go(func() {
				toWalk(*start, limitGoroutines)
			})
			//fmt.Println("Waiting for words to be found...")
			wg.Wait()

			if len(foundWords) > 0 {
				fmt.Printf("found words: ")

				// Print all found words
				foundWordsList := make([]string, 0, len(foundWords))
				for word := range foundWords {
					foundWordsList = append(foundWordsList, word)
				}
				allFoundWordsList = append(allFoundWordsList, foundWordsList...)

				sortAndPrint(foundWordsList, []string{})
				// } else {
				// 	fmt.Println("No words found...")
			} else {
				fmt.Println()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("All found words:")
	sortAndPrint(allFoundWordsList, matrix.specials)
	time.Sleep(5000 * time.Millisecond)
	if len(matrix.specials) > 0 {
		filteredWordsList := make([]string, 0)
		for _, word := range allFoundWordsList {
			for _, special := range matrix.specials {
				if strings.Contains(word, special) {
					filteredWordsList = append(filteredWordsList, word)
					break
				}
			}
		}
		fmt.Println("\n\n\nAll found words in specials:")
		if len(filteredWordsList) == 0 {
			fmt.Printf("no words found... BOOO HOOO")
			time.Sleep(5000 * time.Millisecond)
		} else {
			sortAndPrint(filteredWordsList, matrix.specials)
			time.Sleep(5000 * time.Millisecond)
		}
	}
}

func addFoundWord(word string) {
	foundWordsMutex.Lock()
	defer foundWordsMutex.Unlock()
	foundWords[word] = true
}

func toWalk(word Word, limitGoroutines chan struct{}) {
	//word.PrintBreadCrumb()
	var wg sync.WaitGroup
	for _, dir := range *word.directions {
		//Clone the word
		newWord := Word{
			word:        make([]rune, len(word.word)),
			coordinates: make([]Coord, len(word.coordinates)),
			matrix:      word.matrix,
			dictionary:  word.dictionary,
			directions:  word.directions,
		}
		copy(newWord.word, word.word)
		copy(newWord.coordinates, word.coordinates)
		if newWord.canWalk(dir) {
			//time.Sleep(10 * time.Millisecond)
			//fmt.Printf(">")
			limitGoroutines <- struct{}{}
			wg.Go(func() {
				toWalk(newWord, limitGoroutines)
				//fmt.Printf("\b")
				<-limitGoroutines
				//time.Sleep(10 * time.Millisecond)
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
	stringWord := strings.ToUpper(string(w.word))
	if w.dictionary.IsWord(stringWord) {

		fullPathWalked := ""
		for _, coord := range w.coordinates {
			fullPathWalked += fmt.Sprintf("(%d,%d)", coord.X+1, coord.Y+1)
		}

		wordAndFullPathWalked := fmt.Sprintf("%s %s", stringWord, fullPathWalked)
		addFoundWord(wordAndFullPathWalked)
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

func sortAndPrint(allFoundWordsList []string, preferredCoordinates []string) {
	if len(allFoundWordsList) == 0 {
		fmt.Println("No words found...")
		return
	}

	sort.Slice(allFoundWordsList, func(i, j int) bool {

		found_i := 0
		found_j := 0
		for _, coord := range preferredCoordinates {
			if strings.Contains(allFoundWordsList[i], coord) {
				found_i++
			}
			if strings.Contains(allFoundWordsList[j], coord) {
				found_j++
			}
		}
		if found_i > found_j {
			return true
		}
		if found_i < found_j {
			return false
		}
		return len(allFoundWordsList[i]) > len(allFoundWordsList[j])
	})
	maxLength := len(allFoundWordsList[0])
	count := 0
	for _, word := range allFoundWordsList {
		fmt.Printf("%-*.*s ", maxLength, maxLength, word)
		if maxLength > 70 {
			count = 2
		} else if maxLength > 50 {
			if count == 0 {
				count = 1
			}
		}
		if maxLength > 20 {
			maxLength = len(word)
		}
		if (count+1)%3 == 0 {
			fmt.Println()
		}
		count++
	}
	fmt.Println()
	//fmt.Println()
	//time.Sleep(1000 * time.Millisecond)
}
