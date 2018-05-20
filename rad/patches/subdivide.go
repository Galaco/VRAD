package patches

import (
	"github.com/galaco/vrad/cache"
	"log"
	"github.com/galaco/vrad/common/constants"
	"github.com/galaco/vrad/common/types"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/vmath"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/vrad/vmath/polygon"
	"github.com/galaco/vrad/vmath/vector"
	"math"
	"github.com/galaco/vrad/rad/lightmap"
	"github.com/galaco/vrad/rad/clustertable"
)

// @TODO read this from args
const numBounce = 2
// @TODO replace this with -fast cmd flag
const do_fast = false

var SUbDivideCOunt = 0

func SubdividePatches() {
	var num uint32 //unsigned
	var patch *types.Patch

	if numBounce == 0 {
		return
	}

	uiPatchCount := len(*cache.GetPatches())
	log.Printf("%d patches before subdivision\n", uiPatchCount)

	for i := 0; i < uiPatchCount; i++ {
		pCur := &(*cache.GetPatches())[i]
		pCur.PlaneDist = pCur.Plane.Distance

		/*
		pCur->ndxNextParent = faceParents.Element( pCur->faceNumber );
		faceParents[pCur->faceNumber] = pCur - g_Patches.Base();
		 */
		pCur.NdxNextParent = cache.GetFaceParents()[pCur.FaceNumber]
		// its finding the index of pCur (difference between pCur mem address and [0] mem address in int sizes)
		// @TODO is this correct
		cache.SetFaceParent(pCur.FaceNumber, int(i))
	}

	for i := 0 ; i < uiPatchCount; i++ {
		patch = &((*cache.GetPatches())[i])
		patch.Parent = -1
		if PreventSubdivision(patch) == true {
			continue
		}

		// @TODO Allow FAST compile
		// This true should actually be derived from the -fast cmd flag
		if  false == do_fast {
			if (*cache.GetTargetFaces())[patch.FaceNumber].DispInfo == -1 {
				SubdividePatch(i)
			} else {
				// @TODO Enable displacement support
				//StaticDispMgr()->SubdividePatch( i );
			}
		}
	}

	// fixup next pointers
	for i := 0; i < len(*cache.GetTargetFaces()); i++ {
		cache.GetFacePatches()[i] = constants.CONSTRUCTS_INVALID_INDEX
	}

	uiPatchCount = len(*cache.GetPatches())
	for i := 0; i < uiPatchCount; i++ {
		pCur := &(*cache.GetPatches())[i]
		pCur.NdxNext = cache.GetFacePatches()[pCur.FaceNumber]
		// its finding the index of pCur (difference between pCur mem address and [0] mem address in int sizes)
		// @TODO is this correct
		cache.SetFacePatch(pCur.FaceNumber, i)

/*
#if 0
		var prev *types.Patch
		prev = face_g_Patches[(*cache.GetPatches())[i].FaceNumber]
		(*cache.GetPatches())[i].NdxNext = prev
		face_g_Patches[(*cache.GetPatches())[i].FaceNumber] = &((*cache.GetPatches())[i])
#endif
*/
	}

	// Cache off the leaf number:
	// We have to do this after subdivision because some patches span leaves.
	// (only the faces for model #0 are split by it's BSP which is what governs the PVS, and the leaves we're interested in)
	// Sub models (1-255) are only split for the BSP that their model forms.
	// When those patches are subdivided their origins can end up in a different leaf.
	// The engine will split (clip) those faces at run time to the world BSP because the models
	// are dynamic and can be moved.  In the software renderer, they must be split exactly in order
	// to sort per polygon.
	for i := 0; i < uiPatchCount; i++ {
		(*cache.GetPatches())[i].ClusterNumber = int(clustertable.ClusterFromPoint(&((*cache.GetPatches())[i].Origin)))

		//
		// test for point in solid space (can happen with detail and displacement surfaces)
		//
		if (*cache.GetPatches())[i].ClusterNumber == -1 {
			for j := 0; j < (*cache.GetPatches())[i].Winding.NumPoints; j++ {
				clusterNumber := clustertable.ClusterFromPoint( &(*cache.GetPatches())[i].Winding.Points[j])
				if clusterNumber != -1 {
					(*cache.GetPatches())[i].ClusterNumber = int(clusterNumber)
					break
				}
			}
		}
	}

	// build the list of patches that need to be lit
	for num = 0; num < uint32(uiPatchCount); num++ {
		// do them in reverse order
		i := uint32(uiPatchCount) - num - 1

		// skip patches with children
		pCur := &(*cache.GetPatches())[i]
		if pCur.Child1 == constants.CONSTRUCTS_INVALID_INDEX {
			if pCur.ClusterNumber != - 1 {
				pCur.NdxNextClusterChild = cache.GetClusterChildren()[pCur.ClusterNumber]
				// its finding the index of pCur (difference between pCur mem address and [0] mem address in int sizes)
				// @TODO is this correct (referring to i here...)
				cache.SetClusterChild(pCur.ClusterNumber, int(i))
			}
		}
/*
#if 0
		if g_Patches[i].child1 == g_Patches.InvalidIndex() {
			if g_Patches[i].clusterNumber != -1 {
				g_Patches[i].NextClusterChild = cluster_children[g_Patches[i].clusterNumber];
				cluster_children[g_Patches[i].clusterNumber] = &g_Patches[i];
			}
		}
#endif
*/
	}

	log.Printf("%d patches after subdivision\n", uiPatchCount)
}


