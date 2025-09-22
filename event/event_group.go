package event

import "sync"

// 事件分组
type eventGroupList struct {
	sync.Mutex
	events map[int]map[int]*Event
}

var groupEvents eventGroupList

func init() {
	groupEvents.events = make(map[int]map[int]*Event)
}

func (this *eventGroupList) get(groupId, id int) *Event {
	this.Lock()
	defer this.Unlock()

	g := this.events[groupId]
	if g == nil {
		return nil
	}
	return g[id]
}

func (this *eventGroupList) add(groupId, id int, e *Event) {
	this.Lock()
	defer this.Unlock()

	g := this.events[groupId]
	if g == nil {
		g = make(map[int]*Event)
		this.events[groupId] = g
	}
	g[id] = e
}

// Listener
// 注册响应(事件名,回调) 以首次注册时函数类型为准 必须init时全注册好 动态注册可能造成call时的不确定性
func RegGroup(groupId, id int, cb interface{}) {
	if lockReg {
		panic("static event had lock reg")
	}
	e := groupEvents.get(groupId, id)
	if e == nil {
		e = NewEvent(id, cb)
		groupEvents.add(groupId, id, e)
	}
	e.Reg(cb)
}

// Dispatcher
// 触发事件_同步串行(事件名,参数集)
func CallGroup(groupId, id int, a ...interface{}) {
	e := groupEvents.get(groupId, id)
	if e == nil {
		return
	}
	e.Call(a...)
}

//初始化等一次性用的,用完可以清掉
//func CloseGroup(groupId int) {
//	defer groupEvents.Unlock()
//	groupEvents.Lock()
//	groupEvents.events[groupId] = make(map[int]*Event)
//}
//
//func CloseAllGroup() {
//	defer groupEvents.Unlock()
//	groupEvents.Lock()
//	groupEvents.events = make(map[int]map[int]*Event)
//}
