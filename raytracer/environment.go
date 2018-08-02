package raytracer

import (
	"github.com/galaco/vrad/raytracer/cache"
	"github.com/galaco/vrad/materials/light"
	"github.com/galaco/vrad/raytracer/math"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	math2 "math"
	"github.com/galaco/vrad/common/constants/compiler"
	"github.com/galaco/vrad/raytracer/kdtree"
	"github.com/galaco/vrad/vmath/polygon"
	"github.com/galaco/vrad/raytracer/types"
	"github.com/galaco/vrad/vmath/ssemath/simd"
)

var rayTracerEnvironment Environment

func GetEnvironment() *Environment {
	if &rayTracerEnvironment == nil {
		rayTracerEnvironment = *NewEnvironment()
	}

	return &rayTracerEnvironment
}


type Environment struct {
	Flags uint32
	MinBound mgl32.Vec3
	MaxBound mgl32.Vec3
	BackgroundColour mgl32.Vec4					//< color where no intersection
	OptimizedKDTree []cache.OptimisedKDNode			//< the packed kdtree. root is 0
	OptimizedTriangleList []cache.OptimisedTriangle //< the packed triangles
	TriangleIndexList []int32						//< the list of triangle indices.
	LightList []light.Descriptor					//< the list of lights
	TriangleColors []mgl32.Vec3					//< color of tries
	TriangleMaterials []int							//< material index of tries
}

func (environment *Environment) AddTriangle(id int32, v1 *mgl32.Vec3, v2 *mgl32.Vec3, v3 *mgl32.Vec3, colour *mgl32.Vec3) {
	environment.AddTriangleWithMaterial(id, v1, v2, v3, colour, 0, 0)
}

func (environment *Environment) AddTriangleWithMaterial(id int32, v1 *mgl32.Vec3, v2 *mgl32.Vec3, v3 *mgl32.Vec3,
	colour *mgl32.Vec3, flags uint16, materialIndex int) {
	triangle := cache.OptimisedTriangle{}
	triangle.TriGeometryData.NTriangleID = id
	triangle.TriGeometryData.VertexCoordData[0] = (*v1)[0]
	triangle.TriGeometryData.VertexCoordData[1] = (*v1)[1]
	triangle.TriGeometryData.VertexCoordData[2] = (*v1)[2]
	triangle.TriGeometryData.VertexCoordData[3] = (*v2)[0]
	triangle.TriGeometryData.VertexCoordData[4] = (*v2)[1]
	triangle.TriGeometryData.VertexCoordData[5] = (*v2)[2]
	triangle.TriGeometryData.VertexCoordData[6] = (*v3)[0]
	triangle.TriGeometryData.VertexCoordData[7] = (*v3)[1]
	triangle.TriGeometryData.VertexCoordData[8] = (*v3)[2]
	triangle.TriGeometryData.NFlags = uint8(flags)
	environment.OptimizedTriangleList = append(environment.OptimizedTriangleList, triangle)

	if  0 == (flags & RTE_FLAGS_DONT_STORE_TRIANGLE_COLORS) {
		environment.TriangleColors = append(environment.TriangleColors, *colour)
	}
	if  0 == (flags & RTE_FLAGS_DONT_STORE_TRIANGLE_MATERIALS) {
		environment.TriangleMaterials = append(environment.TriangleMaterials, materialIndex)
	}
	// 	printf("add triangle from (%f %f %f),(%f %f %f),(%f %f %f) id %d\n",
	// 		   XYZ(v1),XYZ(v2),XYZ(v3),id);
}

func (environment *Environment) AddQuad(id int32, v1 *mgl32.Vec3, v2 *mgl32.Vec3, v3 *mgl32.Vec3,
	v4 *mgl32.Vec3, colour *mgl32.Vec3) {
	environment.AddTriangle(id,v1,v2,v3,colour)
	environment.AddTriangle(id+1,v1,v3,v4,colour)
}

