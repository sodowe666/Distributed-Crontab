package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
	"distributedcrontab/master-web/service"
)

type handler func(server *Server)

type Server struct {
	*gin.Engine
	Handlers []handler
}

//new一个server
func NewServer(env string,cfgHandlers... handler) ServerInterface {
	server := new(Server)
	//设置模式
	gin.SetMode(env)
	server.Engine = gin.Default()
	//执行配置项
	for _,cfgHandler := range cfgHandlers{
		cfgHandler(server)
	}
	return server
}

//服务初始化设置
func (server *Server) Bootstrap() *Server {
	//初始化路由
	server.initRouter()
	//初始化service
	service.Init()
	return server
}

//服务启动
func (server *Server) Start() error {
	//debug模式下的处理
	if gin.Mode() == gin.DebugMode {

	}
	return server.run()
}

//server run
func (server *Server) run() error {
	port := viper.Get("server.port")
	//
	readTimerout := viper.Get("server.readTimeout").(int)
	writeTimeout := viper.Get("server.writeTimeout").(int)
	s := &http.Server{
		Handler:        server,
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    time.Duration(readTimerout),
		WriteTimeout:    time.Duration(writeTimeout),
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}
