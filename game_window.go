package robb

import (
	"log"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GameWindow struct {
	BackgroundColor Color
	Scene           *Scene
	Width           int
	Height          int
	keybindings     map[glfw.Key]func()
	keyStates       map[glfw.Key]glfw.Action
	window          *glfw.Window
}

func init() {
	runtime.LockOSThread()
}

func NewGameWindow(width int, height int, title string) (*GameWindow, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, err
	}

	gl.Enable(gl.DEPTH_TEST)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	if os.Getenv("ROBB_DEBUG") == "1" {
		setupDebugCallback()
	}

	return &GameWindow{
		Width:       width,
		Height:      height,
		window:      window,
		keybindings: make(map[glfw.Key]func()),
		keyStates:   make(map[glfw.Key]glfw.Action),
	}, nil
}

func setupDebugCallback() {
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	gl.DebugMessageCallback(func(source uint32, type_ uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
		log.Printf("OpenGL Debug Message: %s\n", message)
	}, nil)
}

func (gw *GameWindow) Close() {
	gw.window.SetShouldClose(true)
}

func (gw *GameWindow) Start() {
	gw.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		gw.keyStates[key] = action
	})

	var lastX, lastY float64 = float64(gw.Width) / 2, float64(gw.Height) / 2
	firstMouse := true

	gw.window.SetCursorPosCallback(func(window *glfw.Window, xpos float64, ypos float64) {
		if firstMouse {
			lastX = xpos
			lastY = ypos
			firstMouse = false
		}

		xOffset := xpos - lastX
		yOffset := ypos - lastY

		lastX = xpos
		lastY = ypos

		gw.Scene.Camera.ProcessMouseMovement(float32(xOffset), float32(yOffset))
	})

	var deltaTime, lastFrame float32
	defer glfw.Terminate()

	for !gw.window.ShouldClose() {
		gl.ClearColor(gw.BackgroundColor[0], gw.BackgroundColor[1], gw.BackgroundColor[2], gw.BackgroundColor[3])
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for key, action := range gw.keyStates {
			if action == glfw.Press || action == glfw.Repeat {
				if handler, found := gw.keybindings[key]; found {
					handler()
				}
			}
		}

		if gw.Scene != nil {
			currentFrame := float32(glfw.GetTime())
			deltaTime = currentFrame - lastFrame
			lastFrame = currentFrame
			gw.Scene.Camera.deltaTime = deltaTime
			gw.Scene.Draw()
		}

		gw.window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (w *GameWindow) AddKeybinding(key glfw.Key, handler func()) {
	w.keybindings[key] = handler
	w.keyStates[key] = glfw.Release
}
