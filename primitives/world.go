package primitives

type World struct {
	Elements []Hittable
}

func (w *World) Hit(r *Ray, tMin float64, tMax float64) (bool, HitRecord) {
	hitAnything := false
	closest := tMax
	record := HitRecord{}

	for _, element := range w.Elements {
		hit, rec := element.Hit(r, tMin, closest)

		if hit {
			hitAnything = true
			closest = rec.T
			record = rec
		}
	}

	return hitAnything, record
}
