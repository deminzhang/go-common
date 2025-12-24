package vec

import (
	"math"

	"github.com/deminzhang/go-common/mathtable"
)

type POSKEY int64

type Vector3Int struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	Z int32 `json:"z"`
}

func (this *Vector3Int) Set(x, y, z int32) {
	this.X, this.Y, this.Z = x, y, z
}

func (this Vector3Int) Clone() Vector3Int {
	return Vector3Int{X: this.X, Y: this.Y, Z: this.Z}
}

func (this *Vector3Int) Add(v Vector3Int) {
	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
}

func (this *Vector3Int) Sub(v Vector3Int) {
	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
}
func (this *Vector3Int) SubVal(v Vector3Int) Vector3Int {
	return Vector3Int{
		X: this.X - v.X,
		Y: this.Y - v.Y,
		Z: this.Z - v.Z,
	}
}
func (this *Vector3Int) AddVal(v Vector3Int) Vector3Int {
	return Vector3Int{
		X: this.X + v.X,
		Y: this.Y + v.Y,
		Z: this.Z + v.Z,
	}
}
func (this *Vector3Int) IsZero() bool {
	return this.X == 0 && this.Y == 0 && this.Z == 0
}

func (this *Vector3Int) Equal(v Vector3Int) bool {
	//return this.X == v.X && this.Y == v.Y && this.Z == v.Z
	return *this == v
}

func (this *Vector3Int) ManhattanDistance(v Vector3Int) int32 {
	return int32(math.Abs(float64(this.X-v.X)) + math.Abs(float64(this.Y-v.Y)) + math.Abs(float64(this.Z-v.Z)))
}

func (this *Vector3Int) Key() POSKEY {
	return POSKEY(this.Z)<<32 | POSKEY(this.Y)<<16 | POSKEY(this.X)
}

func (this *Vector3Int) ShortKey() uint32 {
	return uint32(this.Z)<<16 | uint32(this.X)
}

func (this Vector3Int) ToVector2Int() Vector2Int { return Vector2Int{X: this.X, Y: this.Z} }

func (this *Vector3Int) RotateX(angle int) Vector3Int {
	s := mathtable.SinByAngle(angle)
	c := mathtable.CosByAngle(angle)

	x := float32(this.X)*c - float32(this.Z)*s
	z := float32(this.X)*s + float32(this.Z)*c

	return Vector3Int{X: int32(x), Y: this.Y, Z: int32(z)}
}

// ZeroV3Int 返回：零向量(0,0,0)
func ZeroV3Int() Vector3Int {
	return Vector3Int{X: 0, Y: 0, Z: 0}
}
