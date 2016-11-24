package main

import "github.com/qeedquan/go-media/sdl"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mousePos() (x, y int) {
	mx, my, _ := sdl.GetMouseState()
	sx, sy := screen.Size()
	lx, ly := screen.LogicalSize()
	x = int(float64(mx) / float64(sx) * float64(lx))
	y = int(float64(my) / float64(sy) * float64(ly))
	return
}