func (environment *Environment) AddAxisAlignedRectangularSolid(id int32, minCoord *mgl32.Vec3, maxCoord *mgl32.Vec3,
	colour *mgl32.Vec3) {
	// "far" face
	environment.AddQuad(id,
		&mgl32.Vec3{minCoord[0],maxCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],maxCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],minCoord[1],maxCoord[2]},
		&mgl32.Vec3{minCoord[0],minCoord[1],maxCoord[2]},colour)
	// "near" face
	environment.AddQuad(id,
		&mgl32.Vec3{minCoord[0],maxCoord[1],minCoord[2]},
		&mgl32.Vec3{maxCoord[0],maxCoord[1],minCoord[2]},
		&mgl32.Vec3{maxCoord[0],minCoord[1],minCoord[2]},
		&mgl32.Vec3{minCoord[0],minCoord[1],minCoord[2]},colour)

	// "left" face
	environment.AddQuad(id,
		&mgl32.Vec3{minCoord[0],maxCoord[1],maxCoord[2]},
		&mgl32.Vec3{minCoord[0],maxCoord[1],minCoord[2]},
		&mgl32.Vec3{minCoord[0],minCoord[1],minCoord[2]},
		&mgl32.Vec3{minCoord[0],minCoord[1],maxCoord[2]},colour)
	// "right" face
	environment.AddQuad(id,
		&mgl32.Vec3{maxCoord[0],maxCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],maxCoord[1],minCoord[2]},
		&mgl32.Vec3{maxCoord[0],minCoord[1],minCoord[2]},
		&mgl32.Vec3{maxCoord[0],minCoord[1],maxCoord[2]},colour)

	// "top" face
	environment.AddQuad(id,
		&mgl32.Vec3{minCoord[0],maxCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],maxCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],maxCoord[1],minCoord[2]},
		&mgl32.Vec3{minCoord[0],maxCoord[1],minCoord[2]},colour)
	// "bot" face
	environment.AddQuad(id,
		&mgl32.Vec3{minCoord[0],minCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],minCoord[1],maxCoord[2]},
		&mgl32.Vec3{maxCoord[0],minCoord[1],minCoord[2]},
		&mgl32.Vec3{minCoord[0],minCoord[1],minCoord[2]},colour)
}

func (environment *Environment) SetupAccelerationStructure() {
	root := cache.OptimisedKDNode{}
	environment.OptimizedKDTree = append(environment.OptimizedKDTree, root)
	rootTriangleList := make([]int, len(environment.OptimizedTriangleList))

	for t := 0; t < len(environment.OptimizedTriangleList); t++ {
		rootTriangleList[t] = t
	}

	math.CalculateTriangleListBounds(
		&environment.OptimizedTriangleList,
		&rootTriangleList,
		&environment.MinBound,
		&environment.MaxBound)

	environment.RefineNode(0,&rootTriangleList,len(environment.OptimizedTriangleList), environment.MinBound, environment.MaxBound,0)
	for i := 0; i < len(environment.OptimizedTriangleList); i++ {
		environment.OptimizedTriangleList[i].ChangeIntoIntersectionFormat()
	}
}

func (environment *Environment) Trace4Rays(rays *types.FourRays, TMin simd.Flt4x, TMax simd.Flt4x,
resultOut *types.RayTracingResult,
skipId int, callback *types.ITransparentTriangleCallback) {
	log.Printf("environment: Trace4Rays NOT IMPLEMENTED")

}

func (environment *Environment) Trace4RaysWithoutDirectionSign() {
	log.Printf("environment: Trace4RaysWithoutDirectionSign NOT IMPLEMENTED")

}

func (environment *Environment) ComputeVirtualLightSources() {
	log.Printf("environment: ComputeVirtualLightSources NOT IMPLEMENTED")

}

func (environment *Environment) RenderScene() {
	log.Printf("environment: RenderScene NOT IMPLEMENTED")

}

func (environment *Environment) AddToRayStream() {
	log.Printf("environment: AddToRayStream NOT IMPLEMENTED")

}

func (environment *Environment) FlushStreamEntry() {
	log.Printf("environment: FlushStreamEntry NOT IMPLEMENTED")

}

