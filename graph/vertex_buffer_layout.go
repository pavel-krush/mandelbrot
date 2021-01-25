package graph

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"unsafe"
)

type VertexBufferLayoutElement struct {
	typ        uint32 // type
	typSize    uint32 // type size in bytes
	count      int32  // elements count
	normalized bool
}

type VertexBufferLayout struct {
	elements []VertexBufferLayoutElement
}

func NewVertexBufferLayout() *VertexBufferLayout {
	return &VertexBufferLayout{}
}

func (vbl *VertexBufferLayout) PushFloat(count int32) {
	vbl.elements = append(vbl.elements, VertexBufferLayoutElement{
		count:      count,
		typ:        gl.FLOAT,
		typSize:    uint32(unsafe.Sizeof(float32(0))),
		normalized: false,
	})
}

func (vbl *VertexBufferLayout) GetElements() []VertexBufferLayoutElement {
	return vbl.elements
}

func (vbl *VertexBufferLayout) GetStride() int32 {
	var stride int32

	for _, element := range vbl.elements {
		stride += int32(element.typSize) * element.count
	}

	return stride
}
