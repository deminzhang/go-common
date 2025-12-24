package geom2d

import (
	"testing"
)

func TestPointRectangle(t *testing.T) {
	r := NewRectangle(0, 0, 2, 2, 0)
	p := NewPoint(0.5, 0.5)
	if !p.Intersects(r) {
		t.Fatalf("expected point inside rectangle")
	}
}

func TestPointTriangle(t *testing.T) {
	tri := NewTriangle(0, 0, 5, 0, 0, 5)
	p := NewPoint(1, 1)
	if !p.Intersects(tri) {
		t.Fatalf("expected point inside triangle")
	}
}

func TestCircleSector(t *testing.T) {
	s := NewSector(0, 0, 5, 0, 90)
	c := NewCircle(3, 3, 1)
	if !c.Intersects(s) {
		t.Fatalf("expected circle intersect sector")
	}
}

func TestSectorSector(t *testing.T) {
	s1 := NewSector(0, 0, 5, 0, 90)
	s2 := NewSector(3, 1, 3, 45, 135)
	if !s1.Intersects(s2) || !s2.Intersects(s1) {
		t.Fatalf("expected sectors intersect")
	}
}

func TestRectRect(t *testing.T) {
	r1 := NewRectangle(0, 0, 4, 2, 30)
	r2 := NewRectangle(1, 0, 2, 2, -15)
	if !r1.Intersects(r2) || !r2.Intersects(r1) {
		t.Fatalf("expected rectangles intersect")
	}
}

func TestLineSegmentCircle(t *testing.T) {
	ls := NewLineSegment(-10, 0, 10, 0)
	c := NewCircle(0, 0, 1)
	if !ls.Intersects(c) {
		t.Fatalf("expected line segment intersect circle")
	}
}

func TestTriangleCircle(t *testing.T) {
	tri := NewTriangle(0, 0, 5, 0, 0, 5)
	c := NewCircle(1, 1, 0.5)
	if !tri.Intersects(c) {
		t.Fatalf("expected triangle intersect circle")
	}
}
