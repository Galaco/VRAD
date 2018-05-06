package polygon

import (
	"github.com/galaco/bsp/primitives/common"
	"math"
	"log"
	"github.com/go-gl/mathgl/mgl32"
)

const MAX_POINTS_ON_WINDING = 64
const TEMPCONST_NUM_THREADS = 1

func BaseWindingForPlane(normal *mgl32.Vec3, dist float32) *common.Winding {
	var x int
	var v,max float32
	var org, vright, vup mgl32.Vec3
	var w *common.Winding

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
	//VectorMA (vup, -v, normal, vup);
	vup = vup.Normalize()
	//VectorScale (normal, dist, org);

	vright = vup.Cross(*normal)
	//VectorScale (vup, (MAX_COORD_INTEGER*4), vup);
	//VectorScale (vright, (MAX_COORD_INTEGER*4), vright);

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


func ChopWindingInPlace(inOut **common.Winding, normal *mgl32.Vec3, dist float32, epsilon float32) {
	var in *common.Winding
	var dists [MAX_POINTS_ON_WINDING + 4]float32
	var sides [MAX_POINTS_ON_WINDING + 4]int
	counts := [3]int{0,0,0}
	var dot float32
	var i, j int
	mid := mgl32.Vec3{0,0,0}
	var f *common.Winding
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
		//FreeWinding(in)
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
	if f.NumPoints > MAX_POINTS_ON_WINDING {
		log.Fatal("ClipWinding: MAX_POINTS_ON_WINDING")
	}

	//FreeWinding(in)
	*inOut = f
}

func freeWinding(w *common.Winding) {
	//if w.NumPoints == 0xdeaddead {
	//	log.Fatal("FreeWinding: freed a freed winding")
	//}

	// ThreadLock()
	//w.NumPoints = 0xdeaddead
/**
if (w->numpoints == 0xdeaddead)
		Error ("FreeWinding: freed a freed winding");

	ThreadLock();
	w->numpoints = 0xdeaddead; // flag as freed
	w->next = winding_pool[w->maxpoints];
	winding_pool[w->maxpoints] = w;
	ThreadUnlock();
 */
}


// @TODO This function does a load of stuff to externals.
func NewWinding(n int) *common.Winding {
	return &common.Winding{}
/**
winding_t	*w;

	if (numthreads == 1)
	{
		c_winding_allocs++;
		c_winding_points += points;
		c_active_windings++;
		if (c_active_windings > c_peak_windings)
			c_peak_windings = c_active_windings;
	}
	ThreadLock();
	if (winding_pool[points])
	{
		w = winding_pool[points];
		winding_pool[points] = w->next;
	}
	else
	{
		w = (winding_t *)malloc(sizeof(*w));
		w->p = (Vector *)calloc( points, sizeof(Vector) );
	}
	ThreadUnlock();
	w->numpoints = 0; // None are occupied yet even though allocated.
	w->maxpoints = points;
	w->next = NULL;
	return w;
 */
}