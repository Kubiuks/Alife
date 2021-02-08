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
	//start from 1 coz 0 is empty cell
	for i:=1; i<7; i++ {
		cell, err := model.NewAgent(a, i, rand.Intn(w-1), rand.Intn(h-1), true)
		if err != nil {
			log.Fatal(err)
		}
		a.AddAgent(cell)
		grid2D.SetCell(cell.X(), cell.Y(), cell)
	}

	addFood(5, 5, a, grid2D)

	a.LimitIterations(999)

	ch := make(chan [][]interface{})
	a.SetReportFunc(func(a *lib.ABM) {
		ch <- grid2D.Dump(func(a lib.Agent) int {
			if a == nil {
				return 0
			}
			return a.ID()})
	})

	go func() {
		a.StartSimulation()
		close(ch)
	}()

	ui := model.NewUI(w, h)
	defer ui.Stop()
	ui.AddGrid(ch)
	ui.Loop()
}

func addFood(x, y int, a *lib.ABM, grid2D *model.Grid) {
	cell, err := model.NewFood(a, x, y)
	if err != nil {
		log.Fatal(err)
	}
	a.AddAgent(cell)
	grid2D.SetCell(cell.X(), cell.Y(), cell)
}