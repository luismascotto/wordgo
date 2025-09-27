package main

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
	directions  *[]string
}
