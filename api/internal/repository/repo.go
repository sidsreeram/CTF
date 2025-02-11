package repository

import (
	"errors"
	"time"

	"github.com/ctf/api/internal/models"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}
type Repository interface {
	GetTeams() ([]models.Team, error)
	GetChallenges() ([]models.Challenge, error) 
	GetScores() ([]models.Team, error)
	UpdateScore(teamName, challengeName string, score int) error 
	
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetTeams - Fetch all teams from DB
func (r *repository) GetTeams() ([]models.Team, error) {
	var teams []models.Team
	 err := r.db.Select("id,name, score,is_blocked").Find(&teams).Error
	return teams, err
}

// GetChallenges - Fetch all challenges from DB
func (r *repository) GetChallenges() ([]models.Challenge, error) {
	var challenges []models.Challenge
	err := r.db.Find(&challenges).Error
	return challenges, err
}

// GetScores - Fetch all scores from DB
func (r *repository) GetScores() ([]models.Team, error) {
	var teams []models.Team
	if err := r.db.Select("name, score").Order("score DESC").Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

// UpdateScore - Update a team's score in DB
func (r *repository) UpdateScore(teamName, challengeName string, score int) error {
	var existingScore models.Score
	// Check if the score already exists for this team and challenge
	if err := r.db.Where("team_id = ? ", teamName, challengeName).First(&existingScore).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no existing record, create a new one
			newScore := models.Score{
				TeamName:      teamName,
				Score:         score,
				Timestamp:     time.Now(), // Use time.Now() for current timestamp
			}
			// Using GORM's Create method to insert new score
			if err := r.db.Create(&newScore).Error; err != nil {
				return err
			}
		} else {
			// If any other error occurs
			return err
		}
	} else {
		// If score exists, update it
		existingScore.Score = score
		existingScore.Timestamp = time.Now() // Use time.Now() for current timestamp
		// Using GORM's Save method to update the existing score
		if err := r.db.Save(&existingScore).Error; err != nil {
			return err
		}
	}
	return nil
}

