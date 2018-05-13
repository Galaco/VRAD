package lightmap

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
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