package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	framePath = "res/imgs/frame-%d.png"
)

type bird struct {
	iter     int
	textures []*sdl.Texture
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

	return &bird{textures: textures}, nil
}

func (bird *bird) paint(ren *sdl.Renderer) error {
	rect := &sdl.Rect{W: 100, H: 86, X: 10, Y: winHeight/2 - 43/2}

	bird.iter++
	frameSelector := bird.iter / 10 % len(bird.textures)

	if err := ren.Copy(bird.textures[frameSelector], nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}

func (bird *bird) destroy() {
	for _, texture := range bird.textures {
		texture.Destroy()
	}
}
