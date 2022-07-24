package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	cardHeight = 150
	cardWidth  = 400
	cardX      = winWidth - cardWidth - 20
	cardY      = 10
)

type card struct {
	mu   sync.RWMutex
	tex  *sdl.Texture
	font *ttf.Font

	score int

	x, y, W, H int32
}

func getCardTexture(ren *sdl.Renderer, font *ttf.Font, score int) (*sdl.Texture, error) {
	text := fmt.Sprintf("CTC: %d Crore", score)
	c := sdl.Color{R: 225, G: 200, B: 225, A: 0}

	surface, err := font.RenderUTF8Solid(text, c)
	if err != nil {
		return nil, fmt.Errorf("Could not generate card surface: %v", err)
	}
	defer surface.Free()

	return ren.CreateTextureFromSurface(surface)
}

func newCard(ren *sdl.Renderer) (*card, error) {
	font, err := ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		return nil, fmt.Errorf("Could not open font: %v", err)
	}

	tex, err := getCardTexture(ren, font, 0)
	if err != nil {
		return nil, fmt.Errorf("Could not generate card texture: %v", err)
	}

	return &card{font: font, tex: tex, x: cardX, y: cardY, H: cardHeight, W: cardWidth}, nil
}

func (card *card) reset() {
	card.mu.Lock()
	defer card.mu.Unlock()

	card.score = 0
}

func (card *card) update() {
	card.mu.Lock()
	defer card.mu.Unlock()

	card.score++
}

func (card *card) paint(ren *sdl.Renderer) error {
	card.mu.Lock()
	defer card.mu.Unlock()

	tex, err := getCardTexture(ren, card.font, card.score)
	if err != nil {
		return fmt.Errorf("Could not generate card texture: %v", err)
	}

	card.tex = tex

	rect := &sdl.Rect{W: card.W, H: card.H, X: card.x, Y: card.y}

	if err := ren.Copy(card.tex, nil, rect); err != nil {
		return fmt.Errorf("Could not copy card texture: %v", err)
	}

	return nil
}

func (card *card) destroy() {
	card.mu.Lock()
	defer card.mu.Unlock()

	card.tex.Destroy()
}
