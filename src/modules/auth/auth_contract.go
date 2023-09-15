package auth

import (
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/google/uuid"
)

type signupReq struct {
	Name     *string  `json:"name" validate:"required,gt=0"`
	Username *string  `json:"username" validate:"required,gt=0"`
	Password *string  `json:"password" validate:"required,gt=0"`
	Role     *uc.Role `json:"role" validate:"required,role"`
}

type signinReq struct {
	Username *string `json:"username" validate:"required,gt=0"`
	Password *string `json:"password" validate:"required,gt=0"`
}

type userRes struct {
	Token *string    `json:"token"`
	ID    *uuid.UUID `json:"id"`
	Name  *string    `json:"name"`
	Role  *uc.Role   `json:"role"`
}
