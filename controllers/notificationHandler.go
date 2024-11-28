package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// FCMRequest 구조체 정의
type FCMRequest struct {
	Token string `json:"token" binding:"required"` // FCM 토큰
}

// FCM 푸시 알림 핸들러
func SendNotificationHandler(context *gin.Context) {
	var req FCMRequest

	// JSON 데이터 바인딩
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"errer": "올바른 요청 데이터를 제고해주세요"})
		return
	}

	// 최신 재난 메시지 가져오기
	latestMessage, err := service.FetchLatestDisasterMessage()
	if err != nil {
		log.Printf("재난 안전 문자 API 호출 실패: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "재난 안전 문자를 가져오는 데 실패했습니다.", "details": err.Error()})
		return
	}

	// FCM 푸시 알림 전송
	err = service.SendNotification(req.Token, "긴급 재난 알림", latestMessage)
	if err != nil {
		log.Printf("FCM 푸시 알림 전송 실패: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "푸시 알림 전송에 실패했습니다.", "details": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "알림 전송 성공"})
}
