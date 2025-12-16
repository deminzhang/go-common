package event_test

import (
	"fmt"
	"testing"

	"github.com/deminzhang/go-common/event"
)

func TestEventExample(t *testing.T) {
	//Test测试示例
	fmt.Println(">>event0.sample")
	//定义
	const (
		_EventTest  = -1
		_EventTest2 = -2
		_Game1      = -100
		_Game2      = -101
	)
	//单响应
	event.Reg(_EventTest, func() {
		fmt.Println("_EventTest1")
	})
	event.Call(_EventTest)

	//多响应
	event.Reg(_EventTest2, func(s string, i int) {
		fmt.Println("A", s, i)
	})
	//
	//staticevent.When(nil, 100)
	event.Reg(_EventTest2, func(s string, i int) {
		fmt.Println("B", s, i+1)
	})
	event.Call(_EventTest2, "_EventTest2A", 100)
	//staticevent.Call(_EventTest2, "_EventTest2B", 200)

	//分组
	event.RegGroup(_Game1, _EventTest, func(s string, i int) {
		fmt.Println("C", s, i+1)
	})
	event.RegGroup(_Game1, _EventTest, func(s string, i int) {
		fmt.Println("D", s, i+1)
	})
	event.CallGroup(_Game1, _EventTest, "_EventTest3A", 100)

}

func BenchmarkEvent(b *testing.B) {
	if true { //0.0025~ ns/op
		event.Reg(-1, func(s string, i int) {
			fmt.Println(s, i+1)
		})
		event.Reg(-1, func(s string, i int) {
			fmt.Println(s, i+2)
		})
		event.Call(-1, "_EventTest2A", 100)
	} else { //0.0005~ ns/op
		var list = []func(s string, i int){
			func(s string, i int) {
				fmt.Println(s, i+1)
			},
			func(s string, i int) {
				fmt.Println(s, i+2)
			},
		}
		for _, f := range list {
			f("_EventTest2A", 100)
		}

	}
}

func TestEvent2(t *testing.T) {
	// define.go
	var (
		Event1 = event.Def[func(s string, i int)]()
		Event2 = event.Def[func(s string)]()
		Event3 = event.Def[func()]()
	)

	// listener.go
	// 示例1：两个参数的函数
	Event1.Reg(func(s string, i int) {
		fmt.Println("A", s, i)
	})
	Event1.Reg(func(s string, i int) {
		fmt.Println("B", s, i)
	})

	// 示例2：一个参数的函数
	Event2.Reg(func(s string) {
		fmt.Println("A", s)
	})
	Event2.Reg(func(s string) {
		fmt.Println("B", s)
	})

	// 示例3：没有参数的函数
	Event3.Reg(func() {
		fmt.Println("Callback 1")
	})
	Event3.Reg(func() {
		fmt.Println("Callback 2")
	})

	// trigger.go
	Event1.Call("_EventTest2A", 100)
	Event2.Call("_EventTest2A")
	Event3.Call()

}
