package simd




func subFloat(f *Flt4x, index int) *float32 {
	return &(*f)[index]
}

func swapFloats(a *Flt4x, aIndex int, b *Flt4x, bIndex int) {
	tmp := subFloat(a, aIndex)
	*subFloat(a, aIndex) = *subFloat(b, bIndex)
	*subFloat(b, bIndex) = *tmp
}
