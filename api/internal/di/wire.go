package di

import (
	database "github.com/ctf-api/internal/db"
	"github.com/ctf-api/internal/handlers"
	"github.com/ctf-api/internal/repository"
	"github.com/ctf-api/internal/server"
	"github.com/ctf-api/internal/usecase"

	"github.com/google/wire"
)

func InitializeServer1() *server.Server {
	wire.Build(
		database.InitDB,
		repository.NewRepository,
		usecase.NewUsecase,
		handlers.NewHandlers,
		repository.NewTeamRepository,
		usecase.NewTeamUsecase,
		handlers.NewTeamHandler,
		handlers.NewChallengeHandler,
		usecase.NewChallengeUseCase,
		repository.NewChallengeRepository,
		server.NewServer,
	)
	return &server.Server{}
}
