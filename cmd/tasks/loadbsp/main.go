package loadbsp

import (
	"os"
	"strconv"
	"log"
	"strings"
	"github.com/galaco/vrad/cmd"
	"github.com/galaco/vrad/vmath/matrix"
	"github.com/galaco/bsp"
	"github.com/galaco/vmf"
	"github.com/galaco/bsp/primitives/model"
	"github.com/galaco/vrad/cmd/tasks/loadbsp/brush"
	brush2 "github.com/galaco/bsp/primitives/brush"
	"github.com/galaco/bsp/flags"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/vmath/polygon"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/raytracer"
	"github.com/galaco/vrad/common/parser/lights-rad"
	entitiesparser "github.com/galaco/vrad/common/parser/entities"
	"github.com/galaco/vrad/common/types"
	"github.com/galaco/source-tools-common/constants"
	vrad_constants "github.com/galaco/vrad/common/constants"
	"github.com/galaco/vrad/rad"
	"time"
)

// Main command function.
func Main(args *cmd.Args, transfered interface{}) (interface{}, error) {
	if args.LowPriority == true {
		// Go doesnt really support this...
	}

	if args.DumpPatches == true {
		//InitDumpPatchesFiles()
	}

	// Here we should prep the log file...

	if args.Lights != "" {
		//Try to load default lights.rad

		// @TODO

		// Load designer lights.rad if specified
		// @TODO what to do with response.
		if args.Lights != "" {
			file,_ := os.Open(args.Lights)
			reader := lights_rad.NewReader(file, file.Name())
			*cache.GetTexLightCache() = make([]types.TexLight, lights_rad.MAX_TEXLIGHTS)
			reader.Read(cache.GetTexLightCache())
			file.Close()
		}
	}

	log.Printf("Loading %s\n", args.Filename)
	file,err := loadBSP(args.Filename)
	if err != nil {
		return nil,err
	}
	cache.BuildLumpCache(args.Filename, file)

	//@TODO
	//g_pFullFileSystem->AddSearchPath(source, "GAME", PATH_ADD_TO_HEAD);
	//g_pFullFileSystem->AddSearchPath(source, "MOD", PATH_ADD_TO_HEAD);
	mapFlags := cache.GetLumpCache().MapFlags
	if args.StaticPropLighting {
		mask := 0x00000002
		if args.HDR {
			mask = 0x00000001
		}
		mapFlags.LevelFlags |= uint32(mask)
	} else {
		// @TODO This *could* be incorrect
		mapFlags.LevelFlags &^= uint32(0x00000002 | 0x00000001)
	}

	// Determine face target
	var numFaces = 0
	if args.HDR == true {
		cache.SetTargetFaces(&cache.GetLumpCache().FacesHDR)
		if len(*cache.GetTargetFaces()) == 0 {
			numFaces = len(cache.GetLumpCache().Faces)
		}
	} else {
		cache.SetTargetFaces(&cache.GetLumpCache().Faces)
	}

	entData := cache.GetLumpCache().EntData
	entImportAsVmfBlock,err := parseEntities(&entData)
	entities := entImportAsVmfBlock.Unclassified

	//@TODO THis should not be done this late!
	cache.SetEntityList(entitiesparser.Parse(entities))

	ExtractBrushEntityShadowCasters(&entities)

	//@TODO
	//StaticPropMgr()->Init();
	//StaticDispMgr()->Init();

	if cache.GetLumpCache().Visibility.NumClusters == 0 {
		log.Printf("No vis information, direct lighting only.\n")
		args.Bounce = 0
		cache.GetConfig().Ambient = mgl32.Vec3{0.1, 0.1, 0.1}

		// Equivalent of CountClusters
		for i := 0; i < len(cache.GetLumpCache().Leafs); i++ {
			if int32(cache.GetLumpCache().Leafs[i].Cluster) > cache.GetLumpCache().Visibility.NumClusters {
				cache.GetLumpCache().Visibility.NumClusters = int32(cache.GetLumpCache().Leafs[i].Cluster)
			}
		}
		cache.GetLumpCache().Visibility.NumClusters++
	}

	//
	// patches and referencing data (ensure capacity)
	//
	// TODO: change the maxes to the amount from the bsp!!
	//
	//	g_Patches.EnsureCapacity( MAX_PATCHES );
	for ndx := 0; ndx < constants.MAX_MAP_FACES; ndx++ {
		cache.SetFacePatch(ndx, -1)
		cache.SetFaceParent(ndx, -1)
	}

	for ndx := 0; ndx < constants.MAX_MAP_CLUSTERS; ndx++ {
		cache.SetClusterChild(ndx, -1)
	}

	// Setup ray tracer
	addBrushesForRayTrace()
	// @TODO
	//StaticDispMgr()->AddPolysForRayTrace();
	//StaticPropMgr()->AddPolysForRayTrace();

	// Dump raytracer for glview
	//@TODO
	//if g_bDumpRtEnv {
	//	WriteRTEnv("trace.txt");
	//}

	// Build acceleration structure
	log.Println("Setting up ray-trace acceleration structure... ")

	setupStart := time.Now().UnixNano() / int64(time.Millisecond)
	raytracer.GetEnvironment().SetupAccelerationStructure()
	setupEnd := time.Now().UnixNano() / int64(time.Millisecond)
	log.Printf("Done (%f seconds)\n", float32(setupEnd-setupStart) / 1000)

	rad.Start(args)

	// Setup incremental lighting.
/* @TODO Not supported yet. Is probably useless for all intents and purposes
	if g_pIncremental == true {
		if !g_pIncremental.Init( args.Filename, incrementfile ) {
			log.Println("Unable to load incremental lighting file in %s.\n", incrementfile)
		}
	}
*/
	return numFaces,nil
}

