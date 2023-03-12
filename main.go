package main

// There should be one univeral listening port

import (
	"GuthiNetwork/platform"
	"GuthiNetwork/utility"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

var width int
var height int

type Config struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Complex struct {
	real, imag float64
}

type MandelbrotState struct {
	index       int
	count       int
	total_count int
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

var range_number int // 1 for 100 to 200 and false for 0 to 100

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

var count int

func plot_mandelbrot() {
	net_platform := platform.GetPlatform()
	count++
	log.Printf("Count %d\n", count)
	if count == 2 {
		im := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, width}})
		for i := 0; i < width; i++ {
			for j := 0; j < width; j++ {
				c, err := net_platform.GetDataOfArray("mandelbrot", width*i+j)
				if err != nil {

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

func render_mandelbrot(range_number int) {
	diff := width / 2
	min := 0 + range_number*diff
	max := diff + range_number*diff
	max_iter := 500
	radius := 4.0

	fmt.Println(min, max)
	start := Complex{-1.5, -2}
	end := Complex{1, 2}
	net_platform := platform.GetPlatform()
	curr_time := time.Now().UnixMilli()
	for x := 0; x < width; x++ {
		real := start.real + (float64(x)/float64(width))*(end.real-start.real)
		for y := min; y < max; y++ {
			imag := start.imag + (float64(y)/float64(height))*(end.imag-start.imag)
			z := Complex{real, imag}
			n_iter := does_diverge(&z, radius, max_iter)

			color_element := uint16(utility.Min((float64(n_iter)-math.Log2(z.absolute()/float64(radius)))/float64(max_iter)*255, 255.0))
			color := Color{color_element, utility.Min(255, color_element*2), utility.Min(255, color_element*3)}
			err := net_platform.SetDataOfArray("mandelbrot", width*x+y, color)

			if err != nil {
				log.Printf("Index: %d\n", width*x+y)
				panic(err)
			}
		}
		fmt.Printf("Index completed: %d\n", x)
		time.Sleep(time.Millisecond * 100)
	}

	platform.Send_array_to_nodes("mandelbrot", net_platform)
	fmt.Println("Completed")
	fmt.Printf("Total time taken: %d\n", time.Now().UnixMilli()-curr_time)
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		config.Name = ""
		config.Address = ""
		return config
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func sum(range_number int) {
	range_sum := 1000
	minimum := 1 + range_sum*range_number
	maximum := range_sum + range_sum*range_number
	net_platform := platform.GetPlatform()
	for i := minimum; i <= maximum; i++ {
		platform.Lock(net_platform)
		prev_sum := 0
		prev_sum_interface, _ := net_platform.GetData("total_sum")
		prev_sum = prev_sum_interface.(int)
		prev_sum += i
		net_platform.SetData("total_sum", prev_sum)
		platform.Unlock(net_platform)
		time.Sleep(time.Millisecond)
	}

	sum, _ := net_platform.GetData("total_sum")
	fmt.Println("Total sum: ", sum)
}

type Color struct {
	R uint16
	G uint16
	B uint16
}

func main() {
	count = 0
	width = 256
	height = 256
	port := flag.Int("port", 6969, "Port for the network") // send port using command line argument (-port 6969)
	sum_type := flag.Int("range", 0, "Type of range")

	flag.Parse()
	range_number = *sum_type
	fmt.Println(range_number, *sum_type)

	config := LoadConfiguration("configy.json")
	net_platform, err := platform.CreateNetworkPlatform(config.Name, config.Address, *port)
	if err != nil {
		panic(err)
	}
	if net_platform.Self_node.Socket.Port != 6969 {
		net_platform.ConnectToNode("127.0.0.1:6969") // one of the way to connect to a particular node, request all the nodes information it has
	}
	go platform.ListenForTCPConnection(net_platform)

	var sg sync.WaitGroup
	sg.Add(1)
	c := Color{}
	gob.Register(c)

	net_platform.RegisterFunction(render_mandelbrot)
	if *port == 6969 {
		net_platform.BindFunctionCompletionEventHandler("render_mandelbrot", plot_mandelbrot)
		curr_time := time.Now().UnixMilli()
		net_platform.CreateArray("mandelbrot", width*height, c)
		fmt.Println(time.Now().UnixMilli() - curr_time)
		fmt.Println("Not Debugging process")

		// if !isLaunchedByDebugger() {
		// 	reader := bufio.NewReader(os.Stdin)
		// 	reader.ReadString('\n')
		// } else {
		// 	for len(net_platform.Connected_nodes) == 0 {

		// 	}
		// }
		args := []interface{}{0, 1}
		// net_platform.CallFunction(platform.GetFunctionName(render_mandelbrot), 0, "")
		// net_platform.CallFunction(platform.GetFunctionName(render_mandelbrot), 1, net_platform.Connected_nodes[0].GetAddressString())
		net_platform.DispatchFunction("render_mandelbrot", args)
	}

	sg.Wait()
}
