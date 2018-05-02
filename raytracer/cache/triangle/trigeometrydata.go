package triangle

type TriGeometryData struct {
	NTriangleID int32									// id of the triangle.

	VertexCoordData [9]float32								// can't use a vector in a union

	NFlags uint8											// triangle flags
	NTmpData0 byte								// used by kd-tree builder
	NTmpData1 byte								// used by kd-tree builder
}

