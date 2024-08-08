package robb

import "github.com/go-gl/mathgl/mgl32"

type Drawable interface {
	Draw(view, projection mgl32.Mat4)
}

type CameraMovement uint32

type Camera interface {
	ViewMatrix() mgl32.Mat4
	ProjectionMatrix() mgl32.Mat4
	ProcessKeyboard(movement CameraMovement)
	ProcessMouseMovement(xOffset float32, yOffset float32)
}

const (
	Forward = iota
	Backward
	Left
	Right
	Up
	Down
)
