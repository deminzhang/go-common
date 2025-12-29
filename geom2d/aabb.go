package geom2d

import (
	"github.com/deminzhang/go-common/vec"
)

type AABB struct {
	BaseShape
	Width  float32
	Height float32
}

func NewAABB(x, y, width, height float32) *AABB {
	return &AABB{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x, Y: y}}, Width: width, Height: height}
}

func (a *AABB) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		// treat AABB as OBB with zero rotation
		return pointInRectangle(other.Pos, &OBB{AABB: *a})
	case *Circle:
		return circleRectIntersectSAT(other, a)
	case *LineSegment:
		return segmentRectIntersectSAT(other, a)
	case *Triangle:
		return triangleRectIntersectSAT(other, a)
	case *AABB:
		return rectRectIntersectSAT(&OBB{AABB: *a}, &OBB{AABB: *other})
	case *OBB:
		return rectRectIntersectSAT(&OBB{AABB: *a}, other)
	}
	return false
}

// helper wrappers to reuse OBB implementations
func circleRectIntersectSAT(c *Circle, a *AABB) bool {
	return (&OBB{AABB: *a}).Intersects(c)
}

func segmentRectIntersectSAT(ls *LineSegment, a *AABB) bool {
	return (&OBB{AABB: *a}).Intersects(ls)
}

func triangleRectIntersectSAT(t *Triangle, a *AABB) bool {
	return (&OBB{AABB: *a}).Intersects(t)
}
