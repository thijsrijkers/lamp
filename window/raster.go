package window

import (
	"image"
	"image/color"
	"lamp/config"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/gdamore/tcell/v2"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Raster struct {
	widget.BaseWidget
	screen      tcell.SimulationScreen
	cursorX     *int
	cursorY     *int
	raster      *canvas.Raster
	selectStart fyne.Position
	OnSelect    func(string)
}

func NewRaster(screen tcell.SimulationScreen, cursorX, cursorY *int) *Raster {
	r := &Raster{screen: screen, cursorX: cursorX, cursorY: cursorY}
	r.raster = canvas.NewRaster(r.draw)
	r.ExtendBaseWidget(r)
	return r
}

func (r *Raster) draw(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	cells, sw, sh := r.screen.GetContents()
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
			ru := rune(' ')
			if len(cell.Runes) > 0 && cell.Runes[0] != 0 {
				ru = cell.Runes[0]
			}
			fg, bg, _ := cell.Style.Decompose()
			x, y := col*cw, row*ch

			FillRect(img, x, y, cw, ch, TcellColorToRGBA(bg, false))

			if ru != ' ' {
				(&xfont.Drawer{
					Dst:  img,
					Src:  image.NewUniform(TcellColorToRGBA(fg, true)),
					Face: Face,
					Dot:  fixed.P(x, y+ascent),
				}).DrawString(string(ru))
			}
		}
	}

	cx, cy := *r.cursorX, *r.cursorY
	if cx >= 0 && cx < config.Cols && cy >= 0 && cy < config.Rows {
		FillRect(img, cx*cw, cy*ch, cw, ch,
			color.RGBA{R: 255, G: 255, B: 255, A: 80})
	}

	return img
}

func (r *Raster) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.raster)
}

func (r *Raster) MouseDown(e *desktop.MouseEvent) {
	r.selectStart = e.Position
}

func (r *Raster) MouseUp(e *desktop.MouseEvent) {
	col1 := int(r.selectStart.X / float32(CharW))
	row1 := int(r.selectStart.Y / float32(CharH))
	col2 := int(e.Position.X / float32(CharW))
	row2 := int(e.Position.Y / float32(CharH))

	var buf strings.Builder
	for row := row1; row <= row2; row++ {
		for col := col1; col <= col2; col++ {
			ru, _, _, _ := r.screen.GetContent(col, row)
			buf.WriteRune(ru)
		}
	}

	if r.OnSelect != nil {
		r.OnSelect(buf.String())
	}
}

func (r *Raster) SetMinSize(size fyne.Size) {
	r.raster.SetMinSize(size)
}
