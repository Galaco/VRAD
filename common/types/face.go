package types

import "github.com/go-gl/mathgl/mgl32"

type FaceNeighbour struct {
	NumNeighbours int			// neighboring faces that share vertices
	Neighbour []int				// neighboring face list (max of 64)

	Normal []mgl32.Vec3			// adjusted normal per vertex
	FaceNormal mgl32.Vec3		// face normal

	HasDisp bool				// is this surface a displacement surface???
}
