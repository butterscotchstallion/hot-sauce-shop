package lib

const MaxLevel = 25
const ActivitiesRequiredPerLevel = 12
const CommentExperience = 100

func GetLevelExperienceMap() map[int]int {
	userLevels := make(map[int]int)
	for level := range MaxLevel {
		userLevels[level] = CommentExperience * ActivitiesRequiredPerLevel
	}
	return userLevels
}

func GetUserLevelByExperience(experience int) int {
	levelExperienceMap := GetLevelExperienceMap()
	for level := range levelExperienceMap {
		if experience >= levelExperienceMap[level] {
			return level
		}
	}
	return 1
}
