package main

import (
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Title struct {
	baseState
	hdrFont  *sdlttf.Font
	instFont *sdlttf.Font
}

func newTitle(gs *GameState) *Title {
	t := &Title{}
	t.gs = gs
	t.hdrFont = loadFont("arial.ttf", 30)
	t.instFont = loadFont("arial.ttf", 15)
	return t
}

func (t *Title) event() string {
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
				return "quit"
			case sdl.K_SPACE:
				t.gs.Reset()
				return "instructions"
			}
		}
	}
	return ""
}

func (t *Title) draw() {
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()

	printCenter(t.hdrFont, 50, BrightMagenta, "Godspeed You! Mutated Maze")
	printCenter(t.instFont, 240, BrightMagenta, "Copyright 2011 Nick Sonneveld.")
	printCenter(t.instFont, 270, BrightMagenta, "Press <SPACE> to start.")

	screen.Present()
}
