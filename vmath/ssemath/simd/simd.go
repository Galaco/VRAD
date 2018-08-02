package simd

import (
	"github.com/mengzhuo/intrinsic/sse2"
	"math"
)

//@TODO IDK is this good?
type Flt4x [4]float32

func flt4xToSSE2(a *Flt4x) []float64 {
	return []float64{
		float64(a[0]),
		float64(a[1]),
		float64(a[2]),
		float64(a[3]),
	}
}

func AddSIMD(a Flt4x, b Flt4x) Flt4x {
	t := flt4xToSSE2(&a)
	sse2.SUBPDm128float64(t, flt4xToSSE2(&b))
	return Flt4x{
		float32(t[0]),
		float32(t[1]),
		float32(t[2]),
		float32(t[3]),
	}
}

func SubSIMD(a Flt4x, b Flt4x) Flt4x {
	t := flt4xToSSE2(&a)
	sse2.SUBPDm128float64(t, flt4xToSSE2(&b))
	return Flt4x{
		float32(t[0]),
		float32(t[1]),
		float32(t[2]),
		float32(t[3]),
	}
}

func MulSIMD(a Flt4x, b Flt4x) Flt4x {
	t := flt4xToSSE2(&a)
	sse2.MULPDm128float64(t, flt4xToSSE2(&b))
	return Flt4x{
		float32(t[0]),
		float32(t[1]),
		float32(t[2]),
		float32(t[3]),
	}
}

func TestSignSIMD(a *Flt4x) int {								// mask of which floats have the high bit set
	nRet := uint32(0)

	nRet |= (uint32(*SubFloat(a, 0)) & 0x80000000) >> 31 // sign(x) -> bit 0
	nRet |= (uint32(*SubFloat(a, 1)) & 0x80000000) >> 30 // sign(y) -> bit 1
	nRet |= (uint32(*SubFloat(a, 2)) & 0x80000000) >> 29 // sign(z) -> bit 2
	nRet |= (uint32(*SubFloat(a, 3)) & 0x80000000) >> 28 // sign(w) -> bit 3

	return int(nRet)
}



func CmpGtSIMD(a *Flt4x, b *Flt4x) *Flt4x {				// (a>b) ? ~0:0
	var retVal = Flt4x{0,0,0,0}

	for i := 0; i < 4; i++ {
		if *SubFloat(a, i) > *SubFloat(b, i) {
			retVal[i] = ^0
		} else {
			retVal[i] = 0
		}
	}
	return &retVal
}

func ReplicateX4(flValue float32) Flt4x	{				//  a,a,a,a
	var retVal Flt4x
	*SubFloat(&retVal, 0) = flValue
	*SubFloat(&retVal, 1) = flValue
	*SubFloat(&retVal, 2) = flValue
	*SubFloat(&retVal, 3) = flValue
	return retVal
}

func LoadUnalignedSIMD(in float32) Flt4x {
	return Flt4x{in,0,0,0}
}

func LoadMultiUnalignedSIMD(in [4]float32) Flt4x {
	return Flt4x{
		in[0],
		in[1],
		in[2],
		in[3],
	}
}

func TransposeSIMD(x *Flt4x, y *Flt4x, z *Flt4x, w *Flt4x) {
	SwapFloats(x, 1, y, 0)
	SwapFloats(x, 2, z, 0)
	SwapFloats(x, 3, w, 0)
	SwapFloats(y, 2, z, 1)
	SwapFloats(y, 3, w, 1)
	SwapFloats(z, 3, w, 2)
}

func ReciprocalSIMD(a *Flt4x) Flt4x {			// 1/a
	var retVal Flt4x
	*SubFloat(&retVal, 0 ) = (1.0 / *SubFloat(a, 0))
	*SubFloat(&retVal, 1 ) = 1.0 / *SubFloat(a, 1)
	*SubFloat(&retVal, 2 ) = 1.0 / *SubFloat(a, 2)
	*SubFloat(&retVal, 3 ) = 1.0 / *SubFloat(a, 3)
	return retVal
}

