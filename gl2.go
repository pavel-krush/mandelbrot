package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	_ "image/png"
	"log"
	"runtime"
	"unsafe"

	"mandelbrot/graph"
)

func init() {
	runtime.LockOSThread()
}

func GLDebugMessageCallback(source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {
	fmt.Printf("GLDebugMessage(gltype=%d, id=%d, severity=%d, length=%d, message=\"%s\", userParam=%p)\n",
		gltype, id, severity, length, message, userParam)
}

const (
	windowWidth  = 640
	windowHeight = 480
)

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "gl test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	proj := mgl32.Ortho(0.0, float32(windowWidth), float32(windowHeight), 0, -1, 1)

	shader, err := graph.NewShader("res/basic.shader")
	if err != nil {
		panic(err)
	}
	shader.Bind()
	shader.SetUniformMat4f("u_MVP", proj)

	texture, err := graph.NewTexture("res/texture.png")
	if err != nil {
		panic(err)
	}

	// Create fractal render object
	fractalObject, err := graph.NewObject2d(shader, texture)
	if err != nil {
		panic(err)
	}

	fractalObject.AddVertex(mgl32.Vec2{0, 0}, mgl32.Vec2{0.0, 0.0})
	fractalObject.AddVertex(mgl32.Vec2{100, 0}, mgl32.Vec2{1.0, 0.0})
	fractalObject.AddVertex(mgl32.Vec2{100, 100}, mgl32.Vec2{1.0, 1.0})
	fractalObject.AddVertex(mgl32.Vec2{0, 100}, mgl32.Vec2{0.0, 1.0})

	fractalObject.AddIndexBufferData(0, 1, 2)
	fractalObject.AddIndexBufferData(2, 3, 0)

	err = fractalObject.Compile()
	if err != nil {
		panic(err)
	}

	renderer := graph.NewRenderer()

	x := 0
	delta := 3

	maxX := windowWidth - 100

	for !window.ShouldClose() {
		renderer.Clear()

		renderer.Draw(fractalObject, x, windowHeight / 2 - 100 / 2)
		renderer.Draw(fractalObject, maxX - x, windowHeight / 2 - 100 / 2)

		x += delta
		if x < 0 || x > maxX {
			delta = -delta
			x += delta
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}

	fractalObject.Destroy()

	texture.Destroy()
	shader.Destroy()

	glfw.Terminate()
}
