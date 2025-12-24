package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

type Sector struct {
	Circle
	StartAngle float32 // 起始角度（度）
	EndAngle   float32 // 结束角度（度）
}

func NewSector(x, y, radius, startAngle, endAngle float32) *Sector {
	return &Sector{Circle: Circle{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x, Y: y}}, Radius: radius}, StartAngle: startAngle, EndAngle: endAngle}
}

func (s *Sector) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		return other.intersectsSector(s)
	case *Circle:
		return s.intersectsCircle(other)
	case *Sector:
		return s.intersectsSector(other)
	case *Rectangle:
		return other.Intersects(s)
	case *LineSegment:
		return other.Intersects(s)
	case *Triangle:
		return other.Intersects(s)
	}
	return false
}

func (s *Sector) intersectsCircle(circle *Circle) bool {
	return circle.intersectsSector(s)
}

func (s *Sector) intersectsSector(other *Sector) bool {
	// quick reject by circle distances
	dx := float64(s.Pos.X - other.Pos.X)
	dy := float64(s.Pos.Y - other.Pos.Y)
	d := math.Hypot(dx, dy)
	if d > float64(s.Radius+other.Radius) {
		return false
	}
	// center containment
	if pointInSectorGeneric(s.Pos, other) || pointInSectorGeneric(other.Pos, s) {
		return true
	}
	// check endpoints of s
	startRad := degToRad(s.StartAngle)
	endRad := degToRad(s.EndAngle)
	pStart := vec.Vec2[float32]{X: s.Pos.X + float32(math.Cos(startRad))*s.Radius, Y: s.Pos.Y + float32(math.Sin(startRad))*s.Radius}
	pEnd := vec.Vec2[float32]{X: s.Pos.X + float32(math.Cos(endRad))*s.Radius, Y: s.Pos.Y + float32(math.Sin(endRad))*s.Radius}
	if pointInSectorGeneric(pStart, other) || pointInSectorGeneric(pEnd, other) {
		return true
	}
	// check endpoints of other
	ostart := degToRad(other.StartAngle)
	oend := degToRad(other.EndAngle)
	oStart := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(ostart))*other.Radius, Y: other.Pos.Y + float32(math.Sin(ostart))*other.Radius}
	oEnd := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(oend))*other.Radius, Y: other.Pos.Y + float32(math.Sin(oend))*other.Radius}
	if pointInSectorGeneric(oStart, s) || pointInSectorGeneric(oEnd, s) {
		return true
	}
	// check circle-circle intersection points
	pts, ok := circleCircleIntersections(s.Pos, other.Pos, s.Radius, other.Radius)
	if ok {
		for _, p := range pts {
			if pointInSectorGeneric(p, s) && pointInSectorGeneric(p, other) {
				return true
			}
		}
	}
	// check radial edges vs other circle
	if segmentIntersectsCircle(s.Pos, pStart, other.Pos, other.Radius) || segmentIntersectsCircle(s.Pos, pEnd, other.Pos, other.Radius) || segmentIntersectsCircle(other.Pos, oStart, s.Pos, s.Radius) || segmentIntersectsCircle(other.Pos, oEnd, s.Pos, s.Radius) {
		return true
	}
	return false
}
