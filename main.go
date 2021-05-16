package main

import (
	"errors"
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var ErrTerminated = errors.New("terminated")

type Control interface {
	control()
}

type Game struct {
	player        *GameObject
	controlTarget Control
	objects       []*GameObject
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrTerminated
	}
	g.controlTarget.control()

	for _, object := range g.objects {
		object.update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%+v)", g.player))
	op := &ebiten.DrawImageOptions{}
	for _, object := range g.objects {
		op.GeoM.Translate(object.pos.x, object.pos.y)
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
	player := newGameObject(
		setImg(i),
	)
	objects := []*GameObject{player}

	if err != nil {
		log.Fatal(err)
	}
	game = &Game{player: player, objects: objects, controlTarget: player}
	if err := ebiten.RunGame(game); err != nil {
		if err == ErrTerminated {
			return
		}
		log.Fatal(err)
	}
}
