package usecase

import (
	"errors"
	"fmt"

	"github.com/ctf-api/internal/middleware"
	"github.com/ctf-api/internal/models"
	"github.com/ctf-api/internal/repository"
)

type TeamUsecase interface {
	RegisterTeam(teamName, password string) error
	LoginTeam(teamName, password string) (string, string, int, error)
}

type teamUsecase struct {
	repo repository.TeamRepository
}

func NewTeamUsecase(repo repository.TeamRepository) TeamUsecase {
	return &teamUsecase{repo: repo}
}

func (u *teamUsecase) RegisterTeam(teamName, password string) error {
	// Hash password
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		return err
	}
	fmt.Println(hashedPassword)
	team := &models.Team{
		Name:         teamName,
		PasswordHash: hashedPassword,
	}

	return u.repo.CreateTeam(team)
}

func (u *teamUsecase) LoginTeam(teamName, password string) (string, string, int, error) {
	team, err := u.repo.GetTeamByName(teamName)
	if err != nil {
		return "", "", 0, errors.New("team not found")
	}

	if team.IsBlocked {
		return "", "", 0, errors.New("team is blocked by admin")
	}

	if !models.CheckPassword(team.PasswordHash, password) {
		return "", "", 0, errors.New("invalid credentials")
	}

	var token string
	var errToken error
	var role string

	if team.IsAdmin {
		token, errToken = middleware.GenerateAdminToken(int(team.ID))
		role = "admin"
	} else {
		token, errToken = middleware.GenerateTeamToken(int(team.ID))
		role = "team"
	}

	if errToken != nil {
		return "", "", 0, errToken
	}

	return token, role, int(team.ID), nil
}
