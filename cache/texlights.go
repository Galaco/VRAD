package cache

import "github.com/galaco/vrad/common/types"

var texLightCache []types.TexLight

func AddToTexLightCache(texLight *types.TexLight) {
	texLightCache = append(texLightCache, *texLight)
}

func GetTexLightCache() *[]types.TexLight {
	return &texLightCache
}