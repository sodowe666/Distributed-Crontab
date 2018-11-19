package service

import (
	"reflect"
	"distributedcrontab/common/utils"
	"errors"
)

type ServiceInterface interface {
}

var childServiceContainer = map[string]func() ServiceInterface{
	"orderService": NewOrderService,
	"userService": NewUserService,
}

type Service struct {
	OrderService *OrderService
	UserService  *UserService
}

var service *Service

func GetService() *Service {
	if service == nil {
		service = Init()
	}
	return service
}

func NewService() *Service {
	return &Service{}
}

func Init() *Service {
	service := NewService()
	err := service.setChildServiceByReflect()
	if err != nil {
		panic(err)
	}
	return service
}

func (service *Service) setChildServiceByReflect() error {
	//反射方式赋值service结构体内容
	refV := reflect.Indirect(reflect.ValueOf(service))
	refT := reflect.TypeOf(service).Elem()
	for i := 0; i < refV.NumField(); i++ {
		fieldT := refT.Field(i)
		fieldV := refV.Field(i)
		fieldTypeName := fieldT.Type.Elem().Name()
		newFunc, ok := childServiceContainer[utils.LcFirst(fieldTypeName)]
		if !ok {
			return errors.New("Service "+fieldTypeName + " initialize fail")
		}
		value := reflect.ValueOf(newFunc())
		fieldV.Set(value)
	}
	return nil
}

//func (service *Service)Error() string {
//	return "SetDao Method must set"
//}
