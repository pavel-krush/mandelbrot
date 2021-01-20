package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"strings"
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

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window, err := glfw.CreateWindow(640, 480, "gl test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	positions := []float32{
		-0.5, -0.5,
		0.5, -0.5,
		0.5, 0.5,
		-0.5, 0.5,
	}

	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	va := graph.NewVertexArray()
	vb := graph.NewVertexBuffer(gl.Ptr(positions), 4 * 2 * int(unsafe.Sizeof(float32(0))))

	layout := graph.NewVertexBufferLayout()
	layout.PushFloat(2)

	va.AddBuffer(vb, layout)

	ib := graph.NewIndexBuffer(indices)

	pg, err := glProgram(vs, fs)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(pg)

	uColorLocation := gl.GetUniformLocation(pg, gl.Str("u_Color\x00"))

	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	r := float32(0.0)
	increment := float32(0.05)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(pg)
		gl.Uniform4f(uColorLocation, r, 0.3, 0.8, 1.0)

		va.Bind()
		ib.Bind()

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		err := gl.GetError()
		if err != gl.NO_ERROR {
			fmt.Printf("gl error: %d\n", err)
		}

		if r > 1.0 {
			increment = -0.05
		} else if r < 0.0 {
			increment = 0.05
		}

		r += increment

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteProgram(pg)

	vb.Destroy()
	ib.Destroy()

	glfw.Terminate()
}

var vs = `
#version 330 core

layout(location=0) in vec4 position;

void main() {
	gl_Position = position;
}`

var fs = `
#version 330 core

layout(location=0) out vec4 color;

uniform vec4 u_Color;

void main() {
	color = u_Color;
}`

func glProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := glShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := glShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func glShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
