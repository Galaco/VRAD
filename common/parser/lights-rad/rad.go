package lights_rad

const MAX_TEXLIGHTS = 128

type Rad struct {
	NumTexlights int
	UseHDR bool
	ForcedTextureShadowsModels []string
	NonShadowCastingMaterials []string
}
