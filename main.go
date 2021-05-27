package main

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var ErrTerminated = errors.New("terminated")

type Control interface {
	control()
}

type Game struct {
	player        GameObject
	backGround    *ebiten.Image
	controlTarget Control
	objects       []GameObject
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrTerminated
	}
	g.controlTarget.control()

	for _, object := range g.objects {
		object.update()
	}

	for i, objA := range g.objects {
		if len(g.objects) == i+1 {
			break
		}
		d := objA.(*Actor)
		others := resolv.NewSpace()
		for _, other := range g.objects[i+1:] {
			c := other.(*Actor)
			others.Add(c.hitBox.area)
		}

		// X-axis
		res := others.Resolve(d.hitBox.area, int32(d.vec.x), 0)
		if res.Colliding() {
			d.vec.x = float64(res.ResolveX)
		}

		// Y-axis
		res = others.Resolve(d.hitBox.area, 0, int32(d.vec.y)+1)
		if res.Colliding() {
			d.onFloor = true
			d.vec.y = float64(res.ResolveY)
		} else {
			d.onFloor = false
		}
		objA.afterUpdate()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("(%+v)", g.player))
	screen.DrawImage(g.backGround, nil)

	for _, object := range g.objects {
		object.draw(screen)
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
	player := newPlayer(
		300, 200,
		i,
	)

	backGround := ebiten.NewImage(640, 480)
	backGround.Fill(color.RGBA{100, 100, 100, 255})
	groundImage := ebiten.NewImage(50, 50)
	groundImage.Fill(color.White)
	ground := newActor(100, 300, groundImage)

	floorImage := ebiten.NewImage(600, 20)
	floorImage.Fill(color.White)
	floor := newActor(0, 400, floorImage)
	objects := []GameObject{player, ground, floor}
	if err != nil {
		log.Fatal(err)
	}
	game = &Game{player: player, objects: objects, controlTarget: player, backGround: backGround}
	if err := ebiten.RunGame(game); err != nil {
		if err == ErrTerminated {
			return
		}
		log.Fatal(err)
	}
}
