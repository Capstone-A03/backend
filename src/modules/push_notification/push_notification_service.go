package pushnotification

var sendCh chan getPushNotificationWsRes

func Send(kind int, title string, message string, data interface{}) {
	sendCh <- getPushNotificationWsRes{
		Kind:    kind,
		Title:   title,
		Message: message,
		Data:    data,
	}
}
