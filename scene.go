package main

import (
	"fmt"
	"log"
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

func (s *scene) update() {
	s.bird.update()
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

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true

	case *sdl.MouseButtonEvent:
		s.bird.jump()

	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.CommonEvent, *sdl.AudioDeviceEvent:

	default:
		log.Printf("Event type: %T", event)
	}

	return false
}

func (s *scene) run(ren *sdl.Renderer, events <-chan sdl.Event) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		done := false
		for !done {
			select {
			case event := <-events:
				done = s.handleEvent(event)
			case <-tick:
				s.update()
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
