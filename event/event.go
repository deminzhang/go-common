package Event

import (
	"fmt"
	"log"
	"reflect"
)

// 静态事件,限init时Reg
var events = make(map[string][]interface{})
var types = make(map[string]reflect.Type)

// func TestEvent() {
// 	//Test
// 	Reg("_EventTest", func() {
// 		fmt.Println("_EventTest")
// 	})
// 	Call("_EventTest")

// 	Reg("_EventTest2", func(s string, i int) {
// 		fmt.Println(s, i)
// 	})
// 	Reg("_EventTest2", func(s string, i int) {
// 		fmt.Println(s, i+1)
// 	})
// 	Call("_EventTest2", "_EventTest2", 123)

// }

// 注册响应 以首次注册时函数类型为准
func Reg(name string, foo any) {
	t0 := types[name]
	t := reflect.TypeOf(foo)
	s := fmt.Sprintf("%s", t)
	if len(s) < 5 || s[:5] != "func(" {
		log.Fatalf("Event.Reg %s #2 must be a func*, got %s", name, s)
	}
	if t0 == nil {
		types[name] = t
	} else {
		if t0 != t {
			log.Fatalln("Event.Reg type not equal:", name, t0, t)
		}
	}
	list := events[name]
	list = append(list, foo)
	events[name] = list
}

// 触发事件
func Call(name string, a ...any) {
	list := events[name]
	for _, foo := range list {
		f := reflect.ValueOf(foo)
		if len(a) != f.Type().NumIn() {
			log.Fatalln("Event.Call params number error:", name)
			return
		}
		in := make([]reflect.Value, len(a))
		for k, v := range a {
			in[k] = reflect.ValueOf(v)
		}
		f.Call(in)
	}
}
