package db

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func DeleteAllDocument(c *gin.Context) {
	client := GetMongoClient()
	if client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not initialized"})
		return
	}

	collection := client.Database("SC2_DB").Collection("placeList")
	result, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete documents: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted %d documents", result.DeletedCount),
	})
}
