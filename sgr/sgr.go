package sgr

import (
	"lamp/params"

	"github.com/gdamore/tcell/v2"
)

func HandleSGR(args []int, style *tcell.Style) {
	if len(args) == 0 {
		*style = tcell.StyleDefault
		return
	}

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == 0:
			*style = tcell.StyleDefault
		case a == 1:
			*style = (*style).Bold(true)
		case a == 2:
			*style = (*style).Dim(true)
		case a == 3:
			*style = (*style).Italic(true)
		case a == 4:
			*style = (*style).Underline(true)
		case a == 5:
			*style = (*style).Blink(true)
		case a == 7:
			*style = (*style).Reverse(true)
		case a == 22:
			*style = (*style).Bold(false).Dim(false)
		case a == 23:
			*style = (*style).Italic(false)
		case a == 24:
			*style = (*style).Underline(false)
		case a == 25:
			*style = (*style).Blink(false)
		case a == 27:
			*style = (*style).Reverse(false)
		case 30 <= a && a <= 37:
			*style = (*style).Foreground(params.AnsiColor(a - 30))
		case 40 <= a && a <= 47:
			*style = (*style).Background(params.AnsiColor(a - 40))
		case 90 <= a && a <= 97:
			*style = (*style).Foreground(params.AnsiColor(a - 90 + 8))
		case 100 <= a && a <= 107:
			*style = (*style).Background(params.AnsiColor(a - 100 + 8))
		case a == 39:
			*style = (*style).Foreground(tcell.ColorDefault)
		case a == 49:
			*style = (*style).Background(tcell.ColorDefault)
		case a == 38 && i+2 < len(args) && args[i+1] == 5:
			*style = (*style).Foreground(tcell.Color(args[i+2]))
			i += 2
		case a == 48 && i+2 < len(args) && args[i+1] == 5:
			*style = (*style).Background(tcell.Color(args[i+2]))
			i += 2
		case a == 38 && i+4 < len(args) && args[i+1] == 2:
			*style = (*style).Foreground(tcell.NewRGBColor(
				int32(args[i+2]), int32(args[i+3]), int32(args[i+4]),
			))
			i += 4
		case a == 48 && i+4 < len(args) && args[i+1] == 2:
			*style = (*style).Background(tcell.NewRGBColor(
				int32(args[i+2]), int32(args[i+3]), int32(args[i+4]),
			))
			i += 4
		}
	}
}
