package vec

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Vec3[T constraints.Integer | constraints.Float] struct {
	X T `json:"x"`
	Y T `json:"y"`
	Z T `json:"z"`
}

// Warn: 浮点相等判定确认应用环境是否合适
func (v3 Vec3[T]) Equal(v Vec3[T]) bool {
	return v3 == v
}

func (v3 *Vec3[T]) Set(x, y, z T) {
	v3.X = x
	v3.Y = y
	v3.Z = z
}

func (v3 Vec3[T]) Clone() Vec3[T] {
	return Vec3[T]{X: v3.X, Y: v3.Y, Z: v3.Z}
}

func (v3 *Vec3[T]) Add(v Vec3[T]) {
	v3.X += v.X
	v3.Y += v.Y
	v3.Z += v.Z
}

func (v3 Vec3[T]) Added(v Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v3.X + v.X, Y: v3.Y + v.Y, Z: v3.Z + v.Z}
}

func (v3 *Vec3[T]) Sub(v Vec3[T]) {
	v3.X -= v.X
	v3.Y -= v.Y
	v3.Z -= v.Z
}
func (v3 Vec3[T]) Subtracted(v Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v3.X - v.X, Y: v3.Y - v.Y, Z: v3.Z - v.Z}
}

func (v3 *Vec3[T]) Multiply(scalar T) {
	v3.X *= scalar
	v3.Y *= scalar
	v3.Z *= scalar
}

func (v3 Vec3[T]) Multiplied(scalar T) Vec3[T] {
	return Vec3[T]{X: v3.X * scalar, Y: v3.Y * scalar, Z: v3.Z * scalar}
}

func (v3 *Vec3[T]) Divide(scalar T) {
	if scalar == 0 {
		//panic("v/0！")
		v3.X = T(math.Inf(1))
		v3.Y = T(math.Inf(1))
		v3.Z = T(math.Inf(1))
		return
	}
	v3.X /= scalar
	v3.Y /= scalar
	v3.Z /= scalar
}

func (v3 Vec3[T]) Divided(scalar T) Vec3[T] {
	if scalar == 0 {
		return Vec3[T]{
			X: T(math.Inf(1)),
			Y: T(math.Inf(1)),
			Z: T(math.Inf(1)),
		}
	}
	return Vec3[T]{X: v3.X / scalar, Y: v3.Y / scalar, Z: v3.Z / scalar}
}

func (v3 *Vec3[T]) Scale(v Vec3[T]) {
	v3.X *= v.X
	v3.Y *= v.Y
	v3.Z *= v.Z
}

func (v3 Vec3[T]) Scaled(v Vec3[T]) Vec3[T] {
	return Vec3[T]{X: v3.X * v.X, Y: v3.Y * v.Y, Z: v3.Z * v.Z}
}

func (v3 Vec3[T]) Dot(v Vec3[T]) T {
	return v3.X*v.X + v3.Y*v.Y + v3.Z*v.Z
}

func (v3 Vec3[T]) Magnitude() T {
	return T(math.Sqrt(float64((v3.X * v3.X) + (v3.Y * v3.Y) + (v3.Z * v3.Z))))
}

func (v3 Vec3[T]) MagnitudeSqr() T {
	return v3.X*v3.X + v3.Y*v3.Y + v3.Z*v3.Z
}

func (v3 *Vec3[T]) Normalize() {
	mag := v3.Magnitude()
	if mag == 0 {
		v3.X = 0
		v3.Y = 0
		v3.Z = 0
		return
	}
	v3.X /= mag
	v3.Y /= mag
	v3.Z /= mag
}

func (v3 Vec3[T]) Normalized() Vec3[T] {
	mag := v3.Magnitude()
	if mag == 0 {
		return Vec3[T]{X: 0, Y: 0, Z: 0}
	}
	return Vec3[T]{X: v3.X / mag, Y: v3.Y / mag, Z: v3.Z / mag}
}

