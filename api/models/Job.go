package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Job struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title       string    `gorm:"size:255;not null;unique" json:"title"`
	Description string    `gorm:"size:255;not null;" json:"Description"`
	IsActive    bool      `json:"is_active"`
	Rate        int       `json:"rate"`
	User        User      `json:"user"`
	UserID      uint32    `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Distance    float64   `sql:"_" json:"distance"`
}

func (j *Job) Prepare() {
	j.ID = 0
	j.Title = html.EscapeString(strings.TrimSpace(j.Title))
	j.Description = html.EscapeString(strings.TrimSpace(j.Description))
	j.User = User{}
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
}

func (j *Job) Validate() error {

	if j.Title == "" {
		return errors.New("Required Title")
	}
	if j.Description == "" {
		return errors.New("Required Description")
	}
	if j.UserID < 1 {
		return errors.New("Required User")
	}
	return nil
}

func (j *Job) SaveJob(db *gorm.DB) (*Job, error) {
	var err error
	err = db.Debug().Model(&Job{}).Create(&j).Error
	if err != nil {
		return &Job{}, err
	}
	if j.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", j.UserID).Take(&j.User).Error
		if err != nil {
			return &Job{}, err
		}
	}
	return j, nil
}

func (j *Job) FindAllJobs(db *gorm.DB) (*[]Job, error) {
	var err error
	Jobs := []Job{}
	err = db.Debug().Model(&Job{}).Limit(100).Find(&Jobs).Error
	if err != nil {
		return &[]Job{}, err
	}
	if len(Jobs) > 0 {
		for i, _ := range Jobs {
			err := db.Debug().Model(&User{}).Where("id = ?", Jobs[i].UserID).Take(&Jobs[i].User).Error
			if err != nil {
				return &[]Job{}, err
			}
		}
	}
	return &Jobs, nil
}

func (j *Job) FindJobByID(db *gorm.DB, pid uint64) (*Job, error) {
	var err error
	err = db.Debug().Model(&Job{}).Where("id = ?", pid).Take(&j).Error
	if err != nil {
		return &Job{}, err
	}
	if j.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", j.UserID).Take(&j.User).Error
		if err != nil {
			return &Job{}, err
		}
	}
	return j, nil
}

func (j *Job) UpdateAJob(db *gorm.DB) (*Job, error) {

	var err error

	err = db.Debug().Model(&Job{}).Where("id = ?", j.ID).Updates(Job{Title: j.Title, Description: j.Description, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Job{}, err
	}
	if j.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", j.UserID).Take(&j.User).Error
		if err != nil {
			return &Job{}, err
		}
	}
	return j, nil
}

func (j *Job) DeleteAJob(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Job{}).Where("id = ? and user_id = ?", pid, uid).Take(&Job{}).Delete(&Job{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Job not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (j *Job) SoftDelete(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	err := db.Debug().Model(&Job{}).Where("id = ?", j.ID).Updates(Job{IsActive: false, UpdatedAt: time.Now()}).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Job not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (j *Job) FindJobsByDistance(db *gorm.DB, distance float64, lat float64, lon float64) (*[]Job, error) {
	fmt.Println("func params: ")
	fmt.Print(distance)
	fmt.Print(lat)
	fmt.Println(lon)
	var err error
	Jobs := []Job{}
	err = db.Raw("SELECT *,( 6371 * acos( cos( radians(?) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(?) ) + sin( radians(?) ) * sin( radians( latitude ) ) ) ) AS distance FROM jobs where user_id=? HAVING distance < ? ", lat, lon, lat, j.UserID, distance).Scan(&Jobs).Error
	if err != nil {
		return &[]Job{}, err
	}
	if len(Jobs) > 0 {
		for i, _ := range Jobs {
			err := db.Debug().Model(&User{}).Where("id = ?", Jobs[i].UserID).Take(&Jobs[i].User).Error
			if err != nil {
				return &[]Job{}, err
			}
		}
	}

	return &Jobs, nil
}
