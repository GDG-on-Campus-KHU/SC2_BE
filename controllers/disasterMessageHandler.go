package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDisasterMessagesHandler(context *gin.Context) {
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
