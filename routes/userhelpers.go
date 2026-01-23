package routes

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

type CreateUserRequest struct {
	T              *testing.T
	E              *httpexpect.Expect
	Username       string
	Password       string
	AvatarFilename string
}

func CreateUserAndVerify(request CreateUserRequest) {

}
