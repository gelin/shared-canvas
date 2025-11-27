package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
)

type ImageHolder struct {
	image *image.Paletted
	Draw  chan *DrawMessage
}

type DrawMessage struct {
	Method string            `json:"method"`
	Params DrawMessageParams `json:"params,omitempty"`
}

type DrawMessageParams struct {
	X int    `json:"x,optional"`
	Y int    `json:"y,optional"`
	W int    `json:"w"`
	H int    `json:"h"`
	P string `json:"p"`
}

var Transparent = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
var Black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
var White = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
var Palette = color.Palette{Transparent, Black, White}
var TransparentIndex = Palette.Index(Transparent)
var BlackIndex = Palette.Index(Black)
var WhiteIndex = Palette.Index(White)

func NewImageHolder() *ImageHolder {
	// TODO: make width and height configurable
	img := image.NewPaletted(image.Rect(0, 0, 384, 384), Palette)
	draw.Draw(img, img.Bounds(), image.NewUniform(White), image.Point{}, draw.Src)
	holder := &ImageHolder{
		image: img,
		Draw:  make(chan *DrawMessage),
	}
	go holder.run()
	return holder
}

func (h *ImageHolder) run() {
	for {
		msg := <-h.Draw
		h.draw(msg)
	}
}

func (h *ImageHolder) draw(msg *DrawMessage) {
	drawImage := convertMessageToImage(msg)
	drawRect := image.Rect(msg.Params.X, msg.Params.Y, msg.Params.X+msg.Params.W, msg.Params.Y+msg.Params.H)
	draw.Draw(h.image, drawRect, drawImage, image.Point{}, draw.Over)
}

func convertMessageToImage(msg *DrawMessage) image.Image {
	img := image.NewPaletted(image.Rect(0, 0, msg.Params.W, msg.Params.H), Palette)
	for i, c := range msg.Params.P {
		switch c {
		case '_':
			img.Pix[i] = uint8(TransparentIndex)
			break
		case '0':
			img.Pix[i] = uint8(BlackIndex)
			break
		case '1':
			img.Pix[i] = uint8(WhiteIndex)
			break
		}
	}
	return img
}

func (h *ImageHolder) WriteImagePNG(w io.Writer) {
	err := png.Encode(w, h.image)
	if err != nil {
		log.Printf("error encoding image: %v", err)
	}
}
