// challenge_usecase.go
package usecase

import (
	"github.com/ctf/api/internal/repository"
	"github.com/ctf/api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type challengeUseCase struct {
	challengeRepo repository.ChallengeRepository
}
type ChallengeUsecase interface {
	CreateChallenge(name, description, downloadlink, hint, flagText string, score int) error
	VerifyAndSubmitFlag(teamID, challengeID int, submittedFlag string) (bool, error) 
	GetChallenges() ([]models.Challenge, error)
	DeleteChallenge(id int) error
	
}

func NewChallengeUseCase(repo repository.ChallengeRepository) ChallengeUsecase {
	return &challengeUseCase{
		challengeRepo: repo,
	}
}

// CreateChallenge creates a new challenge with a hashed flag
func (uc *challengeUseCase) CreateChallenge(name, description, downloadlink, hint, flagText string, score int) error {
	// Hash the flag before storing
	hashedFlag, err := bcrypt.GenerateFromPassword([]byte(flagText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	challenge := &models.Challenge{
		Name:         name,
		Description:  description,
		Flag:         string(hashedFlag),
		DownloadLink: downloadlink,
		Hint:         hint,
		Score:        score,
		Count:        0,
	}

	return uc.challengeRepo.CreateChallenge(challenge)
}

// VerifyAndSubmitFlag verifies a flag submission and updates team score
func (uc *challengeUseCase) VerifyAndSubmitFlag(teamID, challengeID int, submittedFlag string) (bool, error) {
	return uc.challengeRepo.VerifyFlag(teamID, challengeID, submittedFlag)
}

// GetChallenges retrieves all available challenges
func (uc *challengeUseCase) GetChallenges() ([]models.Challenge, error) {
	return uc.challengeRepo.GetChallenges()
}


// DeleteChallenge removes a challenge by name
func (uc *challengeUseCase) DeleteChallenge(id int) error {
	return uc.challengeRepo.DeleteChallenge(id)
}