func (environment *Environment) FinishRayStream() {
	log.Printf("environment: FinishRayStream NOT IMPLEMENTED")

}

func (environment *Environment) MakeLeafNode(firstTri int, lastTri int) {
	log.Printf("environment: MakeLeafNode NOT IMPLEMENTED")
}

func (environment *Environment) CalculateCostsOfSplit(splitPlane int, triangleList []int, numTriangles int, minBound mgl32.Vec3,
	maxBound mgl32.Vec3, splitValue float32, nLeft int, nRight int, nBoth int) float32 {

	// determine the costs of splitting on a given axis, and label triangles with respect to
	// that axis by storing the value in coordselect0. It will also return the number of
	// tris in the left, right, and nboth groups, in order to facilitate memory
	nLeft = 0
	nRight = 0
	nBoth = 0
	// now, label each triangle. Since we have not converted the triangles into
	// intersection fromat yet, we can use the CoordSelect0 field of each as a temp.
	minCoord := 1.0e23
	maxCoord := -1.0e23

	for t := 0; t < numTriangles; t++ {
		tri := &environment.OptimizedTriangleList[triangleList[t]]

		for v := 0; v < 3; v++ {
			minCoord = math2.Min(minCoord, float64(tri.Vertex(v)[splitPlane]))
			maxCoord = math2.Max(maxCoord, float64(tri.Vertex(v)[splitPlane]))
		}
		switch tri.ClassifyAgainstAxisSplit(splitPlane, splitValue) {
		case cache.PLANECHECK_NEGATIVE:
			nLeft++
			tri.TriGeometryData.NTmpData0 = cache.PLANECHECK_NEGATIVE
		case cache.PLANECHECK_POSITIVE:
			nRight++
			tri.TriGeometryData.NTmpData0 = cache.PLANECHECK_POSITIVE
		case cache.PLANECHECK_STRADDLING:
			nBoth++
			tri.TriGeometryData.NTmpData0 = cache.PLANECHECK_STRADDLING
		}
	}

	// @TODO Verify this is correct
	if nLeft != 0 && nBoth == 0 && nRight == 0 {
		splitValue = float32(maxCoord)
	}
	if nRight != 0 && nBoth == 0 && nLeft == 0 {
		splitValue = float32(minCoord)
	}

	leftMins := minBound
	leftMaxs := maxBound
	rightMins := minBound
	rightMaxs := maxBound
	leftMaxs[splitPlane] = splitValue
	rightMins[splitPlane] = splitValue
	SA_L := polygon.BoxSurfaceArea(leftMins, leftMaxs)
	SA_R := polygon.BoxSurfaceArea(rightMins, rightMaxs)
	ISA := 1.0 / polygon.BoxSurfaceArea(minBound, maxBound)
	costOfSplit := kdtree.COST_OF_TRAVERSAL + kdtree.COST_OF_INTERSECTION *
		(float32(nBoth) + (SA_L * ISA * float32(nLeft)) + (SA_R * ISA * float32(nRight)))

	return costOfSplit
}

