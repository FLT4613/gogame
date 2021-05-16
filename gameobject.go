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
		isWalk    bool
		direction Direction
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
	obj.acc.y = 0.02

	obj.sm = newStateManager()

	obj.sm.addState("Idle", onUpdate(func() { obj.vec.x = 0 }))
	obj.sm.addState("Move", onUpdate(func() {
		obj.walk()
	}))
	obj.sm.addTransition("Idle", "Move", func() bool {
		return obj.isWalk
	})
	obj.sm.addTransition("Move", "Idle", func() bool {
		return !obj.isWalk
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
	transitions  map[string][]func()
}

func newStateManager() *StateManager {
	var s StateManager
	s.states = map[string]State{}
	s.currentState = ""
	s.transitions = map[string][]func(){}
	return &s
}

func (sm *StateManager) update() {
	for _, transition := range sm.transitions[sm.currentState] {
		transition()
	}
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

func (sm *StateManager) addTransition(oldStateKey string, newStateKey string, condition func() bool) {
	_, ok := sm.transitions[oldStateKey]
	if !ok {
		sm.transitions[oldStateKey] = []func(){}
	}
	sm.transitions[oldStateKey] = append(sm.transitions[oldStateKey], func() {
		if condition() {
			sm.changeState(newStateKey)
		}
	})
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
