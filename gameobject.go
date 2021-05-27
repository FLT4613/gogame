package main

import (
	"image/color"

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
		image       *ebiten.Image
		hitBox      HitBox
		onFloor     bool
		isWalk      bool
		isJumpReady bool
		sm          *StateManager
	}
)

type HitBoxOption func(*HitBox)
type HitBox struct {
	area   *resolv.Rectangle
	offset Point
	image  *ebiten.Image
}

func setSize(size Point) HitBoxOption {
	return func(h *HitBox) {
		h.area.SetXY(int32(size.x), int32(size.y))
	}
}

func setOffset(offset Point) HitBoxOption {
	return func(h *HitBox) {
		h.offset.x = offset.x
		h.offset.y = offset.y
	}
}

func newHitBox(size Point, options ...HitBoxOption) HitBox {
	area := resolv.NewRectangle(0, 0, int32(size.x), int32(size.y))
	hitBox := HitBox{
		area:   area,
		offset: Point{0, 0},
		image:  ebiten.NewImage(int(size.x), int(size.y)),
	}
	hitBox.image.Fill(color.RGBA{255, 0, 0, 100})

	for _, option := range options {
		option(&hitBox)
	}

	return hitBox
}

type GameObject interface {
	update()
	afterUpdate()
	draw(*ebiten.Image)
}

type Option func(*Actor)

func newActor(x float64, y float64, image *ebiten.Image, options ...Option) *Actor {
	obj := Actor{}
	obj.pos.x = x
	obj.pos.y = y
	obj.image = image
	for _, option := range options {
		option(&obj)
	}
	obj.hitBox = newHitBox(Point{float64(obj.image.Bounds().Dx()), float64(obj.image.Bounds().Dy())})
	return &obj
}

func newPlayer(x float64, y float64, image *ebiten.Image, options ...Option) *Actor {
	obj := newActor(x, y, image, options...)
	obj.hitBox = newHitBox(Point{38, 40}, setOffset(Point{10, 10}))
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
	obj.hitBox.area.SetXY(int32(obj.pos.x)+int32(obj.hitBox.offset.x), int32(obj.pos.y)+int32(obj.hitBox.offset.y))
}

func (obj *Actor) afterUpdate() {
	obj.pos.Add(obj.vec)
}

func (obj *Actor) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(obj.pos.x, obj.pos.y)
	screen.DrawImage(obj.image, op)

	opImage := &ebiten.DrawImageOptions{}
	opImage.GeoM.Translate(obj.pos.x+obj.hitBox.offset.x, obj.pos.y+obj.hitBox.offset.y)
	screen.DrawImage(obj.hitBox.image, opImage)
}
