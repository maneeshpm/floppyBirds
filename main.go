package main

import (
	"fmt"
	"os"
	"time"

	img "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

const (
	winWidth  = 1600
	winHeight = 900

	fontPath = "res/fonts/flappy_font.ttf"
	fontSize = 20

	bgPath = "res/imgs/background.png"

	flappyName = "flappy shukle"
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

	if err := drawTitle(ren); err != nil {
		return fmt.Errorf("Could not draw title: %v", err)
	}

	time.Sleep(5 * time.Second)

	if err := drawBackground(ren); err != nil {
		return fmt.Errorf("Could not draw background: %v", err)
	}

	time.Sleep(5 * time.Second)
	return nil
}

// Draw a test title
func drawTitle(ren *sdl.Renderer) error {
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
	surface, err := font.RenderUTF8Solid(flappyName, c)
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

	if err := ren.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	ren.Present()
	return nil
}

func drawBackground(ren *sdl.Renderer) error {
	ren.Clear()

	imageTexture, err := img.LoadTexture(ren, bgPath)
	if err != nil {
		return fmt.Errorf("Could not create background texture: %v", err)
	}
	defer imageTexture.Destroy()

	if err := ren.Copy(imageTexture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}
	ren.Present()
	return nil
}