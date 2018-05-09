package polygon

import "github.com/go-gl/mathgl/mgl32"

func GetEdgeEquation(p1 mgl32.Vec3, p2 mgl32.Vec3, c1 int, c2 int, insidePoint mgl32.Vec3) mgl32.Vec3 {
	nx := p1[c2] - p2[c2]
	ny := p2[c1] - p1[c1]
	d := -(nx * p1[c1] + ny * p1[c2])
	// 	assert(fabs(nx*p1[c1]+ny*p1[c2]+d)<0.01);
	// 	assert(fabs(nx*p2[c1]+ny*p2[c2]+d)<0.01);

	// use the convention that negative is "outside"
	trialDist := insidePoint[c1] * nx + insidePoint[c2] * ny + d
	if trialDist < 0 {
		nx = -nx
		ny = -ny
		d = -d
		trialDist = -trialDist
	}
	nx /= trialDist										// scale so that it will be =1.0 at the oppositve vertex
	ny /= trialDist
	d /= trialDist

	return mgl32.Vec3{nx,ny,d}
}
