package main

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertMessageToImage(t *testing.T) {
	msg := &DrawMessage{
		Params: DrawMessageParams{
			W: 2,
			H: 2,
			P: "01_0",
		},
	}
	img := convertMessageToImage(msg)

	assert.Equal(t, Black, img.At(0, 0))
	assert.Equal(t, White, img.At(1, 0))
	assert.Equal(t, Transparent, img.At(0, 1))
	assert.Equal(t, Black, img.At(1, 1))
}

func Test_convertImageToMessage(t *testing.T) {
	img := image.NewPaletted(image.Rect(0, 0, 2, 2), Palette)
	img.Set(0, 0, White)
	img.Set(1, 1, Black)
	msg := convertImageToMessage(img)

	assert.Equal(t, "init", msg.Method)
	assert.Equal(t, 2, msg.Params.W)
	assert.Equal(t, 2, msg.Params.H)
	assert.Equal(t, "1__0", msg.Params.P)
}
