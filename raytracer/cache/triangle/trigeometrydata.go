package triangle

type TriGeometryData struct {
	NTriangleID int32									// id of the triangle.

	VertexCoordData [9]float32								// can't use a vector in a union

	NFlags uint8											// triangle flags
	// @TODO Signed chars. Is this okay?
	NTmpData0 int8								// used by kd-tree builder
	NTmpData1 int8								// used by kd-tree builder
}

