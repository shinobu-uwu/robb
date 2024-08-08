package robb

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/shinobu-uwu/robb/shaders"
)

type Shader struct {
	Program uint32
}

var textureShader *Shader
var staticShader *Shader

func TextureShader() *Shader {
	if textureShader == nil {
		shader, err := NewShader(string(shaders.TextureVertexShaderSrc), string(shaders.TextureFragmentShaderSrc))
		if err != nil {
			panic(err)
		}

		textureShader = shader
	}

	return textureShader
}

func StaticShader() *Shader {
	if staticShader == nil {
		shader, err := NewShader(string(shaders.StaticVertexShaderSrc), string(shaders.StaticFragmentShaderSrc))
		if err != nil {
			panic(err)
		}

		staticShader = shader
	}

	return staticShader
}

func LoadShaderFromFile(vertexPath, fragmentPath string) (*Shader, error) {
	vertexBytes, err := os.ReadFile(vertexPath)
	if err != nil {
		return nil, err
	}

	fragmentBytes, err := os.ReadFile(fragmentPath)
	if err != nil {
		return nil, err
	}

	return NewShader(string(vertexBytes), string(fragmentBytes))
}

func NewShader(vertexSrc, fragmentSrc string) (*Shader, error) {
	var success int32
	vertex := gl.CreateShader(gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertex)
	str, free := gl.Strs(vertexSrc)
	gl.ShaderSource(vertex, 1, str, nil)
	gl.CompileShader(vertex)
	free()

	gl.GetShaderiv(vertex, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertex, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength))
		gl.GetShaderInfoLog(vertex, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to compile vertex shader: %v", log)
	}

	fragment := gl.CreateShader(gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragment)
	str, free = gl.Strs(fragmentSrc)
	gl.ShaderSource(fragment, 1, str, nil)
	gl.CompileShader(fragment)
	free()

	gl.GetShaderiv(fragment, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(fragment, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength))
		gl.GetShaderInfoLog(fragment, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to compile fragment shader: %v", log)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)

	gl.GetProgramiv(program, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	return &Shader{Program: program}, nil
}

func (s *Shader) Use() {
	gl.UseProgram(s.Program)
}

func (s *Shader) SetBool(name string, value bool) {
	var v int32
	if value {
		v = 1
	}

	s.Use()
	gl.Uniform1i(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), v)
}

func (s *Shader) SetInt(name string, value int32) {
	s.Use()
	gl.Uniform1i(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetFloat(name string, value float32) {
	s.Use()
	gl.Uniform1f(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), value)
}

func (s *Shader) SetVec4f(name string, x, y, z, w float32) {
	s.Use()
	gl.Uniform4f(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), x, y, z, w)
}

func (s *Shader) SetMat4(name string, value mgl32.Mat4) {
	s.Use()
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")), 1, false, &value[0])
}
