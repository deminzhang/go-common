package event

import (
	"fmt"
	"reflect"
)

// 静态事件 禁止后续注册
type StaticEvent[T any] struct {
	cbs  []T
	Call T
}

var lockReg bool //禁止后续注册,静态事件应在init期间注册好

func (e *StaticEvent[T]) Reg(cb T) {
	if lockReg {
		panic("static event had lock reg. must reg in init")
	}
	e.cbs = append(e.cbs, cb)
}

// 全局锁定静态事件注册
func LockReg() {
	lockReg = true
}
func SDef[T any]() *StaticEvent[T] {
	cbType := reflect.TypeFor[T]()
	if cbType.Kind() != reflect.Func {
		panic("Event type parameter must be a function")
	}
	e := &StaticEvent[T]{
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
