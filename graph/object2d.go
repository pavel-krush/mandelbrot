package graph

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
	"unsafe"
)

type Object2D struct {
	vertices []mgl32.Vec4

	va *VertexArray
	ib *IndexBuffer
}

func NewObject2d() *Object2D {
	ret := &Object2D{}
	return ret
}

func (o *Object2D) AddVertex(vCoords mgl32.Vec2, texCoords mgl32.Vec2) {
	o.vertices = append(o.vertices, mgl32.Vec4{vCoords[0], vCoords[1], texCoords[0], texCoords[1]})
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

	// Create index buffer
	var indexBufferData []uint32

	var v0, v1, v2 uint32 = 0, 1, 2
	last := false
	for {
		// Create polygon
		indexBufferData = append(indexBufferData, v0)
		indexBufferData = append(indexBufferData, v1)
		indexBufferData = append(indexBufferData, v2)

		if last {
			break
		}

		// Advance vertex indices
		v0 = v1
		v1 = v2
		v2++
		if v2 == uint32(len(o.vertices) - 1) {
			v2 = 0
			last = true
		}
	}

	o.ib = NewIndexBuffer(indexBufferData)

	return nil
}

func (o *Object2D) Draw() {
	//shader.Bind()
	//texture.Bind(0)
	//shader.SetUniform1i("u_Texture", 0)
	//shader.SetUniformMat4f("u_MVP", proj)
	//renderer.Draw(va, ib, shader)
}
