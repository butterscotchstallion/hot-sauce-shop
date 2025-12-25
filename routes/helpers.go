package routes

import "github.com/google/uuid"

func GenerateUniqueName() string {
	postUUID, postUUIDErr := uuid.NewRandom()
	if postUUIDErr != nil {
		panic("Failed to generate post UUID")
	}
	return postUUID.String()
}
