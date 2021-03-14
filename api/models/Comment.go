package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	Content   string    `gorm:"text;not null" json:"content" `
	Author    User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Comment) Prepare() {
	c.ID = 0
	c.Content = html.EscapeString(strings.TrimSpace(c.Content))
	c.Author = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Comment) Validate() error {
	if c.Content == "" {
		return errors.New("Content Required")
	}

	return nil
}

func (c *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	var err error
	err = db.Debug().Model(&Comment{}).Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.AuthorID).Take(&c.Author).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) GetComments(db *gorm.DB, pid uint64) (*[]Comment, error) {
	var err error

	comments := []Comment{}
	err = db.Debug().Model(&Comment{}).Where("post_id = ?", pid).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}

	if len(comments) > 0 {
		for i, _ := range comments {
			err = db.Debug().Model(&User{}).Where("id = ?", comments[i].AuthorID).Take(&comments[i].Author).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}

	return &comments, nil
}

func (c *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {
	var err error

	err = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Updates(Comment{Content: c.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.AuthorID).Take(&c.Author).Error
		if err != nil {
			return &Comment{}, err
		}
	}

	return c, nil
}

func (c *Comment) DeleteAComment(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Comment not found")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint32) (int64, error) {
	comments := []Comment{}

	db = db.Debug().Model(&Like{}).Where("author_id = ?", uid).Find(&comments).Delete(&comments)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (l *Like) DeletePostComments(db *gorm.DB, pid uint64) (int64, error) {
	comments := []Comment{}

	db = db.Debug().Model(&Comment{}).Where("post_id = ?", pid).Find(&comments).Delete(&comments)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
