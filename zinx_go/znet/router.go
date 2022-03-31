package znet

import "zinx-learn/ziface"

//实现router时，先嵌入这个BaseRouter基类，然后根据这个基类的方法进行重写就好
type BaseRouter struct{}

//这里之所以BaseRouter的方法都为空
//是因为有的Router不希望有PreHandle、PostHandle这两个业务
//在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

//在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
