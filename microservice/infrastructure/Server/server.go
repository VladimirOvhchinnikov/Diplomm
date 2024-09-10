package server

import (
	"context"
	"diplom/infrastructure/config"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	server *http.Server
}

func NewServer(logger *zap.Logger, addr string, handler http.Handler) *Server {

	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

}

func (s *Server) Start() error {

	s.logger.Info(fmt.Sprintf("Server started on %s", s.server.Addr))
	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Error("The server could not start", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("Server exited successfully")
	return nil
}

func (s *Server) Restart() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		s.logger.Error("The server was unable to shut down properly.", zap.Error(err))
		return err
	}

	err = s.Start()
	if err != nil {
		s.logger.Error("The server was unable to start correctly.", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) RestartWithRetries(maxRetries int, PathConfig string) error {

	var serverMirror *Server
	configs, err := config.LoadConfig(PathConfig)
	if err != nil {
		s.logger.Error("Failed to obtain configuration files to start mirror server", zap.Error(err))
	} else {
		var mirrorStarted bool
		for _, config := range configs.Servers {
			serverMirror = NewServer(s.logger, config.Addr, s.server.Handler)
			err := serverMirror.Start()
			if err == nil {
				mirrorStarted = true
				break
			}
		}

		// Логируем только если ни один зеркальный сервер не запустился
		if !mirrorStarted {
			s.logger.Error("Failed to start any mirror server. Proceeding with primary server restart")
		}
	}

	// Перезапуск основного сервера
	var lastErr error
	for retries := 0; retries < maxRetries; retries++ {
		err := s.Restart()
		if err != nil {
			lastErr = err
			s.logger.Warn(fmt.Sprintf("Restart failed, retrying in %d seconds", retries+1), zap.Error(err))
			time.Sleep(time.Duration(retries+1) * time.Second)
		} else {
			s.logger.Info("Server restarted successfully")

			// Задержка перед завершением зеркального сервера
			time.Sleep(15 * time.Minute)

			// Завершаем зеркальный сервер с использованием контекста и таймаута
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			serverMirror.Shutdown(ctx)

			return nil
		}
	}

	return fmt.Errorf("server failed to restart after %d retries: last error: %v", maxRetries, lastErr)
}
