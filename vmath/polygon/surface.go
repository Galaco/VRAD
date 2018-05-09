package polygon

import "github.com/go-gl/mathgl/mgl32"

func BoxSurfaceArea(boxMin mgl32.Vec3, boxMax mgl32.Vec3) float32{
	boxdim := boxMax.Sub(boxMin)
	return 2.0 * ((boxdim[0]*boxdim[2])+(boxdim[0]*boxdim[1])+(boxdim[1]*boxdim[2]))
}
