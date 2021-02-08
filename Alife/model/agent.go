package model


import (
	"errors"
	"math/rand"
	"Alife/lib"
)

// Agent implements abm.Agent and worlds.XY and
// walks randomly over 2D grid.
type Agent struct {
	id 			 int
	x, y         int
	origx, origy int
	grid         *Grid
	trail        bool // leave trail?
}

func (w *Agent) ID() int {
	return w.id
}

func NewAgent(abm *lib.ABM, id, x, y int, trail bool) (*Agent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("Agent needs a Grid world to operate")
	}
	return &Agent{
		id:    id,
		origx: x,
		origy: y,
		x:     x,
		y:     y,
		grid:  grid,
		trail: trail,
	}, nil
}

func (w *Agent) Run() {
	rx := rand.Intn(4)
	oldx, oldy := w.x, w.y
	switch rx {
	case 0:
		w.x++
	case 1:
		w.y++
	case 2:
		w.x--
	case 3:
		w.y--
	}

	var err error
	if w.trail {
		err = w.grid.Copy(oldx, oldy, w.x, w.y)
	} else {
		err = w.grid.Move(oldx, oldy, w.x, w.y)
	}

	if err != nil {
		w.x, w.y = oldx, oldy
	}
}

func (w *Agent) X() int { return w.x }
func (w *Agent) Y() int { return w.y }
