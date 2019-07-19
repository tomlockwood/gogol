package render

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	gol "github.com/tomlockwood/gogol"
)

const (
	width  = 3840
	height = 2160

	vertexShaderSource = `
    #version 410
		in vec3 vp;
    void main() {
				gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
		#version 410
		uniform vec3 v_color;
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(v_color, 1.0);
    }
` + "\x00"
)

var (
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

var nextGame, saveGame, randomizeGame bool

func onKey(w *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	} else if key == glfw.KeySpace && action == glfw.Press {
		nextGame = true
	} else if key == glfw.KeyK && action == glfw.Press {
		saveGame = true
	} else if key == glfw.KeyR && action == glfw.Press {
		randomizeGame = true
	}
}

// Render game of life
func Render(o gol.GameOpts, fps int) {

	runtime.LockOSThread()

	window := initGlfw()
	window.SetKeyCallback(onKey)
	defer glfw.Terminate()

	program := initOpenGL()
	g := gol.MakeGame(o)
	cells := makeCells(g)

	for !window.ShouldClose() {
		t := time.Now()
		if nextGame {
			g = gol.MakeGame(o)
			cells = makeCells(g)
			nextGame = false
		} else if saveGame {
			gol.Save(
				gol.GameSave{Rules: g.Rules.Array, Grid: g.Grid.Array},
				fmt.Sprintf("./%s.json", time.Now().Format(time.RFC3339)))
			saveGame = false
		} else if randomizeGame {
			g.Reset()
			randomizeGame = false
		}
		g.Tick()
		draw(g, cells, window, program)
		deltat := time.Second / time.Duration(fps)
		time.Sleep(deltat - time.Since(t))
	}
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "AUTOMATA", glfw.GetPrimaryMonitor(), nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func draw(g gol.Game, cells [][]*cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	var wg sync.WaitGroup

	wg.Add(g.X * g.Y)

	for y := range cells {
		for x, c := range cells[y] {
			func(c cell) {
				defer wg.Done()
				color := g.Rules.Array[g.Grid.Array[y][x]].Colour
				gl.Uniform3f(0, color.R, color.G, color.B)
				c.draw()
			}(*c)
		}
	}

	wg.Wait()

	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

type cell struct {
	drawable uint32

	x int
	y int
}

func (c *cell) draw() {
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func makeCells(g gol.Game) [][]*cell {
	cells := make([][]*cell, g.Y)
	for y := range cells {
		cells[y] = make([]*cell, g.X)
		for x := range cells[y] {
			cells[y][x] = newCell(g, x, y)
		}
	}
	return cells
}

func newCell(g gol.Game, x, y int) *cell {
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(g.X)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(g.Y)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	return &cell{
		drawable: makeVao(points),

		x: x,
		y: y,
	}
}
