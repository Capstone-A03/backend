package user

import "github.com/google/uuid"

type getUserReqQuery struct {
	ID *uuid.UUID `query:"id" validate:"required"`
}

type updateUserReq struct {
	ID       *uuid.UUID `json:"id" validate:"required"`
	Name     *string    `json:"name"`
	Username *string    `json:"username"`
	Password *string    `json:"password"`
}

type deleteUserReq struct {
	ID *uuid.UUID `json:"id" validate:"required"`
}
