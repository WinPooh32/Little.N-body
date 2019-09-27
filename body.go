package main

import (
	"math"
)

// G is my custom Gravitational constant
// real G is 6.674E-11
const G = 6.674E-8

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
	var f vector
	d := dist2d(a.coord, b.coord)

	if d >= 2 {
		module := number(G*float64(a.mass*b.mass)) / (d)

		f.x = (module * (a.coord.x - b.coord.x)) / (d)
		f.y = (module * (a.coord.y - b.coord.y)) / (d)
	}

	return f
}
