package main

import (
	"distributedcrontab/master-web/server"
	"runtime"
	"os"
	"github.com/gin-gonic/gin"
	"io"
	"github.com/spf13/pflag"
	"distributedcrontab/common/system"
)

func main() {
	//设置线程数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//加载日志
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	//配置加载
	env := pflag.StringP("env", "e", "debug", "run the environment, dev or prod")
	configFilePath := pflag.StringP("config", "c", "", "config file path")
	pflag.Parse()
	if err := system.SetupConfig(*env, *configFilePath); err != nil {
		panic("config load fail")
	}
	//实例化server
	s := server.NewServer(*env)
	//服务初始化，初始化路由
	s.Bootstrap()
	//服务启动
	if err := s.Start(); err != nil {
		panic("server start Error: " + err.Error())
	}
}