func (v3 Vec3[T]) Cross(v Vec3[T]) Vec3[T] {
	return Vec3[T]{
		X: v3.Y*v.Z - v3.Z*v.Y,
		Y: v3.Z*v.X - v3.X*v.Z,
		Z: v3.X*v.Y - v3.Y*v.X,
	}
}

// 线性插值
func (v3 Vec3[T]) Lerp(v Vec3[T], t T) Vec3[T] {
	return Vec3[T]{
		X: v3.X + (v.X-v3.X)*t,
		Y: v3.Y + (v.Y-v3.Y)*t,
		Z: v3.Z + (v.Z-v3.Z)*t,
	}
}

func (v3 Vec3[T]) LerpUnclamped(v Vec3[T], t T) Vec3[T] {
	return v.Lerp(v3, 1-t)
}

// 距离计算
func (v3 Vec3[T]) Distance(v Vec3[T]) T {
	dx := float64(v3.X - v.X)
	dy := float64(v3.Y - v.Y)
	dz := float64(v3.Z - v.Z)
	return T(math.Sqrt(dx*dx + dy*dy + dz*dz))
}

func (v3 Vec3[T]) DistanceSqr(v Vec3[T]) T {
	dx := float64(v3.X - v.X)
	dy := float64(v3.Y - v.Y)
	dz := float64(v3.Z - v.Z)
	return T(dx*dx + dy*dy + dz*dz)
}

// 向量投影
func (v3 Vec3[T]) ProjectOn(v Vec3[T]) Vec3[T] {
	magSqr := v.MagnitudeSqr()
	if magSqr == 0 {
		return Vec3[T]{X: 0, Y: 0, Z: 0}
	}
	dot := v3.Dot(v)
	scale := dot / magSqr
	return v.Multiplied(scale)
}

// 反射向量
func (v3 Vec3[T]) Reflect(normal Vec3[T]) Vec3[T] {
	dot := v3.Dot(normal)
	return Vec3[T]{
		X: v3.X - 2*dot*normal.X,
		Y: v3.Y - 2*dot*normal.Y,
		Z: v3.Z - 2*dot*normal.Z,
	}
}

func (v3 *Vec3[T]) RotateX(angle int) Vec3[T] {
	rad := float64(angle) * math.Pi / 180.0
	// sin := mathtable.SinByAngle(angle)
	// cos := mathtable.CosByAngle(angle)
	cos := T(math.Cos(rad))
	sin := T(math.Sin(rad))
	y := v3.Y*cos - v3.Z*sin
	z := v3.Y*sin + v3.Z*cos
	return Vec3[T]{X: v3.X, Y: y, Z: z}
}

func (v3 *Vec3[T]) RotateY(angle int) Vec3[T] {
	rad := float64(angle) * math.Pi / 180.0
	// sin := mathtable.SinByAngle(angle)
	// cos := mathtable.CosByAngle(angle)
	cos := T(math.Cos(rad))
	sin := T(math.Sin(rad))
	x := v3.Z*sin + v3.X*cos
	z := v3.Z*cos - v3.X*sin
	return Vec3[T]{X: x, Y: v3.Y, Z: z}
}

func (v3 *Vec3[T]) RotateZ(angle int) Vec3[T] {
	rad := float64(angle) * math.Pi / 180.0
	// sin := mathtable.SinByAngle(angle)
	// cos := mathtable.CosByAngle(angle)
	cos := T(math.Cos(rad))
	sin := T(math.Sin(rad))
	x := v3.X*cos - v3.Y*sin
	y := v3.X*sin + v3.Y*cos
	return Vec3[T]{X: x, Y: y, Z: v3.Z}
}

func (v3 Vec3[T]) XYToVec2() Vec2[T] {
	return Vec2[T]{X: v3.X, Y: v3.Y}
}

func (v3 Vec3[T]) XZToVec2() Vec2[T] {
	return Vec2[T]{X: v3.X, Y: v3.Z}
}
