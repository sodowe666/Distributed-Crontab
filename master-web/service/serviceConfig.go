package service

//当有新的childService加入时，往下面两个配置文件添加新的元素

//childService容器，用于配置子service的初始化函数，和下面的结构体元素的类型对应，字符串开头必须小写
var childServiceContainer = map[string]func() ServiceInterface{
	"orderService": NewOrderService,
	"userService": NewUserService,
}

//配置子service，添加子service的结构体元素时，和上面的容器里的初始化函数配置对应
type Service struct {
	OrderService *OrderService
	UserService  *UserService
}
