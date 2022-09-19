package znet

import (
	"fmt"
	"strconv"

	"example.com/m/ziface"
)

// 消息管理的具体实现

type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化，创建msgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), "is NOT FOUND!NEED Register!")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册了
		panic("repeat api,msgid = " + strconv.Itoa(int(msgID)))
	}
	// 2.添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api Msgid = ", msgID, " succ!")
}
