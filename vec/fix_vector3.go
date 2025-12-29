package vec

import (
	"fmt"
	"math"

	. "github.com/deminzhang/go-common/fix64"
)

// / FixVector3
type FixVector3 struct {
	X Fix64
	Y Fix64
	Z Fix64
}

func (this *FixVector3) Equal(v FixVector3) bool {
	//return this.X == v.X && this.Y == v.Y && this.Z == v.Z
	return *this == v
}

func (this *FixVector3) Set(x, y, z Fix64) {
	this.X = x
	this.Y = y
	this.Z = z
}

func (this *FixVector3) SetVec(v FixVector3) {
	this.X = v.X
	this.Y = v.Y
	this.Z = v.Z
}

func (this *FixVector3) Clone() FixVector3 {
	var tmp FixVector3
	tmp.X = this.X
	tmp.Y = this.Y
	tmp.Z = this.Z
	return tmp
}

func (this *FixVector3) Add(b FixVector3) {
	this.X = this.X.Add(b.X)
	this.Y = this.Y.Add(b.Y)
	this.Z = this.Z.Add(b.Z)
}

func (this *FixVector3) Sub(b FixVector3) {
	this.X = this.X.Sub(b.X)
	this.Y = this.Y.Sub(b.Y)
	this.Z = this.Z.Sub(b.Z)
}

func (this *FixVector3) Multiply(b Fix64) {
	this.X = this.X.Mul(b)
	this.Y = this.Y.Mul(b)
	this.Z = this.Z.Mul(b)
}

func (this *FixVector3) Divide(b Fix64) {
	if b == FixZero {
		//panic("v/0！")
		this.X = Fix64(math.Inf(1))
		this.Y = Fix64(math.Inf(1))
		this.Z = Fix64(math.Inf(1))
	}
	this.X = this.X.Div(b)
	this.Y = this.Y.Div(b)
	this.Z = this.Z.Div(b)
}

// 向量：分向量乘
func (this *FixVector3) Scale(v FixVector3) {
	this.X = this.X.Mul(v.X)
	this.Y = this.Y.Mul(v.Y)
	this.Z = this.Z.Mul(v.Z)
}

func (this *FixVector3) Dot(b FixVector3) Fix64 {
	return this.X.Mul(b.X).Add(this.Y.Mul(b.Y)).Add(this.Z.Mul(b.Z))
}

func (this *FixVector3) Cross(b FixVector3) {
	x, y, z := this.X, this.Y, this.Z
	this.X = y.Mul(b.Z).Sub(z.Mul(b.Y))
	this.Y = z.Mul(b.X).Sub(x.Mul(b.Z))
	this.Z = x.Mul(b.Y).Sub(y.Mul(b.X))
}

func (this *FixVector3) Magnitude() Fix64 {
	return (this.SqrMagnitude()).Sqrt()
}

func (this *FixVector3) SqrMagnitude() Fix64 {
	return this.X.Mul(this.X) + this.Y.Mul(this.Y) + this.Z.Mul(this.Z)
}

func (this *FixVector3) Normalize() {
	l := this.Magnitude()
	if l == 0 {
		this.Set(FixOne, 0, 0)
		return
	}
	this.Divide(l)
}

// 三维向量：单位化值 (0向量禁用)
func (this *FixVector3) Normalized() FixVector3 {
	v := this.Clone()
	v.Normalize()
	return v
}

// 朝向不变拉长度
func (this *FixVector3) ScaleToLength(newLength Fix64) {
	this.Normalize()
	this.Multiply(newLength)
}

// 复制朝向定长
func (this *FixVector3) ScaledToLength(newLength Fix64) FixVector3 {
	v := this.Clone()
	v.ScaleToLength(newLength)
	return v
}

func NewFixVector3(x, y, z Fix64) FixVector3 {
	var tmp FixVector3
	tmp.X = x
	tmp.Y = y
	tmp.Z = z
	return tmp
}

// 返回：零向量(0,0,0)
func ZeroFix() FixVector3 {
	return FixVector3{X: 0, Y: 0, Z: 0}
}

// X 轴 单位向量
func XAxisFix() FixVector3 {
	return FixVector3{X: FixOne, Y: 0, Z: 0}
}

func RightFix() FixVector3 {
	return FixVector3{X: FixOne, Y: 0, Z: 0}
}

func LeftFix() FixVector3 {
	return FixVector3{X: NewFix64(-1), Y: 0, Z: 0}
}

