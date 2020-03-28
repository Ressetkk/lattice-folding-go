package main

import (
	"github.com/fogleman/gg"
)

const width, height = 1000, 1000
const proteinRadius = 16.0
const lineWidth = 4.0

func GenerateImage(nodes *[]*Node, out string) error {
	canvas := gg.NewContext(width, height)
	canvas.SetLineWidth(lineWidth)

	var px, py float64
	for i, node := range *nodes {
		x := float64(width)/2.0 + float64(node.x)*proteinRadius*4
		y := float64(height)/2.0 + float64(node.y)*proteinRadius*4

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

	return canvas.SavePNG(out)
}
