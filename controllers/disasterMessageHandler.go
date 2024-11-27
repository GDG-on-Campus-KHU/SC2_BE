package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 클라이언트 요청을 처리하기 위한 Handler 함수
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
