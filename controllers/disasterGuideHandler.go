package controllers

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ActionPlan의 ActRmks 목록 반환 핸들러
func GetActionPlanActRmks(c *gin.Context) {
	actRmksList, err := service.GetActionPlanActRmks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"actRmks": actRmksList})
}
