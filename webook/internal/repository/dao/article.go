package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	UpdateById(ctx context.Context, entity Article) error
	Insert(ctx context.Context, entity Article) (int64, error)
	Sync(ctx context.Context, entity Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (g *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var res PublishedArticle
	err := g.db.WithContext(ctx).
		Where("id = ?", id).
		First(&res).Error
	return res, err
}

func (g *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := g.db.WithContext(ctx).
		Where("id = ?", id).
		First(&art).Error
	return art, err
}

func (g *GORMArticleDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	var articles []Article
	err := g.db.WithContext(ctx).
		Where("author_id = ?", uid).
		Offset(offset).
		Limit(limit).
		Order("utime desc").
		Find(&articles).Error

	return articles, err
}

func (g *GORMArticleDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? and author_id = ?", id, uid).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("ID 不对或者创作者不对")
		}
		return tx.Model(&PublishedArticle{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			}).Error
	})
}

func (g *GORMArticleDAO) Sync(ctx context.Context, entity Article) (int64, error) {
	var id = entity.Id
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)
		// 事务独占一个连接
		dao := NewArticleDao(tx)
		if id > 0 {
			err = dao.UpdateById(ctx, entity)
		} else {
			id, err = dao.Insert(ctx, entity)
		}
		if err != nil {
			return err
		}
		entity.Id = id
		now := time.Now().UnixMilli()
		pubArt := PublishedArticle(entity)
		pubArt.Ctime = now
		pubArt.Utime = now
		err = tx.Clauses(clause.OnConflict{
			// 对MySQL不起效，但是可以兼容别的方言
			// INSERT xxx ON DUPLICATE KEY SET `title`=?
			// 别的方言：
			// sqlite INSERT XXX ON CONFLICT DO UPDATES WHERE
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   pubArt.Title,
				"content": pubArt.Content,
				"utime":   now,
				"status":  pubArt.Status,
			}),
		}).Create(&pubArt).Error
		return err
	})
	return id, err
}

func (g *GORMArticleDAO) UpdateById(ctx context.Context, entity Article) error {
	now := time.Now().UnixMilli()
	res := g.db.WithContext(ctx).Model(&entity).Where("id = ? AND author_id = ?", entity.Id, entity.AuthorId).Updates(
		map[string]any{
			"title":   entity.Title,
			"content": entity.Content,
			"status":  entity.Status,
			"utime":   now,
		})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("ID 不对或者创作者不对")
	}

	return nil
}

func (g *GORMArticleDAO) Insert(ctx context.Context, entity Article) (int64, error) {
	now := time.Now().UnixMilli()
	entity.Ctime = now
	entity.Utime = now
	err := g.db.WithContext(ctx).Create(&entity).Error
	return entity.Id, err
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 我要根据创作者ID来查询
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	// 更新时间
	Utime int64 `bson:"utime,omitempty"`
}

type PublishedArticle Article
