package graph

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
	rendererId uint32
}

func NewVertexArray() *VertexArray {
	ret := VertexArray{}

	gl.GenVertexArrays(1, &ret.rendererId)
	gl.BindVertexArray(ret.rendererId)

	return &VertexArray{}
}

func (va *VertexArray) AddBuffer(vb *VertexBuffer, layout *VertexBufferLayout) {
	va.Bind()
	vb.Bind()

	elements := layout.GetElements()

	stride := layout.GetStride()

	var offset int
	var i uint32
	for i = 0; i < uint32(len(elements)); i ++ {
		element := elements[i]
		gl.EnableVertexAttribArray(i)
		gl.VertexAttribPointer(i, element.count, element.typ, element.normalized, stride, gl.PtrOffset(offset))
		offset += int(element.count) * int(element.typSize)
	}
}

func (va *VertexArray) Bind() {
	gl.BindVertexArray(va.rendererId)
}

func (va *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

func (va *VertexArray) Destroy() {
	gl.DeleteVertexArrays(1, &va.rendererId)
}
