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
	"math"
	"github.com/galaco/vrad/vmath"
	"github.com/galaco/vrad/vmath/quadratic"
	"github.com/galaco/vrad/vmath/quaterion"
	"strconv"
)


const DIRECT_SCALE =  100.0*100.0
// @TODO Set this from -dlight flag
var lightThreshold = float64(0.1)
// @TODO Set this from -scale flag
var lightScale = float32(1.0)
// @TODO Set this from -hdr flag
const useHDR = false
// @TODO Set this from -softsun flag
var sunAngularExtent = 1.0

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
			ParseLightSpot(e, dl)
		} else if name == "light_environment" {
			ParseLightEnvironment(e, dl)
		} else if name == "light" {
			ParseLightPoint(e, dl)
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
		// @NOTE
		// Funky stuff here. activeLights isn't *quite* what it sounds like
		// Its a linked list, where variable is a single light.
		// Notice what we do is replace using dl as the new start of the list, but first referencing the current start as
		// dl.Next
		dl.Next = activeLights
		activeLights = dl
	}

	return dl
}

func AddDLightToActiveList(dl *types.DirectLight) {
	dl.Next = activeLights
	activeLights = dl
}

func SetDLightVis(dl *types.DirectLight, cluster int) {
	if dl.PVS == nil {
		dl.PVS = make([]byte, (cache.GetLumpCache().Visibility.NumClusters / 8) + 1)
	}

	GetVisCache(-1, cluster, &(dl.PVS))
}

//-----------------------------------------------------------------------------
// Various parsing methods
//-----------------------------------------------------------------------------

func ParseLightGeneric(e *types.Entity, dl *types.DirectLight) {
	var e2 *types.Entity
	var target string
	dest := mgl32.Vec3{}

	dl.Light.Style = int32(e.FloatForKey("style"))

	// get intensity
	var err error
	dl.Light.Intensity,err = e.LightForKey("_lightHDR")
	if useHDR == true && err != nil {
	} else {
		dl.Light.Intensity,_ = e.LightForKey("_light")
	}

	// check angle, targets
	target = e.ValueForKey("target")
	if target[0] != 0 {    // point towards target
		e2 = cache.FindTargetEntity(target)
		if e2 == nil {
			log.Printf("WARNING: light at (%d %d %d) has missing target\n",
				int(dl.Light.Origin[0]), int(dl.Light.Origin[1]), int(dl.Light.Origin[2]))
		} else {
			dest = e2.VectorForKey("origin")
			dl.Light.Normal = dest.Sub(dl.Light.Origin)
			dl.Light.Normal = dl.Light.Normal.Normalize()
		}
	} else {
		// point down angle
		angles := e.VectorForKey("angles")
		pitch := e.FloatForKey("pitch")
		angle := e.FloatForKey("angle")
		SetupLightNormalFromProps(&quaterion.QAngle{angles.X(), angles.Y(), angles.Z()}, angle, pitch, &dl.Light.Normal )
	}
	if useHDR == true {
		vector.Scale(&dl.Light.Intensity,
			e.FloatForKeyWithDefault("_lightscaleHDR", 1.0),
			&dl.Light.Intensity)
	}
}


