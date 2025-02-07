// challenge_usecase.go
package usecase

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/ctf-api/internal/models"
	"github.com/ctf-api/internal/repository"
)

type ChallengeUseCase struct {
	challengeRepo repository.ChallengeRepository
}

func NewChallengeUseCase(repo repository.ChallengeRepository) *ChallengeUseCase {
	return &ChallengeUseCase{
		challengeRepo: repo,
	}
}

// CreateChallenge creates a new challenge with a hashed flag
func (uc *ChallengeUseCase) CreateChallenge(name, description, downloadlink, hint, flagText string, score int) error {
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
func (uc *ChallengeUseCase) VerifyAndSubmitFlag(teamID, challengeID int, submittedFlag string) (bool, error) {
	return uc.challengeRepo.VerifyFlag(teamID, challengeID, submittedFlag)
}

// GetChallenges retrieves all available challenges
func (uc *ChallengeUseCase) GetChallenges() ([]models.Challenge, error) {
	return uc.challengeRepo.GetChallenges()
}

// UpdateChallenge allows updating challenge details
func (uc *ChallengeUseCase) UpdateChallenge(challenge *models.Challenge) error {
	return uc.challengeRepo.UpdateChallenge(challenge)
}

// DeleteChallenge removes a challenge by name
func (uc *ChallengeUseCase) DeleteChallenge(id int) error {
	return uc.challengeRepo.DeleteChallenge(id)
}
