package repository

import (
	"errors"
	"time"

	"github.com/ctf-api/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetTeams - Fetch all teams from DB
func (r *Repository) GetTeams() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Find(&teams).Error
	return teams, err
}

// GetChallenges - Fetch all challenges from DB
func (r *Repository) GetChallenges() ([]models.Challenge, error) {
	var challenges []models.Challenge
	err := r.db.Find(&challenges).Error
	return challenges, err
}

// GetScores - Fetch all scores from DB
func (r *Repository) GetScores() ([]models.Score, error) {
	var scores []models.Score
	if err := r.db.Order("timestamp ASC").Find(&scores).Error; err != nil {
		return nil, err
	}
	return scores, nil
}

// UpdateScore - Update a team's score in DB
func (r *Repository) UpdateScore(teamName, challengeName string, score int) error {
	var existingScore models.Score
	// Check if the score already exists for this team and challenge
	if err := r.db.Where("team_name = ? AND challenge_name = ?", teamName, challengeName).First(&existingScore).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no existing record, create a new one
			newScore := models.Score{
				TeamName:      teamName,
				ChallengeName: challengeName,
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

// BlockTeam - Marks a team as blocked
func (r *Repository) BlockTeam(teamName string) error {
	return r.db.Model(&models.Team{}).
		Where("team_name = ?", teamName).
		Update("is_blocked", true).Error
}

// SubmitFlag - Saves a submitted flag (independent of teams)
func (r *Repository) SubmitFlag(flag string) error {
	newFlag := models.Flag{Flag: flag}
	return r.db.Create(&newFlag).Error
}

// UpdateTimerStatus - Updates the CTF timer status
func (r *Repository) UpdateTimerStatus(status string) error {
	return r.db.Model(&models.Timer{}).
		Where("id = ?", 1).
		Update("status", status).Error
}