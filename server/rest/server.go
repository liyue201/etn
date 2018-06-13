package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liyue201/etn/server/config"
	"github.com/liyue201/etn/server/models"
	"github.com/liyue201/go-logger"
	"net/http"
)

type RestSever struct {
	httpServer *http.Server
}

func NewHttpServer(port int) *RestSever {
	gin.SetMode(gin.ReleaseMode)
	engin := gin.Default()

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	httpServer := &http.Server{Addr: addr, Handler: engin}
	server := &RestSever{
		httpServer: httpServer,
	}

	engin.Static("/static", config.Cfg.Static)

	server.initRoute(engin)

	return server
}

func (s *RestSever) Run() {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		logger.Errorf("RestSever.Run %s", err)
	}
}

func (s *RestSever) Stop() {
	s.httpServer.Shutdown(context.Background())
}

func (s *RestSever) initRoute(r gin.IRouter) {
	r.GET("/api/v1/version", s.GetVersion)
	r.GET("/api/v1/files", s.GetFiles)
}

func (s *RestSever) GetVersion(c *gin.Context) {
	RespJson(c, OK, config.Cfg.Version)
}

func (s *RestSever) GetFiles(c *gin.Context) {
	files := models.GetFiles()
	RespJson(c, OK, files)
}
