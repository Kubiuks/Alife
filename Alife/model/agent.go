package model

import (
	"Alife/lib"
	"errors"
	"fmt"
	"math"
	"math/rand"
)

// Agent implements lib.Agent and
// walks randomly over 2D grid.
type Agent struct {
	oxytocin	 float64
	cortisol	 float64
	alive		 bool
	energy		 float64
	id 			 int
	x, y         float64
	origx, origy float64
	stepSize      float64
	grid         *Grid
	trail        bool
	direction    int
	ch 			 chan string
}

func NewAgent(abm *lib.ABM, id int, x, y float64, ch chan string, trail bool) (*Agent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("Agent needs a Grid world to operate")
	}
	return &Agent{
		alive    : true,
		energy   : 1,
		id       : id,
		origx    : x,
		origy    : y,
		x        : x,
		y        : y,
		grid     : grid,
		trail    : trail,
		stepSize : 1,
		direction: rand.Intn(360),
		ch       : ch,
		oxytocin : 1,
		cortisol : 1,
	}, nil
}

func (a *Agent) Run() {
	if a.ch != nil {
		a.ch <- fmt.Sprintf("'[%v, %v, %v, %v]'", a.id, a.energy, a.oxytocin, a.cortisol)
	}
	for _,agent := range a.grid.agentVision[a.id-1] {
		if agent.ID() == -1 {
			//see food
		} else if agent.ID() == -3 {
			//see wall
			a.moveFromWall()
			return
		} else {
			//see agent
		}
	}
	//if !a.alive{
	//	a.id = 0
	//	return
	//}
	a.move()
}

func (a *Agent) motivation(){

}

func (a *Agent) moveFromWall(){
	oldx, oldy := a.x, a.y
	oldDirection := a.direction
	a.direction = mod(oldDirection + rand.Intn(136) - rand.Intn(136), 360)
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
		//a.energy -= 0.5
		if a.energy <= 0 {
			a.alive = false
		}
	}
}

func (a *Agent) move(){
	oldx, oldy := a.x, a.y
	oldDirection := a.direction
	a.direction = mod(oldDirection + rand.Intn(21) - rand.Intn(21), 360)
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
		a.energy -= rand.Float64()
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
