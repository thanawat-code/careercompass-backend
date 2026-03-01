package models

import (
	"time"

	"github.com/google/uuid"
)

// LearningPath represents a career learning path
type LearningPath struct {
	ID          uuid.UUID `json:"id"`
	CareerName  string    `json:"career_name"`
	Description string    `json:"description"`
	TotalStages int       `json:"total_stages"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Stage represents a stage in a learning path
type Stage struct {
	ID               uuid.UUID  `json:"id"`
	LearningPathID   uuid.UUID  `json:"learning_path_id"`
	StageNumber      int        `json:"stage_number"`
	Title            string     `json:"title"`
	Subtitle         string     `json:"subtitle"`
	PositionTop      *string    `json:"position_top"`
	PositionLeft     *string    `json:"position_left"`
	PositionRight    *string    `json:"position_right"`
	PositionBottom   *string    `json:"position_bottom"`
	PositionTransform *string   `json:"position_transform"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Course represents a course in a stage
type Course struct {
	ID        uuid.UUID `json:"id"`
	StageID   uuid.UUID `json:"stage_id"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	URL       *string   `json:"url"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStageProgress represents a user's progress on a stage
type UserStageProgress struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	StageID     uuid.UUID  `json:"stage_id"`
	Status      string     `json:"status"` // locked, in-progress, completed
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// StageWithStatus combines stage info with user's progress status
type StageWithStatus struct {
	ID               uuid.UUID  `json:"id"`
	StageNumber      int        `json:"stage_number"`
	Title            string     `json:"title"`
	Subtitle         string     `json:"subtitle"`
	Status           string     `json:"status"` // locked, in-progress, completed
	PositionTop      *string    `json:"position_top"`
	PositionLeft     *string    `json:"position_left"`
	PositionRight    *string    `json:"position_right"`
	PositionBottom   *string    `json:"position_bottom"`
	PositionTransform *string   `json:"position_transform"`
	Courses          []Course   `json:"courses"`
}

// LearningPathResponse is the full response for a learning path with user progress
type LearningPathResponse struct {
	ID              uuid.UUID         `json:"id"`
	CareerName      string            `json:"career_name"`
	Description     string            `json:"description"`
	TotalStages     int               `json:"total_stages"`
	CompletedStages int               `json:"completed_stages"`
	Stages          []StageWithStatus `json:"stages"`
}

// UpdateProgressRequest is used to update a user's stage progress
type UpdateProgressRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	StageID string `json:"stage_id" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=locked in-progress completed"`
}
