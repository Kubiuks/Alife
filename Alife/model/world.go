package model

import (
	"Alife/lib"
	"errors"
	"fmt"
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
	visionVectors []directionVectors
	walls []directionVectors
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
	g.initialiseVisionVectors(visionLength, visionAngle)
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
	l := len(agents)
	for j := 0; j < l; j++ {
		if agent, ok := agents[j].(*Agent); ok {
			g.agentVision[agent.ID()-1] = nil
			center := vector{agent.x, agent.y}

			leftVisionEnd := vector{g.visionVectors[agent.direction].leftVector.x+center.x,
									g.visionVectors[agent.direction].leftVector.y+center.y}
			rightVisionEnd := vector{ g.visionVectors[agent.direction].rightVector.x+center.x,
									  g.visionVectors[agent.direction].rightVector.y+center.y}
			for i:=0;i<4;i++{
				if wall := g.checkWallInSigth(i, center, leftVisionEnd, rightVisionEnd); wall != nil{
					newWall := NewWall(wall.(vector).x,wall.(vector).y)
					g.agentVision[agent.ID()-1] = append(g.agentVision[agent.ID()-1], newWall)
				}
			}
			for k := 0; k < l; k++ {
				if agents[j] == agents[k] {
					continue
				}
				point := vector{agents[k].X(), agents[k].Y()}
				if isInsideSector(center, point, g.visionVectors[agent.direction].leftVector,
								  g.visionVectors[agent.direction].rightVector, g.visionLength){
					g.agentVision[agent.ID()-1] = append(g.agentVision[agent.ID()-1], agents[k])
				}
			}
		}
	}
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
		return errors.New("x < 0")
	}
	if y <= 0 {
		return errors.New("y < 0")
	}
	if x >= float64(g.width) {
		return errors.New("x > grid width")
	}
	if y >= float64(g.height) {
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

func (g *Grid) initialiseVisionVectors(visionLength, visionAngle int) {
	g.visionVectors = make([]directionVectors, 360)
	for i:=0;i<360;i++ {
		g.visionVectors[i].leftVector  = vector{float64(visionLength) * math.Sin((float64(i)+(float64(visionAngle)+0.00001))*(math.Pi/180.0)),
												float64(visionLength) * math.Cos((float64(i)+(float64(visionAngle)+0.00001))*(math.Pi/180.0))}
		g.visionVectors[i].rightVector = vector{float64(visionLength) * math.Sin((float64(i)-(float64(visionAngle)+0.00001))*(math.Pi/180.0)),
												float64(visionLength) * math.Cos((float64(i)-(float64(visionAngle)+0.00001))*(math.Pi/180.0))}
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (g *Grid) testVision() {
	fmt.Println("Testing Vision")
	center := vector{50, 50}
	point := vector{50, 60}
	for j := 0; j < 8; j++ {
		fmt.Println(isInsideSector(center, point, g.visionVectors[j].leftVector,
			g.visionVectors[j].rightVector, 20))
	}
	fmt.Println("Vision Tested")
}


func (g *Grid) testWalldetection() {
	fmt.Println("Testing WallDetection")
	center := vector{float64(90), float64(20)}
	leftVisionEnd := vector{g.visionVectors[1].leftVector.x+center.x,
		g.visionVectors[1].leftVector.y+center.y}
	rightVisionEnd := vector{ g.visionVectors[1].rightVector.x+center.x,
		g.visionVectors[1].rightVector.y+center.y}
	fmt.Println(center)
	fmt.Println(leftVisionEnd)
	fmt.Println(rightVisionEnd)
	wall := g.checkWallInSigth(2,center, leftVisionEnd, rightVisionEnd)
	newWall := NewWall(wall.(vector).x,wall.(vector).y)
	fmt.Println(wall)
	fmt.Println(newWall.x, newWall.y)
	fmt.Println("WallDetection Tested")
}