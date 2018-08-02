package ssemath

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/vmath/ssemath/simd"
)

/// class FourVectors stores 4 independent vectors for use in SIMD processing. These vectors are
/// stored in the format x x x x y y y y z z z z so that they can be efficiently SIMD-accelerated.
type FourVectors struct {
	X, Y, Z simd.Flt4x
}

/// LoadAndSwizzle - load 4 Vectors into a FourVectors, performing transpose op
func (f *FourVectors) LoadAndSwizzle(a *mgl32.Vec3, b *mgl32.Vec3, c *mgl32.Vec3, d *mgl32.Vec3) {
	// TransposeSIMD has large sub-expressions that the compiler can't eliminate on x360
	// use an unfolded implementation here
	x := simd.LoadUnalignedSIMD(a.X())
	y := simd.LoadUnalignedSIMD(b.X())
	z := simd.LoadUnalignedSIMD(c.X())
	w := simd.LoadUnalignedSIMD(d.X())
	// now, matrix is:
	// x y z ?
	// x y z ?
	// x y z ?
	// x y z ?
	simd.TransposeSIMD(&x, &y, &z, &w)
}

func (f *FourVectors) Add4Vectors(b FourVectors) FourVectors{
	f.X = simd.AddSIMD(f.X, b.X)
	f.Y = simd.AddSIMD(f.Y, b.Y)
	f.Z = simd.AddSIMD(f.Z, b.Z)

	return *f
}

func (f FourVectors) Sub(b FourVectors) FourVectors {
	var retVal = f
	retVal.X = simd.SubSIMD(f.X,b.X)
	retVal.Y = simd.SubSIMD(f.Y,b.Y)
	retVal.Z = simd.SubSIMD(f.Z,b.Z)

	return retVal
}


func (f *FourVectors) MultiplyFV(b FourVectors) *FourVectors {
	f.X = simd.MulSIMD(f.X, b.X)
	f.Y = simd.MulSIMD(f.Y, b.Y)
	f.Z = simd.MulSIMD(f.Z, b.Z)

	return f
}

func (f *FourVectors) MultiplyFlt4x(scale simd.Flt4x) *FourVectors {
	f.X = simd.MulSIMD(f.X, scale)
	f.Y = simd.MulSIMD(f.Y, scale)
	f.Z = simd.MulSIMD(f.Z, scale)

	return f
}

func (f *FourVectors) Multiply(scale float32) FourVectors {
	scalepacked := simd.ReplicateX4(scale)
	f.MultiplyFlt4x(scalepacked)

	return *f
}

func (f *FourVectors) MultiplySelf(b FourVectors) simd.Flt4x {
	dot := simd.MulSIMD(f.X, b.X)
	dot = simd.MaddSIMD(f.Y,b.Y, dot)
	dot = simd.MaddSIMD(f.Z,b.Z, dot)
	return dot
}

/// normalize all 4 vectors in place.
func (f *FourVectors) VectorNormalize() {
	magSq := (*f).MultiplySelf(*f)						// length^2
	f.MultiplyFlt4x(simd.ReciprocalSqrtSIMD(&magSq))	// *(1.0/sqrt(length^2))
}

/// return the approximate length of all 4 vectors. uses the sqrt approximation instruction
func (f FourVectors) Length() simd.Flt4x {
	return simd.SqrtEstSIMD((f.length2()))
}

func (f FourVectors) length2() simd.Flt4x {
	return f.MultiplySelf(f)
}

func (f *FourVectors) DuplicateVector(v *mgl32.Vec3) { //< set all 4 vectors to the same vector value
	f.X = simd.ReplicateX4(v.X())
	f.Y = simd.ReplicateX4(v.Y())
	f.Z = simd.ReplicateX4(v.Z())
}

func (f *FourVectors) Vec(idx int) *mgl32.Vec3 {						//< unpack one of the vectors
	return &mgl32.Vec3{
		f.VecX(idx),
		f.VecY(idx),
		f.VecZ(idx),
	}
}


// X(),Y(),Z() - get at the desired component of the i'th (0..3) vector.
func (f *FourVectors) VecX(idx int) float32 {
// NOTE: if the output goes into a register, this causes a Load-Hit-Store stall (don't mix fpu/vpu math!)
	return *simd.SubFloat(&f.X, idx)
}

