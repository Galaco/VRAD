package lightmap

import (
	"log"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/source-tools-common/constants"
	"github.com/galaco/vrad/common/types"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/common/constants/compiler"
	"github.com/galaco/bsp/primitives/vertnormal"
	"github.com/galaco/vrad/vmath/polygon"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/bsp/primitives/leaf"
)

const SMOOTHING_GROUP_HARD_EDGE	= 0xff000000

var vertexRef [constants.MAX_MAP_VERTS]int
var vertexFace [constants.MAX_MAP_VERTS][]int
var faceNeighbour [constants.MAX_MAP_FACES]types.FaceNeighbour
// @TODO --smoothing parameter should override this
var smoothingThreshold = float32(0.7071067) // cos(45.0*(M_PI/180))
var computedNumVertNormalIndices = int(0)
var activeLights *types.DirectLight
var globalAmbient *types.DirectLight
var globalSkyLight *types.DirectLight
var numDirectLights = 0


func PairEdges() {
	var k, m, n int
	var numNeighbours int
	var tmpNeighbour [64]int
	var f *face.Face
	var fn *types.FaceNeighbour

	if compiler.DEBUG == true {
		log.Println("DEBUG: PairingEdges: Counting faces that reference each vertex..")
	}
	// count number of faces that reference each vertex
	for i := 0; i < len(*cache.GetTargetFaces()); i++ {
		f = &((*cache.GetTargetFaces())[i])

		for j := 0; j < int(f.NumEdges); j++ {
			// Store the count in vertexRef
			vertexRef[EdgeVertex(f,j)]++
		}
	}

	if compiler.DEBUG == true {
		log.Println("DEBUG: PairingEdges: Allocate Room")
	}
	// allocate room
	for i := 0; i < len(cache.GetLumpCache().Vertexes); i++ {
		// use the count from above to allocate a big enough array
		vertexFace[i] = make([]int, vertexRef[i])
		// clear the temporary data
		vertexRef[i] = 0
	}

	if compiler.DEBUG == true {
		log.Println("DEBUG: PairingEdges: Store all faces per vertex")
	}
	// store a list of every face that uses a particular vertex
	for i := 0; i < len(*cache.GetTargetFaces()); i++ {
		f = &((*cache.GetTargetFaces())[i])
		for j := 0; j < len(cache.GetLumpCache().Edges); j++ {
			n := EdgeVertex(f, j)

			for k = 0; k < vertexRef[n]; k++ {
				if vertexFace[n][k] == i {
					break
				}
			}
			if k >= vertexRef[n] {
				// add the face to the list
				vertexFace[n][k] = i
				vertexRef[n]++
			}
		}
	}

	if compiler.DEBUG == true {
		log.Println("DEBUG: PairingEdges: Calculate normals")
	}
	// calc normals and set displacement surface flag
	for i := 0; i < len(*cache.GetTargetFaces()); i++ {
		f = &((*cache.GetTargetFaces())[i])
		fn = &faceNeighbour[i]

		// get face normal
		fn.FaceNormal = cache.GetLumpCache().Planes[f.Planenum].Normal

		// set displacement surface flag
		fn.HasDisp = false
		// @TODO Add Displacement support
		if polygon.ValidDispFace(f) == true {
			fn.HasDisp = true
		}
	}

	if compiler.DEBUG == true {
		log.Println("DEBUG: PairingEdges: Find neighbours")
	}
	// find neighbors
	for i := 0; i < len(*cache.GetTargetFaces()); i++ {
		f = &((*cache.GetTargetFaces())[i])

		numNeighbours = 0
		fn = &faceNeighbour[i]

		// allocate room for vertex normals
		fn.Normal = make([]mgl32.Vec3, f.NumEdges)

		// look up all faces sharing vertices and add them to the list
		for j := 0; j < int(f.NumEdges); j++ {
			n = EdgeVertex(f, j)

			for k := 0; k < vertexRef[n]; k++ {
				var cosNormalsAngle float32
				var neighbourNormal *mgl32.Vec3

				// skip self
				if vertexFace[n][k] == i {
					continue
				}

				// if this face doesn't have a displacement -- don't consider displacement neighbors
				if (!fn.HasDisp) && (faceNeighbour[vertexFace[n][k]].HasDisp) {
					continue
				}

				neighbourNormal = &faceNeighbour[vertexFace[n][k]].FaceNormal;
				cosNormalsAngle = (*neighbourNormal).Dot(fn.FaceNormal)

				// add normal if >= threshold or its a displacement surface (this is only if the original
				// face is a displacement)
				if fn.HasDisp {
					// Always smooth with and against a displacement surface.
					fn.Normal[j] = fn.Normal[j].Add(*neighbourNormal)
				} else {
					// No smoothing - use of method (backwards compatibility).
					if ( f.SmoothingGroups == 0 ) && ( (*cache.GetTargetFaces())[vertexFace[n][k]].SmoothingGroups == 0 ) {
						if cosNormalsAngle >= smoothingThreshold {
							fn.Normal[j] = fn.Normal[j].Add(*neighbourNormal)
						} else {
							// not considered a neighbor
							continue
						}
					} else {
						smoothingGroup := f.SmoothingGroups & (*cache.GetTargetFaces())[vertexFace[n][k]].SmoothingGroups

						// Hard edge.
						if ( smoothingGroup & SMOOTHING_GROUP_HARD_EDGE ) != 0 {
							continue
						}

						if smoothingGroup != 0 {
							fn.Normal[j] = fn.Normal[j].Add(*neighbourNormal)
						} else {
							// not considered a neighbor
							continue
						}
					}
				}

				// look to see if we've already added this one
				for m := 0; m < numNeighbours; m++ {
					if tmpNeighbour[m] == vertexFace[n][k] {
						break
					}
				}

				if m >= numNeighbours {
					// add to neighbor list
					tmpNeighbour[m] = vertexFace[n][k]
					numNeighbours++

					// @TODO assert len(tmpNeighbour)  is appropriate
					if numNeighbours > len(tmpNeighbour) {
						log.Fatalf("Stack overflow in neighbors\n")
					}
				}
			}
		}


		if compiler.DEBUG_LEVEL > 7 {
			log.Printf("DEBUG: PairingEdges: Copy neighbours for face: %d\n", i)
		}
		if numNeighbours > 0 {
			// copy over neighbor list
			fn.NumNeighbours = numNeighbours
			fn.Neighbour = make([]int, numNeighbours)
			for m = 0; m < numNeighbours; m++ {
				fn.Neighbour[m] = tmpNeighbour[m]
			}
		}

		if compiler.DEBUG_LEVEL > 7 {
			log.Printf("DEBUG: PairingEdges: Fixup normals for face: %d\n", i)
		}
		// fixup normals
		for j := 0; j < int(f.NumEdges); j++ {
			fn.Normal[j] = fn.Normal[j].Add(fn.FaceNormal)
			fn.Normal[j] = fn.Normal[j].Normalize()
		}
	}
}

