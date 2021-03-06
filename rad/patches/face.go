package patches

import (
	"github.com/galaco/vrad/cache"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/source-tools-common/constants"
	vrad_constants "github.com/galaco/vrad/common/constants"
	"github.com/galaco/vrad/vmath/polygon"
	"github.com/galaco/bsp/flags"
	"math"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/plane"
	"log"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/texdata"
	"github.com/galaco/vrad/vmath/vector"
	"strings"
)

var reflectivityScale = float32(1.0)
//var fakePlanes = int(0)
var texScale = true
var maxChop = float32(4) // coarsest allowed number of luxel widths for a patch
var minChop = float32(4) // "-chop" tightest number of luxel widths for a patch, used on edges

// @NOTE
// This was updated to remove the need to define fakePlanes
// largely because of variable sized (slice) of bsp Planes lump.
func MakePatchForFace(fn int, w *polygon.Winding) {
	f := &(*cache.GetTargetFaces())[fn]
	var area float32
	var patch *types.Patch
	var tx *texinfo.TexInfo

	// get texture info
	tx = &((cache.GetLumpCache().TexInfo)[f.TexInfo])

	// No patches at all for fog!
/*	if STATIC_FOG {
		if IsFog(f) {
			return
		}
	}
*/
	// the sky needs patches or the form factors don't work out correctly
	// if (IsSky( f ) )
	// 	return;

	area = polygon.WindingArea(w)
	if area <= 0 {
		numDegenerateFaces++
		// Msg("degenerate face\n");
		return
	}

	totalArea += area

	// get a patch

	cache.AddPatchToCache(&types.Patch{})
	ndxPatch := len(*cache.GetPatches()) - 1
	patch = &((*cache.GetPatches())[ndxPatch])
	patch.NdxNext = vrad_constants.CONSTRUCTS_INVALID_INDEX
	patch.NdxNextParent = vrad_constants.CONSTRUCTS_INVALID_INDEX
	patch.NdxNextClusterChild = vrad_constants.CONSTRUCTS_INVALID_INDEX
	patch.Child1 = vrad_constants.CONSTRUCTS_INVALID_INDEX
	patch.Child2 = vrad_constants.CONSTRUCTS_INVALID_INDEX
	patch.Parent = vrad_constants.CONSTRUCTS_INVALID_INDEX
	if tx.Flags & flags.SURF_BUMPLIGHT != 0 {
		patch.NeedsBumpMap = true
	} else {
		patch.NeedsBumpMap = false
	}

	// link and save patch data
	patch.NdxNext = cache.GetFacePatches()[fn]
	cache.SetFacePatch(fn, ndxPatch)
	//	patch->next = face_g_Patches[fn];
	//	face_g_Patches[fn] = patch;

	// compute a separate scale for chop - since the patch "scale" is the texture scale
	// we want textures with higher resolution lighting to be chopped up more
	chopScale := [2]float32{16.0, 16.0}
	if texScale == true {
		// Compute the texture "scale" in s,t
		for i := 0; i < 2; i++ {
			patch.Scale[i] = 0.0
			chopScale[i] = 0.0
			for j := 0; j < 3; j++ {
				patch.Scale[i] +=
					tx.TextureVecsTexelsPerWorldUnits[i][j] *
					tx.TextureVecsTexelsPerWorldUnits[i][j]
				chopScale[i] +=
					tx.LightmapVecsLuxelsPerWorldUnits[i][j] *
					tx.LightmapVecsLuxelsPerWorldUnits[i][j]
			}
			patch.Scale[i] = float32(math.Sqrt(float64(patch.Scale[i])))
			chopScale[i] = float32(math.Sqrt(float64(chopScale[i])))
		}
	} else {
		patch.Scale[0] = 1.0
		patch.Scale[1] = 1.0
	}

	patch.Area = area
	patch.Sky = IsSky(f)

	// chop scaled up lightmaps coarser
	patch.LuxScale = (chopScale[0] + chopScale[1]) / 2
	patch.Chop = maxChop

/*
#ifdef STATIC_FOG
    patch->fog = FALSE;
#endif
 */

	patch.Winding = w

	patch.Plane = &(cache.GetLumpCache().Planes[f.Planenum])

	// make a new plane to adjust for origined bmodels
	if cache.GetFaceOffsets()[fn][0] != 0 ||
		cache.GetFaceOffsets()[fn][1] != 0 ||
		cache.GetFaceOffsets()[fn][2] != 0 {
		var pl *plane.Plane
		numPlanes := len(cache.GetLumpCache().Planes)

		// origin offset faces must create new planes
		if numPlanes >= constants.MAX_MAP_PLANES {
			log.Fatal("numplanes >= MAX_MAP_PLANES")
		}

		//log.Printf("NumPlanes total: %d, numPlanes: %d, fakePlanes: %d\n", len(cache.GetLumpCache().Planes), numPlanes, fakePlanes)
		// Our bsplib uses slices.
		// If we wanna generate new Planes, we need to append a new one to use directly below
		(cache.GetLumpCache()).Planes = append(cache.GetLumpCache().Planes, plane.Plane{})
		//log.Printf("NewSize: %d\n", len(cache.GetLumpCache().Planes))

		pl = &(cache.GetLumpCache().Planes[numPlanes])
		//fakePlanes++

		*pl = *(patch.Plane)
		pl.Distance += cache.GetFaceOffsets()[fn].Dot(pl.Normal)
		patch.Plane = pl
	}

	patch.FaceNumber = fn
	polygon.WindingCenter(w, &(patch.Origin))

	// Save "center" for generating the face normals later.
	(cache.GetFaceCentroids()[fn]) = patch.Origin.Sub(cache.GetFaceOffsets()[fn])
	patch.Normal = patch.Plane.Normal

	polygon.WindingBounds(w, &(patch.FaceMins), &(patch.FaceMaxs))
	patch.Mins = patch.FaceMins
	patch.Maxs = patch.FaceMaxs

	BaseLightForFace(f, &patch.BaseLight, &patch.BaseArea, &patch.Reflectivity)

	// Chop all texlights very fine.
	if patch.BaseLight.ApproxEqual(mgl32.Vec3{0,0,0}) == false {
		// patch->chop = do_extra ? maxchop / 2 : maxchop;
		tx.Flags |= flags.SURF_LIGHT
	}

	// @TODO ADD DISPLACEMENT SUPPORT
	// get rid of do extra functionality on displacement surfaces
	if polygon.ValidDispFace(f){
		patch.Chop = maxChop
	}

	// @TODO GALACOS_PORT_NOTE
	// Below note copied from Valve implementation
	// FIXME: If we wanted to add a dependency from vrad to the material system,
	// we could do this. It would add a bunch of file accesses, though:

	/*
	// Check for a material var which would override the patch chop
	bool bFound;
	const char *pMaterialName = TexDataStringTable_GetString( dtexdata[ tx->texdata ].nameStringTableID );
	MaterialSystemMaterial_t hMaterial = FindMaterial( pMaterialName, &bFound, false );
	if ( bFound )
	{
		const char *pChopValue = GetMaterialVar( hMaterial, "%chop" );
		if ( pChopValue )
		{
			float flChopValue;
			if ( sscanf( pChopValue, "%f", &flChopValue ) > 0 )
			{
				patch->chop = flChopValue;
			}
		}
	}
	*/
}

