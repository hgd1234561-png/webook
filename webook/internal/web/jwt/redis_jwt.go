package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	// access_token key

	AtKey = []byte("bstmmTdM2KFxXcm544kMZzzBsBgwgb6J")
	// refresh_token key

	RtKey = []byte("bstmmTdM2KFxXcm544kMZzzBsBgwgb6B")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: client,
	}
}

func (r RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc := ctx.MustGet("user").(UserClaims)

	return r.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", uc.Ssid), "", time.Hour*24*7).Err()
}

func (r RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	// 用JWT来校验
	tokenHeader := ctx.GetHeader("Authorization")

	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 || segs[0] != "Bearer" {
		return ""
	}
	return segs[1]
}

func (r RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, uid, ssid)
	return err
}

// 没用指针是因为要组合

func (r RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 1 分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := r.cmd.Exists(ctx, fmt.Sprintf("user:ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token 无效")
	}
	return nil
}

func (r RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}
