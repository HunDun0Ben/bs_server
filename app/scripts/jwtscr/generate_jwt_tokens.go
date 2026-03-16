package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

var (
	accessTokenExpireSeconds  int
	refreshTokenExpireSeconds int

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
)

var u user.User

func init() {
	pflag.IntVarP(&accessTokenExpireSeconds, "expire", "t", 10, "access token 过期时间, 单位 s")
	pflag.IntVarP(&refreshTokenExpireSeconds, "refresh-token-expire", "f", 10, "refresh token 过期时间, 单位 s")
	pflag.StringVarP(&u.Username, "username", "u", "username", "token 关联的用户名")
	pflag.StringSliceVarP(&u.Roles, "roles", "r", []string{"admin", "user"}, "用户角色 (逗号分隔)")
}

func main() {
	pflag.Parse()
	if err := conf.InitConfig(); err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}
	slog.Info("Token expiration settings", slog.Int("access_token_expire", accessTokenExpireSeconds),
		slog.Int("refresh_token_expire", refreshTokenExpireSeconds))
	slog.Info("User info for token generation", "user", u)

	conf.AppConfig.JWT.Expire = time.Duration(accessTokenExpireSeconds) * time.Second
	atok, rftok, err, rfErr := bsjwt.GenerateTokenPair(u, false, nil)
	if err != nil || rfErr != nil {
		slog.Error("Failed to generate token pair", "access_token_error", err, "refresh_token_error", rfErr)
		os.Exit(1)
	}
	slog.Info("Generated tokens", slog.String("accessToken", atok),
		slog.String("refreshToken", rftok))
}
