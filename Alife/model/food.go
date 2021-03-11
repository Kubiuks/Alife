package model

import (
	"Alife/lib"
	"errors"
	"sync"
)

type Food struct {
	// model
	alive		 bool
	resource	 float64
	maxResource  float64
	owner		 *Agent
	// immplementation
	mutex 		 sync.Mutex
	id 			 int
	x, y         float64
	grid         *Grid
}

func NewFood(abm *lib.ABM, x, y float64) (*Food, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("agent needs a Grid world to operate")
	}
	return &Food{
		alive: true,
		resource: 4,
		maxResource: 4,
		owner: nil,
		id:    -1,
		x:     x,
		y:     y,
		grid:  grid,
	}, nil
}

func (f *Food) Run() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if !f.alive {
		return
	}
	if f.resource < f.maxResource {
		f.resource += 0.004
	}
}

func (f *Food) reduceResource(amount float64){
	f.mutex.Lock()
	f.resource -= amount
	if f.resource <=0 {
		f.alive = false
		f.id = -4
	}
	f.mutex.Unlock()
}

func (f *Food) Resource() float64{
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.resource
}
func (f *Food) SetOwner(agent *Agent) {
	f.mutex.Lock()
	f.owner = agent
	f.mutex.Unlock()
}
func (f *Food) Owner() *Agent { return f.owner }
func (f *Food) ID() int { return f.id }
func (f *Food) X() float64 { return f.x }
func (f *Food) Y() float64 { return f.y }

