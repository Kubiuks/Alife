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
	// model needed
	oxytocin	 		float64
	cortisol	 		float64
	alive		 		bool
	energy		 		float64
	socialness	 		float64
	rank				int
	rankBasedCortisol	float64
	stressed			bool
	aggresionReceived	int
	aggresionGiven		int
	nutritionChange		float64
	socialChange		float64
	oxytocinChange		float64
	cortisolChange		float64
	adaptiveThreshold	float64
	bondPartners		[]int
	dyadicStrength		[]float64
	foodTimeWaiting		int

	stepSize     		float64
	// implementation needed
	id 			 int
	x, y         float64
	origx, origy float64
	grid         *Grid
	trail        bool
	direction    float64
	ch 			 chan string
}

func NewAgent(abm *lib.ABM, id, rank int, x, y float64, ch chan string, trail bool, CortisolThresholdCondition string) (*Agent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("agent needs a Grid world to operate")
	}
	//rankBasedCortisol, err := rankBasedOnCortisol(rank)
	//if err != nil {
	//	return nil, err
	//}
	adaptiveThreshold, err := cortisolThreshold(rank, CortisolThresholdCondition)
	if err != nil {
		return nil, err
	}
	return &Agent{
		alive    			: true,
		energy   			: 1,
		oxytocin 			: 1,
		cortisol 			: 0,
		socialness			: 1,
		stressed			: false,
		nutritionChange		: 0.0003,
		socialChange		: 0.0003,
		oxytocinChange		: 0.0005,
		cortisolChange		: 0.0005,
		aggresionReceived	: 0,
		aggresionGiven		: 0,
		rank				: rank,
		//rankBasedCortisol	: rankBasedCortisol,
		adaptiveThreshold	: adaptiveThreshold,
		foodTimeWaiting		: 0,


		stepSize 			: 1,
		//----------------
		id       : id,
		origx    : x,
		origy    : y,
		x        : x,
		y        : y,
		grid     : grid,
		trail    : trail,
		direction: rand.Float64()*360,
		ch       : ch,
	}, nil
}

func (a *Agent) Run() {
	// send data
	if a.ch != nil { a.ch <- fmt.Sprintf("'[%v, %v, %v, %v]'", a.id, a.energy, a.oxytocin, a.cortisol) }
	// dont do anything if dead
	if !a.alive{ return }

	a.actionSelection()

	a.updateInternals()

	// check if died in this iteration
	if a.energy <= 0 {
		a.alive = false
		a.grid.ClearCell(a.x, a.y, a.id)
	}
}

func (a *Agent) actionSelection(){
	foods, agents, walls := a.inVision()
	// food
	foodSalience := float64(len(foods))
	energyErr := 1 - a.energy
	eatMotivation := energyErr + (energyErr * foodSalience)
	// social
	socialErr := 1 - a.socialness
	robotSalience := float64(len(agents))
	groomMotivation := socialErr + (socialErr * robotSalience)


	if groomMotivation > eatMotivation {
		a.pickAgent(agents, walls)
	} else {
		a.findEatFood(foods, walls)
	}
}

func (a *Agent) agentVal(agent *Agent) float64{
	rankDiff := float64(a.rank - agent.Rank())/5
	bond := 0.0
	DSI := 0.0
	for i, id := range a.bondPartners{
		if agent.ID() == id {
			// there is a bond
			bond = 1
			DSI = a.dyadicStrength[i]
		}
	}
	return rankDiff + (bond*DSI*a.oxytocin)
}

func (a *Agent) pickAgent(agents, walls []lib.Agent) {
	if agents == nil {
		// dont see any agent, so turn if see wall else random move
		if walls != nil {
			a.turnFromWall()
		} else {
			a.randomMove()
		}
	} else {
		// see some agents, so need to pick groom partner
		// which is the agent with highest agentVal
		groomPartner := agents[0].(*Agent)
		agentVal := a.agentVal(groomPartner)
		for _, temp := range agents[1:] {
			tmpAgentVal := a.agentVal(temp.(*Agent))
			if tmpAgentVal > agentVal {
				agentVal = tmpAgentVal
				groomPartner = temp.(*Agent)
			}
		}
		a.groomOrAggresionOrAvoid(agentVal, groomPartner)
	}
}

func (a *Agent) groomOrAggresionOrAvoid(agentVal float64, agent *Agent) {
	if a.stressed {
		if agentVal < 0 {
			a.avoidAgent()
		} else if agentVal < 1 {
			if a.rank > agent.Rank() {
				a.aggresion(agent)
			} else {
				a.avoidAgent()
			}
		} else {
			a.groom(agent)
		}
	} else {
		a.groom(agent)
	}
}