//-----------------------------------------------------------------------------
// Purpose: does this surface take/emit light
//-----------------------------------------------------------------------------
func PreventSubdivision( patch *types.Patch ) bool {
	f := &((*cache.GetTargetFaces())[patch.FaceNumber])
	tx := &(cache.GetLumpCache().TexInfo[f.TexInfo])

	if tx.Flags & flags.SURF_NOCHOP != 0 {
		return true
	}

	if (tx.Flags & flags.SURF_NOLIGHT) != 0 && 0 == (tx.Flags & flags.SURF_LIGHT) {
		return true
	}

	return false
}

func SubdividePatch(ndxPatch int) {
	var w, o1, o2 *polygon.Winding
	var patch *types.Patch
	var shouldSubDivide bool
	widest := float32(-1)
	widestAxis := -1

	// get the current patch
	//if len(*cache.GetPatches()) < ndxPatch {
	//	return
	//}
	patch = &((*cache.GetPatches())[ndxPatch])
	if patch == nil {
		return
	}

	// never subdivide sky patches
	if patch.Sky == true {
		return
	}

	w = patch.Winding

	// subdivide along the widest axis
	total := patch.Maxs.Sub(patch.Mins)
	vector.Scale(&total, patch.LuxScale, &total)

	for i := 0; i < 3; i++ {
		if total[i] > widest {
			widestAxis = i
			widest = total[i]
		}

		if (total[i] >= patch.Chop) && (total[i] >= minChop) {
			shouldSubDivide = true
		}
	}

	if (!shouldSubDivide) && widestAxis != -1 {
		// make more square
		if (total[widestAxis] > total[(widestAxis + 1) % 3] * 2) && (total[widestAxis] > total[(widestAxis + 2) % 3] * 2) {
			if patch.Chop > minChop {
				shouldSubDivide = true
				patch.Chop = float32(math.Max( float64(minChop), float64(patch.Chop / 2)))
			}
		}
	}

	if false == shouldSubDivide {
		return
	}

	// split the winding
	split := mgl32.Vec3{0,0,0}
	split[widestAxis] = 1
	dist := (patch.Mins[widestAxis] + patch.Maxs[widestAxis]) * 0.5
	ClipWindingEpsilon(w, &split, dist, vmath.ON_EPSILON, &o1, &o2)

	// calculate the area of the patches to see if they are "significant"
	var center1, center2 mgl32.Vec3
	area1 := WindingAreaAndBalancePoint( o1, &center1 )
	area2 := WindingAreaAndBalancePoint( o2, &center2 )

	if area1 == 0 || area2 == 0 {
		log.Printf("zero area child patch\n")
		return
	}

	// create new child patches
	ndxChild1Patch := CreateChildPatch(ndxPatch, o1, area1, &center1)
	ndxChild2Patch := CreateChildPatch(ndxPatch, o2, area2, &center2)

	SUbDivideCOunt++
	cpt := SUbDivideCOunt
	if cpt < 0 {

	}

	// FIXME: This could go into CreateChildPatch if child1, child2 were stored in the patch as child[0], child[1]
	patch = &((*cache.GetPatches())[ndxPatch])
	patch.Child1 = ndxChild1Patch
	patch.Child2 = ndxChild2Patch

	SubdividePatch(ndxChild1Patch)
	SubdividePatch(ndxChild2Patch)
}

