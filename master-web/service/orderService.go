package service

import (
	"fmt"
	"distributedcrontab/master-web/dao"
)

type OrderService struct {
	dao *dao.OrderDao
}

func NewOrderService() ServiceInterface {
	return &OrderService{
		dao: nil,
	}
}

func (order *OrderService) SetDao() {
	order.dao = dao.NewOrderDao()
}

func (order *OrderService) Aaa() {
	fmt.Printf("bbb")
}