func (a *Agent) groom(agent *Agent){
	if distance(a.x, a.y, agent.X(), agent.Y()) < 2 {
		// TODO perform groom
		a.socialness = 1
	} else {
		a.moveTo(agent)
	}
}

func (a *Agent) aggresion(agent *Agent){
	if distance(a.x, a.y, agent.X(), agent.Y()) < 2 {
		// TODO perform aggresion
		a.socialness = 1
	} else {
		a.moveTo(agent)
	}
}

func (a *Agent) avoidAgent() {
	a.move(mod(a.direction - (90*1.5*a.cortisol), 360))
}

func (a *Agent) findEatFood(foods, walls []lib.Agent){
	if foods == nil {
		// dont see food, so turn if see wall else random move
		if walls != nil {
			a.turnFromWall()
		} else {
			a.randomMove()
		}
	} else {
		// see food so approach or eat if close
		// first calculate closest food
		var food lib.Agent
		dist := 100.0
		for _, temp := range foods{
			tmpDist := distance(a.x, a.y, temp.X(), temp.Y())
			if tmpDist < dist {
				food = temp
				dist = tmpDist
			}
		}
		if dist <= 1 {
			// next to food, so can eat
			a.eatFood(food.(*Food))
			if a.energy >= 1{
				a.move(mod(a.direction - 90, 360))
			}
		} else {
			// see food, calculate if can approach
			a.approachOrAvoid(food.(*Food))
		}
	}
}

func (a *Agent) eatFood(f *Food){
	if a.foodTimeWaiting < 5 {
		a.foodTimeWaiting++
	} else {
		f.reduceResource(0.01)
		a.energy += 0.01
		a.foodTimeWaiting++
	}
	if a.foodTimeWaiting >= 6 {
		a.foodTimeWaiting = 0
	}
}

func (a *Agent) turnFromWall(){
	a.direction = mod(a.direction + rand.Float64()*135 - rand.Float64()*135, 360)
}

func (a *Agent) approachOrAvoid(f *Food) {
	agentVal := 1.0
	if f.Owner() != nil {
		agentVal = a.agentVal(f.Owner())
	}
	if agentVal < 0 {
		// cant approach food
		a.move(mod(a.direction - 180, 360))
	} else {
		a.moveTo(f)
	}
}

func (a *Agent) moveTo(agent lib.Agent){
	a.move(math.Atan2(agent.X() - a.x, agent.Y() - a.y)*(180.0/math.Pi))
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
	}
}

func (a *Agent) updateInternals(){
	a.energy -= a.nutritionChange
	if a.oxytocin > 0 {
		a.socialness -= a.socialChange
	}
}

func (a *Agent) inVision() ([]lib.Agent,[]lib.Agent,[]lib.Agent){
	var foods, agents, walls []lib.Agent
	for _,agent := range a.grid.agentVision[a.id-1] {
		if agent.ID() == -1 {
			// sees food
			foods = append(foods, agent)
		} else if agent.ID() == -3 {
			// sees a wall
			walls = append(walls, agent)
		} else {
			// the agent sees another agent
			agents = append(agents, agent)
		}
	}
	return foods, agents, walls
}

/*
________________________________________________________________________________________________________________________
___________________________________________SETUP/GETTERS/SETTERS________________________________________________________
________________________________________________________________________________________________________________________
*/

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

func cortisolThreshold(rank int, CortisolThresholdCondition string) (float64, error) {
	switch CortisolThresholdCondition {
	case "Control":
		return 1.1, nil
	case "Neutral":
		return 0.5, nil
	case "High":
		return 0.7, nil
	case "Low":
		return 0.2, nil
	case "Low-High":
		switch rank {
		case 1:
			return 0.7, nil
		case 2:
			return 0.6, nil
		case 3:
			return 0.5, nil
		case 4:
			return 0.4, nil
		case 5:
			return 0.3, nil
		case 6:
			return 0.2, nil
		}
	case "High-Low":
		switch rank {
		case 1:
			return 0.2, nil
		case 2:
			return 0.3, nil
		case 3:
			return 0.4, nil
		case 4:
			return 0.5, nil
		case 5:
			return 0.6, nil
		case 6:
			return 0.7, nil
		}
	}
	return 0, errors.New("invalid cortisol threshhold condition")
}

func (a *Agent) SetBonds(bonds []int) {
	a.bondPartners = bonds
	for i:=0; i<len(bonds);i++ {
		a.dyadicStrength = append(a.dyadicStrength, 2)
	}
}
func (a *Agent) Rank() int { return a.rank }
func (a *Agent) ID() int { return a.id }
func (a *Agent) Direction() float64 { return a.direction }
func (a *Agent) Alive() bool { return a.alive }
func (a *Agent) X() float64      { return a.x }
func (a *Agent) Y() float64      { return a.y }
