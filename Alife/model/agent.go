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
	direction    float64
	ch 			 chan string
}

func NewAgent(abm *lib.ABM, id int, x, y float64, ch chan string, trail bool) (*Agent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("agent needs a Grid world to operate")
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
		direction: rand.Float64()*360,
		ch       : ch,
		oxytocin : 1,
		cortisol : 1,
	}, nil
}

func (a *Agent) Run() {
	if a.ch != nil {
		a.ch <- fmt.Sprintf("'[%v, %v, %v, %v]'", a.id, a.energy, a.oxytocin, a.cortisol)
	}
	if !a.alive{
		return
	}
	for _,agent := range a.grid.agentVision[a.id-1] {
		if agent.ID() == -1 {
			// the agent sees food
			if a.energy < 0.5 {
				distance := distance(a.x, a.y, agent.X(), agent.Y())
				if distance <= 1 {
					a.eatFood(agent.(*Food))
					return
				} else {
					a.moveToFood(agent.(*Food))
					return
				}
			}
		} else if agent.ID() == -3 {
			// the agent sees a wall
			a.turnFromWall()
			return
		} else if agent.ID() == -4{
			// the agent sees empty food
		} else {
			// the agent sees another agent
		}
	}
	a.randomMove()
}

func (a *Agent) motivation(){

}

func (a *Agent) eatFood(f *Food){
	if f.resource > 0{
		f.reduceResource(0.01)
		a.energy += 0.01
	}
}

func (a *Agent) turnFromWall(){
	a.direction = mod(a.direction + rand.Float64()*135 - rand.Float64()*135, 360)
}

func (a *Agent) moveToFood(f *Food){
	a.move(math.Atan2(f.x - a.x, f.y - a.y)*(180.0/math.Pi))
}

func (a *Agent) randomMove(){
	a.move(mod(a.direction + rand.Float64()*20 - rand.Float64()*20, 360))
}

func (a *Agent) move(direction float64){
	oldx, oldy := a.x, a.y
	oldDirection := a.direction
	a.direction = direction
	a.x = oldx + a.stepSize * math.Sin(a.direction*(math.Pi/180.0))
	a.y = oldy + a.stepSize * math.Cos(a.direction*(math.Pi/180.0))

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
		a.energy -= 0.0003
		if a.energy <= 0 {
			a.alive = false
			a.grid.ClearCell(a.x, a.y, a.id)
		}
	}
}

func mod(a, b float64) float64 {
	c := math.Mod(a, b)
	if c < 0 {
		return c+b
	}
	return c
}

func distance(x1, y1, x2, y2 float64) float64{
	return math.Sqrt(math.Pow(x2-x1,2) + math.Pow(y2-y1,2))
}

func (a *Agent) ID() int { return a.id }
func (a *Agent) Direction() float64 { return a.direction }
func (a *Agent) Alive() bool { return a.alive }
func (a *Agent) X() float64      { return a.x }
func (a *Agent) Y() float64      { return a.y }
