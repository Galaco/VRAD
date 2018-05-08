package lights_rad

import (
	"io"
	"bufio"
	"log"
	"strings"
	"fmt"
	"math"
	"github.com/galaco/vrad/common/types"
	"github.com/go-gl/mathgl/mgl32"
)

type Reader struct {
	file io.Reader
	name string
}

func (reader *Reader) Read(texLights *[]types.TexLight) *Rad {
	useHDR := false
	numTexlights := 0
	fileTexLights := 0
	forcedTextureShadowsModels := []string{}
	nonShadowCastingMaterials := []string{}

	log.Printf("[Reading texlights from '%s']\n", reader.name)

	bufReader := bufio.NewReader(reader.file)

	var buf []byte
	scan := make([]byte, 4)

	var lines []string
	line := make([]byte,1)
	for len(line) != 0 {
		line,_ = bufReader.ReadBytes('\r')
		lines = append(lines, strings.TrimLeft(string(line), "\n"))
	}

	for _,line := range lines {
		// skip if bad
		if len(line) == 0 {
			break
		}

		// Skip hdr section if not compiled with HDR
		if strings.Contains(string(line), "hdr:") {
			buf = append(buf, scan...)
			if useHDR == false {
				continue
			}
		}
		// Skip ldr section if not compiled with LDR
		if strings.Contains(string(line), "ldr:") {
			buf = append(buf, scan...)
			if useHDR == true {
				continue
			}
		}

		noShadName := ""
		if noShadow,_ := fmt.Sscanf(string(line), "noshadow %s", &noShadName); noShadow > 0 {
			if strings.HasSuffix(noShadName, ".") == true {
				noShadName = strings.Split(noShadName, ".")[0]
			}
			log.Printf("add %s as a non shadow casting material\n", noShadName)
			nonShadowCastingMaterials = append(nonShadowCastingMaterials, noShadName)
		} else if forceTexShadow,_ := fmt.Sscanf(string(line), "forcetextureshadow %s", &noShadName); forceTexShadow > 0 {
			log.Printf("add %s as a non shadow casting material\n", noShadName)
			forcedTextureShadowsModels = append(forcedTextureShadowsModels, forceTextureShadowsOnModel(noShadName))
		} else {
			texLight := ""
			value := mgl32.Vec3{}
			if numTexlights == MAX_TEXLIGHTS {
				log.Fatalf("Too many texlights, max = %d\n", MAX_TEXLIGHTS)
			}

			argCount,_ := fmt.Sscanf(string(line), "%s ", &texLight)
			if argCount != 1 {
				if len(line) != 4 {
					log.Printf("ignoring bad texlight '%s'\n", line)
				}
				continue
			}

			strLight := strings.TrimSpace(strings.Replace(string(line), texLight, "", -1))
			lightForString( strLight, &value )

			j := 0
			for j = 0; j < numTexlights; j++ {
				if strings.Compare((*texLights)[j].Name, texLight) == 0 {
					if strings.Compare(*(*texLights)[j].Filename, reader.name) == 0 {
						log.Printf("ERROR\a: Duplication of '%s' in file '%s'!\n", (*texLights)[j].Name, (*texLights)[j].Filename)
					} else if (*texLights)[j].Value[0] != value[0]  ||
						(*texLights)[j].Value[1] != value[1] ||
						(*texLights)[j].Value[2] != value[2] {
						log.Printf("Warning: Overriding '%s' from '%s' with '%s'!\n",
							(*texLights)[j].Name, (*texLights)[j].Filename, reader.name)
					} else {
						log.Printf("Warning: Redundant '%s' def in '%s' AND '%s'!\n",
							(*texLights)[j].Name, (*texLights)[j].Filename, reader.name)
					}
					break
				}
			}
			texLight = (*texLights)[j].Name
			(*texLights)[j].Value = value
			(*texLights)[j].Filename = &reader.name
			fileTexLights++

			numTexlights = int(math.Max( float64(numTexlights), float64(j + 1)))
		}
	}

	log.Printf("[%d texlights parsed from '%s']\n\n", fileTexLights, reader.name)

	return &Rad{
		UseHDR: false,
		NumTexlights:numTexlights,
		ForcedTextureShadowsModels: forcedTextureShadowsModels,
		NonShadowCastingMaterials: nonShadowCastingMaterials,
	}
}

func lightForString(light string, intensity *mgl32.Vec3) bool {

	// FIX THESE!
	lightScale := float32(1.0)
	useHDR := true

	intensity[0] = 0
	intensity[1] = 0
	intensity[2] = 0

	var r, g, b, scaler float32
	var rHDR,gHDR,bHDR,scalerHDR float32

	argCount,_ := fmt.Sscanf(light, "%e %e %e %e %e %e %e %e",
		&r, &g, &b, &scaler, &rHDR,&gHDR,&bHDR,&scalerHDR)

	if argCount == 8 {
		if useHDR == true {
			r = rHDR
			g = gHDR
			b = bHDR
			scaler = scalerHDR
		}
		argCount = 4
	}

	// make sure light is legal
	if r < 0.0 || g < 0.0 || b < 0.0 || scaler < 0.0 {
		// Actually don't need to do this, as intensity is already set as 0,0,0...
		//intensity.Init( 0.0f, 0.0f, 0.0f )
		return false
	}

	// convert to linear
	intensity[0] = float32(math.Pow( float64(r / 255.0), 2.2 ) * 255)				// convert to linear

	switch argCount {
	case 1:
		// The R,G,B values are all equal.
		intensity[2] = intensity[0]
		intensity[1] = intensity[2]
	case 3:
	case 4:
		// Save the other two G,B values.
		intensity[1] = float32(math.Pow( float64(g / 255.0), 2.2 ) * 255)
		intensity[2] = float32(math.Pow( float64(b / 255.0), 2.2 ) * 255)

		// Did we also get an "intensity" scaler value too?
		if argCount == 4 {
			// Scale the normalized 0-255 R,G,B values by the intensity scaler
			intensity[0] = intensity[0] * (scaler / 255.0)
			intensity[1] = intensity[1] * (scaler / 255.0)
			intensity[2] = intensity[2] * (scaler / 255.0)
		}

	default:
		log.Printf("unknown light specifier type - %s\n",light)
		return false
	}

	// scale up source lights by scaling factor
	intensity[0] = intensity[0] * lightScale
	intensity[1] = intensity[1] * lightScale
	intensity[2] = intensity[2] * lightScale

	return true
}

func forceTextureShadowsOnModel(modelName string) string {
	cleanModelName := strings.TrimLeft(modelName, "models/")
	cleanModelName = strings.TrimRight(cleanModelName, ".mdl")
	return cleanModelName
}


func NewReader(file io.Reader, name string) *Reader {
	return &Reader{
		file: file,
		name: name,
	}
}