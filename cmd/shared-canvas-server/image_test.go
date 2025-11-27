package main

import (
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
