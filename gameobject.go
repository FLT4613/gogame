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
	obj.acc.y = 0.05

	obj.sm = newStateManager()

	obj.sm.addState("Idle", onUpdate(func() { obj.vec.x = 0 }))
	obj.sm.addState("Move", onUpdate(func() {
		obj.walk()
	}))
	obj.sm.addState("JumpStart", onUpdate(func() {
		obj.vec.y = -20
	}))
	obj.sm.addState("AirIdle", onUpdate(func() {
		if obj.acc.y == 0 {
			obj.acc.y = 1
		}

		obj.vec.y += obj.acc.y
		log.Print(obj.acc.y)
	}))
	obj.sm.addState("Land", onUpdate(func() {
		obj.acc.y = 0
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
	// if obj.pos.y < 400 {
	obj.vec.Add(obj.acc)
	// } else {
	// 	obj.vec.y = 0
	// 	obj.acc.y = 0
	// }
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

func (sm *StateManager) addTransition(sources []string, destination string, condition func() bool) {
	for _, source := range sources {
		_, ok := sm.transitions[source]
		if !ok {
			sm.transitions[source] = []func(){}
		}
		sm.transitions[source] = append(sm.transitions[source], func() {
			if condition() {
				sm.changeState(destination)
			}
		})
	}
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