func SaveVertexNormals() {
	var fn *types.FaceNeighbour
	var i, j int
	var f *face.Face
	var normalList NormalList


	computedNumVertNormalIndices = 0

	for i = 0; i < len(*cache.GetTargetFaces()); i++ {
		fn = &faceNeighbour[i]
		f = &(*cache.GetTargetFaces())[i]

		for j = 0; j < int(f.NumEdges); j++ {
			var vNormal mgl32.Vec3
			if len(fn.Normal) > 0 {
				vNormal = fn.Normal[j]
			} else {
				// original faces don't have normals
				vNormal = mgl32.Vec3{0,0,0}
			}

			if computedNumVertNormalIndices == constants.MAX_MAP_VERTNORMALINDICES {
				log.Fatal("g_numvertnormalindices == MAX_MAP_VERTNORMALINDICES")
			}

			cache.GetLumpCache().VertNormalIndices[computedNumVertNormalIndices] = uint16(normalList.FindOrAddNormal(&vNormal))
			computedNumVertNormalIndices++
		}
	}

	if len(normalList.Normals) > constants.MAX_MAP_VERTNORMALS {
		log.Fatal("g_numvertnormals > MAX_MAP_VERTNORMALS")
	}

	// Copy the list of unique vert normals into g_vertnormals.
	// @TODO this is a naughty and inefficient way to convert
	// from []Vec3 to []vertNormal.VertNormal, even though they are the same
	// size.
	lazyList := make([]vertnormal.VertNormal, len(normalList.Normals))
	for _,v := range normalList.Normals {
		lazyList = append(lazyList, vertnormal.VertNormal{
			Pos: v,
		})
	}

	cache.GetLumpCache().VertNormals = lazyList
}

