package vec

import (
	"math"

	"github.com/deminzhang/go-common/mathtable"
)

// 三维向量：(x,y,z)
// Deprecated: 使用泛型版本 vec.Vec3[float32] 替代
type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

// Warn: 浮点相等判定确认应用环境是否合适
func (this *Vector3) Equal(v Vector3) bool {
	return *this == v
}

// 三维向量：设值
func (this *Vector3) Set(x, y, z float32) {
	this.X = x
	this.Y = y
	this.Z = z
}
func (this *Vector3) SetVec(v Vector3) {
	this.X = v.X
	this.Y = v.Y
	this.Z = v.Z
}

// 三维向量：拷贝
func (this *Vector3) Clone() Vector3 {
	return Vector3{X: this.X, Y: this.Y, Z: this.Z}
}

// 三维向量：加上
// this = this + v
func (this *Vector3) Add(v Vector3) {
	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
}

// 三维向量：减去
// this = this - v
func (this *Vector3) Sub(v Vector3) {
	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
}

// 三维向量：数乘
// func (this *Vector3) Mul(scalar float32) {
func (this *Vector3) Multiply(scalar float32) {
	this.X *= scalar
	this.Y *= scalar
	this.Z *= scalar
}

// 三维向量：分向量乘
func (this *Vector3) Scale(v Vector3) {
	this.X *= v.X
	this.Y *= v.Y
	this.Z *= v.Z
}

// 三维向量：数除
func (this *Vector3) Divide(scalar float32) {
	if scalar == 0 {
		//panic("v/0！")
		this.X = float32(math.Inf(1))
		this.Y = float32(math.Inf(1))
		this.Z = float32(math.Inf(1))
		return
	}
	this.X /= scalar
	this.Y /= scalar
	this.Z /= scalar
}

// 三维向量：点积
func (this *Vector3) Dot(v Vector3) float32 {
	return this.X*v.X + this.Y*v.Y + this.Z*v.Z
}

// 三维向量：叉积
func (this *Vector3) Cross(v Vector3) {
	x, y, z := this.X, this.Y, this.Z
	this.X = y*v.Z - z*v.Y
	this.Y = z*v.X - x*v.Z
	this.Z = x*v.Y - y*v.X
}

// 三维向量：长度
func (this *Vector3) Magnitude() float32 {
	return float32(math.Sqrt(float64(this.X*this.X + this.Y*this.Y + this.Z*this.Z)))
}

// 三维向量：长度平方
func (this *Vector3) SqrMagnitude() float32 {
	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}

// 三维向量：单位化 (0向量禁用)
func (this *Vector3) Normalize() {
	l := this.Magnitude()
	if l == 0 {
		this.Set(1, 0, 0)
		return
	}
	this.Divide(l)
}

// 三维向量：单位化值 (0向量禁用)
func (this *Vector3) Normalized() Vector3 {
	v := this.Clone()
	v.Normalize()
	return v
}

// 朝向不变拉长度
func (this *Vector3) ScaleToLength(newLength float32) {
	this.Normalize()
	this.Multiply(newLength)
}

// 复制朝向定长
func (this *Vector3) ScaledToLength(newLength float32) Vector3 {
	v := this.Clone()
	v.ScaleToLength(newLength)
	return v
}

//func (this Vector3) ToProtoVec3f() *comm.Vec3F {
//	return &comm.Vec3F{X: this.X, Y: this.Y, Z: this.Z}
//}
//
//func (this *Vector3) SetProtoVec3f(v *comm.Vec3F) {
//	this.Set(v.X, v.Y, v.Z)
//}

func (this *Vector3) Vector2() Vector2 {
	return Vector2{X: this.X, Y: this.Z}
}

// 绕(X,Z)平面绕center点逆时针旋转 anticlockwise
func (this *Vector3) RotateYAnticlockwise(center Vector3, angle int) Vector3 {
	/*
		假设对图片上任意点(x,y)，绕一个坐标点(rx0,ry0)逆时针旋转a角度后的新的坐标设为(x0, y0)，有公式：
		x0= (x - rx0)*cos(a) - (y - ry0)*sin(a) + rx0 ;
		y0= (x - rx0)*sin(a) + (y - ry0)*cos(a) + ry0 ;
	*/
	sin := mathtable.SinByAngle(angle)
	cos := mathtable.CosByAngle(angle)

	x := (this.X-center.X)*cos - (this.Z-center.Z)*sin + center.X
	z := (this.X-center.X)*sin + (this.Z-center.Z)*cos + center.Z
	return Vector3{X: x, Y: this.Y, Z: z}
}

