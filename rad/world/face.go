package world

import (
	"github.com/galaco/bsp/primitives/face"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/common/constants"
	"math"
	"log"
	"github.com/galaco/vrad/vmath/polygon"
)

func CalcFaceExtents(s *face.Face, lightmapTextureMinsInLuxels [2]int32, lightmapTextureSizeInLuxels[2]int32) {
	var mins, maxs [2]float32
	var i,j int
	e := int32(0)

	var v *mgl32.Vec3
	var tex *texinfo.TexInfo

	mins[1] = 1e24
	mins[0] = mins[1]
	maxs[1] = -1e24
	maxs[0] = maxs[1]

	tex = &(cache.GetLumpCache().TexInfo)[s.TexInfo]

	for i = 0; i < int(s.NumEdges); i++ {
		e = (cache.GetLumpCache().SurfEdges)[s.FirstEdge + int32(i)]
		//@TODO how does this work?
		if e >= 0 {
			v = &((cache.GetLumpCache().Vertexes)[cache.GetLumpCache().Edges[e][0]])
		} else {
			v = &((cache.GetLumpCache().Vertexes)[cache.GetLumpCache().Edges[-e][1]])
		}

		for j = 0; j < 2; j++ {
			val := v[0] * tex.LightmapVecsLuxelsPerWorldUnits[j][0] +
				v[1] * tex.LightmapVecsLuxelsPerWorldUnits[j][1] +
				v[2] * tex.LightmapVecsLuxelsPerWorldUnits[j][2] +
				tex.LightmapVecsLuxelsPerWorldUnits[j][3]

			if val < mins[j] {
				mins[j] = val
			}
			if val > maxs[j] {
				maxs[j] = val
			}
		}
	}

	var nMaxLightmapDim int
	 if s.DispInfo == -1 {
		 nMaxLightmapDim = constants.MAX_LIGHTMAP_DIM_WITHOUT_BORDER
	} else {
		 nMaxLightmapDim = constants.MAX_DISP_LIGHTMAP_DIM_WITHOUT_BORDER
	}

	for i = 0; i < 2; i++ {
		mins[i] = float32(math.Floor(float64(mins[i])))
		maxs[i] = float32(math.Ceil(float64(maxs[i])))

		lightmapTextureMinsInLuxels[i] = int32(mins[i])
		lightmapTextureSizeInLuxels[i] = int32(maxs[i] - mins[i])

		if int(lightmapTextureSizeInLuxels[i]) > nMaxLightmapDim + 1 {
			point := mgl32.Vec3{0, 0, 0}
			for j = 0; j < len(cache.GetLumpCache().Edges); j++ {
				e = (cache.GetLumpCache().SurfEdges)[s.FirstEdge + int32(j)]
				if e < 0 {
					v = &(cache.GetLumpCache().Vertexes[cache.GetLumpCache().Edges[-e][1]])
				} else {
					v = &(cache.GetLumpCache().Vertexes[cache.GetLumpCache().Edges[e][0]])
				}
				point = point.Add(*v)
				log.Printf("Bad surface extents point: %f %f %f\n", v.X(), v.Y(), v.Z())
			}
			point = point.Mul(1.0 / float32(s.NumEdges))
/**			@TODO Expose TexDataStringTable
			log.Fatal("Bad surface extents - surface is too big to have a lightmap\n\tmaterial %s around point (%.1f %.1f %.1f)\n\t(dimension: %d, %d>%d)\n",
				TexDataStringTable_GetString( dtexdata[cache.GetLumpCache().TexInfo[s.TexInfo].TexData].NameStringTableID ),
				point.X(), point.Y(), point.Z(),
				int(i),
				int(lightmapTextureSizeInLuxels[i]),
				int(nMaxLightmapDim + 1))
**/
		}
	}
}

func WindingFromFace(f *face.Face, origin *mgl32.Vec3) *polygon.Winding {
	var dv *mgl32.Vec3
	var v uint16
	var w *polygon.Winding

	w = polygon.NewWinding (int(f.NumEdges))
	w.NumPoints = int(f.NumEdges)

	for i := 0 ; i < int(f.NumEdges); i++ {
		se := cache.GetLumpCache().SurfEdges[int(f.FirstEdge) + i]
		if se < 0 {
			v = (cache.GetLumpCache().Edges)[-se][1]
		} else {
			v = (cache.GetLumpCache().Edges)[se][0]
		}

		dv = &(cache.GetLumpCache().Vertexes[v])
		w.Points[i] = dv.Add(*origin)
	}

	RemoveColinearPoints(w)

	return w
}