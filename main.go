package main

import (
	"fmt"
	"os"
	"kenrg.co/rayz/primitives"
	"math"
	"math/rand"
)

const (
	nx = 800
	ny = 400
	numSamples = 10
	c = 255.99
)

var (
	white = primitives.Vector{1, 1, 1}
	blue = primitives.Vector{0.5, 0.7, 1}
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func color(r primitives.Ray, world primitives.Hittable, depth int) primitives.Vector {
	hit, record := world.Hit(r, 0.001, math.MaxFloat64)

	if hit {
		if depth < 50 {
			bounced, bouncedRay := record.Bounce(r, record)
			if bounced {
				newColor := color(bouncedRay, world, depth + 1)
				return record.Material.Color().Multiply(newColor)
			}
		}

		return primitives.Vector{}
	}

	return gradient(r)
}

func gradient(r primitives.Ray) primitives.Vector {
	v := r.Direction.Normalize()
	t := 0.5 * (v.Y + 1.0)
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func createFile() *os.File {
	f, err := os.Create("out.ppm")
	check(err, "Error opening file: %v\n")

	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writing to file: %v\n")
	return f
}

func writePixel(f *os.File, rgb primitives.Vector) {
	ir := int(c * math.Sqrt(rgb.X))
	ig := int(c * math.Sqrt(rgb.Y))
	ib := int(c * math.Sqrt(rgb.Z))

	_, err := fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
	check(err, "Error writing point to file: %v\n")
}

func sample(world *primitives.World, camera *primitives.Camera, i int, j int) primitives.Vector {
	rgb := primitives.Vector{}

	for s := 0; s < numSamples; s++ {
		u := (float64(i) + rand.Float64()) / float64(nx)
		v := (float64(j) + rand.Float64()) / float64(ny)

		r := camera.RayAt(u, v)
		rgb = rgb.Add(color(r, world, 0))
	}

	return rgb.DivideScalar(float64(numSamples))
}

func render(world *primitives.World, camera *primitives.Camera) {
	f := createFile()
	defer f.Close()

	ch := make(chan int)
	defer close(ch)

	go func() {
		for {
			if j, rendering := <-ch; rendering {
				percentComplete := float64(j) / float64(ny) * 100.0
				fmt.Printf("\rRendering: %.2f%%", percentComplete)
			}
		}
	}()

	row := 1
	for j := ny - 1; j >= 0; j-- {
		ch <- row
		row++

		for i := 0; i < nx; i++ {
			rgb := sample(world, camera, i, j)
			writePixel(f, rgb)
		}
	}
}

func main() {
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

	render(&world, &camera)
}
