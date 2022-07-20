package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	framePath = "res/imgs/frame-%d.png"
	gravity   = 0.2
	jumpSpeed = 8
)

type bird struct {
	mu sync.RWMutex

	dead     bool
	iter     int
	textures []*sdl.Texture

	y, yspeed float64
	x, xspeed float64
}

func newBird(ren *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		texture, err := img.LoadTexture(ren, fmt.Sprintf(framePath, i))
		if err != nil {
			return nil, fmt.Errorf("Could not could not load frame: %v", err)
		}
		textures = append(textures, texture)
	}

	return &bird{textures: textures, x: 10, y: winHeight / 2, xspeed: 1}, nil
}

func (bird *bird) isDead() bool {
	bird.mu.RLock()
	defer bird.mu.RUnlock()
	return bird.dead
}

func (bird *bird) reset() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.y = winHeight / 2
	bird.yspeed = 0
	bird.dead = false
}

func (bird *bird) update() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.iter++

	bird.yspeed += gravity
	bird.y -= bird.yspeed
	if bird.y < 0 {
		bird.dead = true
	}

	bird.x += bird.xspeed
}

func (bird *bird) paint(ren *sdl.Renderer) error {
	bird.mu.RLock()
	defer bird.mu.RUnlock()

	rect := &sdl.Rect{W: 200, H: 150, X: int32(bird.x), Y: (winHeight - int32(bird.y)) - 43/2}
	frameSelector := bird.iter / 10 % len(bird.textures)

	if err := ren.Copy(bird.textures[frameSelector], nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}

func (bird *bird) jump() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.yspeed = -jumpSpeed
}

func (bird *bird) destroy() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	for _, texture := range bird.textures {
		texture.Destroy()
	}
}
