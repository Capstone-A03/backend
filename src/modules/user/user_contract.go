package user

import (
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/google/uuid"
)

type getUserListQuery struct {
	SearchByRole *uc.Role   `query:"role"`
	LastID       *uuid.UUID `query:"lastId"`
	Limit        *int       `query:"limit"`
}

type getUserReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateUserReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateUserReq struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type deleteUserReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
