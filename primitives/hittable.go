package primitives

type HitRecord struct {
	T             float64
	Point, Normal Vector
	Material
}

type Hittable interface {
	Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord)
}
