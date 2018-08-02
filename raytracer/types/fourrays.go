package types

import "github.com/galaco/vrad/vmath/ssemath"


// fast SSE-ONLY ray tracing module. Based upon various "real time ray tracing" research.
//#define DEBUG_RAYTRACE 1
type FourRays struct {
	Origin ssemath.FourVectors
	Direction ssemath.FourVectors
}

/*
inline void Check(void) const
{
// in order to be valid to trace as a group, all four rays must have the same signs in all
// of their direction components
#ifndef NDEBUG
for(int c=1;c<4;c++)
{
Assert(direction.X(0)*direction.X(c)>=0);
Assert(direction.Y(0)*direction.Y(c)>=0);
Assert(direction.Z(0)*direction.Z(c)>=0);
}
#endif
}
// returns direction sign mask for 4 rays. returns -1 if the rays can not be traced as a
// bundle.
int CalculateDirectionSignMask(void) const;
*/