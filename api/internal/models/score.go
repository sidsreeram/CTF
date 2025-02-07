package models

import (
	"time"
)


type Score struct {
	TeamName      string    `json:"team_name"`
	ChallengeName string    `json:"challenge_name"`
	ChallengeID   int      `gorm:"not null" json:"challenge_id"`
	Score         int       `json:"score"`
	Timestamp     time.Time `json:"timestamp"`
}
