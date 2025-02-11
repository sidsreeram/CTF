package repository

import (
	"errors"
	"time"

	"github.com/ctf/api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ChallengeRepository interface {
	CreateChallenge(challenge *models.Challenge) error
	GetChallengeByName(name string) (*models.Challenge, error)
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
	tx := r.db.Begin() // Start transaction
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock challenge row to prevent concurrent updates
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&challenge, challengeID).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// Check if team already solved this challenge
	var existingScore models.Score
	if err := tx.Where("team_id = ? AND challenge_id = ?", teamID, challengeID).First(&existingScore).Error; err == nil {
		tx.Rollback()
		return false, errors.New("team has already solved this challenge")
	}

	// Verify flag using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(challenge.Flag), []byte(submittedFlag))
	if err != nil {
		tx.Rollback()
		return false, nil // Incorrect flag
	}

	// **🔹 Calculate the score dynamically based on solver count**
	var solverCount int64
	tx.Model(&models.Score{}).Where("challenge_id = ?", challengeID).Count(&solverCount)

	var score int
	switch solverCount {
	case 0:
		score = challenge.Score // First solver
	case 1:
		score = int(float64(challenge.Score) * 0.9) // Second solver
	case 2:
		score = int(float64(challenge.Score) * 0.8) // Third solver
	default:
		score = int(float64(challenge.Score) * 0.7) // Others
	}

	// **🔹 Increment challenge count**
	challenge.Count++
	if err := tx.Save(&challenge).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// **🔹 Update team score**
	var team models.Team
	if err := tx.First(&team, teamID).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	team.Score += score
	if err := tx.Save(&team).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// **🔹 Store score in DB**
	newScore := models.Score{
		TeamID:      teamID,
		TeamName:    team.Name,
		ChallengeID: challengeID,
		Score:       score,
		Timestamp:   time.Now(),
	}
	if err := tx.Create(&newScore).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()

	return true, nil
}
