package vec

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Vec[T constraints.Integer | constraints.Float] struct {
	X T `json:"x"`
	Y T `json:"y"`
	Z T `json:"z"`
}

func (this *Vec[T]) Equal(v Vec[T]) bool {
	return *this == v
}

func (this *Vec[T]) Set(x, y, z T) {
	this.X = x
	this.Y = y
	this.Z = z
}

func (this *Vec[T]) Clone() Vec[T] {
	return Vec[T]{X: this.X, Y: this.Y, Z: this.Z}
}

func (this *Vec[T]) Add(v Vec[T]) {
	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
}

func (this *Vec[T]) Sub(v Vec[T]) {
	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
}

func (this *Vec[T]) Multiply(scalar T) {
	this.X *= scalar
	this.Y *= scalar
	this.Z *= scalar
}

func (this *Vec[T]) Divide(scalar T) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = T(math.Inf(1))
		this.Y = T(math.Inf(1))
		this.Z = T(math.Inf(1))
		return
	}
	this.X /= scalar
	this.Y /= scalar
	this.Z /= scalar
}

func (this *Vec[T]) Scale(v Vec[T]) {
	this.X *= v.X
	this.Y *= v.Y
	this.Z *= v.Z
}

func (this *Vec[T]) Dot(v Vec[T]) T {
	return this.X*v.X + this.Y*v.Y + this.Z*v.Z
}

func (this *Vec[T]) Magnitude() T {
	return T(math.Sqrt(float64((this.X * this.X) + (this.Y * this.Y) + (this.Z * this.Z))))
}

func (this *Vec[T]) MagnitudeSqr() T {
	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}

func (this *Vec[T]) Normalize() {
	mag := this.Magnitude()
	if mag == 0 {
		this.X = 0
		this.Y = 0
		this.Z = 0
		return
	}
	this.X /= mag
	this.Y /= mag
	this.Z /= mag
}

func (this *Vec[T]) Normalized() Vec[T] {
	mag := this.Magnitude()
	if mag == 0 {
		return Vec[T]{X: 0, Y: 0, Z: 0}
	}
	return Vec[T]{X: this.X / mag, Y: this.Y / mag, Z: this.Z / mag}
}

func (this *Vec[T]) Cross(v Vec[T]) Vec[T] {
	return Vec[T]{
		X: this.Y*v.Z - this.Z*v.Y,
		Y: this.Z*v.X - this.X*v.Z,
		Z: this.X*v.Y - this.Y*v.X,
	}
}

func (this *Vec[T]) Lerp(v Vec[T], t T) Vec[T] {
	return Vec[T]{
		X: this.X + (v.X-this.X)*t,
		Y: this.Y + (v.Y-this.Y)*t,
		Z: this.Z + (v.Z-this.Z)*t,
	}
}
