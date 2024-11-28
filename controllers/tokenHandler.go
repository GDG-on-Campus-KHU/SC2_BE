package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TokenRequest struct {
	UserID int    `json:"user_id"binding:"required"`
	Token  string `json:"token"binding:"required"`
}

// RegisterTokenHandler : 클라이언트에서 토큰을 등록하는 엔드포인트
func RegisterTokenHandler(c *gin.Context) {
	var tokenRequest TokenRequest
	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repository.TokenRepository.SaveToken(tokenRequest.UserID, tokenRequest.Token)
	if err != nil {
		log.Printf("Failed to save token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
