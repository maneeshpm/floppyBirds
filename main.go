package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

const (
	winWidth  = 1600
	winHeight = 900

	fontPath = "res/fonts/flappy_font.ttf"
	fontSize = 20

	bgPath = "res/imgs/background.png"

	flappyName = "flappy shawtty"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not run initializer: %v", err)
		os.Exit(2)
	}
	fmt.Println("hello")
}

// Start a test gui
func run() error {
	// Initialize sdl
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("Could not initialize sdl: %v", err)
	}
	defer sdl.Quit() // Quit at end

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Could not initialize ttf: %v", err)
	}
	defer ttf.Quit()

	// Initialize a windor and renderer
	win, ren, err := sdl.CreateWindowAndRenderer(winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Could not initialize renderer and window: %v", err)
	}
	defer win.Destroy()

	if err := drawText(ren, flappyName); err != nil {
		return fmt.Errorf("Could not draw title: %v", err)
	}

	time.Sleep(1 * time.Second)

	s, err := newScene(ren)
	if err != nil {
		return fmt.Errorf("Could not create scene: %v", err)
	}
	defer s.Destroy()

	// quit := make(chan struct{})

	events := make(chan sdl.Event)
	errc := s.run(ren, events)

	for {
		select {
		case err := <-errc:
			return err
		case events <- sdl.WaitEvent():
		}
	}

}

// Write a spanning string on the window
func drawText(ren *sdl.Renderer, text string) error {
	ren.Clear()
	font, err := ttf.OpenFont(fontPath, fontSize)

	if err != nil {
		return fmt.Errorf("Could not open font: %v", err)
	}
	defer font.Close()

	c := sdl.Color{
		R: 0,
		G: 200,
		B: 225,
		A: 0,
	}

	fmt.Println("works 1")
	surface, err := font.RenderUTF8Solid(text, c)
	if err != nil {
		return fmt.Errorf("Could not render text: %v", err)
	}
	defer surface.Free()
	fmt.Println("works 2")

	t, err := ren.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer t.Destroy()

	rect := &sdl.Rect{W: 1200, H: 700, X: 200, Y: 100}

	if err := ren.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	ren.Present()
	return nil
}
