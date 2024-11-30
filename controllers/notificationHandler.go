package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// 토큰 등록
func RegisterToken(c *gin.Context) {
	var req models.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.SaveToken(models.TokenRequest{Token: os.Getenv("FCM_TOKEN")})
	c.JSON(http.StatusOK, gin.H{"message": "토큰 등록 완료"})
}

// HandlePushNotification: 푸시 알림을 처리하는 함수
func HandlePushNotification(response *models.DisasterGuideResponse) {
	if response == nil {
		log.Println("No response to process for push notification")
		return
	}

	// 푸시 알림 전송
	err := service.SendNotification(response.Results.HotspotResults.PushAlarming)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
	} else {
		log.Println("Push notification sent successfully")
	}
}
