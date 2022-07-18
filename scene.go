package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg   *sdl.Texture
	bird *bird
}

func newScene(ren *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(ren, bgPath)
	if err != nil {
		return nil, fmt.Errorf("Could not create background texture: %v", err)
	}

	bird, err := newBird(ren)
	if err != nil {
		return nil, fmt.Errorf("Could not create bird: %v", err)
	}

	return &scene{bg: bg, bird: bird}, nil
}

func (s *scene) paint(ren *sdl.Renderer) error {
	ren.Clear()

	if err := ren.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	if err := s.bird.paint(ren); err != nil {
		return fmt.Errorf("Could not paint biard: %v", err)
	}

	ren.Present()
	return nil
}

func (s *scene) run(ren *sdl.Renderer, ctx context.Context) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
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
	s.bird.destroy()
}
