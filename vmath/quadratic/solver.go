package quadratic


func SolveInverseQuadraticMonotonic(x1 float32, y1 float32, x2 float32, y2 float32, x3 float32, y3 float32, a *float32, b *float32, c *float32) bool {
	// use SolveInverseQuadratic, but if the sigm of the derivative at the start point is the wrong
	// sign, displace the mid point

	// first, sort parameters
	if x1 > x2 {
		vSwap(&x1, &x2)
		vSwap(&y1, &y2)
	}
	if x2 > x3 {
		vSwap(&x2, &x3)
		vSwap(&y2, &y3)
	}
	if x1 > x2 {
		vSwap(&x1, &x2)
		vSwap(&y1, &y2)
	}
	// this code is not fast. what it does is when the curve would be non-monotonic, slowly shifts
	// the center point closer to the linear line between the endpoints. Should anyone need htis
	// function to be actually fast, it would be fairly easy to change it to be so.
	for blend_to_linear_factor := 0.0; blend_to_linear_factor <= 1.0; blend_to_linear_factor += 0.05 {
		tempy2 := float64(1 - blend_to_linear_factor) * float64(y2) + blend_to_linear_factor * FLerp(y1,y3,x1,x3,x2)
		if !SolveInverseQuadratic(x1, y1, x2, float32(tempy2), x3, y3, a, b, c) {
			return false
		}
		derivative := 2.0 * (*a) + (*b)
		if (y1 < y2) && (y2 < y3) {							// monotonically increasing
			if derivative >= 0.0 {
				return true
			}
		} else {
			if (y1 > y2) && (y2 > y3) {							// monotonically decreasing
				if derivative <= 0.0 {
					return true
				}
			} else {
				return true
			}
		}
	}
	return true
}

// solves for "a, b, c" where "a x^2 + b x + c = y", return true if solution exists
func SolveInverseQuadratic(x1 float32, y1 float32, x2 float32, y2 float32, x3 float32, y3 float32, a *float32, b *float32, c *float32) bool {
	det := (x1 - x2)*(x1 - x3)*(x2 - x3)

	// FIXME: check with some sort of epsilon
	if det == 0.0 {
		return false
	}

	*a = (x3*(-y1 + y2) + x2*(y1 - y3) + x1*(-y2 + y3)) / det

	*b = (x3*x3*(y1 - y2) + x1*x1*(y2 - y3) + x2*x2*(-y1 + y3)) / det

	*c = (x1*x3*(-x1 + x3)*y2 + x2*x2*(x3*y1 - x1*y3) + x2*(-(x3*x3*y1) + x1*x1*y3)) / det

	return true
}

func vSwap(a *float32, b *float32) {
	c := *a
	a = b
	b = &c
}