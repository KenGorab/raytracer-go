package primitives

import "math"

type Sphere struct {
	Center Vector
	Radius float64
	Material
}

func (s *Sphere) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	oc := r.Origin.Subtract(s.Center)
	a := r.Direction.Dot(r.Direction)
	b := oc.Dot(r.Direction)
	c := oc.Dot(oc) - s.Radius * s.Radius
	discriminant := math.Pow(b, 2) - a * c

	rec := HitRecord{Material: s.Material}

	if discriminant > 0 {
		t := (-b - math.Sqrt(discriminant)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.Point = r.Point(t)
			rec.Normal = rec.Point.Subtract(s.Center).DivideScalar(s.Radius)

			return true, rec
		}

		t = (-b + math.Sqrt(discriminant)) / a
		if t < tMax && t > tMin {
			rec.T = t
			rec.Point = r.Point(t)
			rec.Normal = rec.Point.Subtract(s.Center).DivideScalar(s.Radius)

			return true, rec
		}
	}

	return false, rec
}
