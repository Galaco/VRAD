package types

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/vrad/vmath/polygon"
)

type Patch struct {
	Winding *polygon.Winding
	Mins, Maxs, FaceMins, FaceMaxs mgl32.Vec3
	Origin mgl32.Vec3	// adjusted off face by face normal
	Plane *plane.Plane	// plane (corrected for facing)
	IterationKey uint16	// Used to prevent touching the same patch multiple times in the same query.
						// See IncrementPatchIterationKey().

	// these are packed into one dword
	// @TODO These WERE a bitfield.
	// Screw that, the performance loss is worth the memory
	// tradeoff here
	NormalMajorAxis uint8	// the major axis of base face normal
	Sky bool
	NeedsBumpMap bool
	Pad uint32

	Normal mgl32.Vec3			// adjusted for phong shading
	PlaneDist float32			// Fixes up patch planes for brush models with an origin brush

	Chop float32				// smallest acceptable width of patch face
	LuxScale float32			// average luxels per world coord
	Scale [2]float32			// Scaling of texture in s & t

	TotalLight BumpLights		// accumulated by radiosity
	// does NOT include light
	// accounted for by direct lighting
	BaseLight mgl32.Vec3			// emissivity only
	BaseArea float32			// surface per area per baselight instance

	DirectLight mgl32.Vec3		// direct light value
	Area float32

	Reflectivity mgl32.Vec3		// Average RGB of texture, modified by material type.

	SampleLight mgl32.Vec3
	SampleArea float32		// for averaging direct light
	FaceNumber int
	ClusterNumber int

	Parent int			// patch index of parent
	Child1 int			// patch index for children
	Child2 int

	NdxNext int					// next patch index in face
	NdxNextParent int				// next parent patch index in face
	NdxNextClusterChild int		// next terminal child index in cluster
	//	struct		patch_s		*next;					// next in face
	//	struct		patch_s		*nextparent;		    // next in face
	//	struct		patch_s		*nextclusterchild;		// next terminal child in cluster

	NumTransfers int
	Transfers *Transfer

	Indices [3]int16				// displacement use these for subdivision
}
