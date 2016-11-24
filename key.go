package main

import "github.com/qeedquan/go-media/sdl"

type FollowKey interface {
	FollowerPos(k interface{}) (x, y int)
}

type Key struct {
	Blitter
	follow FollowKey
}

func newKey(x, y int) *Key {
	k := &Key{}
	k.x, k.y = x, y
	k.image = loadImage("key.png")
	return k
}

func (k *Key) Hitbox() sdl.Rect {
	return sdl.Rect{int32(k.x) + 4, int32(k.y) + 4, 8, 8}
}

func (k *Key) Follow(f FollowKey) {
	k.follow = f
}

func (k *Key) Update() {
	if k.follow == nil {
		return
	}
	k.x, k.y = k.follow.FollowerPos(k)
}
