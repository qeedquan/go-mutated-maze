package main

import "github.com/qeedquan/go-media/sdl"

const (
	InitDistanceFromPlayer = 5
	NumInitMobs            = 15
	MaxCharge              = 10000
	ChargePerFrame         = 50
	SparkTimeMS            = 300
)

type Spark struct {
	time uint32
	x, y int
}

type GameState struct {
	score     int
	startTime uint32
	time      int
	timeLeft  int
	scoreTime int

	mobs           []*Mob
	mobsSaved      int
	mobsAvail      int
	mobsSavedTotal int
	deathCount     int

	level  int
	charge int

	maze   *Maze
	player *RailsThing
	sparks []*Spark
}

func newGameState() *GameState {
	return &GameState{}
}

func (g *GameState) Reset() {
	g.free()
	g.score = 0
	g.time = 60
	g.timeLeft = 0
	g.level = 0
	g.mobsSaved = 0
	g.mobsAvail = 0
	g.mobsSavedTotal = 0
	g.deathCount = 0
}

func (g *GameState) NextLevel() {
	g.free()
	if g.level++; g.level < 0 {
		g.level = 0
	}
	g.time = max(45, 120-(g.level-1)*20)
	g.mobsAvail = NumInitMobs + (g.level-1)*3
	g.charge = MaxCharge
}

func (g *GameState) AliveMobs() []*Mob {
	var p []*Mob
	for _, m := range g.mobs {
		if !m.isDead {
			p = append(p, m)
		}
	}
	return p
}

func (g *GameState) AddSpark(x, y int) {
	g.sparks = append(g.sparks, &Spark{sdl.GetTicks(), x, y})
}

func (g *GameState) free() {
	if g.player != nil {
		g.player.Free()
	}
	for i := range g.mobs {
		g.mobs[i].Free()
	}
	g.maze = nil
	g.player = nil
	g.mobs = nil
}
