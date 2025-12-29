package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

type Point struct {
	BaseShape
}

func NewPoint(x, y float32) *Point {
	return &Point{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x, Y: y}}}
}

func (p *Point) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		return p.Pos.Equal(other.Pos)
	case *Circle:
		return p.intersectsCircle(other)
	case *Sector:
		return p.intersectsSector(other)
	case *LineSegment:
		return p.isOnLineSegment(other.P1.Pos, other.P2.Pos)
	case *OBB:
		return p.intersectsRectangle(other)
	case *AABB:
		// AABB is axis-aligned rectangle, reuse rectangle point test with zero-rotation OBB
		return pointInRectangle(p.Pos, &OBB{AABB: *other})
	case *Triangle:
		return p.intersectsTriangle(other)
	}
	return false
}

func (p *Point) intersectsCircle(circle *Circle) bool {
	disSq := p.Pos.DistanceSqr(circle.Pos)
	return disSq < circle.Radius*circle.Radius
}

func (p *Point) intersectsSector(sector *Sector) bool {
	disSq := p.Pos.DistanceSqr(sector.Pos)
	inCircle := disSq < sector.Radius*sector.Radius
	if !inCircle {
		return false
	}
	dir := vec.Vec2[float32]{
		X: p.Pos.X - sector.Pos.X,
		Y: p.Pos.Y - sector.Pos.Y,
	}
	angle := float32(math.Atan2(float64(dir.Y), float64(dir.X)) * (180.0 / math.Pi))
	if angle < 0 {
		angle += 360
	}
	startAngle := sector.StartAngle
	endAngle := sector.EndAngle
	if startAngle < 0 {
		startAngle += 360
	}
	if endAngle < 0 {
		endAngle += 360
	}
	if startAngle > endAngle {
		return angle >= startAngle || angle <= endAngle
	} else {
		return angle >= startAngle && angle <= endAngle
	}
}

func (p *Point) isOnLineSegment(start, end vec.Vec2[float32]) bool {
	cross := (end.Y-start.Y)*(p.Pos.X-start.X) - (end.X-start.X)*(p.Pos.Y-start.Y)
	if math.Abs(float64(cross)) > 1e-6 {
		return false
	}
	dot := (p.Pos.X-start.X)*(end.X-start.X) + (p.Pos.Y-start.Y)*(end.Y-start.Y)
	if dot < 0 {
		return false
	}
	squaredLengthBA := (end.X-start.X)*(end.X-start.X) + (end.Y-start.Y)*(end.Y-start.Y)
	if dot > squaredLengthBA {
		return false
	}
	return true
}

func (p *Point) intersectsRectangle(rect *OBB) bool {
	return pointInRectangle(p.Pos, rect)
}

func (p *Point) intersectsTriangle(tri *Triangle) bool {
	return pointInTriangle(p.Pos, tri.A.Pos, tri.B.Pos, tri.C.Pos)
}
