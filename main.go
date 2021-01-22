package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	_ "image/png"
	"log"
	"mandelbrot/fractal"
	"mandelbrot/palette"
	"runtime"
	"unsafe"

	"mandelbrot/graph"
)

func init() {
	runtime.LockOSThread()
}

const (
	windowWidth = 640
	windowHeight = 480
	initPhysWidth = 3.0
	initPhysHeight = 2.0
)

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

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: windowWidth, Y: windowHeight},
	})

	GenerateMandelbrot(img, -0.7, 0, 1)

	texture := graph.NewTexture(windowWidth, windowHeight)
	texture.SetImageData(img.Pix)

	renderer := graph.NewRenderer()

	for !window.ShouldClose() {
		renderer.Clear()

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
}
