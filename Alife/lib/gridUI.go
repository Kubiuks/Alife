package lib

import (
	"image"
	"image/color"
	"image/draw"
)

type GridWidget struct {
	OffsetX, OffsetY int
	SquareSize       int
	Cols, Rows       int

	image *image.RGBA
	grid  [][]int
}

func NewGridWidget(cols, rows, w, h int) *GridWidget {
	squareSize := w / cols
	grid := make([][]int, cols)
	for i := 0; i < cols; i++ {
		grid[i] = make([]int, rows)
	}
	return &GridWidget{
		Cols:       cols,
		Rows:       rows,
		SquareSize: squareSize,

		grid:  grid,
		image: image.NewRGBA(image.Rect(0, 0, cols*squareSize+20, rows*squareSize+20)),
	}
}

func (g *GridWidget) Draw(m *image.RGBA) {
	r := g.image.Bounds()
	draw.Draw(m, r, g.image, image.ZP, draw.Src)
	gridColor := image.NewUniform(color.RGBA{90, 90, 90, 0})
	// Vertical lines.
	x := 10
	y := 10
	wid := 1
	for i := 0; i < g.Cols; i++ {
		r := image.Rect(x, y, x+wid, y+(g.Rows-1)*g.SquareSize)
		draw.Draw(m, r, gridColor, image.ZP, draw.Src)
		x += g.SquareSize
	}
	// Horizontal lines.
	x = 10
	for i := 0; i < g.Rows; i++ {
		r := image.Rect(x, y, x+(g.Cols-1)*g.SquareSize+wid, y+wid)
		draw.Draw(m, r, gridColor, image.ZP, draw.Src)
		y += g.SquareSize
	}

	for i := 0; i < g.Cols; i++ {
		for j := 0; j < g.Rows; j++ {
			g.DrawCell(m, i, j, g.grid[i][j])
		}
	}
}

func (g *GridWidget) SetGrid(dump [][]interface{}) {
	for i := 0; i < g.Cols-1; i++ {
		for j := 0; j < g.Rows-1; j++ {
			g.grid[i][j] = dump[i][j].(int)
		}
	}
}

func (g *GridWidget) DrawCell(m *image.RGBA, x, y int, agent int) {
	offX, offY := 10, 10
	X := offX + x*g.SquareSize
	Y := offY + y*g.SquareSize
	colors := []*image.Uniform{image.NewUniform(color.RGBA{255, 0, 0, 0}),
							   image.NewUniform(color.RGBA{0, 250, 0, 0}),
							   image.NewUniform(color.RGBA{0, 0, 255, 0}),
							   image.NewUniform(color.RGBA{255, 0, 255, 0}),
							   image.NewUniform(color.RGBA{0, 255, 255, 0}),
							   image.NewUniform(color.RGBA{255, 140, 0, 0})}
	col := image.Black
	// agents start from 1
	if agent != 0{
		if agent == -1{
			col = image.NewUniform(color.RGBA{255, 255, 0, 0})
		} else {
			agent--
			col = colors[agent%6]
		}
	}
	r := image.Rect(X+1, Y+1, X+g.SquareSize-1, Y+g.SquareSize-1)
	draw.Draw(m, r, col, image.ZP, draw.Src)
}
