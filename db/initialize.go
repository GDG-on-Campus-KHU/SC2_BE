package db

import (
	"net/http"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

func DeleteAllDocument(c *gin.Context) {
    if GetMongoClient() == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not initialized"})
        return
    }

    collection := GetMongoClient().Database("SC2_DB").Collection("placeList")

    result, err := collection.DeleteMany(context.TODO(), struct{}{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete documents: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("Successfully deleted %d documents", result.DeletedCount),
    })
}