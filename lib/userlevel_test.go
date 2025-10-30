package lib

import "testing"

func TestGetLevelByExperience(t *testing.T) {
	expectedLevel := 2
	experience := (CommentExperience * ActivitiesRequiredPerLevel) * expectedLevel
	actualLevel := GetUserLevelByExperience(experience)
	if actualLevel != expectedLevel {
		t.Fatalf("GetLevelByExperience() expected %d, got %d", expectedLevel, actualLevel)
	}
}
