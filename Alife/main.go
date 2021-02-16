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
	a := model.NewSimulation()
	w, h := 99, 99
	numberOfAgents := 6
	visionLength := 20
	// vision angle is both to the right and left so actually twice this amount
	// vision in degrees and must be smaller than 90 (overall smaller than 180)
	visionAngle := 40
	grid2D := model.NewWorld(w, h, numberOfAgents, visionLength, visionAngle)
	a.SetWorld(grid2D)
	//start from 1 coz id 0 is empty cell
	for i:=1; i<numberOfAgents+1; i++ {
		x, y := rand.Intn(w-1), rand.Intn(h-1)
		addAgent(x, y, i, a, grid2D, false)
	}

	addFood(5, 5, a, grid2D)
	//addFood(98, 98, a, grid2D)
	//addFood(0, 0, a, grid2D)
	//addFood(0, 98, a, grid2D)

	a.LimitIterations(10000000)

	chGrid := make(chan [][]interface{})
	//chAlive := make(chan int)
	a.SetReportFunc(func(a *model.ABM) {
		chGrid <- grid2D.Dump(func(a lib.Agent) int {
			//time.Sleep(10*time.Microsecond)
			if a == nil {
				return 0
			}
			return a.ID()})
	})

	go func() {
		a.StartSimulation()
		close(chGrid)
	}()

	ui := lib.NewUI(w, h)
	defer ui.Stop()
	ui.AddGrid(chGrid)
	ui.Loop()
}

func addAgent(x, y, id int, a *model.ABM, grid2D *model.Grid, trail bool) {
	cell, err := model.NewAgent(a, id, x, y, trail)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}

func addFood(x, y int, a *model.ABM, grid2D *model.Grid) {
	cell, err := model.NewFood(a, x, y)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}