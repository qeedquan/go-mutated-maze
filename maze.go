package main

import (
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/sdl"
)

var (
	MazeColor = BrightGreen
)

const (
	Up = iota
	Down
	Left
	Right
)

type MazeNode struct {
	id     int
	x, y   int
	color  sdl.Color
	parent *Maze
	open   [4]bool
}

func newMazeNode(id, x, y int, parent *Maze) *MazeNode {
	return &MazeNode{
		id:     id,
		x:      x,
		y:      y,
		color:  MazeColor,
		parent: parent,
	}
}

func (m *MazeNode) Hitbox() sdl.Rect {
	return sdl.Rect{int32(m.x) * 16, int32(m.y) * 16, 16, 16}
}

func (m *MazeNode) Node(dir int) *MazeNode {
	switch dir {
	case Up:
		return m.parent.Node(m.x, m.y-1)
	case Down:
		return m.parent.Node(m.x, m.y+1)
	case Left:
		return m.parent.Node(m.x-1, m.y)
	case Right:
		return m.parent.Node(m.x+1, m.y)
	}

	panic("unreachable")
}

func (m *MazeNode) IsOpen(dir int) bool {
	return m.open[dir]
}

func (m *MazeNode) WallCount() int {
	n := 0
	for _, o := range m.open {
		if !o {
			n++
		}
	}
	return n
}

func (m *MazeNode) SetWall(dir int) {
	switch {
	case dir == Up && m.open[Up]:
		m.open[Up] = false
		o := m.Node(Up)
		o.open[Down] = false

	case dir == Down && m.open[Down]:
		m.open[Down] = false
		o := m.Node(Down)
		o.open[Up] = false

	case dir == Left && m.open[Left]:
		m.open[Left] = false
		o := m.Node(Left)
		o.open[Right] = false

	case dir == Right && m.open[Right]:
		m.open[Right] = false
		o := m.Node(Right)
		o.open[Left] = false
	}
}

func (m *MazeNode) HasNode(dir int) bool {
	switch dir {
	case Up:
		return m.y > 0
	case Down:
		return m.y < m.parent.height-1
	case Left:
		return m.x > 0
	case Right:
		return m.x < m.parent.width-1
	default:
		panic("unreachable")
	}
}

func (m *MazeNode) ClearWall(dir int) {
	switch {
	case dir == Up && !m.open[Up] && m.HasNode(Up):
		m.open[Up] = true
		o := m.Node(Up)
		o.open[Down] = true
	case dir == Down && !m.open[Down] && m.HasNode(Down):
		m.open[Down] = true
		o := m.Node(Down)
		o.open[Up] = true
	case dir == Left && !m.open[Left] && m.HasNode(Left):
		m.open[Left] = true
		o := m.Node(Left)
		o.open[Right] = true
	case dir == Right && !m.open[Right] && m.HasNode(Right):
		m.open[Right] = true
		o := m.Node(Right)
		o.open[Left] = true
	}
}

func (m *MazeNode) AvailDirs() []string {
	var n []string

	if m.open[Up] {
		n = append(n, "up")
	}
	if m.open[Down] {
		n = append(n, "down")
	}
	if m.open[Left] {
		n = append(n, "left")
	}
	if m.open[Right] {
		n = append(n, "right")
	}
	return n
}

func (m *MazeNode) OpenNodes() []*MazeNode {
	var n []*MazeNode

	if m.open[Up] {
		n = append(n, m.Node(Up))
	}
	if m.open[Down] {
		n = append(n, m.Node(Down))
	}
	if m.open[Left] {
		n = append(n, m.Node(Left))
	}
	if m.open[Right] {
		n = append(n, m.Node(Right))
	}
	return n
}

func (m *MazeNode) NearbyDirs() []int {
	var n []int

	if m.y > 0 {
		n = append(n, Up)
	}
	if m.y < m.parent.height-1 {
		n = append(n, Down)
	}
	if m.x > 0 {
		n = append(n, Left)
	}
	if m.x < m.parent.width-1 {
		n = append(n, Right)
	}
	return n
}

func (m *MazeNode) OpenAll() {
	for _, dir := range m.NearbyDirs() {
		m.ClearWall(dir)
	}
}

func (m *MazeNode) Pxy() (x, y int) {
	return m.x * 16, m.y * 16
}

func (m *MazeNode) DistanceFrom(n *MazeNode) float64 {
	x := float64(m.x - n.x)
	y := float64(m.y - n.y)
	return math.Sqrt(x*x + y*y)
}

type Maze struct {
	width  int
	height int
	nodes  []*MazeNode
}

func newMaze(width, height int) *Maze {
	m := &Maze{
		width:  width,
		height: height,
	}
	braidReset(m)
	return m
}

