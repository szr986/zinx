package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"example.com/m/utils"
	"example.com/m/ziface"
)

// 解决tcp粘包的封包拆包模块,tlv格式的封装
// 面向TCP连接中的数据流

// 封包拆包具体实现
type DataPack struct{}

// 拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	// Datalen uint32(4字节)+ID uint32（4字节）
	return 8
}

// 封包
//  datalen|msgid|data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将datalen写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgID写进databuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写进databuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	// 返回二进制
	return dataBuff.Bytes(), nil
}

// 拆包,将包的head信息读出来，再根据head里的len再进行一次读取
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioreader
	databuff := bytes.NewReader(binaryData)
	// 只解压head信息,得到dalalen和MsgID
	msg := &Message{}

	// 读datalen
	if err := binary.Read(databuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读MsgID
	if err := binary.Read(databuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断datalen是否已经超出了最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large package!")
	}

	return msg, nil
}
