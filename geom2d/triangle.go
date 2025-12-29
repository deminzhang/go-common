package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

type Triangle struct {
	A, B, C Point
}

func NewTriangle(ax, ay, bx, by, cx, cy float32) *Triangle {
	return &Triangle{A: Point{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: ax, Y: ay}}}, B: Point{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: bx, Y: by}}}, C: Point{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: cx, Y: cy}}}}
}

func (t *Triangle) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		return pointInTriangle(other.Pos, t.A.Pos, t.B.Pos, t.C.Pos)
	case *Circle:
		// circle center inside triangle
		if pointInTriangle(other.Pos, t.A.Pos, t.B.Pos, t.C.Pos) {
			return true
		}
		// any vertex inside circle
		if other.Pos.DistanceSqr(t.A.Pos) <= other.Radius*other.Radius || other.Pos.DistanceSqr(t.B.Pos) <= other.Radius*other.Radius || other.Pos.DistanceSqr(t.C.Pos) <= other.Radius*other.Radius {
			return true
		}
		// any edge intersects circle
		if segmentIntersectsCircle(t.A.Pos, t.B.Pos, other.Pos, other.Radius) || segmentIntersectsCircle(t.B.Pos, t.C.Pos, other.Pos, other.Radius) || segmentIntersectsCircle(t.C.Pos, t.A.Pos, other.Pos, other.Radius) {
			return true
		}
		return false
	case *OBB:
		// delegate to rectangle
		return other.Intersects(t)
	case *AABB:
		return other.Intersects(t)
	case *LineSegment:
		// endpoints in triangle
		if pointInTriangle(other.P1.Pos, t.A.Pos, t.B.Pos, t.C.Pos) || pointInTriangle(other.P2.Pos, t.A.Pos, t.B.Pos, t.C.Pos) {
			return true
		}
		// segment intersects triangle edges
		if segmentIntersectsSegment(other.P1.Pos, other.P2.Pos, t.A.Pos, t.B.Pos) || segmentIntersectsSegment(other.P1.Pos, other.P2.Pos, t.B.Pos, t.C.Pos) || segmentIntersectsSegment(other.P1.Pos, other.P2.Pos, t.C.Pos, t.A.Pos) {
			return true
		}
		return false
	case *Triangle:
		// vertex in other
		if pointInTriangle(t.A.Pos, other.A.Pos, other.B.Pos, other.C.Pos) || pointInTriangle(t.B.Pos, other.A.Pos, other.B.Pos, other.C.Pos) || pointInTriangle(t.C.Pos, other.A.Pos, other.B.Pos, other.C.Pos) {
			return true
		}
		if pointInTriangle(other.A.Pos, t.A.Pos, t.B.Pos, t.C.Pos) || pointInTriangle(other.B.Pos, t.A.Pos, t.B.Pos, t.C.Pos) || pointInTriangle(other.C.Pos, t.A.Pos, t.B.Pos, t.C.Pos) {
			return true
		}
		// edge intersections
		edgesA := [][2]vec.Vec2[float32]{{t.A.Pos, t.B.Pos}, {t.B.Pos, t.C.Pos}, {t.C.Pos, t.A.Pos}}
		edgesB := [][2]vec.Vec2[float32]{{other.A.Pos, other.B.Pos}, {other.B.Pos, other.C.Pos}, {other.C.Pos, other.A.Pos}}
		for _, ea := range edgesA {
			for _, eb := range edgesB {
				if segmentIntersectsSegment(ea[0], ea[1], eb[0], eb[1]) {
					return true
				}
			}
		}
		return false
	case *Sector:
		// any triangle vertex in sector
		if pointInSectorGeneric(t.A.Pos, other) || pointInSectorGeneric(t.B.Pos, other) || pointInSectorGeneric(t.C.Pos, other) {
			return true
		}
		// any sector endpoints inside triangle
		startRad := degToRad(other.StartAngle)
		endRad := degToRad(other.EndAngle)
		pStart := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(startRad))*other.Radius, Y: other.Pos.Y + float32(math.Sin(startRad))*other.Radius}
		pEnd := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(endRad))*other.Radius, Y: other.Pos.Y + float32(math.Sin(endRad))*other.Radius}
		if pointInTriangle(pStart, t.A.Pos, t.B.Pos, t.C.Pos) || pointInTriangle(pEnd, t.A.Pos, t.B.Pos, t.C.Pos) {
			return true
		}
		// edges intersection with radial segments or arc
		if segmentIntersectsSegment(t.A.Pos, t.B.Pos, other.Pos, pStart) || segmentIntersectsSegment(t.B.Pos, t.C.Pos, other.Pos, pStart) || segmentIntersectsSegment(t.C.Pos, t.A.Pos, other.Pos, pStart) {
			return true
		}
		if segmentIntersectsCircle(t.A.Pos, t.B.Pos, other.Pos, other.Radius) || segmentIntersectsCircle(t.B.Pos, t.C.Pos, other.Pos, other.Radius) || segmentIntersectsCircle(t.C.Pos, t.A.Pos, other.Pos, other.Radius) {
			// need to check angle of closest point
			if cp := closestPointOnSegment(other.Pos, t.A.Pos, t.B.Pos); angleBetweenDeg(float32(math.Atan2(float64(cp.Y-other.Pos.Y), float64(cp.X-other.Pos.X))*(180.0/math.Pi)), other.StartAngle, other.EndAngle) {
				return true
			}
			if cp := closestPointOnSegment(other.Pos, t.B.Pos, t.C.Pos); angleBetweenDeg(float32(math.Atan2(float64(cp.Y-other.Pos.Y), float64(cp.X-other.Pos.X))*(180.0/math.Pi)), other.StartAngle, other.EndAngle) {
				return true
			}
			if cp := closestPointOnSegment(other.Pos, t.C.Pos, t.A.Pos); angleBetweenDeg(float32(math.Atan2(float64(cp.Y-other.Pos.Y), float64(cp.X-other.Pos.X))*(180.0/math.Pi)), other.StartAngle, other.EndAngle) {
				return true
			}
		}
		return false
	}
	return false
}
