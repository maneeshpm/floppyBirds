package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	framePath = "res/imgs/frame-%d.png"
	gravity   = 0.2
	jumpSpeed = 8
)

type bird struct {
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

func (bird *bird) paint(ren *sdl.Renderer) error {
	rect := &sdl.Rect{W: 200, H: 150, X: int32(bird.x), Y: (winHeight - int32(bird.y)) - 43/2}

	bird.iter++

	bird.yspeed += gravity
	bird.y -= bird.yspeed
	if bird.y < 0 {
		bird.yspeed *= -1
		bird.y = 0
	}

	bird.x += bird.xspeed
	frameSelector := bird.iter / 10 % len(bird.textures)

	if err := ren.Copy(bird.textures[frameSelector], nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}

func (bird *bird) jump() {
	bird.yspeed = -jumpSpeed
}

func (bird *bird) destroy() {
	for _, texture := range bird.textures {
		texture.Destroy()
	}
}
