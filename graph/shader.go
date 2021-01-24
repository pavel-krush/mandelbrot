package graph

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"

	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/pkg/errors"
)

type Shader struct {
	filePath string
	vertexSource string
	fragmentSource string

	rendererId uint32

	uniformCache map[string]int32
}

func NewShader(filePath string) (*Shader, error) {
	var err error
	ret := &Shader{
		filePath: filePath,
		uniformCache: make(map[string]int32),
	}

	err = ret.readSource()
	if err != nil {
		return nil, err
	}

	err = ret.compile()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *Shader) Destroy() {
	gl.DeleteProgram(s.rendererId)
}

func (s *Shader) Bind() {
	gl.UseProgram(s.rendererId)
}

func (s *Shader) Unbind() {
	gl.UseProgram(0)
}

// read shader source file
// both vertex and fragment shaders must be defined in source file
func (s *Shader) readSource() error {
	type scanMode int
	const (
		modeUnknown scanMode = iota
		modeVertex
		modeFragment
	)

	mode := modeUnknown

	// read all the file
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	// scan data line by line
	for scanner.Scan() {
		line := scanner.Text()
		// check if we have shader type in current line
		if strings.HasPrefix(line, "#shader") {
			// parse shader type and set appropriate parse mode
			if strings.Contains(line, "vertex") {
				mode = modeVertex
			} else if strings.Contains(line, "fragment") {
				mode = modeFragment
			} else {
				return fmt.Errorf("unknown shader type: " + line)
			}
			continue
		}

		// append shader source to current shader
		switch mode {
		case modeVertex:
			s.vertexSource += line + "\n"
		case modeFragment:
			s.fragmentSource += line + "\n"
		default:
			return fmt.Errorf("please define a shader type first")
		}
	}

	// check both shader types are defined
	if s.vertexSource == "" {
		return fmt.Errorf("no vertex shader source given")
	}

	if s.fragmentSource == "" {
		return fmt.Errorf("no fragment shader source given")
	}

	return nil
}

// compile vertex and fragment shaders and link a program
func (s *Shader) compile() error {
	vertexShaderId, err := s.compileShader(gl.VERTEX_SHADER, s.vertexSource)
	if err != nil {
		return errors.Wrapf(err, "vertex shader")
	}

	fragmentShaderId, err := s.compileShader(gl.FRAGMENT_SHADER, s.fragmentSource)
	if err != nil {
		return errors.Wrapf(err, "fragment shader")
	}

	s.rendererId = gl.CreateProgram()

	gl.AttachShader(s.rendererId, vertexShaderId)
	gl.AttachShader(s.rendererId, fragmentShaderId)

	gl.LinkProgram(s.rendererId)

	var status int32
	gl.GetProgramiv(s.rendererId, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32

		gl.GetProgramiv(s.rendererId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.rendererId, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShaderId)
	gl.DeleteShader(fragmentShaderId)

	return nil
}

// compile shader of a given type from a source code
func (s *Shader) compileShader(shaderType uint32, shaderSource string) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	cSources, free := gl.Strs(shaderSource + "\x00")
	gl.ShaderSource(shader, 1, cSources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile: %v", log)
	}

	return shader, nil
}

func (s *Shader) SetUniform1i(name string, v0 int32) {
	gl.Uniform1i(s.getUniformLocation(name), v0)
}

func (s *Shader) SetUniform2f(name string, v0, v1 float32) {
	gl.Uniform2f(s.getUniformLocation(name), v0, v1)
}

func (s *Shader) SetUniform4f(name string, v0, v1, v2, v3 float32) {
	gl.Uniform4f(s.getUniformLocation(name), v0, v1, v2, v3)
}

func (s *Shader) SetUniformMat4f(name string, mat4 mgl32.Mat4) {
	gl.UniformMatrix4fv(s.getUniformLocation(name), 1, false, &mat4[0])
}

func (s *Shader) getUniformLocation(name string) int32 {
	var location int32
	var ok bool

	nameZeroed := name + "\x00"

	location, ok = s.uniformCache[nameZeroed]
	if !ok {
		location = gl.GetUniformLocation(s.rendererId, gl.Str(nameZeroed))
		s.uniformCache[nameZeroed] = location
	}

	return s.uniformCache[nameZeroed]
}
