package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"log"
)

// HandlePushNotification: 푸시 알림을 처리하는 함수
func HandlePushNotification(response *models.DisasterGuideResponse) {
	if response == nil {
		log.Println("No response to process for push notification")
		return
	}

	// 알림 제목과 내용 설정
	title := "재난 경보"
	body := response.PushAlarming

	// 푸시 알림 전송
	err := service.SendNotification(token, title, body)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
	} else {
		log.Println("Push notification sent successfully")
	}
}

func SendDisasterNotification(disasterData models.DisasterMessage) {
	// 1. 재난 데이터 AI 모델로 전송
	response, err := service.SendDisasterMessage(disasterData)
	if err != nil {
		log.Printf("Failed to send disaster notification to AI model: %v", err)
		return
	}

	// 2. AI 모델 응답 데이터를 FCM 푸시 알림으로 전송
	HandlePushNotification(response, token)
}
