package main

// Intersect encapsulates all information of an intersection between a ray and an object.
type Intersect struct {
	ray      *Ray
	geometry Geometry
	position *Vec3
	t        float64
}

// NewIntersect returns a new intersect.
func NewIntersect(ray *Ray, geometry Geometry, t float64) *Intersect {
	return &Intersect{
		ray:      ray,
		geometry: geometry,
		position: ray.At(t),
		t:        t,
	}
}

// Geometry returns the geometry on which the intersection lies.
func (i *Intersect) Geometry() Geometry {
	return i.geometry
}

// T returns the t value of the intersect.
func (i *Intersect) T() float64 {
	return i.t
}
