package primitives

import (
	"math"
	"math/rand"
)

type Vector struct {
	X, Y, Z float64
}

var UnitVector = Vector{1, 1, 1}

func VectorInUnitSphere() Vector {
	for {
		r := Vector{rand.Float64(), rand.Float64(), rand.Float64()}
		p := r.MultiplyScalar(2.0).Subtract(UnitVector)
		if p.SquaredLength() >= 1.0 {
			return p
		}
	}
}

func (v Vector) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

func (v Vector) SquaredLength() float64 {
	return math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2)
}

func (v Vector) Dot(other Vector) float64 {
	return v.X * other.X + v.Y * other.Y + v.Z * other.Z
}

func (v Vector) Normalize() Vector {
	return v.DivideScalar(v.Length())
}

func (v Vector) Add(o Vector) Vector {
	return Vector{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vector) Subtract(o Vector) Vector {
	return Vector{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vector) Multiply(o Vector) Vector {
	return Vector{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v Vector) Divide(o Vector) Vector {
	return Vector{v.X / o.X, v.Y / o.Y, v.Z / o.Z}
}

func (v Vector) AddScalar(t float64) Vector {
	return Vector{v.X + t, v.Y + t, v.Z + t}
}

func (v Vector) SubtractScalar(t float64) Vector {
	return Vector{v.X - t, v.Y - t, v.Z - t}
}

func (v Vector) MultiplyScalar(t float64) Vector {
	return Vector{v.X * t, v.Y * t, v.Z * t}
}

func (v Vector) DivideScalar(t float64) Vector {
	return Vector{v.X / t, v.Y / t, v.Z / t}
}
