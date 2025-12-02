package lib

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

const MaxLevel = 25
const ActivitiesRequiredPerLevel = 12
const CommentExperience = 100

const (
	ActivityComment = iota
	ActivityPost
	ActivityImagePost
)

type UserLevelInfo struct {
	Level                     int     `json:"level"`
	Experience                float64 `json:"experience"`
	PercentageOfLevelComplete float64 `json:"percentageOfLevelComplete"`
}

func AddExperienceToUserId(dbPool *pgxpool.Pool, experience int, userId int) (float64, error) {
	const query = `
		INSERT INTO user_experience (user_id, experience, updated_at) VALUES ($1, $2, NOW())
		ON CONFLICT(user_id)
		    DO UPDATE SET experience = user_experience.experience + $2, updated_at = NOW()
		RETURNING user_experience.experience
	`
	var updatedExperience float64
	updateErr := dbPool.QueryRow(
		context.Background(),
		query,
		userId,
		experience).Scan(&updatedExperience)
	if updateErr != nil {
		return 0, updateErr
	}
	return updatedExperience, nil
}

func AddCommentExperienceToUser(dbPool *pgxpool.Pool, userId int) (float64, error) {
	activityTypeExperienceMap := GetActivityTypeExperienceMap()
	return AddExperienceToUserId(dbPool, activityTypeExperienceMap[ActivityComment], userId)
}

func AddPostExperienceToUser(dbPool *pgxpool.Pool, userId int) (float64, error) {
	activityTypeExperienceMap := GetActivityTypeExperienceMap()
	return AddExperienceToUserId(dbPool, activityTypeExperienceMap[ActivityPost], userId)
}

func AddImagePostExperienceToUser(dbPool *pgxpool.Pool, userId int) (float64, error) {
	activityTypeExperienceMap := GetActivityTypeExperienceMap()
	return AddExperienceToUserId(dbPool, activityTypeExperienceMap[ActivityImagePost], userId)
}

func GetActivityTypeExperienceMap() map[int]int {
	activityTypeExperienceMap := make(map[int]int)
	activityTypeExperienceMap[ActivityComment] = 100
	activityTypeExperienceMap[ActivityPost] = 150
	activityTypeExperienceMap[ActivityImagePost] = 200
	return activityTypeExperienceMap
}

func GetExperienceByActivityType(activityType int) (int, error) {
	activityTypeExperienceMap := GetActivityTypeExperienceMap()
	val, ok := activityTypeExperienceMap[activityType]
	if !ok {
		return 0, errors.New("activity type not found")
	}
	return val, nil
}

func GetLevelExperienceMap() map[int]float64 {
	userLevels := make(map[int]float64)
	for level := range MaxLevel {
		if level > 0 {
			userLevels[level] = (CommentExperience * ActivitiesRequiredPerLevel) * float64(level)
		}
	}
	return userLevels
}

func GetUserLevelByExperience(experience float64) int {
	levelExperienceMap := GetLevelExperienceMap()
	for level := range levelExperienceMap {
		if experience >= levelExperienceMap[level] {
			return level
		}
	}
	return 1
}

func GetPercentageOfLevelComplete(experience float64, level int) float64 {
	levelExperienceMap := GetLevelExperienceMap()
	levelExperience := levelExperienceMap[level]
	return experience / levelExperience
}
