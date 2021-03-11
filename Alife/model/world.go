package model

import (
	"Alife/lib"
	"errors"
	"math"
	"sync"
)

type Grid struct {
	mx            sync.RWMutex
	width, height int
	visionLength int
	visionAngle  int
	cells []lib.Agent
	agentVision [][]lib.Agent
	walls []directionVectors
	worldDynamics string
}

type directionVectors struct {
	leftVector vector
	rightVector vector
}

type vector struct {
	x, y float64
}

func NewWorld(width, height, numberOfAgents, visionLength, visionAngle int) *Grid {
	g := &Grid{
		width:  width,
		height: height,
		visionLength: visionLength,
		visionAngle: visionAngle,
	}
	g.cells = make([]lib.Agent, g.size())
	g.agentVision = make([][]lib.Agent, numberOfAgents)
	for i := 0; i < numberOfAgents; i++ {
		g.agentVision[i] = nil
	}
	g.walls = make([]directionVectors, 4)
	g.initialiseWalls(width, height)
	//g.testVision()
	//g.testIntersection()
	//g.testWalldetection()
	return g
}

// Tick marks beginning of the new time period.
// Implements World interface.
func (g *Grid) Tick(agents []lib.Agent) {
	g.mx.RLock()
	defer g.mx.RUnlock()
	for j := 0; j < len(agents); j++ {
		if agent, ok := agents[j].(*Agent); ok {
			g.checkAgentVision(agents, agent)
		} else if food, ok := agents[j].(*Food); ok {
			g.checkOwnerOfFood(agents, food)
		}
	}
}

func (g *Grid) checkAgentVision(agents []lib.Agent, agent *Agent) {
	g.agentVision[agent.id-1] = nil
	center := vector{agent.x, agent.y}
	vision := g.findVsionVectors(agent.direction, g.visionLength, g.visionAngle)
	leftVisionEnd := vector{vision.leftVector.x+center.x,
		vision.leftVector.y+center.y}
	rightVisionEnd := vector{ vision.rightVector.x+center.x,
		vision.rightVector.y+center.y}
	for i:=0;i<4;i++{
		if wall := g.checkWallInSigth(i, center, leftVisionEnd, rightVisionEnd); wall != nil{
			newWall := NewWall(wall.(vector).x,wall.(vector).y)
			g.agentVision[agent.ID()-1] = append(g.agentVision[agent.ID()-1], newWall)
		}
	}
	for k := 0; k < len(agents); k++ {
		if agent.ID() == agents[k].ID() {
			continue
		}
		point := vector{agents[k].X(), agents[k].Y()}
		if isInsideSector(center, point, vision.leftVector,
			vision.rightVector, g.visionLength){
			g.agentVision[agent.ID()-1] = append(g.agentVision[agent.ID()-1], agents[k])
		}
	}
}

func (g *Grid) checkOwnerOfFood(agents []lib.Agent, food *Food) {
	var highestRankAgent *Agent
	highestRank := 0
	for k := 0; k < len(agents); k++ {
		if agents[k].ID() < 1 {
			continue
		}
		center := vector{food.X(), food.Y()}
		point := vector{agents[k].X(), agents[k].Y()}
		relVector := vector{point.x - center.x, point.y - center.y}
		if isWithinRadius(relVector, 4) && agents[k].(*Agent).Rank() > highestRank {
			highestRankAgent = agents[k].(*Agent)
		}
	}
	food.SetOwner(highestRankAgent)
}

