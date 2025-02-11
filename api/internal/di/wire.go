package di

import (
	"github.com/ctf/api/internal/config"
	database "github.com/ctf/api/internal/db"
	"github.com/ctf/api/internal/handlers"
	"github.com/ctf/api/internal/repository"
	"github.com/ctf/api/internal/server"
	"github.com/ctf/api/internal/usecase"

	"github.com/google/wire"
)

// ProvideConfig loads the configuration
func ProvideConfig1() (config.Config, error) {
    return config.LoadConfig()
}

var providerSet1 = wire.NewSet(
    ProvideConfig1,          // Provide config.Config
    database.InitDB,              // Provide *gorm.DB
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

// InitializeServer generates the Wire dependency injection
func InitializeServer1() (*server.Server, error) {
    wire.Build(providerSet)
    return nil, nil
}
