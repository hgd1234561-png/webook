package web

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/service"
	"GkWeiBook/webook/internal/web/jwt"
	"GkWeiBook/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	l   logger.Logger
	svc service.ArticleService
}

func NewArticleHandler(l logger.Logger, svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		l:   l,
		svc: svc,
	}
}

func (a *ArticleHandler) RegisterRouters(r *gin.Engine) {
	creatorGp := r.Group("/articles")
	creatorGp.POST("/edit", a.Edit)
	creatorGp.POST("/publish", a.Publish)
	creatorGp.POST("/withdraw", a.Withdraw)
	// 创作者接口
	creatorGp.GET("/detail/:id", a.Detail)
	creatorGp.POST("/list", a.List)

	pub := creatorGp.Group("/pub")
	pub.GET("/:id", a.PubDetail)

	// 传入一个参数，true 就是点赞, false 就是不点赞
	pub.POST("/like", a.Like)
	pub.POST("/collect", a.Collect)
}

// 编辑接口

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var editReq EditReq
	err := ctx.BindJSON(&editReq)
	if err != nil {
		a.l.Error("ArticleHandler  Edit  BindJSON  err", logger.Field{Key: "绑定参数错误", Val: err.Error()})
		return
	}

	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := a.svc.Save(ctx, domain.Article{
		Id:      editReq.Id,
		Title:   editReq.Title,
		Content: editReq.Content,
		Authored: domain.Author{
			Id: uc.Uid,
		},
	})

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Edit  Save  err", logger.Field{Key: "保存文章错误", Val: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	type PublishReq struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req PublishReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Publish  BindJSON  err", logger.Field{Key: "绑定参数错误", Val: err.Error()})
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	id, err := a.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Authored: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Publish  Publish  err", logger.Field{Key: "发布文章错误", Val: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type WithdrawReq struct {
		Id int64 `json:"id"`
	}
	var req WithdrawReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Withdraw  BindJSON  err", logger.Field{Key: "绑定参数错误", Val: err.Error()})
		return
	}

	uc := ctx.MustGet("user").(jwt.UserClaims)
	err := a.svc.Withdraw(ctx, uc.Uid, req.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Withdraw  Withdraw  err", logger.Field{Key: "撤回文章错误", Val: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (a *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.BindJSON(&page); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  List  Bind  err", logger.Field{Key: "绑定参数错误", Val: err.Error()})
		return
	}

	us := ctx.MustGet("user").(jwt.UserClaims)
	articles, err := a.svc.GetByAuthor(ctx, us.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  List  List  err", logger.Field{Key: "获取文章列表错误", Val: err}, logger.Field{Key: "用户id", Val: us.Uid},
			logger.Field{Key: "分页参数", Val: page})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: slice.Map[domain.Article, ArticleVo](articles, func(idx int, article domain.Article) ArticleVo {
			return ArticleVo{
				Id:       article.Id,
				Title:    article.Title,
				Abstract: article.Abstract(),
				AuthorId: article.Authored.Id,
				Status:   article.Status.ToUint8(),
				Ctime:    article.Ctime.Format(time.DateTime),
				Utime:    article.Utime.Format(time.DateTime),
			}
		}),
	})
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "id 错误",
		})
		a.l.Error("ArticleHandler  Detail  ParseInt  err", logger.Field{Key: "绑定参数错误", Val: err.Error()})
		return
	}
	art, err := a.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Detail  GetById  err", logger.Field{Key: "获取文章详情错误", Val: err.Error()})
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	if art.Authored.Id != uc.Uid {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		a.l.Error("ArticleHandler  Detail  art.Authored.Id != uc.Uid ", logger.Field{Key: "非法查询文章", Val: uc.Uid})
	}
	vo := ArticleVo{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Authored.Id,
		Status:   art.Status.ToUint8(),
		Ctime:    art.Ctime.Format(time.DateTime),
		Utime:    art.Utime.Format(time.DateTime),
	}
	ctx.JSON(http.StatusOK, Result{
		Data: vo,
	})
}

func (a *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "id 错误",
			Code: 4,
		})
		a.l.Error("ArticleHandler  PubDetail  ParseInt  err", logger.Field{Key: "id", Val: idStr}, logger.Field{Key: "绑定参数错误", Val: err.Error()})
		return
	}
	art, err := a.svc.GetPubById(ctx, domain.Article{Id: id})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		a.l.Error("ArticleHandler  PubDetail  GetPubById  err", logger.Field{Key: "获取文章详情错误", Val: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:         art.Id,
			Title:      art.Title,
			Content:    art.Content,
			AuthorId:   art.Authored.Id,
			AuthorName: art.Authored.Name,
			Status:     art.Status.ToUint8(),
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
		},
	})
}

func (a *ArticleHandler) Like(ctx *gin.Context) {

}

func (a *ArticleHandler) Collect(ctx *gin.Context) {

}

type ArticleVo struct {
	Id         int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	Content    string `json:"content,omitempty"`
	AuthorId   int64  `json:"authorId,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`

	ReadCnt    int64 `json:"readCnt"`
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	Liked      bool  `json:"liked"`
	Collected  bool  `json:"collected"`
}

type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
