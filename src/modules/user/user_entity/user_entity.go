package userentity

import (
	"capstonea03/be/src/libs/db/sql"
	"capstonea03/be/src/libs/env"
	"capstonea03/be/src/libs/hash/argon2"
	applogger "capstonea03/be/src/libs/logger"
	"capstonea03/be/src/libs/validator"
	u "capstonea03/be/src/modules/user/user_constant"
	"reflect"
)

type UserModel struct {
	sql.Model
	Name     *string `gorm:"not null" json:"name,omitempty"`
	Username *string `gorm:"uniqueIndex;not null" json:"username,omitempty"`
	Password *string `gorm:"not null" json:"-"`
	Role     *u.Role `gorm:"not null" json:"role,omitempty" validate:"role"`
}

func (UserModel) TableName() string {
	return "users"
}

type userDB = sql.Service[UserModel]

var userRepo *userDB
var logger = applogger.New("UserModule")

func InitRepository(db *sql.DB) {
	if db == nil {
		logger.Panic("db cannot be nil")
	}

	userRepo = sql.NewService[UserModel](db)
}

func UserRepository() *userDB {
	if userRepo == nil {
		logger.Panic("userRepo is nil")
	}

	return userRepo
}

func CreateInitialUser() {
	role := u.Role(env.Get(env.INITIAL_USER_ROLE, env.Option{MustExist: true}))

	count, err := UserRepository().Count(&sql.CountOptions{
		Where: &[]sql.Where{
			{
				Query: "role = ?",
				Args:  []interface{}{role},
			},
		},
	})
	if err != nil {
		logger.Panic(err)
	}
	if *count > 0 {
		return
	}

	name := env.Get(env.INITIAL_USER_NAME, env.Option{MustExist: true})
	username := env.Get(env.INITIAL_USER_USERNAME, env.Option{MustExist: true})
	password := env.Get(env.INITIAL_USER_PASSWORD, env.Option{MustExist: true})
	encodedHash, err := argon2.GetEncodedHash(&password)
	if err != nil {
		logger.Panic(err)
	}
	if _, err := UserRepository().Create(&UserModel{
		Name:     &name,
		Username: &username,
		Password: encodedHash,
		Role:     &role,
	}); err != nil {
		logger.Panic(err)
	}
}

func RegisterRoleValidation() {
	if err := validator.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() == reflect.Pointer {
			if fl.Field().IsNil() {
				return true
			}
		}
		if fl.Field().String() == "ADMIN" || fl.Field().String() == "JANITOR" {
			return true
		}
		return false
	}, true); err != nil {
		logger.Panic(err)
	}
}
