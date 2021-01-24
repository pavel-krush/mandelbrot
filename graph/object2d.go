package graph

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"unsafe"
)

type Object2D struct {
	shader *Shader
	texture *Texture
	va *VertexArray
	ib *IndexBuffer

	vertices []mgl32.Vec4
	ibData []uint32
}

func NewObject2d(shader *Shader, texture *Texture) (*Object2D, error) {
	if shader == nil {
		return nil, errors.Errorf("shader required")
	}

	if texture == nil {
		return nil, errors.Errorf("texture required")
	}

	ret := &Object2D{
		shader: shader,
		texture: texture,
	}

	return ret, nil
}

func (o *Object2D) AddVertex(vCoords mgl32.Vec2, texCoords mgl32.Vec2) {
	o.vertices = append(o.vertices, mgl32.Vec4{vCoords[0], vCoords[1], texCoords[0], texCoords[1]})
}

func (o *Object2D) AddIndexBufferData(data ...uint32) {
	o.ibData = append(o.ibData, data...)
}

func (o *Object2D) Compile() error {
	if len(o.vertices) < 3 {
		return errors.New("At least 3 vertices required")
	}

	// Setup vertex buffer
	bufferSize := len(o.vertices) * // count of vertices
		4 * // count of components per vertex
		int(unsafe.Sizeof(float32(0))) // size of one component - float32

	vb := NewVertexBuffer(gl.Ptr(o.vertices), bufferSize)
	layout := NewVertexBufferLayout()
	layout.PushFloat(2)
	layout.PushFloat(2)

	// Setup vertex array
	o.va = NewVertexArray()
	o.va.AddBuffer(vb, layout)

	// Setup index buffer
	o.ib = NewIndexBuffer(o.ibData)

	return nil
}

func (o *Object2D) Destroy() {
	o.va.Destroy()
	o.ib.Destroy()
}

// Drawable interface
func (o *Object2D) GetVertexArray() *VertexArray {
	return o.va
}

func (o *Object2D) GetIndexBuffer() *IndexBuffer {
	return o.ib
}

func (o *Object2D) GetShader() *Shader {
	return o.shader
}

func (o *Object2D) GetTexture() *Texture {
	return o.texture
}
