package models

import (
	"time"
)

type Score struct {
    ID          uint      `gorm:"primarykey"`
    TeamID      int       `gorm:"not null"`
    TeamName    string    `gorm:"not null"`
    ChallengeID int       `gorm:"not null"`
    Score       int       `gorm:"not null"`
    Timestamp   time.Time `gorm:"not null"`
    
}
