package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	vertArrays  map[uint32][2]uint32
	lastScrollX float64
	lastScrollY float64
	scale       float32 = 1.0

	camX float64 = -width / 2
	camY float64 = -height / 2

	pressPosX float64
	pressPosY float64

	currentGraphIdx uint32
)

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	// Сглаживание линий
	gl.Enable(gl.LINE_SMOOTH)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)

	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) [3]uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STREAM_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	var pack [3]uint32
	pack[0] = vao
	pack[1] = uint32(len(points) / 2)
	pack[2] = vbo

	return pack
}

func drawLine(x1, y1, x2, y2, width float32) {
	gl.LineWidth(width)

	gl.Begin(gl.LINES)
	{
		gl.Vertex2f(x1, y1)
		gl.Vertex2f(x2, y2)
	}
	gl.End()
}

func renderScene(vao [3]uint32) {
	gl.PointSize(2)
	gl.Color3f(1, 1, 0.)
	gl.BindVertexArray(vao[0])
	gl.DrawArrays(gl.POINTS, 0, int32(vao[1]))
}

func draw(window *glfw.Window, program uint32, vao [3]uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, width, height, 0, -1, 1)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Scalef(scale, scale, 0)

	renderScene(vao)

	gl.Flush()

	glfw.PollEvents()
	window.SwapBuffers()
}
