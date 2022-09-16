package znet

import "example.com/m/ziface"

// 实现router时，先嵌入这个baserouter基类，然后根据需要对这个基类的方法进行重写就好了
type BaseRouter struct {
}

// 这里之所以BaseRouter的方法都为空
// 是因为有的Router不希望有pre，post这两个业务
// 所以Router全部继承BaseRouter的好处就是，不需要实现pre post
func (br *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (br *BaseRouter) Handle(request ziface.IRequest)     {}
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