func (g *Grid) Move(id int, fromX, fromY, toX, toY float64) error {
	if err := g.validateXY(fromX, fromY); err != nil {
		return err
	}
	if err := g.validateXY(toX, toY); err != nil {
		return err
	}
	g.mx.Lock()
	defer g.mx.Unlock()
	indexFrom := g.idx(fromX, fromY)
	indexTo := g.idx(toX, toY)
	if indexFrom == indexTo {
		return nil
	}
	agentFrom := g.cells[indexFrom]
	agentTo := g.cells[indexTo]
	if agentFrom.ID() == id{
		g.cells[indexFrom] = nil
		if agentTo == nil {
			g.cells[indexTo] = agentFrom
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agentFrom)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agentFrom)
			g.cells[indexTo] = holder
		}
	} else {
		agentFromPost, agent := agentFrom.(*HolderAgent).DeleteAgent(id)
		g.cells[indexFrom] = agentFromPost
		if agentTo == nil {
			g.cells[indexTo] = agent
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agent)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agent)
			g.cells[indexTo] = holder
		}
	}
	return nil
}

func (g *Grid) Copy(id int, fromX, fromY, toX, toY float64) error {
	if err := g.validateXY(fromX, fromY); err != nil {
		return err
	}
	if err := g.validateXY(toX, toY); err != nil {
		return err
	}
	g.mx.Lock()
	defer g.mx.Unlock()
	indexFrom := g.idx(fromX, fromY)
	indexTo := g.idx(toX, toY)
	agentFrom := g.cells[indexFrom]
	agentTo := g.cells[indexTo]
	if agentFrom.ID() == id{
		if agentTo == nil {
			g.cells[indexTo] = agentFrom
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agentFrom)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agentFrom)
			g.cells[indexTo] = holder
		}
	} else {
		agentFromPost, agent := agentFrom.(*HolderAgent).DeleteAgent(id)
		g.cells[indexFrom] = agentFromPost
		if agentTo == nil {
			g.cells[indexTo] = agent
		} else if agentTo.ID() == -2 {
			agentTo.(*HolderAgent).AddAgent(agent)
		} else {
			holder := NewHolderAgent(g, toX, toY)
			holder.AddAgent(agentTo)
			holder.AddAgent(agent)
			g.cells[indexTo] = holder
		}
	}
	return nil
}

func (g *Grid) SetCell(x, y float64, c lib.Agent) {
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

func (g *Grid) ClearCell(x, y float64, id int) {
	if err := g.validateXY(x, y); err != nil {
		panic(err)
	}
	g.mx.Lock()
	temp := g.cells[g.idx(x, y)]
	if temp.ID() == id {
		g.cells[g.idx(x, y)] = nil
	} else {
		temp.(*HolderAgent).DeleteAgent(id)
	}
	g.mx.Unlock()
}

func (g *Grid) size() int {
	return g.height * g.width
}

func (g *Grid) idx(x, y float64) int {
	return int(math.Floor(y))*g.width + int(math.Floor(x))
}

func (g *Grid) Width() int {
	return g.width
}

func (g *Grid) Height() int {
	return g.height
}

func (g *Grid) validateXY(x, y float64) error {
	if x <= 0 {
		return errors.New("x <= 0")
	}
	if y <= 0 {
		return errors.New("y <= 0")
	}
	if x >= float64(g.width) {
		return errors.New("x >= grid width")
	}
	if y >= float64(g.height) {
		return errors.New("y >= grid height")
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
			a := g.cells[g.idx(float64(i), float64(j))]
			ret[i][j] = fn(a)
		}
	}
	return ret
}

func isInsideSector(center, point, sectorLeft, sectorRight vector, radius int) bool {
	relVector := vector{point.x - center.x, point.y - center.y}
	return isWithinRadius(relVector, radius) &&
		   !areClockwise(sectorRight, relVector) &&
		   areClockwise(sectorLeft, relVector)
}

func areClockwise(v1, v2 vector) bool {
	return -v1.y*v2.x + v1.x*v2.y > 0
}
func isWithinRadius(v vector, radius int) bool {
	return v.x*v.x + v.y*v.y <= math.Pow(float64(radius), 2)
}

