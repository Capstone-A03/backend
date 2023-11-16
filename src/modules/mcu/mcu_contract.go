package mcu

import "github.com/google/uuid"

type getMcuListReqQuery struct {
	LastID *uuid.UUID `query:"lastId"`
	Limit  *int       `query:"limit"`
}

type getMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addMcuReq struct {
	DumpID *uuid.UUID `json:"dumpId" validate:"required"`
}

type updateMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateMcuReq struct {
	DumpID *uuid.UUID `json:"dumpId"`
}

type deleteMcuReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}
