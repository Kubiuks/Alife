package model

import (
	"Alife/lib"
	"errors"
)

type Food struct {
	resource	 float64
	id 			 int
	x, y         int
	grid         *Grid
}

func (f *Food) ID() int {
	return f.id
}

func NewFood(abm *lib.ABM, x, y int) (*Food, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("agent needs a Grid world to operate")
	}
	return &Food{
		resource: 1,
		id:    -1,
		x:     x,
		y:     y,
		grid:  grid,
	}, nil
}

func (f *Food) Run() {
	f.resource -= 0.01
}

func (f *Food) X() int { return f.x }
func (f *Food) Y() int { return f.y }

