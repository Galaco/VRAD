package raytracer

import (
	"github.com/galaco/vrad/raytracer/cache"
	"github.com/galaco/vrad/materials/light"
	"github.com/galaco/vrad/raytracer/math"
	"github.com/go-gl/mathgl/mgl32"
)

const RTE_FLAGS_FAST_TREE_GENERATION = 1
const RTE_FLAGS_DONT_STORE_TRIANGLE_COLORS = 2				// saves memory if not needed
const RTE_FLAGS_DONT_STORE_TRIANGLE_MATERIALS = 4

const TRACE_ID_SKY        = 0x01000000  // sky face ray blocker
const TRACE_ID_OPAQUE     = 0x02000000  // everyday light blocking face
const TRACE_ID_STATICPROP = 0x04000000  // static prop - lower bits are prop ID

var rayTracerEnvironment Environment

func GetEnvironment() *Environment {
	if &rayTracerEnvironment == nil {
		rayTracerEnvironment = Environment{}
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
	rootTriangleList := []int{}

	for t := 0; t < len(environment.OptimizedTriangleList); t++ {
		rootTriangleList[t] = t
	}

	math.CalculateTriangleListBounds(
		&environment.OptimizedTriangleList,
		&rootTriangleList,
		&environment.MinBound,
		&environment.MaxBound)

	//RefineNode(0,&rootTriangleList,len(environment.OptimizedTriangleList),&environment.MinBound,&environment.MaxBound,0)
	for i := 0; i < len(environment.OptimizedTriangleList); i++ {
		environment.OptimizedTriangleList[i].ChangeIntoIntersectionFormat()
	}
}

func (environment *Environment) Trace4Rays() {

}

func (environment *Environment) Trace4RaysWithoutDirectionSign() {

}

func (environment *Environment) ComputeVirtualLightSources() {

}

func (environment *Environment) RenderScene() {

}

func (environment *Environment) AddToRayStream() {

}

func (environment *Environment) FlushStreamEntry() {

}

func (environment *Environment) FinishRayStream() {

}

func (environment *Environment) MakeLeafNode(firstTri int, lastTri int) {

}

func (environment *Environment) CalculateCostsOfSplit() {

}

func (environment *Environment) RefineNode(nodeNumber int, triangleList *[]int, numTris int,
	MinBound mgl32.Vec3, MaxBound mgl32.Vec3, depth int) {

}

func (environment *Environment) CalculateTriangleListBounds(triangles *[]int, numTriangles int,
	minOut *mgl32.Vec3, maxOut *mgl32.Vec3) {
/**
if (ntris<3)											// never split empty lists
	{
		// no point in continuing
		OptimizedKDTree[node_number].Children=KDNODE_STATE_LEAF+(TriangleIndexList.Count()<<2);
		OptimizedKDTree[node_number].SetNumberOfTrianglesInLeafNode(ntris);

#ifdef DEBUG_RAYTRACE
		OptimizedKDTree[node_number].vecMins = MinBound;
		OptimizedKDTree[node_number].vecMaxs = MaxBound;
#endif

		for(int t=0;t<ntris;t++)
			TriangleIndexList.AddToTail(tri_list[t]);
		return;
	}

	float best_cost=1.0e23;
	int best_nleft=0,best_nright=0,best_nboth=0;
	float best_splitvalue=0;
	int split_plane=0;

	int tri_skip=1+(ntris/10);								// don't try all trinagles as split
															// points when there are a lot of them
	for(int axis=0;axis<3;axis++)
	{
		for(int ts=-1;ts<ntris;ts+=tri_skip)
		{
			for(int tv=0;tv<3;tv++)
			{
				int trial_nleft,trial_nright,trial_nboth;
				float trial_splitvalue;
				if (ts==-1)
					trial_splitvalue=0.5*(MinBound[axis]+MaxBound[axis]);
				else
				{
					// else, split at the triangle vertex if possible
					CacheOptimizedTriangle &tri=OptimizedTriangleList[tri_list[ts]];
					trial_splitvalue = tri.Vertex(tv)[axis];
					if ((trial_splitvalue>MaxBound[axis]) || (trial_splitvalue<MinBound[axis]))
						continue;							// don't try this vertex - not inside

				}
//				printf("ts=%d tv=%d tp=%f\n",ts,tv,trial_splitvalue);
				float trial_cost=
					CalculateCostsOfSplit(axis,tri_list,ntris,MinBound,MaxBound,trial_splitvalue,
										  trial_nleft,trial_nright, trial_nboth);
// 				printf("try %d cost=%f nl=%d nr=%d nb=%d sp=%f\n",axis,trial_cost,trial_nleft,trial_nright, trial_nboth,
// 					   trial_splitvalue);
				if (trial_cost<best_cost)
				{
					split_plane=axis;
					best_cost=trial_cost;
					best_nleft=trial_nleft;
					best_nright=trial_nright;
					best_nboth=trial_nboth;
					best_splitvalue=trial_splitvalue;
					// save away the axis classification of each triangle
					for(int t=0 ; t < ntris; t++)
					{
						CacheOptimizedTriangle &tri=OptimizedTriangleList[tri_list[t]];
						tri.m_Data.m_GeometryData.m_nTmpData1 = tri.m_Data.m_GeometryData.m_nTmpData0;
					}
				}
				if (ts==-1)
					break;
			}
		}

	}
	float cost_of_no_split=COST_OF_INTERSECTION*ntris;
	if ( (cost_of_no_split<=best_cost) || NEVER_SPLIT || (depth>MAX_TREE_DEPTH))
	{
		// no benefit to splitting. just make this a leaf node
		OptimizedKDTree[node_number].Children=KDNODE_STATE_LEAF+(TriangleIndexList.Count()<<2);
		OptimizedKDTree[node_number].SetNumberOfTrianglesInLeafNode(ntris);
#ifdef DEBUG_RAYTRACE
		OptimizedKDTree[node_number].vecMins = MinBound;
		OptimizedKDTree[node_number].vecMaxs = MaxBound;
#endif
		for(int t=0;t<ntris;t++)
			TriangleIndexList.AddToTail(tri_list[t]);
	}
	else
	{
// 		printf("best split was %d at %f (mid=%f,n=%d, sk=%d)\n",split_plane,best_splitvalue,
// 			   0.5*(MinBound[split_plane]+MaxBound[split_plane]),ntris,tri_skip);
		// its worth splitting!
		// we will achieve the splitting without sorting by using a selection algorithm.
		int32 *new_triangle_list;
		new_triangle_list=new int32[ntris];

		// now, perform surface area/cost check to determine whether this split was worth it
		Vector LeftMins=MinBound;
		Vector LeftMaxes=MaxBound;
		Vector RightMins=MinBound;
		Vector RightMaxes=MaxBound;
		LeftMaxes[split_plane]=best_splitvalue;
		RightMins[split_plane]=best_splitvalue;

		int n_left_output=0;
		int n_both_output=0;
		int n_right_output=0;
		for(int t=0;t<ntris;t++)
		{
			CacheOptimizedTriangle &tri=OptimizedTriangleList[tri_list[t]];
			switch( tri.m_Data.m_GeometryData.m_nTmpData1 )
			{
				case PLANECHECK_NEGATIVE:
//					printf("%d goes left\n",t);
					new_triangle_list[n_left_output++]=tri_list[t];
					break;
				case PLANECHECK_POSITIVE:
					n_right_output++;
//					printf("%d goes right\n",t);
					new_triangle_list[ntris-n_right_output]=tri_list[t];
					break;
				case PLANECHECK_STRADDLING:
//					printf("%d goes both\n",t);
					new_triangle_list[best_nleft+n_both_output]=tri_list[t];
					n_both_output++;
					break;


			}
		}
		int left_child=OptimizedKDTree.Count();
		int right_child=left_child+1;
// 		printf("node %d split on axis %d at %f, nl=%d nr=%d nb=%d lc=%d rc=%d\n",node_number,
// 			   split_plane,best_splitvalue,best_nleft,best_nright,best_nboth,
// 			   left_child,right_child);
		OptimizedKDTree[node_number].Children=split_plane+(left_child<<2);
		OptimizedKDTree[node_number].SplittingPlaneValue=best_splitvalue;
#ifdef DEBUG_RAYTRACE
		OptimizedKDTree[node_number].vecMins = MinBound;
		OptimizedKDTree[node_number].vecMaxs = MaxBound;
#endif
		CacheOptimizedKDNode newnode;
		OptimizedKDTree.AddToTail(newnode);
		OptimizedKDTree.AddToTail(newnode);
		// now, recurse!
		if ( (ntris<20) && ((best_nleft==0) || (best_nright==0)) )
			depth+=100;
		RefineNode(left_child,new_triangle_list,best_nleft+best_nboth,LeftMins,LeftMaxes,depth+1);
		RefineNode(right_child,new_triangle_list+best_nleft,best_nright+best_nboth,
				   RightMins,RightMaxes,depth+1);
		delete[] new_triangle_list;
	}
 */
}

func (environment *Environment) AddInfinitePointlight(position mgl32.Vec3, intensity mgl32.Vec3) {

}

func (environment *Environment) InitializeFromLoadedBSP() {

}

func (environment *Environment) AddBSPFace() {

}

func (environment *Environment) MakeRoomForTriangles(numTriangles int) {

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