// YAxisFix Y轴单位向量
func YAxisFix() FixVector3 {
	return FixVector3{X: 0, Y: FixOne, Z: 0}
}

func UpFix() FixVector3 {
	return FixVector3{X: 0, Y: FixOne, Z: 0}
}

func DownFix() FixVector3 {
	return FixVector3{X: 0, Y: -FixOne, Z: 0}
}

// ZAxisFix Z轴 单位向量
func ZAxisFix() FixVector3 {
	return FixVector3{X: 0, Y: 0, Z: FixOne}
}

func ForwardFix() FixVector3 {
	return FixVector3{X: 0, Y: 0, Z: FixOne}
}

func BackFix() FixVector3 {
	return FixVector3{X: 0, Y: 0, Z: -FixOne}
}

func XYAxisFix() FixVector3 {
	return FixVector3{X: FixOne, Y: FixOne, Z: 0}
}

func XZAxisFix() FixVector3 {
	return FixVector3{X: FixOne, Y: 0, Z: FixOne}
}

func YZAxisFix() FixVector3 {
	return FixVector3{X: 0, Y: FixOne, Z: FixOne}
}

func XYZAxisFix() FixVector3 {
	return FixVector3{X: FixOne, Y: FixOne, Z: FixOne}
}

// func PositiveInfinity() FixVector3 {
// 	return FixVector3{X: NewFix64(math.Inf(1)), Y: NewFix64(math.Inf(1)), Z: NewFix64(math.Inf(1))}
// }

// func NegativeInfinity() FixVector3 {
// 	return FixVector3{X: NewFix64(math.Inf(-1)), Y: NewFix64(math.Inf(-1)), Z: NewFix64(math.Inf(-1))}
// }

func AddFV3(a, b FixVector3) FixVector3 {
	return NewFixVector3(a.X.Add(b.X), a.Y.Add(b.Y), a.Z.Add(b.Z))
}

func SubFV3(a, b FixVector3) FixVector3 {
	return NewFixVector3(a.X.Sub(b.X), a.Y.Sub(b.Y), a.Z.Sub(b.Z))
}

// 返回：a X b 向量 (X 叉乘)
func CrossFV3(a, b FixVector3) FixVector3 {
	return FixVector3{X: (a.Y.Mul(b.Z)).Sub(a.Z.Mul(b.Y)), Y: (a.Z.Mul(b.X)).Sub(a.X.Mul(b.Z)),
		Z: (a.X.Mul(b.Y)).Sub(a.Y.Mul(b.X))}
}

// 返回：分向量乘
func ScaleFV3(a, b FixVector3) FixVector3 {
	return FixVector3{X: a.X.Mul(b.X), Y: a.Y.Mul(b.Y), Z: a.Z.Mul(b.Y)}
}

func NormalizedFV3(a FixVector3) FixVector3 {
	var v FixVector3 = a.Clone()
	v.Normalized()
	return v
}

func DistanceFV3(a, b *FixVector3) Fix64 {
	x := a.X.Sub(b.X)
	y := a.Y.Sub(b.Y)
	z := a.Z.Sub(b.Z)
	return (x.Mul(x).Add(y.Mul(y))).Add(z.Mul(z))
}

func DistanceSqrFV3(a, b FixVector3) Fix64 {
	x := a.X.Sub(b.X)
	y := a.Y.Sub(b.Y)
	z := a.Z.Sub(b.Z)
	return ((x.Mul(x).Add(y.Mul(y))).Add(z.Mul(z))).Sqrt()
}

func LerpFV3(from, to FixVector3, factor Fix64) FixVector3 {
	x := from.X.Add((to.X.Sub(from.X)).Mul(factor))
	y := from.Y.Add((to.Y.Sub(from.Y)).Mul(factor))
	z := from.X.Add((to.X.Sub(from.X)).Mul(factor))
	return NewFixVector3(x, y, z)
}

func (this *FixVector3) ToString() string {
	return fmt.Sprintf("X:%s Y:%s Z:%s", this.X.String(), this.Y.String(), this.Z.String())
}

func (this *FixVector3) Float() string {
	return fmt.Sprintf("X:%v Y:%v Z:%v", this.X.Float64(), this.Y.Float64(), this.Z.Float64())
}

func (this FixVector3) GetHashCode() int64 {
	return this.X.GetHashCode() + this.Y.GetHashCode() + this.Z.GetHashCode()
}
