package main

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

func NewDirections() *[]string {

	Directions = make([]string, 8)
	Directions = append(Directions, L)
	Directions = append(Directions, TL)
	Directions = append(Directions, T)
	Directions = append(Directions, TR)
	Directions = append(Directions, R)
	Directions = append(Directions, BR)
	Directions = append(Directions, B)
	Directions = append(Directions, BL)

	return &Directions
}
