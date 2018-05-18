package lightmap


import (
	"log"
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/vrad/common/constants"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/vmath/vector"
	"github.com/galaco/bsp/primitives/worldlight"
	"strings"
	"github.com/galaco/vrad/rad/clustertable"
	"github.com/galaco/vrad/common/constants/compiler"
)


const DIRECT_SCALE =  100.0*100.0
// @TODO Set this from -dlight flag
var lightThreshold = float64(0.1)
// @TODO Set this from -scale flag
var lightScale = float32(1.0)

/*
  =============
  CreateDirectLights
  =============
*/
func CreateDirectLights() {
	var p *types.Patch
	var dl *types.DirectLight //directlight_t
	var e *types.Entity
	var name string

	numDirectLights = 0

	// @TODO idk do we need to run this?
	//FreeDLights();

	//
	// surfaces
	//
	uiPatchCount := len(*cache.GetPatches())
	for i := 0; i< uiPatchCount; i++ {
		p = &((*cache.GetPatches())[i])

		// skip parent patches
		if p.Child1 != constants.CONSTRUCTS_INVALID_INDEX {
			continue
		}

		if p.BaseArea < 1e-6 {
			continue
		}

		if compiler.DEBUG_LEVEL > 6 {
			log.Printf("DirectLight: Baselight average = %e", vector.Avg(p.BaseLight))
		}
		if vector.Avg(p.BaseLight) >= lightThreshold {
			dl = AllocDLight(&p.Origin, true)

			dl.Light.Type = worldlight.EMIT_SURFACE
			dl.Light.Normal = p.Normal
			if 1.0e-20 < p.Normal.Len() {
				log.Fatalf("Patch normal out of bounds during DirectLight creation\n")
				//Assert(VectorLength(p.Normal ) > 1.0e-20)
			}
			// scale intensity by number of texture instances
			vector.Scale(&p.BaseLight, lightScale * p.Area * p.Scale[0] * p.Scale[1] / p.BaseArea, &dl.Light.Intensity)

			// scale to a range that results in actual light
			vector.Scale(&dl.Light.Intensity, DIRECT_SCALE, &(dl.Light.Intensity))
		}
	}

	//
	// entities
	// @NOTE
	// Custom light type definitions could be added here.
	// So long as radiosity implements the light. It could be theoretically possible to use
	// toggleable lights as many light_environments
	for i := 0; i < len(*cache.GetAllEntities()); i++ {
		e = cache.GetEntity(i)
		name = e.ValueForKey("classname")
		if strings.HasPrefix(name, "light") == false {
			continue
		}

		// Light_dynamic is actually a real entity; not to be included here...
		if name == "light_dynamic" {
			continue
		}

		if name == "light_spot" {
//			ParseLightSpot(e, dl)
		} else if name == "light_environment" {
//			ParseLightEnvironment(e, dl)
		} else if name == "light" {
//			ParseLightPoint(e, dl)
		} else {
			log.Printf("unsupported light entity: \"%s\"\n", name)
		}
	}

	log.Printf("%d direct lights\n", numDirectLights)
	// exit(1);
}


func AllocDLight( origin *mgl32.Vec3, addToList bool ) *types.DirectLight {
	var dl *types.DirectLight

	dl = types.NewDirectLight()

	if compiler.DEBUG_LEVEL > 3 {
		log.Printf("Allocated directLight no: %d\n", numDirectLights)
	}

	numDirectLights++
	dl.Index = numDirectLights

	dl.Light.Origin = *origin

	dl.Light.Cluster = int32(clustertable.ClusterFromPoint(&dl.Light.Origin))
	SetDLightVis( dl, int(dl.Light.Cluster) )

	dl.FaceNum = -1

	if addToList == true {
		// @TODO Broke!
		dl.Next = activeLights
		activeLights = dl
	}

	return dl
}

func SetDLightVis(dl *types.DirectLight, cluster int) {
	if dl.PVS == nil {
		dl.PVS = make([]byte, (cache.GetLumpCache().Visibility.NumClusters / 8) + 1)
	}

	GetVisCache(-1, cluster, &(dl.PVS))
}