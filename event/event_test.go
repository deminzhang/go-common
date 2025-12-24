package event_test

import (
	"fmt"
	"testing"

	"github.com/deminzhang/go-common/event"
)

// define*.go 全局静态事件
var (
	Event1 = event.Def[func(s string, i int)]()
	Event2 = event.Def[func(s string)]()
	//Event3 = event.Def[func(ctx context.Context)]() //复杂参数需引用复的的包,define.go专门放一个包
)

// 静态事件注册Reg=AddListener
func init() {
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
}

// logic*.go 动态事件
var (
	Event5 = event.Event[func(s string, i int)]()
)

// main thread
func TestEvent2(t *testing.T) {
	// Event2.Reg(func(s string) { // panic,因为已经锁定注册
	// 	fmt.Println("C", s)
	// })

	Event5.Reg(func(s string, i int) { //动态事件注册不锁定
		fmt.Println("E5", s, i)
	})

	// logic.go Call==Trigger==Dispatch
	Event1.Call("_EventTest2A", 100)
	Event2.Call("_EventTest2A")
	//Event2.Call("_EventTest2A", 2123) // 参数不匹配
	e5 := func(s string, i int) {
		fmt.Println("E6", s, i)
	}
	Event5.Reg(e5)
	Event5.Call("_EventTest2A", 100)
	Event5.UnReg(e5)
	Event5.Call("_EventTest2B", 200)
}
