package graph

import (
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type IndexBuffer struct {
	rendererId uint32
	count int32
}

func NewIndexBuffer(data []uint32) *IndexBuffer {
	ret := IndexBuffer{
		count: int32(len(data)),
	}

	gl.GenBuffers(1, &ret.rendererId)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ret.rendererId)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, int(ret.count) * int(unsafe.Sizeof(uint32(0))), gl.Ptr(data), gl.STATIC_DRAW)

	return &ret
}

func (ib *IndexBuffer) Destroy() {
	gl.DeleteBuffers(1, &ib.rendererId)
}

func (ib *IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib.rendererId)
}

func (ib *IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (ib *IndexBuffer) GetCount() int32 {
	return ib.count
}
