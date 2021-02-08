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
	for i:=0; i<5; i++ {
		cell, err := model.NewAgent(a, i, rand.Intn(w-1), rand.Intn(h-1), true)
		if err != nil {
			log.Fatal(err)
		}
		a.AddAgent(cell)
		grid2D.SetCell(cell.X(), cell.Y(), cell)
	}

	a.LimitIterations(9999)

	ch := make(chan [][]interface{})
	a.SetReportFunc(func(a *lib.ABM) {
		ch <- grid2D.Dump(func(a lib.Agent) int {
			if a == nil {
				return 0
			}
			return 1})
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
