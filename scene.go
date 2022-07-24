package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg *sdl.Texture

	bird  *bird
	pipes *pipes
	card  *card
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

	pipes, err := newPipes(ren)
	if err != nil {
		return nil, fmt.Errorf("Could not create pipe: %v", err)
	}

	card, err := newCard(ren)
	if err != nil {
		return nil, fmt.Errorf("Could not create card: %v", err)
	}

	return &scene{bg: bg, bird: bird, pipes: pipes, card: card}, nil
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update(s.card)

	s.pipes.check(s.bird)
}

func (s *scene) reset() {
	s.bird.reset()
	s.pipes.reset()
	s.card.reset()
}

func (s *scene) paint(ren *sdl.Renderer) error {
	ren.Clear()

	if err := ren.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	if err := s.bird.paint(ren); err != nil {
		return fmt.Errorf("Could not paint bird: %v", err)
	}
	if err := s.pipes.paint(ren); err != nil {
		return fmt.Errorf("Could not paint pipe: %v", err)
	}
	if err := s.card.paint(ren); err != nil {
		return fmt.Errorf("Could not paint pipe: %v", err)
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
				if s.bird.isDead() {
					drawText(ren, "LMAO ded")
					time.Sleep(2 * time.Second)
					s.reset()
				}
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
	s.pipes.destroy()
	s.card.destroy()
}