func (m *Maze) Px() int {
	return m.width * 16
}

func (m *Maze) Py() int {
	return m.height * 16
}

func (m *Maze) Node(x, y int) *MazeNode {
	if y < 0 {
		y = m.height - 1
	} else if y == m.height {
		y = 0
	}

	if x < 0 {
		x = m.width - 1
	} else if x == m.width {
		x = 0
	}

	return m.nodes[y*m.width+x]
}

func (m *Maze) Gen() {
	braidGen(m, nil)
}

func (m *Maze) RegenSelected(selected []*MazeNode) {
	braidRegenSelected(m, selected)
}

func (m *Maze) CollideNodes(r sdl.Rect) []*MazeNode {
	var p []*MazeNode
	for _, n := range m.nodes {
		if r.Collide(n.Hitbox()) {
			p = append(p, n)
		}
	}
	return p
}

func (m *Maze) Blit() {
	screen.SetDrawColor(sdl.Color{0x40, 0x40, 0x40, 0xFF})
	screen.FillRect(nil)

	wp := m.width * 16
	hp := m.height * 16
	screen.SetDrawColor(MazeColor)
	screen.DrawRect(&sdl.Rect{0, 0, int32(wp), int32(hp)})

	for _, n := range m.nodes {
		x, y := n.Pxy()
		screen.SetDrawColor(n.color)
		if !n.IsOpen(Up) {
			screen.DrawLine(x, y, x+15, y)
		}
		if !n.IsOpen(Down) {
			screen.DrawLine(x, y+15, x+15, y+15)
		}
		if !n.IsOpen(Left) {
			screen.DrawLine(x, y, x, y+15)
		}
		if !n.IsOpen(Right) {
			screen.DrawLine(x+15, y, x+15, y+15)
		}
	}
}

// http://www.astrolog.org/labyrnth/algrithm.htm
// Braid: To create a Maze without dead ends, basically add wall segments
// throughout the Maze at random, but ensure that each new segment added will
// not cause a dead end to be made. I make them with four steps: (1) Start with
// the outer wall, (2) Loop through the Maze and add single wall segments
// touching each wall vertex to ensure there are no open rooms or small "pole"
// walls in the Maze, (3) Loop over all possible wall segments in random
// order, adding a wall there if it wouldn't cause a dead end, (4) Either run
// the isolation remover utility at the end to make a legal Maze that has a
// solution, or be smarter in step three and make sure a wall is only added
// if it also wouldn't cause an isolated section.

func braidReset(m *Maze) {
	m.nodes = nil
	id := 0
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			n := newMazeNode(id, x, y, m)
			n.open[Up] = y != 0
			n.open[Down] = y != m.height-1
			n.open[Left] = x != 0
			n.open[Right] = x != m.width-1
			id++
			m.nodes = append(m.nodes, n)
		}
	}
}

func braidGen(m *Maze, nodes []*MazeNode) {
	if nodes == nil {
		nodes = m.nodes
	}

	type wall struct {
		node *MazeNode
		dir  int
	}
	var walls []wall
	for _, n := range nodes {
		for i, o := range n.open {
			if !o {
				continue
			}
			walls = append(walls, wall{n, i})
		}
	}
	for i := len(walls) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		walls[i], walls[j] = walls[j], walls[i]
	}

	for _, w := range walls {
		n := w.node
		if !n.IsOpen(w.dir) {
			continue
		}
		if n.WallCount() >= 2 {
			continue
		}

		o := n.Node(w.dir)
		if o.WallCount() >= 2 {
			continue
		}
		if !braidConnected(n, o) {
			continue
		}

		n.SetWall(w.dir)
	}
}

type braidSet map[*MazeNode]*MazeNode

func (s braidSet) Add(nodes ...*MazeNode) {
	for _, n := range nodes {
		s[n] = n
	}
}

func (s braidSet) Pop() *MazeNode {
	for k, v := range s {
		delete(s, k)
		return v
	}
	return nil
}

func braidConnected(f, s *MazeNode) bool {
	seen := make(braidSet)
	queue := make(braidSet)

	seen.Add(f, s)
	for _, x := range f.OpenNodes() {
		if x != f && x != s {
			queue.Add(x)
		}
	}

	for {
		n := queue.Pop()
		if n == nil {
			break
		}
		seen.Add(n)

		for _, x := range n.OpenNodes() {
			if x == s {
				return true
			}
			if seen[x] != nil {
				continue
			}
			queue.Add(x)
		}
	}
	return false
}

func braidRegenSelected(m *Maze, selected []*MazeNode) {
	for _, n := range selected {
		n.OpenAll()
	}
	braidGen(m, selected)
}
