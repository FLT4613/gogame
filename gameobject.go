package main

import (
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Point struct {
	x float64
	y float64
}

func (p *Point) Add(p2 Point) *Point {
	p.x += p2.x
	p.y += p2.y
	return p
}

type (
	Actor struct {
		pos         Point
		vec         Point
		acc         Point
		direction   Direction
		MoveSpeed   float64
		img         *ebiten.Image
		hitBox      *resolv.Rectangle
		onFloor     bool
		isWalk      bool
		isJumpReady bool
		sm          *StateManager
	}
)

type GameObject interface {
	update()
	afterUpdate()
	draw(*ebiten.Image)
}

type Option func(*Actor)

// func setHitbox(hitBox image.Rectangle) Option {
// 	return func(obj *Actor) {
// 		obj.hitBox = hitBox
// 	}
// }

func newActor(x float64, y float64, image *ebiten.Image, options ...Option) *Actor {
	obj := Actor{}
	obj.pos.x = x
	obj.pos.y = y
	obj.img = image
	for _, option := range options {
		option(&obj)
	}
	obj.hitBox = resolv.NewRectangle(int32(x), int32(y), int32(image.Bounds().Dx()), int32(image.Bounds().Dy()))
	return &obj
}

func newPlayer(x float64, y float64, image *ebiten.Image, options ...Option) *Actor {
	obj := newActor(x, y, image, options...)
	obj.MoveSpeed = 3
	obj.sm = newStateManager()
	obj.sm.addState("Idle", onUpdate(func() { obj.vec.x = 0 }))
	obj.sm.addState("Move", onUpdate(func() {
		obj.walk()
	}))
	obj.sm.addState("JumpStart", onUpdate(func() {
		obj.onFloor = false
		obj.vec.y = -10
	}))
	obj.sm.addState("AirIdle", onUpdate(func() {
		obj.acc.y = 0.4
		obj.vec.y += obj.acc.y
	}))
	obj.sm.addState("Land", onUpdate(func() {
		obj.vec.y = 0
		obj.isJumpReady = false
	}))

	obj.sm.addTransition([]string{"Idle", "Land"}, "Move", func() bool {
		return obj.isWalk
	})
	obj.sm.addTransition([]string{"Move", "Land"}, "Idle", func() bool {
		return !obj.isWalk
	})
	obj.sm.addTransition([]string{"Idle", "Move", "Land"}, "JumpStart", func() bool {
		return obj.isJumpReady
	})
	obj.sm.addTransition([]string{"JumpStart"}, "AirIdle", func() bool {
		return true
	})
	obj.sm.addTransition([]string{"AirIdle"}, "Land", func() bool {
		return obj.onFloor
	})
	obj.sm.addTransition([]string{"Idle", "Move"}, "AirIdle", func() bool {
		return !obj.onFloor
	})

	obj.sm.changeState("Idle")
	return obj
}

type Direction int

const (
	Left Direction = iota
	Right
)

func (obj *Actor) control() {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		obj.direction = Left
		obj.isWalk = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		obj.direction = Right
		obj.isWalk = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		obj.isJumpReady = true
	}
	keys := inpututil.PressedKeys()
	if len(keys) == 0 {
		obj.isWalk = false
	}
}

func (obj *Actor) walk() {
	sign := 1.0
	if obj.direction == Left {
		sign = -1.0
	}

	obj.vec.x = obj.MoveSpeed * sign
}

func (obj *Actor) update() {
	if obj.sm != nil {
		obj.sm.update()
	}
	obj.hitBox.SetXY(int32(obj.pos.x), int32(obj.pos.y))
}

func (obj *Actor) afterUpdate() {
	obj.pos.Add(obj.vec)
}

func (obj *Actor) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(obj.pos.x, obj.pos.y)
	screen.DrawImage(obj.img, op)
}
