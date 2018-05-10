package rad

import (
	"log"
	"github.com/galaco/vrad/raytracer"
	"github.com/galaco/vrad/cmd"
	"github.com/galaco/vrad/common/constants"
	"github.com/galaco/vrad/cache"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"github.com/galaco/vrad/vmath/vector"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/vrad/rad/world"
	"github.com/galaco/vrad/rad/clustertable"
)

func Start(config *cmd.Args) {
	log.Println("RadWorld_Start()")

	e := raytracer.GetEnvironment()
	log.Printf("KDTree Size: %d\n", len(e.OptimizedKDTree))
	log.Printf("TriangleIndexList Size: %d\n", len(e.TriangleIndexList))
	log.Printf("OptimisedTriangleList Size: %d\n", len(e.OptimizedTriangleList))
	log.Printf("TriangleMaterials Size: %d\n", len(e.TriangleMaterials))
	log.Printf("TriangleColours Size: %d\n", len(e.TriangleColors))
	log.Printf("LightList Size: %d\n", len(e.LightList))

	if config.LuxelDensity < 1.0 {
		// Remember the old lightmap vectors.
		oldLightmapVecs := [constants.MAX_MAP_TEXINFO][2][4]float32{}
		for i := 0; i < len(cache.GetLumpCache().TexInfo); i++ {
			for j := 0; j < 2; j++ {
				for k := 0; k < 3; k++ {
					oldLightmapVecs[i][j][k] = cache.GetLumpCache().TexInfo[i].LightmapVecsLuxelsPerWorldUnits[j][k]
				}
			}
		}

		// rescale luxels to be no denser than "luxeldensity"
		for i := 0; i < len(cache.GetLumpCache().TexInfo); i++ {
			tx := &cache.GetLumpCache().TexInfo[i]
			for j := 0; j < 2; j++ {
				tmp := mgl32.Vec3{
					tx.LightmapVecsLuxelsPerWorldUnits[j][0],
					tx.LightmapVecsLuxelsPerWorldUnits[j][1],
					tx.LightmapVecsLuxelsPerWorldUnits[j][2],
				}

				// @TODO this is not right!!
				//float scale = VectorNormalize(tmp)
				scale := tmp.Normalize().Len()
				// only rescale them if the current scale is "tighter" than the desired scale
				// FIXME: since this writes out to the BSP file every run, once it's set high it can't be reset
				// to a lower value.
				if  float32(math.Abs(float64(scale))) > config.LuxelDensity {
					if scale < 0 {
						scale = -config.LuxelDensity
					} else {
						scale = config.LuxelDensity
					}
					vector.Scale(&tmp, scale, &tmp)

					tx.LightmapVecsLuxelsPerWorldUnits[j][0] = tmp.X()
					tx.LightmapVecsLuxelsPerWorldUnits[j][1] = tmp.Y()
					tx.LightmapVecsLuxelsPerWorldUnits[j][2] = tmp.Z()
				}
			}
		}

		UpdateAllFaceLightmapExtents()
	}

	clustertable.MakeParents(0, -1)

	clustertable.BuildClusterTable()

	// turn each face into a single patch
	//MakePatches()
	//PairEdges()

	// store the vertex normals calculated in PairEdges
	// so that the can be written to the bsp file for
	// use in the engine
	//SaveVertexNormals()

	// subdivide patches to a maximum dimension
	//SubdividePatches ()

	// add displacement faces to cluster table
	//AddDispsToClusterTable()

	// create directlights out of patches and lights
	//CreateDirectLights ()

	// set up sky cameras
	//ProcessSkyCameras()

}

func UpdateAllFaceLightmapExtents() {
	for i := 0; i < len(*cache.GetTargetFaces()); i++ {
		pFace := &(*cache.GetTargetFaces())[i]

		if (cache.GetLumpCache().TexInfo[pFace.TexInfo].Flags & (flags.SURF_SKY | flags.SURF_NOLIGHT)) != 0 {
			continue		// non-lit texture
		}

		world.CalcFaceExtents(pFace, pFace.LightmapTextureMinsInLuxels, pFace.LightmapTextureSizeInLuxels)
	}
}