package vec

type Vector2Int struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

func (this *Vector2Int) Equal(v Vector2Int) bool {
	//return this.X == v.X && this.Y == v.Y
	return *this == v
}

func (this *Vector2Int) Clone() Vector2Int {
	return Vector2Int{X: this.X, Y: this.Y}
}

func (this *Vector2Int) Add(v Vector2Int) {
	this.X += v.X
	this.Y += v.Y
}

func (this *Vector2Int) Sub(v Vector2Int) {
	this.X -= v.X
	this.Y -= v.Y
}

func (this *Vector2Int) ToVector3Int() Vector3Int {
	return Vector3Int{this.X, 0, this.Y}
}

func (this *Vector2Int) Key() uint32 {
	return uint32(this.Y)<<16 | uint32(this.X)
}
