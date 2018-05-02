package cmd

import (
	"flag"
	"log"
)

type Args struct {
	Filename string // Added here
	StaticPropLighting bool
	StaticPropNormals bool
	OnlyStaticProps bool
	StaticPropPolys bool
	DisablePropSelfShadowing bool
	TextureShadows bool
	DumpPatches bool
	NoDetailLighting bool
	DumpNormals bool
	RedErrors bool
	DumpTrace bool
	LargeDispSampleRadius bool
	Bounce int
	Verbose bool
	Threads int
	Lights string
	NoExtra bool
	DebugExtra bool
	FastAmbient bool
	Fast bool
	NoSkyboxRecurse bool
	Final bool
	ExtraSky bool
	CenterSamples int
	SmoothingThreshold int
	DLightMap bool
	LuxelDensity int
	LowPriority bool
	LogHashData bool
	OnlyDetail bool
	SoftSun bool
	MaxDispSampleSize int
	StopOnExit bool
	FullMiniDumps bool
	HDR bool
	LDR bool
	MaxChop int
	Chop int
	DispChop int
	DispPatchRadius int
}

func GetArgs() *Args{
	// Parse args here later
	Filename := flag.String("filename", "TESTDATA/ze_bioshock_v6_2.bsp", "filename=<string>") // Added here
	StaticPropLighting := flag.Bool("StaticPropLighting", false, "")
	StaticPropNormals := flag.Bool("StaticPropNormals", false, "")
	OnlyStaticProps := flag.Bool("OnlyStaticProps", false, "")
	StaticPropPolys := flag.Bool("StaticPropPolys", false, "")
	DisablePropSelfShadowing := flag.Bool("nossprops", false, "")
	TextureShadows := flag.Bool("textureshadows", false, "")
	DumpPatches := flag.Bool("dump", false, "")
	NoDetailLighting := flag.Bool("nodetaillight", false, "")
	RedErrors := flag.Bool("rederrors", false, "")
	DumpNormals := flag.Bool("dumpnormals", false, "")
	DumpTrace := flag.Bool("dumptrace", false, "")
	LargeDispSampleRadius := flag.Bool("LargeDispSampleRadius", false, "")
	Bounce := flag.Int("bounce", 1, "")
	Verbose := flag.Bool("verbose", false, "")
	Threads := flag.Int("threads", 1, "")
	Lights := flag.String("lights", "", "")
	NoExtra := flag.Bool("noextra", false, "")
	DebugExtra := flag.Bool("debugextra", false, "")
	FastAmbient := flag.Bool("fastambient", false, "")
	Fast := flag.Bool("fast", false, "")
	NoSkyboxRecurse := flag.Bool("noskyboxrecurse", false, "")
	Final := flag.Bool("final", false, "")
	ExtraSky := flag.Bool("extrasky", false, "")
	CenterSamples := flag.Int("centersamples", 1, "")
	SmoothingThreshold := flag.Int("smooth", 1, "")
	DLightMap := flag.Bool("dlightmap", false, "")
	LuxelDensity := flag.Int("luxeldensity", 1, "")
	LowPriority := flag.Bool("low", false, "")
	LogHashData := flag.Bool("loghash", false, "")
	OnlyDetail := flag.Bool("onlydetail", false, "")
	SoftSun := flag.Bool("softsun", false, "")
	MaxDispSampleSize := flag.Int("maxsampledispsize", 1, "")
	StopOnExit := flag.Bool("StopOnExit", false, "")
	FullMiniDumps := flag.Bool("FullMinidumps", false, "")
	HDR := flag.Bool("hdr", false, "")
	LDR := flag.Bool("ldr", false, "")
	MaxChop := flag.Int("maxchop", 1, "")
	Chop := flag.Int("chop", 1, "")
	DispChop := flag.Int("dispchop", 1, "")
	DispPatchRadius := flag.Int("disppatchradius", 1, "")

	log.Println(*Filename)
	return &Args{
		*Filename,
		*StaticPropLighting,
		*StaticPropNormals,
		*OnlyStaticProps,
		*StaticPropPolys,
		*DisablePropSelfShadowing,
		*TextureShadows,
		*DumpPatches,
		*NoDetailLighting,
		*DumpNormals,
		*RedErrors,
		*DumpTrace,
		*LargeDispSampleRadius,
		*Bounce,
		*Verbose,
		*Threads,
		*Lights,
		*NoExtra,
		*DebugExtra,
		*FastAmbient,
		*Fast,
		*NoSkyboxRecurse,
		*Final,
		*ExtraSky,
		*CenterSamples,
		*SmoothingThreshold,
		*DLightMap,
		*LuxelDensity,
		*LowPriority,
		*LogHashData,
		*OnlyDetail,
		*SoftSun,
		*MaxDispSampleSize,
		*StopOnExit,
		*FullMiniDumps,
		*HDR,
		*LDR,
		*MaxChop,
		*Chop,
		*DispChop,
		*DispPatchRadius,
	}
}