package window

import (
	"image"
	"image/color"
	"lamp/config"

	"fyne.io/fyne/v2/canvas"
	"github.com/gdamore/tcell/v2"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func NewRaster(screen tcell.SimulationScreen, cursorX, cursorY *int) *canvas.Raster {
	return canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))

		cells, sw, sh := screen.GetContents()
		ascent := Face.Metrics().Ascent.Ceil()

		cw := w / config.Cols
		ch := h / config.Rows
		if cw == 0 {
			cw = 1
		}
		if ch == 0 {
			ch = 1
		}

		for row := 0; row < sh && row < config.Rows; row++ {
			for col := 0; col < sw && col < config.Cols; col++ {
				cell := cells[row*sw+col]
				r := rune(' ')
				if len(cell.Runes) > 0 && cell.Runes[0] != 0 {
					r = cell.Runes[0]
				}
				fg, bg, _ := cell.Style.Decompose()
				x, y := col*cw, row*ch

				FillRect(img, x, y, cw, ch, TcellColorToRGBA(bg, false))

				if r != ' ' {
					(&xfont.Drawer{
						Dst:  img,
						Src:  image.NewUniform(TcellColorToRGBA(fg, true)),
						Face: Face,
						Dot:  fixed.P(x, y+ascent),
					}).DrawString(string(r))
				}
			}
		}

		cx, cy := *cursorX, *cursorY
		if cx >= 0 && cx < config.Cols && cy >= 0 && cy < config.Rows {
			FillRect(img, cx*cw, cy*ch, cw, ch,
				color.RGBA{R: 255, G: 255, B: 255, A: 80})
		}

		return img
	})
}
