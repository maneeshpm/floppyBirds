package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	pipeSpeed  = 2
	pipeWidth  = 80
	pipeHeight = 200

	pipePath = "res/imgs/pipe.png"
)

type pipe struct {
	mu sync.RWMutex

	x int32
	h int32
	w int32

	inverted bool
}

type pipes struct {
	texture *sdl.Texture
	mu      sync.RWMutex

	speed int32
	ps    []*pipe
}

func newPipes(ren *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(ren, pipePath)
	if err != nil {
		return nil, fmt.Errorf("Could not could not load pipe: %v", err)
	}

	pipes := &pipes{
		texture: texture,
		speed:   pipeSpeed,
	}

	go func() {
		for {
			pipes.mu.Lock()
			pipes.ps = append(pipes.ps, newPipe())
			pipes.mu.Unlock()

			time.Sleep(2 * time.Second)
		}
	}()

	return pipes, nil
}

func newPipe() *pipe {
	return &pipe{
		x: winWidth,
		h: pipeHeight + int32(rand.Intn(300)),
		w: pipeWidth,

		inverted: (rand.Int63()%2 == 0),
	}
}

func (pipes *pipes) reset() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	pipes.ps = nil
}

func (pipes *pipes) update(card *card) {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	var rem []*pipe
	for _, pipe := range pipes.ps {
		pipe.mu.Lock()
		pipe.x -= pipes.speed
		pipe.mu.Unlock()

		if pipe.x+pipe.w > 0 {
			rem = append(rem, pipe)
		} else {
			card.update()
		}
	}
	pipes.ps = rem
}

func (pipes *pipes) paint(ren *sdl.Renderer) error {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	for _, pipe := range pipes.ps {
		if err := pipe.paint(ren, pipes.texture); err != nil {
			return err
		}
	}

	return nil
}

func (pipe *pipe) paint(ren *sdl.Renderer, tex *sdl.Texture) error {
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	rect := &sdl.Rect{W: pipe.w, H: pipe.h, X: pipe.x, Y: (winHeight - pipe.h)}
	flip := sdl.FLIP_NONE

	if pipe.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := ren.CopyEx(tex, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("Could not copy pipe texture: %v", err)
	}

	return nil
}

func (pipes *pipes) check(bird *bird) {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	for _, pipe := range pipes.ps {
		bird.check(pipe)
	}
}

func (pipes *pipes) destroy() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	pipes.texture.Destroy()
}
