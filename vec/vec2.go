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

func (this *Vec2[T]) Length() float64 {
	return math.Sqrt(float64((this.X * this.X) + (this.Y * this.Y)))
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

func (this *Vec2[T]) Multiply(scalar float64) {
	this.X = T(float64(this.X) * scalar)
	this.Y = T(float64(this.Y) * scalar)
}

func (this *Vec2[T]) Divide(scalar T) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = T(math.Inf(1))
		this.Y = T(math.Inf(1))
		return
	}
	this.X /= scalar
	this.Y /= scalar
}

// 向量：分向量乘
func (this *Vec2[T]) Scale(v Vec2[T]) {
	this.X *= v.X
	this.Y *= v.Y
}

// 向量：点积
func (this *Vec2[T]) Dot(v Vec2[T]) T {
	return this.X*v.X + this.Y*v.Y
}

// 向量：长度
func (this *Vec2[T]) Magnitude() T {
	return T(math.Sqrt(float64(this.X*this.X + this.Y*this.Y)))
}

// 向量：长度平方
func (this *Vec2[T]) SqrMagnitude() T {
	return this.X*this.X + this.Y*this.Y
}

// 向量：单位化 (0向量禁用)
func (this *Vec2[T]) Normalize() {
	l := this.Magnitude()
	if l == 0 {
		this.Set(1, 0)
		return
	}
	this.Divide(l)
}

// 向量：单位化值 (0向量禁用)
func (this *Vec2[T]) Normalized() Vec2[T] {
	v := this.Clone()
	v.Normalize()
	return v
}

// 朝向不变拉长度
func (this *Vec2[T]) ScaleToLength(newLength float64) {
	this.Normalize()
	this.Multiply(newLength)
}

// 复制朝向定长
func (this *Vec2[T]) ScaledToLength(newLength float64) Vec2[T] {
	v := this.Clone()
	v.ScaleToLength(newLength)
	return v
}

func (this *Vec2[T]) MoveTowards(targetX T, targetY T, speed float64) Vec2[T] {
	dir := Vec2[T]{X: targetX - this.X, Y: targetY - this.Y}
	dis := dir.Length()
	if dis <= speed {
		return Vec2[T]{X: targetX, Y: targetY}
	}
	dir.Normalize()
	dir.Multiply(speed)
	newVec := Vec2[T]{X: this.X + dir.X, Y: this.Y + dir.Y}
	return newVec
}
