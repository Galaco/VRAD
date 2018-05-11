package world

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/common/constants"
	"github.com/galaco/vrad/vmath/polygon"
)

// @TODO refactor this away?
var c_removed int

func RemoveColinearPoints(w *polygon.Winding) {
	var p [constants.MAX_POINTS_ON_WINDING]mgl32.Vec3
	var nump int

	for i := 0; i < w.NumPoints; i++ {
		j := (i + 1) % w.NumPoints
		k := (i + w.NumPoints - 1) % w.NumPoints
		v1 := w.Points[j].Sub(w.Points[i])
		v2 := w.Points[i].Sub(w.Points[k])
		v1 = v1.Normalize()
		v2 = v2.Normalize()

		if v1.Dot(v2) < 0.999 {
			p[nump] = w.Points[i]
			nump++
		}
	}

	if nump == w.NumPoints {
		return
	}

	// @TODO Replace with numthreads
	if constants.TEMPCONST_NUM_THREADS == 1 {
		c_removed += w.NumPoints - nump
	}
	w.NumPoints = nump

	tPoints := make([]mgl32.Vec3, nump)
	for i := 0; i < nump; i++ {
		tPoints[i] = p[i]
	}

	w.Points = tPoints
}
