package model

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

// Agent implements lib.Agent and
// walks randomly over 2D grid.
type Agent struct {
	alive		 bool
	energy		 float64
	id 			 int
	x, y         float64
	origx, origy float64
	stepSize      float64
	grid         *Grid
	trail        bool
	direction    int
}

func NewAgent(abm *ABM, id int, x, y float64, trail bool) (*Agent, error) {
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
		stepSize: 1,
		direction: rand.Intn(360),
	}, nil
}

func (a *Agent) Run() {
	//fmt.Println(a.x, a.y)
	//fmt.Printf("I am agent: %v, and I see: %v\n", a.id, a.grid.agentVision[a.id-1])
	for _,agent := range a.grid.agentVision[a.id-1] {
		if agent.ID() == -1 {
			fmt.Printf("I see food at (%v,%v)\n", agent.X(), agent.Y())
		} else if agent.ID() == -3 {
			fmt.Printf("I see the Wall at (%v,%v)\n", agent.X(), agent.Y())
		} else {
			fmt.Printf("I see an agent at (%v,%v)\n", agent.X(), agent.Y())
		}
	}
	//if !a.alive{
	//	a.id = 0
	//	return
	//}
	a.move()
}

func (a *Agent) move(){
	oldx, oldy := a.x, a.y
	oldDirection := a.direction
	a.direction = mod(oldDirection + rand.Intn(41) - 20, 360)
	a.x = oldx + a.stepSize * math.Sin(float64(a.direction)*(math.Pi/180.0))
	a.y = oldy + a.stepSize * math.Cos(float64(a.direction)*(math.Pi/180.0))

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

func mod(a, b int) int {
	c := a % b
	if c < 0 {
		return c+b
	}
	return c
}

func (a *Agent) ID() int { return a.id }
func (a *Agent) Direction() int { return a.direction }
func (a *Agent) Alive() bool { return a.alive }
func (a *Agent) X() float64      { return a.x }
func (a *Agent) Y() float64      { return a.y }
