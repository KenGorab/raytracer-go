package main

import (
	"fmt"
	"os"
	"kenrg.co/rayz/primitives"
	"kenrg.co/rayz/renderer"
)

const (
	nx = 800
	ny = 400
	numSamples = 10
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func createFile() *os.File {
	f, err := os.Create("out.ppm")
	check(err, "Error opening file: %v\n")

	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writing to file: %v\n")
	return f
}

func writePixel(f *os.File, pixel primitives.Vector) {
	r, g, b := pixel.AsRGB()

	_, err := fmt.Fprintf(f, "%d %d %d\n", r, g, b)
	check(err, "Error writing point to file: %v\n")
}

func main() {
	f := createFile()
	defer f.Close()

	camera := primitives.NewCamera(primitives.Vector{0, 1, 2}, primitives.Vector{0, 0, 0}, 75.0, float64(nx) / float64(ny), 0.01)

	floor := primitives.Sphere{Center: primitives.Vector{0, -250.5, 0}, Radius: 250, Material: primitives.Lambertian{primitives.Vector{0.3, 0.3, 0.3}}}
	left := primitives.Sphere{Center: primitives.Vector{-1, 0, 0}, Radius: 0.5, Material: primitives.Metal{primitives.Vector{0.8, 0.8, 0.8}, 0.0}}
	right := primitives.Sphere{Center: primitives.Vector{1, 0, 0}, Radius: 0.5, Material: primitives.Metal{primitives.Vector{0.8, 0.6, 0.2}, 0.0}}
	center := primitives.Sphere{Center: primitives.Vector{0, 0, -1}, Radius: 0.5, Material: primitives.Metal{primitives.Vector{0.1, 0.1, 0.1}, 0.0}}

	world := primitives.World{}
	world.Add(&center)
	world.Add(&floor)
	world.Add(&left)
	world.Add(&right)

	renderer := renderer.Renderer{World: world, Camera: camera}

	pixels := renderer.Render(nx, ny, numSamples)

	for _, pixel := range pixels {
		writePixel(f, pixel)
	}
}
