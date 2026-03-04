package params

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func ParseParams(params string) []int {
	params = strings.TrimSpace(params)
	if params == "" {
		return []int{}
	}
	parts := strings.Split(params, ";")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if n, err := strconv.Atoi(p); err == nil {
			result = append(result, n)
		}
	}
	return result
}

func GetArg(args []int, index int, def int) int {
	if index >= len(args) || args[index] == 0 {
		return def
	}
	return args[index]
}

func AnsiColor(n int) tcell.Color {
	switch n {
	case 0:
		return tcell.ColorBlack
	case 1:
		return tcell.ColorMaroon
	case 2:
		return tcell.ColorGreen
	case 3:
		return tcell.ColorOlive
	case 4:
		return tcell.ColorNavy
	case 5:
		return tcell.ColorPurple
	case 6:
		return tcell.ColorTeal
	case 7:
		return tcell.ColorSilver
	case 8:
		return tcell.ColorGray
	case 9:
		return tcell.ColorRed
	case 10:
		return tcell.ColorLime
	case 11:
		return tcell.ColorYellow
	case 12:
		return tcell.ColorBlue
	case 13:
		return tcell.ColorFuchsia
	case 14:
		return tcell.ColorAqua
	case 15:
		return tcell.ColorWhite
	default:
		return tcell.ColorDefault
	}
}
