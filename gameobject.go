package main

import (
	"fmt"
	"log"

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
		pos       Point
		vec       Point
		acc       Point
		MoveSpeed float64
		img       *ebiten.Image
		sm        *StateManager
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

	obj.sm = newStateManager()
	obj.sm.addState("Idle", onUpdate(func() { obj.vec.x = 0 }))
	obj.sm.addState("MoveLeft", onUpdate(func() { obj.moveLeft() }))
	obj.sm.addState("MoveRight", onUpdate(func() { obj.moveRight() }))
	obj.sm.changeState("Idle")
	return &obj
}

func (obj *GameObject) control() {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		obj.sm.changeState("MoveLeft")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		obj.sm.changeState("MoveRight")
	}
	keys := inpututil.PressedKeys()
	if len(keys) == 0 {
		obj.sm.changeState("Idle")
	}
}

func (obj *GameObject) moveLeft() {
	obj.vec.x = -obj.MoveSpeed
}

func (obj *GameObject) moveRight() {
	obj.vec.x = obj.MoveSpeed
}

func (obj *GameObject) update() {
	obj.sm.update()
	if obj.pos.y < 400 {
		obj.vec.Add(obj.acc)
	} else {
		obj.vec.y = 0
		obj.acc.y = 0
	}
	obj.pos.Add(obj.vec)
}

type StateManager struct {
	states       map[string]State
	currentState string
}

func newStateManager() *StateManager {
	var s StateManager
	s.states = map[string]State{}
	s.currentState = ""
	return &s
}

func (sm *StateManager) update() {
	sm.states[sm.currentState].onUpdate()
}

func (sm *StateManager) addState(newStateKey string, options ...StateOption) {
	s := State{}
	s.onEnter = func() {}
	s.onUpdate = func() {}
	s.onExit = func() {}
	for _, option := range options {
		option(&s)
	}
	sm.states[newStateKey] = s
}

func (sm *StateManager) changeState(key string) {
	if sm.currentState == key {
		return
	}
	if sm.currentState != "" {
		current := sm.states[sm.currentState]
		current.onExit()
	}
	log.Print(fmt.Sprintf("%v -> %v", sm.currentState, key))
	sm.currentState = key
	sm.states[sm.currentState].onEnter()
}

type State struct {
	onEnter  func()
	onUpdate func()
	onExit   func()
}
