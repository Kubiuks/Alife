package model

import (
	"Alife/lib"
	"errors"
	"math"
	"math/rand"
	"sync"
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
	stressed			bool
	nutritionChange		float64
	socialChange		float64
	oxytocinChange		float64
	cortisolChange		float64
	adaptiveThreshold	float64
	bondPartners		[]int
	DSIstrengths		[]float64
	foodTimeWaiting		int
	motivation			float64
	physEffTouch		float64
	touchIntensity		float64
	tactileIntensity	float64
	DSImode				string
	psychEffEatTogether float64
	justEaten			bool
	sharedFoodWith		[]int
	eatingTogetherIntensity	float64

	stepSize     		float64

	// implementation needed
	mutex 		 sync.Mutex
	id 			 int
	x, y         float64
	origx, origy float64
	grid         *Grid
	trail        bool
	direction    float64
	ch 			 chan []float64
	numOfAgents	 int
}

func NewAgent(abm *lib.ABM, id, rank, numOfAgents int, x, y float64, ch chan []float64, trail bool, CortisolThresholdCondition, DSImode string) (*Agent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("agent needs a World defined to operate")
	}
	grid, ok := world.(*Grid)
	if !ok {
		return nil, errors.New("agent needs a Grid world to operate")
	}

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
		nutritionChange		: 0.0006,
		socialChange		: 0.0003,
		oxytocinChange		: 0.0005,
		cortisolChange		: 0.0005,
		rank				: rank,
		adaptiveThreshold	: adaptiveThreshold,
		foodTimeWaiting		: 0,
		motivation			: 0,
		physEffTouch		: 0.1,
		touchIntensity		: 0,
		tactileIntensity	: 0,
		DSImode				: DSImode,
		psychEffEatTogether : 0.1,
		justEaten			: false,
		eatingTogetherIntensity: 0.1,

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
		numOfAgents: numOfAgents,
	}, nil
}

func (a *Agent)  Run() {
	// send data
	data := []float64{float64(a.id), a.energy, a.socialness, a.oxytocin, a.cortisol}
	if a.ch != nil { a.ch <- data }

	// dont do anything if dead
	if !a.alive{ return }

	a.actionSelection()

	a.updateInternals()

	// check if died in this iteration
	if a.energy <= 0 {
		a.alive = false
		a.energy = 0
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
		if a.justEaten {
			a.move(mod(a.direction - 90, 360))
			a.justEaten = false
			a.foodTimeWaiting = 0
			a.sharedEatingFood()
			a.sharedFoodWith = nil
		} else {
			a.motivation = groomMotivation
			a.touchIntensity = a.motivation * a.physEffTouch
			a.tactileIntensity = a.touchIntensity * a.cortisol
			a.pickAgent(agents, walls)
		}
	} else {
		a.motivation = eatMotivation
		a.eatingTogetherIntensity = a.motivation * a.psychEffEatTogether
		a.findEatFood(agents, foods, walls)
	}

	a.updateCT(energyErr + socialErr, agents, foods)
}