func loadBSP(filename string) (*bsp.Bsp,error){
	file,err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	reader := bsp.NewReader(file)
	return reader.Read()
}

// Parse Entity block.
// Vmf lib is actually capable of doing this;
// contents are loaded into Vmf.Unclassified
func parseEntities(data *string) (vmf.Vmf,error) {
	stringReader := strings.NewReader(*data)
	reader := vmf.NewReader(stringReader)

	return reader.Read()
}


// Some brush entities can cast shadows.
// We need to make a note of them
func ExtractBrushEntityShadowCasters(entities *vmf.Node) {
	models := cache.GetLumpCache().Models
	for _,iEntity := range *entities.GetAllValues() {
		entity := iEntity.(vmf.Node)
		if entity.HasProperty("vrad_brush_cast_shadows") == true {
			splOrigin := strings.Split(entity.GetProperty("origin"), " ")
			splAngles := strings.Split(entity.GetProperty("angles"), " ")
			origin := mgl32.Vec3{}
			for i,sf := range splOrigin {
				f,_ := strconv.ParseFloat(sf, 32)
				origin[i] = float32(f)
			}
			angles := mgl32.Vec3{}
			for i,sf := range splAngles {
				f,_ := strconv.ParseFloat(sf, 32)
				angles[i] = float32(f)
			}

			xform := matrix.Mat4{}
			xform.SetupMatrixOrgAngles( &origin, &angles )
			// Adds to raytrace environment
			addBrushes(brushmodelForEntity(&entity, &models), xform)
		}
	}
}

//Find brushmodel for associated index
func brushmodelForEntity(entity *vmf.Node, models *[]model.Model) *model.Model {
	modelName := entity.GetProperty("model")
	if len(modelName) > 1 {
		modelIndex,_ := strconv.ParseInt(modelName, 8, 32)
		modelIndex += 1
		if modelIndex > 0 && int(modelIndex) < len(*models) {
			return &(*models)[modelIndex]
		}
	}

	return nil
}

// Add brushes (NOTE: PLURAL) from a modal to raytracer environment
func addBrushes(model *model.Model, xform matrix.Mat4) {
	if model != nil {
		brushList := []int{}

		brush.GetBrushRecursive(int(model.HeadNode), &brushList)
		for i := 0; i < len(brushList); i++ {
			ndxBrush := brushList[i]
			addBrushToRaytraceEnvironment( &(cache.GetLumpCache().Brushes)[ndxBrush], &xform )
		}
	}
}

