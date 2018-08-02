package cameras

import (
	"github.com/galaco/vrad/cache"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/raytracer/trace"
	"log"
)

func ProcessSkyCameras() {
	numSkyCameras := 0
	for i := 0; i < len(cache.GetLumpCache().Areas); i++ {
		cache.SetAreaSkyCamera(i, -1)
	}

	entities := cache.GetAllEntities();
	for i := 0; i < len(*entities); i++ {
		e := (*entities)[i]
		name := e.ValueForKey("classname")
		if name != "sky_camera" {
			continue
		}

		var origin mgl32.Vec3
		e.VectorForKey("origin")
		node := trace.PointLeafnum(&origin)
		area := -1

		if node >= 0 && node < len(cache.GetLumpCache().Leafs) {
			area = int(cache.GetLumpCache().Leafs[node].Area())
		}
		scale := e.FloatForKey("scale")

		if scale > 0.0 {
			cache.GetSkyCamera(numSkyCameras).Origin = origin
			cache.GetSkyCamera(numSkyCameras).SkyToWorld = scale
			cache.GetSkyCamera(numSkyCameras).WorldToSky = 1.0 / scale
			cache.GetSkyCamera(numSkyCameras).Area = area

			if area >= 0 && area < len(cache.GetLumpCache().Areas) {
				cache.SetAreaSkyCamera(area, numSkyCameras)
			}

			numSkyCameras++
		}
	}

	log.Printf("Found %d camera(s)\n", numSkyCameras)
}
