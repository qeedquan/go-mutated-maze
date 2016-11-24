package main

type Blitter struct {
	x, y  int
	image *Image
}

func (b *Blitter) Blit() {
	b.image.Blit(b.x, b.y)
}

type Entity struct {
	Blitter
	gs        *GameState
	vx, vy    int
	dvx, dvy  int
	isDead    bool
	dir       string
	wantedDir string
}