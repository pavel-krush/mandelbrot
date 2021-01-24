package main

import (
	"fmt"
	"image"
	"log"
	"mandelbrot/fractal"
	"mandelbrot/graph"
	"mandelbrot/palette"
	"runtime"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	initPhysWidth = 3.0
	initPhysHeight = 2.0
)

func init() {
	runtime.LockOSThread()
}

type Application struct {
	window *glfw.Window
	windowWidth int
	windowHeight int
	windowTitle string

	state *State

	renderer *graph.Renderer

	fractalObject *graph.Object2D // Simple textured rectangle. Fractal will be rendered here
	fractalImg *image.RGBA // The image object where fractal will be drawn
	fractalTexture *graph.Texture // OpenGL texture which will be rendered on an fractalObject
	shader *graph.Shader // The main shader
}

func NewApplication(windowWidth, windowHeight int, windowTitle string) *Application {
	ret := &Application{
		windowWidth: windowWidth,
		windowHeight: windowHeight,
		windowTitle: windowTitle,

		state: NewState(),
	}

	return ret
}

func (a *Application) RegenerateFractal() {
	cx, cy, cscale := a.state.GetCoords()

	x, _ := cx.Float64()
	y, _ := cy.Float64()
	scale, _ := cscale.Float64()

	GenerateMandelbrot(a.fractalImg, x, y, scale)
	a.fractalTexture.SetImageData(a.fractalImg.Pix)
}

func (a *Application) Run() {
	a.Start()

	for !a.window.ShouldClose() {
		// Clear scene
		a.renderer.Clear()

		// Render fractal
		a.renderer.Draw(a.fractalObject, 0, 0)

		a.window.SwapBuffers()
		glfw.PollEvents()
	}

	a.Terminate()

	defer glfw.Terminate()
}

// Initialize application
func (a *Application) Start() {
	var err error

	// Initialize glfw
	if err = glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	// Create window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	a.window, err = glfw.CreateWindow(a.windowWidth, a.windowHeight, a.windowTitle, nil, nil)
	if err != nil {
		panic(err)
	}
	a.window.MakeContextCurrent()

	glfw.SwapInterval(1)

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println("OpenGL version", gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Create projection matrix that moves zero coordinates to left top corner and
	// scales width and height to the window size
	proj := mgl32.Ortho(0.0, float32(a.windowWidth), float32(a.windowHeight), 0, -1, 1)

	// Load main shader
	a.shader, err = graph.NewShader("res/basic.shader")
	if err != nil {
		panic(err)
	}
	a.shader.Bind()
	a.shader.SetUniformMat4f("u_MVP", proj)

	// Create image buffer, the target for a fractal generating functions
	a.fractalImg = image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: a.windowWidth, Y: a.windowHeight},
	})

	// Create the texture which fractal will be rendered on
	a.fractalTexture = graph.NewTexture(a.windowWidth, a.windowHeight)

	// Create fractal render object
	a.fractalObject, err = graph.NewObject2d(a.shader, a.fractalTexture)
	if err != nil {
		panic(err)
	}

	a.fractalObject.AddVertex(mgl32.Vec2{0, 0}, mgl32.Vec2{0.0, 0.0})
	a.fractalObject.AddVertex(mgl32.Vec2{float32(a.windowWidth), 0}, mgl32.Vec2{1.0, 0.0})
	a.fractalObject.AddVertex(mgl32.Vec2{float32(a.windowWidth), float32(a.windowHeight)}, mgl32.Vec2{1.0, 1.0})
	a.fractalObject.AddVertex(mgl32.Vec2{0, float32(a.windowHeight)}, mgl32.Vec2{0.0, 1.0})
	a.fractalObject.AddIndexBufferData(0, 1, 2)
	a.fractalObject.AddIndexBufferData(2, 3, 0)

	err = a.fractalObject.Compile()
	if err != nil {
		panic(err)
	}

	// Setup renderer
	a.renderer = graph.NewRenderer()

	// TOOD: remove
	a.RegenerateFractal()
}

func (a *Application) Terminate() {
	a.fractalObject.Destroy()
	a.fractalTexture.Destroy()
	a.shader.Destroy()
}

func GenerateMandelbrot(target *image.RGBA, cx float64, cy float64, scale float64) {
	width := target.Rect.Max.X
	height := target.Rect.Max.Y

	// Calculate physical width and height
	physWidth := initPhysWidth * scale
	physHeight := initPhysHeight * scale

	// Scale physical bounds
	physMinX := cx - (physWidth / 2)
	physMinY := cy - (physHeight / 2)

	// Calculate pixel-to-physical scale
	scaleX := physWidth / float64(width)
	scaleY := physHeight / float64(height)

	pal := palette.CreatePaletteGrayscaleRecursive(256)

	// (x, y) - are pixel coords
	for y := 0; y < height; y++ {
		// (physX, physY) - are physical coordinates
		physY := float64(y)*scaleY + physMinY
		for x := 0; x < width; x++ {
			physX := float64(x)*scaleX + physMinX

			// get fractal value at the point
			value := fractal.Mandelbrot(complex(physX, physY))

			// convert it to the color and set pixel color
			target.Set(x, y, pal[int(float64(len(pal)) * value)])
		}
	}
}
