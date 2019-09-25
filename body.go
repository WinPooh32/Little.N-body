package main

import (
	"math"
)

// G is Gravitational constant
const G = 6.674E-11

type number float64

type vector struct {
	x, y number
}

type body struct {
	coord    vector
	velocity vector
	mass     number
}

func dist2d(a, b vector) number {
	t1 := (a.x - b.x)
	t2 := (a.y - b.y)
	s := float64(t1*t1 + t2*t2)
	return number(math.Sqrt(s))
}

func force2d(a, b body) vector {
	// d := dist2d(a.coord, b.coord)
	// g := (G * a.mass * b.mass)
	// if g < 0.0001 {
	// 	return 0
	// }
	// return g / d

	var f vector
	d := dist2d(a.coord, b.coord)
	// g := number(G * float64(a.mass*b.mass))
	if d > 1 {
		module := number(G*float64(a.mass*b.mass)) / d

		f.x = (module * (a.coord.x - b.coord.x)) / (d * d)
		f.y = (module * (a.coord.y - b.coord.y)) / (d * d)
	}

	// if g > math.SmallestNonzeroFloat32 {
	// 	dx := (a.coord.x - b.coord.x)
	// 	if math.Abs(float64(dx)) > math.SmallestNonzeroFloat32 {
	// 		f.x = g / dx * 10
	// 	}

	// 	dy := (a.coord.y - b.coord.y)
	// 	if math.Abs(float64(dy)) > math.SmallestNonzeroFloat32 {
	// 		f.y = g / dy * 10
	// 	}
	// }

	return f
}
