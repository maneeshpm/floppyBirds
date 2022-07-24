package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	framePath = "res/imgs/frame-%d.png"
	gravity   = 0.3
	jumpSpeed = 8

	birdWidth  = 120
	birdHeight = 120
)

type bird struct {
	mu sync.RWMutex

	dead     bool
	iter     int
	textures []*sdl.Texture

	W, H int32

	x, y   int32
	yspeed float64
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

	return &bird{textures: textures, x: 50, y: winHeight / 2, W: birdWidth, H: birdHeight}, nil
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
	bird.y -= int32(bird.yspeed)
	if bird.y < 0 {
		bird.dead = true
	}
}

func (bird *bird) paint(ren *sdl.Renderer) error {
	bird.mu.RLock()
	defer bird.mu.RUnlock()

	rect := &sdl.Rect{W: bird.W, H: bird.H, X: int32(bird.x), Y: (winHeight - int32(bird.y)) - 43/2}
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

func (bird *bird) check(pipe *pipe) {
	bird.mu.Lock()
	defer bird.mu.Unlock()
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	if bird.x+bird.W < pipe.x {
		return
	}

	if pipe.x+pipe.w < bird.x {
		return
	}

	if pipe.inverted && winHeight-pipe.h > bird.y+bird.H/2 {
		return
	}

	if !pipe.inverted && bird.y-bird.H/2 > pipe.h {
		return
	}

	bird.dead = true
}

func (bird *bird) destroy() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	for _, texture := range bird.textures {
		texture.Destroy()
	}
}
