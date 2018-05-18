package clustertable

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/leaf"
	"github.com/galaco/vrad/vmath"
	"github.com/galaco/vrad/cache"
)

func ClusterFromPoint( point *mgl32.Vec3) int16 {
	return PointInLeaf( 0, point ).Cluster
}

func PointInLeaf(iNode int, point *mgl32.Vec3) *leaf.Leaf {
	if iNode < 0 {
		return &(cache.GetLumpCache().Leafs[(-1-iNode)])
	}

	node := &(cache.GetLumpCache().Nodes[iNode])
	plane := &cache.GetLumpCache().Planes[node.PlaneNum]

	dist := point.Dot(plane.Normal) - plane.Distance

	if dist > vmath.TEST_EPSILON {
		return PointInLeaf( int(node.Children[0]), point )
	} else if dist < -vmath.TEST_EPSILON {
		return PointInLeaf( int(node.Children[1]), point )
	} else {
		pTest := PointInLeaf( int(node.Children[0]), point );
		if pTest.Cluster != -1 {
			return pTest
		}

		return PointInLeaf( int(node.Children[1]), point );
	}
}