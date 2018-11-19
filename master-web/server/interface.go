package server

type ServerInterface interface {
	Start() error
	Bootstrap() *Server
	run() error
}
