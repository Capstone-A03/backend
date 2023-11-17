package contracts

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Error      *Error      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Pagination struct {
	Limit *int `json:"limit,omitempty"`
	Count *int `json:"count,omitempty"`
	Page  *int `json:"page,omitempty"`
	Total *int `json:"total,omitempty"`
}

func NewError(err *fiber.Error, message string) Error {
	return Error{
		Code:    err.Code,
		Status:  err.Error(),
		Message: message,
	}
}

func NewWSRes[T any](c *websocket.Conn, data T) error {
	bytes, err := sonic.ConfigFastest.Marshal(data)
	if err != nil {
		return err
	}
	if err := c.WriteMessage(websocket.TextMessage, bytes); err != nil {
		return err
	}
	return nil
}

func (e Error) Error() string {
	return e.Message
}
