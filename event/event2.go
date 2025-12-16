package event

import "reflect"

type SEvent[T any] struct {
	cbs  []T
	Call T
}

func Def[T any]() *SEvent[T] {
	cbType := reflect.TypeFor[T]()
	if cbType.Kind() != reflect.Func {
		panic("Event type parameter must be a function")
	}
	e := &SEvent[T]{
		cbs: make([]T, 0, 1),
	}
	fn := reflect.MakeFunc(cbType, func(args []reflect.Value) []reflect.Value {
		for _, cb := range e.cbs {
			reflect.ValueOf(cb).Call(args)
		}
		return nil
	})
	e.Call = fn.Interface().(T)
	return e
}

func (e *SEvent[T]) Reg(cb T) {
	e.cbs = append(e.cbs, cb)
}
