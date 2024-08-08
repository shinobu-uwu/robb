package robb

import "github.com/go-gl/mathgl/mgl32"

type Position mgl32.Vec4
type Coordinates mgl32.Vec2

type Vertex struct {
	Position  Position
	TexCoords Coordinates
}
