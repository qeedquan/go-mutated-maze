package main

import (
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

type Mob struct {
	Entity

	kind             string
	state            string
	cvx, cvy         int
	waitCount        int
	exclamationCount int

	good        *Image
	bad         *Image
	exclamation *Image

	follow *RailsThing

	exclaim *sdlmixer.Chunk
}

func newMob(gs *GameState, x, y int) *Mob {
	m := &Mob{}
	m.gs = gs
	m.x, m.y = x, y
	m.good = loadImage("other_dude.png", colorBlackRandom(SpriteColors))
	m.bad = loadImage("ghost.png", colorBlackRandom(BadColors))
	m.exclamation = loadImage("exclamation.png")
	m.setFriendly()
	return m
}

func (m *Mob) Free() {
	m.good.Free()
	m.bad.Free()
	m.exclamation.Free()
}

func (m *Mob) Hitbox() sdl.Rect {
	return sdl.Rect{int32(m.x) + 4, int32(m.y) + 3, 8, 10}
}

func (m *Mob) setFriendly() {
	m.kind = "friendly"
	m.image = m.good
	m.initNothing()
}

func (m *Mob) setEnemy() {
	m.kind = "enemy"
	m.image = m.bad
	m.initNothing()
}

func (m *Mob) initNothing() {
	m.state = "nothing"
}

func (m *Mob) Update() {
	switch {
	case m.isDead:
		m.image.angle = math.Mod(m.image.angle+90, 360)
		m.x, m.y = m.x+m.dvx, m.y+m.dvy
		m.dvy++
	case m.kind == "friendly":
		m.friendlyUpdate()
	case m.kind == "enemy":
		m.enemyUpdate()
	}
}

func (m *Mob) initWait() {
	m.state = "wait"
	m.waitCount = 3 + rand.Intn(6)
}

func (m *Mob) initRandom() {
	m.state = "random"

	// otherwise only change if right in middle of node
	n := m.gs.maze.Node(m.x/16, m.y/16)
	l := n.AvailDirs()
	m.wantedDir = l[rand.Intn(len(l))]
	m.dir = m.wantedDir
}

func (m *Mob) initChase() {
	m.state = "chase"
}

func (m *Mob) friendlyUpdate() {
	if m.state == "nothing" {
		f := []func(){m.initWait, m.initRandom}
		f[rand.Intn(len(f))]()
	}

	switch m.state {
	case "random":
		m.updateRandom()
	case "wait":
		m.updateWait()
	case "follow_player":
		m.updateFollowPlayer()
	}
}

func (m *Mob) enemyUpdate() {
	if m.state == "nothing" {
		f := []func(){m.initWait, m.initRandom}
		f[rand.Intn(len(f))]()
	}
	m.detectPlayer()

	switch m.state {
	case "random":
		m.updateRandom()
	case "wait":
		m.updateWait()
	case "exclamation":
		m.updateExclamation()
	case "chase":
		m.updateChase()
	}
}

func (m *Mob) updateRandom() {
	dirs := map[string]sdl.Point{
		"up":    {0, -1},
		"down":  {0, 1},
		"left":  {-1, 0},
		"right": {1, 0},
	}

	p := dirs[m.dir]
	m.x, m.y = m.x+int(p.X), m.y+int(p.Y)

	// otherwise only change if right in middle of node
	if m.x%16 == 0 && m.y%16 == 0 {
		m.initNothing()
		return
	}
}

func (m *Mob) updateWait() {
	if m.waitCount--; m.waitCount <= 0 {
		m.initNothing()
	}
}

func (m *Mob) updateFollowPlayer() {
	if m.follow == nil {
		return
	}
	m.x, m.y = m.follow.FollowerPos(m)
}

func (m *Mob) updateChase() {
	if m.x%16 == 0 && m.y%16 == 0 {
		n := m.gs.maze.Node(m.x/16, m.y/16)
		if m.cvx < 0 {
			if !n.open[Left] {
				m.initNothing()
				return
			}
		}

		if m.cvx > 0 {
			if !n.open[Right] {
				m.initNothing()
				return
			}
		}

		if m.cvy < 0 {
			if !n.open[Up] {
				m.initNothing()
				return
			}
		}

		if m.cvy > 0 {
			if !n.open[Down] {
				m.initNothing()
				return
			}
		}
	}

	m.x += m.cvx
	m.y += m.cvy
}

func (m *Mob) initExclamation(x, y int) {
	m.state = "exclamation"
	m.x >>= 1
	m.x <<= 1
	m.y >>= 1
	m.y <<= 1
	m.cvx, m.cvy = x, y
	playSound(m.exclaim)
}

func (m *Mob) updateExclamation() {
	if m.exclamationCount--; m.exclamationCount <= 0 {
		m.initChase()
	}
}

func (m *Mob) detectDirFunc(dir int) func(*MazeNode) *MazeNode {
	return func(x *MazeNode) *MazeNode {
		if x.open[dir] {
			return x.Node(dir)
		}
		return nil
	}
}

func (m *Mob) detectPlayer() {
	if m.state == "exclamation" {
		return
	}

	if m.y%16 == 0 {
		if m.findPlayerOn(m.detectDirFunc(Left)) {
			if m.state == "chase" && m.cvx < 0 {
				return
			}
			m.initExclamation(-2, 0)
			return
		}

		if m.findPlayerOn(m.detectDirFunc(Right)) {
			if m.state == "chase" && m.cvx > 0 {
				return
			}
			m.initExclamation(2, 0)
			return
		}
	}

	if m.x%16 == 0 {
		if m.findPlayerOn(m.detectDirFunc(Up)) {
			if m.state == "chase" && m.cvy < 0 {
				return
			}
			m.initExclamation(0, -2)
			return
		}

		if m.findPlayerOn(m.detectDirFunc(Down)) {
			if m.state == "chase" && m.cvy > 0 {
				return
			}
			m.initExclamation(0, 2)
			return
		}
	}
}

func (m *Mob) findPlayerOn(nextNode func(*MazeNode) *MazeNode) bool {
	s := m.gs

	n := s.maze.Node(m.x/16, m.y/16)
	p := s.player.Hitbox()
	for n != nil {
		if n.Hitbox().Collide(p) {
			return true
		}
		n = nextNode(n)
	}

	return false
}

func (m *Mob) ToggleKind() {
	switch m.kind {
	case "friendly":
		m.setEnemy()
	case "enemy":
		m.setFriendly()
	}
}

func (m *Mob) FogMutate(x, y int) {
	m.gs.AddSpark(x, y)
	if m.state == "follow_player" {
		m.follow.Unfollow(m)
	}

	m.initNothing()
	m.x, m.y = x, y
	m.ToggleKind()
}

func (m *Mob) Follow(otherSprite *RailsThing) {
	m.follow = otherSprite
	m.state = "follow_player"
}

func (m *Mob) Die() {
	if m.kind == "enemy" {
		m.gs.score += 200
	}

	m.isDead = true
	m.dvx = 3 - rand.Intn(7)
	m.dvy = -3
}

func (m *Mob) Blit() {
	m.Blitter.Blit()
	if m.state == "exclamation" {
		m.exclamation.Blit(m.x, m.y-8)
	}
}
