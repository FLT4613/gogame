package main

import (
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
	Ground struct {
		pos Point
		img *ebiten.Image
	}

	Character struct {
		pos         Point
		vec         Point
		acc         Point
		isWalk      bool
		isJumpReady bool
		direction   Direction
		MoveSpeed   float64
		img         *ebiten.Image
		sm          *StateManager
	}
)

type GameObject interface {
	update()
	draw(*ebiten.Image)
	setImg(*ebiten.Image)
	setPosition(float64, float64)
}

type Option func(GameObject)

func setImg(i *ebiten.Image) Option {
	return func(obj GameObject) {
		obj.setImg(i)
	}
}

func setPos(x float64, y float64) Option {
	return func(obj GameObject) {
		obj.setPosition(x, y)
	}
}

func newGround(options ...Option) *Ground {
	obj := Ground{}
	for _, option := range options {
		option(&obj)
	}
	return &obj
}

func (obj *Ground) update() {}

func (obj *Ground) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(obj.pos.x, obj.pos.y)
	screen.DrawImage(obj.img, op)
}

func (obj *Ground) setImg(i *ebiten.Image) {
	obj.img = i
}

func newCharacter(options ...Option) *Character {
	obj := Character{}
	for _, option := range options {
		option(&obj)
	}
	obj.MoveSpeed = 3
	// obj.acc.y = 0.005

	obj.sm = newStateManager()

	obj.sm.addState("Idle", onUpdate(func() { obj.vec.x = 0 }))
	obj.sm.addState("Move", onUpdate(func() {
		obj.walk()
	}))
	obj.sm.addState("JumpStart", onUpdate(func() {
		obj.vec.y = -9
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
		return obj.pos.y > 400
	})
	obj.sm.addTransition([]string{"Idle", "Move", "Land"}, "AirIdle", func() bool {
		return obj.pos.y < 400
	})

	obj.sm.changeState("Idle")
	return &obj
}

type Direction int

const (
	Left Direction = iota
	Right
)

func (obj *Character) control() {
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

func (obj *Character) walk() {
	sign := 1.0
	if obj.direction == Left {
		sign = -1.0
	}

	obj.vec.x = obj.MoveSpeed * sign
}

func (obj *Character) update() {
	obj.sm.update()
	// obj.vec.Add(obj.acc)
	obj.pos.Add(obj.vec)
}

func (obj *Character) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(obj.pos.x, obj.pos.y)
	screen.DrawImage(obj.img, op)
}

func (obj *Character) setImg(i *ebiten.Image) {
	obj.img = i
}

func (obj *Character) setPosition(x float64, y float64) {
	obj.pos.x = x
	obj.pos.y = y
}

func (obj *Ground) setPosition(x float64, y float64) {
	obj.pos.x = x
	obj.pos.y = y
}