func ClipWindingEpsilon(in *polygon.Winding, normal *mgl32.Vec3, dist float32,
				epsilon float32, front **polygon.Winding, back **polygon.Winding) {
	var	dists [constants.MAX_POINTS_ON_WINDING + 4]float32
	var	sides [constants.MAX_POINTS_ON_WINDING + 4]int
	var	counts [3]int
	var dot float32
	var i, j, maxpts int
	var mid = mgl32.Vec3{0,0,0}
	var b,f *polygon.Winding

// determine sides for each point
	for i = 0; i < in.NumPoints; i++ {
		dot = in.Points[i].Dot(*normal)
		dot -= dist
		dists[i] = dot
		if dot > epsilon {
			sides[i] = polygon.SIDE_FRONT
		} else if dot < -epsilon {
			sides[i] = polygon.SIDE_BACK
		} else {
			sides[i] = polygon.SIDE_ON
		}
		counts[sides[i]]++
	}
	sides[i] = sides[0]
	dists[i] = dists[0]

	*back = nil
	*front = nil

	if 0 == counts[0] {
		*back = polygon.CopyWinding(in)
		return
	}
	if 0 == counts[1] {
		*front = polygon.CopyWinding(in)
		return
	}

	maxpts = in.NumPoints + 4	// cant use counts[0]+2 because
								// of fp grouping errors

	f = polygon.NewWinding(maxpts)
	*front = f
	b = polygon.NewWinding(maxpts)
	*back = b

	for i = 0; i < in.NumPoints ; i++ {
		p1 := &in.Points[i]

		if sides[i] == polygon.SIDE_ON {
			f.Points[f.NumPoints] = *p1
			f.NumPoints++
			b.Points[b.NumPoints] = *p1
			b.NumPoints++
			continue
		}

		if sides[i] == polygon.SIDE_FRONT {
			f.Points[f.NumPoints] = *p1
			f.NumPoints++
		}
		if sides[i] == polygon.SIDE_BACK {
			b.Points[b.NumPoints] = *p1
			b.NumPoints++
		}

		if sides[i + 1] == polygon.SIDE_ON || sides[i + 1] == sides[i] {
			continue
		}

	// generate a split point
		p2 := &in.Points[(i + 1) % in.NumPoints]

		dot = dists[i] / (dists[i] - dists[i+1])
		for j = 0; j < 3 ; j++ {	// avoid round off error when possible
			if normal[j] == 1 {
				mid[j] = dist
			} else if normal[j] == -1 {
				mid[j] = -dist
			} else {
				mid[j] = p1[j] + dot * (p2[j] - p1[j])
			}
		}

		f.Points[f.NumPoints] = mid
		f.NumPoints++
		b.Points[b.NumPoints] = mid
		b.NumPoints++
	}

	if f.NumPoints > maxpts || b.NumPoints > maxpts {
		log.Fatal ("ClipWinding: points exceeded estimate");
	}
	if f.NumPoints > constants.MAX_POINTS_ON_WINDING || b.NumPoints > constants.MAX_POINTS_ON_WINDING {
		log.Fatal("ClipWinding: MAX_POINTS_ON_WINDING")
	}
}


