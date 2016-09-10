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
	numSamples = 100
	c = 255.99
)

var (
	white = primitives.Vector{1, 1, 1}
	blue = primitives.Vector{0.5, 0.7, 1}

	camera = primitives.NewCamera()

	sphere = primitives.Sphere{Center: primitives.Vector{0, 0, -1}, Radius: 0.5}
	floor = primitives.Sphere{Center: primitives.Vector{0, -100.5, -1}, Radius: 100}

	world = primitives.World{[]primitives.Hittable{&sphere, &floor}}
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func color(r *primitives.Ray, h primitives.Hittable) primitives.Vector {
	hit, record := h.Hit(r, 0.0, math.MaxFloat64)

	if hit {
		return record.Normal.AddScalar(1.0).MultiplyScalar(0.5)
	}

	unitDirection := r.Direction.Normalize()
	return gradient(&unitDirection)
}

func gradient(v *primitives.Vector) primitives.Vector {
	t := 0.5 * (v.Y + 1.0)
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func main() {
	f, err := os.Create("out.ppm")
	defer f.Close()
	check(err, "Error opening file: %v\n")

	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writing to file: %v\n")

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			rgb := primitives.Vector{}

			for s := 0; s < numSamples; s++ {
				u := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)

				r := camera.RayAt(u, v)
				color := color(&r, &world)
				rgb = rgb.Add(color)
			}

			rgb = rgb.DivideScalar(float64(numSamples))

			ir := int(c * rgb.X)
			ig := int(c * rgb.Y)
			ib := int(c * rgb.Z)

			_, err = fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
			check(err, "Error writing point to file: %v\n")
		}
	}
}
