package znet

import (
	"fmt"
	"strconv"

	"example.com/m/utils"
	"example.com/m/ziface"
)

// 消息管理的具体实现

type MsgHandle struct {
	// 存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的工作数量
	WorkerPoolSize uint32
}

// 初始化，创建msgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), "is NOT FOUND!NEED Register!")
		return
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

// 启动一个Worker工作池(开启工作池的动作只能发生一次)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据WorkerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1.当前worker对应的channel消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2.启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "is start!...")
	for {
		select {
		case request := <-taskQueue:
			// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给taskqueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不通过的worker
	// 根据客户端简历的connid来分配

	// 轮询
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add connID = ", request.GetConnection().GetConnID,
		"request MSGID = ", request.GetMsgId(),
		"to workerID = ", workerID)

	// 将消息发送给对应的worker
	mh.TaskQueue[workerID] <- request
}
