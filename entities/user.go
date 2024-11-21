package entities

import (
	"delivery/pkg/utils"
	"errors"
	"time"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SendCodeReq struct {
	PhoneNumber string `json:"phone" validate:"required,phone"`
}

func (req *SendCodeReq) Validate() error {
	if !utils.IsPhoneValid(req.PhoneNumber) {
		return errors.New("invalid phone number: must be in format +99XXXXXXXXXX")
	}
	return nil
}

type RegistrReq struct {
	ID          string
	PhoneNumber string `json:"phone" validate:"required,phone" gorm:"type:varchar(13);not null;unique;"`
	Code        string `json:"code" validate:"required,len=6" gorm:"-"`
	FcmToken    string `json:"fcm_token" gorm:"fcm_token"`
}

func (req *RegistrReq) Validate() error {
	if !utils.IsPhoneValid(req.PhoneNumber) {
		return errors.New("invalid phone number: must be in format +99XXXXXXXXXX")
	}
	if len(req.Code) != 6 {
		return errors.New("invalid code: must be 6 characters long")
	}
	return nil
}

type RegistrRes struct {
	ID     string `json:"id"`
	Tokens Tokens `json:"tokens"`
}

type UserProfile struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	Firstname   string    `json:"firstname" gorm:"column:firstname"`
	Surname     string    `json:"surname" gorm:"column:surname"`
	Middlename  string    `json:"middlename" gorm:"column:middlename"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
	Birthdate   time.Time `json:"birthdate" gorm:"column:birthdate"`
	Gender      string    `json:"gender" gorm:"column:gender"`
	CreatedBy   string    `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   string    `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	//DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserLocation struct {
	ID        string  `gorm:"-"`
	UserID    string  `gorm:"column:user_id" json:"user_id"`
	Name      string  `gorm:"column:name" json:"name"`
	Latitude  float64 `gorm:"column:latitude" json:"latitude"`
	Longitude float64 `gorm:"column:longitude" json:"longitude"`
}

type Favorite struct {
	ID           string    `gorm:"-"`
	UserID       string    `json:"user_id" gorm:"column:user_id"`
	ProductID    string    `json:"product_id" gorm:"column:product_id"`
	Is_favorited bool      `gorm:"default:true"`
	AddedAt      time.Time `json:"added_at" gorm:"column:added_at"`
}
