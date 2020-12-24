package utils

import "strings"

// EventHandler funtion
type EventHandler func(interface{})

// Event represent a event
type Event struct {
	Name     string
	Handlers map[uint]EventHandler
	_counter uint
}

// AddListener adds a listener to the event
func (e *Event) AddListener(handler EventHandler) uint {
	if e.Handlers == nil {
		e.Handlers = make(map[uint]EventHandler)
	}
	e.Handlers[e._counter] = handler
	e._counter++
	return e._counter - 1
}

// RemoveListener removes a listener from the event
func (e *Event) RemoveListener(handlerID uint) {
	if _, found := e.Handlers[handlerID]; found {
		delete(e.Handlers, handlerID)
	}
}

// Emit emits the event to all listeners
func (e *Event) Emit(payload interface{}) {
	for _, h := range e.Handlers {
		go h(payload)
	}
}

// Eventer contains the event container methods
type Eventer interface {
	AddEventListener(string, EventHandler) uint
	RemoveEventListener(string, uint)
	Emit(string, interface{})
}

// EventContainer basic struct
type EventContainer struct {
	events map[string]*Event
}

func (ec *EventContainer) verifyInitialEvents() {
	if ec.events == nil {
		ec.events = make(map[string]*Event)
	}
}

// AddEventListener append a handler for specific event
func (ec *EventContainer) AddEventListener(eventName string, handler EventHandler) uint {
	ec.verifyInitialEvents()
	eventName = strings.ToLower(eventName)
	if e, found := ec.events[eventName]; found {
		return e.AddListener(handler)
	}
	ec.events[eventName] = &Event{Name: eventName}
	return ec.events[eventName].AddListener(handler)
}

// RemoveEventListener remove a handler for specific event
func (ec *EventContainer) RemoveEventListener(eventName string, handlerID uint) {
	ec.verifyInitialEvents()
	eventName = strings.ToLower(eventName)
	if e, found := ec.events[eventName]; found {
		e.RemoveListener(handlerID)
	}
}

// Emit an event with their respective payload
func (ec *EventContainer) Emit(eventName string, payload interface{}) {
	ec.verifyInitialEvents()
	eventName = strings.ToLower(eventName)
	if e, found := ec.events[eventName]; found {
		e.Emit(payload)
	}
}
