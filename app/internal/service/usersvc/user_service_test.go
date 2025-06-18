package usersvc_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
)

func TestFindUser(t *testing.T) {
	user, _ := usersvc.NewUserService().FindByLogin(context.Background(), "alice", "hashed_password_1")
	slog.Info("User found", "user", user)
}
