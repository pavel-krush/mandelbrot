package graph

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {

}

type Drawable interface {
	GetVertexArray() *VertexArray
	GetIndexBuffer() *IndexBuffer
	GetShader() *Shader
	GetTexture() *Texture
}

func NewRenderer() *Renderer {
	return nil
}

func (r *Renderer) Draw(object Drawable, x int, y int) {
	shader := object.GetShader()
	shader.Bind()
	shader.SetUniform1i("u_Texture", 0)
	shader.SetUniform2f("coords", float32(x), float32(y))

	object.GetVertexArray().Bind()
	ib := object.GetIndexBuffer()
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
