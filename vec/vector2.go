package vec

import (
	"math"
)

type Vector2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (this *Vector2) Equal(v1 Vector2) bool {
	//return this.X == v.X && this.Y == v.Y
	return *this == v1
}

func (this *Vector2) Set(x, y float32) {
	this.X = x
	this.Y = y
}

func (this *Vector2) Length() float32 {
	return float32(math.Sqrt(float64((this.X * this.X) + (this.Y * this.Y))))
}

func (this *Vector2) Clone() Vector2 {
	return NewVector2(this.X, this.Y)
}

func (this *Vector2) Add(v1 Vector2) {
	this.X += v1.X
	this.Y += v1.Y
}

func (this *Vector2) Sub(v1 Vector2) {
	this.X -= v1.X
	this.Y -= v1.Y
}

func (this *Vector2) Multiply(scalar float32) {
	this.X *= scalar
	this.Y *= scalar
}

func (this *Vector2) Divide(scalar float32) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = float32(math.Inf(1))
		this.Y = float32(math.Inf(1))
		return
	}
	this.X /= scalar
	this.Y /= scalar
}

// 向量：分向量乘
func (this *Vector2) Scale(v Vector2) {
	this.X *= v.X
	this.Y *= v.Y
}

// 向量：点积
func (this *Vector2) Dot(v Vector2) float32 {
	return this.X*v.X + this.Y*v.Y
}

// 向量：长度
func (this *Vector2) Magnitude() float32 {
	return float32(math.Sqrt(float64(this.X*this.X + this.Y*this.Y)))
}

// 向量：长度平方
func (this *Vector2) SqrMagnitude() float32 {
	return this.X*this.X + this.Y*this.Y
}

// 向量：单位化 (0向量禁用)
func (this *Vector2) Normalize() {
	l := this.Magnitude()
	if l == 0 {
		this.Set(1, 0)
		return
	}
	this.Divide(l)
}

// 向量：单位化值 (0向量禁用)
func (this *Vector2) Normalized() Vector2 {
	v := this.Clone()
	v.Normalize()
	return v
}

// 朝向不变拉长度
func (this *Vector2) ScaleToLength(newLength float32) {
	this.Normalize()
	this.Multiply(newLength)
}

// 复制朝向定长
func (this *Vector2) ScaledToLength(newLength float32) Vector2 {
	v := this.Clone()
	v.ScaleToLength(newLength)
	return v
}

// 返回：新向量
func NewVector2(x, y float32) Vector2 {
	return Vector2{X: x, Y: y}
}

// 返回：零向量(0,0,0)
func ZeroV2() Vector2 {
	return Vector2{X: 0, Y: 0}
}

// X 轴 单位向量
func XAxisV2() Vector2 {
	return Vector2{X: 1, Y: 0}
}

// Y 轴 单位向量
func YAxisV2() Vector2 {
	return Vector2{X: 0, Y: 1}
}

func XYAxisV2() Vector2 {
	return Vector2{X: 1, Y: 1}
}

func PositiveInfinityV2() Vector2 {
	return Vector2{X: float32(math.Inf(1)), Y: float32(math.Inf(1))}
}
func NegativeInfinityV2() Vector2 {
	return Vector2{X: float32(math.Inf(-1)), Y: float32(math.Inf(-1))}
}

// 返回：a + b 向量
func AddV2(a, b Vector2) Vector2 {
	return Vector2{X: a.X + b.X, Y: a.Y + b.Y}
}

// 返回：a - b 向量
func SubV2(a, b Vector2) Vector2 {
	return Vector2{X: a.X - b.X, Y: a.Y - b.Y}
}

func ScaleV2(a, b Vector3) Vector2 {
	return Vector2{X: a.X * b.X, Y: a.Y * b.Y}
}

func AddArrayV2(vs []Vector2, dv Vector2) []Vector2 {
	for i, _ := range vs {
		vs[i].Add(dv)
	}
	return vs
}

func MultiplyV2(v Vector2, scalars []float32) []Vector2 {
	var vs []Vector2
	for _, value := range scalars {
		vector := v.Clone()
		vector.Multiply(value)
		vs = append(vs, vector)
	}
	return vs
}

// 求两点间距离
func DistanceV2(a Vector2, b Vector2) float32 {
	return float32(math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2)))
}

// 求两点间距离平方
func DistanceSqrV2(a Vector2, b Vector2) float32 {
	return float32(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

// 线性插值
func LerpV2(a, b Vector2, t float32) Vector2 {
	return Vector2{X: a.X + (b.X-a.X)*t, Y: a.Y + (b.Y-a.Y)*t}
}

func LerpUnclampedV2(a, b Vector2, t float32) Vector2 {
	return LerpV2(b, a, 1-t)
}

// 如果需要再加
//Max	返回由两个向量的最大分量组成的向量。
//Min	返回由两个向量的最小分量组成的向量。
//MoveTowards	将点 current 移向 /target/。
//Reflect	从法线定义的向量反射一个向量。
//Scale	将两个向量的分量相乘。
//SignedAngle	返回 from 与 to 之间的有符号角度（以度为单位）。
//SmoothDamp	随时间推移将一个向量逐渐改变为所需目标。
