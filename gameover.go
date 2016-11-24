package main

import (
	"fmt"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type GameOver struct {
	baseState
	titleFont *sdlttf.Font
	textFont  *sdlttf.Font
}

func newGameOver(gs *GameState) *GameOver {
	g := &GameOver{}
	g.gs = gs
	g.titleFont = loadFont("arial.ttf", 30)
	g.textFont = loadFont("arial.ttf", 15)
	return g
}

func (g *GameOver) event() string {
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
			case sdl.K_ESCAPE, sdl.K_SPACE:
				return "title"
			}
		}
	}
	return ""
}

func (g *GameOver) draw() {
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()

	printCenter(g.titleFont, 50, BrightMagenta, "Game over man, game over!")

	textList := []string{
		fmt.Sprint("Final score: ", g.gs.score),
		fmt.Sprint("Mobs saved: ", g.gs.mobsSavedTotal),
		"",
		"Press <SPACE> to continue.",
	}

	y := 200
	for _, text := range textList {
		printCenter(g.textFont, y, BrightMagenta, text)
		y += g.textFont.Height()
	}

	screen.Present()
}