func (environment *Environment) RefineNode(nodeNumber int, triangleList *[]int, numTris int,
	minBound mgl32.Vec3, maxBound mgl32.Vec3, depth int) {

	if numTris < 3 {											// never split empty lists
		// no point in continuing
		environment.OptimizedKDTree[nodeNumber].Children = KDNODE_STATE_LEAF + (len(environment.TriangleIndexList) << 2)
		environment.OptimizedKDTree[nodeNumber].SetNumberOfTrianglesInLeafNode(numTris)

		// @TODO: Looks to be legacy; properties no longer exist?
		if compiler.DEBUG_RAYTRACE == true {
		//	environment.OptimizedKDTree[nodeNumber].VecMins = minBound
		//	environment.OptimizedKDTree[nodeNumber].VecMaxs = maxBound
		}

		for t := 0; t < numTris; t++ {
			environment.TriangleIndexList = append(environment.TriangleIndexList, int32((*triangleList)[t]))
		}
		return
	}

	bestCost := float32(1.0e23)
	bestNLeft := 0
	bestNRight := 0
	bestNBoth := 0
	bestSplitValue := float32(0.0)
	splitPlane := 0

	triSkip := 1 + (numTris/10)					// don't try all triangles as split
												// points when there are a lot of them
	for axis := 0; axis < 3; axis++ {
		for ts := -1; ts < numTris; ts += triSkip {
			for tv := 0; tv < 3; tv++ {
				var trialNLeft, trialNRight, trialNBoth int
				var trialSplitValue float32

				if ts == -1 {
					trialSplitValue = 0.5 * minBound[axis] + maxBound[axis]
				} else {
					// else, split at the triangle vertex if possible
					tri := &environment.OptimizedTriangleList[(*triangleList)[ts]]
					trialSplitValue = tri.Vertex(tv)[axis]
					if (trialSplitValue > maxBound[axis]) || trialSplitValue < minBound[axis] {
						continue							// don't try this vertex - not inside
					}
				}
				//				printf("ts=%d tv=%d tp=%f\n",ts,tv,trial_splitvalue);
				trialCost := environment.CalculateCostsOfSplit(axis, *triangleList, numTris, minBound, maxBound, trialSplitValue,
					trialNLeft, trialNRight, trialNBoth)
				// 				printf("try %d cost=%f nl=%d nr=%d nb=%d sp=%f\n",axis,trial_cost,trial_nleft,trial_nright, trial_nboth,
				// 					   trial_splitvalue);
				if trialCost < bestCost {
					splitPlane = axis
					bestCost = trialCost
					bestNLeft = trialNLeft
					bestNRight = trialNRight
					bestNBoth = trialNBoth
					bestSplitValue = trialSplitValue
					// save away the axis classification of each triangle
					for t := 0; t <numTris; t++ {
						tri := &environment.OptimizedTriangleList[(*triangleList)[t]]
						tri.TriGeometryData.NTmpData1 = tri.TriGeometryData.NTmpData0
					}
				}

				if ts == -1 {
					break
				}
			}
		}
	}

	// @TODO. commented code below port here
	costOfNoSplit := kdtree.COST_OF_INTERSECTION * numTris
	if (float32(costOfNoSplit) <= bestCost) || NEVER_SPLIT != 0 || (depth > kdtree.MAX_TREE_DEPTH) {
		// no benefit to splitting. just make this a leaf node
		environment.OptimizedKDTree[nodeNumber].Children = KDNODE_STATE_LEAF + (len(environment.TriangleIndexList) << 2)
		environment.OptimizedKDTree[nodeNumber].SetNumberOfTrianglesInLeafNode(numTris)

/*		if compiler.DEBUG_RAYTRACE == true {
			environment.OptimizedKDTree[nodeNumber].VecMins = minBound
			environment.OptimizedKDTree[nodeNumber].VecMaxs = maxBound
		}
*/
		for t:= 0; t < numTris; t++ {
			environment.TriangleIndexList = append(environment.TriangleIndexList, int32((*triangleList)[t]))
		}
	} else {
// 		printf("best split was %d at %f (mid=%f,n=%d, sk=%d)\n",split_plane,best_splitvalue,
// 			   0.5*(MinBound[split_plane]+MaxBound[split_plane]),ntris,tri_skip);
		// its worth splitting!
		// we will achieve the splitting without sorting by using a selection algorithm.
		newTriangleList := make([]int, numTris)

		// now, perform surface area/cost check to determine whether this split was worth it
		leftMins := minBound
		leftMaxs := maxBound
		rightMins := minBound
		rightMaxs := maxBound
		leftMaxs[splitPlane] = bestSplitValue
		rightMins[splitPlane] = bestSplitValue

		nLeftOutput := 0
		nBothOutput := 0
		nRightOutput := 0

		for t := 0; t < numTris; t++ {
			tri := &environment.OptimizedTriangleList[(*triangleList)[t]]
			switch tri.TriGeometryData.NTmpData1 {
			case cache.PLANECHECK_NEGATIVE:
				nLeftOutput++
				//					printf("%d goes left\n",t);
				newTriangleList[nLeftOutput] = (*triangleList)[t]
			case cache.PLANECHECK_POSITIVE:
				nRightOutput++
				//					printf("%d goes right\n",t);
				newTriangleList[numTris - nRightOutput] = (*triangleList)[t]
			case cache.PLANECHECK_STRADDLING:
				//					printf("%d goes both\n",t);
				newTriangleList[bestNLeft + nBothOutput] = (*triangleList)[t]
				nBothOutput++
			}
		}

		leftChild := len(environment.OptimizedKDTree)
		rightChild := leftChild + 1
		// 		printf("node %d split on axis %d at %f, nl=%d nr=%d nb=%d lc=%d rc=%d\n",node_number,
		// 			   split_plane,best_splitvalue,best_nleft,best_nright,best_nboth,
		// 			   left_child,right_child);
		environment.OptimizedKDTree[nodeNumber].Children = splitPlane + (leftChild << 2)
		environment.OptimizedKDTree[nodeNumber].SplittingPlaneValue = bestSplitValue
/*
#ifdef DEBUG_RAYTRACE
		OptimizedKDTree[node_number].vecMins = MinBound;
		OptimizedKDTree[node_number].vecMaxs = MaxBound;
#endif
 */
 		newNode := cache.OptimisedKDNode{}
 		environment.OptimizedKDTree = append(environment.OptimizedKDTree, newNode)
		environment.OptimizedKDTree = append(environment.OptimizedKDTree, newNode)
		// now, recurse!
		if (numTris < 20) && ((bestNLeft == 0) || bestNRight == 0) {
			depth += 100
		}
		environment.RefineNode(leftChild, &newTriangleList, bestNLeft + bestNBoth, leftMins, leftMaxs, depth + 1)
		// @TODO CHECK MY POINTER ARITHMETIC
		tTriangleList := newTriangleList[bestNLeft:]
		environment.RefineNode(rightChild, &tTriangleList, bestNRight + bestNBoth,
			rightMins, rightMaxs, depth + 1)
	}
}

