package model

import (
	"errors"
	"fmt"
	"math/rand"
)

// Agent implements lib.Agent and
// walks randomly over 2D grid.
type Agent struct {
	alive		 bool
	energy		 float32
	id 			 int
	x, y         int
	origx, origy int
	grid         *Grid
	trail        bool
	direction    int
}

func NewAgent(abm *ABM, id, x, y int, trail bool) (*Agent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("Agent needs a Grid world to operate")
	}
	return &Agent{
		alive:  true,
		energy: 1,
		id:     id,
		origx:  x,
		origy:  y,
		x:      x,
		y:      y,
		grid:   grid,
		trail:  trail,
		direction: rand.Intn(8),
	}, nil
}

func (a *Agent) Run() {
	for _,agent := range a.grid.agentVision[a.id-1] {
		if agent.ID() == -1 {
			fmt.Printf("I see food at (%v,%v)\n", agent.X(), agent.Y())
		} else if agent.ID() == -3 {
			//fmt.Printf("I see the Wall at (%v,%v)\n", agent.X(), agent.Y())
		} else {
			//fmt.Printf("I see an agent at (%v,%v)\n", agent.X(), agent.Y())
		}
	}
	//if !a.alive{
	//	a.id = 0
	//	return
	//}
	a.move()
}

func (a *Agent) move(){
	rx := rand.Intn(8)
	oldx, oldy := a.x, a.y
	oldDirection := a.direction
	switch rx {
	/*
		|5|4|3|
		|6|x|2|
		|7|0|1|
	*/
	case 0:
		a.x++
		a.direction = 2
	case 1:
		a.y++
		a.direction = 0
	case 2:
		a.x--
		a.direction = 6
	case 3:
		a.y--
		a.direction = 4
	case 4:
		a.y--
		a.x++
		a.direction = 3
	case 5:
		a.y--
		a.x--
		a.direction = 5
	case 6:
		a.y++
		a.x--
		a.direction = 7
	case 7:
		a.y++
		a.x++
		a.direction = 1
	}

	var err error
	if a.trail {
		err = a.grid.Copy(a.id, oldx, oldy, a.x, a.y)
	} else {
		err = a.grid.Move(a.id, oldx, oldy, a.x, a.y)
	}

	if err != nil {
		a.x, a.y = oldx, oldy
		a.direction = oldDirection
	} else {
		//a.energy -= 0.001
		if a.energy <= 0 {
			a.alive = false
		}
	}
}
func (a *Agent) ID() int { return a.id }
func (a *Agent) Direction() int { return a.direction }
func (a *Agent) Alive() bool { return a.alive }
func (a *Agent) X() int      { return a.x }
func (a *Agent) Y() int      { return a.y }
