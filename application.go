package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"log"
	"mandelbrot/fractal"
	"mandelbrot/graph"
	"runtime"
	"sync"
	"time"
)

func init() {
	runtime.LockOSThread()
}

type Application struct {
	mu sync.Mutex

	window      *glfw.Window
	windowTitle string

	state *State

	renderer *graph.Renderer

	fractalObject  *graph.Object2D // Simple textured rectangle. Fractal will be rendered here
	fractalImg     *image.RGBA     // The image object where fractal will be drawn
	fractalTexture *graph.Texture  // OpenGL texture which will be rendered on an fractalObject
	shader         *graph.Shader   // The main shader

	refreshTexture   bool // Do we need to refresh opengl texture from the buffer

	generating bool // is generating in progress

	cursorPos struct {
		x  float64
		y  float64
	}

	generator fractal.Generator // Current fractal generator
	zoomer    Zoomer
}

func NewApplication(windowTitle string) *Application {
	state := NewState()

	ret := &Application{
		windowTitle: windowTitle,
		state:       state,
	}

	return ret
}

func (a *Application) Lock() {
	a.mu.Lock()
}

func (a *Application) Unlock() {
	a.mu.Unlock()
}

func (a *Application) scheduleRefreshTexture() {
	a.Lock()
	a.refreshTexture = true
	a.Unlock()
}

func (a *Application) clearRefreshTexture() {
	a.Lock()
	a.refreshTexture = false
	a.Unlock()
}

func (a *Application) needRefreshTexture() bool {
	a.Lock()
	defer a.Unlock()
	return a.refreshTexture
}

func (a *Application) startGenerating() {
	a.Lock()
	a.generating = true
	a.Unlock()
}

func (a *Application) stopGenerating() {
	a.Lock()
	a.generating = false
	a.Unlock()
}

func (a *Application) isGenerating() bool {
	a.Lock()
	defer a.Unlock()
	return a.generating
}

func (a *Application) RegenerateFractal() {
	a.startGenerating()

	fmt.Println("Generating fractal")
	fmt.Println(a.state)

	lastUpdated := time.Now()
	lastUpdatedPtr := &lastUpdated

	started := time.Now()

	progress := func(progress float32) {
		now := time.Now()
		if now.Sub(*lastUpdatedPtr) > time.Millisecond*100 {
			a.scheduleRefreshTexture()
		}
	}

	done := func() {
		a.scheduleRefreshTexture()
		a.stopGenerating()

		end := time.Now()
		genTime := end.Sub(started)
		fmt.Printf("Generation time: %s\n", genTime)
	}

	a.generator.Generate(
		a.fractalImg,
		a.state.GetCX(),
		a.state.GetCY(),
		a.state.GetScale(),
		a.state.GetPhysicalWidth(),
		a.state.GetPhysicalHeight(),
		progress,
		done,
	)
}

func (a *Application) SetGenerator(generator fractal.Generator) {
	a.Lock()
	a.generator = generator
	a.Unlock()
}

func (a *Application) Run() {
	if a.generator == nil {
		panic("generator not set")
	}

	a.Start()

	var fps graph.FPS

	// refresh texture from buffer every 100ms
	// todo: refactor it later somehow.
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		<-ticker.C
		if a.isGenerating() {
			a.scheduleRefreshTexture()
		}
	}()

	for !a.window.ShouldClose() {
		// Refresh GL texture from buffer if requested to do so
		if a.needRefreshTexture() {
			a.fractalTexture.SetImageData(a.fractalImg.Pix)
			a.clearRefreshTexture()
		}

		// Clear scene
		a.renderer.Clear()

		// Render fractal
		a.renderer.Draw(a.fractalObject, 0, 0)

		a.window.SwapBuffers()
		glfw.PollEvents()

		fps.FrameRendered()
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
	a.window, err = glfw.CreateWindow(int(a.state.GetScreenWidth()), int(a.state.GetScreenHeight()), a.windowTitle, nil, nil)
	if err != nil {
		panic(err)
	}
	a.window.MakeContextCurrent()
	a.window.SetMouseButtonCallback(a.MouseButtonCallback)
	a.window.SetCursorPosCallback(a.CursorPosCallback)

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
	proj := mgl32.Ortho(0.0, float32(a.state.GetScreenWidth()), float32(a.state.GetScreenHeight()), 0, -1, 1)

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
		Max: image.Point{X: int(a.state.GetScreenWidth()), Y: int(a.state.GetScreenHeight())},
	})

	// Create the texture which fractal will be rendered on
	a.fractalTexture = graph.NewTexture(int(a.state.GetScreenWidth()), int(a.state.GetScreenHeight()))

	// Create fractal render object
	a.fractalObject, err = graph.NewObject2d(a.shader, a.fractalTexture)
	if err != nil {
		panic(err)
	}

	a.fractalObject.AddVertex(mgl32.Vec2{0, 0}, mgl32.Vec2{0.0, 0.0})
	a.fractalObject.AddVertex(mgl32.Vec2{float32(a.state.GetScreenWidth()), 0}, mgl32.Vec2{1.0, 0.0})
	a.fractalObject.AddVertex(mgl32.Vec2{float32(a.state.GetScreenWidth()), float32(a.state.GetScreenHeight())}, mgl32.Vec2{1.0, 1.0})
	a.fractalObject.AddVertex(mgl32.Vec2{0, float32(a.state.GetScreenHeight())}, mgl32.Vec2{0.0, 1.0})
	a.fractalObject.AddIndexBufferData(0, 1, 2)
	a.fractalObject.AddIndexBufferData(2, 3, 0)

	err = a.fractalObject.Compile()
	if err != nil {
		panic(err)
	}

	// Setup renderer
	a.renderer = graph.NewRenderer()

	a.zoomer = NewZoomerSimple()

	a.RegenerateFractal()
}

func (a *Application) Terminate() {
	a.fractalObject.Destroy()
	a.fractalTexture.Destroy()
	a.shader.Destroy()
}

func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Release {
		a.Lock()
		x := a.cursorPos.x
		y := a.cursorPos.y
		a.Unlock()
		a.OnClick(x, y, button)
	}
}

func (a *Application) CursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	a.Lock()
	a.cursorPos.x = xpos
	a.cursorPos.y = ypos
	a.Unlock()
}

func (a *Application) OnClick(x float64, y float64, button glfw.MouseButton) {
	if a.generating {
		return
	}

	var direction ZoomDirection
	if button == glfw.MouseButton1 {
		direction = ZoomDirectionIn
	} else if button == glfw.MouseButton2 {
		direction = ZoomDirectionOut
	}

	a.zoomer.ZoomAt(a.state, x, y, direction)

	a.RegenerateFractal()
}
