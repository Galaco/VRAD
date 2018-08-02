package trace

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/cache"
)

func PointLeafnum(point *mgl32.Vec3) int {
	return PointLeafnum_r( point, 0 )
}

func PointLeafnum_r(point *mgl32.Vec3, ndxNode int) int {
	// while loop here is to avoid recursion overhead
	for ndxNode >= 0 {
		node := &(cache.GetLumpCache().Nodes[ndxNode])
		plane := &(cache.GetLumpCache().Planes[node.PlaneNum])

		var dist float32
		if plane.AxisType < 3 {
			dist = point[plane.AxisType] - plane.Distance
		} else {
			dist = plane.Normal.Dot(*point) - plane.Distance
		}

		if dist < 0.0 {
			ndxNode = int(node.Children[1])
		} else {
			ndxNode = int(node.Children[0])
		}
	}

	return -1 - ndxNode
}
