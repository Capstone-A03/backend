package pushnotification

type getPushNotificationWsRes struct {
	Kind    int    `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
