package models

import "time"

type Flag struct {
	ID   uint `gorm:"primaryKey"`
	Flag string
}

type Timer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    string    `gorm:"not null" json:"status"` // "running" or "frozen"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Challenge struct {
	ID           int    `json:"id"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	DownloadLink string `json:"download_link"`
	Hint         string `json:"hint"`
	Flag         string `json:"-"` // Hashed, not exposed
	Score        int    `json:"score"`
	Count        int    `json:"count"`
}
type Submission struct {
	ID          int       `gorm:"primarykey" json:"id"`
	TeamID      int       `gorm:"not null" json:"team_id"`
	ChallengeId int       `gorm:"not null" json:"challenge_id"`
	IsCorrect   bool      `json:"is_correct"`
}
