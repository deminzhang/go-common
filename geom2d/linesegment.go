package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

type LineSegment struct {
	P1, P2 Point
}

func NewLineSegment(x1, y1, x2, y2 float32) *LineSegment {
	return &LineSegment{P1: Point{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x1, Y: y1}}}, P2: Point{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x2, Y: y2}}}}
}

func (ls *LineSegment) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		return other.isOnLineSegment(ls.P1.Pos, ls.P2.Pos)
	case *Circle:
		return segmentIntersectsCircle(ls.P1.Pos, ls.P2.Pos, other.Pos, other.Radius)
	case *Sector:
		// endpoint inside sector
		if pointInSectorGeneric(ls.P1.Pos, other) || pointInSectorGeneric(ls.P2.Pos, other) {
			return true
		}
		// intersection with sector circle and angle check
		if segmentIntersectsCircle(ls.P1.Pos, ls.P2.Pos, other.Pos, other.Radius) {
			cp := closestPointOnSegment(other.Pos, ls.P1.Pos, ls.P2.Pos)
			angle := math.Atan2(float64(cp.Y-other.Pos.Y), float64(cp.X-other.Pos.X)) * (180.0 / math.Pi)
			if angle < 0 {
				angle += 360.0
			}
			if angleBetweenDeg(float32(angle), other.StartAngle, other.EndAngle) {
				return true
			}
		}
		// radial edges intersection
		startRad := degToRad(other.StartAngle)
		endRad := degToRad(other.EndAngle)
		pStart := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(startRad))*other.Radius, Y: other.Pos.Y + float32(math.Sin(startRad))*other.Radius}
		pEnd := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(endRad))*other.Radius, Y: other.Pos.Y + float32(math.Sin(endRad))*other.Radius}
		if segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, other.Pos, pStart) || segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, other.Pos, pEnd) {
			return true
		}
	case *Rectangle:
		if pointInRectangle(ls.P1.Pos, other) || pointInRectangle(ls.P2.Pos, other) {
			return true
		}
		rCorners := rectangleCorners(other)
		edges := [][2]vec.Vec2[float32]{{rCorners[0], rCorners[1]}, {rCorners[1], rCorners[2]}, {rCorners[2], rCorners[3]}, {rCorners[3], rCorners[0]}}
		for _, e := range edges {
			if segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, e[0], e[1]) {
				return true
			}
		}
	case *Triangle:
		if pointInTriangle(ls.P1.Pos, other.A.Pos, other.B.Pos, other.C.Pos) || pointInTriangle(ls.P2.Pos, other.A.Pos, other.B.Pos, other.C.Pos) {
			return true
		}
		if segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, other.A.Pos, other.B.Pos) || segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, other.B.Pos, other.C.Pos) || segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, other.C.Pos, other.A.Pos) {
			return true
		}
	case *LineSegment:
		return segmentIntersectsSegment(ls.P1.Pos, ls.P2.Pos, other.P1.Pos, other.P2.Pos)
	}
	return false
}

func (ls *LineSegment) Move(delta vec.Vec2[float32]) {
	ls.P1.Pos.Add(delta)
	ls.P2.Pos.Add(delta)
}
