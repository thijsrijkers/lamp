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
	selectEnd   fyne.Position
	selecting   bool
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

	if r.selecting || r.selectEnd != r.selectStart {
		size := r.raster.Size()
		cellW := size.Width / float32(config.Cols)
		cellH := size.Height / float32(config.Rows)

		col1 := int(r.selectStart.X / cellW)
		row1 := int(r.selectStart.Y / cellH)
		col2 := int(r.selectEnd.X / cellW)
		row2 := int(r.selectEnd.Y / cellH)

		col1 = max(0, min(col1, config.Cols-1))
		col2 = max(0, min(col2, config.Cols-1))
		row1 = max(0, min(row1, config.Rows-1))
		row2 = max(0, min(row2, config.Rows-1))

		for row := row1; row <= row2; row++ {
			for col := col1; col <= col2; col++ {
				FillRect(img, col*cw, row*ch, cw, ch,
					color.RGBA{R: 100, G: 150, B: 255, A: 120})
			}
		}
	}

	return img
}

func (r *Raster) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.raster)
}

func (r *Raster) MouseMoved(e *desktop.MouseEvent) {
	if r.selecting {
		r.selectEnd = e.Position
		r.raster.Refresh()
	}
}

func (r *Raster) MouseDown(e *desktop.MouseEvent) {
	r.selectStart = e.Position
	r.selectEnd = e.Position
	r.selecting = true
	r.raster.Refresh()
}

func (r *Raster) MouseUp(e *desktop.MouseEvent) {
	r.selectEnd = e.Position
	r.selecting = false

	size := r.raster.Size()
	cellW := size.Width / float32(config.Cols)
	cellH := size.Height / float32(config.Rows)

	col1 := int(r.selectStart.X / cellW)
	row1 := int(r.selectStart.Y / cellH)
	col2 := int(e.Position.X / cellW)
	row2 := int(e.Position.Y / cellH)

	col1 = max(0, min(col1, config.Cols-1))
	col2 = max(0, min(col2, config.Cols-1))
	row1 = max(0, min(row1, config.Rows-1))
	row2 = max(0, min(row2, config.Rows-1))

	var buf strings.Builder
	for row := row1; row <= row2; row++ {
		for col := col1; col <= col2; col++ {
			ru, _, _, _ := r.screen.GetContent(col, row)
			buf.WriteRune(ru)
		}
		if row < row2 {
			buf.WriteRune('\n')
		}
	}

	if r.OnSelect != nil && buf.Len() > 0 {
		r.OnSelect(buf.String())
	}

	// clear selection after copy
	r.selectStart = fyne.Position{}
	r.selectEnd = fyne.Position{}
	r.raster.Refresh()
}

func (r *Raster) SetMinSize(size fyne.Size) {
	r.raster.SetMinSize(size)
}
