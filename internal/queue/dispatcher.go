package queue

import (
	"context"
	"fmt"
)

type Handler func(ctx context.Context, message Message) error // used to define handler's type below

type Dispatcher struct {
	handlers map[string]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]Handler),
	}
}

func (dispatcher *Dispatcher) Register(messageType string, handler Handler) {
	dispatcher.handlers[messageType] = handler
}

func (dispatcher *Dispatcher) Dispatch(ctx context.Context, message Message) error {
	handler, ok := dispatcher.handlers[message.Type]
	if !ok {
		return fmt.Errorf(
			"no handler registered for message type %q",
			message.Type,
		)
	}

	return handler(ctx, message)
}
