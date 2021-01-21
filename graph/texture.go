package graph

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/pkg/errors"
	"image"
	"image/draw"
	"os"
)

type Texture struct {
	filePath string

	rendererId uint32

	width int
	height int
}

func NewTexture(filePath string) (*Texture, error) {
	ret := &Texture{
		filePath: filePath,
	}

	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "open texture file")
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, errors.Wrap(err, "decode image")
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, errors.New("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.GenTextures(1, &ret.rendererId)
	ret.Bind(0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return ret, nil
}

func (t *Texture) Bind(slot uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + slot)
	gl.BindTexture(gl.TEXTURE_2D, t.rendererId)
}

func (t Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t Texture) Destroy() {
	gl.DeleteTextures(1, &t.rendererId)
}
