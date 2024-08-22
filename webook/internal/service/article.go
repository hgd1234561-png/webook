package service

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/repository"
	"context"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, art domain.Article) (domain.Article, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	userRepo repository.UserRepository
}

func NewArticleService(repo repository.ArticleRepository, userRepo repository.UserRepository) ArticleService {
	return &articleService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (a *articleService) GetPubById(ctx context.Context, art domain.Article) (domain.Article, error) {
	author, err := a.userRepo.FindById(ctx, art.Authored.Id)
	if err != nil {
		return domain.Article{}, err
	}
	art.Authored.Name = author.Nickname
	return a.repo.GetPubById(ctx, art)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a *articleService) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.GetByAuthor(ctx, uid, offset, limit)
}

func (a *articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	return a.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
}

func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}

	return a.repo.Create(ctx, article)
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusPublished
	user, err := a.userRepo.FindById(ctx, article.Authored.Id)
	if err == nil {
		article.Authored = domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		}
	} else {
		// 记录日志
	}
	return a.repo.Sync(ctx, article)
}