func (environment *Environment) CalculateTriangleListBounds(triangles *[]int, numTriangles int,
	minOut *mgl32.Vec3, maxOut *mgl32.Vec3) {
		minOut = &mgl32.Vec3{ 1.0e23, 1.0e23, 1.0e23}
		maxOut = &mgl32.Vec3{ -1.0e23, -1.0e23, -1.0e23}

		for i := 0; i < numTriangles; i++ {
			tri := &environment.OptimizedTriangleList[(*triangles)[i]]
			for v := 0; v < 3; v++ {
				for c := 0; c < 3; c++ {
					minOut[c] = float32(math2.Min(float64(minOut[c]), float64(tri.Vertex(v)[c])))
					maxOut[c] = float32(math2.Max(float64(maxOut[c]), float64(tri.Vertex(v)[c])))
				}
			}
		}
}

func (environment *Environment) AddInfinitePointlight(position mgl32.Vec3, intensity mgl32.Vec3) {
	log.Printf("environment: AddInfinitePointlight NOT IMPLEMENTED")
}

func (environment *Environment) InitializeFromLoadedBSP() {
	log.Printf("environment: InitializeFromLoadedBSP NOT IMPLEMENTED")
}

func (environment *Environment) AddBSPFace() {
	log.Printf("environment: AddBSPFace NOT IMPLEMENTED")
}

func (environment *Environment) MakeRoomForTriangles(numTriangles int) {
	log.Printf("environment: MakeRoomForTriangles NOT IMPLEMENTED")

}

func (environment *Environment) GetTriangle(triangleId int) *cache.OptimisedTriangle{
	return &environment.OptimizedTriangleList[triangleId]
}

func (environment *Environment) GetTriangleMaterial(triangleId int) int {
	return int(environment.TriangleMaterials[triangleId])
}

func (environment *Environment) GetTriangleColor(triangleId int) *mgl32.Vec3 {
	return &environment.TriangleColors[triangleId]
}

func NewEnvironment() *Environment{
	env := Environment{}
	env.Flags = 0
	env.BackgroundColour = mgl32.Vec4{
		0,0,0,0,
	}

	return &env
}