package geom2d

import (
	"math"

	"github.com/deminzhang/go-common/vec"
)

func degToRad(deg float32) float64 {
	return float64(deg) * math.Pi / 180.0
}

func normalizeAngleDeg(a float64) float64 {
	a = math.Mod(a, 360.0)
	if a < 0 {
		return a + 360.0
	}
	return a
}

func angleBetweenDeg(angle, start, end float32) bool {
	a := normalizeAngleDeg(float64(angle))
	s := normalizeAngleDeg(float64(start))
	e := normalizeAngleDeg(float64(end))
	if s > e {
		return a >= s || a <= e
	}
	return a >= s && a <= e
}

func rotatePointAround(p, origin vec.Vec2[float32], angleDeg float32) vec.Vec2[float32] {
	rx := float64(p.X - origin.X)
	ry := float64(p.Y - origin.Y)
	rad := degToRad(angleDeg)
	c := math.Cos(rad)
	s := math.Sin(rad)
	x := c*rx - s*ry
	y := s*rx + c*ry
	return vec.Vec2[float32]{X: float32(x) + origin.X, Y: float32(y) + origin.Y}
}

func pointInRectangle(p vec.Vec2[float32], r *Rectangle) bool {
	// rotate point into rectangle local space by -angle
	local := rotatePointAround(p, r.Pos, -r.Angle)
	dx := float32(math.Abs(float64(local.X - r.Pos.X)))
	dy := float32(math.Abs(float64(local.Y - r.Pos.Y)))
	hw := r.Width / 2.0
	hh := r.Height / 2.0
	return dx <= hw && dy <= hh
}

func pointInTriangle(p, a, b, c vec.Vec2[float32]) bool {
	// barycentric method
	px := float64(p.X)
	py := float64(p.Y)
	ax := float64(a.X)
	ay := float64(a.Y)
	bx := float64(b.X)
	by := float64(b.Y)
	cx := float64(c.X)
	cy := float64(c.Y)

	v0x := cx - ax
	v0y := cy - ay
	v1x := bx - ax
	v1y := by - ay
	v2x := px - ax
	v2y := py - ay

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y

	invDenom := 1.0 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return u >= 0 && v >= 0 && (u+v) <= 1
}

