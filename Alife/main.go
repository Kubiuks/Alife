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
	grid2D := model.NewWorld(w, h)
	a.SetWorld(grid2D)
	//start from 1 coz id 0 is empty cell
	for i:=1; i<7; i++ {
		x, y := rand.Intn(w-1), rand.Intn(h-1)
		addAgent(x, y, i, a, grid2D, false)
	}

	addFood(5, 5, a, grid2D)

	a.LimitIterations(10000000)

	chGrid := make(chan [][]interface{})
	//chAlive := make(chan int)
	a.SetReportFunc(func(a *lib.ABM) {
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

	ui := model.NewUI(w, h)
	defer ui.Stop()
	ui.AddGrid(chGrid)
	ui.Loop()
}

func addAgent(x, y, id int, a *lib.ABM, grid2D *model.Grid, trail bool) {
	cell, err := model.NewAgent(a, id, x, y, trail)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}

func addFood(x, y int, a *lib.ABM, grid2D *model.Grid) {
	cell, err := model.NewFood(a, x, y)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}