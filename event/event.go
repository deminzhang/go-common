package event

import (
	"fmt"
	"reflect"
)

type ErrorMode int

const (
	PanicOnError ErrorMode = iota
	LogErrors
)

var (
	errorMode ErrorMode = PanicOnError
	lockReg   bool      //禁止后续注册,静态事件应在init期间注册好
)

type Event[T any] struct {
	//mu    sync.RWMutex 如果Reg有多线程,则需要加锁,如果Reg在init,则不需要加锁
	cbs  []T
	Call T
}

func (e *Event[T]) Reg(cb T) {
	if lockReg {
		panic("static event had lock reg. must reg in init")
	}
	e.cbs = append(e.cbs, cb)
}

func LockReg() {
	lockReg = true
}
func SetErrorMode(mode ErrorMode) {
	errorMode = mode
}
func Def[T any]() *Event[T] {
	cbType := reflect.TypeFor[T]()
	if cbType.Kind() != reflect.Func {
		panic("Event type parameter must be a function")
	}
	e := &Event[T]{
		cbs: make([]T, 0, 1),
	}
	fn := reflect.MakeFunc(cbType, func(args []reflect.Value) []reflect.Value {
		var ret []reflect.Value // 返回最后一个回调的返回值
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
