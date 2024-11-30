package service

import (
	"context"
	"errors"
	"github.com/GDG-on-Campus-KHU/SC2_BE/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var ActionCollection *mongo.Collection

func InitActionCollection(client *mongo.Client) {
	tokenCollection = client.Database("SC2_DB").Collection("tokens")
	log.Println("[INFO] Token collection initialized")
}

// 데이터를 MongoDB에 저장
func SaveDisasterGuideResponse(response *models.DisasterGuideResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ActionCollection.InsertOne(ctx, response)
	if err != nil {
		log.Printf("[ERROR] Failed to save disaster response: %v", err)
		return err
	}
	return nil
}

// ActionPlan의 ActRmks 목록을 반환
func GetActionPlanActRmks() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var response models.DisasterGuideResponse
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	err := ActionCollection.FindOne(ctx, bson.M{}, opts).Decode(&response)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("no data found")
		}
		log.Printf("[ERROR] Failed to fetch latest disaster response: %v", err)
		return nil, err
	}

	var actRmksList []string
	for _, plan := range response.Results.HotspotResults.ActionPlan {
		actRmksList = append(actRmksList, plan.ActRmks)
	}

	return actRmksList, nil
}
