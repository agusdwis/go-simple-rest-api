package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Like struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (l *Like) SaveLike(db *gorm.DB) (*Like, error) {

	err := db.Debug().Model(&Like{}).Where("post_id = ? AND author_id = ?", l.PostID, l.AuthorID).Take(&l).Error

	if err != nil {
		if err.Error() == "record not found" {
			err = db.Debug().Model(&Like{}).Create(&l).Error
			if err != nil {
				return &Like{}, nil
			}
		}
	} else {
		err = errors.New("double like")
		return &Like{}, err
	}

	return l, nil
}

func (l *Like) DeleteLike(db *gorm.DB) (*Like, error) {
	var deletedLike *Like

	err := db.Debug().Model(Like{}).Where("id = ?", l.ID).Take(&l).Error

	if err != nil {
		return &Like{}, err
	} else {
		deletedLike = l
		db = db.Debug().Model(&Like{}).Where("id = ?", l.ID).Take(&l).Delete(&Like{})

		if db.Error != nil {
			fmt.Println("cant delete like", db.Error)
			return &Like{}, db.Error
		}
	}

	return deletedLike, nil
}

func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (*[]Like, error) {
	likes := []Like{}

	err := db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Error

	if err != nil {
		return &[]Like{}, nil
	}

	return &likes, err
}

func (l *Like) DeleteUserLikes(db *gorm.DB, uid uint32) (int64, error) {
	likes := []Like{}

	db = db.Debug().Model(&Like{}).Where("author_id = ?", uid).Find(&likes).Delete(&likes)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (l *Like) DeletePostLikes(db *gorm.DB, pid uint64) (int64, error) {
	likes := []Like{}

	db = db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Delete(&likes)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}