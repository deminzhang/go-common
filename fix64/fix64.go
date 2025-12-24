package fix64

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

type Fix64 int64

const (
	fractionalPlaces = 32
	one              = 1 << fractionalPlaces
	mask             = 1<<fractionalPlaces - 1
	pi               = int64(13493037705)
	pi2              = int64(26986075409) //pi*2
	halfPi           = int64(6746518852)  //pi/2
	LutSize          = halfPi
	FixZero          = Fix64(0)
	FixOne           = Fix64(one)
	FixPi            = Fix64(pi)
	FixPi2           = Fix64(pi2)
	FixHalfPi        = Fix64(halfPi)
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

func NewFix64[T constraints.Integer | constraints.Float](val T) Fix64 {
	t := float64(val) * float64(one)
	if t > 0 {
		t += 0.5
	} else if t < 0 {
		t -= 0.5
	}
	return Fix64(int64(t))
}

func FromInt[T constraints.Integer](val T) Fix64 {
	return Fix64(int64(val) * one)
}

func FromFloat[T constraints.Float](val T) Fix64 {
	t := float64(val) * float64(one)
	if t > 0 {
		t += 0.5
	} else if t < 0 {
		t -= 0.5
	}
	return Fix64(int64(t))
}

func (f Fix64) Int() int {
	return int(int64(f) / one)
}

func (f Fix64) Int32() int32 {
	return int32(int64(f) / one)
}

func (f Fix64) Int64() int64 {
	return int64(f) / one
}

func (f Fix64) Float32() float32 {
	return float32(float64(f) / float64(one))
}

func (f Fix64) Float64() float64 {
	return float64(f) / float64(one)
}

func (f Fix64) Float3() (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%.4f", float64(f)/float64(one)), 64)
}

func (f Fix64) GetHashCode() int64 {
	return int64(f)
}

func (f Fix64) String() string {
	if f >= 0 {
		return fmt.Sprintf("%d:%010d", int64(f>>fractionalPlaces), int64(f&mask))
	}
	ff := -f
	if ff >= 0 {
		return fmt.Sprintf("-%d:%010d", int64(ff>>fractionalPlaces), int64(ff&mask))
	}
	return "-2251799813685248:0000" // The minimum value is -(1<<51).
}

func (f Fix64) Floor() int64 { return int64((f + 0x00000000) >> fractionalPlaces) }

func (f Fix64) Round() int64 { return int64((f + 0x80000000) >> fractionalPlaces) }

func (f Fix64) Ceil() int64 { return int64((f + 0xffffffff) >> fractionalPlaces) }

func (f Fix64) Mul(y Fix64) Fix64 {
	const M, N = 32, 32
	lo, hi := muli64(int64(f), int64(y))
	//fmt.Println("lo, hi", f, y, lo, hi)
	ret := Fix64(hi<<M | lo>>N)
	ret += Fix64((lo >> (N - 1)) & 1) // Round to nearest, instead of rounding down.
	return ret
}

func muli64(u, v int64) (lo, hi uint64) {
	//fmt.Println("muli64", u, v)
	const (
		s    = 32
		mask = 1<<s - 1
	)

	u1 := uint64(u >> s)
	u0 := uint64(u & mask)
	v1 := uint64(v >> s)
	v0 := uint64(v & mask)
	//fmt.Println("u1, u0, v1, v0", u1, u0, v1, v0)
	w0 := u0 * v0
	t := u1*v0 + w0>>s
	w1 := t & mask
	w2 := uint64(int64(t) >> s)
	w1 += u0 * v1
	return uint64(u) * uint64(v), u1*v1 + w2 + uint64(int64(w1)>>s)
}

func (f Fix64) Add(val Fix64) Fix64 {
	return f + val
}

func (f Fix64) Sub(val Fix64) Fix64 {
	return f - val
}

func (f Fix64) Mul1(val Fix64) Fix64 {
	return f * val / FixOne
}

func (f Fix64) Div1(val Fix64) Fix64 {
	return f * FixOne / val
}

func (f Fix64) Div(val Fix64) Fix64 {
	tmp := val
	if tmp == 0 {
		return tmp
	}
	var remainder = uint64(0)
	if f >= 0 {
		remainder = uint64(f)
	} else {
		remainder = uint64(-f)
	}
	var divider = uint64(0)
	if tmp >= 0 {
		divider = uint64(tmp)
	} else {
		divider = uint64(-tmp)
	}
	var quotient uint64 = 0
	var bitPos = int64(fractionalPlaces) + 1
	// If the divider is divisible by 2^n, take advantage of it.
	for {
		if (divider&0xF) == 0 && bitPos >= 4 {
			divider >>= 4
			bitPos -= 4
		} else {
			break
		}
	}

	for {
		if remainder != 0 && bitPos >= 0 {
			shift := countLeadingZeroes(remainder)
			if shift > bitPos {
				shift = bitPos
			}
			remainder <<= shift
			bitPos -= shift

			var div = remainder / divider
			remainder = remainder % divider
			quotient += div << bitPos
			remainder <<= 1
			bitPos -= 1
		} else {
			break
		}
	}
	quotient += 1
	var result = int64(quotient >> 1)
	if (uint64(f)^uint64(tmp))&0x8000000000000000 != 0 {
		result = -result
	}
	return Fix64(result)
}

func countLeadingZeroes(x uint64) int64 {
	var result int64 = 0
	for {
		if (x & 0xF000000000000000) == 0 {
			result += 4
			x <<= 4
		} else {
			break
		}
	}
	for {
		if (x & 0x8000000000000000) == 0 {
			result += 1
			x <<= 1
		} else {
			break
		}
	}
	return result
}

func (f Fix64) Abs() Fix64 {
	if f < 0 {
		return -f
	}
	return f
}

// Sqrt 平方根		numberOfIterations 影响影响计算结果，课调整
func (f Fix64) Sqrt() Fix64 {
	var numberOfIterations = 8
	if f > 0x64000 { //100
		numberOfIterations = 12
	}
	if f > 0x3e8000 { // 1000
		numberOfIterations = 16
	}
	if f > 0x2710000 { //10000
		numberOfIterations = 32
	}
	return sqrt(f, numberOfIterations)
}

func sqrt(f Fix64, numberOfIterations int) Fix64 {
	if f == 0 {
		return FixZero
	}
	var k = f + FixOne>>1
	for i := 0; i < numberOfIterations; i++ {
		k = (k.Add(f.Div(k))) >> 1
	}
	return k
}

func (f Fix64) Sin() Fix64 {
	var ff = f
	var j = FixZero
	for ; ff < FixZero; ff += Fix64(26986075409) {
	}
	if ff > Fix64(26986075409) {
		ff = ff % (Fix64(26986075409))
	}
	var k = ff.Mul(Fix64(360)).Div(Fix64(26986075409))
	if ff != FixZero && ff != Fix64(6746518852) && ff != Fix64(13493037705) && ff != Fix64(20239556557) && ff != Fix64(26986075409) {

		v1 := ff.Mul(Fix64(36000000)).Div(Fix64(26986075409))
		v2 := k.Mul(NewFix64(100000))
		j = v1 - v2

	}

	if k <= Fix64(90) {
		return sinLookup(k, j)
	}
	if k <= Fix64(180) {
		return sinLookup(Fix64(180)-k, -j)
	}
	if k <= Fix64(270) {
		return -sinLookup(k-Fix64(180), j)
	} else {
		return -sinLookup(Fix64(360)-k, -j)
	}
}

func sinLookup(i Fix64, j Fix64) Fix64 {
	//fmt.Println("i, j", i, j, i.Float(), j.Float())
	if j < Fix64(100000) && i < Fix64(90) {
		return Fix64(sinTable[i]) +
			(Fix64(sinTable[i+1]) - Fix64(sinTable[i])).Mul(j).Div(Fix64(100000))
	} else {
		return Fix64(sinTable[i])
	}

}

var sinTable = []int{
	// fix52_12
	//0, 71, 142, 214, 285, 356, 428, 499, 570, 640,
	//711, 781, 851, 921, 990, 1060, 1129, 1197, 1265, 1333,
	//1400, 1467, 1534, 1600, 1665, 1731, 1795, 1859, 1922, 1985,
	//2047, 2109, 2170, 2230, 2290, 2349, 2407, 2465, 2521, 2577,
	//2632, 2687, 2740, 2793, 2845, 2896, 2946, 2995, 3043, 3091,
	//3137, 3183, 3227, 3271, 3313, 3355, 3395, 3435, 3473, 3510,
	//3547, 3582, 3616, 3649, 3681, 3712, 3741, 3770, 3797, 3823,
	//3848, 3872, 3895, 3917, 3937, 3956, 3974, 3991, 4006, 4020,
	//4033, 4045, 4056, 4065, 4073, 4080, 4086, 4090, 4093, 4095,
	//4096
	// fix32_32    math.Sin计算，并装换成fix64
	0, 74957515, 149892197, 224781220, 299601773, 374331065, 448946331, 523424844, 597743917, 671880911,
	745813244, 819518395, 892973913, 966157422, 1039046630, 1111619334, 1183853429, 1255726910, 1327217885, 1398304576,
	1468965330, 1539178623, 1608923068, 1678177418, 1746920580, 1815131613, 1882789739, 1949874349, 2016365009, 2082241464,
	2147483648, 2212071688, 2275985909, 2339206844, 2401715233, 2463492036, 2524518436, 2584775843, 2644245902, 2702910498,
	2760751762, 2817752074, 2873894071, 2929160652, 2983534983, 3037000500, 3089540917, 3141140230, 3191782722, 3241452965,
	3290135830, 3337816489, 3384480416, 3430113397, 3474701533, 3518231241, 3560689261, 3602062661, 3642338838, 3681505524,
	3719550787, 3756463039, 3792231035, 3826843882, 3860291035, 3892562305, 3923647864, 3953538241, 3982224333, 4009697400,
	4035949075, 4060971360, 4084756634, 4107297652, 4128587547, 4148619834, 4167388412, 4184887562, 4201111956, 4216056650,
	4229717092, 4242089121, 4253168970, 4262953261, 4271439016, 4278623649, 4284504972, 4289081193, 4292350918, 4294313152,
	4294967296}

func (f Fix64) Cos() Fix64 {
	return (f + Fix64(6746518852)).Sin()
}

func (f Fix64) Tan() Fix64 {
	return f.Sin().Div(f.Cos())
}
