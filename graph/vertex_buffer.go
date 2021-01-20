package graph

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexBuffer struct {
	rendererId uint32
}

func NewVertexBuffer(data unsafe.Pointer, size int) *VertexBuffer {
	ret := VertexBuffer{}

	gl.GenBuffers(1, &ret.rendererId)
	gl.BindBuffer(gl.ARRAY_BUFFER, ret.rendererId)
	gl.BufferData(gl.ARRAY_BUFFER, size, data, gl.STATIC_DRAW)

	return &ret
}

func (vb *VertexBuffer) Destroy() {
	gl.DeleteBuffers(1, &vb.rendererId)
}

func (vb *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.rendererId)
}

func (vb *VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
