package lib

import (
	"testing"
)

func TestGetLevelByExperience(t *testing.T) {
	expectedLevel := 2
	levelExperienceMap := GetLevelExperienceMap()
	experience := levelExperienceMap[expectedLevel]
	actualLevel := GetUserLevelByExperience(experience)
	if actualLevel != expectedLevel {
		t.Fatalf("GetLevelByExperience() expected %d, got %d", expectedLevel, actualLevel)
	}
}

func TestGetPercentageOfLevelComplete(t *testing.T) {
	level := 1
	expectedPercentage := 20.0
	levelExperienceMap := GetLevelExperienceMap()
	levelOneExperience := levelExperienceMap[level]
	experience := levelOneExperience * expectedPercentage
	actualPercentage := GetPercentageOfLevelComplete(experience, level)
	if actualPercentage != expectedPercentage {
		t.Fatalf("GetPercentageOfLevelComplete() expected %v, got %v", expectedPercentage, actualPercentage)
	}
}

func TestGetPercentageOfLevelCompleteErrata(t *testing.T) {
	// test error case
	zeroPercentage := GetPercentageOfLevelComplete(0, 1)
	if zeroPercentage != 0 {
		t.Fatalf("GetPercentageOfLevelComplete() expected 0, got %v", zeroPercentage)
	}
}