//-----------------------------------------------------------------------------
// Purpose: subdivide the "parent" patch
//-----------------------------------------------------------------------------
func CreateChildPatch(nParentIndex int, winding *polygon.Winding, flArea float32, vecCenter *mgl32.Vec3) int {
	parent := &((*cache.GetPatches())[nParentIndex])

	// copy all elements of parent patch to children
	child := *parent

	//NOTE: Copying the parent may be better than creating child and copying contents ( see *child = *parent)
	cache.AddPatchToCache(&child)
	nChildIndex := len(*cache.GetPatches())// - 1

	// Set up links
	child.NdxNext = constants.CONSTRUCTS_INVALID_INDEX
	child.NdxNextParent = constants.CONSTRUCTS_INVALID_INDEX
	child.NdxNextClusterChild = constants.CONSTRUCTS_INVALID_INDEX
	child.Child1 = constants.CONSTRUCTS_INVALID_INDEX
	child.Child2 = constants.CONSTRUCTS_INVALID_INDEX
	child.Parent = nParentIndex
	child.IterationKey = 0

	child.Winding = winding
	child.Area = flArea

	child.Origin = *vecCenter
	//@TODO Displacement disabled
	// Apparently this condition shouldnt be met anyway
	// But WE dont know that yet
	if polygon.ValidDispFace(&(*cache.GetTargetFaces())[child.FaceNumber]) == true {
		// shouldn't get here anymore!!
		log.Printf("SubdividePatch: Error - Should not be here!\n")
		//StaticDispMgr()->GetDispSurfNormal( child.FaceNumber, child.Origin, child.Normal, true )
	} else {
		lightmap.GetPhongNormal(child.FaceNumber, &child.Origin, &child.Normal)
	}

	child.PlaneDist = child.Plane.Distance
	polygon.WindingBounds(child.Winding, &child.Mins, &child.Maxs)

	if false == child.BaseLight.ApproxEqual(mgl32.Vec3{0,0,0}) {
		// don't check edges on surf lights
		return nChildIndex
	}

	// Subdivide patch towards minchop if on the edge of the face
	total := child.Maxs.Sub(child.Mins)
	vector.Scale(&total, child.LuxScale, &total)
	if child.Chop > minChop && (total[0] < child.Chop) && (total[1] < child.Chop) && (total[2] < child.Chop) {
		for i := 0; i < 3; i++ {
			if (child.FaceMaxs[i] == child.Maxs[i] || child.FaceMins[i] == child.Mins[i] ) && total[i] > minChop {
				child.Chop = float32(math.Max( float64(minChop), float64(child.Chop / 2)))
				break
			}
		}
	}

	return nChildIndex
}


func WindingAreaAndBalancePoint(w *polygon.Winding, center *mgl32.Vec3) float32 {
	var	d1, d2, cross mgl32.Vec3
	var total float32

	center = &mgl32.Vec3{0,0,0}
	if w == nil {
		return 0.0
	}

	total = 0
	for i := 2 ; i < int(w.NumPoints); i++ {
		d1 = w.Points[i-1].Sub(w.Points[0])
		d2 = w.Points[i].Sub(w.Points[0])
		cross = d1.Cross(d2)
		area := cross.Len()
		total += area

		// center of triangle, weighed by area
		vector.MA( center, area / 3.0, &w.Points[i-1], center )
		vector.MA( center, area / 3.0, &w.Points[i], center )
		vector.MA( center, area / 3.0, &w.Points[0], center )
	}
	if total != 0 {
		vector.Scale(center, 1.0 / total, center)
	}
	return total * 0.5
}