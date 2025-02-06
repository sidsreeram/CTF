package usecase

import (
	"errors"

	"github.com/ctf-api/internal/models"
	"github.com/ctf-api/internal/repository"
)

type TeamUsecase interface {
	RegisterTeam(teamName, password string) error
	LoginTeam(teamName, password string) (bool, error)
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

	team := &models.Team{
		Name:     teamName,
		PasswordHash: hashedPassword,
	}

	return u.repo.CreateTeam(team)
}

func (u *teamUsecase) LoginTeam(teamName, password string) (bool, error) {
	team, err := u.repo.GetTeamByName(teamName)
	if err != nil {
		return false, errors.New("team not found")
	}

	if team.IsBlocked {
		return false, errors.New("team is blocked by admin")
	}

	if !models.CheckPassword(team.PasswordHash, password) {
		return false, errors.New("invalid credentials")
	}

	return true, nil
}
