package event_test

import (
	"fmt"
	"testing"

	"github.com/deminzhang/go-common/event"
)

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
		panic("testX")
	})
	event.LockReg() // 锁定注册，防止后续注册
	// Event3.Reg(func() { // 此处会panic,因为已经锁定注册
	// 	fmt.Println("Callback 3")
	// })

	// trigger.go
	Event1.Call("_EventTest2A", 100)
	Event2.Call("_EventTest2A")
	Event3.Call()

}
