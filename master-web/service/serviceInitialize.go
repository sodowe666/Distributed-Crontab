package service

import (
	"reflect"
	"distributedcrontab/common/utils"
	"errors"
)

type ServiceInterface interface {
}

var service *Service


//以下为service容器实现，实例化serviceConfig.go文件里Service结构体里的元素并赋值，不用管
func GetInstance() *Service {
	if service == nil {
		service = Init()
	}
	return service
}

func NewService() *Service {
	return &Service{}
}

//初始化
func Init() *Service {
	service := NewService()
	err := service.setChildServiceByReflect()
	if err != nil {
		panic(err)
	}
	return service
}

//反射解析Service结构体里的元素，并赋值
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