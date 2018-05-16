package lightmap

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/vmath/vector"
	"log"
)

const numSubDivs = 8

type NormalList struct {
	Normals []mgl32.Vec3

	// This represents a grid from (-1,-1,-1) to (1,1,1).
	normalGrid[numSubDivs][numSubDivs][numSubDivs][]int
}



// Adds the normal if unique. Otherwise, returns the normal's index into m_Normals.
func (list *NormalList) FindOrAddNormal(vNormal *mgl32.Vec3) int {
	var gi [3]int

	// See which grid element it's in.
	for iDim := 0; iDim < 3; iDim++ {
		gi[iDim] = (int)( ((vNormal[iDim] + 1.0) * 0.5) * numSubDivs - 0.000001 )
		gi[iDim] = int(math.Min(float64(gi[iDim]), numSubDivs))
		gi[iDim] = int(math.Max(float64(gi[iDim]), 0))
	}

	// Look for a matching vector in there.
	gridElement := &list.normalGrid[gi[0]][gi[1]][gi[2]]
	for i := 0; i < len(*gridElement); i++ {
		iNormal := (*gridElement)[i]

		pVec := &(list.Normals[iNormal])
		//if( pVec->DistToSqr(vNormal) < 0.00001f )
		if pVec == vNormal {
			return iNormal
		}
	}
	// Ok, add a new one.
	*gridElement = append(*gridElement, len(list.Normals))
	list.Normals = append(list.Normals, *vNormal)
	return len(list.Normals)
}


func GetPhongNormal( faceNum int, spot *mgl32.Vec3, phongNormal *mgl32.Vec3) {
	f := &((*cache.GetTargetFaces())[faceNum])
//	dplane_t	*p = &dplanes[f->planenum];
	var vspot mgl32.Vec3

	faceNormal := (cache.GetLumpCache().Planes[f.Planenum]).Normal
	*phongNormal = faceNormal

	if smoothingThreshold != 1 {
		fn := &(faceNeighbour[faceNum])

		// Calculate modified point normal for surface
		// Use the edge normals iff they are defined.  Bend the surface towards the edge normal(s)
		// Crude first attempt: find nearest edge normal and do a simple interpolation with facenormal.
		// Second attempt: find edge points+center that bound the point and do a three-point triangulation(baricentric)
		// Better third attempt: generate the point normals for all vertices and do baricentric triangulation.

		for j := 0; j < int(f.NumEdges); j++ {
			var v1, v2 mgl32.Vec3
			//int e = dsurfedges[f->firstedge + j];
			//int e1 = dsurfedges[f->firstedge + ((j+f->numedges-1)%f->numedges)];
			//int e2 = dsurfedges[f->firstedge + ((j+1)%f->numedges)];

			//edgeshare_t	*es = &edgeshare[abs(e)];
			//edgeshare_t	*es1 = &edgeshare[abs(e1)];
			//edgeshare_t	*es2 = &edgeshare[abs(e2)];
			// dface_t	*f2;

			n1 := &(fn.Normal[j])
			n2 := &(fn.Normal[(j + 1) % int(f.NumEdges)])

			/*
			  if (VectorCompare( n1, fn->facenormal )
			  && VectorCompare( n2, fn->facenormal) )
			  continue;
			*/

			vert1 := EdgeVertex( f, j );
			vert2 := EdgeVertex( f, j+1 );

			p1 := &((*cache.GetLumpCache()).Vertexes[vert1])
			p2 := &((*cache.GetLumpCache()).Vertexes[vert2])

			// Build vectors from the middle of the face to the edge vertexes and the sample pos.
			v1 = p1.Sub(cache.GetFaceCentroids()[faceNum])
			v2 = p2.Sub(cache.GetFaceCentroids()[faceNum])
			vspot = spot.Sub(cache.GetFaceCentroids()[faceNum])
			aa := v1.Dot(v1)
			bb := v2.Dot(v2)
			ab := v1.Dot(v2)
			a1 := (bb * v1.Dot(vspot) - ab * vspot.Dot(v2)) / (aa * bb - ab * ab)
			a2 := (vspot.Dot(v2)- a1 * ab) / bb

			// Test center to sample vector for inclusion between center to vertex vectors (Use dot product of vectors)
			if  a1 >= 0.0 && a2 >= 0.0 {
				// calculate distance from edge to pos
				var temp mgl32.Vec3

				// Interpolate between the center and edge normals based on sample position
				scale := float32(1.0 - a1 - a2)
				vector.Scale(&fn.FaceNormal, scale, phongNormal)
				vector.Scale(n1, a1, &temp)
				*phongNormal = phongNormal.Add(temp)
				vector.Scale(n2, a2, &temp)
				*phongNormal = phongNormal.Add(temp)

				if 1.0e-20 > phongNormal.Len() {
					log.Fatalf("Phong normal length out of bounds\n")
					//Assert( VectorLength( phongnormal ) > 1.0e-20 );
				}
				*phongNormal = phongNormal.Normalize()

				/*
				  if (a1 > 1 || a2 > 1 || a1 + a2 > 1)
				  {
				  Msg("\n%.2f %.2f\n", a1, a2 );
				  Msg("%.2f %.2f %.2f\n", v1[0], v1[1], v1[2] );
				  Msg("%.2f %.2f %.2f\n", v2[0], v2[1], v2[2] );
				  Msg("%.2f %.2f %.2f\n", vspot[0], vspot[1], vspot[2] );
				  exit(1);

				  a1 = 0;
				  }
				*/
				/*
				  phongnormal[0] = (((j + 1) & 4) != 0) * 255;
				  phongnormal[1] = (((j + 1) & 2) != 0) * 255;
				  phongnormal[2] = (((j + 1) & 1) != 0) * 255;
				*/
				return
			}
		}
	}
}