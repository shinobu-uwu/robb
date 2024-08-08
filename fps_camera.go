package robb

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const yaw float32 = -90.0
const pitch float32 = 0.0
const speed float32 = 2.5
const sensitivity float32 = 0.1

type FpsCamera struct {
	Position  mgl32.Vec3
	Front     mgl32.Vec3
	Up        mgl32.Vec3
	Right     mgl32.Vec3
	WorldUp   mgl32.Vec3
	Fov       float32
	Aspect    float32
	Near      float32
	Far       float32
	yaw       float32
	pitch     float32
	deltaTime float32
}

func NewFpsCamera(aspectRatio, fov float32) *FpsCamera {
	c := FpsCamera{
		Position: mgl32.Vec3{0, 0, 3},
		Front:    mgl32.Vec3{0, 0, -1},
		WorldUp:  mgl32.Vec3{0, 1, 0},
		yaw:      yaw,
		pitch:    pitch,
		Fov:      fov,
		Aspect:   aspectRatio,
		Near:     0.1,
		Far:      100.0,
	}

	c.updateVectors()
	return &c
}

func (c *FpsCamera) ViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

func (c *FpsCamera) ProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(45), 4/3, 0.1, 100)
}

func (c *FpsCamera) ProcessKeyboard(movement CameraMovement) {
	velocity := speed * c.deltaTime

	if movement == Forward {
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	}

	if movement == Backward {
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	}

	if movement == Left {
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	}

	if movement == Right {
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}

	if movement == Up {
		c.Position = c.Position.Add(c.WorldUp.Mul(velocity))
	}

	if movement == Down {
		c.Position = c.Position.Sub(c.WorldUp.Mul(velocity))
	}
}

func (c *FpsCamera) ProcessMouseMovement(xOffset float32, yOffset float32) {
	xOffset *= sensitivity
	yOffset *= sensitivity

	c.yaw += xOffset
	c.pitch -= yOffset

	if c.pitch > 89 {
		c.pitch = 89
	}

	if c.pitch < -89 {
		c.pitch = -89
	}

	c.updateVectors()
}

func (c *FpsCamera) updateVectors() {
	direction := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(c.yaw))) * math.Cos(float64(mgl32.DegToRad(c.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.yaw))) * math.Cos(float64(mgl32.DegToRad(c.pitch)))),
	}
	c.Front = direction.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}
