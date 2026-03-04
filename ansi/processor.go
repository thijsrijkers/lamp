package ansi

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

func ProcessOutput(screen tcell.Screen, data []byte, cursorX, cursorY *int, state *State) {
	width, height := screen.Size()

	if state.ScrollBottom >= height {
		state.ScrollBottom = height - 1
	}

	i := 0
	for i < len(data) {
		ch := data[i]

		switch ch {
		case '\r':
			*cursorX = 0
		case '\n':
			*cursorX = 0
			*cursorY++
		case '\t':
			*cursorX += 4
		case 0x08, 0x7f:
			if *cursorX > 0 {
				*cursorX--
				screen.SetContent(*cursorX, *cursorY, ' ', nil, state.Style)
			}
		case 0x1b:
			if i+1 < len(data) {
				switch data[i+1] {
				case '[':
					i += 2
					start := i
					for i < len(data) && (data[i] < '@' || data[i] > '~') {
						i++
					}
					if i < len(data) {
						cmd := data[i]
						params := string(data[start:i])
						if cmd == 'h' && params == "?1049" {
							screen.Clear()
							*cursorX = 0
							*cursorY = 0
							state.ScrollTop = 0
							state.ScrollBottom = height - 1
						}
						if cmd == 'l' && params == "?1049" {
							screen.Clear()
							*cursorX = 0
							*cursorY = 0
							state.ScrollTop = 0
							state.ScrollBottom = height - 1
						}
						handleCSI(screen, cmd, params, cursorX, cursorY, &state.Style, width, height, state)
					}
				case ']':
					i += 2
					for i < len(data) && data[i] != 0x07 {
						if data[i] == 0x1b && i+1 < len(data) && data[i+1] == '\\' {
							i++
							break
						}
						i++
					}
				case '(', ')', '*', '+', '-', '.', '/':
					i += 2
					if i < len(data) && data[i] == 'B' {
						i++
					}
					continue
				default:
					i++
				}
			}
		default:
			if ch >= 32 {
				r, size := utf8.DecodeRune(data[i:])
				x, y := *cursorX, *cursorY
				if x >= width {
					x = width - 1
				}
				if y >= height {
					y = height - 1
				}
				screen.SetContent(x, y, r, nil, state.Style)
				*cursorX++
				i += size
				continue
			}
		}

		if *cursorY > state.ScrollBottom {
			scrollUp(screen, state.Style, width, state.ScrollTop, state.ScrollBottom)
			*cursorY = state.ScrollBottom
		}

		i++
	}

	screen.Show()
}
