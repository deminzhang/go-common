package event

import (
	"fmt"
	"reflect"
	"sync"
)

type ErrorMode int

const (
	PanicOnError ErrorMode = iota
	LogErrors
)

var (
	defaultErrorHandler = func(err any) {
		switch errorMode {
		case PanicOnError:
			panic(err)
		case LogErrors:
			fmt.Printf("Event callback panic: %v\n", err)
		}
	}

	errorMode    ErrorMode     = PanicOnError
	errorHandler func(err any) = defaultErrorHandler
)

type EventType[T any] struct {
	mu   sync.RWMutex
	cbs  []reflect.Value //调用时用,Reg时反射过的funcs
	Call T
}

func (e *EventType[T]) Reg(cb T) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cbs = append(e.cbs, reflect.ValueOf(cb))
}
func (e *EventType[T]) UnReg(cb T) {
	e.mu.Lock()
	defer e.mu.Unlock()

	targetPtr := reflect.ValueOf(cb).Pointer()
	for i, registered := range e.cbs {
		if registered.Pointer() == targetPtr {
			lastIdx := len(e.cbs) - 1
			if i < lastIdx {
				e.cbs[i] = e.cbs[lastIdx]
			}
			e.cbs = e.cbs[:lastIdx]
			return
		}
	}
}

func Event[T any]() *EventType[T] {
	cbType := reflect.TypeFor[T]()
	if cbType.Kind() != reflect.Func {
		panic("Event type parameter must be a function")
	}
	e := &EventType[T]{}
	fn := reflect.MakeFunc(cbType, func(args []reflect.Value) []reflect.Value {
		var ret []reflect.Value // 返回最后一个回调的返回值
		e.mu.RLock()
		defer e.mu.RUnlock()
		for _, cb := range e.cbs {
			func() {
				defer func() {
					if r := recover(); r != nil {
						errorHandler(r)
					}
				}()
				ret = cb.Call(args)
			}()
		}
		return ret
	})
	e.Call = fn.Interface().(T)
	return e
}

func SetErrorMode(mode ErrorMode) {
	errorMode = mode
}

func SetErrorHandler(handler func(err any)) {
	if handler == nil {
		errorHandler = defaultErrorHandler
		return
	}
	errorHandler = handler
}
