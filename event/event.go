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
	errorMode ErrorMode = PanicOnError
)

// 动态事件
type Event[T any] struct {
	mu   sync.RWMutex //多线程时用
	cbs  []T
	Call T
}

func (e *Event[T]) Reg(cb T) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cbs = append(e.cbs, cb)
}
func (e *Event[T]) Unreg(cb T) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, registered := range e.cbs {
		if reflect.DeepEqual(
			reflect.ValueOf(registered).Pointer(),
			reflect.ValueOf(cb).Pointer(),
		) {
			e.cbs = append(e.cbs[:i], e.cbs[i+1:]...)
			return
		}
	}
}

func New[T any]() *Event[T] {
	cbType := reflect.TypeFor[T]()
	if cbType.Kind() != reflect.Func {
		panic("Event type parameter must be a function")
	}
	e := &Event[T]{
		cbs: make([]T, 0, 1),
	}
	fn := reflect.MakeFunc(cbType, func(args []reflect.Value) []reflect.Value {
		var ret []reflect.Value // 返回最后一个回调的返回值
		e.mu.RLock()
		defer e.mu.RUnlock()
		for _, cb := range e.cbs {
			func() {
				defer func() {
					switch errorMode {
					case PanicOnError:
						return
					}
					if r := recover(); r != nil {
						switch errorMode {
						case LogErrors:
							fmt.Printf("Event callback panic: %v\n", r)
						}
					}
				}()
				ret = reflect.ValueOf(cb).Call(args)
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
