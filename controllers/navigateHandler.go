package controllers

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/GDG-on-Campus-KHU/SC2_BE/models"
    "github.com/GDG-on-Campus-KHU/SC2_BE/db"
    "github.com/GDG-on-Campus-KHU/SC2_BE/service"
    "go.mongodb.org/mongo-driver/mongo"
)

func GetNavigateHandler(c *gin.Context) {
    if db.GetMongoClient() == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not initialized"})
        return
    }

    clientID := os.Getenv("Naver_Cloud_Client_ID")
    clientSecret := os.Getenv("Naver_Cloud_Client_Secret")

    navService := service.NewNavigateClient(clientID, clientSecret)

    start := c.Query("start")
    locationName := c.Query("location")

    if start == "" || locationName == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "start coordinates and location name are required",
        })
        return
    }

    collection := db.GetMongoClient().Database("SC2_DB").Collection("placeList")
    
    var goalPlace models.PlaceLocation
    err := collection.FindOne(context.TODO(), gin.H{"title": locationName}).Decode(&goalPlace)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{
                "error": fmt.Sprintf("Location '%s' not found in database", locationName),
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    goalMapx := service.FormatCoordinate(goalPlace.Mapx)
    goalMapy := service.FormatCoordinate(goalPlace.Mapy)
    goal := fmt.Sprintf("%s,%s", goalMapx, goalMapy)

    response, err := navService.GetNavigate(start, goal)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if len(response.Route.Traoptimal) > 0 {
        routeResponse := models.RouteResponse{
            Start: response.Route.Traoptimal[0].Summary.Start.Location,
            Goal:  response.Route.Traoptimal[0].Summary.Goal.Location,
            Path:  response.Route.Traoptimal[0].Path,
        }
        c.JSON(http.StatusOK, routeResponse)
    } else {
        c.JSON(http.StatusNotFound, gin.H{"error": "no route found"})
    }
}