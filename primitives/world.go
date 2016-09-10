package primitives

type World struct {
	elements []Hittable
}

func (w *World) Add(h Hittable) {
	w.elements = append(w.elements, h)
}

func (w *World) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	hitAnything := false
	closest := tMax
	record := HitRecord{}

	for _, element := range w.elements {
		hit, rec := element.Hit(r, tMin, closest)

		if hit {
			hitAnything = true
			closest = rec.T
			record = rec
		}
	}

	return hitAnything, record
}
