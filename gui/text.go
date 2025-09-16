package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//	func TextDrawV1(dst *ebiten.Image, text string, face font.Face, x, y int, clr color.Color) {
//		"github.com/hajimehoshi/ebiten/v2/text".Draw(dst, text, face, x, y, clr)
//	}

func TextDraw(dst *ebiten.Image, txt string, xface *text.GoXFace, x, y int, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, txt, xface, op)
}
