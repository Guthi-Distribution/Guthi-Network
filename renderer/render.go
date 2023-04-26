package renderer

import (
	"fmt"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func DrawCheckedTexture(array []byte, comp int, pitch int) {

	const gap = 25

	for h := 0; h < 400; h++ {
		for w := 0; w < 400; w++ {
			if (h/gap+w/gap)%2 == 0 {
				array[h*pitch+w*4+0] = 0xFF
				array[h*pitch+w*4+1] = 0xFF
				array[h*pitch+w*4+2] = 0xFF
				array[h*pitch+w*4+3] = 0xFF
			}
		}
	}

}

var g_renderer *sdl.Renderer
var g_window *sdl.Window
var g_texture *sdl.Texture

// update the line at that with arr of bytes
// the format of the supplied data must match
func UpdateTextureSurfaceOneLine(line int32, arr []byte) {
	if g_texture == nil {
		fmt.Println("Windows already destroyed")
		os.Exit(-5)
	}

	byte_array, pitch, err := g_texture.Lock(nil)
	if err != nil {
		fmt.Println("Failed to map the texture")
		os.Exit(-2)
	}
	fmt.Println("Pitched value : ", pitch)
	copy(byte_array[line*int32(pitch):], arr)
	g_texture.Unlock()
}

func PollEvents() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		// g_renderer.Present()
	}

	g_window.Destroy()
}

func UpdateTextureSurfaceOnePoint(w int32, h int32, r byte, g byte, b byte) {
	if g_texture == nil {
		fmt.Println("Windows already destroyed")
		os.Exit(-5)
	}

	byte_array, pitch, err := g_texture.Lock(nil)
	if err != nil {
		fmt.Println("Failed to map the texture")
		os.Exit(-2)
	}
	// copy(byte_array[h*int32(pitch)+w*3:], arr) // where 3 is the component currently in use
	offset := byte_array[h*int32(pitch)+w*3:]
	offset[0] = r
	offset[1] = g
	offset[2] = b
	g_texture.Unlock()
}

func PollSDLRenderer() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			g_renderer = nil
			sdl.Quit()
		}
	}
}

func PresentSurface() { // pass here the array to be updated
	if err := g_renderer.Clear(); err != nil {
		log.Fatal(err)
	}

	err := g_renderer.Copy(g_texture, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	g_renderer.Present()
	sdl.Delay(2)
}

func InitializeRenderer(width int32, height int32) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	var err error

	g_window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	g_renderer, err = sdl.CreateRenderer(g_window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		os.Exit(-1)
	}

	surface, err := g_window.GetSurface()
	if err != nil {
		fmt.Println("Failed to retrieve window surface")
		panic(err)
	}

	rect := sdl.Rect{0, 0, width, height}
	surface.FillRect(&rect, 0xFFFF0000)
	g_window.UpdateSurface()
	// Nothing here
	// For 800 * 600 draw all pixels
	g_texture, err = g_renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_STREAMING, width, height)

	if err != nil {
		log.Fatal(err)
	}

	go PollEvents()
	// PollSDLRenderer()
}

func StreamMandelbrot() {
	// net_platform := platform.GetPlatform()

}

func init() {
}
