package pushnotification

type getPushNotificationWsReqQuery struct {
	token string `query:"token" validate:"required"`
}

type getPushNotificationWsRes struct {
	Kind    int         `json:"kind"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
