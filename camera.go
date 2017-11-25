package main

// Camera casts rays to the scene through the screen to retrieve the pixel densities.
type Camera struct {
	position Vec3 // position of the camera

	tangent  Vec3 // tangent vector (look direction)
	normal   Vec3 // normal vector (up direction)
	binormal Vec3 // binormal vector (right)

	aspectRatio float64 // aspect ratio of the camera lense
	normHeight  float64 // normalized height of the camera lense
}

func NewCamera(position, tangent, normal Vec3) *Camera {
	return &Camera{
		position:    position,
		tangent:     tangent,
		normal:      normal,
		binormal:    tangent.Cross(normal).Normalize(),
		aspectRatio: 1.0,
		normHeight:  1.0,
	}
}

func (c *Camera) Position() Vec3 {
	return c.position
}

func (c *Camera) BasisVecs() (Vec3, Vec3, Vec3) {
	return c.tangent, c.normal, c.binormal
}

func (c *Camera) LookAt() {

}

func (c *Camera) RayThrough(i, j int) Ray {
	return NewRay(NewVec3(0, 0, 0), NewVec3(0, 0, 0))
}
