package primitives

type Ray struct {
	Origin, Direction Vector
}

func (r Ray) Point(t float64) Vector {
	return r.Origin.Add(r.Direction.MultiplyScalar(t))
}
