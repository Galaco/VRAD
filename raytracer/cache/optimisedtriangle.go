package cache

import "github.com/galaco/vrad/raytracer/cache/triangle"

const PLANECHECK_POSITIVE = 1
const PLANECHECK_NEGATIVE = -1
const PLANECHECK_STRADDLING = 0

type OptimisedTriangle struct {
	//@TODO This WAS a union.
	// Cache lines suggest this isn't necessary anymore.
	// Should be verified.
	TriIntersectData triangle.TriIntersectData
	TriGeometryData triangle.TriGeometryData
}

func (t *OptimisedTriangle) ChangeIntoIntersectionFormat() {
	/*
	// lose the vertices and use edge equations instead

	// grab the whole original triangle to we don't overwrite it
	TriGeometryData_t srcTri = m_Data.m_GeometryData;

	m_Data.m_IntersectData.m_nFlags = srcTri.m_nFlags;
	m_Data.m_IntersectData.m_nTriangleID = srcTri.m_nTriangleID;

	Vector p1 = srcTri.Vertex( 0 );
	Vector p2 = srcTri.Vertex( 1 );
	Vector p3 = srcTri.Vertex( 2 );

	Vector e1 = p2 - p1;
	Vector e2 = p3 - p1;

	Vector N = e1.Cross( e2 );
	N.NormalizeInPlace();
	// now, determine which axis to drop
	int drop_axis = 0;
	for(int c=1 ; c<3 ; c++)
	if ( fabs(N[c]) > fabs( N[drop_axis] ) )
	drop_axis = c;

	m_Data.m_IntersectData.m_flD = N.Dot( p1 );
	m_Data.m_IntersectData.m_flNx = N.x;
	m_Data.m_IntersectData.m_flNy = N.y;
	m_Data.m_IntersectData.m_flNz = N.z;

	// decide which axes to keep
	int nCoordSelect0 = ( drop_axis + 1 ) % 3;
	int nCoordSelect1 = ( drop_axis + 2 ) % 3;

	m_Data.m_IntersectData.m_nCoordSelect0 = nCoordSelect0;
	m_Data.m_IntersectData.m_nCoordSelect1 = nCoordSelect1;


	Vector edge1 = GetEdgeEquation( p1, p2, nCoordSelect0, nCoordSelect1, p3 );
	m_Data.m_IntersectData.m_ProjectedEdgeEquations[0] = edge1.x;
	m_Data.m_IntersectData.m_ProjectedEdgeEquations[1] = edge1.y;
	m_Data.m_IntersectData.m_ProjectedEdgeEquations[2] = edge1.z;

	Vector edge2 = GetEdgeEquation( p2, p3, nCoordSelect0, nCoordSelect1, p1 );
	m_Data.m_IntersectData.m_ProjectedEdgeEquations[3] = edge2.x;
	m_Data.m_IntersectData.m_ProjectedEdgeEquations[4] = edge2.y;
	m_Data.m_IntersectData.m_ProjectedEdgeEquations[5] = edge2.z;

	*/
}

func (t *OptimisedTriangle) ClassifyAgainstAxisSplit(splitPlane int, splitValue int) int {
	/*
	// classify a triangle against an axis-aligned plane
	float minc=Vertex(0)[split_plane];
	float maxc=minc;
	for(int v=1;v<3;v++)
	{
		minc=min(minc,Vertex(v)[split_plane]);
		maxc=max(maxc,Vertex(v)[split_plane]);
	}

	if (minc>=split_value)
		return PLANECHECK_POSITIVE;
	if (maxc<=split_value)
		return PLANECHECK_NEGATIVE;
	if (minc==maxc)
		return PLANECHECK_POSITIVE;
	return PLANECHECK_STRADDLING;
	 */
}




