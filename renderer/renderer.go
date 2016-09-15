package renderer

import (
	"kenrg.co/rayz/primitives"
	"math/rand"
	"math"
	"fmt"
)

var (
	white = primitives.Vector{1, 1, 1}
	blue = primitives.Vector{0.5, 0.7, 1}
)

type Renderer struct {
	World  primitives.World
	Camera primitives.Camera
}

func gradient(r primitives.Ray) primitives.Vector {
	v := r.Direction.Normalize()
	t := 0.5 * (v.Y + 1.0)
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
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

func (r Renderer) sampleAt(x, y, width, height, numSamples int) primitives.Vector {
	rgb := primitives.Vector{}

	for s := 0; s < numSamples; s++ {
		u := (float64(x) + rand.Float64()) / float64(width)
		v := (float64(y) + rand.Float64()) / float64(height)

		ray := r.Camera.RayAt(u, v)
		rgb = rgb.Add(color(ray, &r.World, 0))
	}

	return rgb.DivideScalar(float64(numSamples))
}

func (r Renderer) Render(width, height, numSamples int) []primitives.Vector {
	chRow := make(chan int)
	defer close(chRow)

	go func() {
		for {
			if row, rendering := <-chRow; rendering {
				percentComplete := float64(row) / float64(height) * 100.0
				fmt.Printf("\rRendering: %.2f%%", percentComplete)
			}
		}
	}()

	pixels := make([]primitives.Vector, width * height)

	for j := height - 1; j >= 0; j-- {
		chRow <- height - j

		for i := 0; i < width; i++ {
			// We iterate in reverse when calculating pixels, so reverse the order when storing
			pixelIndex := width * height - (j * width + i) - 1
			pixels[pixelIndex] = r.sampleAt(i, j, width, height, numSamples)
		}
	}

	return pixels
}
