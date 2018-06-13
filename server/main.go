package main

import (
	"flag"
	"github.com/judwhite/go-svc/svc"
	"github.com/liyue201/etn/server/models"
	"github.com/liyue201/etn/server/config"
	"github.com/liyue201/etn/server/rest"
	"github.com/liyue201/go-logger"
	"sync"
	"syscall"
)


type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("WaitGroupWrapper Wrap %s", err)
			}
			w.Done()
		}()
		cb()
	}()
}

type Service struct {
	WaitGroupWrapper
	httpServer *rest.RestSever
}

func (s *Service) Init(env svc.Environment) error {
	s.httpServer = rest.NewHttpServer(config.Cfg.Port)
	logger.Info("service inited")
	return nil
}

func (s *Service) Start() error {
	s.Wrap(func() {
		s.httpServer.Run()
	})
	logger.Info("service start")
	return nil
}

func (s *Service) Stop() error {
	s.httpServer.Stop()
	s.Wait()
	logger.Info("service stopped")
	return nil
}

func main() {
	cfg := flag.String("C", "server.yml", "configuration file")
	flag.Parse()

	err := config.InitConfig(*cfg)
	if err != nil {
		logger.Errorf("[main] %s", err)
		return
	}
	dbCfg := config.Cfg.Db
	err = models.InitDb(dbCfg.User,dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Dbname)
	if err != nil {
		logger.Errorf("[main] %s", err)
		return
	}

	service := &Service{}
	if err := svc.Run(service, syscall.SIGINT, syscall.SIGTERM); err != nil {
		logger.Errorf("main:", err)
	}
}
