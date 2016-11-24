package main

import (
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Instructions struct {
	baseState
	bg *Image
}

func newInstructions() *Instructions {
	return &Instructions{
		bg: loadImage("instructions.png"),
	}
}

func (p *Instructions) event() string {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}

		switch ev := ev.(type) {
		case sdl.QuitEvent:
			return "quit"
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				return "title"
			case sdl.K_SPACE:
				return "playgame"
			}
		}
	}
	return ""
}

func (p *Instructions) draw() {
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()
	p.bg.Blit(0, 0)
	screen.Present()
}
