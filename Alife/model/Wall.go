package model

type Wall struct {
	id 			 int
	x, y         int
}

func NewWall(x, y int) *Wall {
	return &Wall{
		id:    -3,
		x:     x,
		y:     y,
	}
}


func (w *Wall) Run() { return }
func (w *Wall) ID() int { return w.id }
func (w *Wall) X() int { return w.x }
func (w *Wall) Y() int { return w.y }