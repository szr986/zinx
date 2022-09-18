package ziface

type IMessage interface {
	// 获取消息ID
	GetMsgId() uint32
	// 获取消息长度
	GetMsgLen() uint32
	// 获取消息的内容
	GetData() []byte

	// 设置消息ID
	SetMsgId(uint32)
	// 设置消息长度
	SetDataLen(uint32)
	// 设置消息内容
	SetData([]byte)
}
