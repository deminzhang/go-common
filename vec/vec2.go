package vec

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Vec2[T constraints.Integer | constraints.Float] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

// Warn: 浮点相等判定确认应用环境是否合适
func (v Vec2[T]) Equal(v2 Vec2[T]) bool {
	return v == v2
}

func (v *Vec2[T]) Set(x, y T) {
	v.X = x
	v.Y = y
}

func (v Vec2[T]) Length() float64 {
	return math.Sqrt(float64((v.X * v.X) + (v.Y * v.Y)))
}

func (v Vec2[T]) Clone() Vec2[T] {
	return Vec2[T]{X: v.X, Y: v.Y}
}

func (v *Vec2[T]) Add(v2 Vec2[T]) {
	v.X += v2.X
	v.Y += v2.Y
}

// 向量加法（返回新向量）
func (v Vec2[T]) Added(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v *Vec2[T]) Sub(v2 Vec2[T]) {
	v.X -= v2.X
	v.Y -= v2.Y
}

// 向量减法（返回新向量）
func (v Vec2[T]) Subtracted(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X - v2.X, Y: v.Y - v2.Y}
}

func (v *Vec2[T]) Multiply(scalar float64) {
	v.X = T(float64(v.X) * scalar)
	v.Y = T(float64(v.Y) * scalar)
}

// 标量乘法（返回新向量）
func (v Vec2[T]) Multiplied(scalar float64) Vec2[T] {
	return Vec2[T]{
		X: T(float64(v.X) * scalar),
		Y: T(float64(v.Y) * scalar),
	}
}

func (v *Vec2[T]) Divide(scalar T) {
	if scalar == 0 {
		//panic("v/0！")
		v.X = T(math.Inf(1))
		v.Y = T(math.Inf(1))
		return
	}
	v.X /= scalar
	v.Y /= scalar
}

// 标量除法（返回新向量）
func (v Vec2[T]) Divided(scalar T) Vec2[T] {
	if scalar == 0 {
		return Vec2[T]{
			X: T(math.Inf(1)),
			Y: T(math.Inf(1)),
		}
	}
	return Vec2[T]{X: v.X / scalar, Y: v.Y / scalar}
}

// 向量：分向量乘
func (v *Vec2[T]) Scale(v2 Vec2[T]) {
	v.X *= v2.X
	v.Y *= v2.Y
}

func (v Vec2[T]) Scaled(v2 Vec2[T]) Vec2[T] {
	return Vec2[T]{X: v.X * v2.X, Y: v.Y * v2.Y}
}

// 向量：点积
func (v Vec2[T]) Dot(v2 Vec2[T]) T {
	return v.X*v2.X + v.Y*v2.Y
}

// 向量：长度
func (v Vec2[T]) Magnitude() T {
	return T(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// 向量：长度平方
func (v Vec2[T]) SqrMagnitude() T {
	return v.X*v.X + v.Y*v.Y
}

// 向量：单位化 (0向量禁用)
func (v *Vec2[T]) Normalize() {
	l := v.Magnitude()
	if l == 0 {
		v.Set(1, 0)
		return
	}
	v.Divide(l)
}

// 向量：单位化值 (0向量禁用)
func (v Vec2[T]) Normalized() Vec2[T] {
	result := v.Clone()
	result.Normalize()
	return result
}

// 朝向不变拉长度
func (v *Vec2[T]) ScaleToLength(newLength float64) {
	v.Normalize()
	v.Multiply(newLength)
}

// 复制朝向定长
func (v Vec2[T]) ScaledToLength(newLength float64) Vec2[T] {
	result := v.Clone()
	result.ScaleToLength(newLength)
	return result
}

func (v Vec2[T]) MoveTowards(targetX T, targetY T, speed float64) Vec2[T] {
	dir := Vec2[T]{X: targetX - v.X, Y: targetY - v.Y}
	dis := dir.Length()
	if dis <= speed {
		return Vec2[T]{X: targetX, Y: targetY}
	}
	dir.Normalize()
	dir.Multiply(speed)
	newVec := Vec2[T]{X: v.X + dir.X, Y: v.Y + dir.Y}
	return newVec
}

func (v Vec2[T]) Distance(v2 Vec2[T]) T {
	dx := float64(v.X - v2.X)
	dy := float64(v.Y - v2.Y)
	return T(math.Sqrt(dx*dx + dy*dy))
}

// 向量：距离平方 判定用节省开方开销
func (v Vec2[T]) DistanceSqr(v2 Vec2[T]) T {
	dx := float64(v.X - v2.X)
	dy := float64(v.Y - v2.Y)
	return T(dx*dx + dy*dy)
}

func (v Vec2[T]) AngleTo(v2 Vec2[T]) float64 {
	dot := float64(v.Dot(v2))
	magV1 := v.Magnitude()
	magV2 := v2.Magnitude()
	if magV1 == 0 || magV2 == 0 {
		return 0
	}
	cosTheta := dot / (float64(magV1) * float64(magV2))
	if cosTheta > 1 {
		cosTheta = 1
	} else if cosTheta < -1 {
		cosTheta = -1
	}
	angleRad := math.Acos(cosTheta)
	angleDeg := angleRad * (180.0 / math.Pi)
	return angleDeg
}

func (v Vec2[T]) Rotate(angleDeg float64) Vec2[T] {
	angleRad := angleDeg * (math.Pi / 180.0)
	// sinA := mathtable.SinByAngle(angleRad)
	// cosA := mathtable.CosByAngle(angleRad)
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)
	xNew := float64(v.X)*cosA - float64(v.Y)*sinA
	yNew := float64(v.X)*sinA + float64(v.Y)*cosA
	return Vec2[T]{X: T(xNew), Y: T(yNew)}
}
