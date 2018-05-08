package cache

import (
	"github.com/galaco/bsp"
	"github.com/galaco/bsp/primitives/plane"
	"github.com/galaco/bsp/primitives/texdata"
	"github.com/galaco/bsp/primitives/face"
	"github.com/galaco/bsp/primitives/leaf"
	"github.com/galaco/bsp/primitives/model"
	"github.com/galaco/bsp/primitives/brush"
	"github.com/galaco/bsp/primitives/brushside"
	"github.com/galaco/bsp/primitives/area"
	"github.com/galaco/bsp/primitives/areaportal"
	"github.com/galaco/bsp/primitives/mapflags"
	"github.com/galaco/bsp/primitives/texinfo"
	"github.com/galaco/bsp/primitives/node"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/bsp/primitives/visibility"
)

type LumpCache struct {
	EntData string
	Planes []plane.Plane
	TexData []texdata.TexData
	Vertexes []mgl32.Vec3
	Visibility visibility.Vis
	Nodes []node.Node
	TexInfo []texinfo.TexInfo
	Faces []face.Face
	Leafs []leaf.Leaf
	LeafBrushes []uint16
	Models []model.Model
	Brushes []brush.Brush
	BrushSides []brushside.BrushSide
	Areas []area.Area
	AreaPortals []areaportal.AreaPortal
	MapFlags mapflags.MapFlags
	FacesHDR []face.Face
}

//Well this should only be called once
func BuildLumpCache(file *bsp.Bsp) *LumpCache {
	// Safe to assume that cache was built if there are n>0 planes
	if &lumpCache != nil {
		return &lumpCache
	}
	lumpCache = LumpCache{}
	lumpCache.EntData = (*file.GetLump(bsp.LUMP_ENTITIES).GetContents()).GetData().(string)
	lumpCache.Planes = *(*file.GetLump(bsp.LUMP_PLANES).GetContents()).GetData().(*[]plane.Plane)
	lumpCache.TexData = *(*file.GetLump(bsp.LUMP_TEXDATA).GetContents()).GetData().(*[]texdata.TexData)
	lumpCache.Vertexes = *(*file.GetLump(bsp.LUMP_VERTEXES).GetContents()).GetData().(*[]mgl32.Vec3)
	lumpCache.Visibility = *(*file.GetLump(bsp.LUMP_VISIBILITY).GetContents()).GetData().(*visibility.Vis)
	lumpCache.Nodes = *(*file.GetLump(bsp.LUMP_NODES).GetContents()).GetData().(*[]node.Node)
	lumpCache.TexInfo = *(*file.GetLump(bsp.LUMP_TEXINFO).GetContents()).GetData().(*[]texinfo.TexInfo)
	lumpCache.Faces = *(*file.GetLump(bsp.LUMP_FACES).GetContents()).GetData().(*[]face.Face)
	lumpCache.Leafs = *(*file.GetLump(bsp.LUMP_LEAFS).GetContents()).GetData().(*[]leaf.Leaf)
	lumpCache.LeafBrushes = *(*file.GetLump(bsp.LUMP_LEAFBRUSHES).GetContents()).GetData().(*[]uint16)
	lumpCache.Models = *(*file.GetLump(bsp.LUMP_MODELS).GetContents()).GetData().(*[]model.Model)
	lumpCache.Brushes = *(*file.GetLump(bsp.LUMP_BRUSHES).GetContents()).GetData().(*[]brush.Brush)
	lumpCache.BrushSides = *(*file.GetLump(bsp.LUMP_BRUSHSIDES).GetContents()).GetData().(*[]brushside.BrushSide)
	lumpCache.Areas = *(*file.GetLump(bsp.LUMP_AREAS).GetContents()).GetData().(*[]area.Area)
	lumpCache.AreaPortals = *(*file.GetLump(bsp.LUMP_AREAPORTALS).GetContents()).GetData().(*[]areaportal.AreaPortal)
	lumpCache.MapFlags = *(*file.GetLump(bsp.LUMP_MAP_FLAGS).GetContents()).GetData().(*mapflags.MapFlags)
	lumpCache.FacesHDR = *(*file.GetLump(bsp.LUMP_FACES_HDR).GetContents()).GetData().(*[]face.Face)


	return &lumpCache
}

func GetLumpCache() *LumpCache {
	return &lumpCache
}

var lumpCache LumpCache