package main

import (
	"fmt"
	"log"
)

type StateOption func(*State)

func onUpdate(body func()) StateOption {
	return func(state *State) {
		state.onUpdate = body
	}
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