func (f *FourVectors) VecY(idx int) float32 {
	return *simd.SubFloat(&f.Y, idx )
}

func (f *FourVectors) VecZ(idx int) float32 {
	return *simd.SubFloat(&f.Z, idx )
}

/*
class ALIGN16 FourVectors
{
public:
	fltx4 x, y, z;

	FORCEINLINE void DuplicateVector(Vector const &v)			//< set all 4 vectors to the same vector value
	{
		x=ReplicateX4(v.x);
		y=ReplicateX4(v.y);
		z=ReplicateX4(v.z);
	}

	FORCEINLINE fltx4 const & operator[](int idx) const
	{
		return *((&x)+idx);
	}

	FORCEINLINE fltx4 & operator[](int idx)
	{
		return *((&x)+idx);
	}

	FORCEINLINE void operator+=(FourVectors const &b)			//< add 4 vectors to another 4 vectors
	{
		x=AddSIMD(x,b.x);
		y=AddSIMD(y,b.y);
		z=AddSIMD(z,b.z);
	}

	FORCEINLINE void operator-=(FourVectors const &b)			//< subtract 4 vectors from another 4
	{
		x=SubSIMD(x,b.x);
		y=SubSIMD(y,b.y);
		z=SubSIMD(z,b.z);
	}

	FORCEINLINE void operator*=(FourVectors const &b)			//< scale all four vectors per component scale
	{
		x=MulSIMD(x,b.x);
		y=MulSIMD(y,b.y);
		z=MulSIMD(z,b.z);
	}

	FORCEINLINE void operator*=(const fltx4 & scale)			//< scale
	{
		x=MulSIMD(x,scale);
		y=MulSIMD(y,scale);
		z=MulSIMD(z,scale);
	}

	FORCEINLINE void operator*=(float scale)					//< uniformly scale all 4 vectors
	{
		fltx4 scalepacked = ReplicateX4(scale);
		*this *= scalepacked;
	}

	FORCEINLINE fltx4 operator*(FourVectors const &b) const		//< 4 dot products
	{
		fltx4 dot=MulSIMD(x,b.x);
		dot=MaddSIMD(y,b.y,dot);
		dot=MaddSIMD(z,b.z,dot);
		return dot;
	}

	FORCEINLINE fltx4 operator*(Vector const &b) const			//< dot product all 4 vectors with 1 vector
	{
		fltx4 dot=MulSIMD(x,ReplicateX4(b.x));
		dot=MaddSIMD(y,ReplicateX4(b.y), dot);
		dot=MaddSIMD(z,ReplicateX4(b.z), dot);
		return dot;
	}

	FORCEINLINE void VProduct(FourVectors const &b)				//< component by component mul
	{
		x=MulSIMD(x,b.x);
		y=MulSIMD(y,b.y);
		z=MulSIMD(z,b.z);
	}
	FORCEINLINE void MakeReciprocal(void)						//< (x,y,z)=(1/x,1/y,1/z)
	{
		x=ReciprocalSIMD(x);
		y=ReciprocalSIMD(y);
		z=ReciprocalSIMD(z);
	}

	FORCEINLINE void MakeReciprocalSaturate(void)				//< (x,y,z)=(1/x,1/y,1/z), 1/0=1.0e23
	{
		x=ReciprocalSaturateSIMD(x);
		y=ReciprocalSaturateSIMD(y);
		z=ReciprocalSaturateSIMD(z);
	}

	// Assume the given matrix is a rotation, and rotate these vectors by it.
	// If you have a long list of FourVectors structures that you all want
	// to rotate by the same matrix, use FourVectors::RotateManyBy() instead.
	inline void RotateBy(const matrix3x4_t& matrix);

	/// You can use this to rotate a long array of FourVectors all by the same
	/// matrix. The first parameter is the head of the array. The second is the
	/// number of vectors to rotate. The third is the matrix.
	static void RotateManyBy(FourVectors * RESTRICT pVectors, unsigned int numVectors, const matrix3x4_t& rotationMatrix );

	/// Assume the vectors are points, and transform them in place by the matrix.
	inline void TransformBy(const matrix3x4_t& matrix);

	/// You can use this to Transform a long array of FourVectors all by the same
	/// matrix. The first parameter is the head of the array. The second is the
	/// number of vectors to rotate. The third is the matrix. The fourth is the
	/// output buffer, which must not overlap the pVectors buffer. This is not
	/// an in-place transformation.
	static void TransformManyBy(FourVectors * RESTRICT pVectors, unsigned int numVectors, const matrix3x4_t& rotationMatrix, FourVectors * RESTRICT pOut );

	/// You can use this to Transform a long array of FourVectors all by the same
	/// matrix. The first parameter is the head of the array. The second is the
	/// number of vectors to rotate. The third is the matrix. The fourth is the
	/// output buffer, which must not overlap the pVectors buffer.
	/// This is an in-place transformation.
	static void TransformManyBy(FourVectors * RESTRICT pVectors, unsigned int numVectors, const matrix3x4_t& rotationMatrix );

	// X(),Y(),Z() - get at the desired component of the i'th (0..3) vector.
	FORCEINLINE const float & X(int idx) const
	{
		// NOTE: if the output goes into a register, this causes a Load-Hit-Store stall (don't mix fpu/vpu math!)
		return SubFloat( (fltx4 &)x, idx );
	}

	FORCEINLINE const float & Y(int idx) const
	{
		return SubFloat( (fltx4 &)y, idx );
	}

	FORCEINLINE const float & Z(int idx) const
	{
		return SubFloat( (fltx4 &)z, idx );
	}

	FORCEINLINE float & X(int idx)
	{
		return SubFloat( x, idx );
	}

	FORCEINLINE float & Y(int idx)
	{
		return SubFloat( y, idx );
	}

	FORCEINLINE float & Z(int idx)
	{
		return SubFloat( z, idx );
	}

	FORCEINLINE Vector Vec(int idx) const						//< unpack one of the vectors
	{
		return Vector( X(idx), Y(idx), Z(idx) );
	}

	FourVectors(void)
	{
	}

	FourVectors( FourVectors const &src )
	{
		x=src.x;
		y=src.y;
		z=src.z;
	}

	FORCEINLINE void operator=( FourVectors const &src )
	{
		x=src.x;
		y=src.y;
		z=src.z;
	}

	/// LoadAndSwizzle - load 4 Vectors into a FourVectors, performing transpose op
	FORCEINLINE void LoadAndSwizzle(Vector const &a, Vector const &b, Vector const &c, Vector const &d)
	{
		// TransposeSIMD has large sub-expressions that the compiler can't eliminate on x360
		// use an unfolded implementation here
#if _X360
		fltx4 tx = LoadUnalignedSIMD( &a.x );
		fltx4 ty = LoadUnalignedSIMD( &b.x );
		fltx4 tz = LoadUnalignedSIMD( &c.x );
		fltx4 tw = LoadUnalignedSIMD( &d.x );
		fltx4 r0 = __vmrghw(tx, tz);
		fltx4 r1 = __vmrghw(ty, tw);
		fltx4 r2 = __vmrglw(tx, tz);
		fltx4 r3 = __vmrglw(ty, tw);

		x = __vmrghw(r0, r1);
		y = __vmrglw(r0, r1);
		z = __vmrghw(r2, r3);
#else
		x		= LoadUnalignedSIMD( &( a.x ));
		y		= LoadUnalignedSIMD( &( b.x ));
		z		= LoadUnalignedSIMD( &( c.x ));
		fltx4 w = LoadUnalignedSIMD( &( d.x ));
		// now, matrix is:
		// x y z ?
		// x y z ?
		// x y z ?
		// x y z ?
		TransposeSIMD(x, y, z, w);
#endif
	}

	/// LoadAndSwizzleAligned - load 4 Vectors into a FourVectors, performing transpose op.
	/// all 4 vectors must be 128 bit boundary
	FORCEINLINE void LoadAndSwizzleAligned(const float *RESTRICT a, const float *RESTRICT b, const float *RESTRICT c, const float *RESTRICT d)
	{
#if _X360
		fltx4 tx = LoadAlignedSIMD(a);
		fltx4 ty = LoadAlignedSIMD(b);
		fltx4 tz = LoadAlignedSIMD(c);
		fltx4 tw = LoadAlignedSIMD(d);
		fltx4 r0 = __vmrghw(tx, tz);
		fltx4 r1 = __vmrghw(ty, tw);
		fltx4 r2 = __vmrglw(tx, tz);
		fltx4 r3 = __vmrglw(ty, tw);

		x = __vmrghw(r0, r1);
		y = __vmrglw(r0, r1);
		z = __vmrghw(r2, r3);
#else
		x		= LoadAlignedSIMD( a );
		y		= LoadAlignedSIMD( b );
		z		= LoadAlignedSIMD( c );
		fltx4 w = LoadAlignedSIMD( d );
		// now, matrix is:
		// x y z ?
		// x y z ?
		// x y z ?
		// x y z ?
		TransposeSIMD( x, y, z, w );
#endif
	}

	FORCEINLINE void LoadAndSwizzleAligned(Vector const &a, Vector const &b, Vector const &c, Vector const &d)
	{
		LoadAndSwizzleAligned( &a.x, &b.x, &c.x, &d.x );
	}

	/// return the squared length of all 4 vectors
	FORCEINLINE fltx4 length2(void) const
	{
		return (*this)*(*this);
	}

	/// return the approximate length of all 4 vectors. uses the sqrt approximation instruction
	FORCEINLINE fltx4 length(void) const
	{
		return SqrtEstSIMD(length2());
	}

	/// normalize all 4 vectors in place. not mega-accurate (uses reciprocal approximation instruction)
	FORCEINLINE void VectorNormalizeFast(void)
	{
		fltx4 mag_sq=(*this)*(*this);						// length^2
		(*this) *= ReciprocalSqrtEstSIMD(mag_sq);			// *(1.0/sqrt(length^2))
	}

	/// normalize all 4 vectors in place.
	FORCEINLINE void VectorNormalize(void)
	{
		fltx4 mag_sq=(*this)*(*this);						// length^2
		(*this) *= ReciprocalSqrtSIMD(mag_sq);				// *(1.0/sqrt(length^2))
	}

	/// construct a FourVectors from 4 separate Vectors
	FORCEINLINE FourVectors(Vector const &a, Vector const &b, Vector const &c, Vector const &d)
	{
		LoadAndSwizzle(a,b,c,d);
	}

	/// construct a FourVectors from 4 separate Vectors
	FORCEINLINE FourVectors(VectorAligned const &a, VectorAligned const &b, VectorAligned const &c, VectorAligned const &d)
	{
		LoadAndSwizzleAligned(a,b,c,d);
	}

	FORCEINLINE fltx4 DistToSqr( FourVectors const &pnt )
	{
		fltx4 fl4dX = SubSIMD( pnt.x, x );
		fltx4 fl4dY = SubSIMD( pnt.y, y );
		fltx4 fl4dZ = SubSIMD( pnt.z, z );
		return AddSIMD( MulSIMD( fl4dX, fl4dX), AddSIMD( MulSIMD( fl4dY, fl4dY ), MulSIMD( fl4dZ, fl4dZ ) ) );

	}

	FORCEINLINE fltx4 TValueOfClosestPointOnLine( FourVectors const &p0, FourVectors const &p1 ) const
	{
		FourVectors lineDelta = p1;
		lineDelta -= p0;
		fltx4 OOlineDirDotlineDir = ReciprocalSIMD( p1 * p1 );
		FourVectors v4OurPnt = *this;
		v4OurPnt -= p0;
		return MulSIMD( OOlineDirDotlineDir, v4OurPnt * lineDelta );
	}

	FORCEINLINE fltx4 DistSqrToLineSegment( FourVectors const &p0, FourVectors const &p1 ) const
	{
		FourVectors lineDelta = p1;
		FourVectors v4OurPnt = *this;
		v4OurPnt -= p0;
		lineDelta -= p0;

		fltx4 OOlineDirDotlineDir = ReciprocalSIMD( lineDelta * lineDelta );

		fltx4 fl4T = MulSIMD( OOlineDirDotlineDir, v4OurPnt * lineDelta );

		fl4T = MinSIMD( fl4T, Four_Ones );
		fl4T = MaxSIMD( fl4T, Four_Zeros );
		lineDelta *= fl4T;
		return v4OurPnt.DistToSqr( lineDelta );
	}
}
*/