func ParseLightSpot(e *types.Entity, dl *types.DirectLight) {
	dest := e.VectorForKey("origin")
	dl = AllocDLight(&dest, true)

	ParseLightGeneric(e, dl)

	dl.Light.Type = worldlight.EMIT_SPOTLIGHT

	dl.Light.Stopdot = e.FloatForKey("_inner_cone")
	if 0 == dl.Light.Stopdot {
		dl.Light.Stopdot = 10
	}

	dl.Light.Stopdot2 = e.FloatForKey("_cone")
	if 0 == dl.Light.Stopdot2 {
		dl.Light.Stopdot2 = dl.Light.Stopdot
	}
	if dl.Light.Stopdot2 < dl.Light.Stopdot {
		dl.Light.Stopdot2 = dl.Light.Stopdot
	}

	// This is a point light if stop dots are 180...
	if (dl.Light.Stopdot == 180) && (dl.Light.Stopdot2 == 180) {
		dl.Light.Stopdot2 = 0
		dl.Light.Stopdot = 0
		dl.Light.Type = worldlight.EMIT_POINT
		dl.Light.Exponent = 0
	} else {
		// Clamp to 90, that's all DX8 can handle!
		if dl.Light.Stopdot > 90 {
			log.Printf("WARNING: light_spot at (%d %d %d) has inner angle larger than 90 degrees! Clamping to 90...\n",
			int(dl.Light.Origin[0]), int(dl.Light.Origin[1]), int(dl.Light.Origin[2]))
			dl.Light.Stopdot = 90
		}

		if dl.Light.Stopdot2 > 90 {
			log.Printf("WARNING: light_spot at (%d %d %d) has outer angle larger than 90 degrees! Clamping to 90...\n",
				int(dl.Light.Origin[0]), int(dl.Light.Origin[1]), int(dl.Light.Origin[2]))
			dl.Light.Stopdot2 = 90
		}

		dl.Light.Stopdot2 = float32(math.Cos(float64(dl.Light.Stopdot2 / 180 * vmath.PI)))
		dl.Light.Stopdot = float32(math.Cos(float64(dl.Light.Stopdot / 180 * vmath.PI)))
		dl.Light.Exponent = e.FloatForKey("_exponent")
	}

	SetLightFalloffParams(e,dl)
}

func SetLightFalloffParams(e *types.Entity, dl *types.DirectLight) {
	d50 := e.FloatForKey("_fifty_percent_distance")
	dl.StartFadeDistance = 0
	dl.EndFadeDistance = - 1
	dl.CapDistance = 1.0e22
	if 0 != d50 {
		d0 := e.FloatForKey("_zero_percent_distance")
		if d0 < d50 {
			log.Printf( "light has _fifty_percent_distance of %f but _zero_percent_distance of %f\n", d50, d0)
			d0 = 2.0 * d50
		}
		a := float32(0.0)
		b := float32(1.0)
		c := float32(0.0)
		if ! quadratic.SolveInverseQuadraticMonotonic( 0, 1.0, d50, 2.0, d0, 256.0, &a, &b, &c ) {
			log.Printf( "can't solve quadratic for light %f %f\n", d50, d0 )
		}
		// it it possible that the parameters couldn't be used because of enforing monoticity. If so, rescale so at
		// least the 50 percent value is right
//		printf("50 percent=%f 0 percent=%f\n",d50,d0);
// 		printf("a=%f b=%f c=%f\n",a,b,c);
		v50 := c + d50 * ( b + d50 * a )
		scale := 2.0 / v50
		a *= scale
		b *= scale
		c *= scale
// 		printf("scaled=%f a=%f b=%f c=%f\n",scale,a,b,c);
// 		for(float d=0;d<1000;d+=20)
// 			printf("at %f, %f\n",d,1.0/(c+d*(b+d*a)));
		dl.Light.QuadraticAttenuation = float32(a)
		dl.Light.LinearAttenuation = float32(b)
		dl.Light.ConstantAttenuation = float32(c)



		if 0 != e.IntForKey("_hardfalloff") {
			dl.EndFadeDistance = d0
			dl.StartFadeDistance = 0.75 * d0 + 0.25 * d50		// start fading 3/4 way between 50 and 0. could allow adjust
		} else {
			// now, we will find the point at which the 1/x term reaches its maximum value, and
			// prevent the light from going past there. If a user specifes an extreme falloff, the
			// quadratic will start making the light brighter at some distance. We handle this by
			// fading it from the minimum brightness point down to zero at 10x the minimum distance
			if math.Abs(float64(a)) > 0. {
				flMax := float32(b / (- 2.0 * a))				// where f' = 0
				if flMax > 0.0 {
					dl.CapDistance = flMax
					dl.StartFadeDistance = flMax
					dl.EndFadeDistance = 10.0 * flMax
				}
			}
		}
	} else {
		dl.Light.ConstantAttenuation = e.FloatForKey("_constant_attn")
		dl.Light.LinearAttenuation = e.FloatForKey("_linear_attn")
		dl.Light.QuadraticAttenuation = e.FloatForKey("_quadratic_attn")

		dl.Light.Radius = e.FloatForKey("_distance")

		// clamp values to >= 0
		if dl.Light.ConstantAttenuation < vmath.EQUAL_EPSILON {
			dl.Light.ConstantAttenuation = 0
		}

		if dl.Light.LinearAttenuation < vmath.EQUAL_EPSILON {
			dl.Light.LinearAttenuation = 0
		}

		if dl.Light.QuadraticAttenuation < vmath.EQUAL_EPSILON {
			dl.Light.QuadraticAttenuation = 0
		}

		if dl.Light.ConstantAttenuation < vmath.EQUAL_EPSILON && dl.Light.LinearAttenuation < vmath.EQUAL_EPSILON && dl.Light.QuadraticAttenuation < vmath.EQUAL_EPSILON {
			dl.Light.ConstantAttenuation = 1
		}

		// scale intensity for unit 100 distance
		ratio := dl.Light.ConstantAttenuation + 100 * dl.Light.LinearAttenuation + 100 * 100 * dl.Light.QuadraticAttenuation
		if ratio > 0 {
			vector.Scale(&dl.Light.Intensity, ratio, &dl.Light.Intensity )
		}
	}
}

