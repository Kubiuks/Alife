package main

import (
	"Alife/lib"
	"Alife/model"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
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

	//start from 1 coz id 0 is empty cell
	for i:=1; i<numberOfAgents+1; i++ {
		x, y := randomFloat(float64(w)), randomFloat(float64(h))
		addAgent(x, y, i, a, grid2D, chVar, false)
	}

	addFood(5, 5, a, grid2D)
	//addFood(98, 98, a, grid2D)
	//addFood(0, 0, a, grid2D)
	//addFood(0, 98, a, grid2D)

	a.LimitIterations(10000000)

	a.SetReportFunc(func(a *lib.ABM) {
		chGrid <- grid2D.Dump(func(a lib.Agent) int {
			time.Sleep(1*time.Microsecond)
			if a == nil {
				return 0
			}
			return a.ID()})
	})

	go func() {
		a.StartSimulation()
		chVar <- "end"
		close(chGrid)
	}()

	// get current time for datafile name
	//t := time.Now()

	//chan to make sure we finish writing before finishing
	finished := make(chan bool)

	// start writer go routine
	//writer := lib.NewWriter(finished, chVar, t.Format(time.Stamp)+".csv", numberOfAgents)
	writer := lib.NewWriter(finished, chVar, "test.csv", numberOfAgents)
	go writer.Loop()

	ui := lib.NewUI(w, h)
	defer ui.Stop()
	ui.AddGrid(chGrid)
	ui.Loop()

	//wait for writer to finish
	<- finished
	close(chVar)
	close(finished)
}

func addAgent(x, y float64, id int, a *lib.ABM, grid2D *model.Grid, ch chan string, trail bool) {
	cell, err := model.NewAgent(a, id, x, y, ch, trail)
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