package vec

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Vec2[T constraints.Integer | constraints.Float] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

func (this *Vec2[T]) Equal(v Vec2[T]) bool {
	return *this == v
}

func (this *Vec2[T]) Set(x, y T) {
	this.X = x
	this.Y = y
}

func (this *Vec2[T]) Clone() Vec2[T] {
	return Vec2[T]{X: this.X, Y: this.Y}
}

func (this *Vec2[T]) Add(v Vec2[T]) {
	this.X += v.X
	this.Y += v.Y
}

func (this *Vec2[T]) Sub(v Vec2[T]) {
	this.X -= v.X
	this.Y -= v.Y
}

func (this *Vec2[T]) Multiply(scalar T) {
	this.X *= scalar
	this.Y *= scalar
}

func (this *Vec2[T]) Divide(scalar T) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = T(math.Inf(1))
		this.Y = T(math.Inf(1))
	}
	this.X /= scalar
	this.Y /= scalar
}

func (this *Vec2[T]) Scale(v Vec2[T]) {
	this.X *= v.X
	this.Y *= v.Y
}

func (this *Vec2[T]) Dot(v Vec2[T]) T {
	return this.X*v.X + this.Y*v.Y
}

func (this *Vec2[T]) Magnitude() T {
	return T(math.Sqrt(float64((this.X * this.X) + (this.Y * this.Y))))
}

func (this *Vec2[T]) MagnitudeSqr() T {
	return this.X*this.X + this.Y*this.Y
}

func (this *Vec2[T]) Normalize() {
	mag := this.Magnitude()
	if mag == 0 {
		this.X = 0
		this.Y = 0
		return
	}
	this.X /= mag
	this.Y /= mag
}

func (this *Vec2[T]) Normalized() Vec2[T] {
	mag := this.Magnitude()
	if mag == 0 {
		return Vec2[T]{X: 0, Y: 0}
	}
	return Vec2[T]{X: this.X / mag, Y: this.Y / mag}
}

func (this *Vec2[T]) Distance(v Vec2[T]) T {
	return T(math.Sqrt(float64((this.X-v.X)*(this.X-v.X) + (this.Y-v.Y)*(this.Y-v.Y))))
}

func (this *Vec2[T]) DistanceSqr(v Vec2[T]) T {
	return (this.X-v.X)*(this.X-v.X) + (this.Y-v.Y)*(this.Y-v.Y)
}

func (this *Vec2[T]) Angle(v Vec2[T]) T {
	return T(math.Atan2(float64(this.Y-v.Y), float64(this.X-v.X)))
}

func (this *Vec2[T]) AngleDeg(v Vec2[T]) T {
	return T(math.Atan2(float64(this.Y-v.Y), float64(this.X-v.X)) * 180 / math.Pi)
}

func (this *Vec2[T]) Rotate(angle T) {
	cos := T(math.Cos(float64(angle)))
	sin := T(math.Sin(float64(angle)))
	x := this.X*cos - this.Y*sin
	y := this.X*sin + this.Y*cos
	this.X = x
	this.Y = y
}

func (this *Vec2[T]) RotateDeg(angle T) {
	cos := T(math.Cos(float64(angle) * math.Pi / 180))
	sin := T(math.Sin(float64(angle) * math.Pi / 180))
	x := this.X*cos - this.Y*sin
	y := this.X*sin + this.Y*cos
	this.X = x
	this.Y = y
}

func (this *Vec2[T]) ScaleToLength(newLength T) {
	mag := this.Magnitude()
	if mag == 0 {
		this.X = 0
		this.Y = 0
		return
	}
	this.X *= newLength / mag
	this.Y *= newLength / mag
}

