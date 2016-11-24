package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

var (
	DisplaySize = sdl.Point{480, 320}
)

var (
	dataDir    = flag.String("data", "data", "data directory")
	fullscreen = flag.Bool("fullscreen", false, "fullscreen")
	sfx        = flag.Bool("sfx", true, "sfx")
	chargeInf  = flag.Bool("chargeinf", false, "infinite charge")

	screen  *Display
	texture *sdl.Texture
	surface *sdl.Surface
)

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

func newDisplay(w, h int, wflag sdl.WindowFlags) (*Display, error) {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, wflag)
	if err != nil {
		return nil, err
	}
	return &Display{window, renderer}, nil
}
func main() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
	flag.Parse()
	initSDL()

	gameState := newGameState()
	title := newTitle(gameState)
	instructions := newInstructions()
	game := newGame(gameState)
	win := newWin(gameState)
	gameOver := newGameOver(gameState)

	state := "title"
	for state != "quit" {
		switch state {
		case "title":
			state = runState(title)
		case "instructions":
			state = runState(instructions)
		case "playgame":
			state = runState(game)
		case "win":
			state = runState(win)
		case "gameover":
			state = runState(gameOver)
		default:
			panic(fmt.Sprintf("unreachable state: %q", state))
		}
	}
}

func initSDL() {
	log.SetPrefix("sdl: ")
	err := sdl.Init(sdl.INIT_EVERYTHING &^ sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal(err)
	}

	err = sdl.InitSubSystem(sdl.INIT_AUDIO)
	if err != nil {
		log.Println(err)
	}

	err = sdlmixer.OpenAudio(44100, sdl.AUDIO_S16, 2, 8192)
	if err != nil {
		log.Println(err)
	}

	err = sdlttf.Init()
	if err != nil {
		log.Fatal(err)
	}

	sdlmixer.AllocateChannels(128)

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "best")

	w, h := int(DisplaySize.X), int(DisplaySize.Y)
	wflag := sdl.WINDOW_RESIZABLE
	screen, err = newDisplay(w, h, wflag)
	if err != nil {
		log.Fatal(err)
	}

	texture, err = screen.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, w, h)
	if err != nil {
		log.Fatal(err)
	}

	surface, err = sdl.CreateRGBSurfaceWithFormat(sdl.SWSURFACE, w, h, 32, sdl.PIXELFORMAT_ABGR8888)
	if err != nil {
		log.Fatal(err)
	}

	screen.SetTitle("Godspeed You! Mutated Maze")
	screen.SetLogicalSize(w, h)
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()
	screen.Present()

	sdl.ShowCursor(0)
}

type baseState struct {
	gs *GameState
}

func (*baseState) init()          {}
func (*baseState) event() string  { return "" }
func (*baseState) update() string { return "" }
func (*baseState) draw()          {}

type State interface {
	init()
	event() string
	update() string
	draw()
}

func runState(s State) string {
	s.init()
	for {
		if newState := s.event(); newState != "" {
			return newState
		}
		if newState := s.update(); newState != "" {
			return newState
		}
		s.draw()
		sdl.Delay(1000 / 30)
	}
}
