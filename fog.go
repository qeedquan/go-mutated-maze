package main

import (
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

var fogFiles = [][2]string{
	{"fog01.png", "fog01_inverse.png"},
	{"fog02.png", "fog02_inverse.png"},
	{"fog03.png", "fog03_inverse.png"},
	{"fog04.png", "fog04_inverse.png"},
	{"fog05.png", "fog05_inverse.png"},
	{"fog06.png", "fog06_inverse.png"},
}

var mazeColors []sdl.Color

func init() {
	mazeColors = append(mazeColors, SpriteColors...)
}

type Fog struct {
	Blitter
	gs          *GameState
	mutateSound *sdlmixer.Chunk

	mutateCount int
	flashCount  int

	pal     int
	vx, vy  int
	passed  bool
	mutated bool
	fog     *Image
	inverse *Image
}

func newFog(gs *GameState) *Fog {
	f := &Fog{}
	f.gs = gs

	for i := len(mazeColors) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		mazeColors[i], mazeColors[j] = mazeColors[j], mazeColors[i]
	}

	fogChoice := fogFiles[rand.Intn(len(fogFiles))]
	f.fog = loadFog(fogChoice[0], gs.maze)
	f.inverse = loadFog(fogChoice[1], gs.maze)
	f.image = f.fog

	f.mutateSound = loadSound("mutate.wav")
	f.x = gs.maze.Px()
	f.vx = -1
	f.mutateCount = f.x/2 + rand.Intn(f.x*4/3-f.x/2)
	return f
}

func loadFog(name string, m *Maze) *Image {
	return loadImage(name, scaleImage(m.Px(), m.Py()))
}

func (f *Fog) mutate() {
	playSound(f.mutateSound)

	newMazeColor := mazeColors[f.pal]
	f.pal = (f.pal + 1) % len(mazeColors)

	var selectedNodes []*MazeNode
	for _, n := range f.gs.maze.nodes {
		cx := n.x*16 + 8
		cy := n.y*16 + 8
		if f.posCoveredWithFog(cx, cy) {
			n.color = newMazeColor
			selectedNodes = append(selectedNodes, n)
		}
	}

	var mutateMobs []*Mob
	for _, m := range f.gs.AliveMobs() {
		cx := m.x + 8
		cy := m.y + 8
		if f.posCoveredWithFog(cx, cy) {
			mutateMobs = append(mutateMobs, m)
		}
	}

	for i := len(selectedNodes) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		selectedNodes[i], selectedNodes[j] = selectedNodes[j], selectedNodes[i]
	}

	i := 0
	for _, m := range mutateMobs {
		n := selectedNodes[i]
		i = (i + 1) % len(selectedNodes)
		x, y := n.Pxy()
		m.FogMutate(x, y)
	}

	f.gs.maze.RegenSelected(selectedNodes)
}

func (f *Fog) posCoveredWithFog(x, y int) bool {
	fx := x - f.x
	fy := y - f.y

	if fx < 0 || fx >= f.gs.maze.Px() {
		return false
	}
	if fy < 0 || fy >= f.gs.maze.Py() {
		return false
	}

	return f.fog.alpha.AlphaAt(fx, fy).A > 0
}

func (f *Fog) Update() {
	f.x += f.vx
	if f.x < -f.gs.maze.Px() {
		f.passed = true
	}

	if !f.mutated {
		if f.mutateCount--; f.mutateCount < 0 {
			f.mutate()
			f.mutated = true
			f.image = f.inverse
		} else if f.mutateCount < 200 {
			f.flashCount += 10
			if f.flashCount > f.mutateCount {
				f.flashCount = 0
				if f.image == f.inverse {
					f.image = f.fog
				} else {
					f.image = f.inverse
				}
			}
		}
	}
}

func (f *Fog) Free() {
	f.fog.Free()
	f.inverse.Free()
}
