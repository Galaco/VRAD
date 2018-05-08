package matrix

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Mat4 mgl32.Mat4

func (matrix *Mat4) SetupMatrixOrgAngles(origin *mgl32.Vec3, angles *mgl32.Vec3) {
	setupMatrixAnglesInternal( matrix, angles )

	// Add translation
	matrix[3] = origin[0]
	matrix[7] = origin[1]
	matrix[11] = origin[2]
	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = 0.0
	matrix[15] = 1.0
}

func setupMatrixAnglesInternal(matrix *Mat4, angles *mgl32.Vec3) {
	var sr, sp, sy, cr, cp, cy float64

	pitch := float64(mgl32.DegToRad(angles[0]))
	yaw := float64(mgl32.DegToRad(angles[1]))
	roll := float64(mgl32.DegToRad(angles[2]))
	sy = math.Sin(yaw)
	cy = math.Cos(yaw)
	sp = math.Sin(pitch)
	cp = math.Cos(pitch)
	sr = math.Sin(roll)
	cr = math.Cos(roll)

	// matrix = (YAW * PITCH) * ROLL
	matrix[0] = float32(cp*cy)
	matrix[4] = float32(cp*sy)
	matrix[8] = float32(-sp)
	matrix[1] = float32(sr*sp*cy+cr*-sy)
	matrix[5] = float32(sr*sp*sy+cr*cy)
	matrix[9] = float32(sr*cp)
	matrix[2] = float32((cr*sp*cy+-sr*-sy))
	matrix[6] = float32((cr*sp*sy+-sr*cy))
	matrix[10] = float32(cr*cp)
	matrix[3] = 0
	matrix[7] = 0
	matrix[11] = 0
}

func (matrix *Mat4) Mul4x3(vector *mgl32.Vec3) *mgl32.Vec3 {
	return vector3DMultiplyPosition(toInternal(matrix), vector)
}

func (matrix *Mat4) Identity() {
	(*matrix)[0] = 1.0
	(*matrix)[5] = 1.0
	(*matrix)[10] = 1.0
	(*matrix)[15] = 1.0

}

func vector3DMultiplyPosition(matrix *mgl32.Mat4, vector *mgl32.Vec3) *mgl32.Vec3 {
	return &mgl32.Vec3 {
		matrix[0] * vector[0] + matrix[1] * vector[1] + matrix[2] * vector[2] + matrix[3],
		matrix[4] * vector[0] + matrix[5] * vector[1] + matrix[6] * vector[2] + matrix[7],
		matrix[8] * vector[0] + matrix[9] * vector[1] + matrix[10] * vector[2] + matrix[11],
	}
}

// Internal mgl3 Mat4 has a load of functions we want.
// Maybe this is a little lazy, for now
func toInternal(matrix *Mat4) *mgl32.Mat4 {
	return &mgl32.Mat4 {
		matrix[0],
		matrix[1],
		matrix[2],
		matrix[3],
		matrix[4],
		matrix[5],
		matrix[6],
		matrix[7],
		matrix[8],
		matrix[9],
		matrix[10],
		matrix[11],
		matrix[12],
		matrix[13],
		matrix[14],
		matrix[15],
	}
}