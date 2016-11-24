package main

import (
	"image"
	"image/draw"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Text struct {
	text  string
	color sdl.Color
}

type fontKey struct {
	filename string
	ptsize   int
}

var (
	fonts = make(map[fontKey]*sdlttf.Font)
)

func loadFont(name string, ptsize int) *sdlttf.Font {
	log.SetPrefix("font: ")
	filename := filepath.Join(*dataDir, name)

	key := fontKey{filename, ptsize}
	if font, found := fonts[key]; found {
		return font
	}

	font, err := sdlttf.OpenFont(filename, ptsize)
	if err != nil {
		log.Fatal(err)
	}

	fonts[key] = font
	return font
}

func printCenter(font *sdlttf.Font, y int, color sdl.Color, text string) {
	width, _, _ := font.SizeUTF8(text)
	blitText(font, int((int(DisplaySize.X)-width)/2), y, color, text)
}

func blitText(font *sdlttf.Font, x, y int, c sdl.Color, text string) {
	r, err := font.RenderUTF8BlendedEx(surface, text, c)
	if err != nil {
		log.Fatal(err)
	}

	p, err := texture.Lock(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = surface.Lock()
	if err != nil {
		log.Fatal(err)
	}
	s := surface.Pixels()
	for i := 0; i < len(p); i += 4 {
		p[i] = s[i+2]
		p[i+1] = s[i]
		p[i+2] = s[i+1]
		p[i+3] = s[i+3]
	}

	surface.Unlock()
	texture.Unlock()

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	screen.Copy(texture, &sdl.Rect{0, 0, r.W, r.H}, &sdl.Rect{int32(x), int32(y), r.W, r.H})
}

type Image struct {
	texture *sdl.Texture
	alpha   *image.Alpha
	width   int
	height  int
	angle   float64
}

func loadImage(name string, xforms ...func(m image.Image) image.Image) *Image {
	filename := filepath.Join(*dataDir, name)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, xform := range xforms {
		img = xform(img)
	}

	r := img.Bounds()
	alpha := image.NewAlpha(image.Rect(0, 0, r.Dx(), r.Dy()))
	draw.Draw(alpha, alpha.Bounds(), img, image.ZP, draw.Src)

	texture, err := sdlimage.LoadTextureImage(screen.Renderer, img)
	if err != nil {
		log.Fatal(err)
	}

	_, _, width, height, _ := texture.Query()
	return &Image{
		texture: texture,
		alpha:   alpha,
		width:   width,
		height:  height,
	}
}

func makeImage(width, height int) *Image {
	texture, err := screen.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, width, height)
	if err != nil {
		log.Fatal(err)
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	return &Image{
		texture: texture,
		width:   width,
		height:  height,
	}
}

func (m *Image) Free() {
	m.texture.Destroy()
}

func (m *Image) Blit(x, y int) {
	screen.CopyEx(m.texture, nil, &sdl.Rect{int32(x), int32(y), int32(m.width), int32(m.height)}, -m.angle, nil, sdl.FLIP_NONE)
}

func (m *Image) BlitArea(x, y int, src sdl.Rect) {
	dst := sdl.Rect{int32(x), int32(y), src.W, src.H}
	screen.CopyEx(m.texture, &src, &dst, -m.angle, nil, sdl.FLIP_NONE)
}

func (m *Image) Bind() {
	log.SetPrefix("image: ")
	err := screen.SetTarget(m.texture)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Image) Unbind() {
	log.SetPrefix("image: ")
	err := screen.SetTarget(nil)
	if err != nil {
		log.Fatal(err)
	}
}

func colorBlackRandom(palette []sdl.Color) func(m image.Image) image.Image {
	return func(m image.Image) image.Image {
		p := image.NewRGBA(m.Bounds())
		draw.Draw(p, p.Bounds(), m, image.ZP, draw.Src)

		rc := palette[rand.Intn(len(palette))]
		r := p.Bounds()
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				_, _, _, ca := p.RGBAAt(x, y).RGBA()
				if ca != 0 {
					p.Set(x, y, rc)
				}
			}
		}

		return p
	}
}

func scaleImage2x(m image.Image) image.Image {
	r := m.Bounds()
	f := scaleImage(r.Dx(), r.Dy())
	return f(m)
}

func scaleImage(width, height int) func(m image.Image) image.Image {
	return func(m image.Image) image.Image {
		p := image.NewRGBA(image.Rect(0, 0, width, height))

		r := p.Bounds()
		s := m.Bounds()

		dw2 := r.Dx() * 2
		dh2 := r.Dy() * 2
		sw2 := r.Dx() * 2
		sh2 := r.Dy() * 2

		sy := s.Min.Y
		h := sh2 - dh2
		for y := r.Min.Y; y < r.Max.Y; y++ {
			w := sw2 - dw2
			sx := 0
			for x := r.Min.X; x < r.Max.X; x++ {
				p.Set(x, y, m.At(sx, sy))
				for w >= 0 {
					sx++
					w -= dw2
				}
				w += sw2
			}

			for h >= 0 {
				sy++
				h -= dh2
			}
			h += sh2
		}
		return p
	}
}
