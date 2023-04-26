package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/Guthi/guthi_network/platform"
	"github.com/Guthi/guthi_network/renderer"
	"github.com/Guthi/guthi_network/utility"
)

var range_number int // 1 for 100 to 200 and false for 0 to 100
var count int

const (
	block_size = 4
	width      = 256
	height     = 256
)

type MandelbrotParam struct {
	// Provide 4 x 4 square box to render
	X int
	Y int
}

type Complex struct {
	real, imag float64
}

func (c *Complex) absolute() float64 {
	return math.Sqrt(c.real*c.real + c.imag*c.imag)
}

func add(c1 Complex, c2 Complex) Complex {
	return Complex{c1.real + c2.real, c1.imag + c2.imag}
}

func multiply(c1 Complex, c2 Complex) Complex {
	return Complex{c1.real*c2.real - c1.imag*c2.imag, c1.real*c2.imag + c1.imag*c2.real}
}

func WaveColoring(c Complex, max_iter int, radius float64) float64 {
	z := Complex{0, 0}
	iterations := 0
	for i := 0; i < max_iter; i++ {
		z = add(multiply(z, z), c)
		iterations += 1
		if z.absolute() >= radius {
			break
		}
	}
	Amount := 0.2
	return 0.5 * math.Sin(Amount*float64(iterations))
}

func does_diverge(c *Complex, radius float64, max_iter int) int {
	iter := 0
	z := Complex{0, 0}
	for iter < max_iter {
		z = add(multiply(z, z), *c)
		iter += 1
		if z.absolute() > radius {
			break
		}
	}
	*c = z

	return iter
}

var present_mutex sync.Mutex

func plot_mandelbrot(func_name string, pram interface{}, return_value interface{}) {
	net_platform := platform.GetPlatform()

	params := pram.(MandelbrotParam)
	fmt.Println("Returned across the function calls : X -> ", params.X, " Y -> ", params.Y)

	present_mutex.Lock()
	defer present_mutex.Unlock()
	for i := params.X; i < params.X+block_size; i++ {
		for j := params.Y; j < params.Y+block_size; j++ {
			c, err := net_platform.GetDataOfArray("mandelbrot", height*i+j)
			if err != nil {
				panic(err)
			}
			r := byte(utility.Min(c.(Color).R*3, 255))
			g := byte(utility.Min(c.(Color).R*5, 255))
			b := byte(utility.Min(c.(Color).R*7, 255))

			renderer.UpdateTextureSurfaceOnePoint(int32(i), int32(j), r, g, b)
		}
	}

	renderer.PollSDLRenderer()
	renderer.PresentSurface()
	// fmt.Println("Callback function to plot mandelbrot was called")
}

func render_mandelbrot(args_supplied interface{}) {
	max_iter := 100
	radius := 4.0

	start := Complex{-2, -2}
	end := Complex{1, 2}
	net_platform := platform.GetPlatform()
	param := args_supplied.(MandelbrotParam) // the f**k is this syntax
	fmt.Println(param)

	platform.SetState("render_mandelbrot", param)

	for x := param.X; x < param.X+block_size; x++ {
		real := start.real + (float64(x)/float64(width))*(end.real-start.real)
		for y := param.Y; y < param.Y+block_size; y++ {
			imag := start.imag + (float64(y)/float64(height))*(end.imag-start.imag)
			z := Complex{real, imag}
			n_iter := does_diverge(&z, radius, max_iter)

			color_element := uint16(utility.Min((float64(n_iter)-math.Log2(z.absolute()/float64(radius)))/float64(max_iter)*255, 255.0))
			color := Color{color_element, utility.Min(255, color_element*2), utility.Min(255, color_element*3)}
			err := net_platform.SetDataOfArray("mandelbrot", height*x+y, color)

			if err != nil {
				log.Printf("Index: %d\n", width*x+y)
				panic(err)
			}
		}

		// platform.SetState("render_mandelbrot", param)
		// platform.SendIndexedArray("mandelbrot", height*x, block_size, net_platform)
	}
	// platform.Send_array_to_nodes("mandelbrot", net_platform)
	fmt.Println("Completed : ", param.X, " and ", param.Y)
}

func init() {
	present_mutex = sync.Mutex{}
}
