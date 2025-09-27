package main

import (
	"errors"
	"strings"
)

var (
	ErrOutOfBoundariesLeft   = errors.New("out of boundaries <")
	ErrOutOfBoundariesRight  = errors.New("out of boundaries >")
	ErrOutOfBoundariesTop    = errors.New("out of boundaries ^")
	ErrOutOfBoundariesBottom = errors.New("out of boundaries v")
)

type Coord struct {
	X int
	Y int
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
