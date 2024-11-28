package service

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/GDG-on-Campus-KHU/SC2_BE/config"
	"log"
)

// FCM 메시지 전송
func SendNotification(body string) error {
	client, err := config.FirebaseApp.Messaging(context.Background())
	if err != nil {
		log.Printf("FCM 클라이언트 초기화 실패: %v", err)
		return err
	}

	// 메시지 구성
	fcmMessage := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: "재난 안전",
			Body:  body,
		},
	}

	// 메시지 전송
	_, err = client.Send(context.Background(), fcmMessage)
	if err != nil {
		log.Printf("FCM 메시지 전송 실패 %v", err)
		return err
	}

	log.Printf("FCM 메시지 전송 성공")
	return nil
}
