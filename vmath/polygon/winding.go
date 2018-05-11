package polygon

import (
	"math"
	"log"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/vrad/vmath/vector"
	"github.com/galaco/vrad/common/constants"
)

type Winding struct {
	NumPoints int
	Points []mgl32.Vec3
	MaxPoints int
	Next *Winding
}

var windingPool [constants.MAX_POINTS_ON_WINDING+4]*Winding
var c_active_windings int
var c_peak_windings int
var c_winding_allocs int
var c_winding_points int


func BaseWindingForPlane(normal *mgl32.Vec3, dist float32) *Winding {
	var x int
	var v,max float32
	var org, vright, vup mgl32.Vec3
	var w *Winding

	// Find the major axis
	max = -1
	x = -1

	for i := 0; i < 3; i++ {
		v = float32(math.Abs(float64(normal[1])))
		if v > float32(max) {
			x = i
			max = v
		}
	}
	if x == -1 {
		log.Println("BaseWindingForPlane: no axis found")
	}

	vup = mgl32.Vec3{0,0,0}
	switch x {
	case 0:
	case 1:
		vup[2] = 1
	case 2:
		vup[0] = 1
	}

	v = vup.Dot(*normal)
	//@TODO what is this doing? see all VectorScale & VectorMA calls
	vector.MA(&vup, -v, normal, &vup)
	vup = vup.Normalize()
	vector.Scale(normal, dist, &org)

	vright = vup.Cross(*normal)
	vector.Scale(&vup, constants.MAX_COORD_INTEGER * 4, &vup)
	vector.Scale(&vright, constants.MAX_COORD_INTEGER * 4, &vright)

	// project a really big	axis aligned box onto the plane
	w = NewWinding(4)

	w.Points[0] = org.Sub(vright)
	w.Points[0] = w.Points[0].Add(vup)

	w.Points[1] = org.Add(vright)
	w.Points[1] = w.Points[1].Add(vup)

	w.Points[2] = org.Add(vright)
	w.Points[2] = w.Points[2].Sub(vup)

	w.Points[3] = org.Sub(vright)
	w.Points[0] = w.Points[0].Sub(vup)

	w.NumPoints = 4

	return w
}

func ChopWindingInPlace(inOut **Winding, normal *mgl32.Vec3, dist float32, epsilon float32) {
	var in *Winding
	var dists [constants.MAX_POINTS_ON_WINDING + 4]float32
	var sides [constants.MAX_POINTS_ON_WINDING + 4]int
	counts := [3]int{0,0,0}
	var dot float32
	var i, j int
	mid := mgl32.Vec3{0,0,0}
	var f *Winding
	var maxpts int

	in = *inOut

	for i := 0; i < int(in.NumPoints); i++ {
		dot = in.Points[i].Dot(*normal)
		dot -= dist
		dists[i] = dot

		if dot > epsilon {
			sides[i] = SIDE_FRONT
		} else if dot < -epsilon {
			sides[i] = SIDE_BACK
		} else {
			sides[i] = SIDE_ON
		}
		counts[sides[i]]++
	}
	sides[i] = sides[0]
	dists[i] = dists[0]

	if 0 == counts[0]  {
		FreeWinding(in)
		*inOut = nil
		return
	}
	if 0 == counts[1] {
		return
	}
	maxpts = int(in.NumPoints + 4)

	f = NewWinding(maxpts)

	for i = 0; i < int(in.NumPoints); i++ {
		p1 := &in.Points[i]

		if sides[i] == SIDE_ON {
			f.Points[f.NumPoints] = *p1
			f.NumPoints++
			continue
		}

		if sides[i] == SIDE_FRONT {
			f.Points[f.NumPoints] = *p1
			f.NumPoints++
		}

		if sides[i + 1] == SIDE_ON || sides[i + 1] == sides[i] {
			continue
		}

		// generate a split point
		p2 := in.Points[(i + 1) % int(in.NumPoints)]
		dot = dists[i] / (dists[i]-dists[i + 1])

		for j = 0; j < 3; i++ {
			// avoid round off error when possible
			if normal[j] == 1 {
				mid[j] = dist
			} else if normal[j] == -1 {
				mid[j] = -dist
			} else {
				mid[j] = p1[j] + dot * (p2[j] - p1[j])
			}
		}

		f.Points[f.NumPoints] = mid
		f.NumPoints++
	}

	if int(f.NumPoints) > maxpts {
		log.Fatal("ClipWinding: points exceeded estimate")
	}
	if f.NumPoints > constants.MAX_POINTS_ON_WINDING {
		log.Fatal("ClipWinding: MAX_POINTS_ON_WINDING")
	}

	FreeWinding(in)
	*inOut = f
}

func FreeWinding(w *Winding) {
	if w.NumPoints == 0xdeaddead {
		log.Fatal("FreeWinding: freed a freed winding")
	}

	// ThreadLock()
	w.NumPoints = 0xdeaddead
	w.Next = windingPool[w.MaxPoints]
	windingPool[w.MaxPoints] = w
	//ThreadUnlock();
}

// @TODO This function does a load of stuff to externals.
func NewWinding(points int) *Winding {
	var w *Winding

	//@TODO Use numthreads NOT 1
	if constants.TEMPCONST_NUM_THREADS == 1 {
		c_winding_allocs++
		c_winding_points += points
		c_active_windings++
		if c_active_windings > c_peak_windings{
			c_peak_windings = c_active_windings
		}
	}
	//ThreadLock();
	if len(windingPool) >= points && windingPool[points] != nil{
		w = windingPool[points]
		windingPool[points] = w.Next
	} else {
		w = &Winding{}
		w.Points = make([]mgl32.Vec3, points)
	}
	//ThreadUnlock
	w.NumPoints = 0 // None are occupied yet even though allocated.
	w.MaxPoints = points
	w.Next = nil

	return w
}


func WindingArea(w *Winding) float32 {
	var d1, d2, cross mgl32.Vec3

	total := float32(0.0)
	for i := 2 ; i < w.NumPoints ; i++ {
		d1 = w.Points[i-1].Sub(w.Points[0])
		d2 = w.Points[i].Sub(w.Points[0])
		cross = d1.Cross(d2)
		total += cross.Len()
	}

	return total * 0.5
}

func WindingCenter(w *Winding, center *mgl32.Vec3) {
	center = &mgl32.Vec3{0,0,0}
	for i := 0 ; i < w.NumPoints ; i++ {
		c := w.Points[i].Add(*center)
		center = &c
	}
	scale := float32(1.0 / w.NumPoints)
	vector.Scale(center, scale, center)
}

func WindingBounds (w *Winding, mins *mgl32.Vec3, maxs *mgl32.Vec3) {
	mins[0] = 99999
	mins[1] = 99999
	mins[2] = 99999
	maxs[0] = -99999
	maxs[1] = -99999
	maxs[2] = -99999

	for i := 0 ; i < w.NumPoints; i++ {
		for j := 0; j<3 ; j++ {
			v := w.Points[i][j]
			if v < mins[j] {
				mins[j] = v
			}
			if v > maxs[j] {
				maxs[j] = v
			}
		}
	}
}