func (a *Agent) agentVal(agent *Agent) float64{
	rankDiff := float64(a.rank - agent.Rank())/float64(a.numOfAgents-1)
	bond := 0.0
	DSI := 0.0
	for i, id := range a.bondPartners{
		if agent.ID() == id {
			// there is a bond
			bond = 1
			DSI = a.DSIstrengths[i]
			break
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
		var groomPartner *Agent
		agentVal := -1.0
		for _, temp := range agents {
			tmpAgentVal := a.agentVal(temp.(*Agent))
			if tmpAgentVal >= agentVal {
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
		a.IncreaseOT(a.touchIntensity)
		agent.IncreaseOT(a.touchIntensity)
		a.socialness = a.socialness + a.tactileIntensity * 0.25
		agent.ModulateCT(-1*a.tactileIntensity*0.3)
		if a.DSImode == "Variable" {
			a.ModulateDSI(agent.ID(), a.tactileIntensity*0.3)
			agent.ModulateDSI(a.id, a.tactileIntensity*0.3)
		}
		a.randomMove()
	} else {
		a.moveTo(agent)
	}
}

func (a *Agent) aggresion(agent *Agent){
	if distance(a.x, a.y, agent.X(), agent.Y()) < 2 {
		a.socialness = a.socialness + a.tactileIntensity * 0.25
		a.ModulateCT(-1*a.tactileIntensity*0.3)
		agent.ModulateCT(a.tactileIntensity*0.3)
		if a.DSImode == "Variable" {
			a.ModulateDSI(agent.ID(), -1*a.tactileIntensity*0.3)
			agent.ModulateDSI(a.id, -1*a.tactileIntensity*0.3)
		}
		a.randomMove()
	} else {
		a.moveTo(agent)
	}
}

func (a *Agent) avoidAgent() {
	a.move(mod(a.direction - (90*1.5*a.cortisol), 360))
}

func (a *Agent) findEatFood(agents, foods, walls []lib.Agent){
	if foods == nil {
		// dont see food
		// if see agents with AgentVal < 0 turn away
		flag := false
		for _, temp := range agents {
			if a.agentVal(temp.(*Agent)) < 0 {
				a.direction = mod(a.direction - (90*a.cortisol), 360)
				flag = true
				break
			}
		}
		// if dont see agents but see wall turn away else random move
		if !flag {
			if walls != nil {
				a.turnFromWall()
			} else {
				a.randomMove()
			}
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
			a.justEaten = true
			a.checkEatenWithBondPartner(food.(*Food))
			if a.energy >= 1 {
				a.move(mod(a.direction - 90, 360))
				a.energy = 1
				a.justEaten = false
				a.sharedEatingFood()
				a.sharedFoodWith = nil
			}
		} else {
			// see food, calculate if can approach
			a.approachOrAvoid(food.(*Food))
		}
	}
}

func (a *Agent) eatFood(f *Food) {
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
	a.stepSize = (1 + a.cortisol * 0.75)/2
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

func (a *Agent) checkEatenWithBondPartner(food *Food) {
	for _, agent := range food.EatingAgents() {
		for _, id := range a.bondPartners{
			if agent.ID() == id {
				// there is a bond with some other eating agent
				if !inList(agent.ID(), a.sharedFoodWith) {
					a.sharedFoodWith = append(a.sharedFoodWith, agent.ID())
				}
			}
		}
	}
}

func (a *Agent) sharedEatingFood() {
	if a.sharedFoodWith == nil {
		return
	}
	a.IncreaseOT(a.eatingTogetherIntensity)
	if a.DSImode == "Variable" {
		for _, id := range a.sharedFoodWith {
			a.ModulateDSI(id, a.eatingTogetherIntensity*a.cortisol*0.3)
		}
	}
}

func (a *Agent) updateInternals(){
	a.mutex.Lock()
	// lose energy
	a.energy = a.energy - 2 * (a.nutritionChange * a.stepSize)
	// lose and correct socialness
	if a.socialness > 1 {
		a.socialness = 1
	}
	a.socialness -= a.socialChange
	if a.socialness < 0 {
		a.socialness = 0
	}
	// lose and correct oxytocin
	if a.oxytocin > 1 {
		a.oxytocin = 1
	}
	a.oxytocin -= a.oxytocinChange
	if a.oxytocin < 0 {
		a.oxytocin = 0
	}
	// correct DSIstrengts
	for i:=0; i<len(a.DSIstrengths);i++ {
		if a.DSIstrengths[i] > 2 {
			a.DSIstrengths[i] = 2
		}
		if a.DSImode == "Variable" {
			a.DSIstrengths[i] = a.DSIstrengths[i] * 0.9997
		}
		if a.DSIstrengths[i] < 0 {
			a.DSIstrengths[i] = 0
		}
	}
	// correct cortisol
	if a.cortisol < 0 {
		a.cortisol = 0
	}
	if a.cortisol > 1 {
		a.cortisol = 1
	}
	// checked if stressed
	if a.cortisol > a.adaptiveThreshold {
		a.stressed = true
	} else {
		a.stressed = false
	}
	a.mutex.Unlock()
}

func (a *Agent) updateCT(sumOfErrors float64, agents, foods []lib.Agent){
	a.mutex.Lock()
	availableAgents := 0.0
	for _, tempAgent := range agents {
		tmpAgentVal := a.agentVal(tempAgent.(*Agent))
		availableAgents += 1 - tmpAgentVal
	}
	availableFoods := 0.0
	for _, tempFood := range foods {
		if tempFood.(*Food).Owner() == nil {
			availableFoods += 1
		} else if a.agentVal(tempFood.(*Food).Owner()) >= 0 {
			availableFoods += 1
		}
	}
	releaseRateCT := ((sumOfErrors - availableAgents - availableFoods)/2)*2*a.cortisolChange
	if releaseRateCT < 0 {
		releaseRateCT = releaseRateCT / 2
	}
	a.cortisol = a.cortisol + releaseRateCT
	if a.cortisol < 0 {
		a.cortisol = 0
	}
	if a.cortisol > 1 {
		a.cortisol = 1
	}
	if a.cortisol > a.adaptiveThreshold {
		a.stressed = true
	} else {
		a.stressed = false
	}
	a.mutex.Unlock()
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

func inList(element int, list []int) bool {

	for _, temp := range list {
		if temp == element {
			return true
		}
	}
	return false
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
		a.DSIstrengths = append(a.DSIstrengths, 2)
	}
}

func (a *Agent)	ModulateDSI(id int, amount float64){
	a.mutex.Lock()
	for i, partnerID := range a.bondPartners{
		if id == partnerID {
			a.DSIstrengths[i] = a.DSIstrengths[i] + amount
			if a.DSIstrengths[i] > 2 {
				a.DSIstrengths[i] = 2
			} else if a.DSIstrengths[i] < 0 {
				a.DSIstrengths[i] = 0
			}
			break
		}
	}
	a.mutex.Unlock()
}

func (a *Agent) ModulateCT(amount float64) {
	a.mutex.Lock()
	a.cortisol = a.cortisol + amount
	if a.cortisol < 0 {
		a.cortisol = 0
	}
	if a.cortisol > 1 {
		a.cortisol = 1
	}
	a.mutex.Unlock()
}
func (a *Agent) IncreaseOT(intensity float64){
	a.mutex.Lock()
	a.oxytocin = a.oxytocin + intensity
	if a.oxytocin > 1 {
		a.oxytocin = 1
	}
	a.mutex.Unlock()
}

func (a *Agent) Rank() int { return a.rank }
func (a *Agent) ID() int { return a.id }
func (a *Agent) Direction() float64 { return a.direction }
func (a *Agent) Alive() bool { return a.alive }
func (a *Agent) X() float64      { return a.x }
func (a *Agent) Y() float64      { return a.y }
