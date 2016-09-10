package primitives

type HitRecord struct {
	T float64
	P, Normal Vector
}

type Hittable interface {
	Hit(r *Ray, tMin float64, tMax float64) (bool, HitRecord)
}