func SetupLightNormalFromProps(angles *quaterion.QAngle, angle float32, pitch float32, output *mgl32.Vec3 ) {
	if angle == constants.ANGLE_UP {
		output[0] = 0
		output[1] = 0
		output[2] = 1
	} else if angle == constants.ANGLE_DOWN {
		output[0] = 0
		output[1] = 0
		output[2] = -1
	} else {
		// if we don't have a specific "angle" use the "angles" YAW
		if 0 == angle {
			angle = angles[constants.YAW]
		}

		output[2] = 0
		output[0] = float32(math.Cos(float64(angle) / 180 * vmath.PI))
		output[1] = float32(math.Cos(float64(angle) / 180 * vmath.PI))
	}

	if 0 == pitch {
		// if we don't have a specific "pitch" use the "angles" PITCH
		pitch = angles[constants.PITCH]
	}

	output[2] = float32(math.Sin(float64(pitch) / 180 * vmath.PI))
	output[0] *= float32(math.Cos(float64(pitch) / 180 * vmath.PI))
	output[1] *= float32(math.Cos(float64(pitch) / 180 * vmath.PI))
}

func ParseLightEnvironment(e *types.Entity, dl *types.DirectLight) {
	dest := e.VectorForKey ("origin")
	dl = AllocDLight(&dest, false)

	ParseLightGeneric(e, dl)

	angleStr := e.ValueForKeyWithDefault("SunSpreadAngle", "")
	if angleStr != "" {
		sunAngularExtent,_ = strconv.ParseFloat(angleStr, 32)
		sunAngularExtent = math.Sin((vmath.PI/180.0) * sunAngularExtent)
		log.Printf("sun extent from map=%f\n", sunAngularExtent)
	}
	if nil == globalSkyLight {
		// Sky light.
		globalSkyLight = dl
		dl.Light.Type = worldlight.EMIT_SKYLIGHT

		// Sky ambient light.
		globalAmbient = AllocDLight(&dl.Light.Origin, false );
		globalAmbient.Light.Type = worldlight.EMIT_SKYAMBIENT;
		var err error
		globalAmbient.Light.Intensity,err = e.LightForKey("_ambientHDR")
		if useHDR && err == nil {
			// we have a valid HDR ambient light value
		} else {
			globalAmbient.Light.Intensity,err = e.LightForKey("_ambient")
			if err == nil {
				vector.Scale(&dl.Light.Intensity, 0.5, &globalAmbient.Light.Intensity)
			}
		}
		if useHDR == true {
			vector.Scale(&globalAmbient.Light.Intensity,
						 e.FloatForKeyWithDefault("_AmbientScaleHDR", 1.0),
				&globalAmbient.Light.Intensity)
		}

		BuildVisForLightEnvironment()

		// Add sky and sky ambient lights to the list.
		AddDLightToActiveList(globalSkyLight)
		AddDLightToActiveList(globalAmbient)
	}
}

func ParseLightPoint(e *types.Entity, dl *types.DirectLight) {
	dest := e.VectorForKey("origin")
	dl = AllocDLight(&dest, true)

	ParseLightGeneric(e, dl)

	dl.Light.Type = worldlight.EMIT_POINT

	SetLightFalloffParams(e,dl)
}