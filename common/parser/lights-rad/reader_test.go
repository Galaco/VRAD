package lights_rad

import (
	"testing"
	"os"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/vrad/common/types"
)

func TestRead(t *testing.T) {
	file,_ := os.Open("lights.rad")

	reader := NewReader(file, file.Name())

	texLights := *cache.GetTexLightCache()
	for i := 0; i < MAX_TEXLIGHTS; i++ {
		cache.AddToTexLightCache(&types.TexLight{})
	}

	reader.Read(&texLights)
}