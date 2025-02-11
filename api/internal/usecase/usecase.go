package usecase

import (
	"github.com/ctf/api/internal/models"
	"github.com/ctf/api/internal/repository"
)

type usecase struct {
	repo repository.Repository
}
type UseCase interface {
	GetTeams() ([]models.Team, error)
	GetChallenges() ([]models.Challenge, error)
	GetScores() ([]models.Team, error)
	UpdateScore(teamName, challengeName string, score int) error
}

func NewUsecase(repo repository.Repository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) GetTeams() ([]models.Team, error) {
	return u.repo.GetTeams()
}

func (u *usecase) GetChallenges() ([]models.Challenge, error) {
	return u.repo.GetChallenges()
}

func (u *usecase) GetScores() ([]models.Team, error) {
	return u.repo.GetScores()
}

func (u *usecase) UpdateScore(teamName, challengeName string, score int) error {
	return u.repo.UpdateScore(teamName, challengeName, score)
}
