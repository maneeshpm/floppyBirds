package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

func drawEnd(ren *sdl.Renderer, score int) error {
	ren.Clear()
	font, err := ttf.OpenFont(fontPath, fontSize)

	if err != nil {
		return fmt.Errorf("Could not open font: %v", err)
	}
	defer font.Close()

	c1 := sdl.Color{R: 0, G: 200, B: 225, A: 0}
	c2 := sdl.Color{R: 50, G: 160, B: 80, A: 0}

	surface1, err := font.RenderUTF8Solid("LMAO ded", c1)
	if err != nil {
		return fmt.Errorf("Could not render text: %v", err)
	}
	defer surface1.Free()

	surface2, err := font.RenderUTF8Solid(fmt.Sprintf("Shawtty sold at %d crore", score), c2)
	if err != nil {
		return fmt.Errorf("Could not render text: %v", err)
	}
	defer surface2.Free()

	t1, err := ren.CreateTextureFromSurface(surface1)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	t2, err := ren.CreateTextureFromSurface(surface2)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer t1.Destroy()
	defer t2.Destroy()

	rect1 := &sdl.Rect{W: 1000, H: 500, X: 300, Y: 100}
	if err := ren.Copy(t1, nil, rect1); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	rect2 := &sdl.Rect{W: 800, H: 200, X: 400, Y: 700}
	if err := ren.Copy(t2, nil, rect2); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	ren.Present()
	return nil
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
					drawEnd(ren, s.card.score)
					time.Sleep(3 * time.Second)
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
