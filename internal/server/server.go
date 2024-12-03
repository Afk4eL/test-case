package server

import (
	"fmt"
	"net/http"
	"test-case/config"
	"test-case/internal/user"
	user_service "test-case/internal/user/delivery/http"
	"test-case/internal/user/repository"
	"test-case/internal/user/usecase"
	"test-case/internal/utils/logger"

	"gorm.io/gorm"
)

type Server struct {
	cfg        *config.Config
	db         *gorm.DB
	httpServer *http.Server
}

func NewNotesServer(cfg *config.Config, db *gorm.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) httpRun(userUC user.UserUseCase) error {
	const op = "server.httpRun"

	httpHandlers := user_service.NewHandlers(userUC)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.Server.Addr, s.cfg.Server.Port),
		Handler: user_service.SetupRouter(httpHandlers),
	}

	logger.Logger.Info().Msg("Http server started on " + s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatal().Msg(fmt.Sprint("Fatal error", op, err.Error()))
		return err
	}

	return nil
}

func (s *Server) Run() error {
	const op = "grpc_app.Run"

	noteRepo := repository.NewUserRepository(s.db)

	noteUC := usecase.NewUserUsecase(noteRepo)

	go s.httpRun(noteUC)

	return nil
}

func (s *Server) Stop() {
	logger.Logger.Info().Msg("stopping http server")

}
