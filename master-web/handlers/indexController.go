package handlers

import "distributedcrontab/master-web/service"

func aaa()  {
	service.GetInstance().UserService.AAA()
}