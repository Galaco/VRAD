package patches

import (
	"log"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/vrad/vmath/polygon"
	"github.com/galaco/bsp/primitives/model"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/common/parser/entities"
	"strings"
	"github.com/galaco/vmf"
	"strconv"
	"github.com/galaco/vrad/rad/world"
)

var numDegenerateFaces = 0
var totalArea = float32(0)

func MakePatches() {
	var f *face.Face
	var w *polygon.Winding
	var mod *model.Model
	origin := mgl32.Vec3{0,0,0}
	var ent *types.Entity

	// Parse Entdata
	entData := cache.GetLumpCache().EntData
	stringReader := strings.NewReader(entData)
	reader := vmf.NewReader(stringReader)
	entImportAsVmfBlock,_ := reader.Read()
	vmfEntities := entImportAsVmfBlock.Unclassified
	// Need to store these somewhere?
	entList := entities.Parse(vmfEntities)
	log.Printf("%d faces\n", len(*cache.GetTargetFaces()))

	for i := 0; i < len(cache.GetLumpCache().Models); i++ {
		mod = &(cache.GetLumpCache().Models[i])
		ent = entities.EntityForModel(i, &entList)
		origin = mgl32.Vec3{0, 0, 0}

		// bmodels with origin brushes need to be offset into their
		// in-use position
		origin = getVectorFromKey(ent, "origin")


		for j := 0 ; j < int(mod.NumFaces); j++ {
			fn := mod.FirstFace + int32(j)
			(*cache.GetFaceEntities())[fn] = ent
			(*cache.GetFaceOffsets())[fn] = origin
			f = &(*cache.GetTargetFaces())[fn]
			if f.DispInfo == -1 {
				w = world.WindingFromFace(f, &origin )
				MakePatchForFace(int(fn), w)
			}
		}
	}

	if numDegenerateFaces > 0 {
		log.Printf("%d degenerate faces\n", numDegenerateFaces)
	}
	log.Printf("%d square feet [%.2f square inches]\n", int(totalArea/144), totalArea)

	//@TODO ADD DISPLACEMENT SUPPORT
	// make the displacement surface patches
	//StaticDispMgr()->MakePatches()
}


func getVectorFromKey(ent *types.Entity, key string) mgl32.Vec3 {
	e := ent.EPairs
	for e != nil {
		if e.Key == key {
			break
		}

		e = e.Next
	}
	splOrigin := strings.Split(e.Value, " ")
	origin := mgl32.Vec3{}
	for i,sf := range splOrigin {
		f,_ := strconv.ParseFloat(sf, 32)
		origin[i] = float32(f)
	}

	return origin
}