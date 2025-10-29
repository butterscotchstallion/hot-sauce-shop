package lib

import "testing"

func TestGetLevelByExperience(t *testing.T) {
	expectedLevel := 2
	experience := (CommentExperience * ActivitiesRequiredPerLevel) * 2
	actualLevel := GetUserLevelByExperience(experience)

	if actualLevel != expectedLevel {
		t.Fatalf("GetLevelByExperience() expected %d, got %d", expectedLevel, actualLevel)
	}
}
