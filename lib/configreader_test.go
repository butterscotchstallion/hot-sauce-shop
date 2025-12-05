package lib

import (
	"testing"
)

func TestReadExampleConfig(t *testing.T) {
	configData, err := ReadConfig("../config.example.toml")
	if err != nil {
		t.Fatalf("Error reading config: %v", err)
	}
	if configData.Database.Dsn != "postgres://user:example@localhost" {
		t.Fatal("dsn mismatch")
	}
	if configData.TestUsers.BoardAdminUsername != "test-user" {
		t.Fatal("Board Admin username mismatch")
	}
	if configData.TestUsers.BoardAdminPassword != "test-user-password" {
		t.Fatal("Board admin password mismatch")
	}
}
