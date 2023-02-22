package main

// There should be one univeral listening port

import (
	"GuthiNetwork/platform"
	"bufio"
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

	"github.com/mitchellh/go-ps"
)

/*
	TODO:
		- State Management
		- Creation of variable in single node
*/

type Config struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Complex struct {
	real, imag float64
}

func (c *Complex) absolute() float64 {
	return c.real*c.real + c.imag*c.imag
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
	for c.absolute() < radius && iter < max_iter {
		z = add(multiply(z, z), *c)
		iter += 1
	}
	*c = z

	return iter
}

func render_mandelbrot(range_number int) {
	diff := (256 / 2)
	min := 0 + range_number*diff
	max := diff + range_number*diff

	width := 256.0
	height := 256.0
	max_iter := 100
	radius := 4.0

	fmt.Println(min, max)
	start := Complex{-2.5, -2}
	end := Complex{1, 2}
	net_platform := platform.GetPlatform()
	for x := 0; x < 256; x++ {
		real := start.real + (float64(x)/width)*(end.real-start.real)
		for y := min; y < max; y++ {
			imag := start.imag + (float64(y)/height)*(end.imag-start.imag)
			z := Complex{real, imag}
			n_iter := does_diverge(&z, radius, max_iter)
			color_element := uint16(((n_iter - int(math.Log2(z.absolute()/radius))) / max_iter) * 255)

			color := Color{color_element}
			// _, err := net_platform.GetDataOfArray("mandelbrot", 256*x+y)
			err := net_platform.SetDataOfArray("mandelbrot", 256*x+y, color)
			if err != nil {
				log.Printf("Index: %d\n", 256*x+y)
				panic(err)
			}
			// data := _data.(Color)
		}
		fmt.Printf("Index completed: %d\n", x)
	}

	fmt.Println("Completed")
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

func isLaunchedByDebugger() bool {
	pid := os.Getppid()

	// We loop in case there were intermediary processes like the gopls language server.
	for pid != 0 {
		switch p, err := ps.FindProcess(pid); {
		case err != nil:
			return false
		case p.Executable() == "dlv":
			return true
		default:
			pid = p.PPid()
		}
	}
	return false
}

type Color struct {
	R uint16
	// G uint16
	// B uint16
}

func main() {
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
		curr_time := time.Now().UnixMilli()
		net_platform.CreateArray("mandelbrot", 256*256, c)
		fmt.Println(time.Now().UnixMilli() - curr_time)
		fmt.Println("Not Debugging process")

		if !isLaunchedByDebugger() {
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		} else {
			for len(net_platform.Connected_nodes) == 0 {

			}
		}
		net_platform.CallFunction(platform.GetFunctionName(render_mandelbrot), 0, "")
		net_platform.CallFunction(platform.GetFunctionName(render_mandelbrot), 1, net_platform.Connected_nodes[0].GetAddressString())
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	im := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{256, 256}})
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			c, err := net_platform.GetDataOfArray("mandelbrot", 256*i+j)
			if err != nil {

			}
			r := c.(Color).R
			g := c.(Color).R
			b := c.(Color).R
			im.Set(i, j, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
		fmt.Printf("Index completed: %d\n", i)
	}

	output, _ := os.Create("img.png")
	png.Encode(output, im)

	sg.Wait()
}
