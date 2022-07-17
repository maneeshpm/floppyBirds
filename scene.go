package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	framePath = "res/imgs/frame-%d.png"
)

type scene struct {
	iter  int
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func newScene(ren *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(ren, bgPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create background texture: %v", err)
	}

	var birds []*sdl.Texture
	for i := 1; i <= 4; i++ {
		bird, err := img.LoadTexture(ren, fmt.Sprintf(framePath, i))
		if err != nil {
			return nil, fmt.Errorf("Could not could not load frame: %v", err)
		}
		birds = append(birds, bird)
	}

	return &scene{bg: bg, birds: birds}, nil
}

func (s *scene) paint(ren *sdl.Renderer) error {
	ren.Clear()
	s.iter++
	s.iter %= len(s.birds)

	if err := ren.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	rect := &sdl.Rect{W: 100, H: 86, X: 10, Y: winHeight/2 - 43/2}

	if err := ren.Copy(s.birds[s.iter], nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	ren.Present()
	return nil
}

func (s *scene) run(ren *sdl.Renderer, ctx context.Context) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(50 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(ren); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) Destroy() {
	s.bg.Destroy()

	for _, bird := range s.birds {
		bird.Destroy()
	}
}
