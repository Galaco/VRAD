package leafambient

import (
	"log"
	"github.com/galaco/vrad/cache"
)

func ComputePerLeafAmbientLighting() {
	// Figure out which lights should go in the per-leaf ambient cubes.
	nInAmbientCube := 0
	nSurfaceLights := 0
	for i := 0; i < *pNumworldlights; i++ {
		dworldlight_t * wl = &dworldlights[i]
		if IsLeafAmbientSurfaceLight(wl) {
			wl->flags |= DWL_FLAGS_INAMBIENTCUBE
		} else {
			wl->flags &= ~DWL_FLAGS_INAMBIENTCUBE
		}
		if wl.Type == emit_surface {
			nSurfaceLights++
		}

		if wl->flags & DWL_FLAGS_INAMBIENTCUBE {
			nInAmbientCube++
		}
	}

	nPercentage := 0
	if nSurfaceLights > 0 {
		nPercentage = (nInAmbientCube*100) / nSurfaceLights
	}
	log.Printf("%d of %d (%d%% of) surface lights went in leaf ambient cubes.\n", nInAmbientCube, nSurfaceLights, nPercentage)

	g_LeafAmbientSamples.SetCount(numleafs)

	if g_bUseMPI {
		// Distribute the work among the workers.
		VMPI_SetCurrentStage("ComputeLeafAmbientLighting")
		DistributeWork(numleafs, VMPI_DISTRIBUTEWORK_PACKETID, VMPI_ProcessLeafAmbient, VMPI_ReceiveLeafAmbientResults)
	} else {
		RunThreadsOn(numleafs, true, ThreadComputeLeafAmbient)
	}

	// now write out the data
	log.Println("Writing leaf ambient...")
	g_pLeafAmbientIndex->RemoveAll()
	g_pLeafAmbientLighting->RemoveAll()
	g_pLeafAmbientIndex->SetCount( numleafs )
	g_pLeafAmbientLighting->EnsureCapacity( numleafs*4 )
	for leafID := 0; leafID < len(cache.GetLumpCache().Leafs); leafID++ {
		const CUtlVector<ambientsample_t> &list = g_LeafAmbientSamples[leafID]
		g_pLeafAmbientIndex->Element(leafID).ambientSampleCount = list.Count()
		if !list.Count() {
			g_pLeafAmbientIndex->Element(leafID).firstAmbientSample = 0
		} else {
			g_pLeafAmbientIndex->Element(leafID).firstAmbientSample = g_pLeafAmbientLighting->Count()
			// compute the samples in disk format.  Encode the positions in 8-bits using leaf bounds fractions
			for i := 0; i < list.Count(); i++ {
				outIndex := g_pLeafAmbientLighting->AddToTail()
				dleafambientlighting_t &light = g_pLeafAmbientLighting->Element(outIndex)

				light.x = Fixed8Fraction( list[i].pos.x, dleafs[leafID].mins[0], dleafs[leafID].maxs[0] )
				light.y = Fixed8Fraction( list[i].pos.y, dleafs[leafID].mins[1], dleafs[leafID].maxs[1] )
				light.z = Fixed8Fraction( list[i].pos.z, dleafs[leafID].mins[2], dleafs[leafID].maxs[2] )
				light.pad = 0
				for side := 0; side < 6; side++ {
					VectorToColorRGBExp32( list[i].cube[side], light.cube.m_Color[side] );
				}
			}
		}
	}
	for i := 0; i < len(cache.GetLumpCache().Leafs); i++ {
		// UNDONE: Do this dynamically in the engine instead.  This will allow us to sample across leaf
		// boundaries always which should improve the quality of lighting in general
		if g_pLeafAmbientIndex->Element(i).ambientSampleCount == 0 {
			if !(dleafs[i].contents & CONTENTS_SOLID) {
				log.Printf("Bad leaf ambient for leaf %d\n", i)
			}

			refLeaf := NearestNeighborWithLight(i)
			g_pLeafAmbientIndex->Element(i).ambientSampleCount = 0
			g_pLeafAmbientIndex->Element(i).firstAmbientSample = refLeaf
		}
	}
	log.Printf("done\n")
}
