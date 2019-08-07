package render

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	gol "github.com/tomlockwood/gogol"
)

const (
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

// ActionFunction that's run on every tick
var ActionFunction Action
var actionKey glfw.Key
var reacted, closeScreen bool
var cells [][]*cell
var game *gol.Game

// options passed in to the renderer
var options *gol.Options

//// EXTERNALLY ACCESSIBLE options

// Reacted represents if any action has occurred that hasn't been responded to
func Reacted() bool {
	return reacted
}

// Acted lets you indicate you have reacted to a KeyPress
func Acted() {
	reacted = true
}

// ActionKey is the key currently pressed
func ActionKey() glfw.Key {
	return actionKey
}

// SetActionFunction for renderer
func SetActionFunction(a Action) {
	ActionFunction = a
}

// CloseScreen closes the screen
func CloseScreen() {
	closeScreen = true
}

// InitGame allows new game creation from an options object
func InitGame(o gol.Options) {
	options = &o
	g := gol.MakeGame(*options)
	game = &g
	cells = makeCells(*game)
}

// GetGame returns the current game
func GetGame() gol.Game {
	return *game
}

// SetGame sets the current game
func SetGame(g gol.Game) {
	game = &g
	cells = makeCells(*game)
}

//// INTERNAL OBJECTS

func onKey(w *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		actionKey = key
		reacted = false
	}
}

// Renderer is the rendering object
type Renderer struct {
	width, height, fps int
	window             *glfw.Window
	program            uint32
}

// Action is called on every Tick of a Game
type Action func()

// ActionDefault is how the game behaves by default
func ActionDefault() {
	if reacted {
	} else if actionKey == glfw.KeyEscape {
		closeScreen = true
		reacted = true
	} else if actionKey == glfw.KeySpace {
		g := gol.MakeGame(*options)
		game = &g
		cells = makeCells(*game)
		reacted = true
	} else if actionKey == glfw.KeyK {
		gol.Save(
			gol.SaveContent{Rules: game.Rules.Array, Grid: game.Field.Front},
			fmt.Sprintf("./%s.json", time.Now().Format(time.RFC3339)))
		reacted = true
	} else if actionKey == glfw.KeyR {
		game.Reset()
		reacted = true
	} else {
		reacted = true
	}
}

// Make a valid renderer
func Make(width int, height int, fps int) Renderer {
	runtime.LockOSThread()

	window := initGlfw(width, height)
	window.SetKeyCallback(onKey)

	program := initOpenGL()

	return Renderer{
		width,
		height,
		fps,
		window,
		program}
}

// Render game of life
func (r *Renderer) Render() {
	SetActionFunction(ActionDefault)
	r.RenderAction()
}

// RenderAction renders a game with a passed-in actionFunction
func (r *Renderer) RenderAction() {
	g := gol.MakeGame(*options)
	game = &g
	cells = makeCells(*game)

	for !r.window.ShouldClose() {
		t := time.Now()
		ActionFunction()
		if closeScreen {
			return
		}
		game.Tick()
		draw(*game, cells, r.window, r.program)
		deltat := time.Second / time.Duration(r.fps)
		time.Sleep(deltat - time.Since(t))
	}
}

// initGlfw initializes glfw and returns a Window to use
func initGlfw(width int, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "AUTOMATA", glfw.GetPrimaryMonitor(), nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

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
				color := g.Rules.Array[g.Field.Front[y][x]].Colour
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
