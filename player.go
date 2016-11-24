package main

import (
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

const (
	FollowDelay = 6
	MaxFollows  = 16
)

type RailsThing struct {
	Entity

	posHist   [FollowDelay * MaxFollows]sdl.Point
	followers []interface{}

	death  *sdlmixer.Chunk
	pickup *sdlmixer.Chunk
}

func newRailsThing(gs *GameState, x, y int) *RailsThing {
	r := &RailsThing{}
	r.gs = gs
	r.x, r.y = x, y
	r.image = loadImage("hero_dude.png", colorBlackRandom(SpriteColors))
	r.death = loadSound("death.wav")
	r.pickup = loadSound("pickup_friend.wav")

	for i := range r.posHist {
		r.posHist[i] = sdl.Point{int32(r.x), int32(r.y)}
	}

	return r
}

func (r *RailsThing) Free() {
	r.image.Free()
}

func (r *RailsThing) Hitbox() sdl.Rect {
	return sdl.Rect{int32(r.x) + 4, int32(r.y) + 3, 8, 10}
}

func (r *RailsThing) Update() {
	if r.isDead {
		r.updateDead()
	} else {
		r.updateAlive()
	}
	r.addPosHistory(r.x, r.y)
}

func (r *RailsThing) updateDead() {
	r.image.angle = math.Mod(r.image.angle+90, 360)
	r.x, r.y = r.x+r.vx, r.y+r.vy
	r.vy++
}

func (r *RailsThing) updateAlive() {
	s := r.gs
	k := sdl.GetKeyboardState()
	switch {
	case k[sdl.SCANCODE_UP] != 0:
		r.wantedDir = "up"
	case k[sdl.SCANCODE_DOWN] != 0:
		r.wantedDir = "down"
	case k[sdl.SCANCODE_LEFT] != 0:
		r.wantedDir = "left"
	case k[sdl.SCANCODE_RIGHT] != 0:
		r.wantedDir = "right"
	}

	switch {
	case r.dir == "up" && r.wantedDir == "down",
		r.dir == "down" && r.wantedDir == "up",
		r.dir == "left" && r.wantedDir == "right",
		r.dir == "right" && r.wantedDir == "left":
		r.dir = r.wantedDir
	}

	if r.x%16 == 0 && r.y%16 == 0 {
		n := s.maze.Node(r.x/16, r.y/16)
		switch {
		case r.wantedDir == "up" && n.open[Up],
			r.wantedDir == "down" && n.open[Down],
			r.wantedDir == "left" && n.open[Left],
			r.wantedDir == "right" && n.open[Right]:
			r.dir = r.wantedDir
		}
	}

	dx, dy := 0, 0
	mx, my := r.x/16, r.y/16
	n := s.maze.Node(mx, my)
	switch r.dir {
	case "up":
		dy = -1
		if r.y%16 == 0 && !n.open[Up] {
			dy = 0
		}

	case "down":
		dy = 1
		if r.y%16 == 0 && !n.open[Down] {
			dy = 0
		}

	case "left":
		dx = -1
		if r.x%16 == 0 && !n.open[Left] {
			dx = 0
		}

	case "right":
		dx = 1
		if r.x%16 == 0 && !n.open[Right] {
			dx = 0
		}
	}

	r.x, r.y = r.x+dx*2, r.y+dy*2
}

func (r *RailsThing) addPosHistory(x, y int) {
	l := len(r.posHist) - 1
	if p := r.posHist[l]; int32(x) == p.X && int32(y) == p.Y {
		return
	}
	copy(r.posHist[:], r.posHist[1:])
	r.posHist[l] = sdl.Point{int32(x), int32(y)}
}

func (r *RailsThing) RegisterFollower(f interface{}) {
	for i := range r.followers {
		if r.followers[i] == f {
			return
		}
	}
	r.followers = append(r.followers, f)
}

func (r *RailsThing) followerIndex(f interface{}) int {
	for i := range r.followers {
		if r.followers[i] == f {
			return i
		}
	}
	r.followers = append(r.followers, f)
	return len(r.followers)
}

func (r *RailsThing) FollowerPos(f interface{}) (x, y int) {
	i := FollowDelay * (r.followerIndex(f) + 1)
	p := r.posHist[len(r.posHist)-i]
	return int(p.X), int(p.Y)
}

func (r *RailsThing) OnMobHit(m *Mob) {
	switch m.kind {
	case "friendly":
		found := false
		for i := range r.followers {
			if r.followers[i] == m {
				found = true
				break
			}
		}
		if !found {
			playSound(r.pickup)
			m.Follow(r)
		}

	case "enemy":
		playSound(r.death)
		r.Die(m)
	}
}

func (r *RailsThing) Die(killingMob *Mob) {
	// try to kill one of our followers first
	if killingMob != nil {
		for i := range r.followers {
			if m, ok := r.followers[i].(*Mob); ok {
				m.Die()
				m.Free()

				l := len(r.followers) - 1
				r.followers[i], r.followers = r.followers[l], r.followers[:l]

				killingMob.Die()
				return
			}
		}
	}

	r.isDead = true
	r.dvx = 3 - rand.Intn(7)
	r.dvy = -3
}

func (r *RailsThing) Unfollow(x interface{}) {
	for i := range r.followers {
		if r.followers[i] == x {
			l := len(r.followers) - 1
			r.followers[i], r.followers = r.followers[l], r.followers[:l]
		}
	}
}

func (r *RailsThing) MobFollowers() []*Mob {
	var m []*Mob

	for _, f := range r.followers {
		if p, ok := f.(*Mob); ok {
			m = append(m, p)
		}
	}
	return m
}