func pointInSectorGeneric(p vec.Vec2[float32], s *Sector) bool {
	dx := float64(p.X - s.Pos.X)
	dy := float64(p.Y - s.Pos.Y)
	dsq := dx*dx + dy*dy
	if dsq > float64(s.Radius*s.Radius) {
		return false
	}
	angle := math.Atan2(dy, dx) * (180.0 / math.Pi)
	if angle < 0 {
		angle += 360.0
	}
	return angleBetweenDeg(float32(angle), s.StartAngle, s.EndAngle)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func distPointToSegmentSq(p, a, b vec.Vec2[float32]) float64 {
	px := float64(p.X)
	py := float64(p.Y)
	ax := float64(a.X)
	ay := float64(a.Y)
	bx := float64(b.X)
	by := float64(b.Y)

	dx := bx - ax
	dy := by - ay
	if dx == 0 && dy == 0 {
		dx = px - ax
		dy = py - ay
		return dx*dx + dy*dy
	}
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = clamp(t, 0, 1)
	cx := ax + t*dx
	cy := ay + t*dy
	dx = px - cx
	dy = py - cy
	return dx*dx + dy*dy
}

func segmentIntersectsSegment(p1, p2, q1, q2 vec.Vec2[float32]) bool {
	// orientation based
	ox := func(a, b, c vec.Vec2[float32]) float64 {
		return float64((b.Y-a.Y)*(c.X-a.X) - (b.X-a.X)*(c.Y-a.Y))
	}
	o1 := ox(p1, p2, q1)
	o2 := ox(p1, p2, q2)
	o3 := ox(q1, q2, p1)
	o4 := ox(q1, q2, p2)
	if (o1 == 0 && onSegment(p1, p2, q1)) || (o2 == 0 && onSegment(p1, p2, q2)) || (o3 == 0 && onSegment(q1, q2, p1)) || (o4 == 0 && onSegment(q1, q2, p2)) {
		return true
	}
	return (o1*o2 < 0) && (o3*o4 < 0)
}

func onSegment(a, b, c vec.Vec2[float32]) bool {
	// c on segment ab
	minx := math.Min(float64(a.X), float64(b.X))
	maxx := math.Max(float64(a.X), float64(b.X))
	miny := math.Min(float64(a.Y), float64(b.Y))
	maxy := math.Max(float64(a.Y), float64(b.Y))
	if float64(c.X) >= minx-1e-6 && float64(c.X) <= maxx+1e-6 && float64(c.Y) >= miny-1e-6 && float64(c.Y) <= maxy+1e-6 {
		// colinear is caller's responsibility
		return true
	}
	return false
}

func segmentIntersectsCircle(p1, p2, center vec.Vec2[float32], r float32) bool {
	dsq := distPointToSegmentSq(center, p1, p2)
	return dsq <= float64(r*r)+1e-6
}

func closestPointOnSegment(p, a, b vec.Vec2[float32]) vec.Vec2[float32] {
	px := float64(p.X)
	py := float64(p.Y)
	ax := float64(a.X)
	ay := float64(a.Y)
	bx := float64(b.X)
	by := float64(b.Y)
	dx := bx - ax
	dy := by - ay
	if dx == 0 && dy == 0 {
		return a
	}
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = clamp(t, 0, 1)
	cx := ax + t*dx
	cy := ay + t*dy
	return vec.Vec2[float32]{X: float32(cx), Y: float32(cy)}
}

func rectangleCorners(r *Rectangle) []vec.Vec2[float32] {
	hw := r.Width / 2.0
	hh := r.Height / 2.0
	corners := []vec.Vec2[float32]{
		{X: hw, Y: hh},
		{X: -hw, Y: hh},
		{X: -hw, Y: -hh},
		{X: hw, Y: -hh},
	}
	res := make([]vec.Vec2[float32], 0, 4)
	rad := degToRad(r.Angle)
	c := math.Cos(rad)
	s := math.Sin(rad)
	for _, p := range corners {
		x := float64(p.X)*c - float64(p.Y)*s
		y := float64(p.X)*s + float64(p.Y)*c
		res = append(res, vec.Vec2[float32]{X: r.Pos.X + float32(x), Y: r.Pos.Y + float32(y)})
	}
	return res
}

func rectAxes(r *Rectangle) []vec.Vec2[float32] {
	rad := degToRad(r.Angle)
	c := float64(math.Cos(rad))
	s := float64(math.Sin(rad))
	// axis along width and height
	return []vec.Vec2[float32]{
		{X: float32(c), Y: float32(s)},
		{X: float32(-s), Y: float32(c)},
	}
}

func dot(u, v vec.Vec2[float32]) float64 {
	return float64(u.X)*float64(v.X) + float64(u.Y)*float64(v.Y)
}

func projectPointsAxis(points []vec.Vec2[float32], axis vec.Vec2[float32]) (float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, p := range points {
		proj := dot(p, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}
	return min, max
}

func projectionsOverlap(aMin, aMax, bMin, bMax float64) bool {
	return !(aMax < bMin || bMax < aMin)
}

func rectRectIntersectSAT(a, b *Rectangle) bool {
	aCorners := rectangleCorners(a)
	bCorners := rectangleCorners(b)
	axes := append(rectAxes(a), rectAxes(b)...)
	for _, axis := range axes {
		minA, maxA := projectPointsAxis(aCorners, axis)
		minB, maxB := projectPointsAxis(bCorners, axis)
		if !projectionsOverlap(minA, maxA, minB, maxB) {
			return false
		}
	}
	return true
}

// return intersection points of two circles (may be 0,1,2)
func circleCircleIntersections(c1, c2 vec.Vec2[float32], r1, r2 float32) ([]vec.Vec2[float32], bool) {
	x0 := float64(c1.X)
	y0 := float64(c1.Y)
	x1 := float64(c2.X)
	y1 := float64(c2.Y)
	R0 := float64(r1)
	R1 := float64(r2)
	dx := x1 - x0
	dy := y1 - y0
	d := math.Hypot(dx, dy)
	if d > R0+R1+1e-6 || d < math.Abs(R0-R1)-1e-6 {
		return nil, false
	}
	if d == 0 && R0 == R1 {
		return nil, false
	}
	// a = (R0^2 - R1^2 + d^2) / (2d)
	a := (R0*R0 - R1*R1 + d*d) / (2 * d)
	h := math.Sqrt(math.Max(0, R0*R0-a*a))
	xm := x0 + a*(dx)/d
	ym := y0 + a*(dy)/d
	rx := -dy * (h / d)
	ry := dx * (h / d)
	p1 := vec.Vec2[float32]{X: float32(xm + rx), Y: float32(ym + ry)}
	p2 := vec.Vec2[float32]{X: float32(xm - rx), Y: float32(ym - ry)}
	if h == 0 {
		return []vec.Vec2[float32]{p1}, true
	}
	return []vec.Vec2[float32]{p1, p2}, true
}
