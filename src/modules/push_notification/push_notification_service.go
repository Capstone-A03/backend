package pushnotification

var sendCh chan getPushNotificationWsRes

func Send(kind int, title string, message string) {
	sendCh <- getPushNotificationWsRes{
		Kind:    kind,
		Title:   title,
		Message: message,
	}
}
