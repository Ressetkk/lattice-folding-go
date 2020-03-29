package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"math"
	"net/http"
)

const width, height = 1000, 1000
const proteinRadius = 16.0
const lineWidth = 4.0

func GenerateImage(nodes *[]*Node, min int, out string) error {
	canvas := gg.NewContext(width, height)
	canvas.SetLineWidth(lineWidth)

	var txmin, txmax, tymin, tymax float64
	for _, node := range *nodes {
		if float64(node.x) < txmin {
			txmin = float64(node.x)
		}
		if float64(node.x) > txmax {
			txmax = float64(node.x)
		}
		if float64(node.y) < tymin {
			tymin = float64(node.y)
		}
		if float64(node.y) > tymax {
			tymax = float64(node.y)
		}
	}

	offsetX, offsetY := (txmax-txmin)*proteinRadius, (tymax-tymin)*proteinRadius
	if math.Abs(txmin) < txmax {
		offsetX = -offsetX
	}
	if math.Abs(tymin) < tymax {
		offsetY = -offsetY
	}

	var px, py float64
	for i, node := range *nodes {
		x := (float64(width)/2.0 + float64(node.x)*proteinRadius*4) + offsetX
		y := (float64(height)/2.0 + float64(node.y)*proteinRadius*4) + offsetY

		if i > 0 {
			canvas.SetRGB(0, 0, 0)
			canvas.DrawLine(x, y, px, py)
			canvas.Stroke()
		}
		px = x
		py = y

		if node.h {
			canvas.SetRGB(0, 0, 0)
		} else {
			canvas.SetRGB(1, 1, 1)
		}
		canvas.DrawCircle(x, y, proteinRadius)
		canvas.Fill()
	}

	fb, err := downloadFont("https://github.com/google/fonts/raw/master/apache/roboto/Roboto-Regular.ttf")
	if err != nil {
		return err
	}
	f, _ := truetype.Parse(fb)
	canvas.SetFontFace(truetype.NewFace(f, &truetype.Options{Size: 32}))
	canvas.DrawString(fmt.Sprintf("Energy: %v", min), 32, 32)
	return canvas.SavePNG(out)
}

func downloadFont(url string) ([]byte, error) {
	fr, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(fr.Body)
}
