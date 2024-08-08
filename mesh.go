package robb

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type DrawMode uint32

const (
	Triangles = gl.TRIANGLES
	Lines     = gl.LINES
)

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	Texture  *Texture
	Color    Color
	DrawMode DrawMode
	shader   *Shader
	vao      uint32
	vbo      uint32
	ebo      uint32
}

func NewMesh(vertices []Vertex, indices []uint32, texture *Texture) *Mesh {
	shader := TextureShader()
	if texture == nil {
		shader = StaticShader()
	}

	m := &Mesh{
		Vertices: vertices,
		Indices:  indices,
		Texture:  texture,
		DrawMode: Triangles,
		shader:   shader,
	}
	m.setup()

	return m
}

func (m *Mesh) Draw(view, projection mgl32.Mat4) {
	if m.Texture == nil {
		m.shader.SetVec4f("color", m.Color[0], m.Color[1], m.Color[2], m.Color[3])
	} else {
		m.Texture.Bind()
		m.shader.SetInt("ourTexture", 0)
	}

	m.shader.SetMat4("view", view)
	m.shader.SetMat4("projection", projection)

	gl.BindVertexArray(m.vao)
	gl.DrawElementsWithOffset(uint32(m.DrawMode), int32(len(m.Indices)), gl.UNSIGNED_INT, 0)

	gl.BindVertexArray(0)
}

func (m *Mesh) setup() {
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	vertexSize := int(unsafe.Sizeof(m.Vertices[0]))
	gl.BufferData(gl.ARRAY_BUFFER, len(m.Vertices)*vertexSize, gl.Ptr(m.Vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &m.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.Indices)*4, gl.Ptr(m.Indices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, int32(vertexSize), 0)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, int32(vertexSize), unsafe.Offsetof(m.Vertices[0].TexCoords))

	gl.BindVertexArray(0)
}
