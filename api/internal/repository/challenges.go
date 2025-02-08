package repository

import (
	"errors"
	"time"

	"github.com/ctf-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ChallengeRepository interface {
	CreateChallenge(challenge *models.Challenge) error
	GetChallengeByName(name string) (*models.Challenge, error)
	UpdateChallenge(challenge *models.Challenge) error
	DeleteChallenge(id int) error
	GetChallenges() ([]models.Challenge, error)
	VerifyFlag(teamID int, challengeID int, flag string) (bool, error)
}

type ChallengeRepo struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &ChallengeRepo{db: db}
}

// CreateChallenge - Add a new challenge
func (r *ChallengeRepo) CreateChallenge(challenge *models.Challenge) error {
	return r.db.Create(challenge).Error
}

// GetChallengeByName - Fetch a challenge by its name
func (r *ChallengeRepo) GetChallengeByName(name string) (*models.Challenge, error) {
	var challenge models.Challenge
	err := r.db.Where("name = ?", name).First(&challenge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Challenge not found
		}
		return nil, err
	}
	return &challenge, nil
}

// UpdateChallenge - Update an existing challenge
func (r *ChallengeRepo) UpdateChallenge(challenge *models.Challenge) error {
	return r.db.Save(challenge).Error
}

// DeleteChallenge - Delete a challenge by its name
func (r *ChallengeRepo) DeleteChallenge(id int) error {
	return r.db.Where("id = ?", id).Delete(&models.Challenge{}).Error
}

// GetChallenges - Fetch all challenges
func (r *ChallengeRepo) GetChallenges() ([]models.Challenge, error) {
	var challenges []models.Challenge
	err := r.db.Find(&challenges).Error
	return challenges, err
}

// VerifyFlag - Verify the flag submitted by a team and update the team's score

// VerifyFlag checks if the submitted flag is correct and updates scores
func (r *ChallengeRepo) VerifyFlag(teamID int, challengeID int, submittedFlag string) (bool, error) {
	var challenge models.Challenge
	if err := r.db.First(&challenge, challengeID).Error; err != nil {
		return false, err
	}

	// Check if the team has already solved this challenge
	var existingScore models.Score
	if err := r.db.Where("team_id = ? AND challenge_id = ?", teamID, challengeID).First(&existingScore).Error; err == nil {
		return false, errors.New("team has already solved this challenge")
	}

	// Verify the flag using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(challenge.Flag), []byte(submittedFlag))
	if err != nil {
		return false, nil // Incorrect flag
	}

	// Calculate the score based on the number of solvers
	var score int
	switch challenge.Count {
	case 0:
		score = challenge.Score // First solver gets 100% of the base score
	case 1:
		score = int(float64(challenge.Score) * 0.9) // Second solver gets 90%
	case 2:
		score = int(float64(challenge.Score) * 0.8) // Third solver gets 80%
	default:
		score = int(float64(challenge.Score) * 0.7) // Others get 70%
	}

	// Update the challenge count
	challenge.Count++
	if err := r.db.Save(&challenge).Error; err != nil {
		return false, err
	}

	// Update the team's score
	var team models.Team
	if err := r.db.First(&team, teamID).Error; err != nil {
		return false, err
	}
	team.Score += score
	if err := r.db.Save(&team).Error; err != nil {
		return false, err
	}

	// Record the score for this team and challenge
	newScore := models.Score{
		TeamName:    team.Name,
		ChallengeID: challengeID,
		Score:       score,
		Timestamp:   time.Now(),
	}
	if err := r.db.Create(&newScore).Error; err != nil {
		return false, err
	}

	return true, nil
}
