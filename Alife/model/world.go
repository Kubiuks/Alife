package model

import (
	"errors"
	"sync"
	"Alife/lib"
)

type Grid struct {
	mx            sync.RWMutex
	width, height int

	cells, cellsPrev []lib.Agent
}

func NewWorld(width, height int) *Grid {
	g := &Grid{
		width:  width,
		height: height,
	}

	g.initSlices()

	return g
}

// Tick marks beginning of the new time period.
// Implements World interface.
func (g *Grid) Tick() {
	g.mx.RLock()
	defer g.mx.RUnlock()

	for i := 0; i < g.size(); i++ {
		g.cellsPrev[i] = lib.CopyAgent(g.cells[i])

	}

}

func (g *Grid) Move(id, fromX, fromY, toX, toY int) error {
	if err := g.validateXY(fromX, fromY); err != nil {
		return err
	}
	if err := g.validateXY(toX, toY); err != nil {
		return err
	}
	g.mx.Lock()
	defer g.mx.Unlock()

	agentFrom := g.cells[g.idx(fromX, fromY)]
	agentTo := g.cells[g.idx(toX, toY)]
	if agentFrom.ID() == id{
		g.cells[g.idx(fromX, fromY)] = nil
		if agentTo == nil {
			g.cells[g.idx(toX, toY)] = agentFrom
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agentFrom)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agentFrom)
			g.cells[g.idx(toX, toY)] = holder
		}
	} else {
		agentFromPost, agent := agentFrom.(*HolderAgent).DeleteAgent(id)
		g.cells[g.idx(fromX, fromY)] = agentFromPost
		if agentTo == nil {
			g.cells[g.idx(toX, toY)] = agent
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agent)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agent)
			g.cells[g.idx(toX, toY)] = holder
		}
	}
	return nil
}

func (g *Grid) Copy(id, fromX, fromY, toX, toY int) error {
	if err := g.validateXY(fromX, fromY); err != nil {
		return err
	}
	if err := g.validateXY(toX, toY); err != nil {
		return err
	}
	g.mx.Lock()
	defer g.mx.Unlock()

	agentFrom := g.cells[g.idx(fromX, fromY)]
	agentTo := g.cells[g.idx(toX, toY)]
	if agentFrom.ID() == id{
		if agentTo == nil {
			g.cells[g.idx(toX, toY)] = agentFrom
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agentFrom)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agentFrom)
			g.cells[g.idx(toX, toY)] = holder
		}
	} else {
		agentFromPost, agent := agentFrom.(*HolderAgent).DeleteAgent(id)
		g.cells[g.idx(fromX, fromY)] = agentFromPost
		if agentTo == nil {
			g.cells[g.idx(toX, toY)] = agent
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agent)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agent)
			g.cells[g.idx(toX, toY)] = holder
		}
	}
	return nil
}

func (g *Grid) Cell(x, y int) lib.Agent {
	if g.validateXY(x, y) != nil {
		return nil
	}
	g.mx.RLock()
	defer g.mx.RUnlock()
	return g.cells[g.idx(x, y)]
}

func (g *Grid) SetCell(x, y int, c lib.Agent) {
	if err := g.validateXY(x, y); err != nil {
		panic(err)
	}
	g.mx.Lock()
	temp := g.cells[g.idx(x, y)]
	if temp == nil {
		g.cells[g.idx(x, y)] = c
	} else if temp.ID() == -2{
		temp.(*HolderAgent).AddAgent(c)
	} else {
		holder := NewHolderAgent(g, x, y)
		holder.AddAgent(temp)
		holder.AddAgent(c)
		g.cells[g.idx(x, y)] = holder
	}
	g.mx.Unlock()
}

func (g *Grid) Width() int {
	return g.width
}

func (g *Grid) Height() int {
	return g.height
}

func (g *Grid) validateXY(x, y int) error {
	if x < 0 {
		return errors.New("x < 0")
	}
	if y < 0 {
		return errors.New("y < 0")
	}
	if x > g.width-1 {
		return errors.New("x > grid width")
	}
	if y > g.height-1 {
		return errors.New("y > grid height")
	}
	return nil
}

func (g *Grid) Dump(fn func(c lib.Agent) int) [][]interface{} {
	g.mx.RLock()
	defer g.mx.RUnlock()

	var ret = make([][]interface{}, g.width)
	for i := 0; i < g.width; i++ {
		ret[i] = make([]interface{}, g.height)
		for j := 0; j < g.height; j++ {
			a := g.cells[g.idx(i, j)]
			ret[i][j] = fn(a)
		}
	}
	return ret
}

// just move this verbose initialization here for brevity.
func (g *Grid) initSlices() {
	g.cells = make([]lib.Agent, g.size())
	g.cellsPrev = make([]lib.Agent, g.size())
}

func (g *Grid) size() int {
	return g.height * g.width
}

func (g *Grid) idx(x, y int) int {
	return y*g.width + x
}