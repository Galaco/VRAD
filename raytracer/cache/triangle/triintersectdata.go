package triangle

type TriIntersectData struct {
	FlNx, FlNy, FlNz float32							// plane equation
	FlD float32

	NTriangleID int32									// id of the triangle.

	ProjectedEdgeEquations [6]float32						// A,B,C for each edge equation.  a
	// point is inside the triangle if
	// a*c1+b*c2+c is negative for all 3
	// edges.

	NCoordSelect0, NCoordSelect1 uint8					// the triangle is projected onto a 2d
	// plane for edge testing. These are
	// the indices (0..2) of the
	// coordinates preserved in the
	// projection

	NFlags uint8											// triangle flags
	Unused0 uint8										// no longer used
}
