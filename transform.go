package robb

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	X     float32
	Y     float32
	Z     float32
	Scale float32
	Angle float32
}

func NewTransform() *Transform {
	return &Transform{
		Scale: 1,
	}
}

func (t *Transform) GetTransformationMatrix() mgl32.Mat4 {
	translation := mgl32.Translate3D(t.X, t.Y, t.Z)
	rotation := mgl32.HomogRotate3DY(mgl32.DegToRad(t.Angle))
	scale := mgl32.Scale3D(t.Scale, t.Scale, t.Scale)
	model := translation.Mul4(rotation.Mul4(scale))

	return model
}