// 绕(X,Z)平面绕center点顺时针旋转 clockwise 与前端unity一致
func (this *Vector3) RotateYClockwise(center Vector3, angle int) Vector3 {
	return this.RotateYAnticlockwise(center, -angle)
}

// 返回：新向量
func NewVector3(x, y, z float32) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

// 返回：零向量(0,0,0)
func ZeroV3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 0}
}

// X 轴 单位向量
func XAxisV3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 0}
}
func RightV3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 0}
}
func LeftV3() Vector3 {
	return Vector3{X: -1, Y: 0, Z: 0}
}

// Y 轴 单位向量
func YAxisV3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 0}
}
func UpV3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 0}
}
func DownV3() Vector3 {
	return Vector3{X: 0, Y: -1, Z: 0}
}

// Z 轴 单位向量
func ZAxisV3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 1}
}
func ForwardV3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: 1}
}
func BackV3() Vector3 {
	return Vector3{X: 0, Y: 0, Z: -1}
}

func XYAxisV3() Vector3 {
	return Vector3{X: 1, Y: 1, Z: 0}
}
func XZAxisV3() Vector3 {
	return Vector3{X: 1, Y: 0, Z: 1}
}
func YZAxisV3() Vector3 {
	return Vector3{X: 0, Y: 1, Z: 1}
}
func XYZAxisV3() Vector3 {
	return Vector3{X: 1, Y: 1, Z: 1}
}

func PositiveInfinityV3() Vector3 {
	return Vector3{X: float32(math.Inf(1)), Y: float32(math.Inf(1)), Z: float32(math.Inf(1))}
}
func NegativeInfinityV3() Vector3 {
	return Vector3{X: float32(math.Inf(-1)), Y: float32(math.Inf(-1)), Z: float32(math.Inf(-1))}
}

// 返回：a + b 向量
func AddV3(a, b Vector3) Vector3 {
	return Vector3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

// 返回：a - b 向量
func SubV3(a, b Vector3) Vector3 {
	return Vector3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

// 返回：a X b 向量 (X 叉乘)
func CrossV3(a, b Vector3) Vector3 {
	return Vector3{X: a.Y*b.Z - a.Z*b.Y, Y: a.Z*b.X - a.X*b.Z, Z: a.X*b.Y - a.Y*b.X}
}

// ScaleV3 返回：分向量乘
func ScaleV3(a, b Vector3) Vector3 {
	return Vector3{X: a.X * b.X, Y: a.Y * b.Y, Z: a.Z * b.Z}
}

func AddArrayV3(vs []Vector3, dv Vector3) []Vector3 {
	for i, _ := range vs {
		vs[i].Add(dv)
	}
	return vs
}

func MultiplyV3(v Vector3, scalars []float32) []Vector3 {
	vs := []Vector3{}
	for _, value := range scalars {
		vector := v.Clone()
		vector.Multiply(value)
		vs = append(vs, vector)
	}
	return vs
}

// 求两点间距离
func DistanceV3(a, b Vector3) float32 {
	return float32(math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2) + math.Pow(float64(a.Z-b.Z), 2)))
}

// 求两点间距离平方
func DistanceSqrV3(a, b Vector3) float32 {
	return float32(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2) + math.Pow(float64(a.Z-b.Z), 2))
}

// LerpV3 线性插值
func LerpV3(a, b Vector3, t float32) Vector3 {
	return Vector3{X: a.X + (b.X-a.X)*t, Y: a.Y + (b.Y-a.Y)*t, Z: a.X + (b.X-a.X)*t}
}

func LerpUnclampedV3(a, b Vector3, t float32) Vector3 {
	return LerpV3(b, a, 1-t)
}

// Max 返回由两个向量的最大分量组成的向量。
func Max(a, b Vector3) Vector3 {
	x, y, z := a.X, a.Y, a.Z
	if b.X > x {
		x = b.X
	}
	if b.Y > y {
		y = b.Y
	}
	if b.Z > z {
		z = b.Z
	}
	return Vector3{X: x, Y: y, Z: z}
}

// Min	返回由两个向量的最小分量组成的向量。
func Min(a, b Vector3) Vector3 {
	x, y, z := a.X, a.Y, a.Z
	if b.X < x {
		x = b.X
	}
	if b.Y < y {
		y = b.Y
	}
	if b.Z < z {
		z = b.Z
	}
	return Vector3{X: x, Y: y, Z: z}
}
