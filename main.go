package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/kengorab/raytracer-go/primitives"
	"github.com/kengorab/raytracer-go/parser"
	"github.com/kengorab/raytracer-go/renderer"
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func createFile(width, height int) *os.File {
	f, err := os.Create("out.ppm")
	check(err, "Error opening file: %v\n")

	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", width, height)
	check(err, "Error writing to file: %v\n")
	return f
}

func writePixel(f *os.File, pixel primitives.Vector) {
	r, g, b := pixel.AsRGB()

	_, err := fmt.Fprintf(f, "%d %d %d\n", r, g, b)
	check(err, "Error writing point to file: %v\n")
}

func renderScene(camera primitives.Camera, world primitives.World, width, height, numSamples int) {
	f := createFile(width, height)
	defer f.Close()

	renderer := renderer.Renderer{World: world, Camera: camera}
	pixels := renderer.Render(width, height, numSamples)

	for _, pixel := range pixels {
		writePixel(f, pixel)
	}
}

func main() {
	arg := os.Args[1:][0]
	switch arg {
	case "-h":
		fmt.Println("Help")
		os.Exit(0)
	default:
		bytes, err := ioutil.ReadFile(arg)
		check(err, "Error opening scene file: %v\n")

		sceneConfig := parser.ReadConfigFromFile(bytes)

		camera := sceneConfig.CameraFromSceneConfig()
		world, err := sceneConfig.WorldFromSceneConfig()
		if err != nil {
			fmt.Printf("Error parsing input file:\n    %v\n", err)
			os.Exit(1)
		}

		renderScene(camera, world, sceneConfig.Width, sceneConfig.Height, sceneConfig.Samples)
	}
}
