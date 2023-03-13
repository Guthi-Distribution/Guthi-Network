package main

import (
	"GuthiNetwork/platform"
	"GuthiNetwork/utility"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

var width int
var height int

var range_number int // 1 for 100 to 200 and false for 0 to 100
var count int

type MandelbrotParam struct {
	Index         int
	Row_completed int
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

func plot_mandelbrot() {
	net_platform := platform.GetPlatform()
	count++
	if count == 2 {
		im := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				c, err := net_platform.GetDataOfArray("mandelbrot", height*i+j)
				if err != nil {
					panic(err)
				}
				r := utility.Min(c.(Color).R*3, 255)
				g := utility.Min(c.(Color).R*5, 255)
				b := utility.Min(c.(Color).R*7, 255)
				im.Set(i, j, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}
			fmt.Printf("Index completed: %d\n", i)
		}

		output, _ := os.Create("img.png")
		png.Encode(output, im)
	}
}

func render_mandelbrot(param MandelbrotParam) {
	range_number := param.Index
	diff := height / 2
	min := 0 + range_number*diff
	max := diff + range_number*diff
	max_iter := 100
	radius := 4.0

	fmt.Println(min, max)
	start := Complex{-2, -2}
	end := Complex{1, 2}
	net_platform := platform.GetPlatform()
	platform.SetState("render_mandelbrot", param)

	start_index := param.Row_completed
	for x := start_index; x < width; x++ {
		real := start.real + (float64(x)/float64(width))*(end.real-start.real)
		for y := min; y < max; y++ {
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

		param.Row_completed = x
		platform.SetState("render_mandelbrot", param)
		fmt.Printf("Index completed: %d\n", x)
		time.Sleep(time.Millisecond * 50)
		platform.SendIndexedArray("mandelbrot", height*x, height, net_platform)
	}
	// platform.Send_array_to_nodes("mandelbrot", net_platform)
	fmt.Println("Completed")
}
