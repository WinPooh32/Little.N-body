package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
	"runtime"

	"math/rand"
	"sync"
	"time"
)

// world's bodies total count
const bcount = 10000
const (
	width  = 1270
	height = 720
)

type testarray [10]int
type world [bcount]body
type pixels [bcount * 2]float32

func randBetween(min, max int64) number {
	return number(min + (rand.Int63() % int64(max-min+1)))
}

func sumForce(b body, w world) vector {
	var sum vector

	for i := range w {
		f := force2d(b, w[i])

		sum.x -= f.x
		sum.y -= f.y
	}

	return sum
}

func calcGalaxyMass(w *world) number {
	var mass number
	for _, b := range w {
		mass += b.mass
	}

	return mass
}

func makeGenesisState() world {
	var w world

	const dst = 100
	const vel = 0
	const vel_b = 1500
	const mmin = 1
	const mmax = 10
	const rad = 10

	for i := 0; i < len(w); i++ {
		w[i].coord.x = width/2 - rad + randBetween(-dst, +dst)*number(rand.Float64())
		w[i].coord.y = height/2 + rad + randBetween(-dst, +dst)*number(rand.Float64())

		w[i].mass = randBetween(mmin, mmax)

		w[i].velocity.x = randBetween(-vel, vel) + vel
		w[i].velocity.y = randBetween(-vel, vel)
	}

	gm := calcGalaxyMass(&w) * 1000

	b := &w[len(w)-1]
	b.coord.x = width/2 - 200
	b.coord.y = height/2 - 100
	b.mass = gm
	b.velocity.x = vel_b
	b.velocity.y = 0

	c := &w[len(w)-2]
	c.coord.x = width/2 + 50
	c.coord.y = height/2 + 200
	c.mass = gm
	c.velocity.x = -vel_b
	// c.velocity.y = -vel_b

	return w
}

func worldToPixels(w *world, px *pixels) {
	for i := range w {
		j := i * 2

		px[j] = float32(w[i].coord.x)
		px[j+1] = float32(w[i].coord.y)
	}

	// fmt.Println(px)
}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	// seed random by current time
	rand.Seed(time.Now().Unix())
	world := makeGenesisState()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(width, height, "n-body", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	program := initOpenGL()

	var pixelBuf pixels
	vao := makeVao(pixelBuf[:])

	const dt = 1E-4
	// t := number(0.0)

	fmt.Println("Begin")

	threads := runtime.NumCPU()
	taskBatch := int(math.Floor(float64(len(world) / threads)))

	var wg sync.WaitGroup
	worker := func(begin, end int) {
		for i := begin; i < end; i++ {
			b := &world[i]

			// velocity = velocity + ( force / mass ) * dt;
			sf := sumForce(*b, world)
			md := b.mass * dt

			if sf.x != 0 {
				fm := sf.x / md
				b.velocity.x += fm
			}
			if sf.y != 0 {
				fm := sf.y / md
				b.velocity.y += fm
			}
		}

		wg.Done()
	}

	integrator := func(begin, end int) {
		for i := begin; i < end; i++ {
			b := &world[i]

			// position = position + velocity * dt;
			b.coord.x += b.velocity.x * dt
			b.coord.y += b.velocity.y * dt
		}
		wg.Done()
	}

	for {
		// var last = time.Now()

		// calc forces
		{
			wg.Add(threads)
			for group := 0; group < threads; group++ {
				b := group * taskBatch
				e := b + taskBatch
				go worker(b, e)
			}
			wg.Wait() // sync threads
		}

		// integrate coords
		{
			wg.Add(threads)
			for group := 0; group < threads; group++ {
				b := group * taskBatch
				e := b + taskBatch
				go integrator(b, e)
			}
			wg.Wait() // sync threads
		}

		// render
		{
			worldToPixels(&world, &pixelBuf)
			gl.BindBuffer(gl.ARRAY_BUFFER, vao[2])
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(pixelBuf), gl.Ptr(&pixelBuf[0]))
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)

			draw(window, program, vao)
		}

		// fmt.Printf("microseconds per particle: %d\n", time.Now().Sub(last).Nanoseconds()/bcount/1000)
	}

	fmt.Println("End")
	select {}
}
