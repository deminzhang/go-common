package vec

import (
	"math"

	. "github.com/deminzhang/go-common/fix64"
)

type FixVector2 struct {
	X Fix64 `json:"x"`
	Y Fix64 `json:"y"`
}

func (this *FixVector2) Equal(v FixVector2) bool {
	//return this.X == v.X && this.Y == v.Y
	return *this == v
}

func (this *FixVector2) Set(x, y Fix64) {
	this.X = x
	this.Y = y
}

func (this *FixVector2) Clone() FixVector2 {
	return NewFixVector2(this.X, this.Y)
}

func (this *FixVector2) Add(v FixVector2) {
	this.X += v.X
	this.Y += v.Y
}

func (this *FixVector2) Sub(v FixVector2) {
	this.X -= v.X
	this.Y -= v.Y
}

func (this *FixVector2) Multiply(scalar Fix64) {
	this.X = this.X.Mul(scalar)
	this.Y = this.Y.Mul(scalar)
}

func (this *FixVector2) Divide(scalar Fix64) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = Fix64(math.Inf(1))
		this.Y = Fix64(math.Inf(1))
	}
	this.X = this.X.Div(scalar)
	this.Y = this.Y.Div(scalar)
}

// 向量：分向量乘
func (this *FixVector2) Scale(v FixVector2) {
	this.X = this.X.Mul(v.X)
	this.Y = this.Y.Mul(v.Y)
}

// 向量：点积
func (this *FixVector2) Dot(v FixVector2) Fix64 {
	return this.X.Mul(v.X) + this.Y.Mul(v.Y)
}

// 向量：长度
func (this *FixVector2) Magnitude() Fix64 {
	return (this.X.Mul(this.X) + this.Y.Mul(this.Y)).Sqrt()
}

// 向量：长度平方
func (this *FixVector2) SqrMagnitude() Fix64 {
	return this.X.Mul(this.X) + this.Y.Mul(this.Y)
}

// 向量：单位化 (0向量禁用)
func (this *FixVector2) Normalize() {
	l := this.Magnitude()
	if l == 0 {
		this.Set(FixOne, 0)
		return
	}
	this.Divide(l)
}

// 单位化值 (0向量禁用)
func (this *FixVector2) Nomalized() FixVector2 {
	v := this.Clone()
	v.Normalize()
	return v
}

// 朝向不变拉长度
func (this *FixVector2) ScaleToLength(newLength Fix64) {
	this.Normalize()
	this.Multiply(newLength)
}

// 复制朝向定长
func (this *FixVector2) ScaledToLength(newLength Fix64) FixVector2 {
	v := this.Clone()
	v.ScaleToLength(newLength)
	return v
}

func (a *FixVector2) Vector2() Vector2 {
	return Vector2{
		X: a.X.Float32(),
		Y: a.Y.Float32(),
	}
}

// 返回：新向量
func NewFixVector2(x, y Fix64) FixVector2 {
	return FixVector2{X: x, Y: y}
}

// 返回：零向量(0,0,0)
func ZeroFV2() FixVector2 {
	return FixVector2{X: 0, Y: 0}
}

// X 轴 单位向量
func XAxisFV2() FixVector2 {
	return FixVector2{X: NewFix64(1), Y: 0}
}

// Y 轴 单位向量
func YAxisFV2() FixVector2 {
	return FixVector2{X: 0, Y: NewFix64(1)}
}

func XYAxisFV2() FixVector2 {
	return FixVector2{X: NewFix64(1), Y: NewFix64(1)}
}

func PositiveInfinityFV2() FixVector2 {
	return FixVector2{X: Fix64(math.Inf(1)), Y: Fix64(math.Inf(1))}
}
func NegativeInfinityFV2() FixVector2 {
	return FixVector2{X: Fix64(math.Inf(-1)), Y: Fix64(math.Inf(-1))}
}

// 返回：a + b 向量
func AddFV2(a, b FixVector2) FixVector2 {
	return FixVector2{X: a.X + b.X, Y: a.Y + b.Y}
}

// 返回：a - b 向量
func SubFV2(a, b FixVector2) FixVector2 {
	return FixVector2{X: a.X - b.X, Y: a.Y - b.Y}
}

func AddArrayFV2(vs []FixVector2, dv FixVector2) []FixVector2 {
	for i := range vs {
		vs[i].Add(dv)
	}
	return vs
}

func MultiplyFV2(v FixVector2, scalars []Fix64) []FixVector2 {
	var vs []FixVector2
	for _, value := range scalars {
		vector := v.Clone()
		vector.Multiply(value)
		vs = append(vs, vector)
	}
	return vs
}

func ScaleFV2(a, b FixVector2) FixVector2 {
	return FixVector2{X: a.X.Mul(b.X), Y: a.Y.Mul(b.Y)}
}

// 求两点间距离
func DistanceFV2(a FixVector2, b FixVector2) Fix64 {
	return ((a.X - b.X).Mul(a.X-b.X) + (a.Y - b.Y).Mul(a.Y-b.Y)).Sqrt()
}

// 求两点间距离平方
func DistanceSqrFV2(a FixVector2, b FixVector2) Fix64 {
	return (a.X - b.X).Mul(a.X-b.X) + (a.Y - b.Y).Mul(a.Y-b.Y)
}

// 线性插值
func LerpFV2(a, b FixVector2, t Fix64) FixVector2 {
	return FixVector2{X: a.X + (b.X - a.X).Mul(t), Y: a.Y + (b.Y - a.Y).Mul(t)}
}

//func LerpUnclamped(a, b FixVector2, t fix64) FixVector2 {
//	return LerpFV2(b, a, t)
//}

// 如果需要再加
//Max	返回由两个向量的最大分量组成的向量。
//Min	返回由两个向量的最小分量组成的向量。
//MoveTowards	将点 current 移向 /target/。
//Reflect	从法线定义的向量反射一个向量。
//Scale	将两个向量的分量相乘。
//SignedAngle	返回 from 与 to 之间的有符号角度（以度为单位）。
//SmoothDamp	随时间推移将一个向量逐渐改变为所需目标。
