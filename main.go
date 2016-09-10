package main

import (
	"fmt"
	"os"

	"kenrg.co/rayz/primitives"
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func main() {
	nx := 400
	ny := 200

	const color = 255.99

	f, err := os.Create("out.ppm")
	defer f.Close()
	check(err, "Error opening file: %v\n")

	_, err = fmt.Fprintf(f, "P3\n%d %d\n255\n", nx, ny)
	check(err, "Error writing to file: %v\n")

	lowerLeft := primitives.Vector{-2.0, -1.0, -1.0}
	horizontal := primitives.Vector{4.0, 0.0, 0.0}
	vertical := primitives.Vector{0.0, 2.0, 0.0}
	origin := primitives.Vector{0.0, 0.0, 0.0}

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u := float64(i) / float64(nx)
			v := float64(j) / float64(ny)

			position := horizontal.MultiplyScalar(u).Add(vertical.MultiplyScalar(v))

			// direction = lowerLeft + (u * horizontal) + (v * vertical)
			direction := lowerLeft.Add(position)

			rgb := primitives.Ray{origin, direction}.Color()

			ir := int(color * rgb.X)
			ig := int(color * rgb.Y)
			ib := int(color * rgb.Z)

			_, err = fmt.Fprintf(f, "%d %d %d\n", ir, ig, ib)
			check(err, "Error writing point to file: %v\n")
		}
	}
}
