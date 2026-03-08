package window

import (
	"fmt"
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
	screen       tcell.SimulationScreen
	cursorX      *int
	cursorY      *int
	raster       *canvas.Raster
	selectStart  fyne.Position
	selectEnd    fyne.Position
	selecting    bool
	OnSelect     func(string)
	Write        func([]byte)
	MouseEnabled *bool
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
		scaleX := float32(w) / r.raster.Size().Width
		scaleY := float32(h) / r.raster.Size().Height

		col1 := int(r.selectStart.X*scaleX) / cw
		row1 := int(r.selectStart.Y*scaleY) / ch
		col2 := int(r.selectEnd.X*scaleX) / cw
		row2 := int(r.selectEnd.Y*scaleY) / ch

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

func (r *Raster) mouseCol(x float32) int {
	size := r.raster.Size()
	cellW := size.Width / float32(config.Cols)
	return int(x/cellW) + 1
}

func (r *Raster) mouseRow(y float32) int {
	size := r.raster.Size()
	cellH := size.Height / float32(config.Rows)
	return int(y/cellH) + 1
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

	dx := e.Position.X - r.selectStart.X
	dy := e.Position.Y - r.selectStart.Y
	isClick := dx*dx+dy*dy < 25

	if isClick && r.Write != nil && r.MouseEnabled != nil && *r.MouseEnabled {
		col := r.mouseCol(e.Position.X)
		row := r.mouseRow(e.Position.Y)
		btn := 0
		if e.Button == desktop.MouseButtonSecondary {
			btn = 2
		}
		r.Write([]byte(fmt.Sprintf("\x1b[<%d;%d;%dM", btn, col, row)))
		r.Write([]byte(fmt.Sprintf("\x1b[<%d;%d;%dm", btn, col, row)))
		r.selectStart = fyne.Position{}
		r.selectEnd = fyne.Position{}
		r.raster.Refresh()
		return
	}

	// it was a drag — treat as selection
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

	r.selectStart = fyne.Position{}
	r.selectEnd = fyne.Position{}
	r.raster.Refresh()
}

func (r *Raster) Dragged(e *fyne.DragEvent) {
	r.selectEnd = e.Position
	r.raster.Refresh()
}

func (r *Raster) DragEnd() {}

func (r *Raster) SetMinSize(size fyne.Size) {
	r.raster.SetMinSize(size)
}
