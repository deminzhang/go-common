package geom2d

import "github.com/deminzhang/go-common/vec"

type IShape interface {
	Intersects(other IShape) bool
}

type BaseShape struct {
	Pos vec.Vec2[float32]
}

func (bs *BaseShape) Intersects(other IShape) bool {
	return false
}

func (bs *BaseShape) Move(delta vec.Vec2[float32]) {
	bs.Pos.Add(delta)
}
