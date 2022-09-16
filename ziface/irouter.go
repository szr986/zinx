package ziface

// 路由的抽象接口
// 路由的数据都是IRequest请求的

type IRouter interface {
	// 在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	// 在处理conn业务的主方法hook
	Handle(request IRequest)
	// 在处理conn业务之后的子方法hook
	PostHandle(request IRequest)
}
