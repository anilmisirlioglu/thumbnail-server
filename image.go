package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

type Grid struct {
	image *image.RGBA
}

func (g *Grid) Set(x, y int, c color.Color) {
	g.image.Set(x, y, c)
}

func (g *Grid) ColorModel() color.Model {
	return g.image.ColorModel()
}

func (g *Grid) Bounds() image.Rectangle {
	return g.image.Bounds()
}

func (g *Grid) At(x, y int) color.Color {
	return g.image.At(x, y)
}

func DrawBackground(opts ScreenshotOption, img image.Image) image.Image {
	width, height := opts.background.width, opts.background.height

	background := Grid{
		image.NewRGBA(image.Rectangle{
			Min: image.Point{},
			Max: image.Point{
				X: width,
				Y: height,
			},
		}),
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			background.Set(x, y, opts.background.color)
		}
	}

	background.draw(img, image.Point{
		X: opts.background.width/2 - int(opts.width)/2,
		Y: opts.background.height/2 - int(opts.height)/2,
	})

	return background.image
}

func (g *Grid) draw(img image.Image, point image.Point) {
	bounds := g.Bounds()
	rect := image.Rectangle{
		Min: point,
		Max: image.Point{
			X: point.X + bounds.Max.X - bounds.Min.X,
			Y: point.Y + bounds.Max.Y - bounds.Min.Y,
		},
	}

	draw.Draw(g, rect, img, image.Point{}, draw.Src)
}

func ServeFrames(byteArray []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(byteArray))
	if err != nil {
		return nil, err
	}

	out, _ := ioutil.TempFile("", "image.png")
	err = png.Encode(out, img)
	if err != nil {
		return nil, err
	}

	defer os.Remove(out.Name())
	return img, nil
}
