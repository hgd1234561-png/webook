package web

import (
	"GkWeiBook/webook/internal/service"
	"GkWeiBook/webook/internal/service/oauth2/wechat"
	ijwt "GkWeiBook/webook/internal/web/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc         wechat.Service
	userService service.UserService
	ijwt.Handler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func (h *OAuth2WechatHandler) RegisterRouters(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.OAuth2URL)
	// 这边用Any万无一失
	g.Any("/callback", h.Callback)
}

func NewOAuth2WechatHandler(svc wechat.Service, userService service.UserService, cfg WechatHandlerConfig, jwtHdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:         svc,
		userService: userService,
		stateKey:    []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		cfg:         cfg,
		Handler:     jwtHdl,
	}
}

func (h *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造登录url失败",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间，预期用户扫码登录的时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.SetCookie("jwt-state", tokenStr, 60*5, "/oauth2/wechat/callback", "", h.cfg.Secure, true)
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	user, err := h.userService.FindOrCreateByWechat(ctx, info)

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.SetLoginToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		// 有人搞你
		// 做好监控
		return fmt.Errorf("拿不到 state 的 cookie: %w", err)
	}

	var sc StateClaims

	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("token已经 过期了: %w", err)
	}

	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	State string `json:"state"`
	jwt.RegisteredClaims
}
