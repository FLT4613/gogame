package main

import (
	"errors"
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var Terminated = errors.New("terminated")

type Game struct {
	player  *GameObject
	objects []*GameObject
	keys    []ebiten.Key
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return Terminated
	}
	g.keys = inpututil.PressedKeys()
	for _, key := range g.keys {
		switch key {
		// case ebiten.KeyW:
		// g.player.moveLeft()
		case ebiten.KeyA:
			g.player.moveLeft()
		// case ebiten.KeyS:
		// y += 1
		case ebiten.KeyD:
			g.player.moveRight()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%+v)", g.player))
	op := &ebiten.DrawImageOptions{}
	for _, object := range g.objects {
		log.Print(object.x)
		op.GeoM.Translate(object.x, object.y)
		screen.DrawImage(object.img, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

var game *Game

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	i, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	player := GameObject{0, 0, 10, i}
	objects := []*GameObject{&player}

	if err != nil {
		log.Fatal(err)
	}
	game = &Game{player: &player, objects: objects}
	if err := ebiten.RunGame(game); err != nil {
		if err == Terminated {
			// Regular termination
			return
		}
		log.Fatal(err)
	}
}
