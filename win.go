package main

import (
	"fmt"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Win struct {
	baseState
	hdrFont  *sdlttf.Font
	textFont *sdlttf.Font
	textList []Text
}

func newWin(gs *GameState) *Win {
	w := &Win{}
	w.gs = gs
	w.hdrFont = loadFont("arial.ttf", 40)
	w.textFont = loadFont("arial.ttf", 15)
	return w
}

func (w *Win) init() {
	var textList []Text

	// time left
	textList = append(textList, Text{
		fmt.Sprintf("Time left: %d seconds", w.gs.timeLeft),
		BrightMagenta,
	})

	scoreTime := w.gs.timeLeft * 30
	textList = append(textList, Text{
		fmt.Sprintf("+%d", scoreTime),
		BrightGreen,
	})
	textList = append(textList, Text{"", BrightMagenta})

	// people saved
	textList = append(textList, Text{
		fmt.Sprintf("Mobs saved: %d", w.gs.mobsSaved),
		BrightMagenta,
	})
	scoreMobs := w.gs.mobsSaved * 500
	textList = append(textList, Text{
		fmt.Sprintf("+%d", scoreMobs),
		BrightGreen,
	})

	scoreBonus := 0
	if w.gs.mobsSaved > 0 && w.gs.mobsSaved == w.gs.mobsAvail {
		scoreBonus = 2000
		textList = append(textList, Text{"You saved everyone!", BrightGreen})
		textList = append(textList, Text{fmt.Sprint("+", scoreBonus), BrightGreen})
		textList = append(textList, Text{"You are a super player!", BrightGreen})
	}
	textList = append(textList, Text{"", BrightMagenta})

	// score
	w.gs.score += scoreTime + scoreMobs + scoreBonus
	textList = append(textList, Text{fmt.Sprint("Score: ", w.gs.score), BrightMagenta})
	textList = append(textList, Text{"", BrightMagenta})

	w.textList = textList

}

func (w *Win) event() string {
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

func (w *Win) draw() {
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()

	printCenter(w.hdrFont, 20, BrightMagenta, "You ESCAPED!")

	y := 100
	for _, p := range w.textList {
		printCenter(w.textFont, y, p.color, p.text)
		y += w.textFont.Height()
	}

	screen.Present()
}
