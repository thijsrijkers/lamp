package terminal

import (
	"lamp/ansi"
	"lamp/config"
	"os"
	"os/exec"
	"os/user"

	"github.com/creack/pty"
	"github.com/gdamore/tcell/v2"
)

type Terminal struct {
	Screen  tcell.SimulationScreen
	CursorX int
	CursorY int
	ptmx    *os.File
}

func New(u *user.User) (*Terminal, error) {
	cmd := exec.Command("bash", "--login")
	cmd.Dir = u.HomeDir
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
		"PATH=/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin",
		"HOME="+u.HomeDir,
		"USER="+u.Username,
		"LOGNAME="+u.Username,
		"SHELL=/bin/bash",
	)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	screen := tcell.NewSimulationScreen("UTF-8")
	screen.Init()
	screen.SetSize(config.Cols, config.Rows)
	pty.Setsize(ptmx, &pty.Winsize{
		Cols: uint16(config.Cols),
		Rows: uint16(config.Rows),
	})

	return &Terminal{Screen: screen, ptmx: ptmx}, nil
}

func (t *Terminal) ReadLoop(state *ansi.State) {
	buf := make([]byte, 4096)
	for {
		n, err := t.ptmx.Read(buf)
		if err != nil {
			return
		}
		data := buf[:n]
		if len(state.Leftover) > 0 {
			data = append(state.Leftover, data...)
			state.Leftover = nil
		}
		ansi.ProcessOutput(t.Screen, data, &t.CursorX, &t.CursorY, state)
	}
}

func (t *Terminal) Write(b []byte) { t.ptmx.Write(b) }
func (t *Terminal) Close()         { t.ptmx.Close() }