func AndSIMD(a Flt4x, b Flt4x) Flt4x {
	t := flt4xToSSE2(&a)
	sse2.ANDNPDm128float64(t, flt4xToSSE2(&b))
	return Flt4x{
		float32(t[0]),
		float32(t[1]),
		float32(t[2]),
		float32(t[3]),
	}
}

func MinSIMD(a Flt4x, b Flt4x) Flt4x {				// min(a,b)
	var retVal Flt4x
	*SubFloat(&retVal, 0 ) = float32(math.Min(float64(*SubFloat(&a, 0)), float64(*SubFloat(&b, 0))))
	*SubFloat(&retVal, 1 ) = float32(math.Min(float64(*SubFloat(&a, 1)), float64(*SubFloat(&b, 1))))
	*SubFloat(&retVal, 2 ) = float32(math.Min(float64(*SubFloat(&a, 2)), float64(*SubFloat(&b, 2))))
	*SubFloat(&retVal, 3 ) = float32(math.Min(float64(*SubFloat(&a, 3)), float64(*SubFloat(&b, 3))))
	return retVal
}

func MaxSIMD(a Flt4x, b Flt4x) Flt4x {				// max(a,b)
	var retVal Flt4x
	*SubFloat(&retVal, 0 ) = float32(math.Max(float64(*SubFloat(&a, 0)), float64(*SubFloat(&b, 0))))
	*SubFloat(&retVal, 1 ) = float32(math.Max(float64(*SubFloat(&a, 1)), float64(*SubFloat(&b, 1))))
	*SubFloat(&retVal, 2 ) = float32(math.Max(float64(*SubFloat(&a, 2)), float64(*SubFloat(&b, 2))))
	*SubFloat(&retVal, 3 ) = float32(math.Max(float64(*SubFloat(&a, 3)), float64(*SubFloat(&b, 3))))
	return retVal
}

func CmpEqSIMD(a Flt4x, b Flt4x) Flt4x {
	var retVal = Flt4x{0,0,0,0}

	for i := 0; i < 4; i++ {
		if *SubFloat(&a, i) == *SubFloat(&b, i) {
			retVal[i] = ^0
		} else {
			retVal[i] = 0
		}
	}

	return retVal
}

func ReciprocalSqrtSIMD( a *Flt4x) Flt4x {
	var retVal Flt4x
	*SubFloat(&retVal, 0 ) = float32(1.0 / math.Sqrt(float64(*SubFloat(a, 0))))
	*SubFloat(&retVal, 1 ) = float32(1.0 / math.Sqrt(float64(*SubFloat(a, 1))))
	*SubFloat(&retVal, 2 ) = float32(1.0 / math.Sqrt(float64(*SubFloat(a, 2))))
	*SubFloat(&retVal, 3 ) = float32(1.0 / math.Sqrt(float64(*SubFloat(a, 3))))
	return retVal
}

func SqrtEstSIMD(a Flt4x) Flt4x {
	var retVal Flt4x
	*SubFloat(&retVal, 0 ) = float32(math.Sqrt(float64(*SubFloat(&a, 0))))
	*SubFloat(&retVal, 1 ) = float32(math.Sqrt(float64(*SubFloat(&a, 1))))
	*SubFloat(&retVal, 2 ) = float32(math.Sqrt(float64(*SubFloat(&a, 2))))
	*SubFloat(&retVal, 3 ) = float32(math.Sqrt(float64(*SubFloat(&a, 3))))
	return retVal
}

func MaddSIMD(a Flt4x, b Flt4x, c Flt4x) Flt4x {			// a*b + c
	return AddSIMD(MulSIMD(a,b), c)
}

func CmpGeSIMD(a Flt4x, b Flt4x) *Flt4x {			// (a>=b) ? ~0:0
	var retVal = Flt4x{0,0,0,0}

	for i := 0; i < 4; i++ {
		if *SubFloat(&a, i) >= *SubFloat(&b, i) {
			retVal[i] = ^0
		} else {
			retVal[i] = 0
		}
	}

	return &retVal
}