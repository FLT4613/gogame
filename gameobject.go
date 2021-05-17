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
	GameObject struct {
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

type Option func(*GameObject)
type StateOption func(*State)

func onUpdate(body func()) StateOption {
	return func(state *State) {
		state.onUpdate = body
	}
}

func setImg(i *ebiten.Image) Option {
	return func(obj *GameObject) {
		obj.img = i
	}
}

func newGameObject(options ...Option) *GameObject {
	obj := GameObject{}
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

func (obj *GameObject) control() {
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

func (obj *GameObject) walk() {
	sign := 1.0
	if obj.direction == Left {
		sign = -1.0
	}

	obj.vec.x = obj.MoveSpeed * sign
}

func (obj *GameObject) update() {
	obj.sm.update()
	// obj.vec.Add(obj.acc)
	obj.pos.Add(obj.vec)
}
