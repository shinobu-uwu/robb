package robb

import "github.com/go-gl/mathgl/mgl32"

type GameObject struct {
	Mesh
	Transform
}

func NewGameObject(vertices []Vertex, indices []uint32, texture *Texture) *GameObject {
	mesh := NewMesh(vertices, indices, texture)

	return &GameObject{
		Mesh:      *mesh,
		Transform: *NewTransform(),
	}
}

func (g *GameObject) Draw(view, projection mgl32.Mat4) {
	g.Mesh.shader.SetMat4("model", g.GetTransformationMatrix())
	g.Mesh.Draw(view, projection)
	g.Mesh.shader.SetMat4("model", mgl32.Ident4())
}
