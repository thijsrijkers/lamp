package ansi

import "github.com/gdamore/tcell/v2"

func scrollUp(screen tcell.Screen, style tcell.Style, width, scrollTop, scrollBottom int) {
	for y := scrollTop + 1; y <= scrollBottom; y++ {
		for x := 0; x < width; x++ {
			ch, comb, st, _ := screen.GetContent(x, y)
			screen.SetContent(x, y-1, ch, comb, st)
		}
	}
	for x := 0; x < width; x++ {
		screen.SetContent(x, scrollBottom, ' ', nil, style)
	}
}

func scrollDown(screen tcell.Screen, style tcell.Style, width, scrollTop, scrollBottom int) {
	for y := scrollBottom - 1; y >= scrollTop; y-- {
		for x := 0; x < width; x++ {
			ch, comb, st, _ := screen.GetContent(x, y)
			screen.SetContent(x, y+1, ch, comb, st)
		}
	}
	for x := 0; x < width; x++ {
		screen.SetContent(x, scrollTop, ' ', nil, style)
	}
}
