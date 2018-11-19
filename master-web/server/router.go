package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//路由配置
func (server *Server) initRouter() {
	router := server
	initRouter(router)
}

//路由配置区域
func initRouter(router *Server) {
	router.GET("aaa", func(context *gin.Context) {
		context.String(http.StatusOK, "asdfasdfasdfasdf")
	})
}
