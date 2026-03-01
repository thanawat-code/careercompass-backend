package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/models"
)

// GetLearningPath returns the learning path for a specific career with user progress
// GET /api/learning-path/:career_name?user_id=xxx
func GetLearningPath(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		careerName := c.Param("career_name")
		userIDStr := c.Query("user_id")

		ctx := context.Background()

		// 1. Get learning path info
		var lp models.LearningPath
		err := db.Pool.QueryRow(ctx,
			`SELECT id, career_name, description, total_stages, created_at, updated_at
			 FROM learning_paths WHERE career_name = $1`,
			careerName,
		).Scan(&lp.ID, &lp.CareerName, &lp.Description, &lp.TotalStages, &lp.CreatedAt, &lp.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Learning path not found"})
			return
		}

		// 2. Get all stages for this learning path
		stageRows, err := db.Pool.Query(ctx,
			`SELECT id, stage_number, title, subtitle,
			        position_top, position_left, position_right, position_bottom, position_transform
			 FROM stages WHERE learning_path_id = $1 ORDER BY stage_number ASC`,
			lp.ID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stages"})
			return
		}
		defer stageRows.Close()

		var stages []models.StageWithStatus
		for stageRows.Next() {
			var s models.StageWithStatus
			if err := stageRows.Scan(
				&s.ID, &s.StageNumber, &s.Title, &s.Subtitle,
				&s.PositionTop, &s.PositionLeft, &s.PositionRight, &s.PositionBottom, &s.PositionTransform,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan stage"})
				return
			}
			// Default status
			s.Status = "locked"
			s.Courses = []models.Course{}
			stages = append(stages, s)
		}

		// 3. If user_id provided, overlay user progress
		completedStages := 0
		if userIDStr != "" {
			userID, err := uuid.Parse(userIDStr)
			if err == nil {
				progressRows, err := db.Pool.Query(ctx,
					`SELECT stage_id, status FROM user_stage_progress WHERE user_id = $1`,
					userID,
				)
				if err == nil {
					defer progressRows.Close()
					progressMap := make(map[uuid.UUID]string)
					for progressRows.Next() {
						var stageID uuid.UUID
						var status string
						if err := progressRows.Scan(&stageID, &status); err == nil {
							progressMap[stageID] = status
						}
					}
					// Apply progress to stages
					for i := range stages {
						if status, ok := progressMap[stages[i].ID]; ok {
							stages[i].Status = status
						}
						if stages[i].Status == "completed" {
							completedStages++
						}
					}
				}
			}
		} else {
			// No user => first stage is in-progress by default
			if len(stages) > 0 {
				stages[0].Status = "in-progress"
			}
		}

		// 4. Get courses for each stage
		for i := range stages {
			courseRows, err := db.Pool.Query(ctx,
				`SELECT id, stage_id, title, subtitle, url, sort_order, created_at, updated_at
				 FROM courses WHERE stage_id = $1 ORDER BY sort_order ASC`,
				stages[i].ID,
			)
			if err != nil {
				continue
			}
			defer courseRows.Close()

			for courseRows.Next() {
				var course models.Course
				if err := courseRows.Scan(
					&course.ID, &course.StageID, &course.Title, &course.Subtitle,
					&course.URL, &course.SortOrder, &course.CreatedAt, &course.UpdatedAt,
				); err == nil {
					stages[i].Courses = append(stages[i].Courses, course)
				}
			}
		}

		if stages == nil {
			stages = []models.StageWithStatus{}
		}

		response := models.LearningPathResponse{
			ID:              lp.ID,
			CareerName:      lp.CareerName,
			Description:     lp.Description,
			TotalStages:     lp.TotalStages,
			CompletedStages: completedStages,
			Stages:          stages,
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetAllLearningPaths returns all available career learning paths
// GET /api/learning-paths
func GetAllLearningPaths(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		rows, err := db.Pool.Query(ctx,
			`SELECT id, career_name, description, total_stages, created_at, updated_at
			 FROM learning_paths ORDER BY career_name ASC`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch learning paths"})
			return
		}
		defer rows.Close()

		var paths []models.LearningPath
		for rows.Next() {
			var lp models.LearningPath
			if err := rows.Scan(&lp.ID, &lp.CareerName, &lp.Description, &lp.TotalStages, &lp.CreatedAt, &lp.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan learning path"})
				return
			}
			paths = append(paths, lp)
		}

		if paths == nil {
			paths = []models.LearningPath{}
		}

		c.JSON(http.StatusOK, paths)
	}
}

// UpdateUserProgress updates user progress for a specific stage
// POST /api/learning-path/progress
func UpdateUserProgress(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UpdateProgressRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		stageID, err := uuid.Parse(req.StageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stage_id"})
			return
		}

		ctx := context.Background()

		// Upsert user stage progress
		query := `
			INSERT INTO user_stage_progress (id, user_id, stage_id, status, started_at, completed_at)
			VALUES (gen_random_uuid(), $1, $2, $3,
				CASE WHEN $3 = 'in-progress' THEN NOW() ELSE NULL END,
				CASE WHEN $3 = 'completed' THEN NOW() ELSE NULL END
			)
			ON CONFLICT (user_id, stage_id) DO UPDATE
			SET status = EXCLUDED.status,
			    started_at = CASE WHEN EXCLUDED.status = 'in-progress' AND user_stage_progress.started_at IS NULL THEN NOW() ELSE user_stage_progress.started_at END,
			    completed_at = CASE WHEN EXCLUDED.status = 'completed' THEN NOW() ELSE NULL END,
			    updated_at = NOW()
			RETURNING id, user_id, stage_id, status, started_at, completed_at, created_at, updated_at`

		var progress models.UserStageProgress
		err = db.Pool.QueryRow(ctx, query, userID, stageID, req.Status).Scan(
			&progress.ID, &progress.UserID, &progress.StageID, &progress.Status,
			&progress.StartedAt, &progress.CompletedAt, &progress.CreatedAt, &progress.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update progress"})
			return
		}

		c.JSON(http.StatusOK, progress)
	}
}

// GetUserProgress returns all stage progress for a user
// GET /api/learning-path/progress/:user_id
func GetUserProgress(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}

		ctx := context.Background()

		rows, err := db.Pool.Query(ctx,
			`SELECT id, user_id, stage_id, status, started_at, completed_at, created_at, updated_at
			 FROM user_stage_progress WHERE user_id = $1`,
			userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
			return
		}
		defer rows.Close()

		var progresses []models.UserStageProgress
		for rows.Next() {
			var p models.UserStageProgress
			if err := rows.Scan(
				&p.ID, &p.UserID, &p.StageID, &p.Status,
				&p.StartedAt, &p.CompletedAt, &p.CreatedAt, &p.UpdatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan progress"})
				return
			}
			progresses = append(progresses, p)
		}

		if progresses == nil {
			progresses = []models.UserStageProgress{}
		}

		c.JSON(http.StatusOK, progresses)
	}
}
