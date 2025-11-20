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
	expire        int
	refreshExpire int

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
)

var u user.User

func init() {
	pflag.IntVarP(&expire, "expire", "t", 10, "access token 过期时间, 单位 s")
	pflag.IntVarP(&refreshExpire, "refresh token expire", "f", 10, "refresh token 过期时间, 单位 s")
	pflag.StringVarP(&u.Username, "username", "u", "username", "refresh token 过期时间, 单位 s")
	pflag.StringSliceVarP(&u.Roles, "roles", "r", []string{"admin", "user"}, "refresh token 过期时间, 单位 s")
}

func main() {
	pflag.Parse()
	slog.Info("Token expiration settings", slog.Int("access_token_expire", expire),
		slog.Int("refresh_token_expire", refreshExpire))
	slog.Info("User info for token generation", "user", u)

	conf.AppConfig.JWT.Expire = time.Duration(expire) * time.Second
	atok, rftok, _, _ := bsjwt.GenerateTokenPair(u)
	slog.Info("Generated tokens", slog.String("accessToken", atok),
		slog.String("refreshToken", rftok))
}
