package lightmap

import (
	"log"
	"github.com/galaco/vrad/cache"
	"github.com/galaco/bsp/primitives/visibility"
)

func GetVisCache( lastOffset int, cluster int, pvs *[]byte) int {
// get the PVS for the pos to limit the number of checks
	if cache.GetLumpCache().Visibility.NumClusters == 0 {
		for i := range *pvs {
			if i < int((cache.GetLumpCache().Visibility.NumClusters + 7) / 8) {
				(*pvs)[i] = 255
			} else {
				break
			}
		}
		lastOffset = -1
	} else {
		if cluster < 0 {
			// Error, point embedded in wall
			// sampled[0][1] = 255;
			for i := range *pvs {
				if i < int((cache.GetLumpCache().Visibility.NumClusters + 7) / 8) {
					(*pvs)[i] = 255
				} else {
					break
				}
			}
			lastOffset = -1
		} else {
			thisOffset := int(cache.GetLumpCache().Visibility.ByteOffset[cluster][visibility.DVIS_PVS])
			if thisOffset != lastOffset {
				if thisOffset == -1 {
					log.Fatalf("visofs == -1\n")
				}

				visRunlength := cache.GetLumpCache().VisDataRaw[thisOffset:]
				pvs = DecompressVis(&visRunlength, len(*pvs))
			}
			lastOffset = thisOffset
		}
	}
	return lastOffset
}

/*
===================
DecompressVis
===================
*/
func DecompressVis(in *[]byte, length int) *[]byte {
	var c int
	var out = make([]byte, length)
	var row int
	var inOffset = 0
	var outOffset = 0

	row = int((cache.GetLumpCache().Visibility).NumClusters + 7) >> 3

	hasSimulatedDoWhile := false
	for (outOffset < len(out) && int(out[outOffset]) < row) || hasSimulatedDoWhile == false {
		hasSimulatedDoWhile = true

		// @NOTE The ++ lines may need to shiftto the stop
		// This will cause an out-of-bounds unless we compare to len()-1
		if inOffset < len(*in) {
			out[outOffset] = (*in)[inOffset]
			inOffset++
			outOffset++
			continue
		}

		c = int((*in)[1])
		if c == 0 {
			log.Fatalf("DecompressVis: 0 repeat")
		}

		inOffset += 2
		if (int(out[outOffset])) + c > row {
			c = row - int(out[outOffset])
			log.Printf("warning: Vis decompression overrun\n")
		}

		for c > 0 {
			outOffset++
			out[outOffset] = 0
			c--
		}
	}

	return &out
}
