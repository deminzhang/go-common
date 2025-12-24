package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

type Circle struct {
	BaseShape
	Radius float32
}

func NewCircle(x, y, r float32) *Circle {
	return &Circle{BaseShape: BaseShape{Pos: vec.Vec2[float32]{X: x, Y: y}}, Radius: r}
}

func (c *Circle) Intersects(target IShape) bool {
	switch other := target.(type) {
	case *Point:
		return other.intersectsCircle(c)
	case *Circle:
		return c.intersectsCircle(other)
	case *Sector:
		return c.intersectsSector(other)
	case *Rectangle:
		return other.Intersects(c)
	case *LineSegment:
		return other.Intersects(c)
	case *Triangle:
		return other.Intersects(c)
	}
	return false
}

func (c *Circle) intersectsCircle(other *Circle) bool {
	disSq := c.Pos.DistanceSqr(other.Pos)
	return disSq < (c.Radius+other.Radius)*(c.Radius+other.Radius)
}
func (c *Circle) intersectsSector(sector *Sector) bool {
	// if circle center inside sector area
	if pointInSectorGeneric(c.Pos, sector) {
		return true
	}
	// quick reject by circle distances
	dx := float64(c.Pos.X - sector.Pos.X)
	dy := float64(c.Pos.Y - sector.Pos.Y)
	d := math.Hypot(dx, dy)
	if d > float64(c.Radius+sector.Radius) {
		return false
	}
	// check arc closest point
	angle := math.Atan2(float64(c.Pos.Y-sector.Pos.Y), float64(c.Pos.X-sector.Pos.X)) * (180.0 / math.Pi)
	if angle < 0 {
		angle += 360.0
	}
	if angleBetweenDeg(float32(angle), sector.StartAngle, sector.EndAngle) {
		// closest point on arc
		rad := degToRad(float32(angle))
		ax := sector.Pos.X + float32(math.Cos(rad))*sector.Radius
		ay := sector.Pos.Y + float32(math.Sin(rad))*sector.Radius
		d2 := float64((c.Pos.X-ax)*(c.Pos.X-ax) + (c.Pos.Y-ay)*(c.Pos.Y-ay))
		if d2 <= float64(c.Radius*c.Radius)+1e-6 {
			return true
		}
	}
	// check intersection with radial edges
	startRad := degToRad(sector.StartAngle)
	endRad := degToRad(sector.EndAngle)
	pStart := vec.Vec2[float32]{X: sector.Pos.X + float32(math.Cos(startRad))*sector.Radius, Y: sector.Pos.Y + float32(math.Sin(startRad))*sector.Radius}
	pEnd := vec.Vec2[float32]{X: sector.Pos.X + float32(math.Cos(endRad))*sector.Radius, Y: sector.Pos.Y + float32(math.Sin(endRad))*sector.Radius}
	if segmentIntersectsCircle(sector.Pos, pStart, c.Pos, c.Radius) || segmentIntersectsCircle(sector.Pos, pEnd, c.Pos, c.Radius) {
		return true
	}
	return false
}