func (g *Grid) checkWallInSigth(wallId int, center, leftVisionEnd, rightVisionEnd vector) interface{}{
	wallStart := g.walls[wallId].leftVector
	wallEnd := g.walls[wallId].rightVector
	leftIntersection := findIntersection(center, leftVisionEnd, wallStart, wallEnd)
	rightIntersection := findIntersection(center, rightVisionEnd, wallStart, wallEnd)
	if leftIntersection != nil && rightIntersection != nil {
		return pointOnWallWithlowestDistance(center, wallStart, wallEnd, g.visionLength)
	} else if leftIntersection != nil {
		return pointOnWallWithlowestDistance(center, wallEnd, leftIntersection.(vector), g.visionLength)
	} else if rightIntersection != nil {
		return pointOnWallWithlowestDistance(center, wallStart, rightIntersection.(vector), g.visionLength)
	}
	return nil
}

func findIntersection(p0, p1, p2, p3 vector) interface{}{
	s10X := p1.x - p0.x
	s10Y := p1.y - p0.y
	s32X := p3.x - p2.x
	s32Y := p3.y - p2.y
	denom := s10X*s32Y - s32X*s10Y
	denomIsPositive := denom > 0
	s02X := p0.x - p2.x
	s02Y := p0.y - p2.y
	sNumer := s10X*s02Y - s10Y*s02X
	if (sNumer < 0) == denomIsPositive {
		return nil
	}
	tNumer := s32X * s02Y - s32Y * s02X
	if (tNumer < 0) == denomIsPositive {
		return nil
	}
	if (sNumer > denom) == denomIsPositive || (tNumer > denom) == denomIsPositive {
		return nil
	}
	t := tNumer / denom
	intersectionPoint := vector{p0.x + (t * s10X), p0.y + (t * s10Y)}
	return intersectionPoint
}

func pointOnWallWithlowestDistance(point, wallStart, wallEnd vector, visionLength int) interface{}{
	A := point.x - wallStart.x
	B := point.y - wallStart.y
	C := wallEnd.x - wallStart.x
	D := wallEnd.y - wallStart.y
	dotProduct := A*C + B*D
	lenSq := C * C + D * D
	var xx, yy float64
	param := dotProduct / lenSq
	if param < 0 {
		xx = wallStart.x
		yy = wallStart.y
	} else if param > 1 {
		xx = wallEnd.x
		yy = wallEnd.y
	} else {
		xx = wallStart.x + param * C
		yy = wallStart.y + param * D
	}
	dx := point.x - xx
	dy := point.y - yy
	if math.Sqrt(dx * dx + dy * dy) <= float64(visionLength)/2 {
		return vector{xx,yy}
	}
	return nil
}

/*
	___1___
	|     |
	0     2
	|__3__|
 */
func (g *Grid) initialiseWalls(width, height int) {
	g.walls[0].leftVector  = vector{0,float64(height)}
	g.walls[0].rightVector = vector{0, 0}
	g.walls[1].leftVector  = vector{0,0}
	g.walls[1].rightVector = vector{float64(width), 0}
	g.walls[2].leftVector  = vector{float64(width),0}
	g.walls[2].rightVector = vector{float64(width), float64(height)}
	g.walls[3].leftVector  = vector{float64(width),float64(height)}
	g.walls[3].rightVector = vector{0, float64(height)}
}

func (g *Grid) findVsionVectors(direction float64, visionLength, visionAngle int) directionVectors {
	return directionVectors{leftVector : vector{float64(visionLength) * math.Sin((direction+(float64(visionAngle)+0.00001))*(math.Pi/180.0)),
											    float64(visionLength) * math.Cos((direction+(float64(visionAngle)+0.00001))*(math.Pi/180.0))},
							rightVector: vector{float64(visionLength) * math.Sin((direction-(float64(visionAngle)+0.00001))*(math.Pi/180.0)),
												float64(visionLength) * math.Cos((direction-(float64(visionAngle)+0.00001))*(math.Pi/180.0))}}
}

func (g *Grid) SetWorldDynamics(condition string){
	g.worldDynamics = condition
}
