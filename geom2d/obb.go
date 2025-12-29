package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

type OBB struct {
	AABB
	Angle float32 // 旋转角度（度）
}

func NewOBB(x, y, width, height, angle float32) *OBB {
	return &OBB{AABB: AABB{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x, Y: y}}, Width: width, Height: height}, Angle: angle}
}

func (r *OBB) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		return pointInRectangle(other.Pos, r)
	case *Circle:
		// transform circle center into rectangle local coords
		local := rotatePointAround(other.Pos, r.Pos, -r.Angle)
		dx := float64(local.X - r.Pos.X)
		dy := float64(local.Y - r.Pos.Y)
		hw := float64(r.Width / 2.0)
		hh := float64(r.Height / 2.0)
		closestX := clamp(dx, -hw, hw)
		closestY := clamp(dy, -hh, hh)
		dx2 := dx - closestX
		dy2 := dy - closestY
		return dx2*dx2+dy2*dy2 <= float64(other.Radius*other.Radius)+1e-6
	case *OBB:
		return rectRectIntersectSAT(r, other)
	case *AABB:
		// treat AABB as OBB with zero rotation
		return rectRectIntersectSAT(r, &OBB{AABB: *other})
	case *LineSegment:
		// check endpoints
		if pointInRectangle(other.P1.Pos, r) || pointInRectangle(other.P2.Pos, r) {
			return true
		}
		// check edges
		rCorners := rectangleCorners(r)
		edges := [][2]vec.Vec2[float32]{{rCorners[0], rCorners[1]}, {rCorners[1], rCorners[2]}, {rCorners[2], rCorners[3]}, {rCorners[3], rCorners[0]}}
		for _, e := range edges {
			if segmentIntersectsSegment(other.P1.Pos, other.P2.Pos, e[0], e[1]) {
				return true
			}
		}
		return false
	case *Triangle:
		// check triangle vertices
		if pointInRectangle(other.A.Pos, r) || pointInRectangle(other.B.Pos, r) || pointInRectangle(other.C.Pos, r) {
			return true
		}
		// rect corners in triangle
		rCorners2 := rectangleCorners(r)
		if pointInTriangle(rCorners2[0], other.A.Pos, other.B.Pos, other.C.Pos) || pointInTriangle(rCorners2[1], other.A.Pos, other.B.Pos, other.C.Pos) || pointInTriangle(rCorners2[2], other.A.Pos, other.B.Pos, other.C.Pos) || pointInTriangle(rCorners2[3], other.A.Pos, other.B.Pos, other.C.Pos) {
			return true
		}
		// edge intersections
		edges2 := [][2]vec.Vec2[float32]{{rCorners2[0], rCorners2[1]}, {rCorners2[1], rCorners2[2]}, {rCorners2[2], rCorners2[3]}, {rCorners2[3], rCorners2[0]}}
		triEdges := [][2]vec.Vec2[float32]{{other.A.Pos, other.B.Pos}, {other.B.Pos, other.C.Pos}, {other.C.Pos, other.A.Pos}}
		for _, re := range edges2 {
			for _, te := range triEdges {
				if segmentIntersectsSegment(re[0], re[1], te[0], te[1]) {
					return true
				}
			}
		}
		return false
	case *Sector:
		// if sector center in rect
		if pointInRectangle(other.Pos, r) {
			return true
		}
		// if any rect corner in sector
		rC := rectangleCorners(r)
		for _, c := range rC {
			if pointInSectorGeneric(c, other) {
				return true
			}
		}
		// if sector arc points in rect
		startRad := degToRad(other.StartAngle)
		endRad := degToRad(other.EndAngle)
		pStart := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(startRad))*other.Radius, Y: other.Pos.Y + float32(math.Sin(startRad))*other.Radius}
		pEnd := vec.Vec2[float32]{X: other.Pos.X + float32(math.Cos(endRad))*other.Radius, Y: other.Pos.Y + float32(math.Sin(endRad))*other.Radius}
		if pointInRectangle(pStart, r) || pointInRectangle(pEnd, r) {
			return true
		}
		// check edges vs sector radial segments and arc
		edges3 := [][2]vec.Vec2[float32]{{rC[0], rC[1]}, {rC[1], rC[2]}, {rC[2], rC[3]}, {rC[3], rC[0]}}
		// radial segments: check if either radial edge intersects any rect edge
		for _, e := range edges3 {
			if segmentIntersectsSegment(other.Pos, pStart, e[0], e[1]) || segmentIntersectsSegment(other.Pos, pEnd, e[0], e[1]) {
				return true
			}
		}
		// intersection with arc: check if any edge intersects circle and angle at closest point within sector
		for _, e := range edges3 {
			if segmentIntersectsCircle(e[0], e[1], other.Pos, other.Radius) {
				cp := closestPointOnSegment(other.Pos, e[0], e[1])
				angle := math.Atan2(float64(cp.Y-other.Pos.Y), float64(cp.X-other.Pos.X)) * (180.0 / math.Pi)
				if angle < 0 {
					angle += 360.0
				}
				if angleBetweenDeg(float32(angle), other.StartAngle, other.EndAngle) {
					return true
				}
			}
		}
		return false
	}
	return false
}
