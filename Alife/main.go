package main

import (
	"Alife/lib"
	"Alife/model"
	"errors"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// setup
	a := lib.NewSimulation()
	w, h := 99, 99
	numberOfAgents := 6
	visionLength := 20

	// vision angle is both to the right and left so actually twice this amount
	// vision in degrees and must be smaller than 90 (overall smaller than 180)
	visionAngle := 40
	grid2D := model.NewWorld(w, h, numberOfAgents, visionLength, visionAngle)
	a.SetWorld(grid2D)

	// channels for communication with UI and Writer
	chGrid := make(chan [][]interface{})
	chVar := make(chan string, numberOfAgents)

	// initialise agents from 1 to 6
	for i:=1; i<numberOfAgents+1; i++ {
		x, y := randomFloat(float64(w)), randomFloat(float64(h))
		addAgent(x, y, i, i, a, grid2D, chVar, false, "Neutral", "Fixed")
	}

	// set up bonds between agents
	bondedAgents := []int{3,4,5}
	err := initialiseBonds(bondedAgents, numberOfAgents, a)
	if err != nil {
		log.Fatal(err)
	}

	// pick world settings
	setupWorld(a, grid2D, "Seasonal")

	a.LimitIterations(15000)

	// reporting function, does something each iteration
	// in this case updates the UI
	a.SetReportFunc(func(a *lib.ABM) {
		chGrid <- grid2D.Dump(func(a lib.Agent) int {
			//time.Sleep(100*time.Nanosecond)
			if a == nil {
				return 0
			}
			return a.ID()})
	})

	// get current time for datafile name
	//t := time.Now()

	//chan to make sure we finish writing before finishing
	finished := make(chan bool)

	// start writer go routine
	//writer := lib.NewWriter(finished, chVar, t.Format(time.Stamp)+".csv", numberOfAgents)
	writer := lib.NewWriter(finished, chVar, "test.csv", numberOfAgents)
	go writer.Loop()

	go func() {
		a.StartSimulation()
		chVar <- "end"
		close(chGrid)
	}()

	ui := lib.NewUI(w, h)
	defer ui.Stop()
	ui.AddGrid(chGrid)
	ui.Loop()

	//wait for writer to finish
	<- finished
	close(chVar)
	close(finished)
}
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________
//______________________________________________________________________________________________________________________

func addAgent(x, y float64, id, rank int, a *lib.ABM, grid2D *model.Grid, ch chan string,
				trail bool, CortisolThresholdCondition, DSImode string) {
	cell, err := model.NewAgent(a, id, rank, x, y, ch, trail, CortisolThresholdCondition, DSImode)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}

func addFood(x, y float64, a *lib.ABM, grid2D *model.Grid) {
	cell, err := model.NewFood(a, x, y)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}

func setupWorld(a *lib.ABM, grid2D *model.Grid, condition string) {
	err := grid2D.SetWorldDynamics(condition)
	if err != nil {
		log.Fatal(err)
	}
	switch condition{
	case "Four":
		addFood(9, 9, a, grid2D)
		addFood(89, 89, a, grid2D)
		addFood(9, 89, a, grid2D)
		addFood(89, 9, a, grid2D)
	case "Seasonal":
		addFood(9, 9, a, grid2D)
		addFood(89, 89, a, grid2D)
		addFood(9, 89, a, grid2D)
		addFood(89, 9, a, grid2D)
	case "Extreme":
		addFood(9, 9, a, grid2D)
		addFood(89, 89, a, grid2D)
		addFood(9, 89, a, grid2D)
		addFood(89, 9, a, grid2D)
	}
}

func initialiseBonds(bondedAgents []int, numberOfAgents int, a *lib.ABM) error {
	for i:= 0; i < len(bondedAgents); i++ {
		if bondedAgents[i] < 1 || bondedAgents[i] > numberOfAgents { return errors.New("invalid agent id") }
		for j:=0; j < len(bondedAgents); j++ {
			if i != j {
				if bondedAgents[i] == bondedAgents[j]{
					return errors.New("agent bond duplicate; agent cannot bond with itself")
				}
			}
		}
	}
	for i:= 0; i < len(bondedAgents); i++ {
		for _,agent := range a.Agents(){
			if agent.ID() == bondedAgents[i] {
				var bonds []int
				for j:=0; j < len(bondedAgents); j++ {
					if i != j {
						bonds = append(bonds, bondedAgents[j])
					}
				}
				agent.(*model.Agent).SetBonds(bonds)
			}
		}
	}
	return nil
}

// needed to make sure it's never 0
func randomFloat(max float64) float64 {
	var res float64
	for {
		res = rand.Float64()
		if res != 0 {
			break
		}
	}
	return res*max
}