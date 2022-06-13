package models

import (
	"errors"
	_ "log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
)

type InviteToken struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Token     string    `gorm:"size:255;not null;unique" json:"invite_token"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

var letters = []rune("1a2b3c4d5e6f7g8h9i0jklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ=&%$#@#&^*+=")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (it *InviteToken) SaveInvite(db *gorm.DB) (*InviteToken, error) {

	var err error
	rand.Seed(time.Now().UnixNano())
	inviteToken := InviteToken{
		Token: randSeq(20),
	}

	err = db.Debug().Create(&inviteToken).Error
	if err != nil {
		return &InviteToken{}, err
	}
	return &inviteToken, nil
}

func (u *InviteToken) FindInviteTokenByToken(db *gorm.DB, token string) (*InviteToken, error) {
	var err error
	err = db.Debug().Model(InviteToken{}).Where("token = ? and (created_at between CURDATE() AND CURDATE() + INTERVAL 1 DAY)", token).Take(&u).Error
	if err != nil {
		return &InviteToken{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &InviteToken{}, errors.New("token not Not Found")
	}
	return u, err
}
