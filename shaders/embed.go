package shaders

import _ "embed"

//go:embed texture_vertex.glsl
var TextureVertexShaderSrc []byte

//go:embed texture_fragment.glsl
var TextureFragmentShaderSrc []byte

//go:embed static_vertex.glsl
var StaticVertexShaderSrc []byte

//go:embed static_fragment.glsl
var StaticFragmentShaderSrc []byte
