package ziface

// iRequest接口
// 把客户端请求的连接信息 和 数据 包装在 一个request中

type IRequest interface {
	// 得到当前的连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte

	// 得到请求消息的ID
	GetMsgId() uint32
}
