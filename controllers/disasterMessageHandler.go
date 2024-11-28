package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 재난 안전 문자 api를 통해 응답을 받음
func GetDisasterMessagesHandler(context *gin.Context) {
	// 클라이언트가 /disaster-messages 경로를 호출하면 FetchLatestDisasterMessage 를 실행
	message, err := service.FetchLatestDisasterMessage()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to fetch disaster messages",
			"detail": err.Error(),
		})
		return
	}

	if message == nil {
		context.JSON(http.StatusOK, gin.H{
			"message": "No new disaster messages available",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data": message,
	})
}

// AI 모델에 재난 알림 전송
func SendDisasterMessageController(context *gin.Context) {
	var request models.DisasterMessage

	// 요청 데이터 바인딩
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data", "detail": err.Error(),
		})
		return
	}

	// AI 모델에 데이터 전달
	response, err := service.SendDisasterMessage(request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to send disaster message",
			"detail": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, response)

	HandlePushNotification(response, token)
}
