package event

import (
	"errors"
	"log"
)

const (
	EventLinkVisited = "link.visited"
)

type Event struct {
	Type string
	Data any
}

type EventBus struct {
	bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{bus: make(chan Event, 100)}
}

func (e *EventBus) Publish(event Event) error {
	select {
	case e.bus <- event:
		return nil
	default:
		log.Println("[EventBus] Канал переполнен, событие потеряно:", event)
		return errors.New("канал переполнен")
	}
}

func (e *EventBus) Subscribe() <-chan Event {
	newChan := make(chan Event, 10)
	go func() {
		for msg := range e.bus {
			newChan <- msg
		}
	}()
	return newChan
}
