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
)

const SMOOTHING_GROUP_HARD_EDGE	= 0xff000000

var vertexRef [constants.MAX_MAP_VERTS]int
var vertexFace [constants.MAX_MAP_VERTS][]int
var faceNeighbour [constants.MAX_MAP_FACES]types.FaceNeighbour
// @TODO --smoothing parameter should override this
var smoothingThreshold = float32(0.7071067) // cos(45.0*(M_PI/180))
var computedNumVertNormalIndices = int(0)

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
/*		if ValidDispFace( f ) {
			fn.HasDisp = true
		}
*/
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


		if compiler.DEBUG == true {
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

		if compiler.DEBUG == true {
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