func EdgeVertex(f *face.Face, edge int) int {
	if edge < 0 {
 		edge += int(f.NumEdges)
	} else if edge >= int(f.NumEdges) {
		edge = edge % int(f.NumEdges)
	}

	k := cache.GetLumpCache().SurfEdges[int(f.FirstEdge) + edge]
	if k < 0 {
		// Msg("(%d %d) ", dedges[-k].v[1], dedges[-k].v[0] );
		return int(cache.GetLumpCache().Edges[-k][1])
	} else {
		// Msg("(%d %d) ", dedges[k].v[0], dedges[k].v[1] );
		return int(cache.GetLumpCache().Edges[k][0])
	}
}

func BuildVisForLightEnvironment() {
	// Create the vis.
	for iLeaf := 0; iLeaf < len(cache.GetLumpCache().Leafs); iLeaf++ {
		cache.GetLumpCache().Leafs[iLeaf].SetFlags(
			cache.GetLumpCache().Leafs[iLeaf].Flags() & ^(leaf.LEAF_FLAGS_SKY | leaf.LEAF_FLAGS_SKY2D))
		iFirstFace := cache.GetLumpCache().Leafs[iLeaf].FirstLeafFace
		for iLeafFace := uint16(0); iLeafFace < cache.GetLumpCache().Leafs[iLeaf].NumLeafFaces; iLeafFace++ {
			iFace := cache.GetLumpCache().LeafFaces[iFirstFace+iLeafFace]

			tex := &(cache.GetLumpCache().TexInfo[(*cache.GetTargetFaces())[iFace].TexInfo])
			if 0 != tex.Flags & flags.SURF_SKY {
				if 0 != tex.Flags & flags.SURF_SKY2D {
					cache.GetLumpCache().Leafs[iLeaf].SetFlags(
						cache.GetLumpCache().Leafs[iLeaf].Flags() | leaf.LEAF_FLAGS_SKY2D)
				} else {
					cache.GetLumpCache().Leafs[iLeaf].SetFlags(
						cache.GetLumpCache().Leafs[iLeaf].Flags() | leaf.LEAF_FLAGS_SKY)
				}
				MergeDLightVis(globalSkyLight, int(cache.GetLumpCache().Leafs[iLeaf].Cluster))
				MergeDLightVis(globalAmbient, int(cache.GetLumpCache().Leafs[iLeaf].Cluster))
				break
			}
		}
	}

	// Second pass to set flags on leaves that don't contain sky, but touch leaves that
	// contain sky.
	var pvs = make([]byte, constants.MAX_MAP_CLUSTERS / 8)

	nLeafBytes := (len(cache.GetLumpCache().Leafs) >> 3) + 1
	pLeafBits := make([]uint8, nLeafBytes)
	pLeaf2DBits := make([]uint8, nLeafBytes)

	for iLeaf := 0; iLeaf < len(cache.GetLumpCache().Leafs); iLeaf++ {
		// If this leaf has light (3d skybox) in it, then don't bother
		if 0 != cache.GetLumpCache().Leafs[iLeaf].Flags() & leaf.LEAF_FLAGS_SKY {
			continue
		}

		// Don't bother with this leaf if it's solid
		if 0 != cache.GetLumpCache().Leafs[iLeaf].Contents & flags.CONTENTS_SOLID {
			continue
		}

		// See what other leaves are visible from this leaf
		GetVisCache(-1, int(cache.GetLumpCache().Leafs[iLeaf].Cluster), &pvs)

		// Now check out all other leaves
		nByte := iLeaf >> 3
		nBit := 1 << uint(iLeaf & 0x7)
		for iLeaf2 := 0; iLeaf2 < len(cache.GetLumpCache().Leafs); iLeaf2++ {
			if iLeaf2 == iLeaf {
				continue
			}

			if 0 == (cache.GetLumpCache().Leafs[iLeaf2].Flags() & ( leaf.LEAF_FLAGS_SKY | leaf.LEAF_FLAGS_SKY2D )) {
				continue
			}

			// Can this leaf see into the leaf with the sky in it?
			if 0 != PVSCheck(&pvs, int(cache.GetLumpCache().Leafs[iLeaf2].Cluster)) {
				continue
			}

			if 0 != cache.GetLumpCache().Leafs[iLeaf2].Flags() & leaf.LEAF_FLAGS_SKY2D {
				pLeaf2DBits[ nByte ] |= uint8(nBit)
			}
			if 0 != cache.GetLumpCache().Leafs[iLeaf2].Flags() & leaf.LEAF_FLAGS_SKY {
				pLeafBits[ nByte ] |= uint8(nBit)

				// As soon as we know this leaf needs to draw the 3d skybox, we're done
				break
			}
		}
	}

	// Must set the bits in a separate pass so as to not flood-fill LEAF_FLAGS_SKY everywhere
	// pLeafbits is a bit array of all leaves that need to be marked as seeing sky
	for iLeaf := 0; iLeaf < len(cache.GetLumpCache().Leafs); iLeaf++{
		// If this leaf has light (3d skybox) in it, then don't bother
		if 0 != cache.GetLumpCache().Leafs[iLeaf].Flags() & leaf.LEAF_FLAGS_SKY {
			continue
		}

		// Don't bother with this leaf if it's solid
		if 0 != cache.GetLumpCache().Leafs[iLeaf].Contents & flags.CONTENTS_SOLID {
			continue
		}

		// Check to see if this is a 2D skybox leaf
		if 0 != (pLeaf2DBits[ iLeaf >> 3 ] & (1 << uint(iLeaf & 0x7))) {
			cache.GetLumpCache().Leafs[iLeaf].SetFlags(
				cache.GetLumpCache().Leafs[iLeaf].Flags() | leaf.LEAF_FLAGS_SKY2D)
		}

		// If this is a 3D skybox leaf, then we don't care if it was previously a 2D skybox leaf
		if 0 != (pLeafBits[ iLeaf >> 3 ] & (1 << uint(iLeaf & 0x7))) {
			cache.GetLumpCache().Leafs[iLeaf].SetFlags(
				cache.GetLumpCache().Leafs[iLeaf].Flags() | leaf.LEAF_FLAGS_SKY)
			cache.GetLumpCache().Leafs[iLeaf].SetFlags(
				cache.GetLumpCache().Leafs[iLeaf].Flags() & ^leaf.LEAF_FLAGS_SKY2D)
		} else {
			// if radial vis was used on this leaf some of the portals leading
			// to sky may have been culled.  Try tracing to find sky.
			if 0 != cache.GetLumpCache().Leafs[iLeaf].Flags() & leaf.LEAF_FLAGS_RADIAL {
/*				if CanLeafTraceToSky(iLeaf) {
					// FIXME: Should make a version that checks if we hit 2D skyboxes.. oh well.
					cache.GetLumpCache().Leafs[iLeaf].SetFlags(
						cache.GetLumpCache().Leafs[iLeaf].Flags() | leaf.LEAF_FLAGS_SKY)
				}*/
			}
		}
	}
}

