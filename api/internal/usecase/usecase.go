package usecase

import (
	"log"
	"github.com/ctf-api/internal/models"
	"github.com/ctf-api/internal/repository"
)

type Usecase struct {
	repo *repository.Repository
}

func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetTeams() ([]models.Team, error) {
	log.Println("usecase ")
	return u.repo.GetTeams()
}

func (u *Usecase) GetChallenges() ([]models.Challenge, error) {
	return u.repo.GetChallenges()
}

func (u *Usecase) GetScores() ([]models.Team, error) {
	return u.repo.GetScores()
}

func (u *Usecase) UpdateScore(teamName, challengeName string, score int) error {
	return u.repo.UpdateScore(teamName, challengeName, score)
}
