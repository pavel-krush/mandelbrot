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

	// Create fractal render object
	fractalObject := graph.NewObject2d()
	fractalObject.AddVertex(mgl32.Vec2{0, 0}, mgl32.Vec2{0.0})
	fractalObject.AddVertex(mgl32.Vec2{windowWidth, 0}, mgl32.Vec2{1.0, 0.0})
	fractalObject.AddVertex(mgl32.Vec2{windowWidth, windowHeight}, mgl32.Vec2{1.0, 1.0})
	fractalObject.AddVertex(mgl32.Vec2{0, windowHeight}, mgl32.Vec2{0.0, 1.0})
	fractalObject.Compile()

	return

	positions := []float32{
		0, 0, 0.0, 0.0,
		windowWidth, 0, 1.0, 0.0,
		windowWidth, windowHeight, 1.0, 1.0,
		0, windowHeight, 0.0, 1.0,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	va := graph.NewVertexArray()
	vb := graph.NewVertexBuffer(gl.Ptr(positions), 4*4*int(unsafe.Sizeof(float32(0))))
	layout := graph.NewVertexBufferLayout()
	layout.PushFloat(2)
	layout.PushFloat(2)
	va.AddBuffer(vb, layout)

	ib := graph.NewIndexBuffer(indices)

	proj := mgl32.Ortho(0.0, float32(windowWidth), 0, float32(windowHeight), -1.0, 1.0)

	shader, err := graph.NewShader("res/basic.shader")
	if err != nil {
		panic(err)
	}

	texture, err := graph.NewTexture("res/texture.png")
	if err != nil {
		panic(err)
	}

	renderer := graph.NewRenderer()

	for !window.ShouldClose() {
		renderer.Clear()

		fractalObject.Draw()

		shader.Bind()
		texture.Bind(0)
		shader.SetUniform1i("u_Texture", 0)
		shader.SetUniformMat4f("u_MVP", proj)
		renderer.Draw(va, ib, shader)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	texture.Destroy()
	vb.Destroy()
	ib.Destroy()
	shader.Destroy()

	glfw.Terminate()
}
