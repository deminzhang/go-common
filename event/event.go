package event

import (
	"log"
	"reflect"
	"sync"
)

// Event 反射版静态事件
// 因用反射,限init时Def,Reg,不作Unregister
// 以常量为事件名
// 线程安全暂无
var lockReg bool //禁止后续注册,静态事件应在init期间注册好

type Event struct {
	id        int
	tp        reflect.Type
	caller    reflect.Value //总回调检查用
	listeners []listener
	argNum    int //参数数
}

type Context interface {
	Set(key string, val any)
	Get(key string) any
}

type EventContext struct {
	val map[string]any
}

func NewEventContextWithVal(k string, v any) *EventContext {
	return &EventContext{
		val: map[string]any{k: v},
	}
}

func (c EventContext) Set(k string, v any) {
	c.val[k] = v
}
func (c EventContext) Get(k string) any {
	v, b := c.val[k]
	if b {
		return v
	}
	return nil
}

type listener struct {
	caller reflect.Value //回调
}

func (e *Event) Reg(cb interface{}) {
	t := reflect.TypeOf(cb)
	if t.Kind() != reflect.Func {
		log.Fatalf("Event.Reg %d #2 must be a func*, got %v", e.id, cb)
	}
	if e.tp != t {
		log.Fatalln("Event.Reg type not equal:", e.id)
	}
	e.listeners = append(e.listeners, listener{
		caller: reflect.ValueOf(cb),
		//filter: newWhen,
	})
}

// 构建反射调用参数
func makeArgv(a ...interface{}) ([]reflect.Value, int) {
	argNum := len(a)
	in := make([]reflect.Value, argNum)
	for k, v := range a {
		in[k] = reflect.ValueOf(v)
	}
	return in, argNum
}

func (e *Event) Call(a ...interface{}) {
	list := e.listeners
	caller := e.caller
	in, argNum := makeArgv(a...)
	if argNum != e.argNum {
		log.Fatalln("Event.Call params number error:", e.id, argNum, caller.Type().NumIn())
		return
	}
	for _, l := range list {
		l.caller.Call(in)
	}
}

func (e *Event) GoCall(a ...interface{}) {
	list := e.listeners
	caller := e.caller
	in, argn := makeArgv(a...)
	if argn != e.argNum {
		log.Fatalln("Event.Call params number error:", e.id, argn, caller.Type().NumIn())
		return
	}
	for _, l := range list {
		go l.caller.Call(in)
	}
}

func (e *Event) GoCallWaitAll(a ...interface{}) {
	list := e.listeners
	caller := e.caller
	num := len(list)
	if num == 0 {
		return
	}
	in, argNum := makeArgv(a...)
	if argNum != e.argNum {
		log.Fatalln("Event.Call params number error:", e.id, argNum, caller.Type().NumIn())
		return
	}
	for _, l := range list {
		go func() {
			defer func() {
				num--
			}()
			l.caller.Call(in)
		}()
	}
	for num > 0 {
	}
}

func NewEvent(id int, foo interface{}) *Event {
	t := reflect.TypeOf(foo)
	if t.Kind() != reflect.Func {
		log.Fatalf("Event.Reg %d #2 must be a func", id)
	}
	call := reflect.ValueOf(foo)
	return &Event{
		id:     id,
		tp:     t,
		caller: call,
		argNum: call.Type().NumIn(),
	}
}

// 全局事件
type eventList struct {
	sync.Mutex
	events map[int]*Event
}

var globalEvents eventList

func init() {
	globalEvents.events = make(map[int]*Event)
}

func (l *eventList) get(id int) *Event {
	l.Lock()
	defer l.Unlock()
	return l.events[id]
}

func (l *eventList) add(id int, e *Event) {
	l.Lock()
	defer l.Unlock()
	l.events[id] = e
}

// Reg Listener
// 注册响应(事件名,回调) 以首次注册时函数类型为准 尽量init时全注册好 动态注册可能造成call时的不确定性
func Reg(id int, cb interface{}) {
	if lockReg {
		panic("static event had lock reg")
	}
	e := globalEvents.get(id)
	if e == nil {
		e = NewEvent(id, cb)
		globalEvents.add(id, e)
	}
	e.Reg(cb)
}

// Call Dispatcher
// 触发事件_同步串行(事件名,参数集)
func Call(id int, a ...interface{}) {
	e := globalEvents.get(id)
	if e == nil {
		return
	}
	e.Call(a...)
}

// GoCall 触发事件_异步并行(事件名,参数集)确定互不冲突可用
func GoCall(id int, a ...interface{}) {
	e := globalEvents.get(id)
	if e == nil {
		return
	}
	e.GoCall(a...)
}

// GoCallWaitAll 触发事件_同步并行(事件名,参数集)确定互不冲突可用
func GoCallWaitAll(id int, a ...interface{}) {
	e := globalEvents.get(id)
	if e == nil {
		return
	}
	e.GoCallWaitAll(a...)
}

// LockReg 锁定静态事件注册
func LockReg() {
	lockReg = true
}
