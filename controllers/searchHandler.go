package controllers

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/GDG-on-Campus-KHU/SC2_BE/models"
    "github.com/GDG-on-Campus-KHU/SC2_BE/service"
)

func NaverSearchHandler(c *gin.Context) {
    query := c.Query("query")
    display := 5

    // 네이버 API 검색 수행
    result, err := service.NaverSearch(query, display)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // MongoDB에 결과 저장 (HTML 태그가 제거된 상태로 저장됨)
    err = service.SaveSearchResults(result)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save to database: %v", err)})
        return
    }

    // 응답에서도 HTML 태그 제거
    for i := range result.Items {
        result.Items[i].Title = service.RemoveHTMLTags(result.Items[i].Title)
    }

    // 간소화된 응답 생성
    simplifiedResults := make([]models.SimplifiedResponse, len(result.Items))
    for i, item := range result.Items {
        simplifiedResults[i] = models.SimplifiedResponse{
            Title:       service.RemoveHTMLTags(item.Title),
            RoadAddress: item.RoadAddress,
        }
    }

    c.JSON(http.StatusOK, gin.H{"items": simplifiedResults})
}