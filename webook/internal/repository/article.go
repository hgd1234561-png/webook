package repository

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/repository/cache"
	"GkWeiBook/webook/internal/repository/dao"
	"context"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type ArticleRepository interface {
	Update(ctx context.Context, article domain.Article) error
	Create(ctx context.Context, article domain.Article) (int64, error)
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, art domain.Article) (domain.Article, error)
}

type articleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func NewArticleRepository(dao dao.ArticleDAO, c cache.ArticleCache) ArticleRepository {
	return &articleRepository{
		dao:   dao,
		cache: c,
	}
}

func (a *articleRepository) GetPubById(ctx context.Context, art domain.Article) (domain.Article, error) {
	res, err := a.cache.GetPub(ctx, art.Id)
	if err == nil {
		return res, nil
	}
	artNew, err := a.dao.GetPubById(ctx, art.Id)
	if err != nil {
		return domain.Article{}, err
	}
	res = a.toDomain(dao.Article(artNew))
	res.Authored.Name = art.Authored.Name
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := a.cache.SetPub(ctx, res)
		if er != nil {
			// 缓存失败，不panic
			// 记录日志
		}
	}()
	return res, nil
}

func (a *articleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	if art, err := a.cache.Get(ctx, id); err == nil {
		return art, nil
	}
	art, err := a.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	res := a.toDomain(art)
	go func() {
		err := a.cache.Set(ctx, res)
		if err != nil {
			// 缓存失败，不panic
		}
	}()
	return res, nil
}

func (a *articleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 首先第一步，判定要不要查询缓存
	// 事实上， limit <= 100 都可以查询缓存
	if offset == 0 && limit == 100 {
		// 查缓存
		arts, err := a.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return arts, nil
		} else {
			// 考虑记日志， 缓存未命中，可以忽略
		}
	}

	arts, err := a.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}

	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, item dao.Article) domain.Article {
		return a.toDomain(item)
	})
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if offset == 0 && limit == 100 {
			// 缓存回写失败，不一定是大问题，但有可能是大问题
			err = a.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 记录日志
				// 我需要监控这里
			}
		}
	}()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		a.preCache(ctx, res)
	}()
	return res, nil
}

func (a *articleRepository) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	err := a.dao.SyncStatus(ctx, uid, id, status.ToUint8())
	if err == nil {
		er := a.cache.DelFirstPage(ctx, uid)
		if er != nil {
			// 记录日志
		}
	}
	return err
}

func (a *articleRepository) Update(ctx context.Context, article domain.Article) error {
	err := a.dao.UpdateById(ctx, a.toEntity(article))
	if err == nil {
		er := a.cache.DelFirstPage(ctx, article.Authored.Id)
		if er != nil {
			// 记录日志
		}
	}
	return nil
}

func (a *articleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	id, err := a.dao.Insert(ctx, a.toEntity(article))
	if err == nil {
		er := a.cache.DelFirstPage(ctx, article.Authored.Id)
		if er != nil {
			// 记录日志
		}
	}
	return id, err
}

func (a *articleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	id, err := a.dao.Sync(ctx, a.toEntity(article))
	if err == nil {
		er := a.cache.DelFirstPage(ctx, article.Authored.Id)
		if er != nil {
			// 记录日志
		}
	}
	// 在这里尝试，设置缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 你可以灵活设置过期时间

		if article.Authored.Name == "" {
			// 要记录日志
			return
		}
		er := a.cache.SetPub(ctx, article)
		if er != nil {
			// 记录日志
		}
	}()

	return id, err
}

func (a *articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Authored.Id,
		Status:   art.Status.ToUint8(),
	}
}

func (a *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Authored: domain.Author{
			// 这里有一个错误
			Id: art.AuthorId,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
}

func (a *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	// 预缓存，列表的第一篇文章被查看的几率比较大，所以预先缓存
	const size = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) < size {
		err := a.cache.Set(ctx, arts[0])
		if err != nil {
			// 记录缓存
		}
	}
}