// Add a single brush to raytrace environment
func addBrushToRaytraceEnvironment(brush *brush2.Brush, xform *matrix.Mat4) {
	if 0 == brush.Contents & flags.MASK_OPAQUE {
		return
	}
	v0 := mgl32.Vec3{}
	v1 := mgl32.Vec3{}
	v2 := mgl32.Vec3{}

	for i := 0; i < int(brush.NumSides); i++ {
		side := (cache.GetLumpCache().BrushSides)[int(brush.FirstSide) + i]
		plane := (cache.GetLumpCache().Planes)[side.PlaneNum]
		tx := (cache.GetLumpCache().TexInfo)[side.TexInfo]
		w := polygon.BaseWindingForPlane(&plane.Normal, plane.Distance)

		if tx.Flags & flags.SURF_SKY == 0 || side.DispInfo != 0 {
			continue
		}

		for j := 0;  j < int(brush.NumSides) && w != nil; j++ {
			if i == j {
				continue
			}
			otherSide := (cache.GetLumpCache().BrushSides)[int(brush.FirstSide) + j]
			if otherSide.Bevel != 0 {
				continue
			}
			plane := (cache.GetLumpCache().Planes)[otherSide.PlaneNum ^ 1]
			polygon.ChopWindingInPlace(&w, &plane.Normal, plane.Distance, 0)
		}
		if w != nil {
			for j := 2; j < int(w.NumPoints); j++ {
				v0 = *xform.Mul4x3(&w.Points[0])
				v1 = *xform.Mul4x3(&w.Points[j - 1])
				v2 = *xform.Mul4x3(&w.Points[j])
				fullCoverage := mgl32.Vec3{1.0,0,0}
				raytracer.GetEnvironment().AddTriangle(raytracer.TRACE_ID_OPAQUE, &v0, &v1, &v2, &fullCoverage)
			}
			polygon.FreeWinding(w)
		}
	}
}

func addBrushesForRayTrace() {
	if len(cache.GetLumpCache().Models) == 0 {
		return
	}

	identity := matrix.Mat4{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	identity.Identity()

	brushList := []int{}
	brush.GetBrushRecursive(int(cache.GetLumpCache().Models[0].HeadNode), &brushList)

	for i := 0; i < len(brushList); i++ {
		pBrush := &cache.GetLumpCache().Brushes[brushList[i]]
		addBrushToRaytraceEnvironment(pBrush, &identity)
	}

	for i := 0; i < int(cache.GetLumpCache().Models[0].NumFaces); i++ {
		ndxFace := int(cache.GetLumpCache().Models[0].FirstFace) + i
		tFace := &(*cache.GetTargetFaces())[ndxFace]
		tx := &cache.GetLumpCache().TexInfo[tFace.TexInfo]

		if 0 == (tx.Flags & flags.SURF_SKY) {
			continue
		}

		points := [vrad_constants.MAX_POINTS_ON_WINDING]mgl32.Vec3{}

		for j := 0; j < int(tFace.NumEdges); j++ {
			if j > vrad_constants.MAX_POINTS_ON_WINDING {
				log.Fatal("***** ERROR! MAX_POINTS_ON_WINDING reached!")
			}

			if int(tFace.FirstEdge) + j >= len(cache.GetLumpCache().SurfEdges) {
				log.Fatal("***** ERROR! face->firstedge + j >= ARRAYSIZE( dsurfedges )!")
			}

			surfEdge := cache.GetLumpCache().SurfEdges[int(tFace.FirstEdge) + j]
			var v uint16

			if surfEdge < 0 {
				v = cache.GetLumpCache().Edges[-surfEdge][1]
			} else {
				v = cache.GetLumpCache().Edges[surfEdge][0]
			}

			if int(v) >= len(cache.GetLumpCache().Edges) {
				log.Fatalf("***** ERROR! v(%u) >= ARRAYSIZE( dvertexes(%d) )!", v, len(cache.GetLumpCache().Vertexes))
			}

			dv := &cache.GetLumpCache().Vertexes[v]
			points[j] = *dv
		}

		for j := 2; j < int(tFace.NumEdges); j++ {
			fullCoverage := mgl32.Vec3{1.0, 0, 0}
			raytracer.GetEnvironment().AddTriangle(raytracer.TRACE_ID_SKY, &points[0], &points[j - 1], &points[j], &fullCoverage)
		}
	}
}