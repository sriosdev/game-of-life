package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const screenW, screenH float64 = 800, 800
const nCellX, nCellY int = 50, 50

const cellWidth float64 = screenW / float64(nCellX)
const cellHeight float64 = screenH / float64(nCellY)

var bgColor color.Color = pixel.RGB(0.1, 0.1, 0.1)

var gameState [nCellX][nCellY]uint8
var gameStateFrame [nCellX][nCellY]uint8
var pause bool = true

func run() {
	win := initGame()

	for !win.Closed() {
		gameStateFrame = gameState

		win.Clear(bgColor)

		inputEvents(win)

		drawMesh(win)

		win.Update()
	}
}

func initGame() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "Game of life",
		Bounds: pixel.R(0, 0, screenW, screenH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func drawMesh(win *pixelgl.Window) {
	for y := 0; y < int(nCellX); y++ {
		for x := 0; x < int(nCellY); x++ {

			if !pause {
				nNeightbours := gameState[mod((x), nCellX)][mod((y+1), nCellY)] +
					gameState[mod((x+1), nCellX)][mod((y+1), nCellY)] +
					gameState[mod((x+1), nCellX)][mod((y), nCellY)] +
					gameState[mod((x+1), nCellX)][mod((y-1), nCellY)] +
					gameState[mod((x), nCellX)][mod((y-1), nCellY)] +
					gameState[mod((x-1), nCellX)][mod((y-1), nCellY)] +
					gameState[mod((x-1), nCellX)][mod((y), nCellY)] +
					gameState[mod((x-1), nCellX)][mod((y+1), nCellY)]

				if gameState[x][y] == 0 && nNeightbours == 3 {
					// Rule 1: A died cell witch has 3 alive neightbours, comes to life again
					gameStateFrame[x][y] = 1
				} else if gameState[x][y] == 1 && (nNeightbours < 2 || nNeightbours > 3) {
					// Rule 2: An alive cell with less than 2 or greater than 3 alive neightbour, dies.
					gameStateFrame[x][y] = 0
				}
			}

			meshPts := []pixel.Vec{
				pixel.V(float64(x)*cellWidth, float64(y)*cellHeight),
				pixel.V((float64(x)+1)*cellWidth, (float64(y)+1)*cellHeight),
			}

			meshCell := imdraw.New(nil)

			if gameStateFrame[x][y] == 0 {
				meshCell.Color = pixel.RGB(0.5, 0.5, 0.5)
				meshCell.Push(meshPts...)
				meshCell.Rectangle(1)
			} else {
				meshCell.Color = pixel.RGB(1, 1, 1)
				meshCell.Push(meshPts...)
				meshCell.Rectangle(0)
			}

			meshCell.Draw(win)
		}
	}

	gameState = gameStateFrame
}

func inputEvents(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeySpace) {
		pause = !pause
	}

	if win.Pressed(pixelgl.MouseButton1) {
		pos := win.MousePosition()
		setDeadOrAlive(pos, true)
	}

	if win.Pressed(pixelgl.MouseButton2) {
		pos := win.MousePosition()
		setDeadOrAlive(pos, false)
	}
}

func setDeadOrAlive(pos pixel.Vec, alive bool) {
	cellX, cellY := int(math.Floor(pos.X/cellWidth)), int(math.Floor(pos.Y/cellHeight))

	if alive {
		gameStateFrame[cellX][cellY] = 1
	} else {
		gameStateFrame[cellX][cellY] = 0
	}
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

func main() {
	pixelgl.Run(run)
}
