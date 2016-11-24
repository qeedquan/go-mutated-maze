package main

import (
	"fmt"
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

var (
	bgColor  = sdl.Color{0x1c, 0x1c, 0x1c, 0xff}
	mazePos  = sdl.Point{16, 16}
	scorePos = sdl.Point{380, 25}
)

type Game struct {
	gs         *GameState
	fog        *Fog
	lockedDoor *Blitter
	key        *Key

	textFont   *sdlttf.Font
	timeupFont *sdlttf.Font

	maze    *Image
	mutator *Image
	charge  *Image
	spark   *Image

	openDoor  *sdlmixer.Chunk
	pickupKey *sdlmixer.Chunk
}

func newGame(gs *GameState) *Game {
	g := &Game{
		gs: gs,

		maze:    makeImage(22*16, 18*16),
		mutator: makeImage(64, 64),
		charge:  loadImage("chargebar.png", scaleImage2x),
		spark:   loadImage("spark.png"),

		textFont:   loadFont("arial.ttf", 15),
		timeupFont: loadFont("arial.ttf", 30),

		openDoor:  loadSound("open_door.wav"),
		pickupKey: loadSound("pickup_key.wav"),
	}
	return g
}

func (g *Game) init() {
	s := g.gs
	s.NextLevel()

	g.mutator.Bind()
	screen.SetDrawColor(sdl.Color{255, 255, 255, 30})
	screen.FillRect(nil)
	g.mutator.Unbind()

	m := newMaze(22, 18)
	m.Gen()
	s.maze = m

	pn := m.Node(0, 0)
	px, py := pn.Pxy()
	s.player = newRailsThing(s, px, py)

	ln := m.nodes[len(m.nodes)-1]
	lx, ly := ln.Pxy()
	g.lockedDoor = &Blitter{lx, ly, loadImage("lock.png")}

	var availNodes []*MazeNode
	for _, n := range s.maze.nodes {
		if n.DistanceFrom(pn) > InitDistanceFromPlayer {
			availNodes = append(availNodes, n)
		}
	}
	for i := len(availNodes) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		availNodes[i], availNodes[j] = availNodes[j], availNodes[i]
	}

	for i := 0; i < s.mobsAvail; i++ {
		l := len(availNodes) - 1
		n := availNodes[l]
		availNodes = availNodes[:l]

		mx, my := n.Pxy()
		mob := newMob(s, mx, my)
		if rand.Float64() <= 1/3.0 {
			mob.ToggleKind()
		}
		s.mobs = append(s.mobs, mob)
	}

	var availKeyNodes []*MazeNode
	for _, n := range m.nodes {
		if n.DistanceFrom(pn) > InitDistanceFromPlayer &&
			n.DistanceFrom(ln) > InitDistanceFromPlayer {
			availKeyNodes = append(availKeyNodes, n)
		}
	}
	kn := availKeyNodes[rand.Intn(len(availKeyNodes))]
	kx, ky := kn.Pxy()
	g.key = newKey(kx, ky)

	g.fog = newFog(s)

	s.startTime = sdl.GetTicks()
}

func (g *Game) event() string {
	s := g.gs
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
				fmt.Println("No soup for you GruikInc!")
			}
		case sdl.MouseButtonDownEvent:
			if ev.Button == 1 && s.charge == MaxCharge {
				playSound(g.fog.mutateSound)

				cx, cy := mousePos()
				r := sdl.Rect{int32(cx) - mazePos.X, int32(cy) - mazePos.Y,
					int32(g.mutator.width), int32(g.mutator.height)}
				selected := s.maze.CollideNodes(r)
				s.maze.RegenSelected(selected)

				if !*chargeInf {
					s.charge = 0
				}
			}
		}
	}
	return ""
}

func (g *Game) update() string {
	s := g.gs

	current := sdl.GetTicks()
	s.timeLeft = max(s.time-int((current-s.startTime)/1000), 0)
	if s.timeLeft <= 0 && !s.player.isDead {
	}
	s.player.Update()
	g.key.Update()

	for _, m := range s.mobs {
		m.Update()
	}

	for _, m := range s.AliveMobs() {
		if s.player.Hitbox().Collide(m.Hitbox()) {
			s.player.OnMobHit(m)
		}
	}

	if g.key.follow != s.player && s.player.Hitbox().Collide(g.key.Hitbox()) {
		playSound(g.pickupKey)
		s.score += 200
		g.key.Follow(s.player)
	}

	if g.key.follow == s.player && g.lockedDoor.x == s.player.x && g.lockedDoor.y == s.player.y {
		playSound(g.openDoor)
		s.score += 200
		s.mobsSaved = len(s.player.MobFollowers())
		s.mobsSavedTotal += s.mobsSaved
		return "win"
	}

	g.fog.Update()
	if g.fog.passed {
		g.fog.Free()
		g.fog = newFog(s)
	}

	if s.player.isDead {
		if s.deathCount++; s.deathCount > 100 {
			return "gameover"
		}
	}

	s.charge = min(MaxCharge, s.charge+ChargePerFrame)

	curTime := sdl.GetTicks()
	var sparks []*Spark
	for _, s := range s.sparks {
		if s.time-curTime < SparkTimeMS {
			sparks = append(sparks, s)
		}
	}
	s.sparks = sparks

	return ""
}

func (g *Game) draw() {
	s := g.gs

	screen.SetDrawColor(bgColor)
	screen.Clear()

	g.maze.Bind()
	s.maze.Blit()
	g.lockedDoor.Blit()
	g.key.Blit()
	for _, m := range s.mobs {
		m.Blit()
	}
	s.player.Blit()
	g.fog.Blit()

	for _, sp := range s.sparks {
		g.spark.Blit(sp.x, sp.y)
	}

	g.maze.Unbind()
	g.maze.Blit(int(mazePos.X), int(mazePos.Y))

	mx, my := mousePos()
	g.mutator.Blit(mx, my)

	g.blitStats()

	chargeWidth := int(float64(g.charge.width) * (float64(s.charge) / MaxCharge))
	chargeArea := sdl.Rect{W: int32(chargeWidth), H: int32(g.charge.height)}
	g.charge.BlitArea(380, 280, chargeArea)

	screen.Present()
}

func (g *Game) blitStats() {
	s := g.gs
	textList := []Text{
		{fmt.Sprint("Level: ", s.level), BrightMagenta},
		{fmt.Sprint("Score: ", s.score), BrightMagenta},
	}
	c := BrightMagenta
	if s.timeLeft < 15 && s.timeLeft&1 != 0 {
		c = BrightWhite
	}
	textList = append(textList, Text{
		fmt.Sprint("Time: ", s.timeLeft), c,
	})

	x, y := scorePos.X, scorePos.Y
	for _, l := range textList {
		blitText(g.textFont, int(x), int(y), l.color, l.text)
		y += int32(g.textFont.Height())
	}

	if s.timeLeft <= 0 {
		blitText(g.timeupFont, 128, 100, BrightMagenta, "TIME UP!")
	}
}
