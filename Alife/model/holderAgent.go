package model

import (
	"Alife/lib"
	"sync"
)

type HolderAgent struct {
	mx     sync.RWMutex
	id int
	x int
	y int
	grid *Grid
	agents []lib.Agent
}

func NewHolderAgent(grid *Grid, x, y int) *HolderAgent {
	return &HolderAgent{
		id:    -2,
		x:     x,
		y:     y,
		grid:  grid,
	}
}

func (h *HolderAgent) AddAgent(agent lib.Agent) {
	h.mx.Lock()
	h.agents = append(h.agents, agent)
	h.mx.Unlock()
}

func (h *HolderAgent) DeleteAgent(id int) (lib.Agent, lib.Agent){
	h.mx.Lock()
	defer h.mx.Unlock()
	l := len(h.agents)
	for i, e := range h.agents{
		if e.ID() == id {
			h.agents[i] = h.agents[l-1]
			if l == 2{
				return h.agents[0], e
			} else {
				h.agents = h.agents[:l-1]
				return h, e
			}
		}
	}
	return h, nil
}

func (h *HolderAgent) ID() int { return h.id }
func (h *HolderAgent) Agents() []lib.Agent {
	h.mx.RLock()
	defer h.mx.RUnlock()
	return h.agents
}

func (h *HolderAgent) Run(){
	return
}