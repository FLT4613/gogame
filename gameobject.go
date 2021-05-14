package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject struct {
	x         float64
	y         float64
	MoveSpeed float64
	img       *ebiten.Image
}

func (obj *GameObject) moveLeft() {
	log.Print(obj.x)
	obj.x -= obj.MoveSpeed
}

func (obj *GameObject) moveRight() {
	obj.x += obj.MoveSpeed
}

func (obj *GameObject) update() {

	// if(obj.y < )	obj.y += 10

}
