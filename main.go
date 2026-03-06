package main

import (
	"lamp/ansi"
	"lamp/config"
	"lamp/events"
	"lamp/terminal"
	"lamp/window"
	"log"
	"os/user"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/gdamore/tcell/v2"
)

func main() {
	window.InitFont()

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	term, err := terminal.New(u)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Close()

	ansiState := ansi.NewState(config.Rows)
	go term.ReadLoop(ansiState)

	logicalW := float32(config.Cols * window.CharW / 2)
	logicalH := float32(config.Rows * window.CharH / 2)

	raster := window.NewRaster(term.Screen, &term.CursorX, &term.CursorY)
	raster.SetMinSize(fyne.NewSize(logicalW, logicalH))

	a := app.New()
	w := a.NewWindow("Lamp")
	w.SetContent(raster)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(logicalW, logicalH))

	go func() {
		ticker := time.NewTicker(time.Second / 30)
		for range ticker.C {
			fyne.Do(raster.Refresh)
		}
	}()

	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		if ev := window.FyneKeyToTcell(e); ev != nil {
			events.HandleEvent(term.Screen, ev, term.Write)
		}
	})
	w.Canvas().SetOnTypedRune(func(r rune) {
		ev := tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
		events.HandleEvent(term.Screen, ev, term.Write)
	})

	w.Show()
	go func() {
		time.Sleep(50 * time.Millisecond)
		fyne.Do(func() {
			w.Resize(fyne.NewSize(logicalW+1, logicalH+1))
			w.Resize(fyne.NewSize(logicalW, logicalH))
		})
	}()
	w.ShowAndRun()
}
