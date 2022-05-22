package seed

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/naqash/goBlog/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "mns",
		Email:    "naqash.subtain@gmail.com",
		Password: "mns",
	},
	models.User{
		Nickname: "naqash",
		Email:    "naqash.subtainr@hotmail.com",
		Password: "mns",
	},
}

var Jobs = []models.Job{
	models.Job{

		Title:       "First Job",
		Description: "Descp 1",
		IsActive:    true,
		Rate:        0,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
	models.Job{

		Title:       "Second Job",
		Description: "",
		IsActive:    true,
		Rate:        4,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
	models.Job{

		Title:       "Third Job",
		Description: "",
		IsActive:    true,
		Rate:        4,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
	models.Job{

		Title:       "Fourth Job",
		Description: "",
		IsActive:    true,
		Rate:        4,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
	models.Job{

		Title:       "Fivth Job",
		Description: "",
		IsActive:    true,
		Rate:        4,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
	models.Job{

		Title:       "Sixth Job",
		Description: "",
		IsActive:    true,
		Rate:        4,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
	models.Job{

		Title:       "8th Job",
		Description: "",
		IsActive:    true,
		Rate:        4,
		User:        models.User{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Latitude:    0,
		Longitude:   0,
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Job{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Job{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Job{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		Jobs[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Job{}).Create(&Jobs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed Jobs table: %v", err)
		}
	}
}
