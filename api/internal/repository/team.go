package repository

import (
    "errors"
    "fmt"
    "regexp"
    "strings"

    "github.com/ctf-api/internal/models"
    "gorm.io/gorm"
)

type TeamRepository interface {
    CreateTeam(team *models.Team) error
    GetTeamByName(name string) (*models.Team, error)
    ValidateTeamName(name string) error
    BlockTeam(teamName string) error
}

type teamRepo struct {
    db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
    return &teamRepo{db: db}
}

// ValidateTeamName ensures the team name is well-formed
func (r *teamRepo) ValidateTeamName(name string) error {
    name = strings.TrimSpace(name)
    if len(name) < 3 || len(name) > 30 {
        return errors.New("team name must be between 3 and 30 characters")
    }

    // Only allow alphanumeric characters and underscores
    validName := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !validName.MatchString(name) {
        return errors.New("team name can only contain letters, numbers, and underscores")
    }
    return nil
}

func (r *teamRepo) CreateTeam(team *models.Team) error {
    if err := r.ValidateTeamName(team.Name); err != nil {
        return err
    }

    // Use GORM's Create method to insert the team into the database
    if err := r.db.Create(team).Error; err != nil {
        return fmt.Errorf("failed to create team: %w", err)
    }
    return nil
}

func (r *teamRepo) GetTeamByName(name string) (*models.Team, error) {
    if err := r.ValidateTeamName(name); err != nil {
        return nil, err
    }

    // Use GORM's Where and First methods to fetch the team by name
    var team models.Team
    if err := r.db.Where("name = ?", name).First(&team).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil // No team found
        }
        return nil, fmt.Errorf("failed to fetch team: %w", err)
    }
    return &team, nil
}

func (r *teamRepo) BlockTeam(teamName string) error {
    return r.db.Model(&models.Team{}).
        Where("name = ?", teamName).
        Update("is_blocked", true).Error
}