func MergeDLightVis(dl *types.DirectLight, cluster int){
	if dl.PVS == nil{
		SetDLightVis(dl, cluster)
	} else {
		var pvs = make([]byte, constants.MAX_MAP_CLUSTERS / 8)
		GetVisCache(-1, cluster, &pvs)

		// merge both vis graphs
		for i := 0; i < int(cache.GetLumpCache().Visibility.NumClusters / 8) + 1; i++ {
			dl.PVS[i] |= pvs[i]
		}
	}
}

func PVSCheck(pvs *[]byte, iCluster int) byte {
	if iCluster >= 0 {
		return (*pvs)[iCluster >> 3] & ( 1 << (uint(iCluster) & 7 ))
	} else {
		// PointInLeaf still returns -1 for valid points sometimes and rather than
		// have black samples, we assume the sample is in the PVS.
		return 1
	}
}

/*
// NOTE: This is just a heuristic.  It traces a finite number of rays to find sky
// NOTE: Full vis is necessary to make this 100% correct.
bool CanLeafTraceToSky( int iLeaf )
{
	// UNDONE: Really want a point inside the leaf here.  Center is a guess, may not be in the leaf
	// UNDONE: Clip this to each plane bounding the leaf to guarantee
	Vector center = vec3_origin;
	for ( int i = 0; i < 3; i++ )
	{
		center[i] = ( (float)(dleafs[iLeaf].mins[i] + dleafs[iLeaf].maxs[i]) ) * 0.5f;
	}

	FourVectors center4, delta;
	fltx4 fractionVisible;
	for ( int j = 0; j < NUMVERTEXNORMALS; j+=4 )
	{
		// search back to see if we can hit a sky brush
		delta.LoadAndSwizzle( g_anorms[j], g_anorms[min( j+1, NUMVERTEXNORMALS-1 )],
			g_anorms[min( j+2, NUMVERTEXNORMALS-1 )], g_anorms[min( j+3, NUMVERTEXNORMALS-1 )] );
		delta *= -MAX_TRACE_LENGTH;
		delta += center4;

		// return true if any hits sky
		TestLine_DoesHitSky ( center4, delta, &fractionVisible );
		if ( TestSignSIMD ( CmpGtSIMD ( fractionVisible, Four_Zeros ) ) )
			return true;
	}

	return false;
}
 */