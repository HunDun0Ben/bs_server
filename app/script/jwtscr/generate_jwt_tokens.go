package main

import (
	"log/slog"
	"time"

	"github.com/spf13/pflag"

	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

var (
	expire    int
	refExpire int

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
)

var u user.User

func init() {
	pflag.IntVarP(&expire, "expire", "t", 10, "access token 过期时间, 单位 s")
	pflag.IntVarP(&refExpire, "refresh token expire", "f", 10, "refresh token 过期时间, 单位 s")
	pflag.StringVarP(&u.Username, "username", "u", "username", "refresh token 过期时间, 单位 s")
	pflag.StringSliceVarP(&u.Roles, "roles", "r", []string{"admin", "user"}, "refresh token 过期时间, 单位 s")
}

func main() {
	pflag.Parse()
	slog.Info("", slog.Int("expire", expire))
	slog.Info("", slog.Int("refresh token expire", expire))
	slog.Info("user info", "user", u)

	conf.AppConfig.JWT.Expire = time.Duration(expire * time.Now().Second())
	atok, rftok, _, _ := bsjwt.GenerateTokenPair(u)
	slog.Info("access token:", slog.String("", atok))
	slog.Info("access token:", slog.String("", rftok))
}
