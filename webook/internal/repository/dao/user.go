package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Insert(ctx context.Context, u User) error
	Update(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}

	return err
}

func (dao *GORMUserDAO) Update(ctx context.Context, u User) error {
	result := dao.db.WithContext(ctx).Model(&u).Updates(User{
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Utime:    time.Now().UnixMilli(),
	})

	return result.Error
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

// User 直接对应数据库的表结构
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	// 代表这是一个可以为 NULL 的列
	Phone sql.NullString `gorm:"unique"`

	// 创建时间  毫秒数
	Ctime int64
	// 更新时间
	Utime int64
}
