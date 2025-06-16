package user_service_test

import (
	"context"
	"log/slog"
	"testing"

	userService "github.com/HunDun0Ben/bs_server/app/service/user_service"
)

func TestFindUser(t *testing.T) {
	user, _ := userService.NewUserService().FindByLogin(context.Background(), "alice", "hashed_password_1")
	slog.Info("User found", "user", user)
}
