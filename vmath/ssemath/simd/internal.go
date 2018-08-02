package simd


func SubFloat(f *Flt4x, index int) *float32 {
	return &(*f)[index]
}

func SwapFloats(a *Flt4x, aIndex int, b *Flt4x, bIndex int) {
	tmp := SubFloat(a, aIndex)
	*SubFloat(a, aIndex) = *SubFloat(b, bIndex)
	*SubFloat(b, bIndex) = *tmp
}
