package model

import (
	"errors"
	"sync"
)

type Food struct {
	mutex 		 sync.RWMutex
	resource	 float32
	id 			 int
	x, y         int
	grid         *Grid
}

func NewFood(abm *ABM, x, y int) (*Food, error) {
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
	f.mutex.Lock()
	if f.resource < 1 {
		f.resource += 0.01
	}
	f.mutex.Unlock()
}

func (f *Food) Resource() float32{
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.resource
}
func (f *Food) ID() int { return f.id }
func (f *Food) X() int { return f.x }
func (f *Food) Y() int { return f.y }