func IsSky (f *face.Face) bool {
	var tx *texinfo.TexInfo

	tx = &cache.GetLumpCache().TexInfo[f.TexInfo]
	if (tx.Flags & flags.SURF_SKY) != 0 {
		return true
	}
	return false
}


func BaseLightForFace(f *face.Face, light *mgl32.Vec3, parea *float32, reflectivity *mgl32.Vec3) {
	var tx *texinfo.TexInfo
	var td *texdata.TexData

	//
	// check for light emited by texture
	//
	tx = &(cache.GetLumpCache().TexInfo[f.TexInfo])
	td = &(cache.GetLumpCache().TexData[tx.TexData])

	LightForTexture(cache.GetTexDataStringTable().GetString(int(td.NameStringTableID)), light)

	*parea = float32(td.Height * td.Width)
	vector.Scale(&td.Reflectivity, reflectivityScale, reflectivity)

	// always keep this less than 1 or the solution will not converge
	for i := 0; i < 3; i++ {
		if reflectivity[i] > 0.99 {
			reflectivity[i] = 0.99
		}
	}
}

func LightForTexture( name string, result *mgl32.Vec3 ) {
	result[0] = 0
	result[1] = 0
	result[2] = 0

	var baseFilename string

	if strings.HasPrefix(name, "maps/") == true {
		localName := strings.TrimLeft(name, "maps/")
		// this might be a patch texture for cubemaps.  try to parse out the original filename.
		if strings.Index(localName, cache.GetLumpCache().FileName) == 0 {
			base := strings.TrimLeft(localName, cache.GetLumpCache().FileName)
			if string(base[0]) == "/" {
				base = strings.TrimLeft(base, "/") // step past the path separator

				// now we've gotten rid of the 'maps/level_name/' part, so we're left with
				// 'originalName_%d_%d_%d'.
				baseFilename = base
				foundSeparators := true
				for i := 0; i < 3; i++ {
					underscore := strings.LastIndex(baseFilename, "_")
					if underscore != -1 {
						// @TODO FUUUUUUCK.
						// We *CAN* ignore this for now
						// We MUST come back to it later, otherwise
						// texture lights wont work!
						//baseFilename[underscore] = texdatastringtable.StringTableNullTerminator
						//*underscore = texdatastringtable.StringTableNullTerminator
					} else {
						foundSeparators = false
					}
				}

				if foundSeparators == true {
					name = baseFilename
				}
			}
		}
	}

	//log.Printf("%s\n", name)
	for i := 0; i < len(*cache.GetTexLightCache()) ; i++ {
		//log.Printf("%s\n", (*cache.GetTexLightCache())[i].Name)
		if name == (*cache.GetTexLightCache())[i].Name {
			result = &(*cache.GetTexLightCache())[i].Value
			return
		}
	}
}