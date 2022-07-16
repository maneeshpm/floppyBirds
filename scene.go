package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	tex *sdl.Texture
}

func newScene(ren *sdl.Renderer) (*scene, error) {
	imageTex, err := img.LoadTexture(ren, bgPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create background texture: %v", err)
	}
	return &scene{tex: imageTex}, nil
}

func (s *scene) paint(ren *sdl.Renderer) error {
	ren.Clear()

	if err := ren.Copy(s.tex, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	ren.Present()
	return nil
}

func (s *scene) Destroy() {
	s.tex.Destroy()
}
