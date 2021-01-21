package graph

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {

}

func NewRenderer() *Renderer {
	return nil
}

func (r *Renderer) Draw(va *VertexArray, ib *IndexBuffer, shader *Shader) {
	shader.Bind()
	va.Bind()
	ib.Bind()

	gl.DrawElements(gl.TRIANGLES, ib.GetCount(), gl.UNSIGNED_INT, nil)
	err := gl.GetError()
	if err != gl.NO_ERROR {
		fmt.Printf("gl error: %d\n", err)
	}
}

func (r *Renderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
