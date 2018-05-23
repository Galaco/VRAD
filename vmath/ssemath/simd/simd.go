package simd

import "github.com/mengzhuo/intrinsic/sse2"

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
	sse2.ADDSDm64float64(t, flt4xToSSE2(&b))
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

	nRet |= (uint32(*subFloat(a, 0)) & 0x80000000) >> 31 // sign(x) -> bit 0
	nRet |= (uint32(*subFloat(a, 1)) & 0x80000000) >> 30 // sign(y) -> bit 1
	nRet |= (uint32(*subFloat(a, 2)) & 0x80000000) >> 29 // sign(z) -> bit 2
	nRet |= (uint32(*subFloat(a, 3)) & 0x80000000) >> 28 // sign(w) -> bit 3

	return int(nRet)
}



func CmpGtSIMD(a *Flt4x, b *Flt4x) *Flt4x {				// (a>b) ? ~0:0
	var retVal = Flt4x{0,0,0,0}

	for i := 0; i < 4; i++ {
		if *subFloat(a, i) > *subFloat(b, i) {
			retVal[i] = ^0
		} else {
			retVal[i] = 0
		}
	}
	return &retVal
}

func ReplicateX4(flValue float32) Flt4x	{				//  a,a,a,a
	var retVal Flt4x
	*subFloat(&retVal, 0) = flValue
	*subFloat(&retVal, 1) = flValue
	*subFloat(&retVal, 2) = flValue
	*subFloat(&retVal, 3) = flValue
	return retVal
}

func LoadUnalignedSIMD(in float32) Flt4x {
	return Flt4x{in,0,0,0}
}

func TransposeSIMD(x *Flt4x, y *Flt4x, z *Flt4x, w *Flt4x) {
	swapFloats(x, 1, y, 0)
	swapFloats(x, 2, z, 0)
	swapFloats(x, 3, w, 0)
	swapFloats(y, 2, z, 1)
	swapFloats(y, 3, w, 1)
	swapFloats(z, 3, w, 2)
}