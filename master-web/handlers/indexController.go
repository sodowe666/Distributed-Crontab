package handlers

import "distributedcrontab/master-web/service"

func Index()  {
	service.GetInstance().UserService.AAA()

}